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
	"time"

	"github.com/gogs/git-module"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/chapterjason/j3n/mod/version"
	"github.com/chapterjason/j3n/modx/gitx"
)

func Next(r *git.Repository, cv version.Version, nv version.Version, workflow Workflow) error {
	cmd := git.NewCommand("status", "--porcelain")

	b, err := cmd.RunWithTimeout(time.Duration(0))

	if err != nil {
		return err
	}

	if len(b) > 0 {
		return errors.New("uncommitted changes")
	}

	log.Infof("Current version: %s", cv)
	log.Infof("Next version: %s", nv)

	bbs := workflow.GetBranch(cv)
	rbs := workflow.GetBranch(nv)

	if !r.HasBranch(bbs) {
		return fmt.Errorf("missing base branch: %s", bbs)
	}

	if r.HasBranch(rbs) {
		return fmt.Errorf("branch already exists: %s", rbs)
	}

	err = r.Checkout(
		rbs, git.CheckoutOptions{
			BaseBranch: bbs,
		},
	)

	if err != nil {
		return err
	}

	log.Infof("Release branch created: %s", rbs)

	err = version.Set(nv)

	if err != nil {
		return err
	}

	log.Infof("Version set to: %s", nv)

	s, err := gitx.GetSignature(r)

	if err != nil {
		return err
	}

	err = r.Add(git.AddOptions{All: true})

	if err != nil {
		return err
	}

	message := version.Replace(workflow.GetBumpMessageFormat(), nv)

	err = r.Commit(s, message)

	if err != nil {
		return err
	}

	log.Infof("Bumped version: %s", nv)

	return nil
}
