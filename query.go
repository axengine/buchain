package btchain

import (
	"encoding/json"
	"github.com/axengine/btchain/define"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"go.uber.org/zap"
)

const (
	QUERY_TX      = "/tx"
	QUERY_ACCOUNT = "/account"
	QUERY_NONCE   = "/nonce"
)

var (
	ZERO_ADDRESS = ethcmn.Address{}
	ZERO_HASH    = ethcmn.Hash{}
)

func (app *BTApplication) QueryTx(tx []byte) define.Result {
	var query define.TxQuery
	if err := rlp.DecodeBytes(tx, &query); err != nil {
		app.logger.Debug("rlp.DecodeBytes", zap.Error(err))
		return define.NewError(define.CodeType_EncodingError, err.Error())
	}
	app.logger.Debug("QueryTx", zap.String("hash", query.TxHash.Hex()), zap.String("address", query.Account.Hex()))
	if query.Account != ZERO_ADDRESS {
		result, err := app.dataM.QueryAccountTxs(&query.Account, query.Direction, query.Cursor, query.Limit, query.Order)
		if err != nil {
			return define.NewError(define.CodeType_InternalError, err.Error())
		}
		return MakeResultData(result)
	}
	if query.TxHash != ZERO_HASH {
		result, err := app.dataM.QueryTxByHash(&query.TxHash)
		if err != nil {
			return define.NewError(define.CodeType_InternalError, err.Error())
		}
		return MakeResultData(result)
	}

	result, err := app.dataM.QueryAllTxs(query.Cursor, query.Limit, query.Order)
	if err != nil {
		return define.NewError(define.CodeType_InternalError, err.Error())
	}
	return MakeResultData(result)
}

func (app *BTApplication) QueryAccount(addr []byte) define.Result {
	address := ethcmn.HexToAddress(string(addr))

	if !app.stateDup.state.Exist(address) {
		return define.NewError(define.CodeType_UnknownAccount, address.Hex())
	}

	balance := app.stateDup.state.GetBalance(address)
	b := &define.BlockAccount{
		Addr:    address.Hex(),
		Balance: balance.String(),
	}
	return MakeResultData(b)
}

func MakeResultData(i interface{}) define.Result {
	jdata, err := json.Marshal(i)
	if err != nil {
		return define.NewError(define.CodeType_InternalError, err.Error())
	}
	return define.NewResultOK(jdata, "")
}
