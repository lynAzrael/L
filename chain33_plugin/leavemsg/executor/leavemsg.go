package executor

import (
	driver "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	ltypes "github.com/33cn/plugin/plugin/dapp/leavemsg/types"
)

var (
	// 执行交易生成的数据KEY
	KeyPrefixSend = "mavl-leavemsg-send:%s"
	KeyPrefixFind = "mavl-leavemsg-find:%s"
	// 本地执行生成的数据KEY
	KeyPrefixSendLocal = "LODB-leavemsg-send:%s"
	KeyPrefixFindLocal = "LODB-leavemsg-find:%s"

	driverName = ltypes.LeaveMsgX
)

// 初始化时通过反射获取本执行器的方法列表
func init() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&LeaveMsg{}))
}

func Init(name string, sub []byte) {
	driver.Register(driverName, NewLeaveMsg, types.GetDappFork(driverName, "Enable"))
}

type LeaveMsg struct {
	driver.DriverBase
}

func NewLeaveMsg() driver.Driver {
	t := &LeaveMsg{}
	t.SetChild(t)
	t.SetExecutorType(types.LoadExecutorType(driverName))
	return t
}

func (l *LeaveMsg) GetDriverName() string {
	return driverName
}

func GetName() string {
	return NewLeaveMsg().GetName()
}
