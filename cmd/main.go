package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"orijinplus/config"
	"orijinplus/seed"
	"orijinplus/server"
	"orijinplus/utils/logger"
	"os"

	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
)

const (
	envPrefix string = "ORIJINPLUS"
)

func main() {
	var (
		rootFlagSet = flag.NewFlagSet("orijinplus", flag.ExitOnError)
		dbseed      = rootFlagSet.Bool("dbseed", false, "seed the database")
		test        = rootFlagSet.Bool("test", false, "test application")
		run         = rootFlagSet.Bool("run", false, "run the server")
	)

	rootCmd := &ffcli.Command{
		Name:      "root",
		Options:   []ff.Option{ff.WithEnvVarPrefix(envPrefix)},
		ShortHelp: "Run root commands.",
		FlagSet:   rootFlagSet,
		Exec: func(_ context.Context, args []string) error {
			if !*dbseed && !*test && !*run {
				return errors.New("-dbseed or -run is required but not provided ")
			}

			if *dbseed {
				seedDatabase()
			}
			if *test {
				testApplication()
			}
			if *run {
				startApplication()
			}
			return nil
		},
	}

	if err := rootCmd.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		logger.Fatal(err.Error())
	}
}

func seedDatabase() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	seed.SeedData(*conf)
}

func testApplication() {
	logger.Fatal("no test function provided")
}

func startApplication() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	server.StartApplication(*conf)
}
