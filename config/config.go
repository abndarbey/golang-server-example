package config

import (
	"fmt"
	"os"
)

// Config stores all configurations of the application
type Config struct {
	Server   *Server
	Database *Database
	// AWSCredentails *AWSCredentails
}

type Database struct {
	PSQLSource string `mapstructure:"PSQL_SOURCE"`
}

type AWSCredentails struct {
	Region          string `mapstructure:"AWS_DEFAULT_REGION"`
	AccessKeyID     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AccessKeySecret string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
}

type Server struct {
	Address string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig() (*Config, error) {
	// Read server address and db

	serverAddress := os.Getenv("SERVER_ADDRESS")
	psqlSource := os.Getenv("PSQL_SOURCE")

	// awsDefaultRegion := os.Getenv("AWS_DEFAULT_REGION")
	// awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	// awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if serverAddress == "" {
		return nil, fmt.Errorf("server address is required")
	}
	if psqlSource == "" {
		return nil, fmt.Errorf("psql source is required")
	}
	// if awsDefaultRegion == "" {
	// 	return nil, fmt.Errorf("aws default region is required")
	// }
	// if awsAccessKeyID == "" {
	// 	return nil, fmt.Errorf("aws access key id is required")
	// }
	// if awsSecretAccessKey == "" {
	// 	return nil, fmt.Errorf("aws secret access key is required")
	// }

	db := &Database{
		PSQLSource: psqlSource,
	}
	server := &Server{
		Address: serverAddress,
	}
	config := &Config{
		Database: db,
		Server:   server,
	}

	return config, nil
}
