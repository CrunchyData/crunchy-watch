package main

import (
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"

	"k8s.io/client-go/tools/remotecommand"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"bytes"
	"fmt"
	"io"
	log "github.com/sirupsen/logrus"
)

type ExecOptions struct {
	Command []string

	Namespace     string
	PodName       string
	ContainerName string

	Stdin         io.Reader
	CaptureStdout bool
	CaptureStderr bool
	// If false, whitespace in std{err,out} will be removed.
	PreserveWhitespace bool
}

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

	var (
		execOut bytes.Buffer
		execErr bytes.Buffer
	)

	pod, err := client.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})

	if err != nil {
		return  fmt.Errorf("could not get pod info: %v", err)
	}

	if len(pod.Spec.Containers) != 1 {
		return  fmt.Errorf("could not determine which container to use")
	}

	req := client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec")

	req.VersionedParams(&apiv1.PodExecOptions{
		Container: pod.Spec.Containers[0].Name,
		Command:   []string{ "touch", "/tmp/pg-failover-trigger"},
		Stdout:    true,
		Stderr:    true,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(restConfig, "POST", req.URL())
	if err != nil {
		return  fmt.Errorf("failed to init executor: %v", err)
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdout:             &execOut,
		Stderr:             &execErr,
	})

	if err != nil {
		return  fmt.Errorf("could not execute: %v", err)
	}

	if execErr.Len() > 0 {
		log.Info(execErr.String())
		return  fmt.Errorf("stderr: %v", execErr.String())
	}

	log.Info(execOut.String())

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