package project_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/bilus/oya/pkg/project"
	tu "github.com/bilus/oya/testutil"
)

func TestProject_Detect_NoOya(t *testing.T) {
	workDir := "./fixtures/empty_project"
	_, err := project.Detect(workDir)
	tu.AssertErr(t, err, "Expected error trying to detect Oya project in empty dir")
}

func TestProject_Detect_InRootDir(t *testing.T) {
	workDir := "./fixtures/project"
	_, err := project.Detect(workDir)
	tu.AssertNoErr(t, err, "Expected no error trying to detect Oya project in its root dir")
}

func TestProject_Detect_InSubDir(t *testing.T) {
	workDir := "./fixtures/project//subdir"
	_, err := project.Detect(workDir)
	tu.AssertNoErr(t, err, "Expected no error trying to detect Oya project in its root dir")
}

func TestProject_Detect_InEmptySubDir(t *testing.T) {
	workDir := "./fixtures/project/empty_subdir"
	_, err := project.Detect(workDir)
	tu.AssertNoErr(t, err, "Expected no error trying to detect Oya project in its root dir")
}

func TestProject_Run_NoOyafile(t *testing.T) {
	workDir := "./fixtures/project/empty_subdir"
	project, err := project.Detect(workDir)
	tu.AssertNoErr(t, err, "Expected no error trying to detect Oya project in its root dir")
	err = project.Run(workDir, "build", ioutil.Discard, ioutil.Discard)
	tu.AssertErr(t, err, "Expected error trying to run without Oyafile")
}

func TestProject_Run_NoTask(t *testing.T) {
	workDir := "./fixtures/project"
	project, err := project.Detect(workDir)
	tu.AssertNoErr(t, err, "Expected no error trying to detect Oya project in its root dir")
	err = project.Run(workDir, "noSuchTask", ioutil.Discard, ioutil.Discard)
	tu.AssertErr(t, err, "Expected error when trying to run without matching task")
}

func TestProject_Run_NoChanges(t *testing.T) {
	workDir := "./fixtures/empty_changeset_project"
	project, err := project.Detect(workDir)
	tu.AssertNoErr(t, err, "Expected no error trying to detect Oya project in its root dir")
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	err = project.Run(workDir, "build", stdout, stderr)
	tu.AssertNoErr(t, err, "Expected no error running with empty changeset")
	tu.AssertEqual(t, 0, len(stdout.String()))
	tu.AssertEqual(t, 0, len(stderr.String()))
}

func TestProject_Run_WithChanges(t *testing.T) {
	workDir := "./fixtures/project"
	project, err := project.Detect(workDir)
	tu.AssertNoErr(t, err, "Expected no error trying to detect Oya project in its root dir")
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	err = project.Run(workDir, "build", stdout, stderr)
	tu.AssertNoErr(t, err, "Expected no error running non-empty changeset")
	tu.AssertEqual(t, "build run", stdout.String())
}