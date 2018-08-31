package main

import (
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

func Test_sourceFileName(t *testing.T) {
	var f flagParams

	blankNameerr := "the source and destination file names cannot be blank"
	sameNameerr := "the source and destination files must be different"

	var fileNameTests = []struct {
		sourcefile string
		destfile   string
		output     string
	}{
		{"", "destfiletest.tmp", blankNameerr},
		{"sourcefiletest.tmp", "destfiletest.tmp", ""},
		{"sourcefiletest.tmp", "", blankNameerr},
		{"", "", blankNameerr},
		{"sourcefiletest.tmp", "sourcefiletest.tmp", sameNameerr},
	}

	for _, tt := range fileNameTests {
		f.sourceFile = tt.sourcefile
		f.destFile = tt.destfile
		got := f.checkFileParams()
		expected := tt.output

		if got != nil && got.Error() != expected {
			t.Errorf("Got: %v Expected: %v", got, expected)
		}
	}

}

func Test_splitFileLength(t *testing.T) {
	var f flagParams

	lineCountError := "the file cannot be split to less than 1 line per file"

	var lineCountTests = []struct {
		lineCount int
		output    string
	}{
		{100, ""},
		{0, lineCountError},
		{1, ""},
	}

	for _, tt := range lineCountTests {
		f.lineCount = tt.lineCount
		got := f.checkLineCount()
		expected := tt.output

		if got != nil && got.Error() != expected {
			t.Errorf("Got: %v Expected: %v", got, expected)
		}
	}

}

func Test_maxFileCount(t *testing.T) {
	var f flagParams

	maxFileCountError := "maximum file count must be zero (maximum files) or greater"

	var maxFileCountTests = []struct {
		maxFileCount int
		output       string
	}{
		{1, ""},
		{0, ""},
		{-1, maxFileCountError},
	}

	for _, tt := range maxFileCountTests {
		f.maxFiles = tt.maxFileCount
		got := f.checkFileCount()
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
	}{
		{1, "MyFileName1.txt"},
		{2, "MyFileName2.txt"},
		{3, "MyFileName3.txt"},
		{4, "MyFileName4.txt"},
		{5, "MyFileName5.txt"},
	}

	for _, tt := range fileNameTests {
		result := f.incFilename(tt.counter)
		if result != tt.expected {
			t.Errorf("Counter: %v Expected: %v Result: %v", tt.counter, tt.expected, result)
		}
	}

}

func Test_fileExists(t *testing.T) {
	var tests = []struct {
		counter  int
		params   flagParams
		expected string
	}{
		{1, flagParams{"unknown", 10, "unknown", 5}, "open %s: The system cannot find the file specified."},
		{2, flagParams{"testfile.txt", 1, "unknown", 1}, "open %s: The system cannot find the file specified."},
	}

	for _, tt := range tests {
		_, got := tt.params.splitFile()
		if got.Error() != fmt.Sprintf(tt.expected, tt.params.sourceFile) {
			t.Errorf("Counter: %v \nExpected: %s\nResult  : %s", tt.counter, tt.expected, got)
		}
	}
}
