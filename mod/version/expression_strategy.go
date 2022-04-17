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
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/chapterjason/j3n/modx/regexpx"
)

type ExpressionStrategy struct {
	Directories []string `json:"directories"`
	Pattern     string   `json:"pattern"`
	Expression  string   `json:"expression"`
	Replacement string   `json:"replacement"`
}

func (es *ExpressionStrategy) Log() string {
	return fmt.Sprintf(
		"Expression(%s -> %s): %s -> %s",
		strings.Join(es.Directories, ", "),
		es.Pattern,
		es.Expression,
		es.Replacement,
	)
}

func (es *ExpressionStrategy) GetExpression() *regexp.Regexp {
	return ReplaceExpression(es.Expression)
}

func (es *ExpressionStrategy) Get() ([]Version, error) {
	files, err := es.getFiles()

	if err != nil {
		return nil, err
	}

	versions := []Version{}

	for _, file := range files {
		version, err := es.getFile(file)

		if err != nil {
			return nil, err
		}

		versions = append(versions, version)
	}

	return versions, nil
}

func (es *ExpressionStrategy) Set(v Version) error {
	files, err := es.getFiles()

	if err != nil {
		return err
	}

	replacement := Replace(es.Replacement, v)

	for _, file := range files {
		err = es.setFile(file, replacement)

		if err != nil {
			return err
		}
	}

	return nil
}

func (es *ExpressionStrategy) getFiles() ([]string, error) {
	files := []string{}

	for _, directory := range es.Directories {
		pattern := path.Join(directory, es.Pattern)
		matches, err := filepath.Glob(pattern)

		if err != nil {
			return nil, err
		}

		files = append(files, matches...)
	}

	return files, nil
}

func (es *ExpressionStrategy) setFile(path string, replacement string) error {
	f, err := os.OpenFile(path, os.O_RDWR, 0644)

	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	b, err := ioutil.ReadAll(f)

	if err != nil {
		return err
	}

	expr := es.GetExpression()

	if expr.Match(b) {
		b = expr.ReplaceAll(b, []byte(replacement))
	}

	err = f.Truncate(0)

	if err != nil {
		return err
	}

	_, err = f.Seek(0, 0)

	if err != nil {
		return err
	}
	_, err = f.Write(b)

	if err != nil {
		return err
	}

	return nil
}

func (es *ExpressionStrategy) getFile(file string) (Version, error) {
	f, err := os.OpenFile(file, os.O_RDONLY, 0644)

	if err != nil {
		return Version{}, err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	b, err := ioutil.ReadAll(f)

	if err != nil {
		return Version{}, err
	}

	expr := es.GetExpression()

	if expr.Match(b) {
		mn, err := regexpx.MatchNamed(expr, b)

		if err != nil {
			return Version{}, err
		}

		vs, ok := mn["version"]

		if !ok {
			return Version{}, nil
		}

		parsed, err := Parse(vs)

		if err != nil {
			return Version{}, err
		}

		return parsed, nil
	}

	return Version{}, fmt.Errorf("no match found")
}
