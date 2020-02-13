package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/marthjod/k8s-jobs-viz/pkg/job"
	"k8s.io/client-go/kubernetes"
)

// Handler is an HTTP handler.
type Handler struct {
	clientset *kubernetes.Clientset
	namespace string
	index     []byte
	// TODO internal cache
}

// New returns a new Handler.
func New(clientset *kubernetes.Clientset, namespace, index string) (*Handler, error) {
	b, err := ioutil.ReadFile(index)
	if err != nil {
		return nil, err
	}
	return &Handler{
		clientset: clientset,
		namespace: namespace,
		index:     b,
	}, nil
}

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jobs, err := job.Update(h.clientset, h.namespace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		b, err := json.Marshal(jobs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(b)
		return
	}

	w.Write(h.index)
}
