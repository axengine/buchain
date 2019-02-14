package define

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
)

type TxSignature struct {
	Sig []byte
}

// SigHash Hash方法
// 本方法不能随意修改 因为Hash会变，链上也用的这个hash作为交易hash
func (tx *Transaction) SigHash() (h ethcmn.Hash) {
	var signTx Transaction
	for _, v := range tx.Actions {
		var action Action
		action.ID = v.ID
		action.CreatedAt = v.CreatedAt
		action.Src = v.Src
		action.Dst = v.Dst
		action.Amount = v.Amount
		action.Data = v.Data
		signTx.Actions = append(signTx.Actions, &action)
	}

	return rlpHash([]interface{}{
		signTx.Type,
		signTx.Actions,
	})
}

func rlpHash(x interface{}) (h ethcmn.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

func (tx *Transaction) sign(privkeys []*ecdsa.PrivateKey) ([]TxSignature, error) {
	hash := tx.SigHash()
	signatures := make([]TxSignature, len(privkeys))

	for i, v := range privkeys {
		sig, err := crypto.Sign(hash.Bytes(), v)
		if err != nil {
			return signatures, err
		}
		//fmt.Println("sign hash:", tx.SigHash().Hex(), " priv:", v, " sign:", sig)
		signatures[i] = TxSignature{Sig: sig}
	}

	return signatures, nil
}

func (tx *Transaction) Sign(privkeys []*ecdsa.PrivateKey) error {
	signatures, err := tx.sign(privkeys)
	for i, v := range tx.Actions {
		copy(v.SignHex[:], signatures[i].Sig)
	}

	return err
}

func Signer(tx *Transaction, sig []byte) (ethcmn.Address, error) {
	if len(sig) != 65 {
		return ethcmn.Address{}, errors.New("invalid signature length")
	}

	sigHash := tx.SigHash()

	pub, err := crypto.SigToPub(sigHash.Bytes(), sig)
	if err != nil {
		return ethcmn.Address{}, err
	}

	return crypto.PubkeyToAddress(*pub), nil
}

func (tx *Transaction) CheckSig() error {
	for _, v := range tx.Actions {
		address, err := Signer(tx, v.SignHex[:])
		if err != nil {
			return err
		}
		if !bytes.Equal(address.Bytes(), v.Src.Bytes()) {
			return errors.New("signature failed")
		}
	}
	return nil
}
