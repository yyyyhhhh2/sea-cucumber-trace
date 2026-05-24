# SeaTrace Blockchain Architecture

## Goal

This project should use Hyperledger Fabric as an immutable anchor layer, not as the primary business database.

Responsibilities are split as follows:

- `backend + SQLite`: source of truth for business data, users, batches, events, and application queries
- `Fabric chaincode`: immutable anchor for event hash, batch number, stage, transaction metadata
- `frontend`: reads business timeline from backend and displays chain anchor status

This is the correct shape for your current graduation-project style system. It keeps implementation cost under control while still giving you a real on-chain verification path.

## Recommended Topology

- Channel: `mychannel`
- Chaincode: `traceability`
- Organizations:
  - `OrgFarmMSP`: breeding org
  - `OrgProcessMSP`: processing org
  - `OrgLogisticsMSP`: logistics org
  - `OrgRetailMSP`: retail org
  - optional `OrgAuditMSP`: regulator / verifier
- Peers:
  - one peer per org for demo
  - one ordering service cluster or Fabric test-network orderer for local development

For a local demo, you can start with:

- 2 orgs minimum: one business org + one verifier org
- 1 channel
- 1 chaincode

For a thesis/demo environment, 4 orgs is more persuasive because it matches your trace stages.

## Data Boundary

Do not put full event details on chain.

Put on chain:

- `batchNo`
- `eventId`
- `stage`
- `dataHash`
- `occurredAt`
- `orgName`
- `txId`
- `txTimestamp`

Keep off chain:

- full `detailJson`
- operator PII if sensitive
- image/video evidence files
- user accounts and passwords
- large query-oriented business data

Reason:

- Fabric world state is not your analytics database
- chain should anchor integrity, not replace the application DB
- privacy and query performance are both better this way

## Write Path

1. Frontend submits event creation request to backend.
2. Backend writes the event into SQLite.
3. Backend computes `dataHash`.
4. Backend calls Fabric gateway `SubmitTransaction("PutTrace", payload)`.
5. Chaincode stores the immutable anchor by composite key:
   - `TRACE~batchNo~eventId`
6. Backend stores the returned `txId` and chain status in `chain_records`.

If Fabric write fails:

- event stays in SQLite
- `chain_records.status = failed`
- admin can retry anchoring later

This is already aligned with your current backend service behavior.

## Read / Verify Path

Public timeline:

1. Backend reads batch + events from SQLite.
2. Backend joins `chain_records`.
3. Frontend shows full business timeline and on-chain tx summary.

Public verification:

1. Backend recalculates hash from event data in SQLite.
2. Backend compares with stored `trace_events.data_hash`.
3. Optional stronger verification:
   - backend queries chaincode `GetTrace(batchNo, eventId)`
   - compares chaincode `dataHash` with local hash

Your current backend already does the local hash verification part. The next upgrade is adding real chain query in `gateway.go`.

## Chaincode Contract

Implemented in `chaincode/traceability/traceability.go`.

Current chaincode capabilities:

- `PutTrace(payload)`
- `GetTrace(batchNo, eventID)`
- `GetTraceByEventID(batchNo, eventID)`
- `TraceExists(batchNo, eventID)`
- `ListBatchTraces(batchNo)`
- `VerifyTraceHash(batchNo, eventID, expectedHash)`

Design choices:

- uses Fabric composite key instead of plain string key
- rejects duplicate writes for the same `batchNo + eventId`
- stores transaction metadata into the chain record itself

## Backend Integration Design

File: `backend/internal/fabric/gateway.go`

Recommended production implementation:

- use Fabric Gateway Go SDK
- connect with:
  - MSP ID
  - user cert
  - user private key
  - peer endpoint
  - TLS cert
- call:
  - `network.GetContract(cfg.FabricCC)`
  - `contract.SubmitTransaction("PutTrace", payload)`
  - optional `contract.EvaluateTransaction("GetTrace", batchNo, eventID)`

Suggested backend interface evolution:

- `RecordTrace(ctx, req)`
- `GetTrace(ctx, batchNo, eventID)`
- `VerifyTrace(ctx, batchNo, eventID, expectedHash)`

## Suggested Environment Variables

Your current config already supports:

- `FABRIC_ENABLED`
- `FABRIC_MSP_ID`
- `FABRIC_CERT_PATH`
- `FABRIC_KEY_PATH`
- `FABRIC_TLS_PATH`
- `FABRIC_PEER_ENDPOINT`
- `FABRIC_CHANNEL`
- `FABRIC_CHAINCODE`

Recommended additions when you wire the real gateway:

- `FABRIC_PEER_HOST_ALIAS`
- `FABRIC_GATEWAY_PEER`
- `FABRIC_IDENTITY_LABEL`

## Deployment Recommendation

### Local Development

- backend uses SQLite
- Fabric uses test-network
- chaincode installed from `chaincode/traceability`
- `FABRIC_ENABLED=false` by default
- switch to `true` only when test-network and certificates are ready

### Demo / Thesis Presentation

- run backend and frontend locally
- run Fabric test-network or prebuilt docker environment
- preload 2 demo batches and multiple event anchors
- show:
  - batch timeline
  - tx id
  - re-hash verification
  - chaincode query result

## What Is Still Missing

The project is now structurally ready at the chaincode level, but not fully wired end-to-end yet.

Remaining work:

1. implement real Fabric gateway calls in `backend/internal/fabric/gateway.go`
2. optionally add backend query methods that read chain data back
3. add deployment scripts for install / approve / commit chaincode
4. add retry job for failed chain anchors

## Minimal Next Step

If you want the fastest path to a working demo, do this next:

1. keep current backend DB flow unchanged
2. wire real `SubmitTransaction("PutTrace", payload)` in `gateway.go`
3. add one backend API that calls chaincode `GetTrace`
4. show the returned on-chain payload in the frontend trace page

That will give you a complete, defensible blockchain architecture with minimal rewrite.
