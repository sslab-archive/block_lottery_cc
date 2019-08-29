package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("LotteryCC")
var ErrArgsNum = shim.Error("request args num is not matched")
var ErrArgsUnmarshal = shim.Error("request args is not match request object")
var ErrUnknownArgs = shim.Error("request args is not defined")

func ErrArgsRequired(argName string) pb.Response {
	return shim.Error("required arg is not exist : " + argName)
}

type LotteryChaincode struct {
}

func (l *LotteryChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	// Initialize launching lottery
	openTxID := stub.GetTxID()
	chanID := stub.GetChannelID()
	creatorBytes, _ := stub.GetCreator()
	openClientIdentity := fmt.Sprintf("%s", creatorBytes)

	logger.Info("Channel ID: " + chanID)
	logger.Info("LotteryCC Init called")
	logger.Info("Initial test lottery's openTxID: " + openTxID)
	logger.Info("ClientIdentity: " + openClientIdentity)

	return shim.Success(nil)
}

func (l *LotteryChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("Blockchain Lottery Chaincode! invoke method###")
	function, args := stub.GetFunctionAndParameters()
	logger.Info("Invoked method is " + args[0])
	logger.Info(args)
	if function != "invoke" {
		return shim.Error("Unknown function call: " + function)
	}

	switch args[0] {
	case "createLotteryEvent":
		return l.createLotteryEvent(stub, args)
	case "queryLotteryEvent":
		return l.queryLotteryEvent(stub, args)
		/* todo : impl this.
		case "removeLotteryEvent":
			return l.createLotteryEvent(stub, args)
		*/
		/* todo : impl this.
		case "updateLotteryEvent":
			return l.createLotteryEvent(stub, args)
		*/
	case "participateLotteryEvent":
		return l.participateLotteryEvent(stub, args)
	case "drawLotteryEvent":
		return l.drawLotteryEvent(stub, args)
		/* todo: impl this.
		case "verifyLotteryEvent":
			return l.draw(stub, args)
		*/
	}
	return shim.Error("Unknown Invoke Method")
}

// args[0] = function name
// args[1] = json args
func (l *LotteryChaincode) createLotteryEvent(stubInterface shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return ErrArgsNum
	}

	// unmarshal request args
	createLotteryRequest := &CreateLotteryRequest{}

	err := json.Unmarshal([]byte(args[1]), createLotteryRequest)
	logger.Debug(createLotteryRequest)
	if err != nil {
		logger.Error(err)
		return ErrArgsUnmarshal
	}

	// check args is valid
	for _, drawType := range createLotteryRequest.DrawTypes {
		switch drawType {
		case DRAW_BLOCK_HASH:
			if createLotteryRequest.TargetBlock.Height == 0 {
				return ErrArgsRequired("targetBlock")
			}
			break
		case DRAW_SERVICE_PROVIDER_HASH:
			if createLotteryRequest.ServiceProviderHash == "" {
				return ErrArgsRequired("serviceProviderHash")
			}
			break
		default:
			return ErrUnknownArgs
		}
	}
	if len(createLotteryRequest.Prizes) == 0 {
		return ErrArgsRequired("prizes")
	}
	for _, prize := range createLotteryRequest.Prizes {
		if prize.WinnerNum <= 0 {
			return ErrArgsRequired("winnerNum")
		}
	}

	// make event object
	txTimestamp, err := stubInterface.GetTxTimestamp()
	if err != nil {
		logger.Error(err)
		return shim.Error(err.Error())
	}

	txInfo := Transaction{
		ID:               stubInterface.GetTxID(),
		SubmitterID:      createLotteryRequest.SubmitterID,
		SubmitterAddress: createLotteryRequest.SubmitterAddress,
		Timestamp:        txTimestamp.Seconds,
	}
	stubInterface.GetTxTimestamp()
	event := NewEvent(createLotteryRequest, txInfo)

	err = event.SaveToLedger(stubInterface)
	if err != nil {
		logger.Error(err)
		return shim.Error(err.Error())
	}

	responseData, err := json.Marshal(event)
	if err != nil {
		logger.Error(err)
		return shim.Error(err.Error())
	}

	return shim.Success(responseData)
}

// args[0] = function name
// args[1] = json args
func (l *LotteryChaincode) queryLotteryEvent(stubInterface shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return ErrArgsNum
	}

	// unmarshal request args
	queryLotteryRequest := &QueryLotteryRequest{}
	err := json.Unmarshal([]byte(args[1]), queryLotteryRequest)
	if err != nil {
		logger.Error(err)
		return ErrArgsUnmarshal
	}

	switch queryLotteryRequest.QueryType {
	case QUERY_BY_EVENT_ID:
		queryByEventIDRequest := &QueryLotteryByEventIDRequest{}
		err := json.Unmarshal([]byte(args[1]), queryByEventIDRequest)
		if err != nil {
			return ErrArgsUnmarshal
		}

		event, err := LoadEventByUUID(stubInterface, queryByEventIDRequest.EventUUID)
		if err != nil {
			return shim.Error(err.Error())
		}

		responseData, err := json.Marshal(event)
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(responseData)

	case QUERY_BY_DATE_RANGE:
		queryByDateRangeRequest := &QueryLotteryByDateRangeRequest{}
		err := json.Unmarshal([]byte(args[1]), queryByDateRangeRequest)
		if err != nil {
			logger.Error(err)
			return ErrArgsUnmarshal
		}

		events, err := LoadEventByDateRange(
			stubInterface,
			queryByDateRangeRequest.StartDateTimestamp,
			queryByDateRangeRequest.EndDateTimestamp,
		)

		if err != nil {
			logger.Error(err)
			return shim.Error(err.Error())
		}

		responseData, err := json.Marshal(events)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(responseData)

	case QUERY_BY_PARTICIPANT_ID:
	default:
		return ErrUnknownArgs
	}

	return ErrUnknownArgs
}
func (l *LotteryChaincode) participateLotteryEvent(stubInterface shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return ErrArgsNum
	}

	// unmarshal request args
	participateLotteryRequest := &ParticipateLotteryRequest{}
	err := json.Unmarshal([]byte(args[1]), participateLotteryRequest)
	if err != nil {
		return ErrArgsUnmarshal
	}

	event, err := LoadEventByUUID(stubInterface, participateLotteryRequest.EventUUID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if event.UUID == "" {
		return shim.Error("cannot find event, ID :" + participateLotteryRequest.EventUUID)
	}

	txTimestamp, err := stubInterface.GetTxTimestamp()
	if err != nil {
		return shim.Error(err.Error())
	}

	err = event.Participate(participateLotteryRequest.Participant, txTimestamp.Seconds)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = event.SaveToLedger(stubInterface)
	if err != nil {
		return shim.Error(err.Error())
	}

	responseData, err := json.Marshal(event)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(responseData)
}

// args[0] = function name
// args[1] = json args
func (l *LotteryChaincode) drawLotteryEvent(stubInterface shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return ErrArgsNum
	}

	// unmarshal request args
	drawLotteryRequest := &DrawLotteryRequest{}
	err := json.Unmarshal([]byte(args[1]), drawLotteryRequest)
	if err != nil {
		return ErrArgsUnmarshal
	}

	event, err := LoadEventByUUID(stubInterface, drawLotteryRequest.EventUUID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if event.UUID == "" {
		return shim.Error("cannot find event, ID :" + drawLotteryRequest.EventUUID)
	}

	// check required input seed
	for _,drawType := range event.DrawTypes{
		switch drawType {
		case DRAW_BLOCK_HASH:
			if drawLotteryRequest.BlockHash == ""{
				return ErrArgsRequired("blockHash")
			}
			event.TargetBlock.Hash = drawLotteryRequest.BlockHash
		case DRAW_SERVICE_PROVIDER_HASH:
			if drawLotteryRequest.ServiceProviderHash == ""{
				return ErrArgsRequired("serviceProviderHash")
			}
			event.ServiceProviderHash = drawLotteryRequest.ServiceProviderHash
		}
	}

	txTimestamp, err := stubInterface.GetTxTimestamp()
	if err != nil {
		return shim.Error(err.Error())
	}

	txInfo := Transaction{
		ID:               stubInterface.GetTxID(),
		SubmitterID:      drawLotteryRequest.SubmitterID,
		SubmitterAddress: drawLotteryRequest.SubmitterAddress,
		Timestamp:        txTimestamp.Seconds,
	}

	err = event.Draw(txInfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	event.DrawTx = txInfo
	err = event.SaveToLedger(stubInterface)
	if err != nil {
		return shim.Error(err.Error())
	}

	responseData, err := json.Marshal(event)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(responseData)
}

func main() {
	err := shim.Start(new(LotteryChaincode))
	if err != nil {
		logger.Info("Error starting Chaincode: %s", err)
	}
}
