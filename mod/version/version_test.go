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
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Version
		wantErr bool
	}{
		{name: "valid", wantErr: false, input: "0.0.4", want: Version{Major: 0, Minor: 0, Patch: 4, Prerelease: []string{}, Build: []string{}}},
		{name: "valid", wantErr: false, input: "1.2.3", want: Version{Major: 1, Minor: 2, Patch: 3, Prerelease: []string{}, Build: []string{}}},
		{name: "valid", wantErr: false, input: "10.20.30", want: Version{Major: 10, Minor: 20, Patch: 30, Prerelease: []string{}, Build: []string{}}},
		{name: "valid", wantErr: false, input: "1.1.2+meta", want: Version{Major: 1, Minor: 1, Patch: 2, Prerelease: []string{}, Build: []string{"meta"}}},
		{name: "valid", wantErr: false, input: "1.1.2-prerelease+meta", want: Version{Major: 1, Minor: 1, Patch: 2, Prerelease: []string{"prerelease"}, Build: []string{"meta"}}},
		{name: "valid", wantErr: false, input: "1.1.2+meta-valid", want: Version{Major: 1, Minor: 1, Patch: 2, Prerelease: []string{}, Build: []string{"meta-valid"}}},
		{name: "valid", wantErr: false, input: "1.0.0-alpha", want: Version{Major: 1, Minor: 0, Patch: 0, Prerelease: []string{"alpha"}, Build: []string{}}},
		{name: "valid", wantErr: false, input: "1.0.0-beta", want: Version{Major: 1, Minor: 0, Patch: 0, Prerelease: []string{"beta"}, Build: []string{}}},
		{name: "valid", wantErr: false, input: "1.0.0-alpha.beta", want: Version{Major: 1, Minor: 0, Patch: 0, Prerelease: []string{"alpha", "beta"}, Build: []string{}}},
		{name: "valid", wantErr: false, input: "1.0.0-alpha.beta.1", want: Version{Major: 1, Minor: 0, Patch: 0, Prerelease: []string{"alpha", "beta", "1"}, Build: []string{}}},
		{name: "valid", wantErr: false, input: "1.0.0-alpha.1", want: Version{Major: 1, Minor: 0, Patch: 0, Prerelease: []string{"alpha", "1"}, Build: []string{}}},
		{name: "valid", wantErr: false, input: "1.0.0-alpha0.valid", want: Version{Major: 1, Minor: 0, Patch: 0, Prerelease: []string{"alpha0", "valid"}, Build: []string{}}},
		{name: "valid", wantErr: false, input: "1.0.0-alpha.0valid", want: Version{Major: 1, Minor: 0, Patch: 0, Prerelease: []string{"alpha", "0valid"}, Build: []string{}}},
		{name: "valid", wantErr: false, input: "1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay", want: Version{Major: 1, Minor: 0, Patch: 0, Prerelease: []string{"alpha-a", "b-c-somethinglong"}, Build: []string{"build", "1-aef", "1-its-okay"}}},
		{name: "valid", wantErr: false, input: "1.0.0-rc.1+build.1", want: Version{Major: 1, Minor: 0, Patch: 0, Prerelease: []string{"rc", "1"}, Build: []string{"build", "1"}}},
		{name: "valid", wantErr: false, input: "2.0.0-rc.1+build.123", want: Version{Major: 2, Minor: 0, Patch: 0, Prerelease: []string{"rc", "1"}, Build: []string{"build", "123"}}},
		{name: "valid", wantErr: false, input: "1.2.3-beta", want: Version{Major: 1, Minor: 2, Patch: 3, Prerelease: []string{"beta"}, Build: []string{}}},
		{name: "valid", wantErr: false, input: "10.2.3-DEV-SNAPSHOT", want: Version{Major: 10, Minor: 2, Patch: 3, Prerelease: []string{"DEV-SNAPSHOT"}, Build: []string{}}},
		{name: "valid", wantErr: false, input: "1.2.3-SNAPSHOT-123", want: Version{Major: 1, Minor: 2, Patch: 3, Prerelease: []string{"SNAPSHOT-123"}, Build: []string{}}},
		{name: "valid", wantErr: false, input: "1.0.0", want: Version{Major: 1, Minor: 0, Patch: 0, Prerelease: []string{}, Build: []string{}}},
		{name: "valid", wantErr: false, input: "2.0.0", want: Version{Major: 2, Minor: 0, Patch: 0, Prerelease: []string{}, Build: []string{}}},
		{name: "valid", wantErr: false, input: "1.1.7", want: Version{Major: 1, Minor: 1, Patch: 7, Prerelease: []string{}, Build: []string{}}},
		{name: "valid", wantErr: false, input: "2.0.0+build.1848", want: Version{Major: 2, Minor: 0, Patch: 0, Prerelease: []string{}, Build: []string{"build", "1848"}}},
		{name: "valid", wantErr: false, input: "2.0.1-alpha.1227", want: Version{Major: 2, Minor: 0, Patch: 1, Prerelease: []string{"alpha", "1227"}, Build: []string{}}},
		{name: "valid", wantErr: false, input: "1.0.0-alpha+beta", want: Version{Major: 1, Minor: 0, Patch: 0, Prerelease: []string{"alpha"}, Build: []string{"beta"}}},
		{name: "valid", wantErr: false, input: "1.2.3----RC-SNAPSHOT.12.9.1--.12+788", want: Version{Major: 1, Minor: 2, Patch: 3, Prerelease: []string{"---RC-SNAPSHOT", "12", "9", "1--", "12"}, Build: []string{"788"}}},
		{name: "valid", wantErr: false, input: "1.2.3----R-S.12.9.1--.12+meta", want: Version{Major: 1, Minor: 2, Patch: 3, Prerelease: []string{"---R-S", "12", "9", "1--", "12"}, Build: []string{"meta"}}},
		{name: "valid", wantErr: false, input: "1.2.3----RC-SNAPSHOT.12.9.1--.12", want: Version{Major: 1, Minor: 2, Patch: 3, Prerelease: []string{"---RC-SNAPSHOT", "12", "9", "1--", "12"}, Build: []string{}}},
		{name: "valid", wantErr: false, input: "1.0.0+0.build.1-rc.10000aaa-kk-0.1", want: Version{Major: 1, Minor: 0, Patch: 0, Prerelease: []string{}, Build: []string{"0", "build", "1-rc", "10000aaa-kk-0", "1"}}},
		{name: "valid", wantErr: false, input: "9999999999999999999.9999999999999999999.9999999999999999999", want: Version{Major: 9999999999999999999, Minor: 9999999999999999999, Patch: 9999999999999999999, Prerelease: []string{}, Build: []string{}}},
		{name: "valid", wantErr: false, input: "1.0.0-0A.is.legal", want: Version{Major: 1, Minor: 0, Patch: 0, Prerelease: []string{"0A", "is", "legal"}, Build: []string{}}},
		{name: "empty", input: "", wantErr: true},
		{name: "invalid", input: "1", wantErr: true},
		{name: "invalid", input: "1.2", wantErr: true},
		{name: "invalid", input: "1.2.3-0123", wantErr: true},
		{name: "invalid", input: "1.2.3-0123.0123", wantErr: true},
		{name: "invalid", input: "1.1.2+.123", wantErr: true},
		{name: "invalid", input: "+invalid", wantErr: true},
		{name: "invalid", input: "-invalid", wantErr: true},
		{name: "invalid", input: "-invalid+invalid", wantErr: true},
		{name: "invalid", input: "-invalid.01", wantErr: true},
		{name: "invalid", input: "alpha", wantErr: true},
		{name: "invalid", input: "alpha.beta", wantErr: true},
		{name: "invalid", input: "alpha.beta.1", wantErr: true},
		{name: "invalid", input: "alpha.1", wantErr: true},
		{name: "invalid", input: "alpha+beta", wantErr: true},
		{name: "invalid", input: "alpha_beta", wantErr: true},
		{name: "invalid", input: "alpha.", wantErr: true},
		{name: "invalid", input: "alpha..", wantErr: true},
		{name: "invalid", input: "beta", wantErr: true},
		{name: "invalid", input: "1.0.0-alpha_beta", wantErr: true},
		{name: "invalid", input: "-alpha.", wantErr: true},
		{name: "invalid", input: "1.0.0-alpha..", wantErr: true},
		{name: "invalid", input: "1.0.0-alpha..1", wantErr: true},
		{name: "invalid", input: "1.0.0-alpha...1", wantErr: true},
		{name: "invalid", input: "1.0.0-alpha....1", wantErr: true},
		{name: "invalid", input: "1.0.0-alpha.....1", wantErr: true},
		{name: "invalid", input: "1.0.0-alpha......1", wantErr: true},
		{name: "invalid", input: "1.0.0-alpha.......1", wantErr: true},
		{name: "invalid", input: "01.1.1", wantErr: true},
		{name: "invalid", input: "1.01.1", wantErr: true},
		{name: "invalid", input: "1.1.01", wantErr: true},
		{name: "invalid", input: "1.2", wantErr: true},
		{name: "invalid", input: "1.2.3.DEV", wantErr: true},
		{name: "invalid", input: "1.2-SNAPSHOT", wantErr: true},
		{name: "invalid", input: "1.2.31.2.3----RC-SNAPSHOT.12.09.1--..12+788", wantErr: true},
		{name: "invalid", input: "1.2-RC-SNAPSHOT", wantErr: true},
		{name: "invalid", input: "-1.0.3-gamma+b7718", wantErr: true},
		{name: "invalid", input: "+justmeta", wantErr: true},
		{name: "invalid", input: "9.8.7+meta+meta", wantErr: true},
		{name: "invalid", input: "9.8.7-whatever+meta+meta", wantErr: true},
		{name: "invalid", input: "99999999999999999999999.999999999999999999.99999999999999999----RC-SNAPSHOT.12.09.1--------------------------------..12", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := Parse(tt.input)

				if (err != nil) != tt.wantErr {
					t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Parse() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
