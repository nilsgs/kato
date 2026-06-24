#Requires -Version 5.1
$ErrorActionPreference = 'Stop'

$InstallDir = Join-Path $env:USERPROFILE '.kato\bin'
$RepoDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$Version = (Get-Content (Join-Path $RepoDir 'VERSION') -Raw).Trim()
$Commit = try { (git -C $RepoDir rev-parse --short HEAD 2>$null) } catch { 'unknown' }
if (-not $Commit) { $Commit = 'unknown' }
$Ldflags = "-s -w -X kato/cmd.version=$Version -X kato/cmd.commit=$Commit"

Write-Host "Building kg v${Version}+${Commit}..."
Write-Host "Installing to $InstallDir..."
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}
$Dest = Join-Path $InstallDir 'kg.exe'
Push-Location (Join-Path $RepoDir 'src')
try {
    & go build -ldflags $Ldflags -o $Dest .
    if ($LASTEXITCODE -ne 0) { throw 'Build failed' }
} finally {
    Pop-Location
}

$UserPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if ($UserPath -split ';' | Where-Object { $_ -eq $InstallDir }) {
    Write-Host "PATH already contains $InstallDir"
} else {
    $NewPath = $InstallDir + ';' + $UserPath
    [Environment]::SetEnvironmentVariable('Path', $NewPath, 'User')
    $env:PATH = $InstallDir + ';' + $env:PATH
    Write-Host "Added $InstallDir to user PATH"
}

Write-Host 'Done. Restart your terminal, then run: kg --help'
