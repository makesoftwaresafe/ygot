// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/openconfig/ygot/testutil"

	yb "github.com/openconfig/ygot/demo/uncompressed/pkg/demo"
)

const (
	// TestRoot is the path at which this test is running such that it
	// is possible to determine where to load test files from.
	TestRoot string = ""
)

func TestUncompressedDemo(t *testing.T) {
	u, err := BuildDemo()
	if err != nil {
		t.Fatalf("BuildDemo(): got unexpected error when generating, got: %v, want: nil", err)
	}

	gotInternal, err := DemoInternalJSON(u)
	if err != nil {
		t.Fatalf("DemoInternalJSON(%v): got unexpected error generating internal JSON, got: %v, want: nil", u, err)
	}

	gotRFC7951, err := DemoRFC7951JSON(u)
	if err != nil {
		t.Fatalf("DemoRFC7951JSON(%v): got unexpected error generating RFC7951 JSON, got: %v, want: nil", u, err)
	}

	tests := []struct {
		name     string
		got      string
		wantFile string
	}{{
		name:     "internal JSON",
		got:      gotInternal,
		wantFile: "testdata/demo.json",
	}, {
		name:     "RFC7951 JSON",
		got:      gotRFC7951,
		wantFile: "testdata/demo-ietf.json",
	}}

	for _, tt := range tests {
		want, ioerr := os.ReadFile(filepath.Join(TestRoot, tt.wantFile))
		if ioerr != nil {
			t.Errorf("%s: os.ReadFile(%s/%s): error reading file, got: %v, want: nil", tt.name, TestRoot, tt.wantFile, ioerr)
			continue
		}

		if diff := pretty.Compare(tt.got, string(want)); diff != "" {
			if diffl, err := testutil.GenerateUnifiedDiff(tt.got, string(want)); err != nil {
				diff = diffl
			}
			t.Errorf("%s: Demo JSON output of %v, did not get expected JSON, diff(-got,+want):\n%s", tt.name, u, diff)
		}
	}
}

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		in      []byte
		wantErr bool
	}{{
		name:    "unmarshal empty",
		in:      []byte(`{"married": [null]}`),
		wantErr: false,
	}}

	for _, tt := range tests {
		if err := yb.Unmarshal([]byte(`{"married": [null]}`), &yb.Root{}); err != nil {
			fmt.Printf("for test %s got error: %v", tt.name, err)
			if !tt.wantErr {
				t.Errorf("%s: did not get expected unmarshal error, got: %v, want: %v", tt.name, err != nil, tt.wantErr)
			}

		}
	}
}
