package chain

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	eth *ethclient.Client
}

func NewClient(rawURL string) (*Client, error) {
	eth, err := ethclient.Dial(rawURL)
	if err != nil {
		return nil, err
	}
	return &Client{eth: eth}, nil
}

func (c *Client) BlockNumber(ctx context.Context) (uint64, error) {
	return c.eth.BlockNumber(ctx)
}

func (c *Client) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return c.eth.HeaderByNumber(ctx, number)
}

func (c *Client) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	return c.eth.FilterLogs(ctx, q)
}

func (c *Client) FinalizedBlockNumber(ctx context.Context) (uint64, error) {
	header, err := c.eth.HeaderByNumber(ctx, big.NewInt(-3)) // finalized
	if err != nil {
		return 0, err
	}
	return header.Number.Uint64(), nil
}

func (c *Client) SafeBlockNumber(ctx context.Context) (uint64, error) {
	header, err := c.eth.HeaderByNumber(ctx, big.NewInt(-4)) // safe
	if err != nil {
		return 0, err
	}
	return header.Number.Uint64(), nil
}

var TransferEventTopic = common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")

func (c *Client) Close() {
	c.eth.Close()
}
