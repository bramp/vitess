/*
Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sqlparser

import (
	"fmt"
	"testing"
)

var (
	// Simple mapping of token IDs to names. This should really
	// be provided by the goyacc output, but it isn't.
	// https://github.com/golang/go/issues/6577
	tokNames = map[int]string{
		STRING:    "STRING",
		ID:        "ID",
		0:         "EOF",
		LEX_ERROR: "LEX_ERROR",
	}
)

func tokenName(id int) string {
	if name, found := tokNames[id]; found {
		return name
	}

	return fmt.Sprintf("%d", id)
}

func TestID(t *testing.T) {
	testcases := []struct {
		in     string
		id     int
		out    string
		nextId int
	}{{
		in:  "a",
		id:  ID,
		out: "a",
	}, {
		in:  "abc",
		id:  ID,
		out: "abc",
	}, {
		in:  "a1b2",
		id:  ID,
		out: "a1b2",
	}, {
		in:  "`aa`",
		id:  ID,
		out: "aa",
	}, {
		in:  "```a```",
		id:  ID,
		out: "`a`",
	}, {
		in:  "`a``b`",
		id:  ID,
		out: "a`b",
	}, {
		in:     "`a``b`c",
		id:     ID,
		out:    "a`b",
		nextId: ID,
	}, {
		in:  "`a``b",
		id:  LEX_ERROR,
		out: "a`b",
	}, {
		in:  "`a``b``",
		id:  LEX_ERROR,
		out: "a`b`",
	}, {
		in:  "``",
		id:  LEX_ERROR,
		out: "",
	}}

	for _, tcase := range testcases {
		tkn := NewStringTokenizer(tcase.in)
		id, out := tkn.Scan()
		if tcase.id != id || string(out) != tcase.out {
			t.Errorf("Scan(%q) = (%s, %q), want (%s, %q)", tcase.in, tokenName(id), out, tokenName(tcase.id), tcase.out)
		}
		if id, out = tkn.Scan(); id != tcase.nextId {
			t.Errorf("Scan() = (%s, %q), want (%s, ?)", tokenName(id), out, tokenName(tcase.nextId))
		}
	}
}

func TestString(t *testing.T) {
	testcases := []struct {
		in     string
		id     int
		want   string
		nextId int
	}{{
		in:   "''",
		id:   STRING,
		want: "",
	}, {
		in:   "''''",
		id:   STRING,
		want: "'",
	}, {
		in:   "'hello'",
		id:   STRING,
		want: "hello",
	}, {
		in:   "'\\n'",
		id:   STRING,
		want: "\n",
	}, {
		in:   "'\\nhello\\n'",
		id:   STRING,
		want: "\nhello\n",
	}, {
		in:   "'a''b'",
		id:   STRING,
		want: "a'b",
	}, {
		in:   "'a\\'b'",
		id:   STRING,
		want: "a'b",
	}, {
		in:   "'\\'",
		id:   LEX_ERROR,
		want: "'",
	}, {
		in:   "'",
		id:   LEX_ERROR,
		want: "",
	}, {
		in:   "'hello\\'",
		id:   LEX_ERROR,
		want: "hello'",
	}, {
		in:   "'hello",
		id:   LEX_ERROR,
		want: "hello",
	}, {
		in:   "'hello\\",
		id:   LEX_ERROR,
		want: "hello",
	}}

	for _, tcase := range testcases {
		tkn := NewStringTokenizer(tcase.in)
		id, got := tkn.Scan()
		if tcase.id != id || string(got) != tcase.want {
			t.Errorf("Scan(%q) = (%s, %q), want (%s, %q)", tcase.in, tokenName(id), got, tokenName(tcase.id), tcase.want)
		}
		if id, got = tkn.Scan(); id != tcase.nextId {
			t.Errorf("Scan() = (%s, %q), want (%s, ?)", tokenName(id), got, tokenName(tcase.nextId))
		}

	}
}
