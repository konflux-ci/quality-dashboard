package ociloader

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseXUnitReport_Success(t *testing.T) {
	tmpDir := t.TempDir()
	reportPath := filepath.Join(tmpDir, "report.xml")

	content := `
<testsuites>
  <testsuite name="SampleSuite" tests="1" failures="1">
    <testcase name="should do something" classname="SampleSuite" time="0.001">
      <failure message="expected true but got false">Assertion failed</failure>
    </testcase>
  </testsuite>
</testsuites>`

	err := os.WriteFile(reportPath, []byte(content), 0644)
	assert.NoError(t, err)

	suites, err := ParseXUnitReport(reportPath)
	assert.NoError(t, err)
	assert.NotNil(t, suites)
	assert.Len(t, suites.Suites, 1)

	suite := suites.Suites[0]
	assert.Equal(t, "SampleSuite", suite.Name)
	assert.Len(t, suite.TestCases, 1)
	assert.Equal(t, "should do something", suite.TestCases[0].Name)
	assert.NotNil(t, suite.TestCases[0].FailureOutput)
	assert.Equal(t, "expected true but got false", suite.TestCases[0].FailureOutput.Message)
}

func TestParseXUnitReport_InvalidXML(t *testing.T) {
	tmpDir := t.TempDir()
	reportPath := filepath.Join(tmpDir, "invalid.xml")

	content := `<testsuites><bad></xml>`

	err := os.WriteFile(reportPath, []byte(content), 0644)
	assert.NoError(t, err)

	suites, err := ParseXUnitReport(reportPath)
	assert.Error(t, err)
	assert.Nil(t, suites)
}

func TestParseXUnitReport_FileNotFound(t *testing.T) {
	suites, err := ParseXUnitReport("/non/existent/file.xml")
	assert.Error(t, err)
	assert.Nil(t, suites)
}
