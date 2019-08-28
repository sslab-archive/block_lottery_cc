package main

import "encoding/json"

type Participant struct {
	UUID        string `json:"UUID"`
	Information string `json:"information"`
}

func (p Participant) ToLedgerBinary() ([]byte, error) {
	return json.Marshal(p)
}
