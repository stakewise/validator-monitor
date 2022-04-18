package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// GetPubkey return list of pubkeys for specified wallet
func GetPubkeys(wallet string) ([]string, error) {
	var pubkeys []string
	var lastId string
	var resp ValidatorID
	for {
		_, err := gqlClient.client.R().
			SetBody(`{"query":"{\n  validators(first:1000, where: {operator: \"` + strings.ToLower(wallet) + `\", id_gt: \"` + lastId + `\"}) {\n    id\n  }\n}","variables":{}}`).
			SetResult(&resp).
			Post("")
		if err != nil || len(resp.Data.Validators) < 1 {
			log.Error().Msgf("Can't get public keys from graph node %s", err)
			return nil, err
		}
		for _, v := range resp.Data.Validators {
			pubkeys = append(pubkeys, v.ID)
		}
		if len(resp.Data.Validators) >= 1000 {
			lastId = pubkeys[len(pubkeys)-1]
		} else {
			break
		}
	}

	chunkSize := 100
	var indexes []string

	for i := 0; i < len(pubkeys); i += chunkSize {
		end := i + chunkSize
		if end > len(pubkeys) {
			end = len(pubkeys)
		}

		addr := fmt.Sprintf("/eth/v1/beacon/states/head/validators?id=%s", strings.Join(pubkeys[i:end], ","))
		var resp ValidatorIndex
		_, err := ethClient.client.R().
			SetResult(&resp).
			Get(addr)
		if err != nil {
			log.Error().Msgf("Can't convert public keys to indexes %s", err)
			return nil, err
		}

		for i := 0; i < len(resp.Data); i++ {
			if resp.Data[i].Status == "active_ongoing" {
				indexes = append(indexes, resp.Data[i].Index)
			}
		}
	}

	return indexes, nil
}

// ValidatorBalances provides the validator balances for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validatorPubkeys is a list of validator pubkeys to restrict the returned values.  If no validators are supplied no filter
// will be applied.
func ValidatorBalances(stateID int, validatorPubkeys []string) (map[string]int, error) {
	if len(validatorPubkeys) <= 1 {
		log.Error().Msg("No active validators")
		return nil, nil
	}
	if len(validatorPubkeys) > cfg.INDEX_CHUNCK_SIZE {
		return chunkedValidatorBalances(stateID, validatorPubkeys)
	}

	url := fmt.Sprintf("/eth/v1/beacon/states/%d/validator_balances", stateID)

	if len(validatorPubkeys) != 0 {
		ids := make([]string, len(validatorPubkeys))
		for i := range validatorPubkeys {
			ids[i] = validatorPubkeys[i]
		}
		url = fmt.Sprintf("%s?id=%s", url, strings.Join(ids, ","))
	}

	var resp ValidatorBalance

	_, err := ethClient.client.R().
		SetResult(&resp).
		Get(url)
	if err != nil || len(resp.Data) < 1 {
		log.Error().Msgf("Can't get validator balances %e", err)
		return nil, err
	}

	res := make(map[string]int)

	for i := 0; i < len(resp.Data); i++ {
		index := resp.Data[i].Index
		balance, err := strconv.Atoi(resp.Data[i].Balance)
		if err != nil {
			log.Error().Err(err)
		}
		res[index] = balance
	}

	return res, nil

}

// chunkedValidatorBalances obtains the validator balances a chunk at a time.
func chunkedValidatorBalances(stateID int, validatorPubkeys []string) (map[string]int, error) {
	res := make(map[string]int)
	pubkeyChunkSize := cfg.INDEX_CHUNCK_SIZE
	for i := 0; i < len(validatorPubkeys); i += pubkeyChunkSize {
		chunkStart := i
		chunkEnd := i + pubkeyChunkSize
		if len(validatorPubkeys) < chunkEnd {
			chunkEnd = len(validatorPubkeys)
		}
		chunk := validatorPubkeys[chunkStart:chunkEnd]
		chunkRes, err := ValidatorBalances(stateID, chunk)
		if err != nil {
			return nil, errors.Wrap(err, "failed to obtain chunk")
		}
		for k, v := range chunkRes {
			res[k] = v
		}
	}

	return res, nil
}

// getEpoch return current and previous epoch slots
func getEpoch() (int, int, error) {
	var resp FinalityCheckpoints
	addr := fmt.Sprintf("/eth/v1/beacon/states/head/finality_checkpoints")

	_, err := ethClient.client.R().
		SetResult(&resp).
		Get(addr)

	if err != nil {
		log.Error().Msgf("Can't get epoch from beacon node %s", err)
		return 0, 0, err
	}

	currentEpoch, err := strconv.Atoi(resp.Data.Finalized.Epoch)
	if err != nil {
		log.Error().Err(err)
		return 0, 0, err
	}

	slotsPerEpoch, err := slotsPerEpoch()
	if err != nil {
		log.Error().Msgf("Can't get slots per epoch %s", err)
		return 0, 0, err
	}

	currentSlot := currentEpoch * slotsPerEpoch
	previousSlot := (currentEpoch - 1) * slotsPerEpoch

	return currentSlot, previousSlot, nil
}

// slotsPerEpoch
//
func slotsPerEpoch() (int, error) {
	var resp struct {
		Data struct {
			SlotsPerEpoch string `json:"SLOTS_PER_EPOCH"`
		} `json:"data"`
	}
	addr := fmt.Sprintf("/eth/v1/config/spec")

	_, err := ethClient.client.R().
		SetResult(&resp).
		Get(addr)

	if err != nil {
		log.Error().Msgf("Can't get epoch from beacon node %s", err)
		return 0, err
	}

	slotsPerEpoch, err := strconv.Atoi(resp.Data.SlotsPerEpoch)
	if err != nil {
		log.Error().Msgf("Can't parse SLOTS_PER_EPOCH %s", err)
		return 0, err
	}

	return slotsPerEpoch, nil
}
