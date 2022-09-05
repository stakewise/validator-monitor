package main

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
)

var (
	validatorStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "validator_effectiveness_status",
	}, []string{"wallet_id", "validator_index"})
)

func ValidatorEffectiveness() error {
	for _, wallet := range strings.Split(cfg.OPERATOR_WALLETS, ",") {
		pubkeys, err := GetPubkeys(wallet)
		if err != nil {
			log.Error().Err(err)
			return err
		}

		currentSlot, previousSlot, err := getEpoch()
		if err != nil {
			log.Error().Err(err)
			return err
		}

		log.Info().Msgf("Operator efficiency calculation. Wallet ID: %s Slots: %d-%d", wallet, previousSlot, currentSlot)

		currentBalance, err := ValidatorBalances(currentSlot, pubkeys)
		if err != nil {
			log.Error().Err(err)
			return err
		}

		previousBalance, err := ValidatorBalances(previousSlot, pubkeys)
		if err != nil {
			log.Error().Err(err)
			return err
		}

		for k, v := range currentBalance {
			if oldBalance, ok := previousBalance[k]; ok {
				if v <= (oldBalance - 5000) {
					validatorStatus.WithLabelValues(wallet, k).Set(0)
				} else {
					validatorStatus.WithLabelValues(wallet, k).Set(1)
				}
			} else {
				continue
			}
		}
	}

	return nil
}
