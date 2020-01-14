package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/marthjod/k8s-jobs-viz/cmd/config"
	"github.com/marthjod/k8s-jobs-viz/pkg/job"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func main() {

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err, "failed to get config")
	}

	var config *rest.Config
	if cfg.KubeConfig != "" {
		config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: cfg.KubeConfig},
			&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: ""}}).ClientConfig()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatal(err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	k8sJobs, err := clientset.BatchV1().Jobs(cfg.Namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	var jobs = []*job.Job{}
	for _, k8sJob := range k8sJobs.Items {
		log.Printf("%+v", k8sJob.Status)

		j := &job.Job{
			Name:           k8sJob.Name,
			StartTime:      k8sJob.Status.StartTime,
			CompletionTime: k8sJob.Status.CompletionTime,
		}
		if k8sJob.Status.Succeeded > 0 {
			j.State = job.Succeeded
		}
		if k8sJob.Status.Failed > 0 {
			j.State = job.Failed
		}
		if k8sJob.Status.Active > 0 {
			j.State = job.Running
		}

		jobs = append(jobs, j)
	}

	b, err := json.Marshal(jobs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", b)

}
