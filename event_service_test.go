package main

import (
	"testing"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
)

func TestLoadEventByDateRange(t *testing.T) {
	l := new(LotteryChaincode)
	m := shim.NewMockStub("lottery_cc",l)
	m.MockTransactionStart("c5a2e0e7231216c9483c170d03e955bee90d41b8262841ebda5545f5a3eab73e")
	re := l.createLotteryEvent(m,[]string{"",`{
    "UUID": "blj5u1aotin61uf67t90",
    "eventName": "test event name3",
    "status": "REGISTERED",
    "createTime": 1566990085,
    "deadlineTime": 1566990085,
    "maxParticipant": 100,
    "participants": [],
    "drawTypes": [
        "DRAW_BLOCK_HASH"
    ],
    "prizes": [
        {
            "UUID": "blj5u1aotin61uf67t8g",
            "title": "prize3",
            "memo": "memo",
            "winnerNum": 10,
            "winners": null
        }
    ],
    "targetBlock": {
        "blockType": "BITCOIN",
        "hash": "",
        "time": 0,
        "height": 592122
    },
    "serviceProviderHash": "",
    "seedHash": "",
    "eventCreateTx": {
        "ID": "c5a2e0e7231216c9483c170d03e955bee90d41b8262841ebda5545f5a3eab73e",
        "submitterId": "SERVICE_PROVIDER_1",
        "submitterAddress": "not impl",
        "timestamp": 1566990085
    },
    "drawTx": {
        "ID": "",
        "submitterId": "",
        "submitterAddress": "",
        "timestamp": 0
    }
}`})
	fmt.Printf("return : %s", re.Payload)

	fmt.Println(m.State)
}