package main

import (
	"bytes"
	"fmt"
	"testing"
)

func Benchmark_incFilename(b *testing.B) {
	b.ReportAllocs()

	var f flagParams
	for n := 0; n < b.N; n++ {
		f.incFilename(n)
	}
}

func Benchmark_lineCounter(b *testing.B) {
	b.ReportAllocs()

	var in *bytes.Buffer
	var tmpstr string

	for i := 0; i < 9; i++ {
		tmpstr = fmt.Sprintf("%vtest line %v\n", tmpstr, i)
	}
	for n := 0; n < b.N; n++ {
		in = bytes.NewBufferString(tmpstr)
		lineCounter(in)
	}
}

func Test_CLI(t *testing.T) {
	var f flagParams

	blankNameerr := "the source and destination file names cannot be blank"
	sameNameerr := "the source and destination files must be different"
	lineCountError := "the file cannot be split to less than 1 line per file"
	maxFileCountError := "maximum file count must be zero (maximum files) or greater"
	evenfilesplitError := "maximum file count cannot be zero to split file evenly over multiple files"

	var fileNameTests = []struct {
		sourcefile string
		destfile   string
		lineCount  int
		maxfiles   int
		evensplit  bool
		output     string
	}{
		{"", "destfiletest.tmp", 10, 0, false, blankNameerr},
		{"sourcefiletest.tmp", "destfiletest.tmp", 10, 0, false, ""},
		{"sourcefiletest.tmp", "", 10, 0, false, blankNameerr},
		{"", "", 10, 0, false, blankNameerr},
		{"sourcefiletest.tmp", "destfiletest.tmp", 10, 0, false, sameNameerr},
		{"sourcefiletest.tmp", "destfiletest.tmp", 100, 0, false, ""},
		{"sourcefiletest.tmp", "destfiletest.tmp", -1, 0, false, lineCountError},
		{"sourcefiletest.tmp", "destfiletest.tmp", 1, 1, false, ""},
		{"sourcefiletest.tmp", "destfiletest.tmp", 10, 1, false, ""},
		{"sourcefiletest.tmp", "destfiletest.tmp", 10, 0, false, ""},
		{"sourcefiletest.tmp", "destfiletest.tmp", 10, -1, false, maxFileCountError},
		{"sourcefiletest.tmp", "destfiletest.tmp", 0, 0, true, evenfilesplitError},
	}

	for _, tt := range fileNameTests {
		f.sourceFile = tt.sourcefile
		f.destFile = tt.destfile
		f.lineCount = tt.lineCount
		f.maxFiles = tt.maxfiles
		f.evenSplit = tt.evensplit
		got := f.checkFlagErrors()
		expected := tt.output

		if got != nil && got.Error() != expected {
			t.Errorf("Got: %v Expected: %v", got, expected)
		}
	}

}

func Test_incFileName(t *testing.T) {
	var f flagParams

	f.destFile = "MyFileName.txt"
	counter := 1

	got := f.incFilename(counter)
	want := "MyFileName1.txt"

	if got != want {
		t.Errorf("Got %v Expected %v", got, want)
	}

}

func Test_multiincFileName(t *testing.T) {
	var f flagParams

	f.destFile = "MyFileName.txt"

	var fileNameTests = []struct {
		counter  int
		expected string
		altexp   string
	}{
		{1, "MyFileName1.txt", "MyFileName1.txt"},
		{2, "MyFileName2.txt", "MyFileName2.txt"},
		{3, "MyFileName3.txt", "MyFileName3.txt"},
		{4, "MyFileName4.txt", "MyFileName4.txt"},
		{5, "MyFileName5.txt", "MyFileName5.txt"},
	}

	for _, tt := range fileNameTests {
		result := f.incFilename(tt.counter)
		if result != tt.expected && result != tt.altexp {
			t.Errorf("Counter: %v Expected: %v Result: %v", tt.counter, tt.expected, result)
		}
	}

}

func Test_fileExists(t *testing.T) {
	var tests = []struct {
		counter  int
		params   flagParams
		expected string
		altexp   string
	}{
		{1, flagParams{"unknown", 10, "unknown", 5, false, false, false}, "opening source file for reading: open %s: no such file or directory", "source file for reading: open %s: no such file or directory"},
		{2, flagParams{"testfile.txt", 1, "unknown", 1, false, false, false}, "opening source file for reading: open %s: no such file or directory", "source file for reading: open %s: no such file or directory"},
	}

	for _, tt := range tests {
		_, got := tt.params.splitFile()
		if got.Error() != fmt.Sprintf(tt.expected, tt.params.sourceFile) && got.Error() != fmt.Sprintf(tt.altexp, tt.params.sourceFile) {
			t.Errorf("Counter: %v \nExpected: %s\nResult  : %s", tt.counter, tt.expected, got)
		}
	}
}

func Test_lineCounter(t *testing.T) {

	var in *bytes.Buffer
	var tmpstr string
	var expected int

	for i := 0; i < 9; i++ {
		tmpstr = fmt.Sprintf("%vtest line %v\n", tmpstr, i)
		expected = i + 1
	}
	in = bytes.NewBufferString(tmpstr)
	got, _ := lineCounter(in)
	if got != expected+1 {
		t.Errorf("Counter: %v \nExpected: %v", got, expected)
		t.Fail()
	}

}
