package types

type LeaveMsgSendTx struct {
	Msg  string `json:"msg"`
	From string `json:"from"`
	To   string `json:"to"`
}

type LeaveMsgQueryTx struct {
	Addr string `json:"addr"`
}
