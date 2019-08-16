package seed

const (
	BITCOIN BlockType = "BITCOIN"
	EOS     BlockType = "EOS"
)

type BlockType string

type BlockInfo struct {
	Block
	Hash      string `json:"hash"`
	Timestamp int64  `json:"time"`
	Height    int64  `json:"height"`
	PrevBlock string `json:"prev_block"`
}

type Block interface {
	GetBlockType() BlockType                      // Get block type (ex bitcoin, EOS ...)
	ToGeneral() BlockInfo
	GetCalcIntervalHeight() int64 // Get height values to calculate the average block creation time from input height
}

type BlockApi interface {
	GetLatestBlock() (BlockInfo,error)
	GetBlock(height int64) (BlockInfo,error)
	CalcDrawTargetBlockHeight(targetTimestamp int64) (int64, error)
}
