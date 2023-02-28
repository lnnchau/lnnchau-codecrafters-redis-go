package main

import (
	"fmt"
	"strconv"
)

const (
	ARRAY string = "*";
	SIMPLE string = "+";
	BULK_STRING string = "$";
)


type ParsedMessage struct {
	dataType string;

	stringValue string;
	arrayValue []ParsedMessage;
}

func (p ParsedMessage) GetString() (string, error) {
	if (p.dataType != BULK_STRING) {
		return "", fmt.Errorf("Expect %s Got %s", BULK_STRING, p.dataType)
	}
	return p.stringValue, nil
}

func (p ParsedMessage) GetArray() ([]ParsedMessage, error) {
	if (p.dataType != ARRAY) {
		return nil, fmt.Errorf("Expect %s Got %s", ARRAY, p.dataType)
	}
	return p.arrayValue, nil
}


func processArray(s []byte, firstIdx int) (ParsedMessage, int) {
	fmt.Println("Processing array")
	fmt.Printf("First index: %d\tValue: %s\n", firstIdx, string(s[firstIdx]))
	arrayLen, _ := strconv.Atoi(string(s[firstIdx])) 
	fmt.Printf("There are %d elements\n", arrayLen)

	i := firstIdx + 1
	count := 0

	newArray := make([]ParsedMessage, arrayLen)

	for i < len(s) {
		switch string(s[i]) {
		case BULK_STRING:
			retrievedString, nextIdx := processString(s, i+1)
			newArray[count] = retrievedString

			i = nextIdx
			count += 1

		case ARRAY:
			retrievedArray, nextIdx := processArray(s, i+1)
			newArray[count] = retrievedArray

			i = nextIdx
			count += 1
		}

		i++

		if count == arrayLen {
			break
		}
	}
	
	return ParsedMessage{
		dataType: ARRAY,
		arrayValue: newArray,
	}, i
}


func processString(s []byte, firstIdx int) (ParsedMessage, int) {
	fmt.Println("Process string")
	strLength, _ := strconv.Atoi(string(s[firstIdx])) 
	fmt.Printf("There are %d characters\n", strLength)
	
	for i:=firstIdx+1; i < len(s); i++ {
		if string(s[i:i+2]) == "\r\n" {
			firstIdxString := i + 2
			endIdxString := firstIdxString + strLength

			return ParsedMessage{
				dataType: BULK_STRING,
				arrayValue: make([]ParsedMessage, 0),
				stringValue: string(s[firstIdxString:endIdxString]),
			}, endIdxString
		}
	}

	return ParsedMessage{}, -1
}

func ParseRESP(s []byte) ParsedMessage {

	i := 0
	if string(s[i]) != ARRAY {
		fmt.Printf("First byte is not *: %s\n", s)
	}

	res, lastIdx := processArray(s, i+1)

	if lastIdx != (len(s) - 1) {
		fmt.Printf("Message is corrupted: %s\n", s)
	}

	return res
}