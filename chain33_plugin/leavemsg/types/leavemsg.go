package types

import (
	"fmt"
	"github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/chain33/types"
	"encoding/json"
	"github.com/33cn/chain33/common/address"
)

var (
	llog       = log15.New("module", LeaveMsgX)
	actionName = map[string]int32{
		"Send":  LeavemsgSendAction,
		"Query": LeavemsgQueryAction,
	}
)

func init() {
	types.AllowUserExec = append(types.AllowUserExec, []byte(LeaveMsgX))
	types.RegistorExecutor(LeaveMsgX, NewType())
	types.RegisterDappFork(LeaveMsgX, "Enable", 0)
}

func getRealExecName(paraName string) string {
	return types.ExecName(paraName + LeaveMsgX)
}

type LeavemsgType struct {
	types.ExecTypeBase
}

func NewType() *LeavemsgType {
	c := &LeavemsgType{}
	c.SetChild(c)
	return c
}

func (l *LeavemsgType) GetPayload() types.Message {
	return &LeaveMsgAction{}
}

func (l *LeavemsgType) GetLogMap() map[int64]*types.LogInfo {
	return map[int64]*types.LogInfo{}
}

func (l *LeavemsgType) GetTypeMap() map[string]int32 {
	return actionName
}

func (leavemsg LeavemsgType) CreateTx(action string, message json.RawMessage) (*types.Transaction, error) {
	if action == "Send" {
		fmt.Println("In leavemsg create tx, message is :", string(message))
		var parm LeaveMsgSendTx
		err := json.Unmarshal(message, &parm)
		if err != nil {
			fmt.Println(err)
			return nil, types.ErrActionNotSupport
		}
		return CreateRawLeaveMsgSendTx(&parm)
	}
	return nil, nil
}

func CreateRawLeaveMsgSendTx(parm *LeaveMsgSendTx) (*types.Transaction, error) {
	if parm == nil {
		llog.Error("CreateRawLeaveMsgSendTx", "parm", parm)
		return nil, types.ErrInvalidParam
	}
	v := &SendMsgParam{
		Msg:  parm.Msg,
		From: parm.From,
		To:   parm.To,
	}
	send := &LeaveMsgAction{
		Ty:    LeavemsgSendAction,
		Value: &LeaveMsgAction_Send{v},
	}

	tx := &types.Transaction{
		Execer:  []byte(getRealExecName(types.GetParaName())),
		Payload: types.Encode(send),
		//Fee:     parm.Fee,
		To: address.ExecAddress(getRealExecName(types.GetParaName())),
	}
	name := getRealExecName(types.GetParaName())
	tx, err := types.FormatTx(name, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
