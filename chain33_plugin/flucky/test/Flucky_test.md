# Flucky测试


## 1 功能测试
### 1.1 投注操作

```bash
[azrael@localhost build]$ curl --data-binary '{"jsonrpc":"2.0", "id": 1, "method":"Chain33.Query","params":[{"execer":"flucky", "funcName":"QueryBetTimes", "payload":{"addr":"14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"}}]} '         -H 'content-type:text/plain;'         http://localhost:8801
{"id":1,"result":{"times":9},"error":null}
[azrael@localhost build]$ curl --data-binary '{"jsonrpc":"2.0", "id": 1, "method":"Chain33.Query","params":[{"execer":"flucky", "funcName":"QueryBetInfo", "payload":{"addr":"14KEKbYtKKQm4wMthSK9J4La4nAiidGozt", "idx": 9}}]} '         -H 'content-type:text/plain;'         http://localhost:8801
{"id":1,"result":{"index":9,"addr":"14KEKbYtKKQm4wMthSK9J4La4nAiidGozt","time":"1546046852","amount":5,"randNum":["2712","9680","5519","7166","3043","8376","8923","3311","7031","3318"],"maxNum":"9680","bonus":0.24503},"error":null}
[azrael@localhost build]$ curl --data-binary '{"jsonrpc":"2.0", "id": 1, "method":"Chain33.CreateTransaction","params":[{"execer":"flucky", "actionName":"Bet", "payload":{"amount":5}}]} '         -H 'content-type:text/plain;'         http://localhost:8801
{"id":1,"result":"0a06666c75636b79120730e4190a02080520a08d0630d4ad859cf6ddf594343a2231474854716e6d6a66336f3963335a6b537136356158636d6867567a6f4b45446f72","error":null}
[azrael@localhost build]$ sign 0a06666c75636b79120730e4190a02080520a08d0630d4ad859cf6ddf594343a2231474854716e6d6a66336f3963335a6b537136356158636d6867567a6f4b45446f72 -a 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt
0a06666c75636b79120730e4190a0208051a6e0801122102504fa1c28caaf1d5a20fefb87c50a49724ff401043420cb3ba271997eb5a43871a473045022100b9d69e67cf2d764742ca60d58f93ebe1d9b8c68d8f33e09839cbdefbcf4fb73302201ba9106681500f7505c15327a535d63e87dd12804b7d4bb406a9b63e7768885f20a08d0628d2a19be10530d4ad859cf6ddf594343a2231474854716e6d6a66336f3963335a6b537136356158636d6867567a6f4b45446f72
[azrael@localhost build]$ send 0a06666c75636b79120730e4190a0208051a6e0801122102504fa1c28caaf1d5a20fefb87c50a49724ff401043420cb3ba271997eb5a43871a473045022100b9d69e67cf2d764742ca60d58f93ebe1d9b8c68d8f33e09839cbdefbcf4fb73302201ba9106681500f7505c15327a535d63e87dd12804b7d4bb406a9b63e7768885f20a08d0628d2a19be10530d4ad859cf6ddf594343a2231474854716e6d6a66336f3963335a6b537136356158636d6867567a6f4b45446f72
0x5139a64934c8feeb2ed4266230579661a0ff2eef82f6ec6a453bf59826b379da
[azrael@localhost build]$ querytx 0x5139a64934c8feeb2ed4266230579661a0ff2eef82f6ec6a453bf59826b379da
{
    "tx": {
        "execer": "flucky",
        "payload": {
            "bet": {
                "amount": 5
            },
            "ty": 3300
        },
        "rawpayload": "0x30e4190a020805",
        "signature": {
            "ty": 1,
            "pubkey": "0x02504fa1c28caaf1d5a20fefb87c50a49724ff401043420cb3ba271997eb5a4387",
            "signature": "0x3045022100b9d69e67cf2d764742ca60d58f93ebe1d9b8c68d8f33e09839cbdefbcf4fb73302201ba9106681500f7505c15327a535d63e87dd12804b7d4bb406a9b63e7768885f"
        },
        "fee": "0.0010",
        "expire": 1546047698,
        "nonce": 3758771687672338132,
        "to": "1GHTqnmjf3o9c3ZkSq65aXcmhgVzoKEDor",
        "from": "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt",
        "hash": "0x5139a64934c8feeb2ed4266230579661a0ff2eef82f6ec6a453bf59826b379da"
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
                        "balance": "9989999998900000",
                        "frozen": "0",
                        "addr": "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
                    },
                    "current": {
                        "currency": 0,
                        "balance": "9989999998800000",
                        "frozen": "0",
                        "addr": "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
                    }
                },
                "rawLog": "0x0a2d10a0aef689a2bbdf11222231344b454b6259744b4b516d34774d7468534b394a344c61346e41696964476f7a74122d1080a1f089a2bbdf11222231344b454b6259744b4b516d34774d7468534b394a344c61346e41696964476f7a74"
            },
            {
                "ty": 6,
                "tyName": "LogExecTransfer",
                "log": {
                    "execAddr": "1GHTqnmjf3o9c3ZkSq65aXcmhgVzoKEDor",
                    "prev": {
                        "currency": 0,
                        "balance": "9996244203000",
                        "frozen": "0",
                        "addr": "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
                    },
                    "current": {
                        "currency": 0,
                        "balance": "9995744203000",
                        "frozen": "0",
                        "addr": "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
                    }
                },
                "rawLog": "0x0a2231474854716e6d6a66336f3963335a6b537136356158636d6867567a6f4b45446f72122c10f8ebd6f4f6a202222231344b454b6259744b4b516d34774d7468534b394a344c61346e41696964476f7a741a2c10f8a1a186f5a202222231344b454b6259744b4b516d34774d7468534b394a344c61346e41696964476f7a74"
            },
            {
                "ty": 6,
                "tyName": "LogExecTransfer",
                "log": {
                    "execAddr": "1GHTqnmjf3o9c3ZkSq65aXcmhgVzoKEDor",
                    "prev": {
                        "currency": 0,
                        "balance": "3755797000",
                        "frozen": "0",
                        "addr": "1GHTqnmjf3o9c3ZkSq65aXcmhgVzoKEDor"
                    },
                    "current": {
                        "currency": 0,
                        "balance": "4255797000",
                        "frozen": "0",
                        "addr": "1GHTqnmjf3o9c3ZkSq65aXcmhgVzoKEDor"
                    }
                },
                "rawLog": "0x0a2231474854716e6d6a66336f3963335a6b537136356158636d6867567a6f4b45446f72122a1088d4f3fe0d222231474854716e6d6a66336f3963335a6b537136356158636d6867567a6f4b45446f721a2a10889ea9ed0f222231474854716e6d6a66336f3963335a6b537136356158636d6867567a6f4b45446f72"
            },
            {
                "ty": 6,
                "tyName": "LogExecTransfer",
                "log": {
                    "execAddr": "1GHTqnmjf3o9c3ZkSq65aXcmhgVzoKEDor",
                    "prev": {
                        "currency": 0,
                        "balance": "4255797000",
                        "frozen": "0",
                        "addr": "1GHTqnmjf3o9c3ZkSq65aXcmhgVzoKEDor"
                    },
                    "current": {
                        "currency": 0,
                        "balance": "4226539028",
                        "frozen": "0",
                        "addr": "1GHTqnmjf3o9c3ZkSq65aXcmhgVzoKEDor"
                    }
                },
                "rawLog": "0x0a2231474854716e6d6a66336f3963335a6b537136356158636d6867567a6f4b45446f72122a10889ea9ed0f222231474854716e6d6a66336f3963335a6b537136356158636d6867567a6f4b45446f721a2a1094bcafdf0f222231474854716e6d6a66336f3963335a6b537136356158636d6867567a6f4b45446f72"
            },
            {
                "ty": 6,
                "tyName": "LogExecTransfer",
                "log": {
                    "execAddr": "1GHTqnmjf3o9c3ZkSq65aXcmhgVzoKEDor",
                    "prev": {
                        "currency": 0,
                        "balance": "9995744203000",
                        "frozen": "0",
                        "addr": "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
                    },
                    "current": {
                        "currency": 0,
                        "balance": "9995773460972",
                        "frozen": "0",
                        "addr": "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
                    }
                },
                "rawLog": "0x0a2231474854716e6d6a66336f3963335a6b537136356158636d6867567a6f4b45446f72122c10f8a1a186f5a202222231344b454b6259744b4b516d34774d7468534b394a344c61346e41696964476f7a741a2c10ec839b94f5a202222231344b454b6259744b4b516d34774d7468534b394a344c61346e41696964476f7a74"
            },
            {
                "ty": 3333,
                "tyName": "LogFluckyBet",
                "log": {
                    "index": 10,
                    "addr": "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt",
                    "time": "1546047584",
                    "amount": 5,
                    "randNum": [
                        "9799",
                        "5131",
                        "6349",
                        "2379",
                        "9777",
                        "6909",
                        "8588",
                        "1273",
                        "5993",
                        "2009"
                    ],
                    "maxNum": "9799",
                    "bonus": 0.2925797,
                    "action": "3300"
                },
                "rawLog": "0x080a122231344b454b6259744b4b516d34774d7468534b394a344c61346e41696964476f7a7418e0a09be10520052a14c74c8b28cd31cb12b14cfd358c43f909e92ed90f30c74c3d02cd953e40e419"
            }
        ]
    },
    "height": 12,
    "index": 0,
    "blocktime": 1546047584,
    "amount": "0.0000",
    "fromaddr": "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt",
    "actionname": "bet",
    "assets": null
}

```

### 1.2 查询用户投注次数

```bash
[azrael@localhost build]$ curl --data-binary '{"jsonrpc":"2.0", "id": 1, "method":"Chain33.Query","params":[{"execer":"flucky", "funcName":"QueryBetTimes", "payload":{"addr":"14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"}}]} '         -H 'content-type:text/plain;'         http://localhost:8801
{"id":1,"result":{"times":10},"error":null}

```

### 1.3 查询用户投注信息

```bash
[azrael@localhost build]$ curl --data-binary '{"jsonrpc":"2.0", "id": 1, "method":"Chain33.Query","params":[{"execer":"flucky", "funcName":"QueryBetInfo", "payload":{"addr":"14KEKbYtKKQm4wMthSK9J4La4nAiidGozt", "idx": 10}}]} '         -H 'content-type:text/plain;'         http://localhost:8801
{"id":1,"result":{"index":10,"addr":"14KEKbYtKKQm4wMthSK9J4La4nAiidGozt","time":"1546047584","amount":5,"randNum":["9799","5131","6349","2379","9777","6909","8588","1273","5993","2009"],"maxNum":"9799","bonus":0.2925797},"error":null}
```

### 1.4 批量查询查询用户投注信息

```bash
[azrael@localhost build]$ curl --data-binary '{"jsonrpc":"2.0", "id": 1, "method":"Chain33.Query","params":[{"execer":"flucky", "funcName":"QueryBetInfoBatch", "payload":{"addr":"1ZdRvtXY2FAa79BaM12owHMJGK9w4S8Ef", "index": 0, "count": 3, "direction": 0}}]} '         -H 'content-type:text/plain;'         http://localhost:8801
{"id":1,"result":{"bets":[{"index":"4","addr":"1ZdRvtXY2FAa79BaM12owHMJGK9w4S8Ef","time":"1546592716","amount":"10","randNum":["9699","742","644","9195","5511","7407","909","6013","9404","6813"],"maxNum":"9699","bonus":7.6438},{"index":"3","addr":"1ZdRvtXY2FAa79BaM12owHMJGK9w4S8Ef","time":"1546592656","amount":"5","randNum":["1341","8228","1014","371","225","2023","3140","1211","3876","9755","7726","9831"],"maxNum":"9831","bonus":7.62},{"index":"2","addr":"1ZdRvtXY2FAa79BaM12owHMJGK9w4S8Ef","time":"1546589846","amount":"5","randNum":["215","8803","4368","8273","5465","2357","9987","459","9999","6073"],"maxNum":"9999","bonus":257}]},"error":null}
```
### 1.5 查询奖池信息

```bash
[azrael@localhost build]$ curl --data-binary '{"jsonrpc":"2.0", "id": 1, "method":"Chain33.Query","params":[{"execer":"flucky", "funcName":"QueryBonusInfo"}]} '         -H 'content-type:text/plain;'         http://localhost:8801
{"id":1,"result":{"userCount":1,"bonusPool":756.7362},"error":null}
```

