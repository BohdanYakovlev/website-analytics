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

type users struct {
	visitedProducts map[string]int
}

func readDay(fileToRead string) (map[string]users, error) {

	file, err := os.Open(fileToRead)
	if err != nil {
		return nil, err
	}

	result := make(map[string]users)

	reader := csv.NewReader(file)
	reader.Read()
	addedRecordIndex := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(record) != 3 {
			return nil, errors.New(fmt.Sprintf("In file %s incorrect filds in %d record", fileToRead, addedRecordIndex))
		}
		addedRecordIndex++
		userId := record[0]
		productId := record[1]

		var tempUser users
		tempUser, ok := result[productId]

		if ok {
			visitSum, ok := tempUser.visitedProducts[userId]
			if ok {
				tempUser.visitedProducts[userId] = visitSum + 1
			} else {
				tempUser.visitedProducts[userId] = 1
			}
		} else {
			tempUser.visitedProducts = make(map[string]int)
			tempUser.visitedProducts[userId] = 1
			result[productId] = tempUser
		}
	}

	return result, nil
}

func printBothDayVisit(firstDayMap map[string]users, secondDayMap map[string]users) {

	resultUsers := make(map[string]int)

	for firstDayMapKey, firstDayMapValue := range firstDayMap {
		secondDayMapValue, ok := secondDayMap[firstDayMapKey]

		if ok {
			for firstDayMapUsersKey, _ := range firstDayMapValue.visitedProducts {
				_, ok := secondDayMapValue.visitedProducts[firstDayMapUsersKey]
				if ok {
					resultUsers[firstDayMapUsersKey] = 0
				}
			}
		}
	}

	for resultUsersKey, _ := range resultUsers {
		fmt.Println(resultUsersKey)
	}
}

func printSecondDayNewProductsVisit(firstDayMap map[string]users, secondDayMap map[string]users) {

	resultUsers := make(map[string]int)
	for secondDayMapProduct, secondDayMapUsers := range secondDayMap {
		firstDayMapUsers, ok := firstDayMap[secondDayMapProduct]
		if !ok {
			for secondDayMapUser, _ := range secondDayMapUsers.visitedProducts {
				resultUsers[secondDayMapUser] = 0
			}
		} else {
			for secondDayMapUser, _ := range secondDayMapUsers.visitedProducts {
				_, ok := firstDayMapUsers.visitedProducts[secondDayMapUser]
				if !ok {
					resultUsers[secondDayMapUser] = 0
				}
			}
		}
	}
	for resultUsersKey, _ := range resultUsers {
		fmt.Println(resultUsersKey)
	}
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

func getReader(filePath string) (*csv.Reader, *os.File) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	return csv.NewReader(file), file
}

func getRecord(reader *csv.Reader, flag *bool) []string {
	record, err := reader.Read()
	if err == io.EOF {
		*flag = false
		return record
	}
	if err != nil {
		panic(err)
	}
	if len(record) != 3 {
		panic(errors.New("incorrect record"))
	}
	return record
}

func createProductsFile(data []string, directoryPath string) {
	fileDirectory := fmt.Sprintf("%s/UserID:%s-ProductId:%s.txt", directoryPath, data[0], data[1])
	file, err := os.Create(fileDirectory)
	if err != nil {
		panic(err)
	}
	file.Close()
}

func createUserFile(userId string) {
	filePath := fmt.Sprintf("%s/%s.txt", resultPath, userId)
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	file.Close()
}

func readFiles(firstFilePath string, secondFilePath string) {

	firstFileReader, firstFile := getReader(firstFilePath)
	defer firstFile.Close()

	secondFileReader, secondFile := getReader(secondFilePath)
	defer secondFile.Close()

	fileFlag := true
	firstFileRecord := getRecord(firstFileReader, &fileFlag)
	secondFileRecord := getRecord(secondFileReader, &fileFlag)

	for fileFlag {

		firstFileRecord = getRecord(firstFileReader, &fileFlag)
		if fileFlag {
			createProductsFile(firstFileRecord, firstDayFiles)
		}

		secondFileRecord = getRecord(secondFileReader, &fileFlag)
		if fileFlag {
			createProductsFile(secondFileRecord, secondDayFiles)
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

func getDataPaths() (string, string) {
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

	return firstDataFile, secondDataFile
}

func main() {

	/*firstDayMap, err := readDay(firstDayFile)
	if err != nil {
		panic(err)
	}
	secondDayMap, err := readDay(secondDayFile)
	if err != nil {
		panic(err)
	}

	fmt.Println("Users that visited some pages on both days:")
	printBothDayVisit(firstDayMap, secondDayMap)

	fmt.Println("Users who did not visit the page on the first day but visited it on the second day:")
	printSecondDayNewProductsVisit(firstDayMap, secondDayMap)*/

	createDirectories()

	firstDataFile, secondDataFile := getDataPaths()

	readFiles(firstDataFile, secondDataFile)

	fmt.Println("Users that visited some pages on both days:")
	printBothDayVisitUsers()

	fmt.Println("Users who did not visit the page on the first day but visited it on the second day:")
	printSecondDayNewProductsVisitUsers()

	deleteDirectory(userDirectory)

}
