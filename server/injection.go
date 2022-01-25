package server

import (
	"orijinplus/app/api/handlers"
	"orijinplus/app/api/routes"
	"orijinplus/app/master"
	"orijinplus/app/services"
	"orijinplus/app/store/blockchain"
	"orijinplus/app/store/dbstore"
	"orijinplus/app/store/filestore"
	"orijinplus/config"
)

// All dependency injections will go here
func Injection(c *config.Clients) (*dbstore.DBStore, *routes.Routes) {
	dbs := dbstore.NewDBStore(c.PostgresConn)
	blk := blockchain.NewBlockchainStore(c.EthereumClient)
	fs := filestore.NewFilestore(c.AWSSession)
	m := master.NewMaster(dbs)
	s := services.NewService(dbs, blk, m)
	h := handlers.NewHandlers(s, fs)
	rt := routes.NewRoutes(h)

	return dbs, rt
}
