



## 传统密钥管理系统（KMS）

![Centralized KMS](https://i.imgur.com/X3TKudW.png)

## 代理重加密(Proxy re-encryption)

![proxy re-recryption](https://i.imgur.com/oX8Y19m.png)

数据持有方: pkA,skA

数据接收方: pkB,skB



## 基于PRE的密钥管理系统

![PRE KMS](https://i.imgur.com/xxqiWu7.png)

## Nucypher中的应用

![](https://i.imgur.com/eTlUccW.png)

### 数据加密

dek=random()     // 生成随机对称加密密钥

c=encryptSym(dek, d) // 使用密钥对铭文进行对称加密

edek=encryptPke(pka, dek) // 使用发送方的公钥进行加密

### 访问委派

reS->R=reKey(skA, pkB) // 使用发送方的私钥和接收方的公钥构造重加密密钥

### 数据解密

edek'=reEncrypt(reS->R, edek) // 使用重加密密钥对密文进行转换

dek=decryptPke(skB, edek') // 使用接收方的私钥进行解密

d=decryptSym(dek, c) // 使用解出的对称密钥进行解密

### 分布式存储

![](https://i.imgur.com/EgjbaKq.png)

在对密钥片段进行重加密时，会产生多个碎片并分发到多个节点中。

在重加密时，可以制定生成碎片的总数量，以及验证一个片段所需的最小阈值。

```go
kfrags = pre.generate_kfrags(delegating_privkey=alices_private_key,
                             signer=alices_signer,
                             receiving_pubkey=bobs_public_key,
                             threshold=10,
                             N=20)
```

