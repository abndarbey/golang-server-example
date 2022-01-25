package config

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Clients struct {
	PostgresConn   *pgxpool.Pool
	EthereumClient *ethclient.Client
	AWSSession     *session.Session
}
