# MySQL import / app user setup (also available via init.ps1)
# Usage:
#   .\scripts\import-mysql.ps1 -RootPassword 'your-root-password'
#   .\scripts\import-mysql.ps1 -RootPassword 'your-root-password' -Import
#   .\scripts\import-mysql.ps1 -RootPassword 'your-root-password' -Import -ForceImport

param(
    [Parameter(Mandatory = $true)]
    [string]$RootPassword,
    [string]$MySqlBin = '',
    [int]$Port = 0,
    [string]$HostName = '',
    [string]$Database = '',
    [string]$AppUser = '',
    [string]$AppPassword = '',
    [switch]$Import,
    [switch]$ForceImport
)

$ErrorActionPreference = 'Stop'
$ScriptsDir = $PSScriptRoot
$Root = Split-Path $ScriptsDir -Parent

. (Join-Path $ScriptsDir 'lib\dev-env.ps1')

$defaults = Get-DevDefaults -ScriptsDir $ScriptsDir
if ($MySqlBin) { $defaults.mysql | Add-Member -NotePropertyName bin -NotePropertyValue $MySqlBin -Force }
if ($Port -gt 0) { $defaults.mysql | Add-Member -NotePropertyName port -NotePropertyValue $Port -Force }
if ($HostName) { $defaults.mysql | Add-Member -NotePropertyName host -NotePropertyValue $HostName -Force }
if ($Database) { $defaults.mysql | Add-Member -NotePropertyName database -NotePropertyValue $Database -Force }
if ($AppUser) { $defaults.mysql | Add-Member -NotePropertyName appUser -NotePropertyValue $AppUser -Force }
if ($AppPassword) { $defaults.mysql | Add-Member -NotePropertyName appPassword -NotePropertyValue $AppPassword -Force }

$mysqlExe = Find-MySqlClient -Defaults $defaults
if (-not $mysqlExe) {
    throw 'mysql.exe not found. Set scripts/dev.local.json or add mysql to PATH'
}

$hostName = $defaults.mysql.host
$port = Get-MySqlPort -Defaults $defaults
$database = $defaults.mysql.database
$appUser = $defaults.mysql.appUser
$appPassword = $defaults.mysql.appPassword

if ($Import) {
    Invoke-DatabaseBootstrap -Root $Root -MysqlExe $mysqlExe -HostName $hostName -Port $port `
        -RootPassword $RootPassword -Database $database -AppUser $appUser -AppPassword $appPassword `
        -ForceImport:$ForceImport
} else {
    Write-DevStep "Connect MySQL ${hostName}:${port}"
    if (-not (Test-MySqlConnection -MysqlExe $mysqlExe -HostName $hostName -Port $port -User root -Password $RootPassword)) {
        throw 'MySQL root connection failed'
    }
    $userSql = @(
        "CREATE USER IF NOT EXISTS '$appUser'@'localhost' IDENTIFIED BY '$appPassword';"
        "CREATE USER IF NOT EXISTS '$appUser'@'127.0.0.1' IDENTIFIED BY '$appPassword';"
        "ALTER USER '$appUser'@'localhost' IDENTIFIED BY '$appPassword';"
        "ALTER USER '$appUser'@'127.0.0.1' IDENTIFIED BY '$appPassword';"
        "GRANT ALL PRIVILEGES ON ``$database``.* TO '$appUser'@'localhost';"
        "GRANT ALL PRIVILEGES ON ``$database``.* TO '$appUser'@'127.0.0.1';"
        "FLUSH PRIVILEGES;"
    ) -join "`n"
    Invoke-MySqlCli -MysqlExe $mysqlExe -HostName $hostName -Port $port -User root -Password $RootPassword -Sql $userSql
    Write-DevOk "User $appUser granted on $database"
}

$configPath = Join-Path $Root 'server\manifest\config\config.yaml'
Update-ServerDevConfig -ConfigPath $configPath -Defaults $defaults -HostName $hostName -Port $port -Database $database

Write-Host ''
Write-Host 'Done.' -ForegroundColor Green
