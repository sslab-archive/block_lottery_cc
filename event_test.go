package main

import (
	"testing"
	"encoding/json"
	"fmt"
)

func TestEvent_SaveToLedger(t *testing.T) {
	a := []byte(string("{\"UUID\":\"uuid test\" , \"Information\":\"infotest\", \"unexpected\":\"uuu\"}"))

	tp := &Participant{}

	type Test2 struct {
		Participant
		unexpected  string
	}
	err := json.Unmarshal(a, tp)

	fmt.Println(tp)
	fmt.Println(err)

}

func b() []byte {
	return nil
}
