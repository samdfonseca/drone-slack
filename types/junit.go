// Based on https://github.com/rodrigodiez/go-junit
package types

import (
	"encoding/xml"
)

type (
	JUnitFailure struct {
		XMLName xml.Name `xml:"failure"`
		Message string   `xml:"message,attr"`
		Type    string   `xml:"type,attr"`
		Text    string   `xml:",chardata"`
	}
	JUnitTestcase struct {
		XMLName  xml.Name   `xml:"testcase"`
		Id       string     `xml:"id,attr"`
		Name     string     `xml:"name,attr"`
		Time     float32    `xml:"time,attr"`
		Failures []*JUnitFailure `xml:"failure"`
	}
	JUnitTestsuite struct {
		XMLName   xml.Name    `xml:"testsuite"`
		Id        string      `xml:"id,attr"`
		Name      string      `xml:"name,attr"`
		Tests     int         `xml:"tests,attr"`
		Failures  int         `xml:"failures,attr"`
		Time      float32     `xml:"time,attr"`
		Testcases []*JUnitTestcase `xml:"testcase"`
	}
	JUnitTestsuites struct {
		XMLName    xml.Name     `xml:"testsuites"`
		Id         string       `xml:"id,attr"`
		Name       string       `xml:"name,attr"`
		Tests      int          `xml:"tests,attr"`
		Failures   int          `xml:"failures,attr"`
		Time       float32      `xml:"time,attr"`
		Testsuites []*JUnitTestsuite `xml:"testsuite"`
	}
)
