package ociloader

import (
	"encoding/xml"
	"os"

	"github.com/konflux-ci/quality-dashboard/api/server/router/prow"
)

// ParseXUnitReport reads and parses an xUnit XML report from the specified file path.
// It takes the file path of the xUnit report as input and returns a pointer to a
// prow.TestSuites struct, which represents the parsed test suites, or an error if
// the file cannot be read or unmarshaled.
func ParseXUnitReport(path string) (*prow.TestSuites, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var suites prow.TestSuites

	if err := xml.Unmarshal(data, &suites); err != nil {
		return nil, err
	}

	return &suites, nil
}
