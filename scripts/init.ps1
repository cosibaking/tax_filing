# XYGo local dev bootstrap
# Usage:
#   .\scripts\init.ps1
#   .\scripts\init.ps1 -RootPassword 'your-root-password'
#   .\scripts\init.ps1 -RootPassword 'your-root-password' -Start
#   .\scripts\init.ps1 -SkipDatabase -SkipMigrate
#
# Override mysql bin/port in scripts/dev.local.json (see dev.local.json.example)

param(
    [string]$RootPassword = '',
    [switch]$SkipDatabase,
    [switch]$SkipMigrate,
    [switch]$ForceImport,
    [switch]$FreshConfig,
    [switch]$Start
)

$ErrorActionPreference = 'Stop'
$ScriptsDir = $PSScriptRoot
$Root = Split-Path $ScriptsDir -Parent
$ServerDir = Join-Path $Root 'server'
$WebDir = Join-Path $Root 'web'
$StartScript = Join-Path $Root 'start.ps1'

. (Join-Path $ScriptsDir 'lib\dev-env.ps1')

Write-Host ''
Write-Host '========== XYGo Bootstrap ==========' -ForegroundColor Cyan
Write-Host ''

$defaults = Get-DevDefaults -ScriptsDir $ScriptsDir
$hostName = $defaults.mysql.host
$port = Get-MySqlPort -Defaults $defaults
$database = $defaults.mysql.database
$appUser = $defaults.mysql.appUser
$appPassword = $defaults.mysql.appPassword

Write-DevStep 'Check prerequisites'
Test-DevPrerequisites
Write-DevOk 'Go / Node / pnpm ready'

Write-DevStep 'Prepare config files'
Ensure-ProjectConfig -Root $Root -Fresh:$FreshConfig
$configPath = Join-Path $ServerDir 'manifest\config\config.yaml'

if (-not $SkipDatabase) {
    $mysqlExe = Find-MySqlClient -Defaults $defaults
    if (-not $mysqlExe) {
        throw @"
mysql client not found. Either:
  1. Add mysql bin directory to PATH
  2. Copy scripts/dev.local.json.example to scripts/dev.local.json and set mysql.bin / mysql.port
"@
    }
    Write-DevOk "MySQL client: $mysqlExe (port $port)"

    if (-not $RootPassword -and $defaults.mysql.rootPassword) {
        $RootPassword = $defaults.mysql.rootPassword
    }
    if (-not $RootPassword) {
        $secure = Read-Host 'MySQL root password' -AsSecureString
        $ptr = [Runtime.InteropServices.Marshal]::SecureStringToBSTR($secure)
        try {
            $RootPassword = [Runtime.InteropServices.Marshal]::PtrToStringBSTR($ptr)
        } finally {
            [Runtime.InteropServices.Marshal]::ZeroFreeBSTR($ptr)
        }
    }

    Invoke-DatabaseBootstrap -Root $Root -MysqlExe $mysqlExe -HostName $hostName -Port $port `
        -RootPassword $RootPassword -Database $database -AppUser $appUser -AppPassword $appPassword `
        -ForceImport:$ForceImport
} else {
    Write-DevWarn 'Skipped database setup (-SkipDatabase)'
}

$configUpdated = $FreshConfig -or -not $SkipDatabase
Update-ServerDevConfig -ConfigPath $configPath -Defaults $defaults -HostName $hostName -Port $port -Database $database -Force:$configUpdated

Ensure-GoFrameCli | Out-Null
Install-BackendDeps -ServerDir $ServerDir
Install-FrontendDeps -WebDir $WebDir

if (-not $SkipMigrate) {
    Invoke-ProjectMigrate -ServerDir $ServerDir
} else {
    Write-DevWarn 'Skipped migrations (-SkipMigrate)'
}

Write-Host ''
Write-Host '========== Bootstrap Complete ==========' -ForegroundColor Green
Write-Host ''
Write-Host 'Default accounts:' -ForegroundColor Cyan
Write-Host '  Admin   Super / 123456  -> http://localhost:5173/admin'
Write-Host '  Member  register        -> http://localhost:5173/user/register'
Write-Host ''
Write-Host 'Start dev servers:' -ForegroundColor Cyan
Write-Host '  .\start.ps1'
Write-Host ''

if ($Start) {
    Write-DevStep 'Starting dev servers'
    & $StartScript
}
