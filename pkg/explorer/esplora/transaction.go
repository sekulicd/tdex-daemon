package esplora

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tdex-network/tdex-daemon/pkg/explorer"
	"github.com/tdex-network/tdex-daemon/pkg/httputil"
)

func (e *esplora) GetTransactionHex(hash string) (string, error) {
	url := fmt.Sprintf(
		"%s/tx/%s/hex",
		e.apiURL,
		hash,
	)
	status, resp, err := httputil.NewHTTPRequest("GET", url, "", nil)
	if err != nil {
		return "", err
	}
	if status != http.StatusOK {
		return "", fmt.Errorf(resp)
	}

	return resp, nil
}

func (e *esplora) IsTransactionConfirmed(txID string) (bool, error) {
	trxStatus, err := e.GetTransactionStatus(txID)
	if err != nil {
		return false, err
	}

	var isConfirmed bool
	switch confirmed := trxStatus["confirmed"].(type) {
	case bool:
		isConfirmed = confirmed
	}

	return isConfirmed, nil
}

func (e *esplora) GetTransactionStatus(txID string) (map[string]interface{}, error) {
	url := fmt.Sprintf(
		"%s/tx/%s/status",
		e.apiURL,
		txID,
	)
	status, resp, err := httputil.NewHTTPRequest("GET", url, "", nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, err
	}

	var trxStatus map[string]interface{}
	err = json.Unmarshal([]byte(resp), &trxStatus)
	if err != nil {
		return nil, err
	}

	return trxStatus, nil
}

func (e *esplora) GetTransactionsForAddress(address string, _ []byte) ([]explorer.Transaction, error) {
	url := fmt.Sprintf("%s/address/%s/txs", e.apiURL, address)
	status, resp, err := httputil.NewHTTPRequest("GET", url, "", nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf(resp)
	}

	return parseTransactions(resp)
}

func (e *esplora) BroadcastTransaction(txHex string) (string, error) {
	url := fmt.Sprintf("%s/tx", e.apiURL)
	headers := map[string]string{
		"Content-Type": "text/plain",
	}

	status, resp, err := httputil.NewHTTPRequest(
		"POST",
		url,
		txHex,
		headers,
	)
	if err != nil {
		return "", err
	}
	if status != http.StatusOK {
		return "", fmt.Errorf("electrs: %s", resp)
	}

	return resp, nil
}

func (e *esplora) Faucet(address string) (string, error) {
	url := fmt.Sprintf("%s/faucet", e.apiURL)
	payload := map[string]string{"address": address}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	respBody := map[string]string{}
	err = json.Unmarshal(data, &respBody)
	if err != nil {
		return "", err
	}

	return respBody["txId"], nil
}

func (e *esplora) Mint(address string, amount int) (string, string, error) {
	url := fmt.Sprintf("%s/mint", e.apiURL)
	payload := map[string]interface{}{"address": address, "quantity": amount}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", "", err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	respBody := map[string]interface{}{}
	err = json.Unmarshal(data, &respBody)
	if err != nil {
		return "", "", err
	}

	return respBody["txId"].(string), respBody["asset"].(string), nil
}

func parseTransactions(txList string) ([]explorer.Transaction, error) {
	txInterfaces := make([]interface{}, 0)
	if err := json.Unmarshal([]byte(txList), &txInterfaces); err != nil {
		return nil, err
	}
	txs := make([]explorer.Transaction, 0, len(txInterfaces))

	for _, txi := range txInterfaces {
		t, err := json.Marshal(txi)
		if err != nil {
			return nil, err
		}
		trx, err := NewTxFromJSON(string(t))
		if err != nil {
			return nil, err
		}
		txs = append(txs, trx)
	}
	return txs, nil
}
