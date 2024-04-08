package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

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

func getDiskResult(files filePaths) (bothDay map[string]void, onlySecondDay map[string]void) {
	createDirectories()

	readFiles(files.firstFilePath, files.secondFilePath)

	bothDayVisitUsers()
	bothDay = getResult(resultPath)

	secondDayNewProductsVisitUsers()
	onlySecondDay = getResult(resultPath)

	deleteDirectory(userDirectory)

	return
}

func bothDayVisitUsers() {

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
}

func secondDayNewProductsVisitUsers() {

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
}

func getResult(directoryPath string) map[string]void {
	userFiles, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		panic(err)
	}

	res := make(map[string]void)

	for _, file := range userFiles {
		if !file.IsDir() {
			_, err = os.Stat(directoryPath + "/" + file.Name())
			if err == nil {
				temp := strings.Replace(file.Name(), ".txt", "", -1)
				res[temp] = void{}
			}
		}
	}
	return res
}

func getUserIdFromFile(fileName string) string {
	fileName = strings.Replace(fileName, "UserId:", "", -1)
	index := strings.Index(fileName, "-")
	return fileName[len("UserId:"):index]
}
