package define

import (
	"crypto/ecdsa"
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"log"
	"testing"
	"time"
)

func TestTransaction_Sign(t *testing.T) {
	tx := new(Transaction)
	var action Action
	action.CreatedAt = uint64(time.Now().UnixNano())
	action.ID = 0

	privkey, _ := crypto.ToECDSA(ethcmn.Hex2Bytes("a256f9ce4ec985070cd69f4523fb06e5b0cd9d50233eb39a5d423315f855f10f"))
	action.Src = ethcmn.HexToAddress("0x7AF7769E025c8F139D1dF636840029aDa8f82fD4")
	action.Dst = ethcmn.HexToAddress("0xF93d00f90375dDB00FFcE69d76f94B062aAADA6b")
	tx.Actions = append(tx.Actions, &action)

	var privkeys []*ecdsa.PrivateKey
	privkeys = append(privkeys, privkey)
	if err := tx.Sign(privkeys); err != nil {
		t.Fatal("sign err", err)
	}

	log.Println("after sign:", tx.Actions[0])

	if err := tx.CheckSig(); err != nil {
		t.Fatal(err)
	}

}

func Test_sign(t *testing.T) {
	privkey, _ := crypto.ToECDSA(ethcmn.Hex2Bytes("8328d801e557119bf0b16452726150401945f8b838e4526f01b5a01a89c2ab46"))
	//address := ethcmn.HexToAddress("0xb36537aE9B731ff649FE3c766Fc0e789d5b50B4D")
	toSign := []byte("hello world")
	var hash []byte
	{
		hw := sha3.NewKeccak256()
		hw.Write(toSign)
		hash = hw.Sum(nil)
	}
	fmt.Println(len(hash))

	sign, err := crypto.Sign(hash, privkey)
	if err != nil {
		panic(err)
	}
	fmt.Println("sign =", ethcmn.Bytes2Hex(sign))

	pub, err := crypto.SigToPub(hash, sign)
	if err != nil {
		panic(err)
	}

	fmt.Println(crypto.PubkeyToAddress(*pub).Hex())

	//crypto.VerifySignature(privkey.Public(), toSign, sign)
}
