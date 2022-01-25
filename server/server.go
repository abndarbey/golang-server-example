package server

import (
	"log"
	"orijinplus/config"
	"orijinplus/settings/cloud"
	"orijinplus/settings/database/postgres"
)

func StartApplication(conf config.Config) {
	postgresConn, err := postgres.ConnectPostgres(conf.Database.PSQLSource)
	if err != nil {
		log.Fatal(err)
	}
	if postgresConn == nil {
		log.Fatal("unable to connect with postgres db")
	}
	defer postgresConn.Close()

	awsSession := cloud.NewAWSSession()

	c := &config.Clients{
		PostgresConn: postgresConn,
		AWSSession:   awsSession,
	}

	restServer := NewRestServer(c)

	restServer.Start(conf.Server.Address)
}
