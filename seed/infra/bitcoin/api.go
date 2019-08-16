package bitcoin

import (
	"github.com/go-resty/resty"
	"github.com/sslab-archive/block_lottery_cc/seed"
	"strconv"
	"encoding/json"
	"time"
)

//GetLatestBlockURL() string                    // Get URL to get latest block
//GetTargetBlockURL(height int64) string        // Get URL to get target block
type Api struct {
	requestClient *resty.Client
}

func NewApi() *Api {
	return &Api{
		requestClient: resty.New(),
	}
}

func (a *Api) GetLatestBlock() (seed.BlockInfo, error) {
	blockInfo := BlockInfo{}

	resp, err := a.requestClient.R().Get("https://blockchain.info/ko/latestblock")
	if err != nil {
		return seed.BlockInfo{}, err
	}

	err = json.Unmarshal(resp.Body(), blockInfo)
	if err != nil {
		return seed.BlockInfo{}, err
	}

	return blockInfo.ToGeneral(), nil
}

func (a *Api) GetBlock(height int64) (seed.BlockInfo, error) {
	blockInfo := BlockInfo{}

	resp, err := a.requestClient.R().Get("https://api.blockcypher.com/v1/btc/main/blocks/" + strconv.FormatInt(height, 10))
	if err != nil {
		return seed.BlockInfo{}, err
	}

	err = json.Unmarshal(resp.Body(), blockInfo)
	if err != nil {
		return seed.BlockInfo{}, err
	}

	return blockInfo.ToGeneral(), nil
}

func (a *Api) CalcDrawTargetBlockHeight(targetTimestamp int64) (int64, error) {
	latestBlock, err := a.GetLatestBlock()
	if err != nil {
		return 0, err
	}

	calcHeight := latestBlock.GetCalcIntervalHeight()
	calcBlock, err := a.GetBlock(calcHeight)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	passedTime := latestBlock.Timestamp - calcBlock.Timestamp
	passedBlockHeight := latestBlock.Height - calcBlock.Height

	timePerBlock := passedTime / passedBlockHeight

	needBlockForTarget := (targetTimestamp - now.Unix()) / timePerBlock

	return needBlockForTarget, nil

}
