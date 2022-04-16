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
	"github.com/spf13/viper"
)

type VersionStrategy struct {
	directory string
}

func NewVersionStrategy(directory string) *VersionStrategy {
	return &VersionStrategy{
		directory,
	}
}

func (vs *VersionStrategy) Get() ([]Version, error) {
	vst := viper.GetString("version")

	v, err := Parse(vst)

	if err != nil {
		return nil, err
	}

	return []Version{v}, nil
}

func (vs *VersionStrategy) Set(v Version) error {
	viper.Set("version", v.String())

	err := viper.SafeWriteConfig()

	if err != nil {
		return err
	}

	return nil
}
