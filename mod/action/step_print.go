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

package action

import (
	"fmt"
	"os"
)

func init() {
	MustRegister(
		"print", func(input any, params map[string]any) (any, error) {
			out := os.Stdout

			if params["stream"] == "stderr" {
				out = os.Stderr
			}

			if input == nil {
				return nil, fmt.Errorf("input is nil")
			}

			switch input.(type) {
			case string:
				fmt.Fprint(out, input)
			case []byte:
				fmt.Fprintf(out, "%s", input)
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				fmt.Fprintf(out, "%d", input)
			case float32, float64:
				fmt.Fprintf(out, "%f", input)
			case bool:
				fmt.Fprintf(out, "%t", input)
			default:
				return nil, fmt.Errorf("unsupported type %T", input)
			}

			return nil, nil
		},
	)
}
