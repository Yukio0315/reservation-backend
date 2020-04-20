package util

import (
	"io/ioutil"
	"log"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// ConvertUtf8ToISOHelper convert utf8 to ISO2022JP
func ConvertUtf8ToISOHelper(str string) string {
	iostr := strings.NewReader(str)
	rio := transform.NewReader(iostr, japanese.ISO2022JP.NewEncoder())
	ret, err := ioutil.ReadAll(rio)
	if err != nil {
		log.Fatal("ConvertUtf8ToISOHelper doesn't work")
	}
	return string(ret)
}

// UniqueID make the uint slice unique
func UniqueID(slice []uint) (result []uint) {
	keys := make(map[uint]bool)
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			result = append(result, entry)
		}
	}
	return result
}
