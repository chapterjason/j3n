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
	"github.com/go-git/go-git/v5"

	"github.com/chapterjason/j3n/mod/version"
)

type MultiBranchWorkflow struct {
	BranchFormat        string `json:"branch_format,omitempty"`
	BranchExpression    string `json:"branch_expression,omitempty"`
	UpdateMessageFormat string `json:"update_message_format,omitempty"`
	BumpMessageFormat   string `json:"bump_message_format,omitempty"`
}

func (m *MultiBranchWorkflow) PreRelease(r *git.Repository, v version.Version, rt ReleaseType) error {
	if rt == ReleaseTypePatch {
		return m.preReleasePatch(r, v)
	} else if rt == ReleaseTypeMinor {
		return m.preReleaseMinor(r, v)
	} else if rt == ReleaseTypeMajor {
		return m.preReleaseMajor(r, v)
	}

	return nil
}

func (m *MultiBranchWorkflow) PostRelease(r *git.Repository, v version.Version, rt ReleaseType) error {
	if rt == ReleaseTypePatch {
		return m.postReleasePatch(r, v)
	} else if rt == ReleaseTypeMinor {
		return m.postReleaseMinor(r, v)
	} else if rt == ReleaseTypeMajor {
		return m.postReleaseMajor(r, v)
	}

	return nil
}

// preReleasePatch checkout the release branch, update version, and commit
func (m *MultiBranchWorkflow) preReleasePatch(r *git.Repository, v version.Version) error {
	rbs := GitReleaseBranchFormatter(v)

	b, err := r.Branch(rbs)

	if err != nil {
		return err
	}

	w, err := r.Worktree()

	if err != nil {
		return err
	}

	err = w.Checkout(
		&git.CheckoutOptions{
			Branch: b.Merge,
		},
	)

	if err != nil {
		return err
	}

	err = version.Set(v)

	if err != nil {
		return err
	}

	_, err = w.Commit(
		version.Replace(m.UpdateMessageFormat, v),
		&git.CommitOptions{
			All: true,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

// preReleaseMinor bump version, and commit
func (m *MultiBranchWorkflow) postReleasePatch(r *git.Repository, v version.Version) error {
	w, err := r.Worktree()

	if err != nil {
		return err
	}

	v.Patch++
	v.Prerelease = []string{"DEV"}

	err = version.Set(v)

	if err != nil {
		return err
	}

	_, err = w.Commit(
		version.Replace(m.BumpMessageFormat, v),
		&git.CommitOptions{
			All: true,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *MultiBranchWorkflow) preReleaseMinor(r *git.Repository, v version.Version) error {
	panic("implement me")
}

func (m *MultiBranchWorkflow) postReleaseMinor(r *git.Repository, v version.Version) error {
	panic("implement me")
}

func (m *MultiBranchWorkflow) preReleaseMajor(r *git.Repository, v version.Version) error {
	panic("implement me")
}

func (m *MultiBranchWorkflow) postReleaseMajor(r *git.Repository, v version.Version) error {
	panic("implement me")
}
