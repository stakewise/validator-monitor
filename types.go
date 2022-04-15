package main

type ValidatorID struct {
	Data struct {
		Validators []struct {
			ID string `json:"id"`
		} `json:"validators"`
	} `json:"data"`
}

type ValidatorIndex struct {
	Data []struct {
		Index string `json:"index"`
	} `json:"data"`
}

type ValidatorBalance struct {
	Data []struct {
		Index   string `json:"index"`
		Balance string `json:"balance"`
	} `json:"data"`
}

type FinalityCheckpoints struct {
	Data struct {
		PreviousJustified struct {
			Epoch string `json:"epoch"`
			Root  string `json:"root"`
		} `json:"previous_justified"`
		CurrentJustified struct {
			Epoch string `json:"epoch"`
			Root  string `json:"root"`
		} `json:"current_justified"`
		Finalized struct {
			Epoch string `json:"epoch"`
			Root  string `json:"root"`
		} `json:"finalized"`
	} `json:"data"`
}

type WalletBalances struct {
	Index   int
	Balance int
}
