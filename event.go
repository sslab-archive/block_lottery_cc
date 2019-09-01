package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/rs/xid"
	"time"
	"strconv"
)

type Status string
type DrawType string

const (
	STATUS_REGISTERD Status = "REGISTERED"
	STATUS_DRAWN     Status = "DRAWN"
	STATUS_REMOVED   Status = "REMOVED"
)

const (
	DRAW_BLOCK_HASH            DrawType = "DRAW_BLOCK_HASH"         // make seed using block hash
	DRAW_SERVICE_PROVIDER_HASH DrawType = "DRAW_SERVICE_PROVIDER"   // make seed using service provider hash
	DRAW_BYZNATINE_PROTOCOL    DrawType = "DRAW_BYZANTINE_PROTOCOL" // make seed using BFT tolerance seed making protocol (paper)
)

type Prize struct {
	UUID      string        `json:"UUID"`
	Title     string        `json:"title"`
	Memo      string        `json:"memo"`
	WinnerNum int64         `json:"winnerNum"`
	Winners   []Participant `json:"winners"`
}
type EventKeyInfo struct {
	UUID       string `json:"UUID"`
	CreateTime int64  `json:"createTime"`
}

type Event struct {
	// event data
	UUID           string        `json:"UUID"`
	EventName      string        `json:"eventName"` // must be UUID
	Status         Status        `json:"status"`
	Contents       string        `json:"contents"`
	CreateTime     int64         `json:"createTime"`     // create timestamp
	DeadlineTime   int64         `json:"deadlineTime"`   // UNIX timestamp
	MaxParticipant int64         `json:"maxParticipant"` // Max number of members
	Participants   []Participant `json:"participants"`
	DrawTypes      []DrawType    `json:"drawTypes"`
	Prizes         []Prize       `json:"prizes"`

	// block hash
	TargetBlock BlockInfo `json:"targetBlock"`

	// service provider hash
	ServiceProviderHash string `json:"serviceProviderHash"`

	// result data
	SeedHash string `json:"seedHash"`

	// transaction info
	EventCreateTx Transaction `json:"eventCreateTx"`
	DrawTx        Transaction `json:"drawTx"`
}

func (e *Event) GetKey() string {
	return MakeKeyByUUID(e.UUID)
}

func (e *Event) SaveToLedger(stubInterface shim.ChaincodeStubInterface) error {
	b, err := e.ToLedgerBinary()
	if err != nil {
		return err
	}

	// save e information if e is changed
	val, err := stubInterface.GetState(e.GetKey())

	if bytes.Compare(b, val) != 0 {
		err = stubInterface.PutState(e.GetKey(), b)
		if err != nil {
			return err
		}
	}

	// save participants if participant data is changed
	for _, participant := range e.Participants {
		key, err := stubInterface.CreateCompositeKey(e.GetKey(), []string{"participants", participant.UUID})
		if err != nil {
			return err
		}

		val, err := stubInterface.GetState(key)
		if err != nil {
			return err
		}

		// val == nil -> new data
		// err != nil && recordedData != participant -> data is changed
		recordedData := &Participant{}
		if err := json.Unmarshal(val, recordedData); (err != nil && *recordedData != participant) || val == nil {
			b, err := participant.ToLedgerBinary()
			if err != nil {
				return err
			}

			err = stubInterface.PutState(key, b)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// it will remove participants data in event data.
// participants data will be record in composite key ( event_{tx timestamp}_{eventUUID}~participants~{participantsUUID} )
func (e Event) ToLedgerBinary() ([]byte, error) {
	e.Participants = nil
	return json.Marshal(e)
}

func (e *Event) Participate(participant Participant, txTimestamp int64) error {
	if int64(len(e.Participants)) >= e.MaxParticipant {
		return errors.New("this event has exceeded the number of participants")
	}

	if txTimestamp > e.DeadlineTime {
		return errors.New("participation is closed")
	}

	if e.containParticipant(participant.UUID) {
		return errors.New("duplicate participate")
	}

	e.Participants = append(e.Participants, participant)
	return nil
}

func (e *Event) Draw(tx Transaction) error {
	if tx.Timestamp < e.DeadlineTime {
		return errors.New("deadline is not passed")
	}

	if e.Status != STATUS_REGISTERD {
		return errors.New("status is not registered. check is removed or already drawn")
	}
	e.Status = STATUS_DRAWN
	concatSeed := e.SeedHash + "_PLUS_" + e.ServiceProviderHash

	var usingParticipants []Participant

	if int64(len(e.Participants)) > e.MaxParticipant {
		usingParticipants = e.Participants[0:e.MaxParticipant]
	} else {
		usingParticipants = e.Participants
	}
	shuffledParticipant := FisherYatesShuffle(usingParticipants, concatSeed)

	totalPrizeNum := int64(0)
	for _, prize := range e.Prizes {
		totalPrizeNum += prize.WinnerNum
	}
	logger.Debug("total winner num : " + strconv.FormatInt(totalPrizeNum, 10))

	// # of participant < # of winner
	if totalPrizeNum > int64(len(shuffledParticipant)) {
		logger.Debug("winner is too big...")
		participantIdx := 0
		for idx := range e.Prizes {
			e.Prizes[idx].Winners = make([]Participant, 0)
			if int64(len(shuffledParticipant)-participantIdx) >= e.Prizes[idx].WinnerNum {
				e.Prizes[idx].Winners = shuffledParticipant[participantIdx : participantIdx+int(e.Prizes[idx].WinnerNum)]
				participantIdx += int(e.Prizes[idx].WinnerNum)
			} else {
				for i := participantIdx; i < len(shuffledParticipant); i++ {
					e.Prizes[idx].Winners = append(e.Prizes[idx].Winners, shuffledParticipant[i])
				}
				break
			}
		}
	} else {
		passedIdx := 0
		for idx := range e.Prizes {
			e.Prizes[idx].Winners = shuffledParticipant[passedIdx : passedIdx+int(e.Prizes[idx].WinnerNum)]
			passedIdx += int(e.Prizes[idx].WinnerNum)
		}
	}
	return nil
}

func (e *Event) containParticipant(participantUUID string) bool {
	for _, p := range e.Participants {
		if p.UUID == participantUUID {
			return true
		}
	}
	return false
}

func NewEvent(request *CreateLotteryRequest, createEventTX Transaction) Event {
	requestPrize := request.Prizes
	for idx := range requestPrize {
		requestPrize[idx].UUID = xid.New().String()
	}
	return Event{
		UUID:                xid.New().String(),
		EventName:           request.EventName,
		Status:              STATUS_REGISTERD,
		Contents:            request.Contents,
		CreateTime:          time.Now().Unix(),
		DeadlineTime:        request.DeadlineTime,
		MaxParticipant:      request.MaxParticipant,
		Participants:        make([]Participant, 0),
		DrawTypes:           request.DrawTypes,
		Prizes:              requestPrize,
		TargetBlock:         request.TargetBlock,
		ServiceProviderHash: request.ServiceProviderHash,
		SeedHash:            "",
		EventCreateTx:       createEventTX,
		DrawTx:              Transaction{},
	}
}

func MakeKeyByUUID(UUID string) string {
	return "event_UUID_" + UUID
}
