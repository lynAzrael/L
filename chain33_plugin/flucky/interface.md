# Flucky游戏

<!-- TOC -->

- [Flucky游戏](#flucky游戏)
    - [1 投注](#1-投注)
    - [2 查询接口](#2-查询接口)
        - [2.1 用户投注信息查询](#21-用户投注信息查询)
        - [2.2 用户投注信息批量查询](#22-用户投注信息批量查询)
        - [2.3 奖池信息查询](#23-奖池信息查询)

<!-- /TOC -->

## 1 投注
请求报文

```json
{
    "jsonrpc": "2.0",
    "method":"flucky.FluckyBetTx",
    "params":[
    {
        "amount":1
    }],
    "id":0
}
```
参数说明

|参数|类型|说明|
|----|----|----|
|amount|int64|用户投注的额度|

响应报文

```bash
{
    "id":1,
	"result":"0a06666c75636b79120730e4190a02080520e80730b4f997edf6e2c6db093a2231474854716e6d6a66336f3963335a6b537136356158636d6867567a6f4b45446f72",
    "error":null
}
```

## 2 查询接口

### 2.1 用户投注信息查询

```json
{
    "jsonrpc":"2.0",
    "id": 1, 
    "method":"Chain33.Query",
    "params":[
    {
        "execer":"flucky", 
        "funcName":"QueryBetInfo", 
        "payload":{
            "addr":"14KEKbYtKKQm4wMthSK9J4La4nAiidGozt", 
            "idx": 10
        }
    }]
}
```

参数说明

|参数|类型|说明|
|----|----|----|
|addr|string|查询的用户地址|
|idx|int|用户投注的index|

响应报文
```json
{
    "id":1,
    "result":{
        "index":10,
        "addr":"14KEKbYtKKQm4wMthSK9J4La4nAiidGozt",
        "time":"1546047584",
        "amount":5,
        "randNum":[
            "9799",
            "5131",
            "6349",
            "2379",
            "9777",
        ],
        "maxNum":"9799",
        "bonus":0.2925797
    },
    "error":null
}
```

参数说明

|参数|类型|说明|
|----|----|----|
|index|int|用户投注的index|
|addr|string|投注用户的地址|
|time|string|用户投注时间|
|amount|int|用户投注的额度|
|randNum|string数组|用户本次抽奖获取到的随机数|
|maxNum|string|用户本次获取到的最大随机数|
|bonus|int|用户本次获取到的奖金|

### 2.2 用户投注信息批量查询

请求报文

```json
{
    "jsonrpc":"2.0", 
    "id": 1, 
    "method":"Chain33.Query",
    "params":[{
        "execer":"flucky", 
        "funcName":"QueryBetInfoBatch", 
        "payload":{
            "addr":"1ZdRvtXY2FAa79BaM12owHMJGK9w4S8Ef", 
            "index": 0, 
            "count": 3, 
            "direction": 0
        }
    }]
} 
```

参数说明

|参数|类型|说明|
|----|----|----|
|addr|string|需要查询的用户地址|
|index|int|批量查询的起始index|
|count|int|查询的数量|
|direction|int|查询的方向，0：逆序，1：正序 默认为1|

响应报文

```json
{
    "id":1,
    "result":{
        "bets":[{
            "index":"4",
            "addr":"1ZdRvtXY2FAa79BaM12owHMJGK9w4S8Ef",
            "time":"1546592716",
            "amount":"10",
            "randNum":["9699","742","644","9195","5511","7407","909","6013","9404","6813"],
            "maxNum":"9699",
            "bonus":7.6438
        },
        {
            "index":"3",
            "addr":"1ZdRvtXY2FAa79BaM12owHMJGK9w4S8Ef",
            "time":"1546592656",
            "amount":"5",
            "randNum":["1211","3876","9755","7726","9831"],
            "maxNum":"9831",
            "bonus":7.62
        },
        {
            "index":"2",
            "addr":"1ZdRvtXY2FAa79BaM12owHMJGK9w4S8Ef",
            "time":"1546589846",
            "amount":"5",
            "randNum":["2357","9987","459","9999","6073"],
            "maxNum":"9999",
            "bonus":257
        }]
    },
    "error":null
}
```

参数说明

|参数|类型|说明|
|----|----|----|
|index|int|用户投注的index|
|addr|string|投注用户的地址|
|time|string|用户投注时间|
|amount|int|用户投注的额度|
|randNum|string数组|用户本次抽奖获取到的随机数|
|maxNum|string|用户本次获取到的最大随机数|
|bonus|int|用户本次获取到的奖金|



### 2.3 奖池信息查询

请求报文

```json
{
    "jsonrpc":"2.0", 
    "id": 1, 
    "method":"Chain33.Query",
    "params":[{
        "execer":"flucky", 
        "funcName":"QueryBonusInfo"
    }]
}
```

响应报文
```json
{
    "id":1,
    "result":{
        "userCount":1,
        "bonusPool":756.7362},
    "error":null
}
```

参数说明

|参数|类型|说明|
|----|----|----|
|userCount|int|奖池中的用户数量|
|bonusPool|int|奖池中累计的奖金数量|