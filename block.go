package main

const (
	BITCOIN BlockType = "BITCOIN"
	EOS     BlockType = "EOS"
)

type BlockType string

type BlockInfo struct {
	BlockType BlockType `json:"blockType"`
	Hash      string    `json:"hash"`
	Timestamp int64     `json:"time"`
	Height    int64     `json:"height"`
}

func NewBitcoinBlock() BlockInfo {
	return BlockInfo{
		BlockType: BITCOIN,
		Hash:      "",
		Timestamp: 0,
		Height:    0,
	}
}
