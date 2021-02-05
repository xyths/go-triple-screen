
```mermaid
stateDiagram-v2
    [*] --> Buying: Long
    Buying --> LongPosition: Buy Order Full-Filled
    LongPosition --> Selling: Short
    Selling --> [*]: Sell Order Full-Filled
```

```mermaid
stateDiagram-v2
    [*] --> Selling: Short
    Selling --> ShortPosition: Sell Order Full-Filled
    ShortPosition --> Buying: Long
    Buying --> [*]: Buy Order Full-Filled
```