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
	"encoding/json"
	"fmt"

	"github.com/chapterjason/j3n/modx/viperx"
)

type Config struct {
	Workflow Workflow `json:"workflow"`
}

func (d *Config) UnmarshalJSON(data []byte) error {
	var typ struct {
		Workflow map[string]any `json:"workflow"`
	}

	if err := json.Unmarshal(data, &typ); err != nil {
		return err
	}

	tp, ok := typ.Workflow["type"]

	if !ok {
		return fmt.Errorf("missing workflow type")
	}

	switch tp {
	case "multi_branch":
		d.Workflow = &MultiBranchWorkflow{
			BranchFormat:        "release/{{VERSION_MAJOR}}.{{VERSION_MINOR}}",
			BranchExpression:    "release/\\d+\\.\\d+",
			UpdateMessageFormat: "Update version for {{VERSION}}",
			BumpMessageFormat:   "Bump version to {{VERSION}}",
		}
	}

	return viperx.Transcode(typ.Workflow, d.Workflow)

}
