package main

import (
	"testing"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"strconv"
)

func TestFisherYatesShuffle(t *testing.T) {

	arr := make([]string,0)
	arr = append(arr, "1")
	arr = append(arr, "2")
	arr = append(arr, "3")
	arr = append(arr, "4")

	shuffledData := make([]string, len(arr))

	randomSource := "testSource"

	for i := 0; i < len(arr); i++ {
		shuffledData[i] = arr[i]
	}

	fmt.Println(arr)

	for j := len(arr) - 1; j > 0; j-- {
		hash := sha256.Sum256([]byte(randomSource+strconv.Itoa(j)))
		var h []byte
		h = hash[:]
		seed := binary.BigEndian.Uint64(h)
		k := seed % uint64(len(arr))
		shuffledData[j], shuffledData[k] = shuffledData[k], shuffledData[j]
	}

	fmt.Println(shuffledData)

}
