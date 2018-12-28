# Feel Lucky 游戏合约

## 1 存储
### 1.1 statedb

```proto
// 投注信息
message BetInfo {
	int32 index = 1;
}

// 当前奖池信息
message BonusInfo {
	int32 userCount = 1;
	float bonusPool = 2; 
}
```

|key|val|说明|
|----|----|----|
|mavl-flucky-user-times:{address}|BetInfo|用户购买次数|
|mavl-flucky-bonul-info|BonusInfo|当前奖池信息|

### 1.2 localdb

```proto
// 用户投注信息
message BetInfo {
	int32 index = 1;
    uint64 time = 2;
    int32 amount = 3;
    repeated int32 randNum = 4;
    int64 bonus = 5;
}
```

|key|val|说明|
|----|----|----|
|LODB-flucky-user-history:{address}:{idx}|BetInfo|用户投注信息|



## 2 接口

```proto
// 投注
message BetInfo {
    int32 index = 1;
}
```

|字段|类型|说明|
|amount|int32|投注的数量|


对外查询接口：

```proto
// 用户投注信息批量查询
message QueryBetInfo {
	string addr = 1;
	int32 startIdx = 2;
	int32 endIdx = 3;
}

// 奖金信息查询
message QueryRewardInfo {
	string addr = 1;
	int32 idx = 2;
}

// 当前奖池信息查询
message QueryBonusInfo {
	
}

// 用户投注次数查询
message QueryBetTimes {
	string addr = 1;
}
```

响应数据
```proto
// 奖金信息
message RewardInfo{
   	repeated int32 randNum = 1;
    float bouns = 2;
}

// 奖池信息
message BonusInfo {
    int32 userCount = 1;
    float bonusPool = 2;
}

// 投注次数信息
message BetTimes {
	int32 times = 1;
}
```

## 3 奖池维护
另外起一个进程进行维护，根据改配置文件中的奖池最小、最大限额进行

## 4 合约配置
```toml
# 管理员地址
managerAddr="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"

# 奖池地址
poolAddr="j4pFNSjfemZOyn37A1Cnt2FrEq"

# 用于生成随机的块数
randLuckyBlockNum=5

```