package prow

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"go.uber.org/zap"
)

// The below types are directly marshalled into XML. The types correspond to jUnit
// XML schema, but do not contain all valid fields. For instance, the class name
// field for test cases is omitted, as this concept does not directly apply to Go.
// For XML specifications see http://help.catchsoftware.com/display/ET/JUnit+Format
// or view the XSD included in this package as 'junit.xsd'

// TestSuites represents a flat collection of jUnit test suites.
type TestSuites struct {
	XMLName xml.Name `xml:"testsuites"`

	// Suites are the jUnit test suites held in this collection
	Suites []*TestSuite `xml:"testsuite"`
}

// TestSuite represents a single jUnit test suite, potentially holding child suites.
type TestSuite struct {
	XMLName xml.Name `xml:"testsuite"`

	// Name is the name of the test suite
	Name string `xml:"name,attr"`

	// NumTests records the number of tests in the TestSuite
	NumTests uint `xml:"tests,attr"`

	// NumSkipped records the number of skipped tests in the suite
	NumSkipped uint `xml:"skipped,attr"`

	// NumFailed records the number of failed tests in the suite
	NumFailed uint `xml:"failures,attr"`

	// Duration is the time taken in seconds to run all tests in the suite
	Duration float64 `xml:"time,attr"`

	// Properties holds other properties of the test suite as a mapping of name to value
	Properties []*TestSuiteProperty `xml:"properties,omitempty"`

	// TestCases are the test cases contained in the test suite
	TestCases []*TestCase `xml:"testcase"`

	// Children holds nested test suites
	Children []*TestSuite `xml:"testsuite"`
}

// TestSuiteProperty contains a mapping of a property name to a value
type TestSuiteProperty struct {
	XMLName xml.Name `xml:"properties"`

	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// TestCase represents a jUnit test case
type TestCase struct {
	XMLName xml.Name `xml:"testcase"`

	// Name is the name of the test case
	Name string `xml:"name,attr"`

	// Classname is an attribute set by the package type and is required
	Classname string `xml:"classname,attr,omitempty"`

	// Duration is the time taken in seconds to run the test
	Duration float64 `xml:"time,attr"`

	// Status show if the test case was succed or not
	Status string `xml:"status,attr"`

	// SkipMessage holds the reason why the test was skipped
	SkipMessage *SkipMessage `xml:"skipped"`

	// FailureOutput holds the output from a failing test
	FailureOutput *FailureOutput `xml:"failure"`

	// SystemOut is output written to stdout during the execution of this test case
	SystemOut string `xml:"system-out,omitempty"`

	// SystemErr is output written to stderr during the execution of this test case
	SystemErr string `xml:"system-err,omitempty"`
}

// SkipMessage holds a message explaining why a test was skipped
type SkipMessage struct {
	XMLName xml.Name `xml:"skipped"`

	// Message explains why the test was skipped
	Message string `xml:"message,attr,omitempty"`
}

// FailureOutput holds the output from a failing test
type FailureOutput struct {
	XMLName xml.Name `xml:"failure"`

	// Message holds the failure message from the test
	Message string `xml:"message,attr"`

	// Output holds verbose failure output from the test
	Output string `xml:",chardata"`
}

// TestResult is the result of a test case
type TestResult string

func parseFileFromRequest(r *http.Request, logger *zap.Logger) (TestSuites, error) {
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `file`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("file")
	if err != nil {
		logger.Sugar().Infof("Failed to get file from header %s", err)
		return suitesXml, err
	}

	// copy example
	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Sugar().Infof("Failed to create temporary file to archive the tests results %s", err)
		return suitesXml, err
	}

	io.Copy(f, file)
	defer file.Close()

	xmlFile, err := os.Open(f.Name())
	// if we os.Open returns an error then handle it
	if err != nil {
		logger.Sugar().Infof("Failed to open xml file %s", err)
		return suitesXml, err
	}

	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// Remove the file after reading
	os.Remove(handler.Filename)

	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'suitesXml' which we defined above
	if err := xml.Unmarshal(byteValue, &suitesXml); err != nil {
		if err != nil {
			logger.Sugar().Infof("Failed convert xml file to golang bytes %s", err)
			return suitesXml, err
		}
	}

	return suitesXml, nil
}
