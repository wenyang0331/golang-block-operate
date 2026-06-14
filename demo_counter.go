// This is a guide-only file - not a real test file.
// abigen 生成的绑定代码用法完整演示

package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	counterpkg "golang-block-operate/pkg/counter"
)

// 这个文件展示了使用 abigen 生成绑定的完整流程
// 设置好环境变量后即可运行

// 使用步骤：
//
// 1. 部署合约（用 Remix 或通过代码）
// 2. 设置环境变量：
//    $env:ETH_RPC_URL = "https://sepolia.infura.io/v3/YOUR_KEY"
//    $env:SENDER_PRIVATE_KEY = "你的私钥"
//    $env:CONTRACT_ADDRESS = "0x部署的合约地址"
// 3. go run examples/usage.go

func demoCounter() {
	// ============================================================
	// 第1步：连接到 Sepolia 测试网
	// ============================================================
	rpcURL := os.Getenv("ETH_RPC_URL")
	if rpcURL == "" {
		log.Fatal("请设置 ETH_RPC_URL")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer client.Close()

	fmt.Println("✅ 已连接到 Sepolia 测试网")

	// ============================================================
	// 第2步：绑定合约（abigen 生成！）
	// ============================================================
	contractAddr := os.Getenv("CONTRACT_ADDRESS")
	if contractAddr == "" {
		log.Fatal("请设置 CONTRACT_ADDRESS")
	}

	counter, err := counterpkg.NewCounter(common.HexToAddress(contractAddr), client)
	if err != nil {
		log.Fatalf("绑定合约失败: %v", err)
	}

	fmt.Printf("✅ 已绑定合约: %s\n\n", contractAddr)

	// ============================================================
	// 第3步：只读查询 - GetCount
	// ============================================================
	fmt.Println("━━━━━━━━━━ 查询计数器值 ━━━━━━━━━━")
	count, err := counter.GetCount(&bind.CallOpts{Context: ctx})
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}
	fmt.Printf("当前值: %s\n\n", count.String())

	// ============================================================
	// 第4步：发送交易 - Increment
	// ============================================================
	privateKey := os.Getenv("SENDER_PRIVATE_KEY")
	if privateKey == "" {
		fmt.Println("⚠️  未设置私钥，跳过写操作测试")
		return
	}

	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatalf("私钥解析失败: %v", err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("获取 chainID 失败: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privKey, chainID)
	if err != nil {
		log.Fatalf("创建签名器失败: %v", err)
	}
	auth.Context = ctx

	// --- Increment ---
	fmt.Println("━━━━━━━━━━ increment ⏫ ━━━━━━━━━━")
	tx, err := counter.Increment(auth)
	if err != nil {
		log.Fatalf("increment 失败: %v", err)
	}
	fmt.Printf("交易已发送!\n")
	fmt.Printf("交易哈希: %s\n", tx.Hash().Hex())
	fmt.Printf("查看详情: https://sepolia.etherscan.io/tx/%s\n\n", tx.Hash().Hex())

	// --- 再次查询 ---
	fmt.Println("━━━━━━━━━━ 再次查询 ⏳ ━━━━━━━━━━")
	// 等待交易确认（简化为等待3秒）
	fmt.Println("等待交易确认...")
	time.Sleep(3 * time.Second)

	count, err = counter.GetCount(&bind.CallOpts{Context: ctx})
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}
	fmt.Printf("当前值: %s (应该是 1)\n\n", count.String())

	// --- SetCount ---
	fmt.Println("━━━━━━━━━━ setCount(42) 🔢 ━━━━━━━━━━")
	tx, err = counter.SetCount(auth, big.NewInt(42))
	if err != nil {
		log.Fatalf("setCount 失败: %v", err)
	}
	fmt.Printf("交易已发送!\n")
	fmt.Printf("交易哈希: %s\n", tx.Hash().Hex())

	time.Sleep(3 * time.Second)

	count, _ = counter.GetCount(&bind.CallOpts{Context: ctx})
	fmt.Printf("当前值: %s (应该是 42)\n\n", count.String())

	// --- Reset ---
	fmt.Println("━━━━━━━━━━ reset 🔄 ━━━━━━━━━━")
	tx, err = counter.Reset(auth)
	if err != nil {
		log.Fatalf("reset 失败: %v", err)
	}
	fmt.Printf("交易已发送!\n")
	fmt.Printf("交易哈希: %s\n", tx.Hash().Hex())

	time.Sleep(3 * time.Second)

	count, _ = counter.GetCount(&bind.CallOpts{Context: ctx})
	fmt.Printf("当前值: %s (应该是 0)\n\n", count.String())

	fmt.Println("所有测试完成！")
}
