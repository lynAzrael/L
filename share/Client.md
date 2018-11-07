# Client模块
## 1. 模块介绍
Cli模块通过RPC接口，直接调用chain33内部的接口实现系统服务。RPC是通过提供一系列协议方式，对外部应用提供各种系统服务的，在chain33中主要采用了用protobuf定义协议的grpc服务和用json定义协议的jsonrpc服务，分别为不同的前端应用提供相同的系统服务。可以简单的将Cli模块理解成一个前端应用。

## 2. 逻辑架构及上下文
### 2.1 模块关系图
* chain33中的位置

### 2.2 处理逻辑
#### 2.2.1 指令的创建
chain33中使用的cobra进行指令集的创建，此处声明的rootCmd是所有指令集统一的入口。

	var rootCmd = &cobra.Command{
		Use:   "chain33-cli",
		Short: "chain33 client tools",
	}

 在编码Chain33使用的指令集时，主要使用cobra.Command结构体中的两个元素：commands 和flags

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

>  commands是通过AddCommand将子指令集添加到rootCmd中的：
>  
>  命令添加到rootCmd中
>
	rootCmd.AddCommand(
		commands.AccountCmd(),
		commands.BlockCmd(),
		commands.BTYCmd(),
		commands.CoinsCmd(),
		...
	)

> 操作指令添加到各个命令中:
>
	cmd.AddCommand(
		DumpKeyCmd(),
		GetAccountListCmd(),
		GetBalanceCmd(),
		ImportKeyCmd(),
		NewAccountCmd(),
		SetLabelCmd(),
	)
}

> flags 是在实现子命令时，通过Flags设置的。比如[block header](#323-block)中的flag, 并用MarkFlagRequired()函数来将flag设置为必填。
>
	func addBlockHeaderFlags(cmd *cobra.Command) {
		cmd.Flags().Int64P("start", "s", 0, "block start height")
		cmd.MarkFlagRequired("start")
		cmd.Flags().Int64P("end", "e", 0, "block end height")
		cmd.MarkFlagRequired("end")
		cmd.Flags().StringP("detail", "d", "f", "whether print header detail info (0/f/false for No; 1/t/true for Yes)")
	}

#### 2.2.2 指令的注册
执行指令实际上是调用chain33内部已经注册好的一些接口函数。例如block last_header命令中最终调用的是chain33内部的GetLastHeader()函数

	func lastHeader(cmd *cobra.Command, args []string) {
		rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
		var res jsonrpc.Header
		ctx := NewRpcCtx(rpcLaddr, "Chain33.GetLastHeader", nil, &res)
		ctx.Run()
	}


1. 这些供外部组件使用的接口函数同意声明在types/proto目录下的rpc.proto文件中

		service chain33 {
		    // chain33 对外提供服务的接口
		    //区块链接口
		    rpc GetBlocks(ReqBlocks) returns (Reply) {}
		    //获取最新的区块头
		    rpc GetLastHeader(ReqNil) returns (Header) {}
		    //交易接口
		    rpc CreateRawTransaction(CreateTx) returns (UnsignTx) {}
		    rpc CreateRawTxGroup(CreateTransactionGroup) returns (UnsignTx) {}
		    //发送签名后交易
		    rpc SendRawTransaction(SignedTx) returns (Reply) {}
		    // 根据哈希查询交易
		    rpc QueryTransaction(ReqHash) returns (TransactionDetail) {}
		    // 发送交易
		    rpc SendTransaction(Transaction) returns (Reply) {}
		
		    //通过地址获取交易信息
		    rpc GetTransactionByAddr(ReqAddr) returns (ReplyTxInfos) {}
		
		    //通过哈希数组获取对应的交易
		    rpc GetTransactionByHashes(ReqHashes) returns (TransactionDetails) {}
		
		    //缓存接口
		    rpc GetMemPool(ReqNil) returns (ReplyTxList) {}
		
		    ...
		}

* client/queueprotocolapi.go文件中进行实现

		type QueueProtocolAPI interface {
			...
			// types.EventGetLastHeader
			GetLastHeader() (*types.Header, error)
			...
		}
		...
	
		func (q *QueueProtocol) GetLastHeader() (*types.Header, error) {
			msg, err := q.query(blockchainKey, types.EventGetLastHeader, &types.ReqNil{})
			if err != nil {
				log.Error("GetLastHeader", "Error", err.Error())
				return nil, err
			}
			if reply, ok := msg.GetData().(*types.Header); ok {
				return reply, nil
			}
			err = types.ErrTypeAsset
			log.Error("GetLastHeader", "Error", err.Error())
			return nil, err
		}

* 启动进程初始化grpc和jsonrpc的时侯，注册到消息队列中。

	JsonRPC注册：

		func NewJSONRPCServer(c queue.Client) *JSONRPCServer {
			j := &JSONRPCServer{}
			j.jrpc.cli.Init(c)
			server := rpc.NewServer()
			j.s = server
			server.RegisterName("Chain33", &j.jrpc)
			return j
		}

	gRPC注册：

		func NewGRpcServer(c queue.Client) *Grpcserver {
			s := &Grpcserver{}
			s.grpc.cli.Init(c)
			var opts []grpc.ServerOption
			//register interceptor
			//var interceptor grpc.UnaryServerInterceptor
			interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
				if err := auth(ctx, info); err != nil {
					return nil, err
				}
				// Continue processing the request
				return handler(ctx, req)
			}
			opts = append(opts, grpc.UnaryInterceptor(interceptor))
			server := grpc.NewServer(opts...)
			s.s = server
			types.RegisterChain33Server(server, &s.grpc)
			return s
		}

* 初始化完成之后，则开始监听客户端的连接

		func (r *RPC) SetQueueClient(c queue.Client) {
			gapi := NewGRpcServer(c)
			japi := NewJSONRPCServer(c)
			r.gapi = gapi
			r.japi = japi
			r.c = c
			//注册系统rpc
			pluginmgr.AddRPC(r)
			go gapi.Listen()
			go japi.Listen()
		}

#### 2.2.3 指令的处理
目前Chain33中命令集在接收响应时均使用的json编码，所以针对jsonRPC看下Server是如何处理命令集的rpc请求的。

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

http.Serve中，每accept一个连接就会创建一个goroutine去处理。

	func (srv *Server) Serve(l net.Listener) error {
		...
		for {
			rw, e := l.Accept()
			...
			c := srv.newConn(rw)
			c.setState(c.rwc, StateNew) // before Serve can return
			go c.serve(ctx)
		}
	}

在创建的协程 中会通过ServeHTTP调用上文中传入的handler。

	func (c *conn) serve(ctx context.Context) {
		...
		for {
			w, err := c.readRequest(ctx)
			...
			serverHandler{c.server}.ServeHTTP(w, w.req)
			...
		}
	}

handler的ServeRequest函数中，从接收到的请求中解码出需要调用的rpc函数进行调用

	func (server *Server) readRequest(codec ServerCodec) (service *service, mtype *methodType, req *Request, argv, replyv reflect.Value, keepReading bool, err error) {
		service, mtype, req, keepReading, err = server.readRequestHeader(codec)
		if err != nil {
			if !keepReading {
				return
			}
			// discard body
			codec.ReadRequestBody(nil)
			return
		}
	
		// Decode the argument value.
		argIsValue := false // if true, need to indirect before calling.
		if mtype.ArgType.Kind() == reflect.Ptr {
			argv = reflect.New(mtype.ArgType.Elem())
		} else {
			argv = reflect.New(mtype.ArgType)
			argIsValue = true
		}
		// argv guaranteed to be a pointer now.
		if err = codec.ReadRequestBody(argv.Interface()); err != nil {
			return
		}
		if argIsValue {
			argv = argv.Elem()
		}
	
		replyv = reflect.New(mtype.ReplyType.Elem())
	
		switch mtype.ReplyType.Elem().Kind() {
		case reflect.Map:
			replyv.Elem().Set(reflect.MakeMap(mtype.ReplyType.Elem()))
		case reflect.Slice:
			replyv.Elem().Set(reflect.MakeSlice(mtype.ReplyType.Elem(), 0, 0))
		}
		return
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

### 3.9 stat

### 3.10 tx

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
