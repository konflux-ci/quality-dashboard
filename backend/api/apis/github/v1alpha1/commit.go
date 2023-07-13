package v1alpha1

import (
	"time"
)

type Commit struct {
	// OID identifies the commit sha
	OID           string
	CommittedDate time.Time
}
