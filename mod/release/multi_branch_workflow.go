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

	"github.com/gogs/git-module"
	"github.com/pkg/errors"

	"github.com/chapterjason/j3n/mod/version"
	"github.com/chapterjason/j3n/modx/gitx"
)

type MultiBranchWorkflow struct {
	UpdateMessageFormat string `json:"update_message_format,omitempty"`
	BumpMessageFormat   string `json:"bump_message_format,omitempty"`
}

func (mbw *MultiBranchWorkflow) PreRelease(r *git.Repository, v version.Version) error {
	rbs := BranchFormatter(v)

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

	message := version.Replace(mbw.UpdateMessageFormat, v)

	err = r.Commit(sig, message)

	if err != nil {
		return errors.Wrap(err, "failed to commit")
	}

	return nil
}

func (mbw *MultiBranchWorkflow) Release(r *git.Repository, v version.Version) error {
	err := mbw.PreRelease(r, v)

	if err != nil {
		return errors.Wrap(err, "failed to pre-release")
	}

	tn := TagFormatter(v)

	currentRevision, err := r.RevParse("HEAD")

	if err != nil {
		return errors.Wrap(err, "failed to get current revision")
	}

	err = r.CreateTag(tn, currentRevision)

	if err != nil {
		return errors.Wrap(err, "failed to create tag")
	}

	err = mbw.PostRelease(r, v)

	if err != nil {
		return errors.Wrap(err, "failed to post-release")
	}

	return nil
}

func (mbw *MultiBranchWorkflow) PostRelease(r *git.Repository, v version.Version) error {
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

	message := version.Replace(mbw.BumpMessageFormat, v)

	err = r.Commit(sig, message)

	if err != nil {
		return errors.Wrap(err, "failed to commit")
	}

	return nil
}
