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
	"strings"
	"time"
)

var (
	Strategies = []Strategy{}
)

type Strategy interface {
	Get() ([]Version, error)
	Set(Version) error
}

func preparePlacement(p string, v Version) string {
	// version
	p = strings.ReplaceAll(p, "{{VERSION_MAJOR}}", fmt.Sprintf("%d", v.Major))
	p = strings.ReplaceAll(p, "{{VERSION_MINOR}}", fmt.Sprintf("%d", v.Minor))
	p = strings.ReplaceAll(p, "{{VERSION_PATCH}}", fmt.Sprintf("%d", v.Patch))
	p = strings.ReplaceAll(p, "{{VERSION_PRERELEASE}}", strings.Join(v.Prerelease, "."))
	p = strings.ReplaceAll(p, "{{VERSION_BUILD}}", strings.Join(v.Build, "."))
	p = strings.ReplaceAll(p, "{{VERSION}}", v.String())

	// time
	now := time.Now()

	p = strings.ReplaceAll(p, "{{TIME_YEAR}}", fmt.Sprintf("%d", now.Year()))
	p = strings.ReplaceAll(p, "{{TIME_MONTH}}", fmt.Sprintf("%02d", int(now.Month())))
	p = strings.ReplaceAll(p, "{{TIME_DAY}}", fmt.Sprintf("%02d", now.Day()))
	p = strings.ReplaceAll(p, "{{TIME_HOUR}}", fmt.Sprintf("%02d", now.Hour()))
	p = strings.ReplaceAll(p, "{{TIME_MINUTE}}", fmt.Sprintf("%02d", now.Minute()))
	p = strings.ReplaceAll(p, "{{TIME_SECOND}}", fmt.Sprintf("%02d", now.Second()))
	p = strings.ReplaceAll(p, "{{TIME_NANOSECOND}}", fmt.Sprintf("%d", now.Nanosecond()))
	p = strings.ReplaceAll(p, "{{TIME_RFC3339}}", now.Format(time.RFC3339))
	p = strings.ReplaceAll(p, "{{TIME_RFC3339_NANO}}", now.Format(time.RFC3339Nano))

	return p
}
