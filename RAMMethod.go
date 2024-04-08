package main

import "fmt"

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
	if !is {
		secondDayProductMap, is := data.secondDayData[userId]
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
}

func getRAMResult(files filePaths) (bothDay map[string]void, onlySecondDay map[string]void) {

	bothHandler := bothDayDataHandler{make(dayMapType), make(dayMapType), make(map[string]void)}
	handleResult(files, &bothHandler)
	bothDay = bothHandler.getResult()

	onlySecondHandler := onlySecondDayDataHandler{make(dayMapType), make(dayMapType)}
	handleResult(files, &onlySecondHandler)
	onlySecondDay = onlySecondHandler.getResult()

	return
}
