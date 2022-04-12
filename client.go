package main

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

var ethClient *RestClient
var gqlClient *RestClient

type RestClient struct {
	client *resty.Client
}

func NewRestClient(address string) *RestClient {
	client := resty.New()

	client.SetTimeout(1 * time.Minute)
	client.SetHostURL(address)
	client.SetHeader("Content-Type", "application/json")

	log.Info().Msgf("Initialize new client connected to: %s", address)

	return &RestClient{
		client: client,
	}
}
