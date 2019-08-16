package lottery

import "github.com/hyperledger/fabric/core/chaincode/shim"

type EventApi struct {

}

func NewEventApi(stub shim.ChaincodeStubInterface) *EventApi{
	return &EventApi{

	}
}
