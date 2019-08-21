package main

type QueryType string

const (
	QUERY_BY_EVENT_ID       QueryType = "QUERY_BY_EVENT_ID"
	QUERY_BY_PARTICIPANT_ID QueryType = "QUERY_BY_PARTICIPANT_ID"
	QUERY_BY_DATE_RANGE     QueryType = "QUERY_BY_DATE_RANGE"
)

type CreateLotteryRequest struct {
	EventName        string     `json:"eventName"`      // must be UUID
	DeadlineTime     int64      `json:"deadlineTime"`   // UNIX timestamp
	MaxParticipant   int64      `json:"maxParticipant"` // Max number of members
	DrawTypes        []DrawType `json:"drawTypes"`
	Prizes           []Prize    `json:"prizes"`
	SubmitterID      string     `json:"submitterID"`
	SubmitterAddress string     `json:"submitterAddress"`

	TargetBlock         BlockInfo `json:"targetBlock"`
	ServiceProviderHash string    `json:"serviceProviderHash"`
}

type QueryLotteryRequest struct {
	QueryType QueryType `json:"queryType"`
}

type QueryLotteryByEventIDRequest struct {
	QueryLotteryRequest
	EventUUID string `json:"eventUUID"`
}

type QueryLotteryByParticipantIDRequest struct {
	QueryLotteryRequest
	ParticipantUUID string `json:"participantUUID"`
}

type QueryLotteryByDateRangeRequest struct {
	QueryLotteryRequest
	StartDateTimestamp int64 `json:"startDateTimestamp"`
	EndDateTimestamp   int64 `json:"endDateTimestamp"`
}

type ParticipateLotteryRequest struct {
	EventUUID   string      `json:"eventUUID"`
	Participant Participant `json:"participant"`
}

type DrawLotteryRequest struct {
	EventUUID           string `json:"eventUUID"`
	BlockHash           string `json:"blockHash"`
	ServiceProviderHash string `json:"serviceProviderHash"`

	SubmitterID      string     `json:"submitterID"`
	SubmitterAddress string     `json:"submitterAddress"`
}

type VerifyLotteryRequest struct {
	EventUUID string `json:"eventUUID"`
	InputHash string `json:"inputHash"`
}
