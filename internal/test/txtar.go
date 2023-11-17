// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/txtar"
)

func WriteTxtar(filename string, files []txtar.File, comment string) error {
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile(filename, txtar.Format(
		&txtar.Archive{
			Comment: []byte(addBoilerplate(currentYear(), comment)),
			Files:   files,
		},
	), 0666); err != nil {
		return err
	}

	return nil
}

// addBoilerplate adds the copyright string for the given year to the
// given comment, and some additional spacing for readability.
func addBoilerplate(year int, comment string) string {
	return fmt.Sprintf(`Copyright %d The Go Authors. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.

%s

`, year, comment)
}

func currentYear() int {
	year, _, _ := time.Now().Date()
	return year
}

var copyrightRE = regexp.MustCompile(`Copyright (\d+)`)

// findCopyrightYear returns the copyright year in this comment,
// or an error if none is found.
func findCopyrightYear(comment string) (int, error) {
	matches := copyrightRE.FindStringSubmatch(comment)
	if len(matches) != 2 {
		return 0, errors.New("comment does not contain a copyright year")
	}
	year, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}
	return year, nil
}

// CheckComment checks the validity of a txtar comment.
// It checks that the "got" comment is the same as would be generated
// by WriteTxtar(..., wantComment), but allows any copyright year.
//
// For testing.
func CheckComment(wantComment, got string) error {
	year, err := findCopyrightYear(got)
	if err != nil {
		return err
	}

	want := addBoilerplate(year, wantComment)
	if diff := cmp.Diff(want, got); diff != "" {
		return fmt.Errorf("comment mismatch (-want, +got):\n%s", diff)
	}

	return nil
}
