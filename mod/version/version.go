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
	"fmt"
	"math"
	"regexp"
	"strconv"
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
	return json.Marshal(v.String())
}

func (v *Version) UnmarshalJSON(bytes []byte) error {
	var s string

	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}

	parsed, err := Parse(s)

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
		return Version{}, ErrInvalidVersion
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

func (v Version) Compare(v2 Version) int {
	if v.Major > v2.Major {
		return 1
	} else if v.Major < v2.Major {
		return -1
	}

	if v.Minor > v2.Minor {
		return 1
	} else if v.Minor < v2.Minor {
		return -1
	}

	if v.Patch > v2.Patch {
		return 1
	} else if v.Patch < v2.Patch {
		return -1
	}

	return comparePrerelease(v, v2)
}

func comparePrerelease(v1 Version, v2 Version) int {
	v1p := v1.Prerelease
	v2p := v2.Prerelease

	v1pl := len(v1p)
	v2pl := len(v2p)

	if v1pl == 0 && v2pl != 0 {
		return 1
	} else if v1pl != 0 && v2pl == 0 {
		return -1
	}

	l := int(math.Max(float64(v1pl), float64(v2pl)))

	for i := 0; i < l; i++ {
		v1pp := ""

		if i < len(v1p) {
			v1pp = v1p[i]
		}

		v2pp := ""

		if i < len(v2p) {
			v2pp = v2p[i]
		}

		if cv := comparePrereleasePart(v1pp, v2pp); cv != 0 {
			return cv
		}
	}

	return 0
}

func comparePrereleasePart(v1pp string, v2pp string) int {
	if v1pp == v2pp {
		return 0
	}

	// 4. A larger set of pre-release fields has a higher precedence than a smaller set, if all of the preceding identifiers are equal.

	if v1pp == "" {
		if v2pp != "" {
			return -1
		}

		return 1
	}

	if v2pp == "" {
		if v1pp != "" {
			return 1
		}

		return -1
	}

	sppi, spperr := strconv.ParseUint(v1pp, 10, 64)
	oppi, opperr := strconv.ParseUint(v2pp, 10, 64)

	// 1. Identifiers consisting of only digits are compared numerically.
	if spperr == nil && opperr == nil {
		if sppi > oppi {
			return 1
		}

		return -1
	}

	// 2. Identifiers with letters or hyphens are compared lexically in ASCII sort order.
	if spperr != nil && opperr != nil {

		if v1pp > v2pp {
			return 1
		}

		return -1
	}

	// 3. Numeric identifiers always have lower precedence than non-numeric identifiers.
	if opperr != nil {
		return -1
	}

	// spperr != nil
	return 1
}
