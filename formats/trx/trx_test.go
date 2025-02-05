package trx

import (
	"encoding/json"
	"flag"
	"github.com/google/go-cmp/cmp"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type testcase struct {
	Filename   string
	XFail      bool
	GoldenFile string
}

/* https://github.com/jstemmer/go-junit-report/blob/master/go-junit-report_test.go */

var testset = []testcase{
	{
		Filename:   "trx-sample-1.trx",
		XFail:      false,
		GoldenFile: "trx-sample-1.json",
	},
	{
		Filename:   "gh-596.trx",
		XFail:      false,
		GoldenFile: "gh-596.json",
	},
	{
		Filename:   "gh-749.trx",
		XFail:      false,
		GoldenFile: "gh-749.json",
	},
	{
		Filename:   "sample.trx",
		XFail:      false,
		GoldenFile: "sample.json",
	},
}

var update = flag.Bool("update", false, "update .golden files")

func TestParser(t *testing.T) {
	for _, test := range testset {
		if test.XFail {
			t.Logf("XFail: %s", test.Filename)
			continue
		}
		t.Logf("Running: %s", test.Filename)

		file, err := os.Open("testdata/" + test.Filename)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("Match structs", test.Filename)

		report, err := NewParser(file)
		if err != nil {
			t.Fatalf("Error parsing %s: %s", test.Filename, err)
		}

		goldenPath := filepath.Join("testdata", test.GoldenFile)
		if *update {
			t.Log("Update", goldenPath)
			data, _ := json.MarshalIndent(report, "", " ")
			if err := ioutil.WriteFile(goldenPath, data, 0644); err != nil {
				t.Fatalf("Failed to update golden file: %s", err)
			}
		}

		rawJson, err := ioutil.ReadFile(goldenPath)
		if err != nil {
			t.Fatal(err.Error())
		}
		var goldenReport TRXReport
		json.Unmarshal(rawJson, &goldenReport)
		if diff := cmp.Diff(*report, goldenReport); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}
