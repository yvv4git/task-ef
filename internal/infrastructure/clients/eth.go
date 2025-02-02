package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/yvv4git/task-ef/internal/config"
)

const (
	reqEthGetBlockByNumber   = "%s/api?module=proxy&action=eth_getBlockByNumber&tag=%s&boolean=true&apikey=%s"
	reqEthGetLastBlockNumber = "%s/api?module=proxy&action=eth_blockNumber&apikey=%s"
)

type Block struct {
	Difficulty       string        `json:"difficulty"`
	ExtraData        string        `json:"extraData"`
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
	Hash             string        `json:"hash"`
	LogsBloom        string        `json:"logsBloom"`
	Miner            string        `json:"miner"`
	MixHash          string        `json:"mixHash"`
	Nonce            string        `json:"nonce"`
	Number           string        `json:"number"`
	ParentHash       string        `json:"parentHash"`
	ReceiptsRoot     string        `json:"receiptsRoot"`
	Sha3Uncles       string        `json:"sha3Uncles"`
	Size             string        `json:"size"`
	StateRoot        string        `json:"stateRoot"`
	Timestamp        string        `json:"timestamp"`
	Transactions     []Transaction `json:"transactions"`
	TransactionsRoot string        `json:"transactionsRoot"`
	Uncles           []string      `json:"uncles"`
}

type Transaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	Type             string `json:"type"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

type ClientETH struct {
	cfg        *config.Config
	httpClient *http.Client
}

func NewClient(cfg *config.Config) *ClientETH {
	return &ClientETH{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: time.Duration(cfg.ExternalAPI.Timeout) * time.Second},
	}
}

func (c *ClientETH) FetchBlockByNumber(blockNumber uint64) (*Block, error) {
	hexBlockNumber := fmt.Sprintf("0x%x", blockNumber)
	url := fmt.Sprintf(reqEthGetBlockByNumber, c.cfg.ExternalAPI.URL, hexBlockNumber, c.cfg.ExternalAPI.Key)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	var response struct {
		Result  *Block `json:"result"`
		Jsonrpc string `json:"jsonrpc"`
		ID      int    `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if response.Result == nil {
		return nil, fmt.Errorf("block not found or other error")
	}

	return response.Result, nil
}

func (c *ClientETH) FetchLastBlockNumber() (uint64, error) {
	url := fmt.Sprintf(reqEthGetLastBlockNumber, c.cfg.ExternalAPI.URL, c.cfg.ExternalAPI.Key)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	var response struct {
		Result string `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("decode response: %w", err)
	}

	blockNumber, err := strconv.ParseUint(response.Result[2:], 16, 64)
	if err != nil {
		return 0, fmt.Errorf("parse block number: %w", err)
	}

	return blockNumber, nil
}
