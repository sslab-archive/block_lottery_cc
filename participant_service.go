package main

import (
	"crypto/sha256"
	"encoding/binary"
)

func FisherYatesShuffle(arr []Participant ,randomSource string) []Participant {

	shuffledData := make([]Participant, len(arr))

	for i := 0; i < len(arr); i++ {
		shuffledData[i] = arr[i]
	}

	for j := len(arr) - 1; j > 0; j-- {
		hash := sha256.New()
		hash.Write([]byte(randomSource+string(j)))
		h := hash.Sum(nil)
		seed := binary.BigEndian.Uint64(h)
		k := seed % uint64(len(arr))
		shuffledData[j], shuffledData[k] = shuffledData[k], shuffledData[j]
	}

	return shuffledData
}

