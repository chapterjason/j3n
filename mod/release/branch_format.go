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
	"regexp"

	"github.com/chapterjason/j3n/mod/version"
	"github.com/chapterjason/j3n/modx/regexpx"
	"github.com/chapterjason/j3n/modx/strconvx"
)

var (
	BranchFormat     = "{{VERSION_MAJOR}}.{{VERSION_MINOR}}"
	BranchExpression = "^(?<major>\\d+)\\.(?P<minor>\\d+)$"
	BranchFormatter  = func(v version.Version) string {
		return version.Replace(BranchFormat, v)
	}
	BranchMatcher = func(s string) bool {
		return regexp.MustCompile(BranchExpression).MatchString(s)
	}
	BranchVersionExtractor = func(s string) (version.Version, error) {
		expr := regexp.MustCompile(BranchExpression)
		groups, err := regexpx.MatchNamed(expr, s)

		if err != nil {
			return version.Version{}, err
		}

		return version.Version{
			Major: strconvx.MustParseUint(groups["major"]),
			Minor: strconvx.MustParseUint(groups["minor"]),
			Patch: 0,
		}, nil
	}
)
