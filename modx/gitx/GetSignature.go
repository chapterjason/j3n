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

package gitx

import (
	"time"

	"github.com/gogs/git-module"
)

func GetSignature(r *git.Repository) (*git.Signature, error) {
	sig := &git.Signature{}

	cmd := git.NewCommand("config", "--get", "user.name")

	b, err := cmd.RunInDirWithTimeout(time.Duration(0), r.Path())

	if err != nil {
		return nil, err
	}

	sig.Name = string(b)

	cmd = git.NewCommand("config", "--get", "user.email")

	b, err = cmd.RunInDirWithTimeout(time.Duration(0), r.Path())

	if err != nil {
		return nil, err
	}

	sig.Email = string(b)

	return sig, nil
}
