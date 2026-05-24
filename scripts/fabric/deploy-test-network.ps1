param(
  [string]$FabricSamplesPath = "",
  [string]$Channel = "mychannel",
  [string]$ChaincodeName = "traceability",
  [string]$ChaincodeVersion = "1.0",
  [int]$ChaincodeSequence = 1,
  [switch]$DownloadSamples
)

$ErrorActionPreference = "Stop"

function Resolve-GitBash {
  $bash = Get-Command bash -ErrorAction SilentlyContinue
  if ($bash) {
    return $bash.Source
  }

  $gitBash = "C:\Program Files\Git\bin\bash.exe"
  if (Test-Path $gitBash) {
    return $gitBash
  }

  throw "Git Bash was not found. Install Git for Windows or add bash.exe to PATH."
}

function Convert-ToGitBashPath([string]$Path) {
  $full = [System.IO.Path]::GetFullPath($Path)
  $drive = $full.Substring(0, 1).ToLowerInvariant()
  $rest = $full.Substring(2).Replace("\", "/")
  return "/$drive$rest"
}

function Assert-Command([string]$Name, [string]$InstallHint) {
  if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
    throw "$Name was not found. $InstallHint"
  }
}

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..\..")
$chaincodePath = Join-Path $repoRoot "chaincode\traceability"

Assert-Command "docker" "Install Docker Desktop, enable WSL 2 backend, and restart this terminal."
Assert-Command "git" "Install Git for Windows."

$dockerVersion = docker --version
Write-Host "Using $dockerVersion"

try {
  docker info | Out-Null
} catch {
  throw "Docker is installed but the daemon is not running. Start Docker Desktop first."
}

if ([string]::IsNullOrWhiteSpace($FabricSamplesPath)) {
  $FabricSamplesPath = Join-Path (Split-Path $repoRoot -Parent) "fabric-samples"
}

if (-not (Test-Path $FabricSamplesPath)) {
  if (-not $DownloadSamples) {
    throw "fabric-samples was not found at '$FabricSamplesPath'. Re-run with -DownloadSamples or pass -FabricSamplesPath."
  }
  git clone https://github.com/hyperledger/fabric-samples.git $FabricSamplesPath
}

$testNetwork = Join-Path $FabricSamplesPath "test-network"
if (-not (Test-Path (Join-Path $testNetwork "network.sh"))) {
  throw "network.sh was not found under '$testNetwork'. Check FabricSamplesPath."
}

$bash = Resolve-GitBash
$testNetworkBash = Convert-ToGitBashPath $testNetwork
$chaincodeBash = Convert-ToGitBashPath $chaincodePath

Write-Host "Stopping any previous Fabric test-network..."
& $bash -lc "cd '$testNetworkBash' && ./network.sh down"

Write-Host "Starting Fabric test-network and creating channel '$Channel'..."
& $bash -lc "cd '$testNetworkBash' && ./network.sh up createChannel -ca -c '$Channel'"

Write-Host "Deploying chaincode '$ChaincodeName' from '$chaincodePath'..."
& $bash -lc "cd '$testNetworkBash' && ./network.sh deployCC -c '$Channel' -ccn '$ChaincodeName' -ccp '$chaincodeBash' -ccl go -ccv '$ChaincodeVersion' -ccs '$ChaincodeSequence'"

$org1 = Join-Path $testNetwork "organizations\peerOrganizations\org1.example.com"
$userMsp = Join-Path $org1 "users\User1@org1.example.com\msp"
$cert = Get-ChildItem (Join-Path $userMsp "signcerts") -Filter "*.pem" | Select-Object -First 1
$key = Get-ChildItem (Join-Path $userMsp "keystore") | Select-Object -First 1
$tls = Join-Path $org1 "peers\peer0.org1.example.com\tls\ca.crt"
$envFile = Join-Path $repoRoot "backend\.env.fabric"

@"
FABRIC_ENABLED=true
FABRIC_MSP_ID=Org1MSP
FABRIC_CHANNEL=$Channel
FABRIC_CHAINCODE=$ChaincodeName
FABRIC_PEER_ENDPOINT=localhost:7051
FABRIC_CERT_PATH=$($cert.FullName)
FABRIC_KEY_PATH=$($key.FullName)
FABRIC_TLS_PATH=$tls
"@ | Set-Content -Encoding UTF8 $envFile

Write-Host ""
Write-Host "Fabric test-network is running."
Write-Host "Chaincode: $ChaincodeName on channel $Channel"
Write-Host "Backend Fabric env written to: $envFile"
Write-Host ""
Write-Host "Important: backend/internal/fabric/gateway.go still needs real Fabric Gateway calls before FABRIC_ENABLED=true can be used by the Go API."
