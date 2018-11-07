# Client模块
## 1. 模块介绍
Cli模块通过RPC接口，直接调用chain33内部的接口实现系统服务。RPC是通过提供一系列协议方式，对外部应用提供各种系统服务的，在chain33中主要采用了用protobuf定义协议的grpc服务和用json定义协议的jsonrpc服务，分别为不同的前端应用提供相同的系统服务。可以简单的将Cli模块理解成一个前端应用。

## 2. 逻辑架构及上下文
### 2.1 模块关系图
* Cli模块与Chain33的交互
![](https://i.imgur.com/SxrfHMB.jpg)

1、根据输入的指令的不同，调用相的Rpc接口

2、Rpc模块在接收到rpc请求之后，会通过设置topic将消息发送到指定的模块

3、各个模块处理完毕之后，将结果返回给rpc模块

4、最终rpc模块根据响应中的信息，构造成cli需要的结构并返回。

### 2.2 处理逻辑
#### 2.2.1 指令的创建
chain33中使用的cobra进行指令集的创建，rootCmd是所有指令集统一的入口。

	var rootCmd = &cobra.Command{
		Use:   "chain33-cli",
		Short: "chain33 client tools",
	}

 rootCmd涉及cobra.Command结构体中的两个元素：commands 和flags

	type Command struct {
		...	

		// commands is the list of commands supported by this program.
		commands []*Command
		...

		// flags is full set of flags.
		flags *flag.FlagSet
		...
	}

 * commands：表示要执行的动作或指令，而每一个指令又可以包含子命令。
 * flags： 指令可以执行的动作或者过滤条件

 commands通过AddCommand添加:

	// AddCommand adds one or more commands to this parent command.
	func (c *Command) AddCommand(cmds ...*Command) {
		for i, x := range cmds {
			...
			c.commands = append(c.commands, x)
			...
		}
	}

 flags 通过AddFlag设置：

	// AddFlag will add the flag to the FlagSet
	func (f *FlagSet) AddFlag(flag *Flag) {
	}

 可以通过MarkFlagRequired将flag设置为必填项：

	func (c *Command) MarkFlagRequired(name string) error {
		return MarkFlagRequired(c.Flags(), name)
	}

#### 2.2.2 指令的调用

	type RpcCtx struct {
		Addr   string			// 对端rpc地址
		Method string			// 调用函数
		Params interface{}		// 入参
		Res    interface{}		// 响应
		cb     Callback			// 回调函数
	}

RpcCtx的创建

	func NewRpcCtx(laddr, method string, params, res interface{}) *RpcCtx {
		return &RpcCtx{
			Addr:   laddr,
			Method: method,
			Params: params,
			Res:    res,
		}
	}

RpcCtx的执行

	func (c *RpcCtx) Run() {
		result, err := c.RunResult()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		data, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Println(string(data))
	}


#### 2.2.3 指令的注册
chain33目前支持的rpc接口

	type QueueProtocolAPI interface {
		Version() (*types.Reply, error)
		Close()
		NewMessage(topic string, msgid int64, data interface{}) queue.Message
		Notify(topic string, ty int64, data interface{}) (queue.Message, error)
		// +++++++++++++++ mempool interfaces begin
		// 同步发送交易信息到指定模块，获取应答消息 types.EventTx
		SendTx(param *types.Transaction) (*types.Reply, error)
		// types.EventTxList
		GetTxList(param *types.TxHashList) (*types.ReplyTxList, error)
		// types.EventGetMempool
		GetMempool() (*types.ReplyTxList, error)
		// types.EventGetLastMempool
		GetLastMempool() (*types.ReplyTxList, error)
		// +++++++++++++++ execs interfaces begin
		// types.EventBlockChainQuery
		Query(driver, funcname string, param types.Message) (types.Message, error)
		QueryConsensus(param *types.ChainExecutor) (types.Message, error)
		QueryConsensusFunc(driver string, funcname string, param types.Message) (types.Message, error)
		QueryChain(param *types.ChainExecutor) (types.Message, error)
		ExecWalletFunc(driver string, funcname string, param types.Message) (types.Message, error)
		ExecWallet(param *types.ChainExecutor) (types.Message, error)
		// --------------- execs interfaces end
	
		// +++++++++++++++ p2p interfaces begin
		// types.EventPeerInfo
		PeerInfo() (*types.PeerList, error)
		// types.EventGetNetInfo
		GetNetInfo() (*types.NodeNetInfo, error)
		// --------------- p2p interfaces end
		// +++++++++++++++ wallet interfaces begin
		// types.EventLocalGet
		LocalGet(param *types.LocalDBGet) (*types.LocalReplyValue, error)
		// types.EventLocalList
		LocalList(param *types.LocalDBList) (*types.LocalReplyValue, error)
		// types.EventWalletGetAccountList
		WalletGetAccountList(req *types.ReqAccountList) (*types.WalletAccounts, error)
		// types.EventNewAccount
		NewAccount(param *types.ReqNewAccount) (*types.WalletAccount, error)
		// types.EventWalletTransactionList
		WalletTransactionList(param *types.ReqWalletTransactionList) (*types.WalletTxDetails, error)
		// types.EventWalletImportprivkey
		WalletImportprivkey(param *types.ReqWalletImportPrivkey) (*types.WalletAccount, error)
		// types.EventWalletSendToAddress
		WalletSendToAddress(param *types.ReqWalletSendToAddress) (*types.ReplyHash, error)
		// types.EventWalletSetFee
		WalletSetFee(param *types.ReqWalletSetFee) (*types.Reply, error)
		// types.EventWalletSetLabel
		WalletSetLabel(param *types.ReqWalletSetLabel) (*types.WalletAccount, error)
		// types.EventWalletMergeBalance
		WalletMergeBalance(param *types.ReqWalletMergeBalance) (*types.ReplyHashes, error)
		// types.EventWalletSetPasswd
		WalletSetPasswd(param *types.ReqWalletSetPasswd) (*types.Reply, error)
		// types.EventWalletLock
		WalletLock() (*types.Reply, error)
		// types.EventWalletUnLock
		WalletUnLock(param *types.WalletUnLock) (*types.Reply, error)
		// types.EventGenSeed
		GenSeed(param *types.GenSeedLang) (*types.ReplySeed, error)
		// types.EventSaveSeed
		SaveSeed(param *types.SaveSeedByPw) (*types.Reply, error)
		// types.EventGetSeed
		GetSeed(param *types.GetSeedByPw) (*types.ReplySeed, error)
		// types.EventGetWalletStatus
		GetWalletStatus() (*types.WalletStatus, error)
		// types.EventDumpPrivkey
		DumpPrivkey(param *types.ReqString) (*types.ReplyString, error)
		// types.EventSignRawTx
		SignRawTx(param *types.ReqSignRawTx) (*types.ReplySignRawTx, error)
		GetFatalFailure() (*types.Int32, error)
		// types.EventCreateTransaction 由服务器协助创建一个交易
		WalletCreateTx(param *types.ReqCreateTransaction) (*types.Transaction, error)
		// types.EventGetBlocks
		GetBlocks(param *types.ReqBlocks) (*types.BlockDetails, error)
		// types.EventQueryTx
		QueryTx(param *types.ReqHash) (*types.TransactionDetail, error)
		// types.EventGetTransactionByAddr
		GetTransactionByAddr(param *types.ReqAddr) (*types.ReplyTxInfos, error)
		// types.EventGetTransactionByHash
		GetTransactionByHash(param *types.ReqHashes) (*types.TransactionDetails, error)
		// types.EventGetHeaders
		GetHeaders(param *types.ReqBlocks) (*types.Headers, error)
		// types.EventGetBlockOverview
		GetBlockOverview(param *types.ReqHash) (*types.BlockOverview, error)
		// types.EventGetAddrOverview
		GetAddrOverview(param *types.ReqAddr) (*types.AddrOverview, error)
		// types.EventGetBlockHash
		GetBlockHash(param *types.ReqInt) (*types.ReplyHash, error)
		// types.EventIsSync
		IsSync() (*types.Reply, error)
		// types.EventIsNtpClockSync
		IsNtpClockSync() (*types.Reply, error)
		// types.EventGetLastHeader
		GetLastHeader() (*types.Header, error)
	
		//types.EventGetLastBlockSequence:
		GetLastBlockSequence() (*types.Int64, error)
		//types.EventGetBlockSequences:
		GetBlockSequences(param *types.ReqBlocks) (*types.BlockSequences, error)
		//types.EventGetBlockByHashes:
		GetBlockByHashes(param *types.ReqHashes) (*types.BlockDetails, error)
	
		// --------------- blockchain interfaces end
	
		// +++++++++++++++ store interfaces begin
		StoreGet(*types.StoreGet) (*types.StoreReplyValue, error)
		StoreGetTotalCoins(*types.IterateRangeByStateHash) (*types.ReplyGetTotalCoins, error)
		// --------------- store interfaces end
	
		// +++++++++++++++ other interfaces begin
		// close chain33
		CloseQueue() (*types.Reply, error)
		// --------------- other interfaces end
	}

* rpc服务注册:

JsonRPC服务注册：

	func (server *Server) RegisterName(name string, rcvr interface{}) error {
		return server.register(rcvr, name, true)
	}

gRPC服务注册：

	func (s *Server) RegisterService(sd *ServiceDesc, ss interface{}) {
		ht := reflect.TypeOf(sd.HandlerType).Elem()
		st := reflect.TypeOf(ss)
		if !st.Implements(ht) {
			grpclog.Fatalf("grpc: Server.RegisterService found the handler of type %v that does not satisfy %v", st, ht)
		}
		s.register(sd, ss)
	}


#### 2.2.3 指令的处理
jRpc的处理：

	func (j *JSONRPCServer) Listen() {
		...
		var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			...
				serverCodec := jsonrpc.NewServerCodec(&HTTPConn{in: ioutil.NopCloser(bytes.NewReader(data)), out: w, r: r})
				w.Header().Set("Content-type", "application/json")
				if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
					w.Header().Set("Content-Encoding", "gzip")
				}
				w.WriteHeader(200)
				err = j.s.ServeRequest(serverCodec)
				if err != nil {
					log.Debug("Error while serving JSON request: %v", err)
					return
				}
			}
		})
	
		handler = co.Handler(handler)
		http.Serve(listener, handler)
	}


## 3. 指令介绍

### 3.1 account
Account management

	Usage:
	  chain33-cli account [command]
	
	Available Commands:
	  balance     Get balance of a account address
	  create      Create a new account with label
	  dump_key    Dump private key for account address
	  import_key  Import private key with label
	  list        Get account list
	  set_label   Set label for account address


#### 3.1.1 account balance 地址余额查询
cli account balance -a "查询地址" -e "执行器地址"

	[lyn@localhost build]$ ./chain33-cli account balance -a 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt
	{
	    "addr": "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt",
	    "execAccount": [
	        {
	            "execer": "coins",
	            "account": {
	                "balance": "100000000.0000",
	                "frozen": "0.0000"
	            }
	        }
	    ]
	}

#### 3.1.2 account create 新建钱包地址
cli account create -l "自定义地址标签"
	
	[lyn@localhost build]$ ./chain33-cli account create -l test
	{
	    "acc": {
	        "balance": "0.0000",
	        "frozen": "0.0000",
	        "addr": "1RackwdGHK5CzdP8oytmRvdb5EGQP3YUX"
	    },
	    "label": "test"
	}

#### 3.1.3 account dump_key 钱包地址私钥导出
cli account dump_key -a "需要导出的账户地址"

	[lyn@localhost build]$ ./chain33-cli account dump_key -a 1RackwdGHK5CzdP8oytmRvdb5EGQP3YUX
	{
	    "replystr": "0x1a8ba8d001fe0a11b02622297ab599f7a1c1116e272ad759d602c7ba708c55d4"
	}
#### 3.1.4 account import_key 钱包地址私钥导入
cli account import_key -k "外部私钥" -l "地址标签"

	[lyn@localhost build]$ ./chain33-cli account import_key -k "0xa830cd3b4b4b236153c9b67bc161076f0e43eb002fec71e283c0ee0d1f644623" -l test
	{
	    "acc": {
	        "balance": "0.0000",
	        "frozen": "0.0000",
	        "addr": "183BMp5Qcjx52e5yGERGs97DPCioChW7gj"
	    },
	    "label": "test"
	}
#### 3.1.5 account list 获取账户列表
cli account list

	[lyn@localhost build]$ ./chain33-cli account list
	{
	    "wallets": [
	        {
	            "acc": {
	                "balance": "0.0000",
	                "frozen": "0.0000",
	                "addr": "16mMKG3h8yGJxUji6pUQGQWDd29jHA6e2Q"
	            },
	            "label": "node award"
	        },
	        {
	            "acc": {
	                "balance": "0.0000",
	                "frozen": "0.0000",
	                "addr": "183BMp5Qcjx52e5yGERGs97DPCioChW7gj"
	            },
	            "label": "test"
	        }
	    ]
	}
#### 3.1.6 account set_label 设置账户地址标签名
cli account set_label -a "账户地址" -l "地址表签名"

	[lyn@localhost build]$ ./chain33-cli account set_label -a 183BMp5Qcjx52e5yGERGs97DPCioChW7gj -l test1
	{
	    "acc": {
	        "balance": "0.0000",
	        "frozen": "0.0000",
	        "addr": "183BMp5Qcjx52e5yGERGs97DPCioChW7gj"
	    },
	    "label": "test1"
	}


### 3.2 block
Get block header or body info

	Usage:
	  chain33-cli block [command]
	
	Available Commands:
	  get           Get blocks between [start, end]
	  hash          Get hash of block at height
	  headers       Get block headers between [start, end]
	  last_header   View last block header
	  last_sequence View last block sequence
	  query_hashs   Query block by hashs
	  sequences     Get block sequences between [start, end]
	  view          View block info by block hash
	
#### 3.2.1 block get 获取指定区块高度区间的区块详细信息
cli block get -s "起始查询高度" -e "结束查询高度" -d "是否选择显示详情"(可选)

	Usage:
	  chain33-cli block get [flags]
	
	Flags:
	  -d, --detail string   whether print block detail info (0/f/false for No; 1/t/true for Yes) (default "f")
	  -e, --end int         block end height
	  -h, --help            help for get
	  -s, --start int       block start height


例：

	[lyn@localhost build]$ ./chain33-cli block get -s 5765 -e 5765
	{
	    "items": [
	        {
	            "block": {
	                "version": 0,
	                "parenthash": "0x5153191bb51dacb01311f6ef15726fed82f84b756603fbf0d49a473562e45672",
	                "txhash": "0x22849a81d554b4f914d7d65c1080d8cc98e5d1de0fd2be4f0db34b439c6a0173",
	                "statehash": "0x0352bdf0ddb0d51d3e68aa75d402a880c4d877683c2a6393971ca80462af8efd",
	                "height": 5765,
	                "blocktime": 1541135879,
	                "txs": [
	                    {
	                        "execer": "norm",
	                        "payload": {
	                            "rawlog": "0x28010a7c0a1445565641784b61484d797a7175586e494856424f1264774f4a7457625a4846464b5a6e67766a756b6d665179784558526c4f56536b75626659446c736a6951434166695a44424e65506150657841644b42667a5559514548684f6349594b5a576f48447a68726156795a6f73786b44424844734c504458656765"
	                        },
	                        "rawpayload": "0x28010a7c0a1445565641784b61484d797a7175586e494856424f1264774f4a7457625a4846464b5a6e67766a756b6d665179784558526c4f56536b75626659446c736a6951434166695a44424e65506150657841644b42667a5559514548684f6349594b5a576f48447a68726156795a6f73786b44424844734c504458656765",
	                        "signature": {
	                            "ty": 1,
	                            "pubkey": "0x03fe25b1a4261c4b98ad1b81307c7b20b776be6c503b95afb5d73ca9d42daecd7a",
	                            "signature": "0x3044022078d751b8e1d0daeb0aa0eb264ce7b54e1320c7017953a397f38d8cb6c920b93902200267edab907706fa2cbb4bb8606e8677135ee7be5fb820c66ae0004056c2cca8"
	                        },
	                        "fee": "0.0100",
	                        "expire": 0,
	                        "nonce": 6935173039743788574,
	                        "to": "1CnmrBJcpTiY6TphmuAiz7HoYSsGwgYgho",
	                        "from": "1FbZaK5HJRwDaN2sQ5oRoUwi7fHJA1QJT1"
	                    }
	                ]
	            },
	            "receipts": null
	        }
	    ]
	}


#### 3.2.2 block hash 获取指定区块高度的区块哈希
cli block hash -t {区块高度}

	Usage:
	  chain33-cli block hash [flags]
	
	Flags:
	  -t, --height int   block height
	  -h, --help         help for hash


例:

	[lyn@localhost build]$ ./chain33-cli block hash -t 10
    {
        "hash": "0x09055102ecb36033adde0fc9c0d523500c1d81693e7ad12181cfc2497b407da9"
    }

#### 3.2.3 block headers 获取指定区块高度区间内的区块头信息
cli block headers -s 起始查询高度 -e 结束查询高度

	Usage:
	  chain33-cli block headers [flags]
	
	Flags:
	  -d, --detail string   whether print header detail info (0/f/false for No; 1/t/true for Yes) (default "f")
	  -e, --end int         block end height
	  -h, --help            help for headers
	  -s, --start int       block start height


例：

	[lyn@localhost build]$ ./chain33-cli block headers -s 5765 -e 5765
	{
	    "items": [
	        {
	            "version": 0,
	            "parentHash": "0x5153191bb51dacb01311f6ef15726fed82f84b756603fbf0d49a473562e45672",
	            "txHash": "0x22849a81d554b4f914d7d65c1080d8cc98e5d1de0fd2be4f0db34b439c6a0173",
	            "stateHash": "0x0352bdf0ddb0d51d3e68aa75d402a880c4d877683c2a6393971ca80462af8efd",
	            "height": 5765,
	            "blockTime": 1541135879,
	            "txCount": 1,
	            "hash": "0xeab471cf4957253ac991bc56744ef7ae9a6b249236e97db2c5c4d998d3e787f2",
	            "difficulty": 0
	        }
	    ]
	}

※ cli命令中有-d的可选参数项，但无实际作用。可省略。

#### 3.2.4 block last_header 获取本钱包当前已同步的最新区块的区块头信息
cli block last_header

	[lyn@localhost build]$ ./chain33-cli block last_header
	{
	    "version": 0,
	    "parentHash": "0xb6bedeb8b7bcd52348f162bed80bf26420df3e886a13ff9aaf7eb2b733cf392d",
	    "txHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	    "stateHash": "0x0352bdf0ddb0d51d3e68aa75d402a880c4d877683c2a6393971ca80462af8efd",
	    "height": 5772,
	    "blockTime": 1541136720,
	    "txCount": 0,
	    "hash": "0x08bcf1957beb722ca874d46f7f7e35b4dea032ac4c0bd519c2eee214e1661284",
	    "difficulty": 0
	}

#### 3.2.5 block last_sequence 获取本钱包当前已同步的最新顺序编号
cli block last_sequence

	[lyn@localhost build]$ ./chain33-cli block last_sequence
	33

※ sequence是否开启的配置项在配置文件中，默认应为关闭，返回为0。

#### 3.2.6 block query_hashs 根据区块哈希获取区块详情
cli block query_hashs "哈希1" "哈希2"

	[lyn@localhost build]$ ./chain33-cli block query_hashs -s 0x08bcf1957beb722ca874d46f7f7e35b4dea032ac4c0bd519c2eee214e1661284
	{
	    "items": [
	        {
	            "block": {
	                "parentHash": "tr7euLe81SNI8WK+2AvyZCDfPohqE/+ar36ytzPPOS0=",
	                "txHash": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
	                "stateHash": "A1K98N2w1R0+aKp11AKogMTYd2g8KmOTlxyoBGKvjv0=",
	                "height": 5772,
	                "blockTime": 1541136720
	            },
	            "prevStatusHash": "A1K98N2w1R0+aKp11AKogMTYd2g8KmOTlxyoBGKvjv0="
	        }
	    ]
	}


#### 3.2.7 block sequences 根据顺序编号查询对应的区块哈希（仅用于平行链相关）
cli block sequences -s "起始查询高度" -e "结束查询高度"

	[lyn@localhost build]$ ./chain33-cli block sequences -s 0 -e 1
	{
	    "blkseqInfos": [
	        {
	            "hash": "0x67c58d6ba9175313f0468ae4e0ddec946549af7748037c2fdd5d54298afd20b6",
	            "type": 1
	        }
	    ]
	}

#### 3.2.8 block view 根据区块哈希查询指定区块的区块头信息
cli block view -s 区块哈希

	[lyn@localhost build]$ ./chain33-cli block view -s 0xeab471cf4957253ac991bc56744ef7ae9a6b249236e97db2c5c4d998d3e787f2
	{
	    "head": {
	        "version": 0,
	        "parentHash": "0x5153191bb51dacb01311f6ef15726fed82f84b756603fbf0d49a473562e45672",
	        "txHash": "0x22849a81d554b4f914d7d65c1080d8cc98e5d1de0fd2be4f0db34b439c6a0173",
	        "stateHash": "0x0352bdf0ddb0d51d3e68aa75d402a880c4d877683c2a6393971ca80462af8efd",
	        "height": 5765,
	        "blockTime": 1541135879,
	        "txCount": 1,
	        "hash": "0xeab471cf4957253ac991bc56744ef7ae9a6b249236e97db2c5c4d998d3e787f2",
	        "difficulty": 0
	    },
	    "txCount": 1,
	    "txHashes": [
	        "0x22849a81d554b4f914d7d65c1080d8cc98e5d1de0fd2be4f0db34b439c6a0173"
	    ]
	}

### 3.3 bty
Construct BTY transactions

	Usage:
	  chain33-cli bty [command]
	
	Available Commands:
	  priv2priv   Create a privacy to privacy transaction
	  priv2pub    Create a privacy to public transaction
	  pub2priv    Create a public to privacy transaction
	  send_exec   Create a send to executor transaction
	  transfer    Create a transfer transaction
	  txgroup     Create a transaction group
	  withdraw    Create a withdraw transaction

#### 3.3.1 priv2priv
cli bty priv2priv
#### 3.3.2 priv2pub
#### 3.3.3 pub2priv
#### 3.3.4 send_exec 构造指定地址向指定执行器转账的交易
cli bty send_exec -a "发送额度" -e "目的执行器" -n "交易备注信息"  

	[lyn@localhost build]$ ./chain33-cli bty send_exec -a 1000 -e coins
	0a05636f696e731212180a2a0e1080d0dbc3f4022205636f696e7320a08d0630e892958d98bbb6fa7b3a22314761485970576d71414a7371527772706f4e6342385676674b7453776a63487174

※ 返回的是未签名的交易信息，后续参考参考wallet中的sign以及send操作

#### 3.3.5 transfer
cli bty transfer -a "转账额度" -t "接收方账户地址" -n "交易备注信息"

	[lyn@localhost build]$ ./chain33-cli bty transfer -a 5210314 -t "1CRkrCJHHqgQFm59AiHe35T5vJ1o5QpPW3"
	0a05636f696e73123018010a2c1080948e9c81bc7622223143526b72434a4848716751466d35394169486533355435764a316f35517050573320a08d06309de3af99a6c087bc7c3a223143526b72434a4848716751466d35394169486533355435764a316f355170505733

※ 返回的是未签名的交易信息，后续参考参考wallet中的sign以及send操作

#### 3.3.6 txgroup
cli bty txgroup

#### 3.3.7 withdraw
cli bty withdraw -a "交易额度" -e "发送交易的执行器名称" -n "交易备注信息"

	[lyn@localhost build]$ ./chain33-cli bty withdraw -a 1000 -e coins
	0a05636f696e7312121803220e1080d0dbc3f4022205636f696e7320a08d0630f3eaba82b5b098c56e3a22314761485970576d71414a7371527772706f4e6342385676674b7453776a63487174

### 3.4 exec 
Executor operation

	Usage:
	  chain33-cli exec [command]
	
	Available Commands:
	  addr        Get address of executor
	  userdata    Write data to user defined executor

#### 3.4.1 addr 获取执行器的地址
cli exec addr -e "执行器名称"
	
	[lyn@localhost build]$ ./chain33-cli exec addr  -e coins
	1GaHYpWmqAJsqRwrpoNcB8VvgKtSwjcHqt

#### 3.4.2 userdata 写数据到用户自定义执行器

### 3.5 mempool
Mempool management

	Usage:
	  chain33-cli mempool [command]
	
	Available Commands:
	  last_txs    Get latest mempool txs
	  list        List mempool txs

#### 3.5.1 last_txs 获取内存池中最近十笔交易
cli mempool last_txs 

	[azrael@localhost build]$ ./chain33-cli mempool last_txs
	{
	    "txs": [
	        {
	            "execer": "norm",
	            "payload": {
	                "nput": {
	                    "key": "vrgMFUSFeEWNdJBUEWTJ",
	                    "value": "0x6c5a706d44566b756251766e69757547494375645341416c564b5967426c504c797158414169456e7577754d6d4b477363726c6369714b43707978414865426579445963484b774274664f65706249464970635351524c744a42544d66554351745a7965"
	                },
	                "ty": 1
	            },
	            "rawpayload": "0x28010a7c0a147672674d465553466545574e644a42554557544a12646c5a706d44566b756251766e69757547494375645341416c564b5967426c504c797158414169456e7577754d6d4b477363726c6369714b43707978414865426579445963484b774274664f65706249464970635351524c744a42544d66554351745a7965",
	            "signature": {
	                "ty": 1,
	                "pubkey": "0x027bdcc95bc051df0e047cd76f8488b40fd513e429c1b284ddd990f531e4cf42af",
	                "signature": "0x304402202c671aac350f56bd07a3563e3c3bebfbc7ffae2d10ce7c7ea4e14953c7f6770b02203b4fe2185481ba17d089d309150e01e68682c0a5b3d46d3532ad93c1304961ab"
	            },
	            "fee": "0.0100",
	            "expire": 0,
	            "nonce": 8872632247761828259,
	            "to": "1CnmrBJcpTiY6TphmuAiz7HoYSsGwgYgho",
	            "from": "17DjBft6j9VBxJddRxe9eSCiN1Y2NiDe7L"
	        }
	    ]
	}



#### 3.5.2 list 获取内存池中的交易列表
cli mempool list 

	[azrael@localhost build]$ ./chain33-cli mempool list
	{
	    "txs": [
	        {
	            "execer": "norm",
	            "payload": {
	                "nput": {
	                    "key": "vrgMFUSFeEWNdJBUEWTJ",
	                    "value": "0x6c5a706d44566b756251766e69757547494375645341416c564b5967426c504c797158414169456e7577754d6d4b477363726c6369714b43707978414865426579445963484b774274664f65706249464970635351524c744a42544d66554351745a7965"
	                },
	                "ty": 1
	            },
	            "rawpayload": "0x28010a7c0a147672674d465553466545574e644a42554557544a12646c5a706d44566b756251766e69757547494375645341416c564b5967426c504c797158414169456e7577754d6d4b477363726c6369714b43707978414865426579445963484b774274664f65706249464970635351524c744a42544d66554351745a7965",
	            "signature": {
	                "ty": 1,
	                "pubkey": "0x027bdcc95bc051df0e047cd76f8488b40fd513e429c1b284ddd990f531e4cf42af",
	                "signature": "0x304402202c671aac350f56bd07a3563e3c3bebfbc7ffae2d10ce7c7ea4e14953c7f6770b02203b4fe2185481ba17d089d309150e01e68682c0a5b3d46d3532ad93c1304961ab"
	            },
	            "fee": "0.0100",
	            "expire": 0,
	            "nonce": 8872632247761828259,
	            "to": "1CnmrBJcpTiY6TphmuAiz7HoYSsGwgYgho",
	            "from": "17DjBft6j9VBxJddRxe9eSCiN1Y2NiDe7L"
	        }
	    ]
	}

### 3.6 net
Net operation

	Usage:
	  chain33-cli net [command]
	
	Available Commands:
	  fault         Get system fault
	  info          Get net information
	  is_clock_sync Get ntp clock synchronization status
	  is_sync       Get blockchain synchronization status
	  peer_info     Get remote peer nodes
	  time          Get time status

#### 3.6.1 fault 查询本节点出现重大故障的次数
cli net fault

	[lyn@localhost build]$ ./chain33-cli net fault
	0

#### 3.6.2 info 查询本节点的网络信息
cli net info 

	[lyn@localhost build]$ ./chain33-cli net info
	{
	    "externalAddr": "192.168.0.147:13802",
	    "localAddr": "192.168.0.147:13802",
	    "service": true,
	    "outbounds": 0,
	    "inbounds": 0
	}

#### 3.6.3 is_clock_sync 查询本节点时间是否同步
cli net is\_clock_syn

	[lyn@localhost build]$ ./chain33-cli net is_clock_sync
	true

#### 3.6.4 is_sync 查询本节点时间是否同步
cli net is_sync

	[lyn@localhost build]$ ./chain33-cli net  is_sync
	true

#### 3.6.5 peer_info 查询本节点时间是否同步
cli net peer_info

	[lyn@localhost build]$ ./chain33-cli net peer_info
	{
	    "peers": [
	        {
	            "addr": "192.168.0.147",
	            "port": 13802,
	            "name": "02e466e00b8db4e67de85d7c667dabeda92faea9fd06f72c29b2c851eb106fefa4",
	            "mempoolSize": 0,
	            "self": true,
	            "header": {
	                "version": 0,
	                "parentHash": "0x5be67d5026e2bc4aede0feda9ce7cea214b7dc89eff481f1f2c4d33e6a7a3c0a",
	                "txHash": "0x0a382a097d9b7e2d0971832e1b7d49f41059cd6340a5c241f7a4c91bbeaf896f",
	                "stateHash": "0x98f10651192e2df48a0b48200c9024896e1351a342403d27227ee0f37025f15b",
	                "height": 5,
	                "blockTime": 1541393636,
	                "txCount": 100,
	                "hash": "0xab67feeb56baf68eef300ef1f01a2fbe026623e3846729277d74de008e149e58",
	                "difficulty": 0
	            }
	        }
	    ]
	}




#### 3.6.6 time 获取本地时间与ntp时间
cli net time
	
	[lyn@localhost build]$ ./chain33-cli net time
	{
	    "ntpTime": "2018-11-05 14:20:08",
	    "localTime": "2018-11-05 12:55:08",
	    "diff": -5100
	}

### 3.7 seed
Seed management

	Usage:
	  chain33-cli seed [command]
	
	Available Commands:
	  generate    Generate seed
	  get         Get seed by password
	  save        Save seed and encrypt with passwd

#### 3.7.1 generate 生成种子
cli seed generate -l "seed语言种类 0:'English' 1:'简体中文'"

	[lyns@localhost build]$ ./chain33-cli seed generate -l 0
	{
	    "seed": "melt inflict dose foam tuna whip fruit boil scrub rude puzzle length ask cruise embody"
	}

#### 3.7.2 get 获取本钱包的种子
cli seed get -p "获取seed的密码"

	[lyn@localhost build]$ ./chain33-cli seed get -p fzm
	{
	    "seed": "melt inflict dose foam tuna whip fruit boil scrub rude puzzle length ask cruise embody"
	}

#### 3.7.3 save 将种子保存为本钱包的种子
cli seed save -s "以空格分隔的seed(15个字符或单词)" -p "加密seed使用的密码"

	[lyn@localhost build]$ ./chain33-cli seed save -s "melt inflict dose foam tuna whip fruit boil scrub rude puzzle length ask cruise embody" -p fzm
	{
	    "isOK": true,
	    "msg": ""
	}

### 3.8 send
cli send bty transfer -a "转账额度" -n "备注信息" -t "接收方地址" -k "私钥/发送方地址"

	[lyn@localhost build]$ ./chain33-cli send bty transfer -a "1000" -n "transfer tx" -t "1LBKc8mA7s57sVoij6AhL7CG6pS3TFkqBu" -k "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
	0x037513739705f899d92add16b3bbcc077076d81e20074b61ca7b391875514480

### 3.9 stat
Coin statistic

	Usage:
	  chain33-cli stat [command]
	
	Available Commands:
	  miner            Get miner statistic
	  ticket_info      Get ticket info by ticket_id
	  ticket_info_list Get ticket info list by ticket_id
	  ticket_stat      Get ticket statistics by addr
	  total_coins      Get total amount of a token (default: bty of current height)

#### 3.9.1 
cli stat miner 

### 3.10 tx
Transaction management

Usage:
  chain33-cli tx [command]

Available Commands:
  addr_overview View transactions of address
  decode        Decode a hex format transaction
  get_hex       Get transaction hex by hash
  query         Query transaction by hash
  query_addr    Query transaction by account address
  query_hash    Get transactions by hashes

#### 3.10.1 addr_overview 查询指定地址在coins合约下的相关交易
cli tx addr_overview -a "账户地址"
	
	[azrael@localhost build]$ ./chain33-cli tx addr_overview  -a 1EcE1nzwRzhrUUVy9k2xDMH5kw9bymcJPc
	{
	    "receiver": "1000.0000",
	    "balance": "100.9960",
	    "txCount": 5
	}
	
#### 3.10.2 decode 对交易字符串进行解析
cli tx decode -d "需要解析的交易信息"

	[azrael@localhost build]$ ./chain33-cli tx decode -d 0a05636f696e7312121803220e1080d0dbc3f4022205636f696e7320a08d0630eaceade1e389c1a5293a22314761485970576d71414a7371527772706f4e6342385676674b7453776a63487174
	{
	    "execer": "coins",
	    "payload": {
	        "withdraw": {
	            "cointoken": "",
	            "amount": "100000000000",
	            "note": "",
	            "execName": "coins",
	            "to": ""
	        },
	        "ty": 3
	    },
	    "rawpayload": "0x1803220e1080d0dbc3f4022205636f696e73",
	    "signature": {
	        "ty": 0,
	        "pubkey": "",
	        "signature": ""
	    },
	    "fee": "0.0010",
	    "expire": 0,
	    "nonce": 2975476712871782250,
	    "to": "1GaHYpWmqAJsqRwrpoNcB8VvgKtSwjcHqt",
	    "from": "1HT7xU2Ngenf7D4yocz2SAcnNLW7rK8d4E"
	}

#### 3.10.3 get_hex 根据交易获取十六进制字符串
cli tx get_hex -s "交易哈希"

	[azrael@localhost build]$ ./chain33-cli tx get_hex -s 0x57fbde79b74a0c3286ca5947c05a5b95efd14f1411669bedabae1c68390237a3
	0a05636f696e73122f18010a2b1080b081daaf14222231456345316e7a77527a687255555679396b3278444d48356b773962796d634a50631a6e0801122102504fa1c28caaf1d5a20fefb87c50a49724ff401043420cb3ba271997eb5a43871a47304502210081d32f66581903960f5fedc79ad742f43a4198209b766e1d5e6668288c3a9bc9022062bb2aa67c5db22aff68cc179fd7af7aad92b7f8cfd634857ff9b77c9c97d86320a08d062887e789df05309bd0b0d793f4cfca6f3a2231456345316e7a77527a687255555679396b3278444d48356b773962796d634a5063

#### 3.10.4 query 根据交易哈希获取交易详情
cli tx query -s "交易哈希" 

	[azrael@localhost build]$ ./chain33-cli  tx query -s 0x57fbde79b74a0c3286ca5947c05a5b95efd14f1411669bedabae1c68390237a3
	{
	    "tx": {
	        "execer": "coins",
	        "payload": {
	            "transfer": {
	                "cointoken": "",
	                "amount": "700000000000",
	                "note": "",
	                "to": "1EcE1nzwRzhrUUVy9k2xDMH5kw9bymcJPc"
	            },
	            "ty": 1
	        },
	        "rawpayload": "0x18010a2b1080b081daaf14222231456345316e7a77527a687255555679396b3278444d48356b773962796d634a5063",
	        "signature": {
	            "ty": 1,
	            "pubkey": "0x02504fa1c28caaf1d5a20fefb87c50a49724ff401043420cb3ba271997eb5a4387",
	            "signature": "0x304502210081d32f66581903960f5fedc79ad742f43a4198209b766e1d5e6668288c3a9bc9022062bb2aa67c5db22aff68cc179fd7af7aad92b7f8cfd634857ff9b77c9c97d863"
	        },
	        "fee": "0.0010",
	        "expire": 1541567367,
	        "nonce": 8040402671450728475,
	        "to": "1EcE1nzwRzhrUUVy9k2xDMH5kw9bymcJPc",
	        "from": "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
	    },
	    "receipt": {
	        "ty": 2,
	        "tyName": "ExecOk",
	        "logs": [
	            {
	                "ty": 2,
	                "tyName": "LogFee",
	                "log": {
	                    "prev": {
	                        "currency": 0,
	                        "balance": "9999699999700000",
	                        "frozen": "0",
	                        "addr": "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
	                    },
	                    "current": {
	                        "currency": 0,
	                        "balance": "9999699999600000",
	                        "frozen": "0",
	                        "addr": "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
	                    }
	                },
	                "rawLog": "0x0a2d10a0e8deb2c9d5e111222231344b454b6259744b4b516d34774d7468534b394a344c61346e41696964476f7a74122d1080dbd8b2c9d5e111222231344b454b6259744b4b516d34774d7468534b394a344c61346e41696964476f7a74"
	            },
	            {
	                "ty": 3,
	                "tyName": "LogTransfer",
	                "log": {
	                    "prev": {
	                        "currency": 0,
	                        "balance": "9999699999600000",
	                        "frozen": "0",
	                        "addr": "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
	                    },
	                    "current": {
	                        "currency": 0,
	                        "balance": "9998999999600000",
	                        "frozen": "0",
	                        "addr": "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
	                    }
	                },
	                "rawLog": "0x0a2d1080dbd8b2c9d5e111222231344b454b6259744b4b516d34774d7468534b394a344c61346e41696964476f7a74122d1080abd7d899c1e111222231344b454b6259744b4b516d34774d7468534b394a344c61346e41696964476f7a74"
	            },
	            {
	                "ty": 3,
	                "tyName": "LogTransfer",
	                "log": {
	                    "prev": {
	                        "currency": 0,
	                        "balance": "300000000000",
	                        "frozen": "0",
	                        "addr": "1EcE1nzwRzhrUUVy9k2xDMH5kw9bymcJPc"
	                    },
	                    "current": {
	                        "currency": 0,
	                        "balance": "1000000000000",
	                        "frozen": "0",
	                        "addr": "1EcE1nzwRzhrUUVy9k2xDMH5kw9bymcJPc"
	                    }
	                },
	                "rawLog": "0x0a2b1080f092cbdd08222231456345316e7a77527a687255555679396b3278444d48356b773962796d634a5063122b1080a094a58d1d222231456345316e7a77527a687255555679396b3278444d48356b773962796d634a5063"
	            }
	        ]
	    },
	    "height": 38,
	    "index": 0,
	    "blocktime": 1541567247,
	    "amount": "7000.0000",
	    "fromaddr": "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt",
	    "actionname": "transfer"
	}


#### 3.10.5 query_addr 根据账户地址获取交易列表
cli tx query_addr -a "账户地址" -t "区块高度" -c "最大返回条数" -d "查询方式" -f "交易类型" 

	[azrael@localhost build]$  ./chain33-cli  tx query_addr -a 1EcE1nzwRzhrUUVy9k2xDMH5kw9bymcJPc -t 36 -c 1 -d 0 -f 0
	{
	    "txInfos": [
	        {
	            "hash": "0x79811df8d4ea2fabc1a14e57e48014d71a95dcca34b336e44318920445886862",
	            "height": 35,
	            "index": 0,
	            "assets": [
	                {
	                    "exec": "coins",
	                    "symbol": "BTY"
	                }
	            ]
	        }
	    ]
	}

#### 3.10.6 query_hash 根据交易哈希（组）获取交易详情
cli tx query_hash -s "交易哈希" "交易哈希" "交易哈希" ...

### 3.11 version
Version info

	[lyn@localhost build]$ ./chain33-cli version
	5.3.0-04cad269


### 3.12 wallet
Wallet management

	Usage:
	  chain33-cli wallet [command]
	
	Available Commands:
	  auto_mine   Set auto mine on/off
	  list_txs    List transactions in wallet
	  lock        Lock wallet
	  merge       Merge accounts' balance into address
	  nobalance   Create nobalance transaction
	  send        Send a transaction
	  set_fee     Set transaction fee
	  set_pwd     Set password
	  sign        Sign transaction
	  status      Get wallet status
	  unlock      Unlock wallet

#### 3.23.1 auto_mine 获取本钱包的交易列表
cli wallet auto_mine -f "是否开启自动挖矿(0:off 1:on)"

	[lyn@localhost build]$ ./chain33-cli  wallet auto_mine -f 0
	{
	    "isOK": true,
	    "msg": ""
	}
#### <font color=red>3.23.2 list_txs 获取本钱包的交易列表</font>

cli wallet list_txs -c "交易数量" -d "查询方式" -f "查询的起始交易地址"



#### 3.23.3 lock 锁定钱包
cli wallet lock

	[lyn@localhost build]$ ./chain33-cli wallet lock
	{
	    "isOK": true,
	    "msg": ""
	}

#### 3.23.4 merge 钱包中地址余额合并
cli wallet merge -t "钱包中某个目标账户地址"

	[lyn@localhost build]$ ./chain33-cli wallet merge -t 1669FvjdPUAqaNLz8yBpXGnfDcqFL7zozG
	{
	    "hashes": [
	        "0x5664ed144c6d202d6085eb1348e5d6e4c2f3e3d9745806af1e5189b1798a9800"
	    ]
	}

#### 3.23.5 nobalance  对不涉及转账金额的交易构造不需要手续费的交易组
cli wallet nobalance -d "未签名的交易数据" -k "签名私钥"
	
	[azrael@localhost build]$ ./chain33-cli wallet nobalance -d "0a14757365722e702e67756f64756e2e7469636b6574" -k 0x2660c263b11dbdc1c78e8183230ceec1d0204b00f9cfc220c68b9df3aedc116c
	0a046e6f6e6512126e6f2d6665652d7472616e73616374696f6e1a6e08011221036874ba6f252d49c6b90e586ef8fadcaa8b70a7300de3c8ef0cd91dc26eb7d72d1a4730450221009ff944a8118258088f3521a7dadd139d43cf981b76479e2b83f10ee6f9a44aba022054068759504c3112f8fffeea3a0c64edbe2b56831d5525d07fb56c50ed4ce6c520d00f30c6eed8caf287a0e94c3a2231447a5464544c61354a50704c644e4e50325072563161364a4374554c413747735440024ac0020a81020a046e6f6e6512126e6f2d6665652d7472616e73616374696f6e1a6e08011221036874ba6f252d49c6b90e586ef8fadcaa8b70a7300de3c8ef0cd91dc26eb7d72d1a4730450221009ff944a8118258088f3521a7dadd139d43cf981b76479e2b83f10ee6f9a44aba022054068759504c3112f8fffeea3a0c64edbe2b56831d5525d07fb56c50ed4ce6c520d00f30c6eed8caf287a0e94c3a2231447a5464544c61354a50704c644e4e50325072563161364a4374554c413747735440024a2064aa77679686ff056fa1403cc71accc34ade6b668a6971d0cbdd3e3d5a727e895220d8e9b0eaa0a3a6c46cebe56ae8cd662457bf4fa32e866c4ceca36c4ea984b5640a3a0a14757365722e702e67756f64756e2e7469636b657440024a2064aa77679686ff056fa1403cc71accc34ade6b668a6971d0cbdd3e3d5a727e895220d8e9b0eaa0a3a6c46cebe56ae8cd662457bf4fa32e866c4ceca36c4ea984b564

#### 3.23.6 send 对已完成签名的交易进行发送
cli wallet send -d "已签名的交易信息" -t "发送的Token名称(默认为BTY)"

	[lyn@localhost build]$ ./chain33-cli wallet send -d 0a13757365722e702e67756f64756e2e746f6b656e1236380422320a03434e59100a1a05313233313222223136363946766a6450554171614e4c7a3879427058476e66446371464c377a6f7a471a6d0801122102504fa1c28caaf1d5a20fefb87c50a49724ff401043420cb3ba271997eb5a43871a463044022067f735e00946957fde90a7a3ba4418c3e3771977f13b2dc510f6172b0cfb97d902202d2a11e7a69aeee8adf42607d1e32e42c6a43b10ca0702f704cd54e151c045fa20a08d0628cae2ffde053085caa48bc7c594be113a223144527535423766505961776179414b524648315857536e354b66415654386d6848
	0x66a9cd227a8483a20c66032d7c071f93de4bf23b430adbed4ee6aa91ec349f15 

※ 返回结果即区块中记录的交易hash

#### 3.23.7 set_fee 设置交易手续费
cli wallet set_fee -a "交易手续费"

	[lyn@localhost build]$ ./chain33-cli wallet set_fee -a 100
	{
	    "isOK": true,
	    "msg": ""
	}

#### 3.23.8 set_pwd 设置钱包密码
cli wallet set_pwd -o "旧密码" -n "新密码"

	[lyn@localhost build]$ ./chain33-cli wallet set_pwd -o fzm -n fzm123
	{
	    "isOK": true,
	    "msg": ""
	}

 
#### 3.23.9 sign 对已构造的交易进行签名
cli wallet sign -d "未签名的交易信息" -a 签名地址/-k 签名私钥 -e 超时时间（默认为120秒）

	[lyn@localhost build]$ ./chain33-cli wallet sign -d 0a13757365722e702e67756f64756e2e746f6b656e1236380422320a03434e59100a1a05313233313222223136363946766a6450554171614e4c7a3879427058476e66446371464c377a6f7a4720a08d063085caa48bc7c594be113a223144527535423766505961776179414b524648315857536e354b66415654386d6848 -a 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt
	0a13757365722e702e67756f64756e2e746f6b656e1236380422320a03434e59100a1a05313233313222223136363946766a6450554171614e4c7a3879427058476e66446371464c377a6f7a471a6d0801122102504fa1c28caaf1d5a20fefb87c50a49724ff401043420cb3ba271997eb5a43871a463044022067f735e00946957fde90a7a3ba4418c3e3771977f13b2dc510f6172b0cfb97d902202d2a11e7a69aeee8adf42607d1e32e42c6a43b10ca0702f704cd54e151c045fa20a08d0628cae2ffde053085caa48bc7c594be113a223144527535423766505961776179414b524648315857536e354b66415654386d6848
※ 返回结果为加密后的交易信息

#### 3.23.10 status 获取钱包状态
cli wallet status 

	[lyn@localhost build]$ ./chain33-cli wallet status
	{
	    "isWalletLock": false,
	    "isAutoMining": false,
	    "isHasSeed": true,
	    "isTicketLock": true
	}

#### 3.23.11 unlock 获取钱包状态
cli wallet unlock -p "密码" -t "持续时间" -s "解锁范围(默认解锁wallet)"

	[lyn@localhost build]$ ./chain33-cli wallet unlock -p fzm -t 0
	{
	    "isOK": true,
	    "msg": ""
	}

---

#RPC模块
## 1. 模块介绍

## 2. 
##4. 二次开发
本章节将会以一个简单的例子来介绍Cli的二次开发过程，例子实现的功能比较简单，但包含了实现一个cli所需要的各方面的方法。

###4.1 
