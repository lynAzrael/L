# Flucky游戏

<!-- TOC -->

- [Flucky游戏](#flucky游戏)
    - [1 投注](#1-投注)
    - [2 查询接口](#2-查询接口)
        - [2.1 用户投注信息查询](#21-用户投注信息查询)
        - [2.2 用户投注信息批量查询](#22-用户投注信息批量查询)
        - [2.3 奖池信息查询](#23-奖池信息查询)
        - [2.4 用户投注次数查询](#24-用户投注次数查询)

<!-- /TOC -->

## 1 投注
请求报文

```json
{
    "jsonrpc":"2.0", 
    "id": 1, 
    "method":"Chain33.CreateTransaction",
    "params":[{
        "execer":"flucky", 
        "actionName":"Bet", 
        "payload":{
            "amount":5
        }
    }]
}
```
参数说明

|参数|类型|说明|
|----|----|----|
|amount|int|用户投注的额度|

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
**根据index，查询指定用户某一次的投注信息**

请求报文

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
|idx|int|指定用户投注的次数，index为10则查询用户第10次投注的信息。(不同用户之间index相互独立)|

响应报文
```json
{
    "id":1,
    "result":{
        "index":"10",
        "addr":"1Gk428DWuhg9kpUJtoNLCadjCLsgiDLGMU",
        "time":"1548063710",
        "amount":"5",
        "randNum":["6029"],
        "maxNum":"6029",
        "bonus":2.5,
        "action":"3300"
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
|bonus|float|用户本次获取到的奖金|

### 2.2 用户投注信息批量查询
**根据指定的起始index,addr,count,direction等信息批量查询用户信息**

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
            "index":"13756",
            "addr":"1Gk428DWuhg9kpUJtoNLCadjCLsgiDLGMU",
            "time":"1485051134",
            "amount":"5",
            "randNum":["2183"],
            "maxNum":"2183",
            "bonus":1,
            "action":"3300"
        },
        {
            "index":"13755",
            "addr":"1Gk428DWuhg9kpUJtoNLCadjCLsgiDLGMU",
            "time":"1485051134",
            "amount":"5",
            "randNum":["5504"],
            "maxNum":"5504",
            "bonus":2.5,
            "action":"3300"
        },
        {
            "index":"13754",
            "addr":"1Gk428DWuhg9kpUJtoNLCadjCLsgiDLGMU",
            "time":"1485051134",
            "amount":"5",
            "randNum":["1082"],
            "maxNum":"1082",
            "bonus":1,
            "action":"3300"
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
|maxNum|string|用户本次获取到的最大随机数，用于计算奖金|
|bonus|float|用户本次获取到的奖金|



### 2.3 奖池信息查询

**查询奖池中投注的用户数量以及总奖金额度**

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
        "bonusPool":756.7362
    },
    "error":null
}
```

参数说明

|参数|类型|说明|
|----|----|----|
|userCount|int|奖池中的用户数量|
|bonusPool|float|奖池中累计的奖金数量|


### 2.4 用户投注次数查询

**查询指定用户投注的总次数**

请求报文

```json
{
    "jsonrpc":"2.0", 
    "id": 1, 
    "method":"Chain33.Query",
    "params":[{
        "execer":"flucky", 
        "funcName":"QueryBetTimes", 
        "payload":{
            "addr":"14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
        }
    }]
} 
```

参数说明

|参数|类型|说明|
|----|----|----|
|addr|string|查询的用户地址|

响应报文
```json
{
    "id":1,
    "result":{
        "times":10
    },
    "error":null
}
```

参数说明

|参数|类型|说明|
|----|----|----|
|times|int|用户投注的总次数|

