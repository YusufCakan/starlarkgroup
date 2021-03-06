// Copyright 2020 Edward McFarlane. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package starlarkgroup

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarktest"
)

func load(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	if module == "assert.star" {
		return starlarktest.LoadAssertModule()
	}
	return nil, fmt.Errorf("unknown module %s", module)
}

func TestExecFile(t *testing.T) {
	thread := &starlark.Thread{Load: load}
	starlarktest.SetReporter(thread, t)
	globals := starlark.StringDict{
		"group": starlark.NewBuiltin("group", Make),
	}

	files, err := filepath.Glob("testdata/*.star")
	if err != nil {
		t.Fatal(err)
	}

	for _, filename := range files {
		src, err := ioutil.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}

		_, err = starlark.ExecFile(thread, filename, src, globals)
		switch err := err.(type) {
		case *starlark.EvalError:
			var found bool
			for i := range err.CallStack {
				posn := err.CallStack.At(i).Pos
				if posn.Filename() == filename {
					linenum := int(posn.Line)
					msg := err.Error()

					t.Errorf("\n%s:%d: unexpected error: %v", filename, linenum, msg)
					found = true
					break
				}
			}
			if !found {
				t.Error(err.Backtrace())
			}
		case nil:
			// success
		default:
			t.Errorf("\n%s", err)
		}

	}
}
