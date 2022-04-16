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

	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"

	"github.com/chapterjason/j3n/mod/version"
)

var (
	ErrAlreadyReleased = errors.New("already released")
)

func Release(r *git.Repository, v version.Version, wf Workflow) error {
	has, err := HasTag(r, v)

	if err != nil {
		return errors.Wrap(err, "failed to check if tag exists")
	}

	if has {
		return errors.Wrap(ErrAlreadyReleased, fmt.Sprintf("tag %s already exists", v))
	}

	if v.Patch != 0 {
		return doRelease(r, v, wf, ReleaseTypePatch)
	} else if v.Minor != 0 {
		return doRelease(r, v, wf, ReleaseTypeMinor)
	} else {
		return doRelease(r, v, wf, ReleaseTypeMajor)
	}
}

func doRelease(r *git.Repository, v version.Version, wf Workflow, rt ReleaseType) error {
	err := wf.PreRelease(r, v, rt)

	if err != nil {
		return err
	}

	h, err := r.Head()

	if err != nil {
		return err
	}

	tn := GitTagFormatter(v)

	_, err = r.CreateTag(tn, h.Hash(), nil)

	if err != nil {
		return err
	}

	err = wf.PostRelease(r, v, rt)

	if err != nil {
		return err
	}

	return nil
}
