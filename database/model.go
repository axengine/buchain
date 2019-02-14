package database

import (
	"database/sql"
	"time"
)

// LedgerHeader object for db
type LedgerHeader struct {
	LedgerID         uint64    `db:"ledgerid"`
	Sequence         string    `db:"sequence"`
	Hash             string    `db:"hash"`
	PrevHash         string    `db:"prevhash"`
	StateRoot        string    `db:"stateroot"`
	TransactionCount uint64    `db:"transactioncount"`
	OperactionCount  uint64    `db:"operactioncount"`
	ClosedAt         time.Time `db:"closedat"`
	TotalCoins       string    `db:"totalcoins"`
	FeePool          string    `db:"feepool"`
	BaseFee          string    `db:"basefee"`
	BaseReserve      string    `db:"basereserve"`
	InflationSeq     uint64    `db:"inflationseq"`
	MaxTxSetSize     uint64    `db:"maxtxsetsize"`
}

// TxData object for db
type TxData struct {
	TxID        uint64 `db:"txid"`
	TxHash      string `db:"txhash"`
	BlockHeight uint64 `db:"blockheight"`
	BlockHash   string `db:"blockhash"`
	ActionCount uint32 `db:"actioncount"`
	ActionID    uint32 `db:"actionid"`
	Src         string `db:"src"`
	Dst         string `db:"dst"`
	Nonce       uint64 `db:"nonce"`
	Amount      string `db:"amount"`
	ResultCode  uint   `db:"resultcode"`
	ResultMsg   string `db:"resultmsg"`
	CreateAt    uint64 `db:"createdat"`
	JData       string `db:"jdata"`
	Memo        string `db:"memo"`
}

// Action object for db
type Action struct {
	ActionID    uint64         `db:"actionid"`
	Typei       int            `db:"typei"`
	Type        string         `db:"type"`
	LedgerSeq   string         `db:"ledgerseq"`
	TxHash      string         `db:"txhash"`
	FromAccount sql.NullString `db:"fromaccount"`
	ToAccount   sql.NullString `db:"toaccount"`
	CreateAt    uint64         `db:"createat"`
	JData       string         `db:"jdata"`
}

// Effect object for db
type Effect struct {
	EffectID  uint64 `db:"effectid"`
	Typei     int    `db:"typei"`
	Type      string `db:"type"`
	LedgerSeq string `db:"ledgerseq"`
	TxHash    string `db:"txhash"`
	ActionID  uint64 `db:"actionid"`
	Account   string `db:"account"`
	CreateAt  uint64 `db:"createat"`
	JData     string `db:"jdata"`
}
