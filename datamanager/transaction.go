package datamanager

import (
	"database/sql"
	"github.com/axengine/btchain/database"
	"github.com/axengine/btchain/define"
	"math/big"

	ethcmn "github.com/ethereum/go-ethereum/common"
)

func (m *DataManager) PrepareTransaction() (*sql.Stmt, error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	fields := []database.Feild{
		database.Feild{Name: "txhash"},
		database.Feild{Name: "blockheight"},
		database.Feild{Name: "blockhash"},
		database.Feild{Name: "actioncount"},
		database.Feild{Name: "actionid"},
		database.Feild{Name: "src"},
		database.Feild{Name: "dst"},
		database.Feild{Name: "nonce"},
		database.Feild{Name: "amount"},
		database.Feild{Name: "resultcode"},
		database.Feild{Name: "resultmsg"},
		database.Feild{Name: "createdat"},
		database.Feild{Name: "jdata"},
		database.Feild{Name: "memo"},
	}

	return m.qdb.Prepare(database.TableTransactions, fields)
}

func (m *DataManager) AddTransactionStmt(stmt *sql.Stmt, data *define.TransactionData) (err error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	fields := []database.Feild{
		database.Feild{Name: "txhash", Value: data.TxHash.Hex()},
		database.Feild{Name: "blockheight", Value: data.BlockHeight},
		database.Feild{Name: "blockhash", Value: data.BlockHash.Hex()},
		database.Feild{Name: "actioncount", Value: data.ActionCount},
		database.Feild{Name: "actionid", Value: data.ActionID},
		database.Feild{Name: "src", Value: data.Src.Hex()},
		database.Feild{Name: "dst", Value: data.Dst.Hex()},
		database.Feild{Name: "nonce", Value: data.Nonce},
		database.Feild{Name: "amount", Value: data.Amount.String()},
		database.Feild{Name: "resultcode", Value: data.ResultCode},
		database.Feild{Name: "resultmsg", Value: data.ResultMsg},
		database.Feild{Name: "createdat", Value: data.CreateAt},
		database.Feild{Name: "jdata", Value: data.JData},
		database.Feild{Name: "memo", Value: data.Memo},
	}
	_, err = m.qdb.Excute(stmt, fields)
	return err
}

// AddTransaction insert a tx record
func (m *DataManager) AddTransaction(data *define.TransactionData) (uint64, error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	fields := []database.Feild{
		database.Feild{Name: "txhash", Value: data.TxHash.Hex()},
		database.Feild{Name: "blockheight", Value: data.BlockHeight},
		database.Feild{Name: "blockhash", Value: data.BlockHash.Hex()},
		database.Feild{Name: "actioncount", Value: data.ActionCount},
		database.Feild{Name: "actionid", Value: data.ActionID},
		database.Feild{Name: "src", Value: data.Src.Hex()},
		database.Feild{Name: "dst", Value: data.Dst.Hex()},
		database.Feild{Name: "nonce", Value: data.Nonce},
		database.Feild{Name: "amount", Value: data.Amount.String()},
		database.Feild{Name: "resultcode", Value: data.ResultCode},
		database.Feild{Name: "resultmsg", Value: data.ResultMsg},
		database.Feild{Name: "createdat", Value: data.CreateAt},
		database.Feild{Name: "jdata", Value: data.JData},
		database.Feild{Name: "memo", Value: data.Memo},
	}

	sqlRes, err := m.qdb.Insert(database.TableTransactions, fields)
	if err != nil {
		return 0, err
	}

	id, err := sqlRes.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

// QuerySingleTx query single tx record
func (m *DataManager) QuerySingleTx(txhash *ethcmn.Hash) (*define.TransactionData, error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	where := []database.Where{
		database.Where{Name: "txhash", Value: txhash.Hex()},
	}

	var result []database.TxData
	err := m.qdb.SelectRows(database.TableTransactions, where, nil, nil, &result)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}
	if len(result) > 1 {
		// panic ?
	}

	r := result[0]
	td := define.TransactionData{
		TxID:        r.TxID,
		TxHash:      ethcmn.HexToHash(r.TxHash),
		BlockHeight: r.BlockHeight,
		BlockHash:   ethcmn.HexToHash(r.BlockHash),
		ActionCount: r.ActionCount,
		ActionID:    r.ActionID,
		Src:         ethcmn.HexToAddress(r.Src),
		Dst:         ethcmn.HexToAddress(r.Dst),
		Nonce:       r.Nonce,
		Amount:      Str2Big(r.Amount),
		ResultCode:  r.ResultCode,
		ResultMsg:   r.ResultMsg,
		CreateAt:    r.CreateAt,
		JData:       r.JData,
		Memo:        r.Memo,
	}
	return &td, nil
}

func (m *DataManager) QueryTxByHash(txhash *ethcmn.Hash) ([]define.TransactionData, error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	where := []database.Where{
		database.Where{Name: "txhash", Value: txhash.Hex()},
	}

	orderT, err := database.MakeOrder("", "txid")
	if err != nil {
		return nil, err
	}

	var result []database.TxData
	if err := m.qdb.SelectRows(database.TableTransactions, where, orderT, nil, &result); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var res []define.TransactionData
	for _, r := range result {
		td := define.TransactionData{
			TxID:        r.TxID,
			TxHash:      ethcmn.HexToHash(r.TxHash),
			BlockHeight: r.BlockHeight,
			BlockHash:   ethcmn.HexToHash(r.BlockHash),
			ActionCount: r.ActionCount,
			ActionID:    r.ActionID,
			Src:         ethcmn.HexToAddress(r.Src),
			Dst:         ethcmn.HexToAddress(r.Dst),
			Nonce:       r.Nonce,
			Amount:      Str2Big(r.Amount),
			ResultCode:  r.ResultCode,
			ResultMsg:   r.ResultMsg,
			CreateAt:    r.CreateAt,
			JData:       r.JData,
			Memo:        r.Memo,
		}

		res = append(res, td)
	}
	return res, nil
}

func Str2Big(num string) *big.Int {
	n := new(big.Int)
	n.SetString(num, 0)
	return n
}

// QueryAccountTxs query account's tx records
func (m *DataManager) QueryAccountTxs(accid *ethcmn.Address, direction uint8, cursor, limit uint64, order string) ([]define.TransactionData, error) {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	where := []database.Where{
		database.Where{Name: "1", Value: 1},
	}
	if accid != nil {
		if direction == 1 {
			where = append(where, database.Where{Name: "dst", Value: accid.Hex()})
		} else {
			where = append(where, database.Where{Name: "src", Value: accid.Hex()})
		}
	}
	orderT, err := database.MakeOrder(order, "txid")
	if err != nil {
		return nil, err
	}
	paging := database.MakePaging("txid", cursor, limit)

	var result []database.TxData
	err = m.qdb.SelectRows(database.TableTransactions, where, orderT, paging, &result)
	if err != nil {
		return nil, err
	}

	var res []define.TransactionData
	for _, r := range result {
		td := define.TransactionData{
			TxID:        r.TxID,
			TxHash:      ethcmn.HexToHash(r.TxHash),
			BlockHeight: r.BlockHeight,
			BlockHash:   ethcmn.HexToHash(r.BlockHash),
			ActionCount: r.ActionCount,
			ActionID:    r.ActionID,
			Src:         ethcmn.HexToAddress(r.Src),
			Dst:         ethcmn.HexToAddress(r.Dst),
			Nonce:       r.Nonce,
			Amount:      Str2Big(r.Amount),
			ResultCode:  r.ResultCode,
			ResultMsg:   r.ResultMsg,
			CreateAt:    r.CreateAt,
			JData:       r.JData,
			Memo:        r.Memo,
		}

		res = append(res, td)
	}

	return res, nil
}

// QueryAllTxs query all tx records
func (m *DataManager) QueryAllTxs(cursor, limit uint64, order string) ([]define.TransactionData, error) {
	return m.QueryAccountTxs(nil, 0, cursor, limit, order)
}
