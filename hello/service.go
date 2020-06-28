package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/yalo/grpc-connect/hello/calculator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func call(request Request) (*calculator.Response, error) {
	req := new(calculator.Request)

	// get operators
	for k, v := range request.QueryStringParameters {
		if k == "operator_one" {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, err
			}

			req.OperatorOne = f
		}

		if k == "operator_two" {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, err
			}

			req.OperatorTwo = f
		}
	}

	connection, err := setup()
	if err != nil {
		return nil, err
	}
	defer connection.Close()

	client := calculator.NewOperationsClient(connection)

	return client.Addition(context.Background(), req)
}

func setup() (*grpc.ClientConn, error) {
	// Setup environments for the service
	port := os.Getenv("PORT")
	url := os.Getenv("URL")

	log.Printf("ENV: url:%s, port:%s\n", url, port)

	service := fmt.Sprintf("%s:%s", url, port)

	return setupConnection(service)
}

func setupConnection(service string) (*grpc.ClientConn, error) {
	secure := os.Getenv("TLS")
	tls, err := strconv.ParseBool(secure)
	if err != nil {
		return nil, err
	}

	opts := make([]grpc.DialOption, 0)
	if tls {
		dir, _ := os.Getwd()
		file := filepath.Join(dir, "bin", "fullchain.pem")

		sslDomain := os.Getenv("SSL_DOMAIN")

		creds, err := credentials.NewClientTLSFromFile(file, sslDomain)
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	return grpc.Dial(service, opts...)
}
