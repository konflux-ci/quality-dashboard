package prow

import (
	"regexp"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ProwJobState specifies whether the job is running
type ProwJobState string

// Various job states.
const (
	// TriggeredState means the job has been created but not yet scheduled.
	TriggeredState ProwJobState = "triggered"
	// PendingState means the job is currently running and we are waiting for it to finish.
	PendingState ProwJobState = "pending"
	// SuccessState means the job completed without error (exit 0)
	SuccessState ProwJobState = "success"
	// FailureState means the job completed with errors (exit non-zero)
	FailureState ProwJobState = "failure"
	// AbortedState means prow killed the job early (new commit pushed, perhaps).
	AbortedState ProwJobState = "aborted"
	// ErrorState means the job could not schedule (bad config, perhaps).
	ErrorState ProwJobState = "error"
)

var RegexpCompiler = regexp.MustCompile("(-main-|-master-)(.*?)(\\/)")

type Extrarefs struct {
	// Org is something like kubernetes or k8s.io
	Org string `json:"org"`
	// Repo is something like test-infra
	Repo string `json:"repo"`
}

type ProwJobSpec struct {
	Type    string `json:"type,omitempty"`
	Cluster string `json:"cluster,omitempty"`
	Job     string `json:"job,omitempty"`
	Refs    Refs   `json:"refs"`

	// Refs is the code under test, determined at runtime by Prow itself
	Extrarefs []Extrarefs `json:"extra_refs"`
}

type Refs struct {
	// Org is something like kubernetes or k8s.io
	Org string `json:"org"`
	// Repo is something like test-infra
	Repo string `json:"repo"`
	// RepoLink links to the source for Repo.
	RepoLink string `json:"repo_link,omitempty"`

	BaseRef string `json:"base_ref,omitempty"`
	BaseSHA string `json:"base_sha,omitempty"`
	// BaseLink is a link to the commit identified by BaseSHA.
	BaseLink string `json:"base_link,omitempty"`
}

type ProwJobStatus struct {
	StartTime      time.Time    `json:"startTime,omitempty"`
	PendingTime    *time.Time   `json:"pendingTime,omitempty"`
	CompletionTime *time.Time   `json:"completionTime,omitempty"`
	State          ProwJobState `json:"state,omitempty"`
	Description    string       `json:"description,omitempty"`
	URL            string       `json:"url,omitempty"`
	PodName        string       `json:"pod_name,omitempty"`
	BuildID        string       `json:"build_id,omitempty"`
}

type ProwJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProwJobSpec   `json:"spec,omitempty"`
	Status ProwJobStatus `json:"status,omitempty"`
}
