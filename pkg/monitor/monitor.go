package monitor

import (
	"celestia-lightnode-monitor/pkg/config"
	"celestia-lightnode-monitor/pkg/log"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Performance struct {
	NodeURL string
	Sync    struct {
		Synced bool
		Error  string
	}
	MinBalance struct {
		Enough bool
		Error  string
	}
}

type CelestiaRpcStatusResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		NodeInfo struct {
			ProtocolVersion struct {
				P2P   string `json:"p2p"`
				Block string `json:"block"`
				App   string `json:"app"`
			} `json:"protocol_version"`
			ID         string `json:"id"`
			ListenAddr string `json:"listen_addr"`
			Network    string `json:"network"`
			Version    string `json:"version"`
			Channels   string `json:"channels"`
			Moniker    string `json:"moniker"`
			Other      struct {
				TxIndex    string `json:"tx_index"`
				RPCAddress string `json:"rpc_address"`
			} `json:"other"`
		} `json:"node_info"`
		SyncInfo struct {
			LatestBlockHash     string    `json:"latest_block_hash"`
			LatestAppHash       string    `json:"latest_app_hash"`
			LatestBlockHeight   string    `json:"latest_block_height"`
			LatestBlockTime     time.Time `json:"latest_block_time"`
			EarliestBlockHash   string    `json:"earliest_block_hash"`
			EarliestAppHash     string    `json:"earliest_app_hash"`
			EarliestBlockHeight string    `json:"earliest_block_height"`
			EarliestBlockTime   time.Time `json:"earliest_block_time"`
			CatchingUp          bool      `json:"catching_up"`
		} `json:"sync_info"`
		ValidatorInfo struct {
			Address string `json:"address"`
			PubKey  struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"pub_key"`
			VotingPower string `json:"voting_power"`
		} `json:"validator_info"`
	} `json:"result"`
}

type BalanceResponse struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

func getLatestBlockFromRPC(standardRPC string) (int, error) {
	timeout := 5
	req, err := http.NewRequest("GET", standardRPC+"/status", nil)
	if err != nil {
		return 0, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	var response CelestiaRpcStatusResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}
	standardrpcLatestBlock, _ := strconv.Atoi(response.Result.SyncInfo.LatestBlockHeight)
	return standardrpcLatestBlock, nil
}

type getewayAPIheadResonpse struct {
	Header struct {
		Version struct {
			Block string `json:"block"`
			App   string `json:"app"`
		} `json:"version"`
		ChainID     string    `json:"chain_id"`
		Height      string    `json:"height"`
		Time        time.Time `json:"time"`
		LastBlockID struct {
			Hash  string `json:"hash"`
			Parts struct {
				Total int    `json:"total"`
				Hash  string `json:"hash"`
			} `json:"parts"`
		} `json:"last_block_id"`
		LastCommitHash     string `json:"last_commit_hash"`
		DataHash           string `json:"data_hash"`
		ValidatorsHash     string `json:"validators_hash"`
		NextValidatorsHash string `json:"next_validators_hash"`
		ConsensusHash      string `json:"consensus_hash"`
		AppHash            string `json:"app_hash"`
		LastResultsHash    string `json:"last_results_hash"`
		EvidenceHash       string `json:"evidence_hash"`
		ProposerAddress    string `json:"proposer_address"`
	} `json:"header"`
	ValidatorSet struct {
		Validators []struct {
			Address string `json:"address"`
			PubKey  struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"pub_key"`
			VotingPower      string `json:"voting_power"`
			ProposerPriority string `json:"proposer_priority"`
		} `json:"validators"`
		Proposer struct {
			Address string `json:"address"`
			PubKey  struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"pub_key"`
			VotingPower      string `json:"voting_power"`
			ProposerPriority string `json:"proposer_priority"`
		} `json:"proposer"`
	} `json:"validator_set"`
	Commit struct {
		Height  int `json:"height"`
		Round   int `json:"round"`
		BlockID struct {
			Hash  string `json:"hash"`
			Parts struct {
				Total int    `json:"total"`
				Hash  string `json:"hash"`
			} `json:"parts"`
		} `json:"block_id"`
		Signatures []struct {
			BlockIDFlag      int       `json:"block_id_flag"`
			ValidatorAddress string    `json:"validator_address"`
			Timestamp        time.Time `json:"timestamp"`
			Signature        string    `json:"signature"`
		} `json:"signatures"`
	} `json:"commit"`
	Dah struct {
		RowRoots    []string `json:"row_roots"`
		ColumnRoots []string `json:"column_roots"`
	} `json:"dah"`
}

func getLatestBlockFromGatewayAPI(gatewayAPI string) (int, error) {
	timeout := 5
	req, err := http.NewRequest("GET", gatewayAPI+"/head", nil)
	if err != nil {
		return 0, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	var response getewayAPIheadResonpse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}
	gatewayAPILatestBlock, _ := strconv.Atoi(response.Header.Height)
	return gatewayAPILatestBlock, nil
}

// checkSyncStatus checks if the gateway API is synced with the standard RPC, false means not synced, true means synced
func checkSyncStatus(gatewayAPI string, standardRPC string) (bool, error) {
	gatewayBlock, err := getLatestBlockFromGatewayAPI(gatewayAPI)
	if err != nil {
		return false, fmt.Errorf("failed to get latest block from Gateway API: %v", err)
	}

	rpcBlock, err := getLatestBlockFromRPC(standardRPC)
	if err != nil {
		return false, fmt.Errorf("failed to get latest block from RPC: %v", err)
	}
	// We consider the max behind block is 5
	if (rpcBlock - gatewayBlock) > 5 {
		log.Log.Error("Node " + gatewayAPI + " is not synced, gateway height is " + strconv.Itoa(gatewayBlock) + " standard RPC block is " + strconv.Itoa(rpcBlock))
		return false, nil
	}
	return true, nil
}

// minBalanceCheck checks if the gateway API has enough balance
func minBalanceCheck(gatewayAPI string, minBalance int) (bool, error) {
	timeout := 5
	req, err := http.NewRequest("GET", gatewayAPI+"/balance", nil)
	if err != nil {
		return false, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}
	var response BalanceResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, err
	}
	balance, _ := strconv.Atoi(response.Amount)
	if balance < minBalance {
		log.Log.Error("Node " + gatewayAPI + " does not have enough balance, balance is " + response.Amount + " minimum balance is " + strconv.Itoa(minBalance))
		return false, nil
	}
	return true, nil
}

func CheckNodes(config config.Config) ([]Performance, error) {
	var NodesPerformance []Performance
	for _, gatewayAPI := range config.Node.GatewayAPI {
		log.Log.Info("Checking node: " + gatewayAPI + "")
		Node := Performance{
			NodeURL: gatewayAPI,
		}

		// check sync status
		isSynced, err := checkSyncStatus(gatewayAPI, config.Node.StandardRPC)
		if err != nil {
			log.Log.Error("Failed to check node " + gatewayAPI + "sync status: " + err.Error())
			Node.Sync.Synced = false
			Node.Sync.Error = err.Error()
		} else {
			Node.Sync.Synced = isSynced
		}

		// check balance
		hasEnoughBalance, err := minBalanceCheck(gatewayAPI, config.Node.MinimumBalance)
		if err != nil {
			log.Log.Error("Failed to check node" + gatewayAPI + err.Error())
			Node.MinBalance.Enough = false
			Node.MinBalance.Error = err.Error()
		} else {
			Node.MinBalance.Enough = hasEnoughBalance
		}
		NodesPerformance = append(NodesPerformance, Node)

	}
	return NodesPerformance, nil
}
