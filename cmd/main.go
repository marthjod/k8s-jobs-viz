package main

import (
	"log"
	"net/http"

	"github.com/marthjod/k8s-jobs-viz/cmd/config"
	"github.com/marthjod/k8s-jobs-viz/pkg/handler"
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

	hdlr, err := handler.New(clientset, cfg.Namespace, cfg.IndexHTML)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", hdlr)
	log.Printf("listening on %s\n", cfg.ListenAddress)
	log.Fatal(http.ListenAndServe(cfg.ListenAddress, nil))

}
