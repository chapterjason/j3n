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

func Release(r *git.Repository, v version.Version) error {
	has, err := HasTag(r, v)

	if err != nil {
		return errors.Wrap(err, "failed to check if tag exists")
	}

	if has {
		return errors.Wrap(ErrAlreadyReleased, fmt.Sprintf("tag %s already exists", v))
	}

	if v.Patch != 0 {
		return Patch(r, v)
	} else if v.Minor != 0 {
		return Minor(r, v)
	} else {
		return Major(r, v)
	}

	return nil
}
