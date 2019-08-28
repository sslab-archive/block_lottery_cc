package main

type Transaction struct {
	ID               string `json:"ID"`
	SubmitterID      string `json:"submitterId"`
	SubmitterAddress string `json:"submitterAddress"`
	Timestamp        int64  `json:"timestamp"`
}
