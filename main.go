package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type dataHandler interface {
	handleFirstDayRecord(record []string)
	handleSecondDayRecord(record []string)
	getResult() map[string]void
	printResult()
}

type void struct {
}

type (
	dayMapType map[string]map[string]void
)

type bothDayDataHandler struct {
	firstDayData  dayMapType
	secondDayData dayMapType
	ignoreUsersId map[string]void
}

type filePaths struct {
	firstFilePath  string
	secondFilePath string
}

func (data *bothDayDataHandler) handleFirstDayRecord(record []string) {
	userId := record[0]
	productId := record[1]

	_, is := data.ignoreUsersId[userId]
	//fmt.Println(data.ignoreUsersId)
	if !is {
		secondDayProductMap, is := data.secondDayData[userId]
		//fmt.Println(data.secondDayData)
		if is {
			_, is = secondDayProductMap[productId]
			if is {
				data.ignoreUsersId[userId] = void{}
				delete(data.firstDayData, userId)
				delete(data.secondDayData, userId)
				return
			}
			addRecord(data.firstDayData, record)

		}
		addRecord(data.firstDayData, record)
	}
}

func (data *bothDayDataHandler) handleSecondDayRecord(record []string) {
	data.firstDayData, data.secondDayData = data.secondDayData, data.firstDayData
	data.handleFirstDayRecord(record)
	data.firstDayData, data.secondDayData = data.secondDayData, data.firstDayData
}

func (data *bothDayDataHandler) getResult() map[string]void {
	return data.ignoreUsersId
}

func (data *bothDayDataHandler) printResult() {
	res := data.getResult()

	fmt.Println("Users that visited some pages on both days:")
	for key, _ := range res {
		fmt.Println(key)
	}
}

type onlySecondDayDataHandler struct {
	firstDayData  dayMapType
	secondDayData dayMapType
}

func (data *onlySecondDayDataHandler) handleFirstDayRecord(record []string) {
	addRecord(data.firstDayData, record)
}

func (data *onlySecondDayDataHandler) handleSecondDayRecord(record []string) {
	addRecord(data.secondDayData, record)
}

func (data *onlySecondDayDataHandler) getResult() map[string]void {
	resultMap := make(map[string]void)
	for secondDayUserId, secondDayProductMap := range data.secondDayData {
		firstDayProductMap, is := data.firstDayData[secondDayUserId]
		if is {
			for productId, _ := range secondDayProductMap {
				_, is = firstDayProductMap[productId]
				if !is {
					resultMap[secondDayUserId] = void{}
				}
			}
		} else {
			resultMap[secondDayUserId] = void{}
		}
	}

	return resultMap
}

func (data *onlySecondDayDataHandler) printResult() {
	res := data.getResult()

	fmt.Println("Users who did not visit the page on the first day but visited it on the second day:")
	for key, _ := range res {
		fmt.Println(key)
	}
}

func addRecord(dayMap dayMapType, record []string) {
	userId := record[0]
	productId := record[1]

	productMap, is := dayMap[userId]
	if is {
		_, is := productMap[productId]
		if !is {
			productMap[productId] = void{}
		}
		return
	}
	tempMap := make(map[string]void)
	tempMap[productId] = void{}
	dayMap[userId] = tempMap
}

func handleResult(files filePaths, handler dataHandler) {

	firstFileReader, firstDataFile := getReader(files.firstFilePath)
	defer firstDataFile.Close()

	secondFileReader, secondDataFile := getReader(files.secondFilePath)
	defer secondDataFile.Close()

	_, _ = getRecord(firstFileReader)
	_, _ = getRecord(secondFileReader)

	for {

		firstFileRecord, firstFileFlag := getRecord(firstFileReader)
		if firstFileFlag {
			handler.handleFirstDayRecord(firstFileRecord)
		}

		secondFileRecord, secondFileFlag := getRecord(secondFileReader)
		if secondFileFlag {
			handler.handleSecondDayRecord(secondFileRecord)
		}

		if !firstFileFlag && !secondFileFlag {
			break
		}
	}

	handler.printResult()
}

func printRAMResult(files filePaths) {

	bothHandler := bothDayDataHandler{make(dayMapType), make(dayMapType), make(map[string]void)}
	handleResult(files, &bothHandler)

	onlySecondHandler := onlySecondDayDataHandler{make(dayMapType), make(dayMapType)}
	handleResult(files, &onlySecondHandler)
}

///////////////////////////////////////////////////////////////////////////////

const (
	userDirectory  = "./Users"
	firstDayFiles  = "./Users/FirstDay"
	secondDayFiles = "./Users/SecondDay"
	resultPath     = "./Users/Result"
)

func createDirectoryIfExist(directoryPath string) {
	_ = os.Mkdir(directoryPath, 0755)
}

func createDirectories() {
	createDirectoryIfExist(userDirectory)
	createDirectoryIfExist(firstDayFiles)
	createDirectoryIfExist(secondDayFiles)
	createDirectoryIfExist(resultPath)
}

func deleteDirectory(directoryPath string) {
	os.RemoveAll(directoryPath)
}

func recreateDirectory(directoryPath string) {
	deleteDirectory(directoryPath)
	createDirectoryIfExist(directoryPath)
}

func createFile(fileDirectory string) {
	file, err := os.Create(fileDirectory)
	if err != nil {
		panic(err)
	}
	file.Close()
}

func createProductsFile(data []string, directoryPath string) {
	fileDirectory := fmt.Sprintf("%s/UserID:%s-ProductId:%s.txt", directoryPath, data[0], data[1])
	createFile(fileDirectory)
}

func createUserFile(userId string) {
	filePath := fmt.Sprintf("%s/%s.txt", resultPath, userId)
	createFile(filePath)
}

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

func printBothDayVisitUsers() {

	recreateDirectory(resultPath)

	firstFiles, err := ioutil.ReadDir(firstDayFiles)
	if err != nil {
		panic(err)
	}

	for _, file := range firstFiles {
		if !file.IsDir() {
			_, err = os.Stat(secondDayFiles + "/" + file.Name())
			if err == nil {
				createUserFile(getUserIdFromFile(file.Name()))
			}
		}
	}

	printResult("./Users/Result")
}

func printSecondDayNewProductsVisitUsers() {

	recreateDirectory(resultPath)

	secondFiles, err := ioutil.ReadDir(secondDayFiles)
	if err != nil {
		panic(err)
	}

	for _, file := range secondFiles {
		if !file.IsDir() {
			_, err = os.Stat(firstDayFiles + "/" + file.Name())
			if os.IsNotExist(err) {
				createUserFile(getUserIdFromFile(file.Name()))
			}
		}
	}

	printResult("./Users/Result")
}

func printResult(directoryPath string) {
	userFiles, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		panic(err)
	}

	for _, file := range userFiles {
		if !file.IsDir() {
			_, err = os.Stat(directoryPath + "/" + file.Name())
			if err == nil {
				fmt.Println(strings.Replace(file.Name(), ".txt", "", -1))
			}
		}
	}
}

func getUserIdFromFile(fileName string) string {
	fileName = strings.Replace(fileName, "UserId:", "", -1)
	index := strings.Index(fileName, "-")
	return fileName[len("UserId:"):index]
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

func printDiskResult(files filePaths) {
	createDirectories()

	readFiles(files.firstFilePath, files.secondFilePath)

	fmt.Println("Users that visited some pages on both days:")
	printBothDayVisitUsers()

	fmt.Println("Users who did not visit the page on the first day but visited it on the second day:")
	printSecondDayNewProductsVisitUsers()

	deleteDirectory(userDirectory)
}

func main() {

	files := filePaths{"./Input/first_day.csv", "./Input/second_day.csv"}
	//firstDataFile, secondDataFile := getDataPaths()

	fmt.Println("RAM method:")
	printRAMResult(files)

	fmt.Println("\nDisk method")
	printDiskResult(files)

}
