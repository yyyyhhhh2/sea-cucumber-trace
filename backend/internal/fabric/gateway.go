package fabric

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

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
	if err := g.validateConfig(); err != nil {
		return nil, err
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal fabric payload: %w", err)
	}

	conn, err := newGrpcConnection(g.cfg)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	id, err := newIdentity(g.cfg)
	if err != nil {
		return nil, err
	}
	sign, err := newSign(g.cfg)
	if err != nil {
		return nil, err
	}

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(conn),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(15*time.Second),
		client.WithCommitStatusTimeout(30*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("connect fabric gateway: %w", err)
	}
	defer gw.Close()

	network := gw.GetNetwork(g.cfg.FabricChannel)
	contract := network.GetContract(g.cfg.FabricCC)

	_, commit, err := contract.SubmitAsyncWithContext(ctx, "PutTrace", client.WithArguments(string(payload)))
	if err != nil {
		return nil, fmt.Errorf("submit PutTrace: %w", err)
	}

	status, err := commit.StatusWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("wait for commit status tx=%s: %w", commit.TransactionID(), err)
	}
	if !status.Successful {
		return nil, fmt.Errorf("transaction %s failed to commit with code %d", status.TransactionID, int32(status.Code))
	}

	return &RecordResult{
		TxID:        status.TransactionID,
		BlockNumber: status.BlockNumber,
	}, nil
}

func (g *Gateway) validateConfig() error {
	missing := make([]string, 0)
	if g.cfg.FabricMSP == "" {
		missing = append(missing, "FABRIC_MSP_ID")
	}
	if g.cfg.FabricCertPath == "" {
		missing = append(missing, "FABRIC_CERT_PATH")
	}
	if g.cfg.FabricKeyPath == "" {
		missing = append(missing, "FABRIC_KEY_PATH")
	}
	if g.cfg.FabricTLSPath == "" {
		missing = append(missing, "FABRIC_TLS_PATH")
	}
	if g.cfg.FabricPeer == "" {
		missing = append(missing, "FABRIC_PEER_ENDPOINT")
	}
	if g.cfg.FabricChannel == "" {
		missing = append(missing, "FABRIC_CHANNEL")
	}
	if g.cfg.FabricCC == "" {
		missing = append(missing, "FABRIC_CHAINCODE")
	}
	if len(missing) > 0 {
		return fmt.Errorf("fabric config missing: %v", missing)
	}
	return nil
}

func newGrpcConnection(cfg *config.Config) (*grpc.ClientConn, error) {
	tlsCertificatePEM, err := os.ReadFile(cfg.FabricTLSPath)
	if err != nil {
		return nil, fmt.Errorf("read fabric TLS cert: %w", err)
	}

	tlsCertificate, err := identity.CertificateFromPEM(tlsCertificatePEM)
	if err != nil {
		return nil, fmt.Errorf("parse fabric TLS cert: %w", err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(tlsCertificate)

	transportCredentials := credentials.NewClientTLSFromCert(certPool, cfg.FabricPeerHost)
	return grpc.NewClient("dns:///"+cfg.FabricPeer, grpc.WithTransportCredentials(transportCredentials))
}

func newIdentity(cfg *config.Config) (*identity.X509Identity, error) {
	certificatePEM, err := os.ReadFile(cfg.FabricCertPath)
	if err != nil {
		return nil, fmt.Errorf("read fabric identity cert: %w", err)
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		return nil, fmt.Errorf("parse fabric identity cert: %w", err)
	}

	id, err := identity.NewX509Identity(cfg.FabricMSP, certificate)
	if err != nil {
		return nil, fmt.Errorf("create fabric identity: %w", err)
	}
	return id, nil
}

func newSign(cfg *config.Config) (identity.Sign, error) {
	privateKeyPEM, err := os.ReadFile(cfg.FabricKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read fabric private key: %w", err)
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("parse fabric private key: %w", err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		return nil, fmt.Errorf("create fabric signer: %w", err)
	}
	return sign, nil
}
