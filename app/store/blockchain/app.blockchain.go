package blockchain

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

type BlockchainStore struct {
	ethClient *ethclient.Client
}

func NewBlockchainStore(ethClient *ethclient.Client) *BlockchainStore {
	return &BlockchainStore{ethClient}
}
