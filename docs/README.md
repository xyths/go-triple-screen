# 三重滤网交易系统(triple screen trading system)

## 逻辑

1. 定时扫描K线
    1. 如果周线状态没有初始化，则初始化周线状态
    2. 如果是周线时间(周一0时)，重新判断周线状态
    3. 判断是否可以交易，尝试发出交易信号
2. 异步处理交易信号，从短周期寻找交易机会
    1. 如果遇到终止交易信号，则停止所有买单挂单，并清空已有头寸
    2. 如果是买入交易信号，则重复判断短周期K线，尝试买入（可能不成功）
3. 下单并监听成交情况，实时调整状态
    1. 参考状态机

**买多状态机**

![state_machine_long](state_machine_long.png)

**卖空状态机**

![state_machine_short](state_machine_short.png)


**周线状态**

上涨、下跌或中性，以及止盈位。

**日线状态**

观望、买入和卖出，以及止损位。

此时已经得到止盈位和止损位，按照盈亏比2:1的原则，可以计算出入场价位，该价位可以用作后续买入技术的参考。

### 状态机


## 策略


### 第一重滤网-市场潮流(First Screen - Market Tide)

第一重滤网是在长周期的K线图上判断大趋势，顺应趋势的交易才被允许，违背趋势的交易将被禁止。

### 第二重滤网-市场波浪(Second Screen - Market Wave)

第二重滤网实际是寻找回调的机会，并发出交易信号。

### 第三重滤网-买入方法(Third Screen - Entry Technique)

#### 方法1：第1版的简易方法

- 做多时，在日线向上突破前一交易日高点的时候买入；
- 卖空时，在日线向下突破前一交易日低点的时候卖出。

**缺点**
这种方法的缺点是止损点的位置太远了。在突破前一交易日高点的位置买入，同时在前一日的低点位置设置止损，意味着如果前一日的价格振幅很大，则止损价离最新的股价很远。这样要承担很大的风险，否则只能用很小的头寸进行交易。
另一种风险是，当突破前一交易日振幅很窄，将止损点刚好设在前一交易日的低点之下时，当日的市场噪声就可能会触发止损

#### 方法2：平均EMA下跌穿透


1. 计算平均穿透值
2. 用今日的EMA值减去昨日的EMA值，将其结果加回今日的EMA值：这是对明日EMA值的一个估算。用估算的明日EMA值减去你计算的平均穿透值，作为明日设置买入订单的触发价位。
3. 你将利用回调以折扣价完成买入交易——避免了在突破时买入须支付的溢价。

**总结**

| 周趋势 | 日趋势 | 行动 | 指令 |
| --- | --- | --- | :--- |
| 上行 | 上行 | 观望 | 无 |
| 上行 | 下行 | 买入 | EMA穿透或者向上突破 |
| 下行 | 下行 | 观望 | 无 |
| 下行 | 上行 | 卖出 | EMA穿透或者向下突破 |


## 回测

回测分BTC现货（单向做多）和永续合约（多空），对比二者结果，看差异多少。