package btchain

import (
	"fmt"
	"github.com/axengine/btchain/define"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/tendermint/go-amino"
	"math/big"
	"testing"
)

func Test_commit(t *testing.T) {
	db := ethdb.NewMemDatabase()

	statedb, _ := state.New(ethcmn.Hash{}, state.NewDatabase(db))
	addr := ethcmn.HexToAddress("0x02865c395bfd104394b786a264662d02177897391aba1155f854cb1065b6a444e5")
	statedb.AddBalance(addr, big.NewInt(1000000000000))
	//root := statedb.IntermediateRoot(false)
	appHash, err := statedb.Commit(false)
	if err != nil {
		t.Fatal(err)
	}
	//t.Log("root hash:", root.Hex(), " appHash:", appHash.Hex())
	err = statedb.Database().TrieDB().Commit(appHash, true)
	if err != nil {
		t.Fatal(err)
	}
	statedb, err = state.New(appHash, state.NewDatabase(db))
	if err != nil {
		t.Fatal(err)
	}
}

func Test_validators(t *testing.T) {
	db := ethdb.NewMemDatabase()

	var vls []*define.Validator
	vls = append(vls, &define.Validator{
		"xxxxxxx",
		10,
	})
	vls = append(vls, &define.Validator{
		"xxxxxxx2",
		20,
	})

	b, err := amino.MarshalBinaryBare(&vls)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Put([]byte("validatorsKey"), b); err != nil {
		t.Fatal(err)
	}

	bb, err := db.Get([]byte("validatorsKey"))
	if err != nil {
		t.Fatal(err)
	}

	var xx []*define.Validator
	err = amino.UnmarshalBinaryBare(bb, &xx)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(xx[0].PubKey, xx[0].Power)
	fmt.Println(xx[1].PubKey, xx[1].Power)
}
