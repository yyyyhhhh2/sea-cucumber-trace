package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const (
	traceObjectType = "TRACE"
)

type TraceContract struct {
	contractapi.Contract
}

type TraceRecord struct {
	ObjectType  string `json:"objectType"`
	BatchNo     string `json:"batchNo"`
	EventID     uint   `json:"eventId"`
	Stage       string `json:"stage"`
	DataHash    string `json:"dataHash"`
	OccurredAt  string `json:"occurredAt"`
	OrgName     string `json:"orgName"`
	TxID        string `json:"txId"`
	TxTimestamp string `json:"txTimestamp"`
}

type BatchTraceSummary struct {
	BatchNo string        `json:"batchNo"`
	Total   int           `json:"total"`
	Items   []TraceRecord `json:"items"`
}

func (t *TraceContract) PutTrace(ctx contractapi.TransactionContextInterface, payload string) error {
	rec, key, err := parseTracePayload(ctx, payload)
	if err != nil {
		return err
	}

	exists, err := t.TraceExists(ctx, rec.BatchNo, rec.EventID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("trace already exists for batch=%s eventId=%d", rec.BatchNo, rec.EventID)
	}

	raw, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("marshal trace record: %w", err)
	}
	return ctx.GetStub().PutState(key, raw)
}

func (t *TraceContract) GetTrace(ctx contractapi.TransactionContextInterface, batchNo string, eventID string) (*TraceRecord, error) {
	eventIDUint, err := strconv.ParseUint(eventID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid event id: %w", err)
	}
	return t.GetTraceByEventID(ctx, batchNo, uint(eventIDUint))
}

func (t *TraceContract) GetTraceByEventID(ctx contractapi.TransactionContextInterface, batchNo string, eventID uint) (*TraceRecord, error) {
	key, err := traceStateKey(ctx, batchNo, eventID)
	if err != nil {
		return nil, err
	}
	raw, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("get trace state: %w", err)
	}
	if len(raw) == 0 {
		return nil, fmt.Errorf("trace not found for batch=%s eventId=%d", batchNo, eventID)
	}

	var rec TraceRecord
	if err := json.Unmarshal(raw, &rec); err != nil {
		return nil, fmt.Errorf("unmarshal trace state: %w", err)
	}
	return &rec, nil
}

func (t *TraceContract) TraceExists(ctx contractapi.TransactionContextInterface, batchNo string, eventID uint) (bool, error) {
	key, err := traceStateKey(ctx, batchNo, eventID)
	if err != nil {
		return false, err
	}
	raw, err := ctx.GetStub().GetState(key)
	if err != nil {
		return false, fmt.Errorf("get trace state: %w", err)
	}
	return len(raw) > 0, nil
}

func (t *TraceContract) ListBatchTraces(ctx contractapi.TransactionContextInterface, batchNo string) (*BatchTraceSummary, error) {
	iter, err := ctx.GetStub().GetStateByPartialCompositeKey(traceObjectType, []string{batchNo})
	if err != nil {
		return nil, fmt.Errorf("query batch traces: %w", err)
	}
	defer iter.Close()

	items := make([]TraceRecord, 0)
	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			return nil, fmt.Errorf("iterate batch traces: %w", err)
		}
		var rec TraceRecord
		if err := json.Unmarshal(kv.Value, &rec); err != nil {
			return nil, fmt.Errorf("unmarshal batch trace: %w", err)
		}
		items = append(items, rec)
	}

	return &BatchTraceSummary{
		BatchNo: batchNo,
		Total:   len(items),
		Items:   items,
	}, nil
}

func (t *TraceContract) VerifyTraceHash(ctx contractapi.TransactionContextInterface, batchNo string, eventID string, expectedHash string) (bool, error) {
	rec, err := t.GetTrace(ctx, batchNo, eventID)
	if err != nil {
		return false, err
	}
	return rec.DataHash == expectedHash, nil
}

func parseTracePayload(ctx contractapi.TransactionContextInterface, payload string) (*TraceRecord, string, error) {
	var rec TraceRecord
	if err := json.Unmarshal([]byte(payload), &rec); err != nil {
		return nil, "", fmt.Errorf("invalid json: %w", err)
	}
	if rec.BatchNo == "" {
		return nil, "", fmt.Errorf("batchNo required")
	}
	if rec.EventID == 0 {
		return nil, "", fmt.Errorf("eventId required")
	}
	if rec.Stage == "" {
		return nil, "", fmt.Errorf("stage required")
	}
	if rec.DataHash == "" {
		return nil, "", fmt.Errorf("dataHash required")
	}
	if rec.OccurredAt == "" {
		return nil, "", fmt.Errorf("occurredAt required")
	}

	txTime, err := txTimestampRFC3339(ctx)
	if err != nil {
		return nil, "", err
	}
	rec.ObjectType = traceObjectType
	rec.TxID = ctx.GetStub().GetTxID()
	rec.TxTimestamp = txTime

	key, err := traceStateKey(ctx, rec.BatchNo, rec.EventID)
	if err != nil {
		return nil, "", err
	}
	return &rec, key, nil
}

func traceStateKey(ctx contractapi.TransactionContextInterface, batchNo string, eventID uint) (string, error) {
	return ctx.GetStub().CreateCompositeKey(traceObjectType, []string{batchNo, strconv.FormatUint(uint64(eventID), 10)})
}

func txTimestampRFC3339(ctx contractapi.TransactionContextInterface) (string, error) {
	ts, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", fmt.Errorf("get tx timestamp: %w", err)
	}
	return time.Unix(ts.Seconds, int64(ts.Nanos)).UTC().Format(time.RFC3339Nano), nil
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
