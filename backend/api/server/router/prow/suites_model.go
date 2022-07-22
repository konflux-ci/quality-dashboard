package prow

import (
	"encoding/xml"
)

// our struct which contains the complete
// array of all Users in the file
type TestXml struct {
	XMLName    xml.Name   `xml:"testsuites"`
	TestSuites TestSuites `xml:"testsuite"`
}

// a simple struct which contains all our
// test suites
type TestSuites struct {
	XMLName   xml.Name    `xml:"testsuite"`
	TestSuite []TestCases `xml:"testcase"`
}

type TestCases struct {
	XMLName xml.Name `xml:"testcase"`
	Failure Failure  `xml:"failure,omitempty"`
	Name    string   `xml:"name,attr"`
	Status  string   `xml:"status,attr"`
	Time    string   `xml:"time,attr"`
}

type Failure struct {
	XMLName xml.Name `xml:"failure,omitempty"`
	Message string   `xml:"message,attr"`
}
