// Copyright (c) 2018-2024, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package sources

import (
	"path/filepath"
	"testing"

	"github.com/sylabs/singularity/v4/internal/pkg/test"
	"github.com/sylabs/singularity/v4/internal/pkg/util/fs"
)

func testWithGoodDir(t *testing.T, f func(d string) error) {
	if err := f(t.TempDir()); err != nil {
		t.Fatalf("Unexpected failure: %v", err)
	}
}

func testWithBadDir(t *testing.T, f func(d string) error) {
	if err := f("/this/will/be/a/problem"); err == nil {
		t.Fatalf("Unexpected success with bad directory")
	}
}

func TestMakeDirs(t *testing.T) {
	test.DropPrivilege(t)
	defer test.ResetPrivilege(t)

	testWithGoodDir(t, makeDirs)
	testWithBadDir(t, makeDirs)
}

func TestMakeSymlinks(t *testing.T) {
	test.DropPrivilege(t)
	defer test.ResetPrivilege(t)

	testWithGoodDir(t, makeSymlinks)
	testWithBadDir(t, makeSymlinks)
}

func TestMakeFiles(t *testing.T) {
	test.DropPrivilege(t)
	defer test.ResetPrivilege(t)

	testWithGoodDir(t, func(d string) error {
		if err := makeDirs(d); err != nil {
			return err
		}
		return makeFiles(d, false)
	})
	testWithBadDir(t, func(d string) error { return makeFiles(d, false) })
	// #4532 - Check that we can succeed with an existing file that doesn't have
	// write permission.
	testWithGoodDir(t, func(d string) error {
		if err := makeDirs(d); err != nil {
			return err
		}
		err := fs.EnsureFileWithPermission(filepath.Join(d, "etc", "hosts"), 0o400)
		if err != nil {
			t.Fatalf("Failed to make test hosts file: %s", err)
		}
		return makeFiles(d, false)
	})
}

func TestMakeBaseEnv(t *testing.T) {
	test.DropPrivilege(t)
	defer test.ResetPrivilege(t)

	testWithGoodDir(t, func(d string) error { return makeBaseEnv(d, false) })
	testWithBadDir(t, func(d string) error { return makeBaseEnv(d, false) })
}
