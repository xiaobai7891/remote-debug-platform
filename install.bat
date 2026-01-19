@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

title 远程调试 SaaS 平台 - Windows 部署

echo.
echo ========================================
echo   远程调试 SaaS 平台 - Windows 部署
echo ========================================
echo.

:: 检查 Docker
where docker >nul 2>nul
if %errorlevel% neq 0 (
    echo [ERROR] 未检测到 Docker，请先安装 Docker Desktop
    echo 下载地址: https://www.docker.com/products/docker-desktop
    pause
    exit /b 1
)

:: 检查 Docker 是否运行
docker info >nul 2>nul
if %errorlevel% neq 0 (
    echo [ERROR] Docker 未运行，请启动 Docker Desktop
    pause
    exit /b 1
)

echo [INFO] Docker 已就绪

:: 获取 IP
for /f "tokens=2 delims=:" %%a in ('ipconfig ^| findstr /c:"IPv4"') do (
    set IP=%%a
    set IP=!IP: =!
    goto :gotip
)
:gotip

:: 生成密码
set "CHARS=ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
set "PASSWORD="
for /L %%i in (1,1,12) do (
    set /a "idx=!random! %% 62"
    for %%j in (!idx!) do set "PASSWORD=!PASSWORD!!CHARS:~%%j,1!"
)

echo.
echo 配置信息:
echo   IP 地址: %IP%
echo   管理员密码: %PASSWORD%
echo.

set /p CONFIRM="确认开始安装? [Y/n]: "
if /i "%CONFIRM%"=="n" (
    echo 安装已取消
    pause
    exit /b 0
)

:: 创建目录
if not exist "deploy\data" mkdir "deploy\data"
if not exist "deploy\ssl" mkdir "deploy\ssl"

:: 生成配置文件
echo [INFO] 生成配置文件...

:: 生成 JWT Secret (简化版)
set "JWT_SECRET=%RANDOM%%RANDOM%%RANDOM%%RANDOM%"

(
echo server:
echo   http_port: 3000
echo   ws_port: 3001
echo.
echo database:
echo   path: "./data/app.db"
echo.
echo port_pool:
echo   min: 10000
echo   max: 60000
echo.
echo jwt:
echo   secret: "%JWT_SECRET%"
echo   expire_hours: 168
echo.
echo pricing:
echo   port_per_day: 1.0
echo.
echo admin:
echo   username: "admin"
echo   password: "%PASSWORD%"
) > deploy\config.yaml

:: 生成 Nginx 配置
(
echo server {
echo     listen 80;
echo     server_name %IP%;
echo.
echo     location / {
echo         proxy_pass http://app:3000;
echo         proxy_set_header Host $host;
echo         proxy_set_header X-Real-IP $remote_addr;
echo     }
echo.
echo     location /ws {
echo         proxy_pass http://app:3000;
echo         proxy_http_version 1.1;
echo         proxy_set_header Upgrade $http_upgrade;
echo         proxy_set_header Connection "upgrade";
echo     }
echo }
) > deploy\nginx.conf

:: 构建和启动
echo [INFO] 正在构建项目 (首次可能需要几分钟)...
docker-compose -f docker-compose.yml build

echo [INFO] 启动服务...
docker-compose -f docker-compose.yml up -d

echo.
echo ========================================
echo   部署成功！
echo ========================================
echo.
echo   访问地址: http://%IP%
echo   管理后台: http://%IP%/admin
echo.
echo   管理员账号: admin
echo   管理员密码: %PASSWORD%
echo.
echo   常用命令:
echo     查看日志: docker-compose logs -f
echo     重启服务: docker-compose restart
echo     停止服务: docker-compose down
echo.
echo ========================================
echo.

pause
