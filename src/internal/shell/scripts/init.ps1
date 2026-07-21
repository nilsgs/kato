# kato shell integration
function kato {
  if ($args[0] -eq 'nav') {
    $dir = & kato.exe nav
    if ($dir) { Set-Location $dir }
  } else {
    & kato.exe @args
  }
}

# kato alias loader — reads ~/.kato/aliases at shell startup
$_kato_aliases = Join-Path $HOME '.kato\aliases'
if (Test-Path $_kato_aliases) {
  Get-Content $_kato_aliases | ForEach-Object {
    if ($_ -match '^\s*$' -or $_ -match '^\s*#') { return }
    $p = $_ -split '=', 2
    $n = $p[0].Trim(); $c = $p[1].Trim()
    Invoke-Expression "function global:$n { & $c @args }"
  }
}
