# pbft

PBFT是一类状态机拜占庭系统，要求整个系统共同维护一个状态，所有节点采取的行动一致。为此，需要运行三类基本协议，包括一致性协议、检查点协议和视图更换协议。

## 一致性协议

一致性协议要求来自客户端的请求在每个服务节点上都按照一个确定的顺序执行。这个协议把服务器节点分为两类：主节点和从节点，其中主节点仅一个。在协议中，主节点负责将客户端的请求排序，从节点按照主节点提供的顺序执行请求。每个服务器节点在同样的配置信息下工作，该配置信息被称为视图（view），主节点更换，视图也随之变化。

 一致性协议至少包含五个阶段：请求(request)、序号分配(pre-prepare)、相互交互(prepare) 、 序号确认(commit)和响应(reply)等阶段。

![pbft流程图](/home/lyn/Desktop/work/L/share/img/pbft.png)

 Request阶段：客户端发送请给主节点，请求消息m=[op, ts, c-id, c-sig], 其中包含需要执行的操作op, 时间戳ts, cs-id是客户端ID，以及客户端签名 c-sig。时间戳是为了保证命令只被执行一次，客户端的签名是为了客户认证和权限控制。

 Pre-Prepare 阶段：当主节点接收请求后，主节点给请求赋值一个序列号sn，广播序号分配消息和客户端的请求消息m，并将构造PRE-PREPARE消息[PP, vn, sn, D(m), p-sig, m]给各从节点，其中PP表示Pre-Prepare消息，vn是视图号，D(m) 是客户消息的摘要（哈希值），p-sig是主节点的签名，m是客户消息。序列号是为了保证命令执行的顺序，视图号让从节点记录当前视图。主节点签名是为了让从节点认证主节点身份。消息摘要保证消息没有被篡改。

 Prepare阶段：从节点接收PRE-PREPARE消息，向其他服务节点广播PREPARE消息[P, vn, sn, D(m), b-id, b-sig]，其中P表示Prepare消息，b-id是从节点id，b-sig是从节点签名。

 Commit阶段：从节点在收到2f + 1个Prepare消息后，对视图内的请求和次序进行验证，广播COMMIT消息[C, vn, sn, D(m), b-id, b-sig]，其中C表示Commit消息。执行收到的客户端的请求并给客户端以响应

 Response阶段：（1）当各节点收到2f + 1个COMMIT消息后，它将执行操作，并提交。同时把回复返回给客户端。回复消息是[R, vn, ts, b-sig], 其中R表示回复消息。（2）客户端等待来自不同节点的响应，若有f + 1个响应相同，则该响应即为运算的结果。

## 检查点协议


## 视图更换协议
视图转换协议保证共识协议的活性（liveness）。当主节点出故障时能保证共识能继续进行。每个备份节点收到一个请求是都会开始一个定时器。如果在一个视图内，时钟超时，该备份节点会发起一个视图转换的消息。它将广播一个<view-change, v+1, n, C, P, i>加上签名的消息给各节点。其中n是当前节点i所知道的上个稳定检查点状态s的序号，C是一个集合，里面有2f+1个检查点证明s是正确消息，然后P是一个Pm的集合，Pm里面有每个从i节点中大于n的请求消息m的pre-prepare, 和2f个匹配的被各节点签名的prepare消息。

> 定时器是在收到什么请求是开始的？如果是pre-prepare，那在收到后续消息时是否会重置定时器。

当v+1视图的主节点收到其它2f个节点的正确的view-change消息后，它将广播<new-view, v+1, V, O>加上签名的消息给所有节点。其中V是主节点收到的view-change消息，加上它自己发的view-change消息，O是一个pre-prepare的消息集合。O是为每个在从低水位标志到高水位标志之间序号n建立的新的pre-prepare的消息。这样新的主节点可以继续的推进共识。

> v+1视图的主节点是如何确定的？

在这里，我们可以看到，PBFT的视图转换协议是非常复杂的，涉及到很多消息的重传。HotStuff的最重要的改进，主要是针对视图更换的协议。





# HotStuff

**HotStuff协议分为基本HotStuff（Basic HotStuff）和链式HotStuff（Chained HotStuff）协议**。



## 基本HotStuff



## 链式HotStuff



### QcTree

| 字段名称   | 类型          | 含义                                                         |
| ---------- | ------------- | ------------------------------------------------------------ |
| Genesis    | *ProposalNode | 创世节点，共识启动时根据共识配置相应高度的区块进行填充。     |
| Root       | *ProposalNode | 初始化节点，启动时根据当前节点的最大高度进行加载以及后续更新Commit node时更新为 |
| HighQC     | *ProposalNode | Tree中最高的QC指针                                           |
| GenericQC  | *ProposalNode |                                                              |
| LockedQC   |               |                                                              |
| CommitQC   |               |                                                              |
| OrphanList |               |                                                              |
| OrphanMap  |               |                                                              |







# HSL2.0 