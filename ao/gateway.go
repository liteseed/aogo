package ao

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const getTransactionQuery = `
query GetTransactions ($transactionIds: [ID!]!) {
	transactions(ids: $transactionIds) {
		edges {
			node {
				owner {
					address
				}
				tags {
					name
					value
				}
				block {
					id
					height
					timestamp
				}
			}
		}
	}
}`

func LoadTransactionMeta(id string) (map[string]any, error) {
	var result map[string]any

	body := map[string]any{
		"query":     getTransactionQuery,
		"variables": map[string]any{"transactionIds": []string{id}},
	}
	data, err := json.Marshal(body)
	if err != nil {
		return result, err
	}

	resp, err := http.Post(GATEWAY+"/graphql", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return result, err
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(res, &result)

	if err != nil {
		return result, err
	}

	return result["data"].(map[string]any)["transactions"].(map[string]any), err
}
