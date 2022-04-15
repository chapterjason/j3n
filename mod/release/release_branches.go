/*
 * Copyright © 2022 Jason Schilling
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package release

import (
	"fmt"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"

	"github.com/chapterjason/j3n/mod/version"
)

type GitReleaseBranchFormatter = func(v version.Version) string
type GitReleaseBranchChecker = func(s string) bool

var (
	DefaultGitReleaseBranchExpression                           = regexp.MustCompile("^release/\\d+\\.\\d+$")
	DefaultGitReleaseBranchFormatter  GitReleaseBranchFormatter = func(v version.Version) string {
		return fmt.Sprintf("release/%d.%d", v.Major, v.Minor)
	}
	DefaultGitReleaseBranchChecker GitReleaseBranchChecker = func(s string) bool {
		return DefaultGitReleaseBranchExpression.MatchString(s)
	}
)

func Branches(r *git.Repository) (storer.ReferenceIter, error) {
	branchIter, err := r.Branches()

	if err != nil {
		return nil, err
	}

	return storer.NewReferenceFilteredIter(
		func(ref *plumbing.Reference) bool {
			return ref.Name().IsBranch() && IsReleaseBranch(ref.Name().Short())
		}, branchIter,
	), nil
}

func IsReleaseBranch(short string) bool {
	return DefaultGitReleaseBranchChecker(short)
}
