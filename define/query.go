package define

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
)

type QueryBase struct {
	Order  string
	Limit  uint64
	Cursor uint64

	Begin uint64
	End   uint64
}
type TxQuery struct {
	QueryBase
	TxHash    ethcmn.Hash
	Account   ethcmn.Address
	Direction uint8
}

type BlockAccount struct {
	Addr    string `json:"address"`
	Balance string `json:"balance"`
}
