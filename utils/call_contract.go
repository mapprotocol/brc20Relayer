package utils

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
)

const DefaultGasLimit = 4500000

var (
	Test, _ = abi.JSON(strings.NewReader(TestABIJSON))
)

func Dail(url string) *ethclient.Client {
	conn, err := ethclient.Dial(url)
	if err != nil {
		log.Crit("dial failed", "url", url)
	}
	return conn
}

func PackInput(a abi.ABI, abiMethod string, params ...interface{}) ([]byte, error) {
	input, err := a.Pack(abiMethod, params...)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func SendTransaction(client *ethclient.Client, from, to common.Address, value *big.Int, privateKey *ecdsa.PrivateKey, input []byte, gasLimitSetting uint64) (common.Hash, error) {
	// Ensure a valid value field and resolve the account nonce
	nonce, err := client.PendingNonceAt(context.Background(), from)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get pending noce, error: %s", err.Error())
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get suggest gas price, error: %s", err.Error())
	}
	gasLimit := uint64(DefaultGasLimit)

	//If the contract surely has code (or code is not needed), estimate the transaction
	msg := ethereum.CallMsg{From: from, To: &to, GasPrice: gasPrice, Value: value, Data: input}
	gasLimit, err = client.EstimateGas(context.Background(), msg)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get estimategest gas, error: %s", err.Error())
	}
	if gasLimit < 1 {
		gasLimit = 2100000
	}
	gasLimit = uint64(DefaultGasLimit)

	if gasLimitSetting != 0 {
		gasLimit = gasLimitSetting // in units
	}

	// Create the transaction, sign it and schedule it for execution
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     input,
	})

	chainID, _ := client.ChainID(context.Background())
	log.Debug("tx info", "nonce ", nonce, " gasLimit ", gasLimit, " gasPrice ", gasPrice, " chainID ", chainID)

	signer := types.LatestSignerForChainID(chainID)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to sign tx, error: %s", err.Error())
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to send transaction, error: %s", err.Error())

	}
	return signedTx.Hash(), nil
}

func QueryTx(conn *ethclient.Client, txHash common.Hash) {
	logger := log.New("func", "QueryTx")
	logger.Info("Please waiting ", " txHash ", txHash.String())
	for {
		time.Sleep(time.Millisecond * 200)
		_, isPending, err := conn.TransactionByHash(context.Background(), txHash)
		if err != nil {
			logger.Warn("failed to get transaction by hash", "error", err)
		}
		if !isPending {
			break
		}
	}

	queryTx(conn, txHash)
}

func queryTx(conn *ethclient.Client, txHash common.Hash) {
	logger := log.New("func", "queryTx")
	var (
		err     error
		receipt *types.Receipt
	)
	for {
		receipt, err = conn.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			break
		}
		logger.Warn("failed to get transaction receipt", "error", err)
		time.Sleep(time.Millisecond * 200)
	}

	if receipt.Status == types.ReceiptStatusSuccessful {
		logger.Info("Transaction Success", "height", receipt.BlockNumber.Uint64(), "hash", txHash)
	} else if receipt.Status == types.ReceiptStatusFailed {
		logger.Error("Transaction Failed", "height", receipt.BlockNumber.Uint64(), "hash", txHash)
	}
}
