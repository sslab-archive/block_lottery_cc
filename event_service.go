package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"strconv"
)

func LoadEventByKey(stubInterface shim.ChaincodeStubInterface, key string) (*Event, error) {
	event := &Event{}
	b, err := stubInterface.GetState(key)
	if b == nil {
		return &Event{}, nil
	}
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, event)
	if err != nil {
		return nil, err
	}

	// load participant data
	event.Participants = make([]Participant, 0)
	participantIterator, err := stubInterface.GetStateByPartialCompositeKey(event.GetKey(), []string{"participants"})
	defer participantIterator.Close()
	if err != nil {
		return nil, err
	}

	for participantIterator.HasNext() {
		b, err := participantIterator.Next()
		if err != nil {
			return nil, err
		}

		p := &Participant{}
		err = json.Unmarshal(b.Value, p)
		if err != nil {
			return nil, err
		}

		event.Participants = append(event.Participants, *p)
	}
	return event, nil
}

func LoadEventByDateRange(stubInterface shim.ChaincodeStubInterface, startDate int64, endDate int64) ([]Event, error) {
	startKey := "event_"+strconv.FormatInt(startDate,10)
	endKey := "event_"+strconv.FormatInt(endDate,10)

	eventIter, err := stubInterface.GetStateByRange(startKey,endKey)
	defer eventIter.Close()
	if err !=nil{
		return nil, err
	}

	events := make([]Event,0)
	for eventIter.HasNext(){
		kv,err := eventIter.Next()
		if err!= nil{
			return nil, err
		}

		event := Event{}
		err = json.Unmarshal(kv.Value,event)

		if err !=nil{
			return nil,err
		}

		events = append(events, event)
	}

	return events,nil
}
