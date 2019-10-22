# 以太坊创建Token
***本文档基于Rinkeby测试网络***

## 以太坊账户

### 通过MetaMask创建账户

![](https://github.com/lynAzrael/L/blob/master/share/img/create_account_metamask.png)

### 私钥导出

进入到处私钥界面

![step 1](https://github.com/lynAzrael/L/blob/master/share/img/prepare_dump_key.png)

输入创建用户时设置的密码

![](https://github.com/lynAzrael/L/blob/master/share/img/passwd_dump_key.png)

记录好自己的私钥信息

![](https://github.com/lynAzrael/L/blob/master/share/img/dump_key_info.png)

***创建账户地址之后，将私钥保管好。私钥一旦丢失，账户中的财产就无法找回。***

## 测试币的获取
[https://www.rinkeby.io/#faucet](https://www.rinkeby.io/#faucet)

使用收币地址生成链接

![](https://github.com/lynAzrael/L/blob/master/share/img/create_public_gist.png)

使用生成的链接获取测试Ether

![](https://github.com/lynAzrael/L/blob/master/share/img/rinkeby_authenticated_faucet.png)

视频参考：[https://www.youtube.com/watch?v=wKFz5c3TU4s](https://www.youtube.com/watch?v=wKFz5c3TU4s)


## Token合约

### 合约部署
使用remix进行合约的编译以及部署

![compile contract](https://github.com/lynAzrael/L/blob/master/share/img/contract_compile.png)

设置总发行量，精度，token名称以及token的标记

|字段|取值|
|----|----|
|totalSupply|21000000|
|decimalUnits|3|
|tokenName|ASD|
|tokenSymbol|asd|

部署合约的交易

![](https://github.com/lynAzrael/L/blob/master/share/img/contract_transfer_tx.png)

通过区块浏览器查看交易的具体信息

浏览器地址: [https://ropsten.etherscan.io](https://ropsten.etherscan.io)

> 本次合约部署位于Rinkeby测试网，因此使用Rinkeby的区块链浏览器。主网的区块链浏览器地址 [https://cn.etherscan.com](https://cn.etherscan.com)

![contract_tx_detail](https://github.com/lynAzrael/L/blob/master/share/img/contract_tx_detail.png)



### 合约调用



