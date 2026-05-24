param(
  [string]$FabricSamplesPath = ""
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

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..\..")
if ([string]::IsNullOrWhiteSpace($FabricSamplesPath)) {
  $FabricSamplesPath = Join-Path (Split-Path $repoRoot -Parent) "fabric-samples"
}

$testNetwork = Join-Path $FabricSamplesPath "test-network"
if (-not (Test-Path (Join-Path $testNetwork "network.sh"))) {
  throw "network.sh was not found under '$testNetwork'. Check FabricSamplesPath."
}

$bash = Resolve-GitBash
$testNetworkBash = Convert-ToGitBashPath $testNetwork

& $bash -lc "cd '$testNetworkBash' && ./network.sh down"
Write-Host "Fabric test-network stopped."
