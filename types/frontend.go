package types

import (
	"database/sql/driver"
	"encoding/json"
	"math/big"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/theQRL/go-zond/common"
	"github.com/theQRL/go-zond/common/hexutil"
)

type Tag string

const (
	ValidatorTagsWatchlist Tag = "watchlist"
)

type ErrorResponse struct {
	Status string // e.g. "200 OK"
	Body   string
}

func (e *ErrorResponse) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &e)
}

func (a ErrorResponse) Value() (driver.Value, error) {
	return json.Marshal(a)
}

type EnsSearchPageData = struct {
	Error  string
	Search string
	Result *EnsDomainResponse
}

type GasNowPageData struct {
	Code int `json:"code"`
	Data struct {
		Rapid     *big.Int `json:"rapid"`
		Fast      *big.Int `json:"fast"`
		Standard  *big.Int `json:"standard"`
		Slow      *big.Int `json:"slow"`
		Timestamp int64    `json:"timestamp"`
		Price     float64  `json:"price,omitempty"`
		PriceUSD  float64  `json:"priceUSD"`
		Currency  string   `json:"currency,omitempty"`
	} `json:"data"`
}

type Eth1AddressSearchItem struct {
	Address string `json:"address"`
	Name    string `json:"name"`
	Token   string `json:"token"`
}

type RawMempoolResponse struct {
	Pending map[string]map[string]*RawMempoolTransaction `json:"pending"`
	Queued  map[string]map[string]*RawMempoolTransaction `json:"queued"`
	// BaseFee map[string]map[string]*RawMempoolTransaction `json:"baseFee"`

	TxsByHash map[common.Hash]*RawMempoolTransaction
}

func (mempool RawMempoolResponse) FindTxByHash(txHashString string) *RawMempoolTransaction {
	return mempool.TxsByHash[common.HexToHash(txHashString)]
}

type RawMempoolTransaction struct {
	Hash             common.Hash     `json:"hash"`
	From             *common.Address `json:"from"`
	To               *common.Address `json:"to"`
	Value            *hexutil.Big    `json:"value"`
	Gas              *hexutil.Big    `json:"gas"`
	GasFeeCap        *hexutil.Big    `json:"maxFeePerGas,omitempty"`
	GasTipCap        *hexutil.Big    `json:"maxPriorityFeePerGas,omitempty"`
	GasPrice         *hexutil.Big    `json:"gasPrice"`
	Nonce            *hexutil.Big    `json:"nonce"`
	Input            *string         `json:"input"`
	TransactionIndex *hexutil.Big    `json:"transactionIndex"`
}

type MempoolTxPageData struct {
	RawMempoolTransaction
	TargetIsContract   bool
	IsContractCreation bool
}

type SyncCommitteesStats struct {
	ParticipatedSlots uint64 `db:"participated_sync" json:"participatedSlots"`
	MissedSlots       uint64 `db:"missed_sync" json:"missedSlots"`
	OrphanedSlots     uint64 `db:"orphaned_sync" json:"-"`
	ScheduledSlots    uint64 `json:"scheduledSlots"`
}

type SignatureType string

const (
	MethodSignature SignatureType = "method"
	EventSignature  SignatureType = "event"
)

type SignatureImportStatus struct {
	LatestTimestamp *string `json:"latestTimestamp"`
	NextPage        *string `json:"nextPage"`
	HasFinished     bool    `json:"hasFinished"`
}

type Signature struct {
	Id        int64  `json:"id"`
	CreatedAt string `json:"created_at"`
	Text      string `json:"text_signature"`
	Hex       string `json:"hex_signature"`
	Bytes     string `json:"bytes_signature"`
}

type SearchValidatorsByEth1Result []struct {
	ZondAddress      string        `db:"from_address_text" json:"zond_address"`
	ValidatorIndices pq.Int64Array `db:"validatorindices" json:"validator_indices"`
	Count            uint64        `db:"count" json:"-"`
}
