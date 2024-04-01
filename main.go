package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	firstDayFile  = "Input/first_day.csv"
	secondDayFile = "Input/second_day.csv"
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
func main() {
	firstDayMap, err := readDay(firstDayFile)
	if err != nil {
		panic(err)
	}
	secondDayMap, err := readDay(secondDayFile)
	if err != nil {
		panic(err)
	}

	fmt.Println("Users that visited some pages on both days:")
	printBothDayVisit(firstDayMap, secondDayMap)

	fmt.Println("Users who did not visit the page on the first day but visited it on the second day")
	printSecondDayNewProductsVisit(firstDayMap, secondDayMap)
}
