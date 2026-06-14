package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 使用示例：
//
// 1. 查询指定区块号的区块信息：
//    export ETH_RPC_URL="https://sepolia.infura.io/v3/YOUR-PROJECT-ID"

//测试网中的ETH_RPC_URL格式
//$env:ETH_RPC_URL = "https://sepolia.infura.io/v3/9eca82c0xxxxxxxxfb3a62e6aeb77a02e"
//注意:测试网中的 $env:ETH_RPC_URL =  "https://sepolia.infura.io/v3/YOUR-PROJECT-ID"
//YOUR-PROJECT-ID为API KEY值

//    go run main.go --mode query-block --block 12345678
//
// 2. 查询最新区块信息：
//    export ETH_RPC_URL="https://sepolia.infura.io/v3/YOUR-PROJECT-ID"
//    go run main.go --mode query-block
//
// 3. 发送以太币转账交易：
//    export ETH_RPC_URL="https://sepolia.infura.io/v3/YOUR-PROJECT-ID"
//    export SENDER_PRIVATE_KEY="your_private_key_hex"
//    go run main.go --mode send-tx --to 0xRecipientAddress --amount 0.01
//
// 注意事项：
// - 私钥可带或不带 0x 前缀
// - 仅在测试网使用，不要在主网使用包含真实资产的私钥

func main() {
	mode := flag.String("mode", "query-block", "operation mode: query-block or send-tx")
	blockNum := flag.Int64("block", -1, "block number to query (default: latest)")
	toAddr := flag.String("to", "", "recipient address (for send-tx)")
	amount := flag.String("amount", "", "amount of ETH to send (for send-tx)")
	flag.Parse()

	rpcURL := os.Getenv("ETH_RPC_URL")
	if rpcURL == "" {
		log.Fatal("ETH_RPC_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		log.Fatalf("failed to connect to Ethereum node: %v", err)
	}
	defer client.Close()

	switch *mode {
	case "query-block":
		queryBlock(ctx, client, *blockNum)
	case "send-tx":
		sendTransaction(ctx, client, *toAddr, *amount)
	case "demo-counter":
		demoCounter()
	default:
		log.Fatalf("unknown mode: %s (use: query-block, send-tx, or demo-counter)", *mode)
	}
}

// ============================================================
// 功能一：查询区块信息
// ============================================================

func queryBlock(ctx context.Context, client *ethclient.Client, blockNum int64) {
	var block *types.Block
	var err error

	if blockNum < 0 {
		// 查询最新区块
		block, err = client.BlockByNumber(ctx, nil)
	} else {
		// 查询指定区块号
		block, err = client.BlockByNumber(ctx, big.NewInt(blockNum))
	}

	if err != nil {
		log.Fatalf("failed to query block: %v", err)
	}

	// 将区块时间戳转换为可读格式
	blockTime := time.Unix(int64(block.Time()), 0)

	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("Block Information\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("Block Number    : %d\n", block.Number().Uint64())
	fmt.Printf("Block Hash      : %s\n", block.Hash().Hex())
	fmt.Printf("Timestamp       : %d (%s)\n", block.Time(), blockTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("Transactions    : %d\n", len(block.Transactions()))
	fmt.Printf("Parent Hash     : %s\n", block.ParentHash().Hex())
	fmt.Printf("State Root      : %s\n", block.Root().Hex())
	fmt.Printf("Gas Limit       : %d\n", block.GasLimit())
	fmt.Printf("Gas Used        : %d\n", block.GasUsed())
	fmt.Printf("Base Fee        : %s Wei\n", block.BaseFee().String())
	fmt.Printf("Miner/Validator : %s\n", block.Coinbase().Hex())

	// 显示 Nonce (仅 Legacy 区块有，EIP-1559 后为随机数)
	fmt.Printf("Nonce           : %d\n", block.Nonce())

	// 显示前 5 笔交易的哈希
	txCount := len(block.Transactions())
	if txCount > 0 {
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("Transaction Hashes (%d total, showing first 5):\n", txCount)
		showCount := txCount
		if showCount > 5 {
			showCount = 5
		}
		for i := 0; i < showCount; i++ {
			tx := block.Transactions()[i]
			fmt.Printf("  [%d] %s\n", i+1, tx.Hash().Hex())
		}
		if txCount > 5 {
			fmt.Printf("  ... and %d more transactions\n", txCount-5)
		}
	}
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

// ============================================================
// 功能二：发送以太币转账交易
// ============================================================

func sendTransaction(ctx context.Context, client *ethclient.Client, toHex, amountStr string) {
	if toHex == "" || amountStr == "" {
		log.Fatal("missing --to or --amount flag for send-tx mode")
	}

	// 读取私钥
	privKeyHex := os.Getenv("SENDER_PRIVATE_KEY")
	if privKeyHex == "" {
		log.Fatal("SENDER_PRIVATE_KEY is not set (required for send-tx mode)")
	}

	// 解析私钥
	privKey, err := crypto.HexToECDSA(trim0x(privKeyHex))
	if err != nil {
		log.Fatalf("invalid private key: %v", err)
	}

	// 从私钥派生公钥和发送方地址
	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddr := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 解析接收方地址
	toAddr := common.HexToAddress(toHex)

	// 解析转账金额（单位：ETH -> Wei）
	amountFloat, ok := new(big.Float).SetString(amountStr)
	if !ok {
		log.Fatalf("invalid amount: %s", amountStr)
	}

	// 1 ETH = 10^18 Wei
	weiPerEth := new(big.Float).SetFloat64(1e18)
	amountWeiFloat := new(big.Float).Mul(amountFloat, weiPerEth)
	amountWei, _ := amountWeiFloat.Int(nil)

	if amountWei.Cmp(big.NewInt(0)) <= 0 {
		log.Fatal("amount must be greater than 0")
	}

	// 获取链 ID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("failed to get chain id: %v", err)
	}

	// 获取 nonce
	nonce, err := client.PendingNonceAt(ctx, fromAddr)
	if err != nil {
		log.Fatalf("failed to get nonce: %v", err)
	}

	// 获取建议的 Gas 价格（EIP-1559 动态费用）
	gasTipCap, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		log.Fatalf("failed to get gas tip cap: %v", err)
	}

	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Fatalf("failed to get header: %v", err)
	}

	baseFee := header.BaseFee
	if baseFee == nil {
		// 不支持 EIP-1559，回退到传统 gas price
		gasPrice, err := client.SuggestGasPrice(ctx)
		if err != nil {
			log.Fatalf("failed to get gas price: %v", err)
		}
		baseFee = gasPrice
	}

	// fee cap = base fee * 2 + tip cap
	gasFeeCap := new(big.Int).Add(
		new(big.Int).Mul(baseFee, big.NewInt(2)),
		gasTipCap,
	)

	// 估算 Gas Limit（普通 ETH 转账固定 21000）
	gasLimit := uint64(21000)

	// 检查 ETH 余额是否足够（转账金额 + Gas 费用）
	balance, err := client.BalanceAt(ctx, fromAddr, nil)
	if err != nil {
		log.Fatalf("failed to get balance: %v", err)
	}

	totalGasCost := new(big.Int).Mul(gasFeeCap, big.NewInt(int64(gasLimit)))
	totalNeeded := new(big.Int).Add(amountWei, totalGasCost)

	if balance.Cmp(totalNeeded) < 0 {
		log.Fatalf("insufficient balance: have %s wei, need %s wei (amount: %s + gas: %s)",
			balance.String(), totalNeeded.String(), amountWei.String(), totalGasCost.String())
	}

	// 构造 EIP-1559 动态费用交易
	txData := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        &toAddr,
		Value:     amountWei,
		Data:      nil, // 普通 ETH 转账不需要 data
	}
	tx := types.NewTx(txData)

	// 签名交易
	signer := types.NewLondonSigner(chainID)
	signedTx, err := types.SignTx(tx, signer, privKey)
	if err != nil {
		log.Fatalf("failed to sign transaction: %v", err)
	}

	// 发送交易
	if err := client.SendTransaction(ctx, signedTx); err != nil {
		log.Fatalf("failed to send transaction: %v", err)
	}

	// 输出交易信息
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("ETH Transfer Transaction Sent\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("From          : %s\n", fromAddr.Hex())
	fmt.Printf("To            : %s\n", toAddr.Hex())
	fmt.Printf("Amount        : %s ETH (%s Wei)\n", amountStr, amountWei.String())
	fmt.Printf("Gas Limit     : %d\n", gasLimit)
	fmt.Printf("Gas Tip Cap   : %s Wei\n", gasTipCap.String())
	fmt.Printf("Gas Fee Cap   : %s Wei\n", gasFeeCap.String())
	fmt.Printf("Estimated Cost: %s Wei (~%.9f ETH)\n", totalGasCost.String(), float64(totalGasCost.Uint64())/1e18)
	fmt.Printf("Nonce         : %d\n", nonce)
	fmt.Printf("Chain ID      : %s\n", chainID.String())
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("\n✅ Transaction Hash:\n")
	fmt.Printf("   %s\n", signedTx.Hash().Hex())
	fmt.Printf("\n")
	fmt.Printf("Check on Sepolia Explorer:\n")
	fmt.Printf("   https://sepolia.etherscan.io/tx/%s\n", signedTx.Hash().Hex())
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

// trim0x 移除十六进制字符串前缀 "0x"
func trim0x(s string) string {
	if len(s) >= 2 && s[0:2] == "0x" {
		return s[2:]
	}
	return s
}
