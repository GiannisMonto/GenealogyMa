@echo off
REM ===================================
REM 族谱数字化管理平台 - 启动脚本 (Windows)
REM ===================================

setlocal enabledelayedexpansion

REM 服务端口
set BACKEND_PORT=8080
set ADMIN_PORT=3000
set VISUALIZATION_PORT=5174

REM 日志目录
set LOG_DIR=.\logs
set TIMESTAMP=%date:~0,4%%date:~5,2%%date:~8,2%_%time:~0,2%%time:~3,2%%time:~6,2%
set TIMESTAMP=%TIMESTAMP: =0%

REM 创建日志目录
if not exist "%LOG_DIR%" mkdir "%LOG_DIR%"

REM 颜色代码 (通过调用 ANSI_BS 实现)
set NC=
set GREEN=[92m
set YELLOW=[93m
set BLUE=[94m
set RED=[91m

:print_title
echo.
echo ========================================
echo    族谱数字化管理平台 - 启动脚本
echo ========================================
echo.

:check_ports
echo [%BLUE%检查%NC%] 检查端口占用情况...
netstat -ano | findstr :%BACKEND_PORT% >nul
if %errorlevel% equ 0 (
    echo [%YELLOW%警告%NC%] 端口 %BACKEND_PORT% (后端API) 已被占用
) else (
    echo [%GREEN%OK%NC%] 端口 %BACKEND_PORT% (后端API) 可用
)

netstat -ano | findstr :%ADMIN_PORT% >nul
if %errorlevel% equ 0 (
    echo [%YELLOW%警告%NC%] 端口 %ADMIN_PORT% (管理后台) 已被占用
) else (
    echo [%GREEN%OK%NC%] 端口 %ADMIN_PORT% (管理后台) 可用
)

netstat -ano | findstr :%VISUALIZATION_PORT% >nul
if %errorlevel% equ 0 (
    echo [%YELLOW%警告%NC%] 端口 %VISUALIZATION_PORT% (可视化前端) 已被占用
) else (
    echo [%GREEN%OK%NC%] 端口 %VISUALIZATION_PORT% (可视化前端) 可用
)

:start_backend
echo.
echo [%BLUE%1/3%NC%] 启动后端服务 (Go)...
cd backend
start "GenealogyMa_Backend" cmd /c "go run cmd/api/main.go > ..\%LOG_DIR%\backend_%TIMESTAMP%.log 2>&1"
cd ..
echo [%GREEN%启动%NC%] 后端服务已启动
echo [%GREEN%日志%NC%] %LOG_DIR%\backend_%TIMESTAMP%.log

:start_admin
echo.
echo [%BLUE%2/3%NC%] 启动管理后台 (React)...
cd frontend\admin
start "GenealogyMa_Admin" cmd /c "npm run dev > ..\..\%LOG_DIR%\admin_%TIMESTAMP%.log 2>&1"
cd ..\..
echo [%GREEN%启动%NC%] 管理后台已启动
echo [%GREEN%日志%NC%] %LOG_DIR%\admin_%TIMESTAMP%.log

:start_visualization
echo.
echo [%BLUE%3/3%NC%] 启动可视化前端 (Vue 3)...
cd frontend\visualization
start "GenealogyMa_Visualization" cmd /c "npm run dev > ..\..\%LOG_DIR%\visualization_%TIMESTAMP%.log 2>&1"
cd ..\..
echo [%GREEN%启动%NC%] 可视化前端已启动
echo [%GREEN%日志%NC%] %LOG_DIR%\visualization_%TIMESTAMP%.log

:wait_services
echo.
echo [%BLUE%等待服务启动...%NC%]

REM 等待后端 (最多30秒)
echo 等待后端服务 (端口 %BACKEND_PORT%)...
for /L %%i in (1,1,30) do (
    curl -s http://localhost:%BACKEND_PORT%/api/health >nul 2>&1
    if !errorlevel! equ 0 goto :backend_ready
    timeout /t 1 >nul
)
:backend_ready
echo  [%GREEN%OK%NC%]

REM 等待前端 (最多20秒)
echo 等待管理后台 (端口 %ADMIN_PORT%)...
for /L %%i in (1,1,20) do (
    curl -s http://localhost:%ADMIN_PORT% >nul 2>&1
    if !errorlevel! equ 0 goto :admin_ready
    timeout /t 1 >nul
)
:admin_ready
echo  [%GREEN%OK%NC%]

echo 等待可视化前端 (端口 %VISUALIZATION_PORT%)...
for /L %%i in (1,1,20) do (
    curl -s http://localhost:%VISUALIZATION_PORT% >nul 2>&1
    if !errorlevel! equ 0 goto :vis_ready
    timeout /t 1 >nul
)
:vis_ready
echo  [%GREEN%OK%NC%]

:print_info
echo.
echo ========================================
echo    服务已启动!
echo ========================================
echo.
echo   管理后台:  http://localhost:%ADMIN_PORT%
echo   可视化:    http://localhost:%VISUALIZATION_PORT%
echo   API文档:   http://localhost:%BACKEND_PORT%/swagger/index.html
echo.
echo   停止服务:   scripts^\stop.bat
echo.
goto :end

:stop
echo.
echo [%YELLOW%正在停止所有服务...%NC%]
taskkill /FI "WINDOWTITLE eq GenealogyMa_*" /F >nul 2>&1
echo [%GREEN%所有服务已停止%NC%]
goto :end

:status
echo [%BLUE%检查服务状态...%NC%]
tasklist /FI "WINDOWTITLE eq GenealogyMa_Backend" | findstr /i "cmd.exe" >nul
if %errorlevel% equ 0 (echo [%GREEN%运行中%NC%] 后端服务) else (echo [%RED%已停止%NC%] 后端服务)
tasklist /FI "WINDOWTITLE eq GenealogyMa_Admin" | findstr /i "cmd.exe" >nul
if %errorlevel% equ 0 (echo [%GREEN%运行中%NC%] 管理后台) else (echo [%RED%已停止%NC%] 管理后台)
tasklist /FI "WINDOWTITLE eq GenealogyMa_Visualization" | findstr /i "cmd.exe" >nul
if %errorlevel% equ 0 (echo [%GREEN%运行中%NC%] 可视化前端) else (echo [%RED%已停止%NC%] 可视化前端)
goto :end

:usage
echo 用法: start.bat {start^|stop^|restart^|status}
echo.
echo   start   - 启动所有服务 (默认)
echo   stop    - 停止所有服务
echo   restart - 重启所有服务
echo   status  - 查看服务状态

:end
endlocal