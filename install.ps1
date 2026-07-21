#Requires -Version 5.1
$ErrorActionPreference = 'Stop'

$InstallDir = Join-Path $env:USERPROFILE '.kato\bin'
$RepoDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$Version = (Get-Content (Join-Path $RepoDir 'VERSION') -Raw).Trim()
$Commit = try { (git -C $RepoDir rev-parse --short HEAD 2>$null) } catch { 'unknown' }
if (-not $Commit) { $Commit = 'unknown' }
$Ldflags = "-s -w -X kato/cmd.version=$Version -X kato/cmd.commit=$Commit"

Write-Host "Building kato v${Version}+${Commit}..."
Write-Host "Installing to $InstallDir..."
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}
$Dest = Join-Path $InstallDir 'kato.exe'
$LegacyDest = Join-Path $InstallDir 'kg.exe'
Push-Location (Join-Path $RepoDir 'src')
try {
    & go build -ldflags $Ldflags -o $Dest .
    if ($LASTEXITCODE -ne 0) { throw 'Build failed' }
} finally {
    Pop-Location
}
if (Test-Path $LegacyDest) {
    Remove-Item -LiteralPath $LegacyDest -Force
    Write-Host "Removed obsolete $LegacyDest"
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

Write-Host 'Done. Restart your terminal, then run: kato --help'

$KatoDir = Join-Path $env:USERPROFILE '.kato'
$InitPs1 = Join-Path $KatoDir 'init.ps1'

# Copy the kato-owned PowerShell init script from the repo to ~/.kato/.
if (-not (Test-Path $KatoDir)) {
    New-Item -ItemType Directory -Path $KatoDir -Force | Out-Null
}
$ScriptSrc = Join-Path $RepoDir 'src\internal\shell\scripts\init.ps1'
Copy-Item -Path $ScriptSrc -Destination $InitPs1 -Force
Write-Host "Wrote $InitPs1"

# Inject (or replace) the kato block in the PowerShell profile using BEGIN/END markers.
$ProfileDir = Split-Path -Parent $PROFILE
if (-not (Test-Path $ProfileDir)) {
    New-Item -ItemType Directory -Path $ProfileDir -Force | Out-Null
}
if (-not (Test-Path $PROFILE)) {
    New-Item -ItemType File -Path $PROFILE -Force | Out-Null
}
$Block = "# BEGIN kato`n. `"$InitPs1`"`n# END kato"
$ProfileContent = Get-Content $PROFILE -Raw -ErrorAction SilentlyContinue
if ($ProfileContent -match '# BEGIN kato') {
    $ProfileContent = $ProfileContent -replace '(?s)# BEGIN kato.*?# END kato', $Block
    Set-Content -Path $PROFILE -Value $ProfileContent -Encoding UTF8 -NoNewline
    Write-Host "Updated kato block in $PROFILE"
} else {
    Add-Content -Path $PROFILE -Value "`n$Block" -Encoding UTF8
    Write-Host "Added kato block to $PROFILE"
}
