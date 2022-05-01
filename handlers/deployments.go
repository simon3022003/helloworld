package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Datadeploy struct {
	Deployment string `json:Deployment,string`
	Namespace  string `json:Namespace,string`
	Portcount  string `json:Portcount,string`
}
type Basedeploy struct {
	Status int
	Data   []Datadeploy `json:Data`
}

func Deployments(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	var clientset *kubernetes.Clientset
	retd := new(Basedeploy)
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
		deployments, err := clientset.AppsV1().Deployments(ns.Name).List(context.TODO(), v1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		for _, deployment := range deployments.Items {
			data := Datadeploy{Deployment: deployment.Name, Namespace: deployment.Namespace, Portcount: strconv.FormatInt(int64(deployment.Status.Replicas), 10)}
			retd.Data = append(retd.Data, data)
		}
	}
	ret_json, err := json.MarshalIndent(retd, "", "\t")
	if err != nil {
		panic(err.Error())
	}
	w.Write([]byte(ret_json))
	return
}
