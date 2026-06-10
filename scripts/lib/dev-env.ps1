function Write-DevStep {
    param([string]$Message)
    Write-Host ">> $Message" -ForegroundColor Cyan
}

function Write-DevOk {
    param([string]$Message)
    Write-Host "OK $Message" -ForegroundColor Green
}

function Write-DevWarn {
    param([string]$Message)
    Write-Host "!! $Message" -ForegroundColor Yellow
}

function Test-CommandExists {
    param([string]$Name)
    return [bool](Get-Command $Name -ErrorAction SilentlyContinue)
}

function Merge-JsonSection {
    param($Target, $Source)
    if (-not $Source) { return }
    foreach ($prop in $Source.PSObject.Properties) {
        $Target | Add-Member -NotePropertyName $prop.Name -NotePropertyValue $prop.Value -Force
    }
}

function Get-DevDefaults {
    param([string]$ScriptsDir)
    $defaultsPath = Join-Path $ScriptsDir 'dev.defaults.json'
    $localPath = Join-Path $ScriptsDir 'dev.local.json'
    if (-not (Test-Path $defaultsPath)) {
        throw "Missing defaults file: $defaultsPath"
    }
    $defaults = Get-Content $defaultsPath -Raw -Encoding UTF8 | ConvertFrom-Json
    if (Test-Path $localPath) {
        $local = Get-Content $localPath -Raw -Encoding UTF8 | ConvertFrom-Json
        Merge-JsonSection -Target $defaults.mysql -Source $local.mysql
        Merge-JsonSection -Target $defaults.dev -Source $local.dev
    }
    return $defaults
}

function Get-GoToolPaths {
    if (-not (Test-CommandExists 'go')) {
        return @{ GoBin = ''; GoRootBin = '' }
    }
    return @{
        GoBin     = Join-Path (go env GOPATH) 'bin'
        GoRootBin = Join-Path (go env GOROOT) 'bin'
    }
}

function Ensure-GoFrameCli {
    $paths = Get-GoToolPaths
    $gfExe = Join-Path $paths.GoBin 'gf.exe'
    if (Test-Path $gfExe) {
        return $gfExe
    }
    if (-not (Test-CommandExists 'go')) {
        throw 'Go not found. Install Go 1.24+ from https://go.dev/dl/'
    }
    Write-DevStep 'Install GoFrame CLI (gf)'
    & go install github.com/gogf/gf/cmd/gf/v2@latest
    if (-not (Test-Path $gfExe)) {
        throw 'gf install failed. Run: go install github.com/gogf/gf/cmd/gf/v2@latest'
    }
    Write-DevOk 'gf ready'
    return $gfExe
}

function Find-MySqlClient {
    param($Defaults)
    if ($Defaults.mysql.bin) {
        $candidate = Join-Path $Defaults.mysql.bin 'mysql.exe'
        if (Test-Path $candidate) { return $candidate }
    }
    $fromPath = Get-Command mysql -ErrorAction SilentlyContinue
    if ($fromPath) { return $fromPath.Source }
    foreach ($dir in $Defaults.mysql.searchBinPaths) {
        $candidate = Join-Path $dir 'mysql.exe'
        if (Test-Path $candidate) { return $candidate }
    }
    return $null
}

function Get-MySqlPort {
    param($Defaults)
    if ($Defaults.mysql.port) {
        return [int]$Defaults.mysql.port
    }
    foreach ($iniPath in $Defaults.mysql.myIniPaths) {
        if (-not (Test-Path $iniPath)) { continue }
        $match = Select-String -Path $iniPath -Pattern '^\s*port\s*=\s*(\d+)\s*$' -AllMatches |
            ForEach-Object { $_.Matches } |
            Select-Object -Last 1
        if ($match) { return [int]$match.Groups[1].Value }
    }
    return 3306
}

function Invoke-MySqlCli {
    param(
        [string]$MysqlExe,
        [string]$HostName,
        [int]$Port,
        [string]$User,
        [string]$Password,
        [string]$Sql,
        [string]$Database = ''
    )
    $args = @('-h', $HostName, '-P', $Port, '-u', $User, "-p$Password")
    if ($Database) { $args += $Database }
    $args += @('-e', $Sql)
    & $MysqlExe @args 2>&1 | ForEach-Object {
        if ($_ -match 'Using a password on the command line interface can be insecure') { return }
        if ($_ -is [System.Management.Automation.ErrorRecord]) { throw $_.Exception.Message }
        Write-Host $_
    }
    if ($LASTEXITCODE -ne 0) {
        throw "MySQL command failed (user=$User, db=$Database)"
    }
}

function Test-MySqlConnection {
    param(
        [string]$MysqlExe,
        [string]$HostName,
        [int]$Port,
        [string]$User,
        [string]$Password
    )
    try {
        Invoke-MySqlCli -MysqlExe $MysqlExe -HostName $HostName -Port $Port -User $User -Password $Password -Sql 'SELECT 1 AS ok;'
        return $true
    } catch {
        return $false
    }
}

function Test-DatabaseInitialized {
    param(
        [string]$MysqlExe,
        [string]$HostName,
        [int]$Port,
        [string]$User,
        [string]$Password,
        [string]$Database
    )
    try {
        $sql = "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = '$Database' AND table_name = 'xy_admin_user'"
        $output = & $MysqlExe -h $HostName -P $Port -u $User "-p$Password" -N -e $sql 2>$null
        return ($output -match '^\s*1\s*$')
    } catch {
        return $false
    }
}

function Ensure-ProjectConfig {
    param(
        [string]$Root,
        [switch]$Fresh
    )
    $serverDir = Join-Path $Root 'server'
    $webDir = Join-Path $Root 'web'
    $serverConfig = Join-Path $serverDir 'manifest\config\config.yaml'
    $serverExample = Join-Path $serverDir 'manifest\config\config.yaml.example'
    $webEnvLocal = Join-Path $webDir '.env.local'
    $webEnvDev = Join-Path $webDir '.env.development'

    if ($Fresh -and (Test-Path $serverConfig)) {
        Remove-Item $serverConfig -Force
    }
    if (-not (Test-Path $serverConfig)) {
        if (-not (Test-Path $serverExample)) {
            throw "Missing config template: $serverExample"
        }
        Copy-Item $serverExample $serverConfig
        Write-DevOk 'Created server/manifest/config/config.yaml'
    }

    if (-not (Test-Path $webEnvLocal)) {
        if (-not (Test-Path $webEnvDev)) {
            throw "Missing web env template: $webEnvDev"
        }
        Copy-Item $webEnvDev $webEnvLocal
        Write-DevOk 'Created web/.env.local'
    }
}

function Test-ConfigNeedsDevUpdate {
    param([string]$ConfigPath)
    if (-not (Test-Path $ConfigPath)) { return $true }
    $content = Get-Content $ConfigPath -Raw -Encoding UTF8
    $linkMatch = [regex]::Match($content, '(?m)^\s*link:\s*"([^"]*)"')
    $activeLink = if ($linkMatch.Success) { $linkMatch.Groups[1].Value } else { '' }
    return (
        $activeLink -match 'your-password' -or
        $activeLink -match '^mysql:root:' -or
        $content -match '(?m)^\s*adapter:\s*"redis"' -or
        $content -match '(?m)^\s*driver:\s*"redis"'
    )
}

function Update-ServerDevConfig {
    param(
        [string]$ConfigPath,
        $Defaults,
        [string]$HostName,
        [int]$Port,
        [string]$Database,
        [switch]$Force
    )
    if (-not (Test-Path $ConfigPath)) { return }
    if (-not $Force -and -not (Test-ConfigNeedsDevUpdate -ConfigPath $ConfigPath)) {
        Write-DevWarn 'Keep existing server config (use -FreshConfig to regenerate)'
        return
    }
    $content = Get-Content $ConfigPath -Raw -Encoding UTF8
    $appUser = $Defaults.mysql.appUser
    $appPassword = $Defaults.mysql.appPassword
    $cacheAdapter = $Defaults.dev.cacheAdapter
    $queueDriver = $Defaults.dev.queueDriver
    $dbLink = "mysql:${appUser}:${appPassword}@tcp(${HostName}:${Port})/${Database}?loc=Local&parseTime=true"

    $content = [regex]::Replace($content, '(?m)^(\s*adapter:\s*")[^"]*(")', ('${1}' + $cacheAdapter + '${2}'), 1)
    $content = [regex]::Replace($content, '(?m)^(\s*driver:\s*")[^"]*(")', ('${1}' + $queueDriver + '${2}'), 1)
    $content = [regex]::Replace($content, '(?m)^(\s*link:\s*")[^"]*(")', ('${1}' + $dbLink + '${2}'), 1)
    Set-Content -Path $ConfigPath -Value $content -Encoding UTF8 -NoNewline
    Write-DevOk "Database link: ${HostName}:${Port}/${Database}"
}

function Invoke-DatabaseBootstrap {
    param(
        [string]$Root,
        [string]$MysqlExe,
        [string]$HostName,
        [int]$Port,
        [string]$RootPassword,
        [string]$Database,
        [string]$AppUser,
        [string]$AppPassword,
        [switch]$ForceImport
    )
    $sqlFile = Join-Path $Root 'mysql_install.sql'
    if (-not (Test-Path $sqlFile)) {
        throw "Missing SQL file: $sqlFile"
    }

    Write-DevStep "Connect MySQL ${HostName}:${Port}"
    if (-not (Test-MySqlConnection -MysqlExe $MysqlExe -HostName $HostName -Port $Port -User root -Password $RootPassword)) {
        throw 'Cannot connect as root. Check service, port and password. Override in scripts/dev.local.json'
    }

    $initialized = Test-DatabaseInitialized -MysqlExe $MysqlExe -HostName $HostName -Port $Port -User root -Password $RootPassword -Database $Database
    if ($initialized -and -not $ForceImport) {
        Write-DevWarn "Database $Database already initialized, skip import (use -ForceImport to re-import)"
    } else {
        Write-DevStep "Create database $Database and import mysql_install.sql"
        Invoke-MySqlCli -MysqlExe $MysqlExe -HostName $HostName -Port $Port -User root -Password $RootPassword -Sql "CREATE DATABASE IF NOT EXISTS ``$Database`` DEFAULT CHARSET utf8mb4;"
        $importCmd = "`"$MysqlExe`" -h $HostName -P $Port -u root -p$RootPassword $Database < `"$sqlFile`""
        cmd /c $importCmd
        if ($LASTEXITCODE -ne 0) {
            throw 'SQL import failed'
        }
        Write-DevOk 'SQL import done'
    }

    Write-DevStep "Create app user $AppUser"
    $userSql = @(
        "CREATE USER IF NOT EXISTS '$AppUser'@'localhost' IDENTIFIED BY '$AppPassword';"
        "CREATE USER IF NOT EXISTS '$AppUser'@'127.0.0.1' IDENTIFIED BY '$AppPassword';"
        "ALTER USER '$AppUser'@'localhost' IDENTIFIED BY '$AppPassword';"
        "ALTER USER '$AppUser'@'127.0.0.1' IDENTIFIED BY '$AppPassword';"
        "GRANT ALL PRIVILEGES ON ``$Database``.* TO '$AppUser'@'localhost';"
        "GRANT ALL PRIVILEGES ON ``$Database``.* TO '$AppUser'@'127.0.0.1';"
        "FLUSH PRIVILEGES;"
    ) -join "`n"
    Invoke-MySqlCli -MysqlExe $MysqlExe -HostName $HostName -Port $Port -User root -Password $RootPassword -Sql $userSql
    Write-DevOk "User $AppUser granted on $Database"
}

function Test-DevPrerequisites {
    $missing = @()
    if (-not (Test-CommandExists 'go')) { $missing += 'Go 1.24+ (https://go.dev/dl/)' }
    if (-not (Test-CommandExists 'node')) { $missing += 'Node.js 20.19+ (https://nodejs.org/)' }
    if (-not (Test-CommandExists 'pnpm')) { $missing += 'pnpm 8.8+ (npm install -g pnpm)' }
    if ($missing.Count -gt 0) {
        throw "Missing dependencies:`n - $($missing -join "`n - ")"
    }
}

function Install-FrontendDeps {
    param([string]$WebDir)
    if (Test-Path (Join-Path $WebDir 'node_modules')) {
        Write-DevWarn 'node_modules exists, skip pnpm install'
        return
    }
    Write-DevStep 'Install frontend deps (pnpm install)'
    Push-Location $WebDir
    try {
        pnpm install
        if ($LASTEXITCODE -ne 0) { throw 'pnpm install failed' }
    } finally {
        Pop-Location
    }
    Write-DevOk 'Frontend deps ready'
}

function Install-BackendDeps {
    param([string]$ServerDir)
    Write-DevStep 'Download backend deps (go mod download)'
    Push-Location $ServerDir
    try {
        go mod download
        if ($LASTEXITCODE -ne 0) { throw 'go mod download failed' }
    } finally {
        Pop-Location
    }
    Write-DevOk 'Backend deps ready'
}

function Invoke-ProjectMigrate {
    param([string]$ServerDir)
    Write-DevStep 'Run migrations (go run tools.go migrate up)'
    Push-Location $ServerDir
    try {
        go run tools.go migrate up
        if ($LASTEXITCODE -ne 0) { throw 'Database migration failed' }
    } finally {
        Pop-Location
    }
    Write-DevOk 'Migrations done'
}
