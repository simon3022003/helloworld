package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Datastatus struct {
	Pod_name   string `json:Deployment,string`
	Namespace  string `json:Namespace,string`
	Deployment string `json:Deployment,string`
	Status     string `json:Status,string`
	K8s_node   string `json:K8s_node,string`
}
type Basestatus struct {
	Status int
	Data   []Datastatus `json:Data`
}

func Podstatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.URL.RawQuery
	key := strings.Split(ctx, "=")[0]
	if key != "status" {
		fmt.Fprintln(w, "you key is '", key, "'. Please enter 'status' as query key and value in 'Pending, Running, Succeeded, Failed, Unknown'.")
		return
	}
	value := strings.Split(ctx, "=")[1]
	w.WriteHeader(http.StatusOK)
	var clientset *kubernetes.Clientset
	retd := new(Basestatus)
	retd.Status = http.StatusOK
	// path-to-kubeconfig -- for example, /root/.kube/config
	config, err := clientcmd.BuildConfigFromFlags("", "/app/config")
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	// access the API to list deployments
	namespace, err := clientset.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, ns := range namespace.Items {
		pods, err := clientset.CoreV1().Pods(ns.Name).List(context.TODO(), v1.ListOptions{})
		deployments, err := clientset.AppsV1().Deployments(ns.Name).List(context.TODO(), v1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		for _, pod := range pods.Items {
			if pod.ObjectMeta.OwnerReferences[0].Kind == "ReplicaSet" && string(pod.Status.Phase) == value {
				for _, deployment := range deployments.Items {
					if strings.Contains(pod.ObjectMeta.OwnerReferences[0].Name, deployment.Name) {
						datastatus := Datastatus{Pod_name: pod.ObjectMeta.Name, Namespace: pod.ObjectMeta.Namespace, Deployment: deployment.Name, Status: string(pod.Status.Phase), K8s_node: pod.Spec.NodeName}
						retd.Data = append(retd.Data, datastatus)
					}
				}
			}
		}
	}
	ret_json, err := json.MarshalIndent(retd, "", "\t")
	if err != nil {
		panic(err.Error())
	}
	w.Write([]byte(ret_json))
	return
}
