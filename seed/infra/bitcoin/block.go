package bitcoin

import (
	"github.com/sslab-archive/block_lottery_cc/seed"
)

const (
	DIFFICULT_CALC_CYCLE = 2016
)

type BlockInfo struct {
	seed.BlockInfo
	BlockIndex int64  `json:"block_index"`
	Nonce      string `json:"nonce"`
}

func NewBlock() *BlockInfo {
	return &BlockInfo{
		BlockInfo: seed.BlockInfo{
			Hash:      "",
			Timestamp: 0,
			Height:    0,
			PrevBlock: "",
		},
		BlockIndex: 0,
		Nonce:      "",
	}
}

func (b *BlockInfo) ToGeneral() seed.BlockInfo {
	return seed.BlockInfo{
		Hash:      b.Hash,
		Timestamp: b.Timestamp,
		Height:    b.Height,
		PrevBlock: b.PrevBlock,
	}
}

func (b *BlockInfo) GetBlockType() seed.BlockType {
	return seed.BITCOIN
}

func (b *BlockInfo) GetCalcIntervalHeight() int64 {
	if b.Height == 0 {
		return 1
	}
	afterRetarget := b.Height % DIFFICULT_CALC_CYCLE
	if afterRetarget < 10 {
		return b.Height - (afterRetarget + DIFFICULT_CALC_CYCLE)
	}

	return b.Height - afterRetarget
}
