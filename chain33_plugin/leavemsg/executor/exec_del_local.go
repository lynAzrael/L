package executor
import (
	"fmt"
	ltypes "github.com/33cn/plugin/plugin/dapp/leavemsg/types"
	"github.com/33cn/chain33/types"
)
// 交易执行成功，将本消息对应的数值减1
func (l *LeaveMsg) ExecDelLocal_Send(send *ltypes.SendMsgParam, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	// 这里简化处理，不做基本的零值及错误检查了
	var sendLog ltypes.SendLog
	types.Decode(receipt.Logs[0].Log, &sendLog)
	localKey := []byte(fmt.Sprintf(KeyPrefixSendLocal, sendLog.Msg))
	oldValue, err := l.GetLocalDB().Get(localKey)
	if err != nil {
		return nil, err
	}
	types.Decode(oldValue, &sendLog)

	val := types.Encode(&sendLog)

	kv := []*types.KeyValue{{Key:localKey, Value:val}}
	return &types.LocalDBSet{KV:kv}, nil
}
