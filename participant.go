package main

import "encoding/json"

type Participant struct {
	UUID            string      `json:"UUID"`
	Information     string      `json:"information"`
	AuthInformation string      `json:"authInformation"`
	ParticipateTx   Transaction `json:"participateTx"`
}

func (p Participant) ToLedgerBinary() ([]byte, error) {
	return json.Marshal(p)
}
