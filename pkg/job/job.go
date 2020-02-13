package job

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// State represents Job state.
type State int

//go:generate jsonenums -type=State
//go:generate stringer -type=State
const (
	Succeeded State = iota
	Failed    State = iota
	Running   State = iota
)

// Job represents a K8s Job abstraction.
type Job struct {
	Name           string       `json:"name"`
	State          State        `json:"state"`
	StartTime      *metav1.Time `json:"start_time"`
	CompletionTime *metav1.Time `json:"completion_time"`
}

// Jobs is a list of Job.
type Jobs []*Job

// Update queries the K8s cluster for jobs.
func Update(clientset *kubernetes.Clientset, namespace string) (Jobs, error) {
	k8sJobs, err := clientset.BatchV1().Jobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var jobs = Jobs{}
	for _, k8sJob := range k8sJobs.Items {
		j := &Job{
			Name:           k8sJob.Name,
			StartTime:      k8sJob.Status.StartTime,
			CompletionTime: k8sJob.Status.CompletionTime,
		}
		if k8sJob.Status.Succeeded > 0 {
			j.State = Succeeded
		}
		if k8sJob.Status.Failed > 0 {
			j.State = Failed
		}
		if k8sJob.Status.Active > 0 {
			j.State = Running
		}

		jobs = append(jobs, j)
	}

	return jobs, nil
}
