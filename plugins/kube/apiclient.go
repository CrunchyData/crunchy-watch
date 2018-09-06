package main

import (
	"bytes"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crunchydata/crunchy-watch/util"
)

func deletePrimaryDeployment(namespace string, name string) error {
	delOptions := metav1.DeleteOptions{}
	var delProp metav1.DeletionPropagation
	delProp = metav1.DeletePropagationForeground
	delOptions.PropagationPolicy = &delProp

	err := client.ExtensionsV1beta1().Deployments(namespace).Delete(name, &delOptions)

	return err
}

func deletePrimaryPod(namespace string, name string) error {
	podsClient := client.CoreV1().Pods(namespace)

	err := podsClient.DeleteCollection(
		&metav1.DeleteOptions{},
		metav1.ListOptions{
			LabelSelector: fmt.Sprintf("name=%s", name),
		},
	)

	return err
}

/*
get replica by namespace and name
*/
func promoteReplica(namespace string, name string) error {

	pod, err := client.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})

	if err != nil {
		return fmt.Errorf("could not get pod info: %v", err)
	}

	if len(pod.Spec.Containers) != 1 {
		return fmt.Errorf("could not determine which container to use")
	}

	var stderr string

	cmd := []string{fmt.Sprintf("/opt/cpm/bin/promote.sh")}

	log.Debugf("executing cmd: %s on pod %s in namespace %s container: %s", cmd, pod.Name, pod.Namespace, pod.Spec.Containers[0].Name)

	stdout, stderr, err := util.ExecWithOptions(restConfig, *client, util.ExecOptions{
		Command:       cmd,
		Namespace:     pod.Namespace,
		PodName:       pod.Name,
		ContainerName: pod.Spec.Containers[0].Name,

		Stdin:              nil,
		CaptureStdout:      true,
		CaptureStderr:      true,
		PreserveWhitespace: false,
	})

	if err != nil {
		log.Errorf("Error executing cmd: %s stderr: %s stdout: %s", cmd, stderr, stdout)
		return fmt.Errorf("could not execute: %v", err)
	}

	return err
}

func drainDeployment(namespace, primaryName, replicaName string) error {
	origDeployment, err := client.ExtensionsV1beta1().Deployments(namespace).Get(replicaName, metav1.GetOptions{})

	if err != nil {
		log.Error(err)
		log.Error("error getting Deployment " + replicaName)
		return err
	}

	var val int32
	val = 0
	origDeployment.Spec.Replicas = &val
	//also update or relabel the deployment so new pods will inherit
	//the correct name label, that of the primary, this allows
	//the primary service to match against the new pod
	_, err = client.ExtensionsV1beta1().Deployments(namespace).Update(origDeployment)

	if err != nil {
		log.Error(err)
		log.Error("error updating Deployment " + replicaName)
		return err
	}
	log.Info("drain deployment " + replicaName)
	log.Info("sleeping till the replica pod is drained...")
	time.Sleep(time.Second * 8)

	//you have to get the most recent copy for Update to work
	origDeployment, err = client.ExtensionsV1beta1().Deployments(namespace).Get(replicaName, metav1.GetOptions{})

	if err != nil {
		log.Error(err)
		log.Error("error getting Deployment " + replicaName)
		return err
	}
	accessor, err2 := meta.Accessor(origDeployment)
	if err2 != nil {
		log.Error(err)
		log.Error("error getting accessor to deployment")
		return err
	}

	objLabels := accessor.GetLabels()
	if objLabels == nil {
		objLabels = make(map[string]string)
	}
	log.Debugf("current labels are %v\n", objLabels)

	objLabels["name"] = primaryName

	log.Debugf("updated labels are %v\n", objLabels)

	accessor.SetLabels(objLabels)
	origDeployment.Spec.Template.ObjectMeta.Labels["name"] = primaryName
	origDeployment.Spec.Selector.MatchLabels["name"] = primaryName

	//also update the PG_MODE env var to primary
	containers := origDeployment.Spec.Template.Spec.Containers
	for i := 0; i < len(containers); i++ {
		if containers[i].Name == "postgres" {
			log.Debugf("container %s \n", containers[i].Name)
			for e := 0; e < len(containers[i].Env); e++ {
				log.Printf("env %s = %s\n", containers[i].Env[e].Name, containers[i].Env[e].Value)
				if containers[i].Env[e].Name == "PG_MODE" {
					containers[i].Env[e].Value = "primary"
					log.Println("setting PG_MODE to primary")
				}
			}

		}
	}

	_, err = client.ExtensionsV1beta1().Deployments(namespace).Update(origDeployment)

	if err != nil {
		log.Error(err)
		log.Error("error updating Deployment " + replicaName)
		return err
	}

	log.Info("sleeping a bit till replica pod labels are updated...")
	time.Sleep(time.Second * 8)

	//you have to get the most recent copy for Update to work
	origDeployment, err = client.ExtensionsV1beta1().Deployments(namespace).Get(replicaName, metav1.GetOptions{})

	if err != nil {
		log.Error(err)
		log.Error("error getting Deployment " + replicaName)
		return err
	}

	log.Info("setting replica deployment replicas back to 1")
	val = 1
	origDeployment.Spec.Replicas = &val

	_, err = client.ExtensionsV1beta1().Deployments(namespace).Update(origDeployment)

	if err != nil {
		log.Error(err)
		log.Error("error updating Deployment replicas to 1 " + replicaName)
		return err
	}

	return err

}

func relabelReplica(namespace string, replica string, primary string) error {
	podsClient := client.CoreV1().Pods(namespace)
	p, err := podsClient.Get(replica, metav1.GetOptions{})

	if err != nil {
		log.Error(err.Error())
		return err
	}

	for label := range p.Labels {
		log.Debugf("label: %s ", label)
	}

	p.Labels["name"] = primary
	p, err = podsClient.Update(&apiv1.Pod{
		p.TypeMeta,
		p.ObjectMeta,
		p.Spec,
		p.Status,
	})

	if err != nil {

		log.Error(err.Error())
	}

	return err

}

func getPod(namespace string, name string) (error, *apiv1.Pod) {

	podsClient := client.CoreV1().Pods(namespace)
	p, err := podsClient.Get(name, metav1.GetOptions{})

	if err != nil {
		log.Error(err.Error())
		return err, nil
	}
	return err, p
}
func getPods(namespace string, name *string, selectors map[string]string) (*apiv1.PodList, error) {

	var b bytes.Buffer
	var options string
	var i = 0
	var l = len(selectors)

	for option, value := range selectors {
		b.WriteString(option)
		b.WriteString("=")
		b.WriteString(value)
		i++
		if i < l {
			b.WriteString(",")
		}

	}
	options = b.String()

	podsClient := client.CoreV1().Pods(namespace)
	return podsClient.List(metav1.ListOptions{
		LabelSelector: options,
	})

}

func getDeployments(namespace string, name *string, selectors map[string]string) (*v1beta1.DeploymentList, error) {

	var b bytes.Buffer
	var options string
	var i = 0
	var l = len(selectors)

	for option, value := range selectors {
		b.WriteString(option)
		b.WriteString("=")
		b.WriteString(value)
		i++
		if i < l {
			b.WriteString(",")
		}

	}
	options = b.String()

	deployments := client.ExtensionsV1beta1().Deployments(namespace)
	return deployments.List(metav1.ListOptions{
		LabelSelector: options,
	})

}
