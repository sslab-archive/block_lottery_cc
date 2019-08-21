package main

import "encoding/json"

type Participant struct {
	UUID        string
	Information string
}

func (p Participant) ToLedgerBinary() ([]byte,error) {
	return json.Marshal(p)
}

