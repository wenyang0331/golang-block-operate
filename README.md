# golang-block-operate

使用 Go 语言和 `go-ethereum` 库与 Sepolia 测试网络交互，实现区块查询和 ETH 转账交易。

## 环境配置（Windows）

### 1. 设置 RPC 地址

```cmd
set ETH_RPC_URL=https://sepolia.infura.io/v3/YOUR-PROJECT-ID
```

> 可以从 [Infura](https://www.infura.io/) 或 [Alchemy](https://www.alchemy.com/) 获取免费的 Sepolia RPC 地址。

### 2. 设置发送交易私钥（仅发送交易时需要）

```cmd
set SENDER_PRIVATE_KEY=your_private_key_hex
```

```powershell
# 设置变量
$env:ETH_RPC_URL = "https://sepolia.infura.io/v3/YOUR-PROJECT-ID"

#测试网中的ETH_RPC_URL格式
#$env:ETH_RPC_URL = "https://sepolia.infura.io/v3/9eca82c008204b7fb3a62e6aeb77a02e"
# 注意:测试网中的 $env:ETH_RPC_URL =  "https://sepolia.infura.io/v3/YOUR-PROJECT-ID"
#YOUR-PROJECT-ID为API KEY值


# 查询变量
echo $env:ETH_RPC_URL

# 输出：https://sepolia.infura.io/v3/YOUR-PROJECT-ID
```

> 私钥支持带或不带 `0x` 前缀。**仅在测试网使用，切勿泄露主网私钥。**

### 3. 查看环境变量是否设置成功

```cmd
echo %ETH_RPC_URL%
echo %SENDER_PRIVATE_KEY%
```

## 功能一：查询区块信息

### 查询最新区块

```cmd
go run main.go --mode query-block
```

### 查询指定区块号

```cmd
go run main.go --mode query-block --block 7780000
```

### 输出示例

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Block Information
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Block Number    : 7780000
Block Hash      : 0xabc123...def456
Timestamp       : 1748888888 (2026-06-02 18:28:08)
Transactions    : 85
Parent Hash     : 0x...
State Root      : 0x...
Gas Limit       : 30000000
Gas Used        : 12345678
Base Fee        : 12345678 Wei
Miner/Validator : 0x...
Nonce           : 0
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Transaction Hashes (85 total, showing first 5):
  [1] 0x...
  [2] 0x...
  [3] 0x...
  [4] 0x...
  [5] 0x...
  ... and 80 more transactions
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## 功能二：发送 ETH 转账交易

### 发送交易

```cmd
go run main.go --mode send-tx --to 0xRecipientAddress --amount 0.01
```

### 输出示例

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
ETH Transfer Transaction Sent
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
From          : 0xSenderAddress...
To            : 0xRecipientAddress...
Amount        : 0.01 ETH (10000000000000000 Wei)
Gas Limit     : 21000
Gas Tip Cap   : 1500000000 Wei
Gas Fee Cap   : 30000000000 Wei
Estimated Cost: 630000000000000 Wei (~0.000630000 ETH)
Nonce         : 5
Chain ID      : 11155111
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ Transaction Hash:
   0xtxhash...

Check on Sepolia Explorer:
   https://sepolia.etherscan.io/tx/0xtxhash...
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## 参数说明

| 参数 | 说明 | 默认值 | 适用模式 |
|------|------|--------|----------|
| `--mode` | 操作模式：`query-block` 或 `send-tx` | `query-block` | 两者 |
| `--block` | 要查询的区块号（-1 表示最新） | `-1` | `query-block` |
| `--to` | 接收方地址 | - | `send-tx` |
| `--amount` | 转账金额（单位：ETH） | - | `send-tx` |

## 环境变量

| 变量 | 说明 | 必需 |
|------|------|------|
| `ETH_RPC_URL` | Sepolia 测试网 RPC 地址 | 始终必需 |
| `SENDER_PRIVATE_KEY` | 发送方私钥（支持 `0x` 前缀） | `send-tx` 模式必需 |

## 注意事项

- 所有操作在 **Sepolia 测试网** 执行，不影响主网资产
- 私钥务必妥善保管，不要在公开场合泄露
- Windows 命令行中 `set` 设置的环境变量仅当前会话有效，关闭窗口后需重新设置
