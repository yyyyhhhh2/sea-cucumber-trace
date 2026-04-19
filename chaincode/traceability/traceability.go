package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type TraceContract struct {
	contractapi.Contract
}

type TraceRecord struct {
	BatchNo    string `json:"batchNo"`
	EventID    uint   `json:"eventId"`
	Stage      string `json:"stage"`
	DataHash   string `json:"dataHash"`
	OccurredAt string `json:"occurredAt"`
	OrgName    string `json:"orgName"`
}

// PutTrace stores an immutable trace anchor keyed by batchNo|eventId.
func (t *TraceContract) PutTrace(ctx contractapi.TransactionContextInterface, payload string) error {
	var rec TraceRecord
	if err := json.Unmarshal([]byte(payload), &rec); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}
	key := fmt.Sprintf("TRACE_%s_%d", rec.BatchNo, rec.EventID)
	return ctx.GetStub().PutState(key, []byte(payload))
}

// GetTrace returns stored JSON for a batch/event key.
func (t *TraceContract) GetTrace(ctx contractapi.TransactionContextInterface, batchNo string, eventID string) (string, error) {
	key := fmt.Sprintf("TRACE_%s_%s", batchNo, eventID)
	b, err := ctx.GetStub().GetState(key)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func main() {
	cc, err := contractapi.NewChaincode(&TraceContract{})
	if err != nil {
		panic(err)
	}
	if err := cc.Start(); err != nil {
		panic(err)
	}
}
