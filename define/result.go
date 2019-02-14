package define

import (
	"encoding/json"
)

const (
	CodeType_OK uint32 = 0
	// General response codes, 0 ~ 99
	CodeType_InternalError     uint32 = 1  //内部错误
	CodeType_EncodingError     uint32 = 2  //编解码错误
	CodeType_BadNonce          uint32 = 3  //nonce错误
	CodeType_Unauthorized      uint32 = 4  //未授权
	CodeType_InsufficientFunds uint32 = 5  //资金不足
	CodeType_UnknownRequest    uint32 = 6  //未知请求
	CodeType_InvalidTx         uint32 = 7  //交易不合法
	CodeType_UnknownAccount    uint32 = 8  //未知帐户
	CodeType_AccountExist      uint32 = 9  //帐户已存在
	CodeType_AccountNotFound   uint32 = 10 //帐户不存在
	CodeType_OutOfOrder        uint32 = 11 //action顺序错误
	CodeType_UnknownError      uint32 = 12 //未知错误
	CodeType_SignerFaild       uint32 = 13 //签名错误
)

// CONTRACT: a zero Result is OK.
type Result struct {
	Code uint32 `json:"Code"`
	//Data []byte   `json:"Data"`
	Data []byte `json:"Data"`
	Log  string `json:"Log"` // Can be non-deterministic
}

type NewRoundResult struct {
}

type CommitResult struct {
	AppHash      []byte
	ReceiptsHash []byte
}

type ExecuteInvalidTx struct {
	Bytes []byte
	Error error
}

type ExecuteResult struct {
	ValidTxs   [][]byte
	InvalidTxs []ExecuteInvalidTx
	Error      error
}

func NewResult(code uint32, data []byte, log string) Result {
	return Result{
		Code: code,
		Data: data,
		Log:  log,
	}
}

func (res Result) ToJSON() string {
	j, err := json.Marshal(res)
	if err != nil {
		return res.Log
	}
	return string(j)
}

func (res *Result) FromJSON(j string) *Result {
	err := json.Unmarshal([]byte(j), res)
	if err != nil {
		res.Code = CodeType_InternalError
		res.Log = j
	}
	return res
}

func (res Result) IsOK() bool {
	return res.Code == CodeType_OK
}

func (res Result) IsErr() bool {
	return res.Code != CodeType_OK
}

func (res Result) Error() string {
	// return fmt.Sprintf("{code:%v, data:%X, log:%v}", res.Code, res.Data, res.Log)
	return res.ToJSON()
}

func (res Result) String() string {
	// return fmt.Sprintf("{code:%v, data:%X, log:%v}", res.Code, res.Data, res.Log)
	return res.ToJSON()
}

func (res Result) PrependLog(log string) Result {
	return Result{
		Code: res.Code,
		Data: res.Data,
		Log:  log + ";" + res.Log,
	}
}

func (res Result) AppendLog(log string) Result {
	return Result{
		Code: res.Code,
		Data: res.Data,
		Log:  res.Log + ";" + log,
	}
}

func (res Result) SetLog(log string) Result {
	return Result{
		Code: res.Code,
		Data: res.Data,
		Log:  log,
	}
}

func (res Result) SetData(data []byte) Result {
	return Result{
		Code: res.Code,
		Data: data,
		Log:  res.Log,
	}
}

//----------------------------------------

// NOTE: if data == nil and log == "", same as zero Result.
func NewResultOK(data []byte, log string) Result {
	return Result{
		Code: CodeType_OK,
		Data: data,
		Log:  log,
	}
}

func NewError(code uint32, log string) Result {
	return Result{
		Code: code,
		Log:  log,
	}
}
