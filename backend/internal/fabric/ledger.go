package fabric

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type RecordRequest struct {
	BatchNo    string `json:"batchNo"`
	EventID    uint   `json:"eventId"`
	Stage      string `json:"stage"`
	DataHash   string `json:"dataHash"`
	OccurredAt string `json:"occurredAt"`
	OrgName    string `json:"orgName"`
}

type RecordResult struct {
	TxID        string `json:"txId"`
	BlockNumber uint64 `json:"blockNumber,omitempty"`
}

type Ledger interface {
	RecordTrace(ctx context.Context, req RecordRequest) (*RecordResult, error)
}

type MockLedger struct{}

func NewMockLedger() *MockLedger { return &MockLedger{} }

func (m *MockLedger) RecordTrace(ctx context.Context, req RecordRequest) (*RecordResult, error) {
	raw, _ := json.Marshal(req)
	sum := sha256.Sum256(raw)
	txID := "mock_" + hex.EncodeToString(sum[:16])
	return &RecordResult{TxID: txID, BlockNumber: 0}, nil
}

func HashTraceEvent(batchNo, stage, detailJSON, location, operator string, occurredAt time.Time, evidenceURLs string) string {
	payload := map[string]string{
		"batchNo":      batchNo,
		"stage":        stage,
		"detailJson":   detailJSON,
		"location":     location,
		"operatorName": operator,
		"evidenceUrls": evidenceURLs,
		"occurredAt":   occurredAt.UTC().Format(time.RFC3339Nano),
	}
	b, _ := json.Marshal(payload)
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func ShortTxID(tx string) string {
	if len(tx) <= 16 {
		return tx
	}
	return fmt.Sprintf("%s…%s", tx[:8], tx[len(tx)-6:])
}
