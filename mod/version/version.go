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
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/chapterjason/j3n/modx/regexpx"
	"github.com/chapterjason/j3n/modx/slicex"
	"github.com/chapterjason/j3n/modx/strconvx"
)

var (
	ErrInvalidVersion = errors.New("invalid version")

	SemverExpression        = regexp.MustCompile("(?P<major>0|[1-9]\\d*)\\.(?P<minor>0|[1-9]\\d*)\\.(?P<patch>0|[1-9]\\d*)(?:-(?P<prerelease>(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?")
	SemverExpressionPartial = regexp.MustCompile("(?P<versiob>(?P<major>0|[1-9]\\d*)\\.(?P<minor>0|[1-9]\\d*)\\.(?P<patch>0|[1-9]\\d*)(?:-(?P<prerelease>(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*)))?")
	SemverExpressionFull    = regexp.MustCompile("^" + SemverExpression.String() + "$")
)

type Version struct {
	Major uint64
	Minor uint64

	Patch      uint64
	Prerelease []string
	Build      []string
}

func (v Version) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, v.String())), nil
}

func (v *Version) UnmarshalJSON(bytes []byte) error {
	parsed, err := Parse(strings.Trim(string(bytes), `"`))

	if err != nil {
		return err
	}

	v.Major = parsed.Major
	v.Minor = parsed.Minor
	v.Patch = parsed.Patch
	v.Prerelease = parsed.Prerelease
	v.Build = parsed.Build

	return nil
}

func Parse[T string | []byte](text T) (Version, error) {
	match, err := regexpx.MatchNamed(SemverExpressionFull, text)

	if err != nil {
		return Version{}, errors.Wrap(err, "failed to parse version")
	}

	return Version{
		Major:      strconvx.MustParseUint(match["major"]),
		Minor:      strconvx.MustParseUint(match["minor"]),
		Patch:      strconvx.MustParseUint(match["patch"]),
		Prerelease: slicex.FilterEmpty(strings.Split(match["prerelease"], ".")),
		Build:      slicex.FilterEmpty(strings.Split(match["buildmetadata"], ".")),
	}, nil
}

func MustParse(text string) Version {
	v, err := Parse(text)

	if err != nil {
		panic(err)
	}

	return v
}

func (v Version) String() string {
	p := ""
	b := ""

	if len(v.Prerelease) > 0 {
		p = "-" + strings.Join(v.Prerelease, ".")
	}

	if len(v.Build) > 0 {
		b = "+" + strings.Join(v.Build, ".")
	}

	return fmt.Sprintf("%d.%d.%d%s%s", v.Major, v.Minor, v.Patch, p, b)
}
