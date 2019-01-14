package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
)

type flagParams struct {
	sourceFile string
	lineCount  int
	destFile   string
	maxFiles   int
	evenSplit  bool
	countLines bool
	repeatHeader bool
}

var params flagParams

func init() {
	flag.StringVar(&params.sourceFile, "i", "", "File to be split")
	flag.IntVar(&params.lineCount, "l", 0, "maximum lines files to be split into")
	flag.StringVar(&params.destFile, "o", "", "Destination file name")
	flag.IntVar(&params.maxFiles, "m", 0, "Maximum number of files to be output. (0 for all)")
	flag.BoolVar(&params.evenSplit, "e", false, "Split the file evenly over the maximum number of files")
	flag.BoolVar(&params.repeatHeader, "h", false, "Repeat the header line of the first file in every file output.")
	flag.BoolVar(&params.countLines, "c", false, "Count the number of lines in the file")
}

func main() {

	flag.Parse()

	// Check for errors with the cli parameters
	err := params.checkFlagErrors()
	if err != nil {
		flag.PrintDefaults()
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Source file: %v \n", params.sourceFile)

	if !params.countLines {
		fmt.Printf("Destination file: %v \n", params.destFile)

		//fileCount, err := params.splitFile(params.countLines)
		fileCount, err := params.splitFile()
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}

		fmt.Printf("Split complete. %v files.\n", fileCount)
	} else {
		lineCount, err := params.splitFile()
		//lineCount, err := params.splitFile(params.countLines)
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}

		fmt.Printf("Totals lines in file %v: %v.\n", params.sourceFile, lineCount)
	}
}

func (param flagParams) checkFlagErrors() error {
	var err error

	if param.sourceFile == "" && param.countLines == true {
		return fmt.Errorf("the source file names cannot be blank when counting lines in the file")
	}

	if (param.sourceFile == "" || param.destFile == "") && param.countLines == false {
		return fmt.Errorf("the source and destination file names cannot be blank")
	}

	if param.sourceFile == param.destFile && err == nil {
		return fmt.Errorf("the source and destination files must be different")
	}

	if param.lineCount < 1 && param.countLines == false && param.evenSplit == false {
		return fmt.Errorf("the file cannot be split to less than 1 line per file")
	}

	if param.maxFiles < 0 {
		return fmt.Errorf("maximum file count must be zero (maximum files) or greater")
	}

	if param.evenSplit && param.maxFiles == 0 {
		return fmt.Errorf("maximum file count cannot be zero to split file evenly over multiple files")
	}

	if param.countLines && param.repeatHeader {
		return fmt.Errorf("line count only, line headers cannot be repeated when there is no file output")
	}

	return err

}

func (param flagParams) incFilename(counter int) string {

	extn := filepath.Ext(param.destFile)
	return param.destFile[0:len(param.destFile)-len(extn)] + strconv.Itoa(counter) + extn

}

//func (param flagParams) splitFile(countOnly bool) (int, error) {
func (param flagParams) splitFile() (int, error) {

	var err error
	var eoferr error
	var fileCount int
	var fileLines string
	var headerString string

	// Open the source file for reading
	fileReader, err := os.Open(param.sourceFile)
	if err != nil {
		return fileCount, fmt.Errorf("opening source file for reading: %v", err)
	}
	defer fileReader.Close()

	if param.countLines {
		return lineCounter(fileReader)
	}

	if param.evenSplit {
		lineCount, err := lineCounter(fileReader)
		if err != nil {
			return 0, fmt.Errorf("unable to get number of lines in the file :%v", err)
		}
		param.lineCount = int(math.Ceil(float64(lineCount) / float64(param.maxFiles)))
		if param.lineCount < 1 {
			param.lineCount = 1
			param.maxFiles = param.maxFiles - 1
		}
		fileReader.Seek(0, 0) // Reset file pointer as we may have already read through it to get the number of lines.
	}

	if param.repeatHeader {
		headerString, err = getHeader(fileReader)
		if err != nil {
			return 0, fmt.Errorf("error reading the header line: %v", err)
		}
		fileReader.Seek(0, 0) // Reset file pointer as we may have already read through it to get the number of lines.
	}

	bufioReader := bufio.NewReader(fileReader)
	for {
		fileCount++
		outputWriter, err := os.Create(param.incFilename(fileCount))
		if err != nil {
			return fileCount, fmt.Errorf("creating output file: %v", err)
		}

		owriter := bufio.NewWriter(outputWriter)

		// Write the header line if the repeatHeader parameter is set to true and we are not on the first file
		if param.repeatHeader && fileCount > 1 {
			if _, err := owriter.WriteString(headerString); err != nil {
				return fileCount, fmt.Errorf("writing header to output file: %v", err)
			}
		}
		for i := 0; i < param.lineCount; i++ {
			fileLines, eoferr = bufioReader.ReadString('\n')
			if eoferr == io.EOF {
				break
			}
			if _, err := owriter.WriteString(fileLines); err != nil {
				return fileCount, fmt.Errorf("writing to output file: %v", err)
			}
		}
		owriter.Flush()
		outputWriter.Close()

		if (param.maxFiles > 0 && fileCount >= param.maxFiles) || eoferr == io.EOF {
			return fileCount, err
		}
	}
}

func lineCounter(rdr io.Reader) (int, error) {
	var lineCount int

	br := bufio.NewReader(rdr)
	for {
		lineCount = lineCount + 1
		_, readerr := br.ReadString('\n')
		if readerr == io.EOF {
			return lineCount, nil
		}
		if readerr != nil {
			return 0, fmt.Errorf("counting lines in source file: %v", readerr)
		}
	}

}

func getHeader(rdr io.Reader) (string, error) {
	br := bufio.NewReader(rdr)
	hdr, readerr := br.ReadString('\n')
	if readerr != nil {
		return "", fmt.Errorf("counting lines in source file: %v", readerr)
	}

	return hdr, nil

}