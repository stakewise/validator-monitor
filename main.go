package main

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
)

var cfg struct {
	BEACON_NODE_URL   string `env:"BEACON_NODE_URL,default=http://localhost:5052"`
	GRAPH_NODE_URL    string `env:"GRAPH_NODE_URL,default=https://api.thegraph.com/subgraphs/name/stakewise/stakewise-mainnet"`
	INDEX_CHUNCK_SIZE int    `env:"INDEX_CHUNCK_SIZE,default=500"`
	OPERATOR_WALLETS  string `env:"OPERATOR_WALLETS,default=0x102f792028a56f13d6d99ed4ec8a6125de98582a"`
	BIND_ADDRESS      string `env:"BIND_ADDRESS,default=0.0.0.0:9000"`
}

func main() {
	ctx := context.Background()
	if err := envconfig.Process(ctx, &cfg); err != nil {
		log.Error().Err(err)
	}

	ethClient = NewRestClient(cfg.BEACON_NODE_URL)
	gqlClient = NewRestClient(cfg.GRAPH_NODE_URL)

	go func() {
		for {
			ValidatorEffectiveness()
			time.Sleep(time.Second * 360)
		}
	}()

	log.Info().Msgf("Starting server: %s/metrics", cfg.BIND_ADDRESS)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(cfg.BIND_ADDRESS, nil)

}
