package basesql

// GetInitSQLs get database initialize sqls
//	qt  sqls to create query tables
//	qi  sqls to create query table-indexs
func (bs *Basesql) GetInitSQLs() (qt, qi []string) {
	qt = []string{
		createEffectSQL,
		createLedgerSQL,
		creatActionSQL,
		createTransactionSQL,
	}
	qi = createQIndex
	return
}

/*
	txid        INTEGER        PRIMARY KEY AUTOINCREMENT,
	txhash      VARCHAR( 66 )  NOT NULL,
	blockheight INTEGER        NOT NULL,
	blockhash   VARCHAR( 64 )  NOT NULL,
	actioncount INT            NOT NULL,
	actionid    INT            NOT NULL,
	uid         VARCHAR( 66 )  NOT NULL,
	relateduid  VARCHAR( 66 )  NOT NULL,
	direction   INT            NOT NULL,
	nonce       INTEGER        NOT NULL,
	amount      TEXT           NOT NULL,
	resultcode  INT            NOT NULL,
	resultmsg   TEXT           NOT NULL,
	createdat   INTEGER		   NOT NULL,
	jdata       TEXT           NOT NULL,
	memo        TEXT
*/
var (
	createQIndex = []string{
		// indexs for table ledger
		"CREATE INDEX IF NOT EXISTS ledgersequence ON ledgerheaders (sequence)",

		// indexs for table transactions
		"CREATE INDEX IF NOT EXISTS idx_txhash ON transactions (txhash)",
		"CREATE INDEX IF NOT EXISTS idx_height ON transactions (blockheight)",
		"CREATE INDEX IF NOT EXISTS idx_src ON transactions (src)",
		"CREATE INDEX IF NOT EXISTS idx_dst ON transactions (dst)",
		"CREATE INDEX IF NOT EXISTS idx_createdat ON transactions (createdat)",

		// indexs for table actions
		"CREATE INDEX IF NOT EXISTS actionstxhash ON actions (txhash)",
		"CREATE INDEX IF NOT EXISTS actionsfromaccount ON actions (fromaccount)",
		"CREATE INDEX IF NOT EXISTS actionstoaccount ON actions (toaccount)",
		"CREATE INDEX IF NOT EXISTS actionscreateat ON actions (createat)",

		// indexs for table effects
		"CREATE INDEX IF NOT EXISTS effectstxhash ON effects (txhash)",
		"CREATE INDEX IF NOT EXISTS effectsaccount ON effects (account)",
		"CREATE INDEX IF NOT EXISTS effectscreateat ON effects (createat)",
	}
)

const (
	creatActionSQL = `CREATE TABLE IF NOT EXISTS actions
	(
		actionid			INTEGER	PRIMARY KEY	AUTOINCREMENT,
		typei				INT			NOT NULL,
		type				VARCHAR(32)	NOT NULL,
		ledgerseq			INT			NOT NULL,
		txhash				VARCHAR(64)	NOT NULL,
		fromaccount			VARCHAR(66),			-- only used in payment
		toaccount			VARCHAR(66),			-- only used in payment
		createat			INT			NOT NULL,
		jdata				TEXT		NOT NULL
	);`

	createEffectSQL = `CREATE TABLE IF NOT EXISTS effects
	(
		effectid          	INTEGER	PRIMARY KEY	AUTOINCREMENT,
		typei				INT,
		type				VARCHAR(32)	NOT NULL,
		ledgerseq			INT			NOT NULL,
		txhash				VARCHAR(64)	NOT NULL,
		actionid			INT			NOT NULL,
		account				VARCHAR(66)	NOT NULL,
		createat			INT			NOT NULL,
		jdata				TEXT		NOT NULL
	);`

	createLedgerSQL = `CREATE TABLE IF NOT EXISTS ledgerheaders
    (
		ledgerid			INTEGER	PRIMARY KEY	AUTOINCREMENT,
		sequence			TEXT 		UNIQUE,
		hash				VARCHAR(64) NOT NULL,
		prevhash			VARCHAR(64) NOT NULL,
		stateroot			VARCHAR(64) NOT NULL,
		transactioncount	INT			NOT NULL,
		operactioncount		INT 		NOT NULL,
		closedat			TIMESTAMP 	NOT NULL,
		totalcoins			TEXT	 	NOT NULL,
		feepool				TEXT 		NOT NULL,
		basefee				TEXT 		NOT NULL,
		inflationseq		INT			NOT NULL,
		basereserve			TEXT 		NOT NULL,
		maxtxsetsize		INT 		NOT NULL
	);`

	createTransactionSQL = `CREATE TABLE IF NOT EXISTS transactions
	(
		txid        INTEGER        PRIMARY KEY AUTOINCREMENT,
		txhash      VARCHAR( 66 )  NOT NULL,
		blockheight INTEGER        NOT NULL,
		blockhash   VARCHAR( 64 )  NOT NULL,
		actioncount INT            NOT NULL,
		actionid    INT            NOT NULL,
		src        VARCHAR( 42 )  NOT NULL,
		dst          VARCHAR( 42 )  NOT NULL,
		nonce       INTEGER        NOT NULL,
		amount      TEXT           NOT NULL,
		resultcode  INT            NOT NULL,
		resultmsg   TEXT           NOT NULL,
		createdat   INTEGER		   NOT NULL,
		jdata       TEXT           NOT NULL,
		memo        TEXT
	);`
)
