package seed

import (
	"context"
	"log"
	"orijinplus/app/master"
	"orijinplus/app/services"
	"orijinplus/app/store/blockchain"
	"orijinplus/app/store/dbstore"
	"orijinplus/config"
	"orijinplus/settings/database/postgres"
	"orijinplus/utils/logger"
)

func SeedData(conf config.Config) {
	PostgresConn, err := postgres.ConnectPostgres(conf.Database.PSQLSource)
	if err != nil {
		log.Fatal(err)
	}
	if PostgresConn == nil {
		log.Fatal("unable to connect with postgres db")
	}
	defer PostgresConn.Close()

	c := &config.Clients{
		PostgresConn: PostgresConn,
		// EthereumClient: ethereumClient,
	}

	dbStore := dbstore.NewDBStore(c.PostgresConn)
	blk := blockchain.NewBlockchainStore(c.EthereumClient)
	m := master.NewMaster(dbStore)
	s := services.NewService(dbStore, blk, m)

	superadmin, userErr := InsertAdmin(s)
	if userErr != nil {
		log.Fatal(userErr.Message)
	}
	if err := InsertOrganizations(s); err != nil {
		log.Fatal(err.Message)
	}

	// Initiate db transactions
	ctx := context.Background()
	tx, txnErr := dbStore.DBTX.BeginTx(ctx)
	if txnErr != nil {
		log.Fatal(txnErr)
	}
	defer dbStore.DBTX.RollbackTx(ctx, tx)

	if err := InsertContainers(ctx, tx, dbStore, superadmin); err != nil {
		log.Fatal(err.Message)
	}
	if err := InsertPallets(ctx, tx, dbStore, superadmin); err != nil {
		log.Fatal(err.Message)
	}

	// Commit transactions to db
	if txnErr := dbStore.DBTX.CommitTx(ctx, tx); err != nil {
		log.Fatal(txnErr)
	}

	logger.Success("Database seeded successfully")
}
