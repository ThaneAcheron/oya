package testutil

import (
	"testing"

	"github.com/bilus/oya/pkg/oyafile"
	"github.com/bilus/oya/pkg/project"
)

func MustListOyafiles(t *testing.T, rootDir string) []*oyafile.Oyafile {
	project, err := project.Load(rootDir)
	AssertNoErr(t, err, "Error detecting project")
	oyafiles, err := project.Oyafiles()
	AssertNoErr(t, err, "Error listing Oyafiles")
	AssertTrue(t, len(oyafiles) > 0, "No Oyafiles found")
	return oyafiles
}

func MustLoadOyafile(t *testing.T, dir, rootDir string) *oyafile.Oyafile {
	o, found, err := oyafile.LoadFromDir(dir, rootDir)
	AssertNoErr(t, err, "Error loading root Oyafile")
	AssertTrue(t, found, "Root Oyafile not found")
	return o
}
