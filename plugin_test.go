package main_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	droneSlack "github.com/drone-plugins/drone-slack"
)


const (
	junitXml = `<testsuites name="plugin_test" tests="6" failures="3" time="6.006">
  <testsuite name="suite1" errors="0" failures="1" skipped="1" timestamp="2018-10-17T17:27:23" time="3.003" tests="3">
    <testcase classname="suite1 test1" name="suite1 test1" time="1.001">
      <skipped/>
    </testcase>
    <testcase classname="suite1 test2" name="suite1 test2" time="1.001">
      <failure>Error: expect(received).toHaveLength(length)

Expected value to have length:
  1
Received:
  []
received.length:
  0
      </failure>
    </testcase>
    <testcase classname="suite1 test3" name="suite1 test3" time="1.001">
    </testcase>
  </testsuite>
  <testsuite name="suite2" errors="0" failures="1" skipped="1" timestamp="2018-10-17T17:27:23" time="3.003" tests="3">
    <testcase classname="suite2 test1" name="suite2 test1" time="1.001">
      <skipped/>
    </testcase>
    <testcase classname="suite2 test2" name="suite2 test2" time="1.001">
      <failure>Error: expect(received).toHaveLength(length)

Expected value to have length:
  1
Received:
  []
received.length:
  0
      </failure>
    </testcase>
    <testcase classname="suite2 test3" name="suite2 test3" time="1.001">
    </testcase>
  </testsuite>
</testsuites>`
)

func TestPlugin_PrepPayload_WithJUnitReport(t *testing.T) {
	slackTemplate := `{{#success build.status}}
build {{build.number}} of {{repo.owner}}/{{repo.name}} succeeded
{{else}}
build {{build.number}} of {{repo.owner}}/{{repo.name}} failed
{{/success}}
{{#if build.jUnitReport}}
  {{#each build.jUnitReport.testsuites}}
testsuite {{name}} has failed tests
    {{#each testcases}}
      {{#if failures}}
{{name}} failed
      {{/if}}
    {{/each}}
  {{/each}}
{{/if}}
`
	f, err := ioutil.TempFile("", "drone-slack-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := os.Remove(f.Name())
		if err != nil {
			t.Log(err)
		}
	}()
	_, err = f.WriteString(junitXml)
	if err != nil {
		t.Fatal(err)
	}
	p := &droneSlack.Plugin{
		Repo: droneSlack.Repo{"drone-plugins", "drone-slack"},
		Build: droneSlack.Build{
			Tag: "v0.0.0",
			Event: "push",
			Number: 1,
			Commit: "62126a02ffea3dabd7789e5c5407553490973665",
			Ref: "refs/heads/master",
			Branch: "master",
			Author: "drone-plugins",
			Pull: "123",
			Message: "Test Build",
			DeployTo: "",
			Status: "pending",
			Link: "https://github.com/drone-plugins/drone-slack/commit/62126a02ffea3dabd7789e5c5407553490973665",
			Started: 0,
			Created: time.Now().Unix(),
		},
		Config: droneSlack.Config{
			Channel: "#drone-slack",
			Template: slackTemplate,
			JUnitResults: f.Name(),
		},
	}
	payload, err := p.PrepPayload()
	if err != nil {
		t.Fatal(err)
	}
	expectedLines := []string{
		fmt.Sprintf("build %d of %s/%s failed", p.Build.Number, p.Repo.Owner, p.Repo.Name),
		"testsuite suite1 has failed tests",
		"    suite1 test2 failed",
		"",
		"testsuite suite2 has failed tests",
		"    suite2 test2 failed",
	}
	actualLines := strings.Split(payload.Attachments[0].Text, "\n")
	for i := range expectedLines {
		if expectedLines[i] != "" && actualLines[i] != expectedLines[i] {
			t.Fatalf("Line %d does not match\nActual: %s\nExpected: %s", i, actualLines[i], expectedLines[i])
		}
	}
}
