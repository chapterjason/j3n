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

package version

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

type packageJson struct {
	Version string `json:"version"`
}

type NpmStrategy struct {
	directory string
}

func NewNpmStrategy(directory string) *NpmStrategy {
	return &NpmStrategy{
		directory,
	}
}

func (ns *NpmStrategy) Get() ([]Version, error) {
	packagePath := ns.directory + "/package.json"

	file, err := os.OpenFile(packagePath, os.O_RDONLY, 0644)

	if err != nil {
		return nil, errors.Wrap(err, "failed to open package.json")
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, errors.Wrap(err, "failed to read contents of package.json")
	}

	var pkg packageJson

	err = json.Unmarshal(bytes, &pkg)

	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal package.json")
	}

	parsed, err := Parse(pkg.Version)

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse version from package.json")
	}

	return []Version{parsed}, nil
}

func (ns *NpmStrategy) Set(v Version) error {
	ok, err := ns.trySetYarn(v)

	if err != nil {
		return errors.Wrap(err, "failed to set version")
	}

	if !ok {
		return ns.trySetNpm(v)
	}

	return nil
}

func (ns *NpmStrategy) trySetYarn(v Version) (bool, error) {
	yarn, err := exec.LookPath("yarn")

	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return false, nil
		}

		return false, errors.Wrap(err, "failed to find yarn")
	}

	cmd := exec.Command(yarn, "version", "--no-git-tag-version", "--new-version", v.String())
	cmd.Dir = ns.directory

	if err := cmd.Run(); err != nil {
		return false, errors.Wrap(err, "failed to run yarn")
	}

	return true, nil
}

func (ns *NpmStrategy) trySetNpm(v Version) error {
	npm, err := exec.LookPath("npm")

	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return errors.Wrap(err, "strategy requires npm or yarn to be used")
		}

		return err
	}

	cmd := exec.Command(npm, "version", "--no-git-tag-version", v.String())
	cmd.Dir = ns.directory

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed running npm")
	}

	return nil
}
