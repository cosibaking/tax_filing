# XYGo Admin 本地开发启动脚本
# 用法: .\start.ps1 [-Init] [-Migrate] [-Restart] [-Stop] [-BackendOnly] [-FrontendOnly] [-Help]

param(
    [switch]$Init,
    [switch]$Migrate,
    [switch]$Restart,
    [switch]$Stop,
    [switch]$BackendOnly,
    [switch]$FrontendOnly,
    [switch]$Help
)

$ErrorActionPreference = 'Stop'
$Root = $PSScriptRoot
$ServerDir = Join-Path $Root 'server'
$WebDir = Join-Path $Root 'web'
$BackendPort = 4096
$FrontendPort = 5173

function Show-Help {
    Write-Host 'XYGo Admin 本地开发启动脚本' -ForegroundColor Cyan
    Write-Host ''
    Write-Host '用法:'
    Write-Host '  .\start.ps1              启动后端 + 前端（各开独立窗口）'
    Write-Host '  .\start.ps1 -Init        从零初始化（配置、数据库、依赖、迁移）'
    Write-Host '  .\scripts\init.ps1 -Start  同上，完成后自动启动'
    Write-Host '  .\start.ps1 -Migrate     启动前先执行数据库迁移'
    Write-Host '  .\start.ps1 -Restart     重启服务（先停后启）'
    Write-Host '  .\start.ps1 -Stop        停止服务（按端口 4096/5173）'
    Write-Host '  .\start.ps1 -BackendOnly 仅操作后端 (http://localhost:4096)'
    Write-Host '  .\start.ps1 -FrontendOnly 仅操作前端 (http://localhost:5173)'
    Write-Host ''
    Write-Host '示例:'
    Write-Host '  .\start.ps1 -Restart -BackendOnly   仅重启后端'
    Write-Host '  .\start.ps1 -Stop                     停止全部'
    Write-Host ''
    Write-Host '环境要求: Go 1.24+, gf CLI, Node.js 20.19+, pnpm 8.8+, MySQL 8, Redis 7'
    Write-Host '详见 README.zh-CN.md'
}

function Test-CommandExists {
    param([string]$Name)
    return [bool](Get-Command $Name -ErrorAction SilentlyContinue)
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

function Get-BackendRunCommand {
    $paths = Get-GoToolPaths
    $gfExe = Join-Path $paths.GoBin 'gf.exe'
    if (Test-Path $gfExe) {
        return "& '$gfExe' run main.go"
    }

    if (Test-CommandExists 'go') {
        Write-Host '提示: 未找到 gf，使用 go run main.go（无热重载）。安装 gf: go install github.com/gogf/gf/cmd/gf/v2@latest' -ForegroundColor Yellow
        return 'go run main.go'
    }

    throw '未找到 Go 环境。请安装 Go 并执行: go install github.com/gogf/gf/cmd/gf/v2@latest'
}

function Stop-PortListener {
    param([int]$Port)

    $pids = @()
    $pattern = ":\s*$Port\s+"

    $netstat = netstat -ano -p tcp 2>$null
    if ($netstat) {
        foreach ($line in $netstat) {
            if ($line -match 'LISTENING' -and $line -match $pattern) {
                $targetPid = ($line -split '\s+')[-1]
                if ($targetPid -match '^\d+$') { $pids += [int]$targetPid }
            }
        }
    }

    if ($pids.Count -eq 0) {
        $conns = Get-NetTCPConnection -LocalPort $Port -State Listen -ErrorAction SilentlyContinue
        if ($conns) {
            $pids = $conns | Select-Object -ExpandProperty OwningProcess -Unique
        }
    }

    foreach ($procId in ($pids | Select-Object -Unique)) {
        try {
            $proc = Get-Process -Id $procId -ErrorAction Stop
            Stop-Process -Id $procId -Force -ErrorAction Stop
            Write-Host "已停止进程 $($proc.ProcessName) (PID $procId, 端口 $Port)" -ForegroundColor Yellow
        } catch {
            # 进程可能已退出，忽略
        }
    }
}

function Test-PortListening {
    param([int]$Port)
    return [bool](Get-NetTCPConnection -LocalPort $Port -State Listen -ErrorAction SilentlyContinue)
}

function Wait-PortListening {
    param(
        [int]$Port,
        [string]$Label,
        [int]$TimeoutSec = 30
    )
    for ($i = 0; $i -lt $TimeoutSec; $i++) {
        if (Test-PortListening -Port $Port) {
            return $true
        }
        Start-Sleep -Seconds 1
    }
    Write-Host "$Label 未在 ${TimeoutSec}s 内监听端口 $Port，请查看弹出的 PowerShell 窗口中的报错。" -ForegroundColor Red
    return $false
}

function Start-EncodedPowerShell {
    param([string]$Script)
    $bytes = [System.Text.Encoding]::Unicode.GetBytes($Script)
    $encoded = [Convert]::ToBase64String($bytes)
    Start-Process powershell -ArgumentList '-NoExit', '-EncodedCommand', $encoded | Out-Null
}

function Stop-DevServices {
    param(
        [bool]$Backend = $true,
        [bool]$Frontend = $true
    )

    if ($Backend) {
        Write-Host "停止后端 (端口 $BackendPort)..." -ForegroundColor Yellow
        Stop-PortListener -Port $BackendPort
    }
    if ($Frontend) {
        Write-Host "停止前端 (端口 $FrontendPort)..." -ForegroundColor Yellow
        Stop-PortListener -Port $FrontendPort
        foreach ($legacyPort in @(3006, 3007)) {
            Stop-PortListener -Port $legacyPort
        }
    }
}

function Ensure-Config {
    $serverConfig = Join-Path $ServerDir 'manifest\config\config.yaml'
    $serverExample = Join-Path $ServerDir 'manifest\config\config.yaml.example'
    if (-not (Test-Path $serverConfig)) {
        if (-not (Test-Path $serverExample)) {
            throw "缺少后端配置模板: $serverExample"
        }
        Copy-Item $serverExample $serverConfig
        Write-Host '配置: 已生成 server/manifest/config/config.yaml，请修改数据库与 Redis 连接' -ForegroundColor Yellow
    }

    $webEnvLocal = Join-Path $WebDir '.env.local'
    $webEnvDev = Join-Path $WebDir '.env.development'
    if (-not (Test-Path $webEnvLocal)) {
        if (-not (Test-Path $webEnvDev)) {
            throw "缺少前端环境变量模板: $webEnvDev"
        }
        Copy-Item $webEnvDev $webEnvLocal
        Write-Host '配置: 已生成 web/.env.local' -ForegroundColor Yellow
    }
}

function Invoke-Init {
    $initScript = Join-Path $Root 'scripts\init.ps1'
    if (-not (Test-Path $initScript)) {
        throw "缺少初始化脚本: $initScript"
    }
    & $initScript
}

function Invoke-Migrate {
    Write-Host '迁移: 执行数据库迁移...' -ForegroundColor Yellow
    Push-Location $ServerDir
    try {
        go run tools.go migrate up
        if ($LASTEXITCODE -ne 0) { throw '数据库迁移失败' }
    } finally {
        Pop-Location
    }
    Write-Host '迁移: 完成' -ForegroundColor Green
}

function Start-Backend {
    $runCmd = Get-BackendRunCommand
    $paths = Get-GoToolPaths
    $pathPrefix = (@($paths.GoBin, $paths.GoRootBin) | Where-Object { $_ }) -join ';'
    Write-Host '后端: 启动中 -> http://localhost:4096' -ForegroundColor Green
    $script = @"
if ('$pathPrefix') { `$env:Path = '$pathPrefix;' + `$env:Path }
Set-Location -LiteralPath '$ServerDir'
$runCmd
"@
    Start-EncodedPowerShell -Script $script
}

function Start-Frontend {
    if (-not (Test-CommandExists 'pnpm')) {
        throw '未找到 pnpm，请先安装: npm install -g pnpm'
    }
    if (-not (Test-Path (Join-Path $WebDir 'node_modules'))) {
        throw '前端依赖未安装，请先运行: .\start.ps1 -Init'
    }
    Write-Host "前端: 启动中 -> http://localhost:$FrontendPort" -ForegroundColor Green
    $script = @"
Set-Location -LiteralPath '$WebDir'
pnpm dev --port $FrontendPort --strictPort
"@
    Start-EncodedPowerShell -Script $script
}

function Show-StartupSummary {
    param(
        [bool]$Backend,
        [bool]$Frontend,
        [bool]$IsRestart
    )

    Write-Host ''
    Write-Host '========== 访问地址 ==========' -ForegroundColor Cyan
    if ($Backend) {
        Write-Host "  后端 API    http://localhost:$BackendPort" -ForegroundColor Green
    }
    if ($Frontend) {
        Write-Host "  用户界面    http://localhost:$FrontendPort/" -ForegroundColor Green
        Write-Host "  管理界面    http://localhost:$FrontendPort/admin" -ForegroundColor Green
    }
    Write-Host ''

    if ($IsRestart) {
        Write-Host '已重启服务，关闭对应窗口即可停止。' -ForegroundColor Cyan
    } elseif ($Backend -and $Frontend) {
        Write-Host '已在新窗口启动后端与前端，关闭对应窗口即可停止服务。' -ForegroundColor Cyan
    } elseif ($Frontend) {
        Write-Host '已在新窗口启动前端，关闭对应窗口即可停止服务。' -ForegroundColor Cyan
    } else {
        Write-Host '已在新窗口启动后端，关闭对应窗口即可停止服务。' -ForegroundColor Cyan
    }
}

if ($Help) {
    Show-Help
    exit 0
}

if ($Init) {
    Invoke-Init
    exit 0
}

$stopBackend = -not $FrontendOnly
$stopFrontend = -not $BackendOnly

if ($Stop) {
    Stop-DevServices -Backend $stopBackend -Frontend $stopFrontend
    Write-Host '服务已停止。' -ForegroundColor Green
    exit 0
}

Ensure-Config

if ($Migrate) {
    Invoke-Migrate
}

$startBackend = $stopBackend
$startFrontend = $stopFrontend

if ($Restart) {
    Write-Host '========== 重启服务 ==========' -ForegroundColor Cyan
    Stop-DevServices -Backend $stopBackend -Frontend $stopFrontend
    Start-Sleep -Seconds 2
}

if ($startBackend) { Start-Backend }
if ($startFrontend) { Start-Frontend }

if ($startBackend -or $startFrontend) {
    Start-Sleep -Seconds 2
    $backendOk = -not $startBackend -or (Wait-PortListening -Port $BackendPort -Label '后端' -TimeoutSec 60)
    $frontendOk = -not $startFrontend -or (Wait-PortListening -Port $FrontendPort -Label '前端')
    Show-StartupSummary -Backend $startBackend -Frontend $startFrontend -IsRestart:$Restart
    if (-not $backendOk -or -not $frontendOk) {
        Write-Host '部分服务未成功启动。请检查新弹出的 PowerShell 窗口中的错误信息。' -ForegroundColor Red
        exit 1
    }
}
