package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Podpod struct {
	Pod_name string `json:Deployment,string`
	Status   string `json:Status,string`
	K8s_node string `json:K8s_node,string`
}
type Datapod struct {
	Deployment string   `json:Deployment,string`
	Namespace  string   `json:Namespace,string`
	Pod        []Podpod `json:Pod`
}
type Basepod struct {
	Status int
	Data   []Datapod `json:Data`
}

func Pods(w http.ResponseWriter, r *http.Request) {
	deploy := chi.URLParam(r, "deploy")
	w.WriteHeader(http.StatusOK)
	var clientset *kubernetes.Clientset
	retd := new(Basepod)
	retp := new(Datapod)
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
			if pod.ObjectMeta.OwnerReferences[0].Kind == "ReplicaSet" && strings.Contains(pod.ObjectMeta.OwnerReferences[0].Name, deploy) {
				podpod := Podpod{Pod_name: pod.ObjectMeta.Name, Status: string(pod.Status.Phase), K8s_node: pod.Spec.NodeName}
				retp.Pod = append(retp.Pod, podpod)
			}
		}
		for _, deployment := range deployments.Items {
			if deployment.Name == deploy {
				datapod := Datapod{Deployment: deployment.Name, Namespace: deployment.Namespace, Pod: retp.Pod}
				retd.Data = append(retd.Data, datapod)
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
