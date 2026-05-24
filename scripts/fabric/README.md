# Fabric Test-Network Deployment

This project uses Hyperledger Fabric as an immutable trace anchor layer. The recommended local demo network is the official `fabric-samples/test-network`.

## Prerequisites

- Docker Desktop with WSL 2 backend enabled
- Git for Windows
- Go installed and available in the shell used by Fabric scripts
- Internet access for the first run, because Fabric samples, Docker images, and Go modules may need to be downloaded

## Deploy

From the repository root:

```powershell
.\scripts\fabric\deploy-test-network.ps1 -DownloadSamples
```

The script does the following:

- locates or downloads `fabric-samples`
- starts the Fabric test network with Certificate Authorities
- creates channel `mychannel`
- deploys `chaincode/traceability` as chaincode name `traceability`
- writes backend Fabric variables to `backend/.env.fabric`

Use an existing `fabric-samples` checkout:

```powershell
.\scripts\fabric\deploy-test-network.ps1 -FabricSamplesPath F:\fabric-samples
```

## Stop

```powershell
.\scripts\fabric\stop-test-network.ps1
```

## Backend Integration Note

The test network can deploy the chaincode, but the Go backend still needs a real Fabric Gateway implementation in:

```text
backend/internal/fabric/gateway.go
```

Until that is implemented, keep:

```text
FABRIC_ENABLED=false
```

for the backend API. With `FABRIC_ENABLED=false`, the app uses the built-in mock ledger for demo transaction IDs.
