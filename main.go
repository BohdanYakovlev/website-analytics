package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
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
	secondFileReader, secondFile := getReader(secondFilePath)

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

	firstFile.Close()
	secondFile.Close()

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

func refresh(filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	writer := csv.NewWriter(file)
	n := rand.Intn(10000)
	for i := 0; i < n; i++ {
		writer.Write([]string{strconv.Itoa(rand.Intn(1000)), strconv.Itoa(rand.Intn(1000)), "1"})
		writer.Flush()
	}
	file.Close()
}

func megaTest() {
	files := filePaths{"./Input/first_day.csv", "./Input/second_day.csv"}

	for {
		refresh(files.firstFilePath)
		refresh(files.secondFilePath)
		printResults(files)
		os.Remove(files.firstFilePath)
		os.Remove(files.secondFilePath)
	}
}

func main() {

	//megaTest()

	files := getDataPaths()

	printResults(files)

}
