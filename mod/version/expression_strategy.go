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
	directories []string
	pattern     string
	expression  *regexp.Regexp
	replacement string
}

func NewExpressionStrategy(directories []string, pattern, expression, replacement string) *ExpressionStrategy {
	return &ExpressionStrategy{
		directories: directories,
		pattern:     pattern,
		expression:  regexp.MustCompile(strings.ReplaceAll(expression, "{{VERSION}}", SemverExpressionPartial.String())),
		replacement: replacement,
	}
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

	replacement := preparePlacement(es.replacement, v)

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

	for _, directory := range es.directories {
		pattern := path.Join(directory, es.pattern)
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

	if es.expression.Match(b) {
		b = es.expression.ReplaceAll(b, []byte(replacement))
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

	if es.expression.Match(b) {
		mn, err := regexpx.MatchNamed(es.expression, b)

		if err != nil {
			return Version{}, err
		}

		parsed, err := Parse(mn["version"])

		if err != nil {
			return Version{}, err
		}

		return parsed, nil
	}

	return Version{}, fmt.Errorf("no match found")
}
