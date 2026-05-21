#Requires -Version 5.1
# ===================================
# 族谱数字化管理平台 - 启动脚本 (PowerShell)
# ===================================

[CmdletBinding()]
param(
    [Parameter(Position = 0)]
    [ValidateSet("start", "stop", "restart", "status")]
    [string]$Action = "start"
)

# 服务配置
$Config = @{
    BackendPort         = 8080
    AdminPort           = 3000
    VisualizationPort   = 5174
    LogDir              = ".\logs"
    BackendProcessName  = "GenealogyMa_Backend"
    AdminProcessName    = "GenealogyMa_Admin"
    VisualizationProcessName = "GenealogyMa_Visualization"
}

# 颜色
function Get-ColorText {
    param([string]$Text, [string]$Color)
    $colors = @{
        "Red"    = "`e[91m"
        "Green"  = "`e[92m"
        "Yellow" = "`e[93m"
        "Blue"   = "`e[94m"
        "NC"     = "`e[0m"
    }
    return "$($colors[$Color])$Text$($colors['NC'])"
}

# 检查端口占用
function Test-PortAvailable {
    param([int]$Port)
    $connections = Get-NetTCPConnection -LocalPort $Port -ErrorAction SilentlyContinue
    return $null -eq $connections
}

# 启动后端服务
function Start-Backend {
    Write-Host ""
    Write-Host (Get-ColorText "[1/3]" "Blue") -NoNewline
    Write-Host " 启动后端服务 (Go)..."

    $logFile = Join-Path $Config.LogDir "backend_$(Get-Date -Format 'yyyyMMdd_HHmmss').log"

    $process = Start-Process -FilePath "go" `
        -ArgumentList "run", "cmd/api/main.go" `
        -WorkingDirectory ".\backend" `
        -PassThru `
        -RedirectStandardOutput $logFile `
        -RedirectStandardError $logFile `
        -WindowStyle Hidden `
        -ErrorAction SilentlyContinue

    if ($process) {
        $process | ForEach-Object { $_.StartInfo.WindowStyle = 'Hidden' } | Out-Null
        Write-Host (Get-ColorText "[启动]" "Green") -NoNewline
        Write-Host " 后端服务 (PID: $($process.Id))"
        Write-Host (Get-ColorText "[日志]" "Green") " $logFile"

        # 保存 PID
        $pidFile = Join-Path $Config.LogDir "backend.pid"
        $process.Id | Set-Content -Path $pidFile -Force
    }
}

# 启动管理后台
function Start-Admin {
    Write-Host ""
    Write-Host (Get-ColorText "[2/3]" "Blue") -NoNewline
    Write-Host " 启动管理后台 (React)..."

    $logFile = Join-Path $Config.LogDir "admin_$(Get-Date -Format 'yyyyMMdd_HHmmss').log"

    $process = Start-Process -FilePath "npm" `
        -ArgumentList "run", "dev" `
        -WorkingDirectory ".\frontend\admin" `
        -PassThru `
        -RedirectStandardOutput $logFile `
        -RedirectStandardError $logFile `
        -WindowStyle Hidden `
        -ErrorAction SilentlyContinue

    if ($process) {
        Write-Host (Get-ColorText "[启动]" "Green") -NoNewline
        Write-Host " 管理后台 (PID: $($process.Id))"
        Write-Host (Get-ColorText "[日志]" "Green") " $logFile"

        $pidFile = Join-Path $Config.LogDir "admin.pid"
        $process.Id | Set-Content -Path $pidFile -Force
    }
}

# 启动可视化前端
function Start-Visualization {
    Write-Host ""
    Write-Host (Get-ColorText "[3/3]" "Blue") -NoNewline
    Write-Host " 启动可视化前端 (Vue 3)..."

    $logFile = Join-Path $Config.LogDir "visualization_$(Get-Date -Format 'yyyyMMdd_HHmmss').log"

    $process = Start-Process -FilePath "npm" `
        -ArgumentList "run", "dev" `
        -WorkingDirectory ".\frontend\visualization" `
        -PassThru `
        -RedirectStandardOutput $logFile `
        -RedirectStandardError $logFile `
        -WindowStyle Hidden `
        -ErrorAction SilentlyContinue

    if ($process) {
        Write-Host (Get-ColorText "[启动]" "Green") -NoNewline
        Write-Host " 可视化前端 (PID: $($process.Id))"
        Write-Host (Get-ColorText "[日志]" "Green") " $logFile"

        $pidFile = Join-Path $Config.LogDir "visualization.pid"
        $process.Id | Set-Content -Path $pidFile -Force
    }
}

# 等待服务就绪
function Wait-ForServices {
    Write-Host ""
    Write-Host (Get-ColorText "等待服务启动..." "Blue")

    # 等待后端
    Write-Host -NoNewline "  等待后端服务 (端口 $($Config.BackendPort))... "
    for ($i = 0; $i -lt 30; $i++) {
        try {
            $response = Invoke-WebRequest -Uri "http://localhost:$($Config.BackendPort)/api/health" -UseBasicParsing -TimeoutSec 1 -ErrorAction SilentlyContinue
            if ($response.StatusCode -eq 200 -or $response.StatusCode -eq 404) {
                Write-Host (Get-ColorText "OK" "Green")
                break
            }
        } catch {
            Start-Sleep -Milliseconds 500
        }
    }

    # 等待管理后台
    Write-Host -NoNewline "  等待管理后台 (端口 $($Config.AdminPort))... "
    for ($i = 0; $i -lt 20; $i++) {
        try {
            $response = Invoke-WebRequest -Uri "http://localhost:$($Config.AdminPort)" -UseBasicParsing -TimeoutSec 1 -ErrorAction SilentlyContinue
            if ($response.StatusCode -eq 200) {
                Write-Host (Get-ColorText "OK" "Green")
                break
            }
        } catch {
            Start-Sleep -Milliseconds 500
        }
    }

    # 等待可视化前端
    Write-Host -NoNewline "  等待可视化前端 (端口 $($Config.VisualizationPort))... "
    for ($i = 0; $i -lt 20; $i++) {
        try {
            $response = Invoke-WebRequest -Uri "http://localhost:$($Config.VisualizationPort)" -UseBasicParsing -TimeoutSec 1 -ErrorAction SilentlyContinue
            if ($response.StatusCode -eq 200) {
                Write-Host (Get-ColorText "OK" "Green")
                break
            }
        } catch {
            Start-Sleep -Milliseconds 500
        }
    }
}

# 停止所有服务
function Stop-AllServices {
    Write-Host ""
    Write-Host (Get-ColorText "正在停止所有服务..." "Yellow")

    # 读取 PID 并停止
    $services = @("backend", "admin", "visualization")
    foreach ($service in $services) {
        $pidFile = Join-Path $Config.LogDir "$service.pid"
        if (Test-Path $pidFile) {
            $pid = Get-Content $pidFile -Raw
            $process = Get-Process -Id $pid -ErrorAction SilentlyContinue
            if ($process) {
                Stop-Process -Id $pid -Force -ErrorAction SilentlyContinue
                Write-Host (Get-ColorText "[停止]" "Green") " $service 服务 (PID: $pid)"
            }
            Remove-Item $pidFile -Force -ErrorAction SilentlyContinue
        }
    }

    # 清理残留进程
    Get-Process | Where-Object {
        $_.MainWindowTitle -like "GenealogyMa_*"
    } | ForEach-Object {
        Stop-Process -Id $_.Id -Force -ErrorAction SilentlyContinue
        Write-Host (Get-ColorText "[清理]" "Yellow") " 残留进程 (PID: $($_.Id))"
    }

    Write-Host ""
    Write-Host (Get-ColorText "所有服务已停止" "Green")
}

# 查看服务状态
function Get-ServiceStatus {
    Write-Host ""
    Write-Host (Get-ColorText "检查服务状态..." "Blue")

    $services = @{
        "backend"        = @{ Display = "后端服务"; Port = $Config.BackendPort }
        "admin"          = @{ Display = "管理后台"; Port = $Config.AdminPort }
        "visualization"  = @{ Display = "可视化前端"; Port = $Config.VisualizationPort }
    }

    foreach ($service in $services.GetEnumerator()) {
        $pidFile = Join-Path $Config.LogDir "$($service.Key).pid"
        $status = (Get-ColorText "[已停止]" "Red")
        $running = $false

        if (Test-Path $pidFile) {
            $pid = Get-Content $pidFile -Raw
            $process = Get-Process -Id $pid -ErrorAction SilentlyContinue
            if ($process) {
                $status = (Get-ColorText "[运行中]" "Green")
                $running = $true
            }
        }

        $portStatus = if (Test-PortAvailable $service.Value.Port) { "端口可用" } else { "端口占用" }
        Write-Host "  $($status) $($service.Value.Display) - $portStatus"
    }
}

# 显示启动信息
function Show-StartupInfo {
    Write-Host ""
    Write-Host (Get-ColorText "========================================" "Green")
    Write-Host (Get-ColorText "   服务已启动!" "Green")
    Write-Host (Get-ColorText "========================================" "Green")
    Write-Host ""
    Write-Host "  $($(Get-ColorText "管理后台:" "Blue"))  http://localhost:$($Config.AdminPort)"
    Write-Host "  $($(Get-ColorText "可视化:" "Blue"))    http://localhost:$($Config.VisualizationPort)"
    Write-Host "  $($(Get-ColorText "API文档:" "Blue"))   http://localhost:$($Config.BackendPort)/swagger/index.html"
    Write-Host ""
    Write-Host "  $($(Get-ColorText "停止服务:" "Yellow"))  .\scripts\Start-PowerShell.ps1 stop"
    Write-Host ""
}

# 主函数
function Start-AllServices {
    Write-Host ""
    Write-Host (Get-ColorText "========================================" "Blue")
    Write-Host (Get-ColorText "   族谱数字化管理平台 - 启动脚本" "Blue")
    Write-Host (Get-ColorText "========================================" "Blue")

    # 创建日志目录
    if (-not (Test-Path $Config.LogDir)) {
        New-Item -ItemType Directory -Path $Config.LogDir -Force | Out-Null
    }

    Write-Host ""
    Write-Host (Get-ColorText "[检查]" "Blue") " 检查端口占用..."
    if (-not (Test-PortAvailable $Config.BackendPort)) {
        Write-Host (Get-ColorText "[警告]" "Yellow") " 端口 $($Config.BackendPort) (后端API) 已被占用"
    } else {
        Write-Host (Get-ColorText "[OK]" "Green") " 端口 $($Config.BackendPort) (后端API) 可用"
    }
    if (-not (Test-PortAvailable $Config.AdminPort)) {
        Write-Host (Get-ColorText "[警告]" "Yellow") " 端口 $($Config.AdminPort) (管理后台) 已被占用"
    } else {
        Write-Host (Get-ColorText "[OK]" "Green") " 端口 $($Config.AdminPort) (管理后台) 可用"
    }
    if (-not (Test-PortAvailable $Config.VisualizationPort)) {
        Write-Host (Get-ColorText "[警告]" "Yellow") " 端口 $($Config.VisualizationPort) (可视化前端) 已被占用"
    } else {
        Write-Host (Get-ColorText "[OK]" "Green") " 端口 $($Config.VisualizationPort) (可视化前端) 可用"
    }

    Write-Host ""
    Write-Host (Get-ColorText "开始启动服务..." "Blue")

    Start-Backend
    Start-Sleep -Seconds 2
    Start-Admin
    Start-Sleep -Seconds 1
    Start-Visualization

    Wait-ForServices
    Show-StartupInfo
}

# 执行
switch ($Action) {
    "start"   { Start-AllServices }
    "stop"    { Stop-AllServices }
    "restart" { Stop-AllServices; Start-Sleep -Seconds 2; Start-AllServices }
    "status"  { Get-ServiceStatus }
}