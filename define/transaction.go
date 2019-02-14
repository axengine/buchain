package define

import (
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"math/big"
)

// TransactionData 交易表
// 一次HTTP请求算一笔交易，拥有唯一的交易HASH
// 单笔交易可以携带多个ACTION，每次ACTION有一个ID，ACTIONID由调用方指定，注意JSON Array没有顺序
// 每个ACTION保护出让方和受让方，“创建帐户”不是ACTION
// 单笔交易多个ACTION时可以部分失败，部分失败返回{成功的action列表，失败的action列表}
type TransactionData struct {
	TxID        uint64         `json:"tx_id"`        //交易ID
	TxHash      ethcmn.Hash    `json:"tx_hash"`      //交易HASH - 重复交易不会被处理
	BlockHeight uint64         `json:"block_height"` //区块高度
	BlockHash   ethcmn.Hash    `json:"block_hash"`   //区块HASH
	ActionCount uint32         `json:"action_count"` //一笔交易多个action
	ActionID    uint32         `json:"action_id"`    //action id
	Src         ethcmn.Address `json:"src"`          //用户ID (if dir==0,uid 表示转入方，否则表示转出方)
	Dst         ethcmn.Address `json:"dst"`          //关联的用户ID
	Nonce       uint64         `json:"nonce"`        //对应操作源帐户(转出方)NONCE
	Amount      *big.Int       `json:"amount"`       //金额
	ResultCode  uint           `json:"result_code"`  //应答码 0-success
	ResultMsg   string         `json:"result_msg"`   //应答消息
	CreateAt    uint64         `json:"created_at"`   //入库时间
	JData       string         `json:"jdata"`        //数据部分 建议JSON序列化
	Memo        string         `json:"memo"`         //交易备注
}

// 通过Transaction 计算txhash
type Transaction struct {
	Type    uint8     //0-默认 1-special op
	Actions []*Action //有序的action 按照ID ASC序
}

func (p *Transaction) String() string {
	return fmt.Sprintf("tx:%v actions:%v",
		p.SigHash().Hex(), p.Len())
}

type Action struct {
	ID        uint8          // 最大支持255笔交易
	CreatedAt uint64         // 时间
	Src       ethcmn.Address //
	Dst       ethcmn.Address //
	Amount    *big.Int       //
	Data      string         // 用户数据
	Memo      string         // 备注
	SignHex   [65]byte       // 签名 65 bytes
}

func (p Transaction) Len() int {
	return len(p.Actions)
}

func (p Transaction) Less(i, j int) bool {
	return p.Actions[i].ID < p.Actions[j].ID
}

func (p Transaction) Swap(i, j int) {
	p.Actions[i], p.Actions[j] = p.Actions[j], p.Actions[i]
}
