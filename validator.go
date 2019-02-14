package btchain

import (
	"encoding/base64"
	"github.com/axengine/btchain/define"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/pkg/errors"
	"github.com/tendermint/go-amino"
	"go.uber.org/zap"
	"golang.org/x/crypto/ed25519"
	"strconv"
	"sync/atomic"
)

var validatorsKey = []byte("validatorsKey")

// validators 说明
// InitChain时参数中的validatorset是创世配置文件中配置的内容，本方法只调一次，程序重启不再调用
// BeginBlock中的validator是空的，无法使用
// validator持久化到leveldb了，包含从创世以来所有的validator信息，power=0时不移除
// 更新注意：如果power double set0，TMCore会出现一致性错误，必须重新创世
//			 如果power不为0，如果网络不通也有可能出现问题，因此需要谨慎更新validator
type ValidatorMgr struct {
	flag       uint32              // 数据更新标识 atomic
	toUpdate   *define.Validator   // 待更新validator
	validators []*define.Validator // 整个区块应用实时的validator信息，包括power=0的validator,必须去重
	chainDb    ethdb.Database      // 持久化
}

func NewValidatorMgr(chainDb ethdb.Database) *ValidatorMgr {
	return &ValidatorMgr{
		chainDb: chainDb,
	}
}

func (p *ValidatorMgr) Get(pubkey string) *define.Validator {
	for _, v := range p.validators {
		if v.PubKey == pubkey {
			return v
		}
	}
	return nil
}

func (p *ValidatorMgr) Save(vl *define.Validator) {
	for _, v := range p.validators {
		if v.PubKey == vl.PubKey {
			v.Power = vl.Power
			return
		}
	}

	p.validators = append(p.validators, vl)
}

// Load 从磁盘加载validators信息
func (p *ValidatorMgr) Load() {
	buf, _ := p.chainDb.Get(validatorsKey)
	if len(buf) != 0 {
		var vls []*define.Validator
		if err := amino.UnmarshalBinaryBare(buf, &vls); err != nil {
			panic(err)
		}
		p.validators = vls
	}
}

// Dump 将Validator信息持久化
func (p *ValidatorMgr) Dump() {
	b, err := amino.MarshalBinaryBare(p.validators)
	if err != nil {
		panic(err)
	}
	if err = p.chainDb.Put(validatorsKey, b); err != nil {
		panic(err)
	}
}

func (p *ValidatorMgr) String() string {
	var str string
	for _, v := range p.validators {
		str = str + "[" + v.PubKey + ":" + strconv.Itoa(int(v.Power)) + "]"
	}
	return str
}

func (app *BTApplication) CheckValidatorUpdate(tx *define.Transaction) error {
	action := tx.Actions[0]
	if action == nil {
		return errors.New("CodeType_InvalidTx")
	}
	pubkey, err := base64.StdEncoding.DecodeString(action.Data)
	if err != nil {
		return err
	}
	if len(pubkey) != 32 {
		return errors.New("error validator pubkey")
	}
	power, err := strconv.Atoi(action.Memo)
	if err != nil {
		return err
	}

	sign := action.SignHex[:64]
	{
		if !ed25519.Verify(pubkey, pubkey, sign) {
			return errors.New("signature failed")
		}
	}

	v := app.valMgr.Get(action.Data)
	if v == nil {
		//不容许 设置一个不存在的validator power=0
		if power == 0 {
			return errors.New("set power=0 for unknow validator")
		}
	} else {
		// 不容许 设置一个power为0的validator power=0
		if v.Power == 0 && power == 0 {
			return errors.New("double set power=0 for validator")
		}
	}
	return nil
}

func (app *BTApplication) DeliverValidatorUpdate(tx *define.Transaction) error {
	action := tx.Actions[0]
	power, _ := strconv.Atoi(action.Memo)

	if atomic.LoadUint32(&app.valMgr.flag) == 1 {
		app.logger.Info("DeliverValidatorUpdate", zap.String("SP Flag", "== 1"))
		return errors.New("Locked")
	}

	app.valMgr.toUpdate = &define.Validator{
		PubKey: action.Data,
		Power:  int64(power),
	}

	atomic.StoreUint32(&app.valMgr.flag, 1)
	return nil
}
