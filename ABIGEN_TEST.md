# abigen 合约交互测试验证

## 前置条件

- abigen安装
- 已部署 Counter 合约到 Sepolia 测试网

## 环境变量

```powershell
$env:ETH_RPC_URL = "https://sepolia.infura.io/v3/YOUR_KEY"
$env:SENDER_PRIVATE_KEY = "你的私钥（不带0x前缀）"
$env:CONTRACT_ADDRESS = "0x部署后的合约地址"
```

## 部署合约（Remix）

1. 打开 [Remix](https://remix.ethereum.org)
2. 新建 `Counter.sol`，粘贴 `contracts/Counter/Counter.sol` 的内容
3. 编译（Solidity Compiler → Compile）
4. 部署（Deploy → Environment 选 "Browser Extension" → 选 MetaMask 账户 → Deploy）
5. 复制部署后的合约地址，设为 `CONTRACT_ADDRESS`

> ⚠️ 钱包需要有 Sepolia 测试 ETH，可从 https://www.alchemy.com/faucets/ethereum-sepolia 领取

## 运行与调用

### 方式一：go run 直接运行（推荐）

```powershell
cd E:\workplace\Web3\homework\golang-block-operate
go run . -mode demo-counter
```

### 方式二：编译后运行

```powershell
cd E:\workplace\Web3\homework\golang-block-operate
go build -o golang-block-operate.exe .
.\golang-block-operate.exe -mode demo-counter
```

### 只测试只读查询（不需私钥）

如果只设了 `ETH_RPC_URL` 和 `CONTRACT_ADDRESS`，程序会执行 `GetCount` 后跳过写操作，输出：

```
✅ 已连接到 Sepolia 测试网
✅ 已绑定合约: 0x...

━━━━━━━━━━ 查询计数器值 ━━━━━━━━━━
当前值: 0

⚠️  未设置私钥，跳过写操作测试
```

### 完整测试（含写入交易）

设置全部三个环境变量后运行，将依次执行：GetCount → Increment → GetCount → SetCount(42) → GetCount → Reset → GetCount

## 预期输出

```
✅ 已连接到 Sepolia 测试网
✅ 已绑定合约: 0x...

━━━━━━━━━━ 查询计数器值 ━━━━━━━━━━
当前值: 0

━━━━━━━━━━ increment ⏫ ━━━━━━━━━━
交易已发送!
交易哈希: 0x...
查看详情: https://sepolia.etherscan.io/tx/...

━━━━━━━━━━ 再次查询 ⏳ ━━━━━━━━━━
等待交易确认...
当前值: 1 (应该是 1)

━━━━━━━━━━ setCount(42) 🔢 ━━━━━━━━━━
交易已发送!
交易哈希: 0x...
当前值: 42 (应该是 42)

━━━━━━━━━━ reset 🔄 ━━━━━━━━━━
交易已发送!
交易哈希: 0x...
当前值: 0 (应该是 0)

🎉 所有测试完成！
```

## 测试验证清单

| 步骤 | 操作 | 验证点 |
|------|------|--------|
| 1 | 连接 RPC | `ethclient.Dial` 成功 |
| 2 | 绑定合约 | `NewCounter` 返回实例无错误 |
| 3 | GetCount | 返回当前计数值（只读，不需私钥） |
| 4 | Increment | 交易发送成功，值 +1 |
| 5 | SetCount(42) | 交易发送成功，值变为 42 |
| 6 | Reset | 交易发送成功，值归零 |

## abigen 生成的 API 速查

```go
// 创建合约实例
counter, _ := counterpkg.NewCounter(common.HexToAddress(addr), client)

// 只读调用
count, _ := counter.GetCount(&bind.CallOpts{Context: ctx})  // → *big.Int

// 写入交易
auth, _ := bind.NewKeyedTransactorWithChainID(privKey, chainID)
tx, _ := counter.Increment(auth)       // → *types.Transaction
tx, _ := counter.Reset(auth)           // → *types.Transaction
tx, _ := counter.SetCount(auth, value) // → *types.Transaction（value *big.Int）
```
