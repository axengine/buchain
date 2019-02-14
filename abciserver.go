package btchain

import (
	"encoding/base64"
	"fmt"
	"github.com/axengine/btchain/define"
	"github.com/axengine/btchain/version"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	abciversion "github.com/tendermint/tendermint/version"
	"go.uber.org/zap"
	"sort"
	"strconv"
	"sync/atomic"
)

func (app *BTApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	app.logger.Info("======>>InitChain", zap.String("chainId", req.ChainId), zap.Time("genis", req.Time))
	for _, v := range req.Validators {
		var validator define.Validator
		validator.PubKey = base64.StdEncoding.EncodeToString(v.PubKey.Data)
		validator.Power = v.Power
		app.valMgr.Save(&validator)
	}
	app.valMgr.Dump()
	app.logger.Info("======>>InitChain", zap.String("validators", app.valMgr.String()))
	return abcitypes.ResponseInitChain{}
}

// Info
// TM Core启动时会查询chain信息 需要返回lastBlock相关信息，否则会从第一个块replay
func (app *BTApplication) Info(req abcitypes.RequestInfo) (resInfo abcitypes.ResponseInfo) {
	if app.currentHeader.Height != 0 {
		resInfo.LastBlockHeight = int64(app.currentHeader.Height)
		resInfo.LastBlockAppHash = app.currentHeader.PrevHash.Bytes()
	}
	resInfo.Data = fmt.Sprintf("{\"size\":%v}", app.currentHeader.Height)
	resInfo.Version = abciversion.ABCIVersion
	resInfo.AppVersion = version.APPVersion
	app.logger.Info("ABCI Info", zap.Uint64("height", app.currentHeader.Height), zap.String("PrevHash", app.currentHeader.PrevHash.Hex()))
	return resInfo
}

// CheckTx
// 初步检查，如果check失败，将不会被打包
func (app *BTApplication) CheckTx(tx []byte) abcitypes.ResponseCheckTx {
	var t define.Transaction
	if err := rlp.DecodeBytes(tx, &t); err != nil {
		app.logger.Warn("rlp unmarshal", zap.Error(err), zap.ByteString("tx", tx))
		return abcitypes.ResponseCheckTx{Code: define.CodeType_EncodingError, Log: "CodeTypeEncodingError"}
	}
	sort.Sort(t)
	app.logger.Debug("ABCI CheckTx", zap.String("tx", t.String()))

	if t.Type == 1 {
		if err := app.CheckValidatorUpdate(&t); err != nil {
			return abcitypes.ResponseCheckTx{Code: define.CodeType_InvalidTx, Log: err.Error()}
		}
		return abcitypes.ResponseCheckTx{Code: define.CodeType_OK, Data: t.SigHash().Bytes()}
	}

	//检查每个操作是否合法
	for i, action := range t.Actions {
		if i != int(action.ID) {
			app.logger.Warn("ABCI CheckTx", zap.String("err", "CodeType_OutOfOrder:"+strconv.Itoa(int(action.ID))))
			return abcitypes.ResponseCheckTx{Code: define.CodeType_OutOfOrder, Log: "CodeType_OutOfOrder:" + strconv.Itoa(int(action.ID))}
		}
		if !app.stateDup.state.Exist(action.Src) {
			app.logger.Warn("ABCI CheckTx", zap.String("err", "CodeType_AccountNotFound:"+action.Src.Hex()))
			return abcitypes.ResponseCheckTx{Code: define.CodeType_AccountNotFound, Log: "CodeType_AccountNotFound:" + action.Src.Hex()}
		}

		balance := app.stateDup.state.GetBalance(action.Src)
		if balance.Cmp(action.Amount) < 0 {
			app.logger.Warn("ABCI CheckTx", zap.String("err", "CodeType_InsufficientFunds:"+balance.String()))
			return abcitypes.ResponseCheckTx{Code: define.CodeType_InsufficientFunds, Log: "CodeType_InsufficientFunds:" + balance.String()}
		}
	}

	//检查签名是否合法
	if err := t.CheckSig(); err != nil {
		app.logger.Warn("ABCI CheckTx", zap.String("err", "CodeType_SignerFaild:"+err.Error()))
		return abcitypes.ResponseCheckTx{Code: define.CodeType_SignerFaild, Log: "CodeType_SignerFaild:" + err.Error()}
	}

	app.logger.Info("ABCI CheckTx", zap.String("hash", t.SigHash().Hex()))
	return abcitypes.ResponseCheckTx{GasWanted: 1, Data: t.SigHash().Bytes()}
}

// BeginBlock
// 区块开始，记录区块高度和hash
func (app *BTApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	app.logger.Debug("ABCI BeginBlock", zap.Int64("height", req.Header.Height), zap.String("hash", ethcmn.Bytes2Hex(req.Hash)))
	//log.Println("===============>", req.ByzantineValidators)
	//var tmp string
	//for _, v := range req.ByzantineValidators {
	//	address := base64.StdEncoding.EncodeToString(v.Validator.Address)
	//	tmp = tmp + "=>[" + address + "]:[" + strconv.Itoa(int(v.Validator.Power)) + "]"
	//}
	//app.logger.Debug("ABCI BeginBlock", zap.String("validators", tmp))
	app.tempHeader.Height = uint64(req.Header.Height)
	app.tempHeader.BlockHash = ethcmn.BytesToHash(req.Hash)
	return abcitypes.ResponseBeginBlock{}
}

func (app *BTApplication) DeliverTx(tx []byte) abcitypes.ResponseDeliverTx {
	var (
		t define.Transaction
	)
	if err := rlp.DecodeBytes(tx, &t); err != nil {
		app.logger.Warn("rlp unmarshal", zap.Error(err), zap.ByteString("tx", tx))
		return abcitypes.ResponseDeliverTx{Code: define.CodeType_EncodingError, Log: "CodeType_EncodingError"}
	}
	sort.Sort(t)
	app.logger.Info("ABCI DeliverTx", zap.String("tx", t.String()))

	if t.Type == 1 {
		if err := app.DeliverValidatorUpdate(&t); err != nil {
			return abcitypes.ResponseDeliverTx{Code: define.CodeType_InvalidTx, Log: err.Error()}
		}
		return abcitypes.ResponseDeliverTx{Code: define.CodeType_OK, Data: t.SigHash().Bytes()}
	}

	//创建快照
	stateSnapshot := app.stateDup.state.Snapshot()

	app.tempHeader.TxCount = app.tempHeader.TxCount + 1
	app.tempHeader.OpCount = app.tempHeader.OpCount + uint64(t.Len())

	txHash := t.SigHash()
	actionCount := t.Len()
	for _, action := range t.Actions {
		//自动创建to
		if !app.stateDup.state.Exist(action.Dst) {
			app.stateDup.state.CreateAccount(action.Dst)
		}

		//取nonce
		nonce := app.stateDup.state.GetNonce(action.Src)

		//必须再次校验余额
		balance := app.stateDup.state.GetBalance(action.Src)
		if balance.Cmp(action.Amount) < 0 {
			app.stateDup.state.RevertToSnapshot(stateSnapshot)
			app.logger.Warn("ABCI DeliverTx", zap.String("err", "CodeType_InsufficientFunds"), zap.String("src", action.Src.Hex()), zap.String("amount", balance.String()))
			return abcitypes.ResponseDeliverTx{Code: define.CodeType_InsufficientFunds, Log: "not enough money"}
		}

		//资金操作
		app.stateDup.state.SubBalance(action.Src, action.Amount)
		app.stateDup.state.AddBalance(action.Dst, action.Amount)
		app.stateDup.state.SetNonce(action.Src, nonce+1)

		var txData define.TransactionData
		txData.TxHash = txHash
		txData.BlockHeight = app.tempHeader.Height
		txData.BlockHash = app.tempHeader.BlockHash
		txData.ActionCount = uint32(actionCount)
		txData.ActionID = uint32(action.ID)
		txData.Src = action.Src
		txData.Dst = action.Dst
		txData.Nonce = nonce
		txData.Amount = action.Amount
		txData.CreateAt = action.CreatedAt
		txData.JData = action.Data
		txData.Memo = action.Memo

		app.blockExeInfo.txDatas = append(app.blockExeInfo.txDatas, &txData)
	}

	return abcitypes.ResponseDeliverTx{Data: txHash.Bytes()}
}

func (app *BTApplication) commitState() (ethcmn.Hash, error) {
	var (
		stateRoot ethcmn.Hash
		err       error
	)

	app.stateDup.lock.Lock()
	defer app.stateDup.lock.Unlock()

	// 更新stateRoot
	stateRoot = app.stateDup.state.IntermediateRoot(false)
	if _, err = app.stateDup.state.Commit(false); err != nil {
		return stateRoot, err
	}
	if err = app.stateDup.state.Database().TrieDB().Commit(stateRoot, true); err != nil {
		return stateRoot, err
	}

	return stateRoot, nil
}

func (app *BTApplication) Commit() abcitypes.ResponseCommit {
	var (
		stateRoot ethcmn.Hash
		err       error
	)

	// 更新stateRoot
	stateRoot, err = app.commitState()
	if err != nil {
		app.logger.Error("ABCI Commit state commit", zap.Error(err))
		app.SaveLastBlock(app.currentHeader.Hash(), app.currentHeader)
		return abcitypes.ResponseCommit{Data: app.currentHeader.Hash()}
	}

	//计算 appHash
	app.tempHeader.StateRoot = stateRoot
	appHash := app.tempHeader.Hash()

	//保存lastblock
	app.SaveLastBlock(appHash, app.tempHeader)

	//更新currentHeader 保护现场
	app.currentHeader = app.tempHeader

	//SQLITE3保存记录
	if err := app.SaveDBData(); err != nil {
		app.logger.Error("ABCI SaveDBData", zap.Error(err), zap.Uint64("height", app.tempHeader.Height))
	}

	//清理区块执行现场
	app.blockExeInfo = &blockExeInfo{}
	app.tempHeader = &define.AppHeader{}

	//  -------------准备新区块---------
	//下个区块需要上个的hash
	app.tempHeader.PrevHash = ethcmn.BytesToHash(appHash)

	app.logger.Info("ABCI Commit", zap.String("hash", ethcmn.BytesToHash(appHash).Hex()), zap.Uint64("height", app.currentHeader.Height))
	return abcitypes.ResponseCommit{Data: appHash}
}

func (app *BTApplication) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	app.logger.Debug("ABCI Query:", zap.String("path", reqQuery.Path), zap.String("data", string(reqQuery.Data)),
		zap.Int64("Height", reqQuery.Height),
		zap.Bool("Prove", reqQuery.Prove))

	switch reqQuery.Path {
	case QUERY_TX:
		result := app.QueryTx(reqQuery.Data)
		b, err := rlp.EncodeToBytes(&result)
		if err != nil {
			return abcitypes.ResponseQuery{Code: define.CodeType_EncodingError, Log: err.Error()}
		}
		return abcitypes.ResponseQuery{Value: b}
	case QUERY_ACCOUNT:
		result := app.QueryAccount(reqQuery.Data)
		b, err := rlp.EncodeToBytes(&result)
		if err != nil {
			return abcitypes.ResponseQuery{Code: define.CodeType_EncodingError, Log: err.Error()}
		}
		return abcitypes.ResponseQuery{Value: b}
	default:
		app.logger.Warn("ABCI Query", zap.String("code", "CodeType_UnknownRequest"))
		return abcitypes.ResponseQuery{Code: define.CodeType_UnknownRequest, Log: "CodeType_UnknownRequest"}
	}
	return abcitypes.ResponseQuery{Value: []byte(fmt.Sprintf("%v", app.currentHeader.PrevHash))}
}

// EndBlock
// https://tendermint.com/docs/spec/abci/apps.html#validator-updates
func (app *BTApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	if atomic.LoadUint32(&app.valMgr.flag) == 1 {
		defer atomic.StoreUint32(&app.valMgr.flag, 0)

		toUpdate := app.valMgr.toUpdate
		if toUpdate == nil {
			return abcitypes.ResponseEndBlock{}
		}
		pubkey, _ := base64.StdEncoding.DecodeString(toUpdate.PubKey)

		app.logger.Info("EndBlock set validator", zap.String("pubkey", toUpdate.PubKey), zap.Int64("power", toUpdate.Power))
		var validator abcitypes.ValidatorUpdate
		validator.Power = toUpdate.Power
		validator.PubKey.Type = "ed25519"
		validator.PubKey.Data = pubkey

		//先更新本地的validator
		{
			v := app.valMgr.Get(toUpdate.PubKey)
			if v == nil {
				app.valMgr.Save(toUpdate)
			} else {
				v.Power = toUpdate.Power
			}
			app.valMgr.Dump()
			app.logger.Info("EndBlock", zap.String("validators", app.valMgr.String()))
		}

		var ValidatorUpdates []abcitypes.ValidatorUpdate
		ValidatorUpdates = append(ValidatorUpdates, validator)
		return abcitypes.ResponseEndBlock{ValidatorUpdates: ValidatorUpdates}
	}
	return abcitypes.ResponseEndBlock{}
}
