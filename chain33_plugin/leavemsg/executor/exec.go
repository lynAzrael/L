package executor

import (
	"fmt"
	ltypes "github.com/33cn/plugin/plugin/dapp/leavemsg/types"
)
import "github.com/33cn/chain33/types"

func (l *LeaveMsg) Exec_Send(send *ltypes.SendMsgParam, tx *types.Transaction, index int) (*types.Receipt, error) {
	msg := send.Msg
	from := send.From
	to := send.To
	res := from + " send a message to " + to
	xx := &ltypes.SendLog{Msg: msg, Echo: res, From: from, To: to}
	receiptLog := types.Encode(xx)
	fmt.Println("xx ====> ", xx)
	fmt.Println("after encode ====>", string(receiptLog))
	log := &types.ReceiptLog{Ty: ltypes.TyLogSend, Log: receiptLog}
	var logs []*types.ReceiptLog
	logs = append(logs, log)
	kv := []*types.KeyValue{{Key: []byte(fmt.Sprintf(KeyPrefixSend, msg)), Value: []byte(res)}}
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	fmt.Println("In exec, the bytes is ", receipt.Logs[0].Log)
	return receipt, nil
}
