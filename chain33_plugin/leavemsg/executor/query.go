package executor

import (
	"fmt"
	ltypes "github.com/33cn/plugin/plugin/dapp/leavemsg/types"
	"github.com/33cn/chain33/types"
)

func (l *LeaveMsg) Query_GetMsg(in *ltypes.QueryMsgParam) (types.Message, error) {
	var sendLog ltypes.SendLog
	localKey := []byte(fmt.Sprintf(KeyPrefixSendLocal, in.Addr))
	fmt.Println("In query , the key is ", localKey)
	value, err := l.GetLocalDB().Get(localKey)
	if err != nil {
		return nil, err
	}
	types.Decode(value, &sendLog)
	res := ltypes.QueryMsgResult{Msg: sendLog.Msg, From: sendLog.From, To: sendLog.To}
	return &res, nil
}
