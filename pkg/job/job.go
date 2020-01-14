package job

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
type Jobs []Job
