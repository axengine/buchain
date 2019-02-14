package define

import (
	"encoding/binary"
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/crypto/merkle"
	"time"
)

type AppHeader struct {
	PrevHash  ethcmn.Hash // global, just used to calculate header-hash
	StateRoot ethcmn.Hash // fill after statedb commit
	BlockHash ethcmn.Hash //区块hash
	Height    uint64      // refresh by new block
	TxCount   uint64      // fill when ready to save
	OpCount   uint64      // fill when ready to save
	ClosedAt  time.Time   // refresh by new block
}

func (h *AppHeader) String() string {
	return fmt.Sprintf("prevhash:%v,stateRoot:%v,height:%v,txCount:%v,opCount:%v,ClosedAt:%v",
		h.PrevHash.Hex(),
		h.StateRoot.Hex(),
		h.Height,
		h.TxCount,
		h.OpCount,
		h.ClosedAt)
}

// Hash hash
func (h *AppHeader) Hash() []byte {
	b0, b1, b2 := make([]byte, 8), make([]byte, 8), make([]byte, 8)
	binary.LittleEndian.PutUint64(b0, h.Height)
	binary.LittleEndian.PutUint64(b1, h.TxCount)
	binary.LittleEndian.PutUint64(b2, h.OpCount)
	m := map[string][]byte{
		"PrevHash":  h.PrevHash.Bytes(),
		"StateRoot": h.StateRoot.Bytes(),
		"Height":    b0,
		"TxCount":   b1,
		"OpCount":   b2,
	}

	return merkle.SimpleHashFromMap(m)
}
