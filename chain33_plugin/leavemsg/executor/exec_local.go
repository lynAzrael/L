package executor

import (
	"fmt"
	ltypes "github.com/33cn/plugin/plugin/dapp/leavemsg/types"
	"github.com/33cn/chain33/types"
)

// 交易执行成功，将本消息对应的数值加1
func (l *LeaveMsg) ExecLocal_Send(send *ltypes.SendMsgParam, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	fmt.Println("ExecLocal send...") // 这里简化处理，不做基本的零值及错误检查了
	for _, log := range receipt.Logs {
		if log.Ty == ltypes.TyLogSend {
			var sendLog ltypes.SendLog
			types.Decode(log.Log, &sendLog)
			localKey := []byte(fmt.Sprintf(KeyPrefixSendLocal, sendLog.To))
			fmt.Println("In exec local, the key is ", localKey)
			oldValue, err := l.GetLocalDB().Get(localKey)
			if err != nil && err != types.ErrNotFound {
				return nil, err
			}
			if err == nil {
				types.Decode(oldValue, &sendLog)
			}
			kv := []*types.KeyValue{{Key: localKey, Value: types.Encode(&sendLog)}}
			return &types.LocalDBSet{KV: kv}, nil
		}
	}
	return nil, nil
}
