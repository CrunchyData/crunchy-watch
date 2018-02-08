package main

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/crunchydata/crunchy-watch/util"
)


func deletePrimaryPod(namespace string, name string ) error {

	podsClient := client.CoreV1().Pods(namespace)
	err := podsClient.Delete(name, nil)

	if err != nil {

		log.Error(err.Error())
	}

	return err
}

/*
get replica by namespace and name
 */
func promoteReplica(namespace string, name string) error {


	pod, err := client.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})

	if err != nil {
		return  fmt.Errorf("could not get pod info: %v", err)
	}

	if len(pod.Spec.Containers) != 1 {
		return  fmt.Errorf("could not determine which container to use")
	}


	err = util.Exec(restConfig, pod.Namespace, pod.Name, pod.Spec.Containers[0].Name, []string{ "touch", "/tmp/pg-failover-trigger"})

	if err != nil {
		return  fmt.Errorf("could not execute: %v", err)
	}

	return err
}

func relabelReplica(namespace string, replica string, primary string) error {
	podsClient := client.CoreV1().Pods(namespace)
	p,err := podsClient.Get(replica,metav1.GetOptions{})

	if err != nil {
		log.Error(err.Error())
		return err
	}

	p,err = podsClient.Update(&apiv1.Pod {
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
	p,err := podsClient.Get(name,metav1.GetOptions{})

	if err != nil {
		log.Error(err.Error())
		return err, nil
	}
	return err, p
}
func getPods(namespace string, name *string, selectors map[string]string )(*apiv1.PodList, error) {

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
	return podsClient.List( metav1.ListOptions{
		LabelSelector: options,
	})



}