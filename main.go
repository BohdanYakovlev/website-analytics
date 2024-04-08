package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

func getReader(filePath string) (*csv.Reader, *os.File) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	return csv.NewReader(file), file
}

func getRecord(reader *csv.Reader) ([]string, bool) {
	record, err := reader.Read()
	if err == io.EOF {
		return record, false
	}
	if err != nil {
		panic(err)
	}
	if len(record) != 3 {
		panic(errors.New("incorrect record"))
	}
	return record, true
}

func readFiles(firstFilePath, secondFilePath string) {

	firstFileReader, firstFile := getReader(firstFilePath)
	defer firstFile.Close()

	secondFileReader, secondFile := getReader(secondFilePath)
	defer secondFile.Close()

	// Read record headers in files
	_, _ = getRecord(firstFileReader)
	_, _ = getRecord(secondFileReader)

	for {

		firstFileRecord, firstFileFlag := getRecord(firstFileReader)
		if firstFileFlag {
			createProductsFile(firstFileRecord, firstDayFiles)
		}

		secondFileRecord, secondFileFlag := getRecord(secondFileReader)
		if secondFileFlag {
			createProductsFile(secondFileRecord, secondDayFiles)
		}

		if !firstFileFlag && !secondFileFlag {
			break
		}
	}
}

func getDataPaths() filePaths {
	var firstDataFile string
	var secondDataFile string

	fmt.Println("Enter data paths:")
	fmt.Println("First day data:")

	_, err := fmt.Scan(&firstDataFile)
	if err != nil {
		fmt.Println("Can not read first data")
	}

	fmt.Println("Second day data:")

	_, err = fmt.Scan(&secondDataFile)
	if err != nil {
		fmt.Println("Can not read second data")
	}

	return filePaths{firstDataFile, secondDataFile}
}

func printResult(m1 map[string]void, m2 map[string]void) {
	for key, _ := range m1 {
		_, is := m2[key]
		if is {
			fmt.Println(key + ", " + key)
		} else {
			fmt.Println(key + ", ")
		}
	}

	for key, _ := range m2 {
		_, is := m1[key]
		if !is {
			fmt.Println(" , " + key)
		}
	}
}

func printResults(files filePaths) {

	startTime := time.Now()
	bothDayRAM, onlySecondDayRAM := getRAMResult(files)
	endTime := time.Now()
	fmt.Println("RAM method running time: " + endTime.Sub(startTime).String())

	startTime = time.Now()
	bothDayDisk, onlySecondDayDisk := getDiskResult(files)
	endTime = time.Now()
	fmt.Println("DISK method running time: " + endTime.Sub(startTime).String())

	fmt.Println("Users that visited some pages on both days:")
	fmt.Println("RAM, DISK")
	printResult(bothDayRAM, bothDayDisk)

	fmt.Println("Users who did not visit the page on the first day but visited it on the second day:")
	fmt.Println("RAM, DISK")
	printResult(onlySecondDayRAM, onlySecondDayDisk)
}

func main() {

	//files := filePaths{"./Input/first_day.csv", "./Input/second_day.csv"}
	files := getDataPaths()

	printResults(files)

}
