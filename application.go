package btchain

import (
	"github.com/axengine/btchain/config"
	"github.com/axengine/btchain/database"
	"github.com/axengine/btchain/database/basesql"
	"github.com/axengine/btchain/datamanager"
	"github.com/axengine/btchain/define"
	"github.com/axengine/btchain/log"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/tendermint/go-amino"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"go.uber.org/zap"
	"math/big"
	"path"
	"path/filepath"
	"sync"
)

var _ abcitypes.Application = (*BTApplication)(nil)

type BTApplication struct {
	abcitypes.BaseApplication

	currentHeader *define.AppHeader
	tempHeader    *define.AppHeader // for executing tx

	blockExeInfo *blockExeInfo

	stateDup *stateDup
	chainDb  ethdb.Database
	dataM    *datamanager.DataManager
	logger   *zap.Logger
	cfg      *config.Config

	valMgr *ValidatorMgr
}

type blockExeInfo struct {
	txDatas        []*define.TransactionData
	inflationOccur bool
}

type stateDup struct {
	lock  sync.RWMutex
	state *state.StateDB
}

func NewBTApplication() *BTApplication {
	app := new(BTApplication)
	app.init()
	return app
}

const (
	LDatabaseCache   = 128
	LDatabaseHandles = 1024
)

var (
	EmptyTrieRoot = ethcmn.HexToHash("0000000000000000000000000000000000000000000000000000000000000000")
	lastBlockKey  = []byte("lastblock")
)

func (app *BTApplication) init() {
	var err error
	var logger *zap.Logger
	cfg := config.New()
	if err := cfg.Init("./config.toml"); err != nil {
		panic("On init yaml:" + err.Error())
	}
	app.cfg = cfg

	logger = log.Initialize("file", cfg.Log.Env, path.Join(cfg.Log.Path, "node.debug.log"))
	app.logger = logger

	//level db
	if app.chainDb, err = OpenDatabase(cfg.DB.Path, "chaindata", LDatabaseCache, LDatabaseHandles); err != nil {
		panic(err)
	}

	//加载validator信息
	app.valMgr = NewValidatorMgr(app.chainDb)
	app.valMgr.Load()
	app.logger.Info("init", zap.String("validators", app.valMgr.String()))

	//加载最新区块信息
	lastBlock := app.LoadLastBlock()
	trieRoot := EmptyTrieRoot
	if len(lastBlock.StateRoot) > 0 {
		trieRoot = ethcmn.BytesToHash(lastBlock.StateRoot)
	}

	//初始化 state
	app.stateDup = new(stateDup)
	if app.stateDup.state, err = state.New(trieRoot, state.NewDatabase(app.chainDb)); err != nil {
		panic(err)
	}

	//创世块
	if trieRoot == EmptyTrieRoot {
		addr := ethcmn.HexToAddress(cfg.Genesis.Account)
		app.stateDup.state.CreateAccount(addr)
		amount, _ := new(big.Int).SetString(cfg.Genesis.Amount, 10)
		app.stateDup.state.AddBalance(addr, amount)

		logger.Info("Genesis", zap.String("account", addr.Hex()), zap.String("amount", amount.String()))

		root := app.stateDup.state.IntermediateRoot(false)
		_, err := app.stateDup.state.Commit(false)
		if err != nil {
			panic(err)
		}

		if err := app.stateDup.state.Database().TrieDB().Commit(root, true); err != nil {
			panic(err)
		}
	}

	//sqlite3
	app.dataM, err = datamanager.NewDataManager(cfg, logger, func(dbname string) database.Database {
		dbi := &basesql.Basesql{}
		err := dbi.Init(dbname, cfg, logger)
		if err != nil {
			panic(err)
		}
		return dbi
	})
	if err != nil {
		panic(err)
	}

	//初始化数据
	app.blockExeInfo = &blockExeInfo{}
	app.currentHeader = &define.AppHeader{
		PrevHash:  ethcmn.BytesToHash(lastBlock.AppHash),
		StateRoot: trieRoot,
		Height:    lastBlock.Height,
	}
	app.tempHeader = app.currentHeader

	logger.Info("BT Application init ok", zap.Uint64("height", app.currentHeader.Height))
}

func OpenDatabase(datadir string, name string, cache int, handles int) (ethdb.Database, error) {
	return ethdb.NewLDBDatabase(filepath.Join(datadir, name), cache, handles)
}

type LastBlockInfo struct {
	Height    uint64
	StateRoot []byte
	AppHash   []byte
	PrevHash  []byte
}

func (app *BTApplication) LoadLastBlock() (lastBlock LastBlockInfo) {
	buf, _ := app.chainDb.Get(lastBlockKey)
	if len(buf) != 0 {
		if err := amino.UnmarshalBinaryBare(buf, &lastBlock); err != nil {
			panic(err)
		}
	}

	return lastBlock
}

func (app *BTApplication) SaveLastBlock(appHash []byte, header *define.AppHeader) {
	lastBlock := LastBlockInfo{
		Height:    header.Height,
		StateRoot: header.StateRoot.Bytes(),
		AppHash:   appHash,
		PrevHash:  header.PrevHash.Bytes(),
	}

	buf, err := amino.MarshalBinaryBare(&lastBlock)
	if err != nil {
		panic(err)
	}

	if err := app.chainDb.Put(lastBlockKey, buf); err != nil {
		panic(err)
	}
}

func (app *BTApplication) SaveDBData() error {
	// begin dbtx
	err := app.dataM.QTxBegin()
	if err != nil {
		return err
	}

	stmt, err := app.dataM.PrepareTransaction()
	if err != nil {
		app.dataM.QTxRollback()
		return err
	}
	for _, v := range app.blockExeInfo.txDatas {
		err = app.dataM.AddTransactionStmt(stmt, v)
		if err != nil {
			app.dataM.QTxRollback()
			return err
		}
	}
	stmt.Close()

	// commit dbtx
	err = app.dataM.QTxCommit()
	if err != nil {
		return err
	}

	return nil
}

func (app *BTApplication) DBSwitch() {
	//c := cron.New()
	//spec := "0 0 0 * * ?"
	//err := c.AddFunc(spec, func() {
	//	app.dBSwitch()
	//})
	//if err != nil {
	//	app.logger.Error("cron", zap.Error(err))
	//}
	//
	//c.Start()
}

func (app *BTApplication) dBSwitch() {
	dataM, err := datamanager.NewDataManager(app.cfg, app.logger, func(dbname string) database.Database {
		dbi := &basesql.Basesql{}
		err := dbi.Init(dbname, app.cfg, app.logger)
		if err != nil {
			panic(err)
		}
		return dbi
	})
	if err != nil {
		app.logger.Error("DBSwitch", zap.Error(err))
		return
	}
	app.dataM = dataM
}
