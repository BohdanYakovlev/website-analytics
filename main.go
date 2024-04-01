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

type user struct {
	visitedProducts map[string]int
}

func readDay(fileToRead string) (map[string]user, error) {

	file, err := os.Open(fileToRead)
	if err != nil {
		return nil, err
	}

	result := make(map[string]user)

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
			return nil, errors.New(fmt.Sprintf("In file %s incorrect filds in %n record", fileToRead, addedRecordIndex))
		}
		addedRecordIndex++
		userId := record[0]
		productId := record[1]

		var tempUser user
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

func main() {

}
