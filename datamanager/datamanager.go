package datamanager

import (
	"github.com/axengine/btchain/config"
	"github.com/axengine/btchain/database"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"sync"
)

// DBCreator to create db instance
type DBCreator func(dbname string) database.Database

// DataManager data access between app and database
type DataManager struct {
	qdb       database.Database
	qNeedLock bool
	qLock     sync.Mutex
}

// NewDataManager create data manager
func NewDataManager(cfg *config.Config, logger *zap.Logger, dbc DBCreator) (*DataManager, error) {
	qdb := dbc("bt_query.db")

	qt, qi := qdb.GetInitSQLs()
	err := qdb.PrepareTables(qt, qi)
	if err != nil {
		return nil, err
	}

	dm := &DataManager{
		qdb: qdb,
	}
	switch cfg.DB.Type {
	case database.DBTypeSQLite3:
		dm.qNeedLock = true
	default:
		dm.qNeedLock = true
	}

	return dm, nil
}

// Close close all dbs
func (m *DataManager) Close() {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}
	if m.qdb != nil {
		m.qdb.Close()
		m.qdb = nil
	}
}

// QTxBegin start database transaction of qdb
func (m *DataManager) QTxBegin() error {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	return m.qdb.Begin()
}

// QTxCommit commit database transaction of qdb
func (m *DataManager) QTxCommit() error {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	return m.qdb.Commit()
}

// QTxRollback rollback database transaction of qdb
func (m *DataManager) QTxRollback() error {
	if m.qNeedLock {
		m.qLock.Lock()
		defer m.qLock.Unlock()
	}

	return m.qdb.Rollback()
}
