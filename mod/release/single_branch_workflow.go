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
	"strings"

	"github.com/gogs/git-module"
	"github.com/pkg/errors"

	"github.com/chapterjason/j3n/mod/version"
	"github.com/chapterjason/j3n/modx/gitx"
)

type SingleBranchWorkflow struct {
	Branch              string `json:"branch,omitempty"`
	TagFormat           string `json:"tag_format,omitempty"`
	UpdateMessageFormat string `json:"update_message_format,omitempty"`
	BumpMessageFormat   string `json:"bump_message_format,omitempty"`
}

func (sbw *SingleBranchWorkflow) GetUpdateMessageFormat() string {
	return sbw.UpdateMessageFormat
}

func (sbw *SingleBranchWorkflow) GetBumpMessageFormat() string {
	return sbw.BumpMessageFormat
}

func (sbw *SingleBranchWorkflow) GetBranch(_ version.Version) string {
	return sbw.Branch
}

func (sbw *SingleBranchWorkflow) GetTag(v version.Version) string {
	return version.Replace(sbw.TagFormat, v)
}

func (sbw *SingleBranchWorkflow) PreRelease(r *git.Repository, v version.Version) error {
	rbs := sbw.GetBranch(v)

	if !r.HasBranch(rbs) {
		return fmt.Errorf("missing branch %s", rbs)
	}

	err := r.Checkout(rbs)

	if err != nil {
		return errors.Wrap(err, "failed to checkout branch")
	}

	err = version.Set(v)

	if err != nil {
		return errors.Wrap(err, "failed to set version")
	}

	sig, err := gitx.GetSignature(r)

	if err != nil {
		return errors.Wrap(err, "failed to get signature")
	}

	message := version.Replace(sbw.UpdateMessageFormat, v)

	err = r.Commit(sig, message)

	if err != nil {
		return errors.Wrap(err, "failed to commit")
	}

	return nil
}

func (sbw *SingleBranchWorkflow) Release(r *git.Repository, v version.Version) error {
	tags, err := r.Tags()

	if err != nil {
		return errors.Wrap(err, "failed to get tags")
	}

	for _, tag := range tags {
		tagVersion, err := version.Parse(strings.TrimPrefix(tag, "v"))

		if err != nil {
			return errors.Wrap(err, "failed to parse tag")
		}

		if tagVersion.Compare(v) > 0 {
			return fmt.Errorf("tag %s is higher than version %s", tag, v)
		}
	}

	err = sbw.PreRelease(r, v)

	if err != nil {
		return errors.Wrap(err, "failed to pre-release")
	}

	tn := sbw.GetTag(v)

	currentRevision, err := r.RevParse("HEAD")

	if err != nil {
		return errors.Wrap(err, "failed to get current revision")
	}

	err = r.CreateTag(tn, currentRevision)

	if err != nil {
		return errors.Wrap(err, "failed to create tag")
	}

	err = sbw.PostRelease(r, v)

	if err != nil {
		return errors.Wrap(err, "failed to post-release")
	}

	return nil
}

func (sbw *SingleBranchWorkflow) PostRelease(r *git.Repository, v version.Version) error {
	v.Patch++
	v.Prerelease = []string{"DEV"}

	err := version.Set(v)

	if err != nil {
		return errors.Wrap(err, "failed to set version")
	}

	sig, err := gitx.GetSignature(r)

	if err != nil {
		return errors.Wrap(err, "failed to get signature")
	}

	message := version.Replace(sbw.BumpMessageFormat, v)

	err = r.Commit(sig, message)

	if err != nil {
		return errors.Wrap(err, "failed to commit")
	}

	return nil
}
