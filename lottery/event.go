package lottery

import "github.com/sslab-archive/block_lottery_cc/participant"

type Status string
type DrawType string

const (
	STATUS_REGISTERD = "REGISTERED"
	STATUS_DRAWN     = "DRAWN"
)

const (
	DRAW_BITCOIN               = "DRAW_BITCOIN"            // make seed using seed hash
	DRAW_SERVICE_PROVIDER_HASH = "DRAW_SERVICE_PROVIDER"   // make seed using service provider hash
	DRAW_BYZNATINE_PROTOCOL    = "DRAW_BYZANTINE_PROTOCOL" // make seed using BFT tolerance seed making protocol (paper)
)

type Prize struct {
	UUID      string
	Title     string
	Memo      string
	WinnerNum int64
	Winners   []participant.ID
}

type Event struct {
	// event data
	EventName      string     `json:"eventName"` // must be UUID
	Status         Status     `json:"status"`
	CreateTime     string     `json:"createTime"`       // create timestamp
	DeadlineTime   string     `json:"deadlineTime"`   // UNIX timestamp
	MaxParticipant string     `json:"maxParticipant"` // Max number of members
	DrawTypes      []DrawType `json:"drawTypes"`
	Prizes         []Prize    `json:"prizes"`

	// block hash
	TargetBlockHeight string `json:"targetBlockHeight"`
	TargetBlockHash   string `json:"targetBlockHash"`

	// service provider hash
	ServiceProviderPubKey string `json:"serviceProviderPubKey"`
	ServiceProviderHash   string `json:"serviceProviderHash"`

	// result data
	SeedHash        string `json:"seedHash"`
	DrawTxID        string `json:"drawTxID"`
	DrawTxTimestamp string `json:"drawTxTimestamp"`
}
