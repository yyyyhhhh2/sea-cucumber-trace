package fabric

import (
	"context"
	"encoding/json"
	"fmt"

	"sea-cucumber-trace/backend/internal/config"
)

type Gateway struct {
	cfg *config.Config
}

func NewGateway(cfg *config.Config) *Gateway {
	if !cfg.FabricEnabled {
		return nil
	}
	return &Gateway{cfg: cfg}
}

func ResolveLedger(enabled bool, gateway *Gateway) Ledger {
	if enabled && gateway != nil {
		return gateway
	}
	return NewMockLedger()
}

func (g *Gateway) RecordTrace(ctx context.Context, req RecordRequest) (*RecordResult, error) {
	if g == nil {
		return nil, fmt.Errorf("fabric gateway not configured")
	}
	args, _ := json.Marshal(req)
	_ = g.cfg.FabricChannel
	return nil, fmt.Errorf("fabric-gateway not wired: implement SubmitTransaction PutTrace in gateway.go (payload=%s)", string(args))
}
