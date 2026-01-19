#!/bin/bash

# 远程调试 SaaS 平台 - 一键部署脚本
# 支持: Ubuntu 20.04+, Debian 10+, CentOS 7+

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# 检查 root 权限
check_root() {
    if [ "$EUID" -ne 0 ]; then
        log_error "请使用 root 用户运行此脚本"
        exit 1
    fi
}

# 检测系统
detect_os() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$ID
        VERSION=$VERSION_ID
    else
        log_error "无法检测操作系统"
        exit 1
    fi
    log_info "检测到系统: $OS $VERSION"
}

# 安装 Docker
install_docker() {
    if command -v docker &> /dev/null; then
        log_info "Docker 已安装"
        return
    fi

    log_info "正在安装 Docker..."

    case $OS in
        ubuntu|debian)
            apt-get update
            apt-get install -y apt-transport-https ca-certificates curl gnupg lsb-release
            curl -fsSL https://download.docker.com/linux/$OS/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
            echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/$OS $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
            apt-get update
            apt-get install -y docker-ce docker-ce-cli containerd.io
            ;;
        centos|rhel|fedora)
            yum install -y yum-utils
            yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
            yum install -y docker-ce docker-ce-cli containerd.io
            ;;
        *)
            log_error "不支持的操作系统: $OS"
            exit 1
            ;;
    esac

    systemctl start docker
    systemctl enable docker
    log_success "Docker 安装完成"
}

# 安装 Docker Compose
install_docker_compose() {
    if command -v docker-compose &> /dev/null; then
        log_info "Docker Compose 已安装"
        return
    fi

    log_info "正在安装 Docker Compose..."
    curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
    log_success "Docker Compose 安装完成"
}

# 配置防火墙
configure_firewall() {
    log_info "配置防火墙..."

    if command -v ufw &> /dev/null; then
        ufw allow 80/tcp
        ufw allow 443/tcp
        ufw allow 3000/tcp
        ufw allow 10000:60000/tcp
        log_success "UFW 防火墙已配置"
    elif command -v firewall-cmd &> /dev/null; then
        firewall-cmd --permanent --add-port=80/tcp
        firewall-cmd --permanent --add-port=443/tcp
        firewall-cmd --permanent --add-port=3000/tcp
        firewall-cmd --permanent --add-port=10000-60000/tcp
        firewall-cmd --reload
        log_success "Firewalld 防火墙已配置"
    else
        log_warn "未检测到防火墙，请手动开放端口 80, 443, 3000, 10000-60000"
    fi
}

# 生成配置
generate_config() {
    local install_dir=$1
    local domain=$2
    local jwt_secret=$(openssl rand -hex 32)
    local admin_password=$3

    log_info "生成配置文件..."

    # 后端配置
    cat > "$install_dir/config.yaml" << EOF
server:
  http_port: 3000
  ws_port: 3001

database:
  path: "./data/app.db"

port_pool:
  min: 10000
  max: 60000

jwt:
  secret: "$jwt_secret"
  expire_hours: 168

pricing:
  port_per_day: 1.0

admin:
  username: "admin"
  password: "$admin_password"
EOF

    # Docker Compose 配置
    cat > "$install_dir/docker-compose.yml" << EOF
version: '3.8'

services:
  app:
    image: ghcr.io/your-username/remote-debug-platform:latest
    container_name: remote-debug-platform
    restart: unless-stopped
    ports:
      - "3000:3000"
      - "10000-60000:10000-60000"
    volumes:
      - ./data:/app/data
      - ./config.yaml:/app/config.yaml
    environment:
      - TZ=Asia/Shanghai

  nginx:
    image: nginx:alpine
    container_name: remote-debug-nginx
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/conf.d/default.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - app
EOF

    # Nginx 配置
    cat > "$install_dir/nginx.conf" << EOF
server {
    listen 80;
    server_name $domain;

    location / {
        proxy_pass http://app:3000;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
    }

    location /ws {
        proxy_pass http://app:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host \$host;
    }
}
EOF

    log_success "配置文件已生成"
}

# 本地构建部署
build_and_deploy() {
    local install_dir=$1

    log_info "正在构建项目..."

    cd "$install_dir"

    # 检查是否有源码
    if [ ! -d "backend" ] || [ ! -d "frontend" ]; then
        log_error "未找到源码目录，请确保在项目根目录运行"
        exit 1
    fi

    # 使用本地构建的 docker-compose
    cat > "$install_dir/docker-compose.yml" << 'EOF'
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: deploy/Dockerfile
    container_name: remote-debug-platform
    restart: unless-stopped
    ports:
      - "3000:3000"
      - "10000-60000:10000-60000"
    volumes:
      - ./deploy/data:/app/data
      - ./deploy/config.yaml:/app/config.yaml
    environment:
      - TZ=Asia/Shanghai

  nginx:
    image: nginx:alpine
    container_name: remote-debug-nginx
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./deploy/nginx.conf:/etc/nginx/conf.d/default.conf
      - ./deploy/ssl:/etc/nginx/ssl
    depends_on:
      - app
EOF

    docker-compose build
    docker-compose up -d

    log_success "部署完成"
}

# 显示部署信息
show_info() {
    local domain=$1
    local admin_password=$2

    echo ""
    echo "========================================"
    echo -e "${GREEN}  部署成功！${NC}"
    echo "========================================"
    echo ""
    echo "  访问地址: http://$domain"
    echo "  管理后台: http://$domain/admin"
    echo ""
    echo "  管理员账号: admin"
    echo "  管理员密码: $admin_password"
    echo ""
    echo "  数据目录: /opt/remote-debug/data"
    echo "  配置文件: /opt/remote-debug/config.yaml"
    echo ""
    echo "  常用命令:"
    echo "    查看日志: docker-compose logs -f"
    echo "    重启服务: docker-compose restart"
    echo "    停止服务: docker-compose down"
    echo ""
    echo "========================================"
}

# 快速安装（从源码）
quick_install() {
    local install_dir=$(pwd)
    local domain=${1:-$(hostname -I | awk '{print $1}')}
    local admin_password=${2:-$(openssl rand -base64 12)}

    log_info "开始快速安装..."
    log_info "安装目录: $install_dir"
    log_info "访问域名: $domain"

    # 生成配置
    mkdir -p "$install_dir/deploy/data"
    mkdir -p "$install_dir/deploy/ssl"

    # 生成 JWT Secret
    local jwt_secret=$(openssl rand -hex 32)

    # 后端配置
    cat > "$install_dir/deploy/config.yaml" << EOF
server:
  http_port: 3000
  ws_port: 3001

database:
  path: "./data/app.db"

port_pool:
  min: 10000
  max: 60000

jwt:
  secret: "$jwt_secret"
  expire_hours: 168

pricing:
  port_per_day: 1.0

admin:
  username: "admin"
  password: "$admin_password"
EOF

    # Nginx 配置
    cat > "$install_dir/deploy/nginx.conf" << EOF
server {
    listen 80;
    server_name $domain;

    location / {
        proxy_pass http://app:3000;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
    }

    location /ws {
        proxy_pass http://app:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host \$host;
    }
}
EOF

    # 构建部署
    build_and_deploy "$install_dir"

    show_info "$domain" "$admin_password"
}

# 交互式安装
interactive_install() {
    echo ""
    echo "========================================"
    echo "  远程调试 SaaS 平台 - 一键部署"
    echo "========================================"
    echo ""

    # 获取域名
    local default_ip=$(hostname -I | awk '{print $1}')
    read -p "请输入域名或IP [$default_ip]: " domain
    domain=${domain:-$default_ip}

    # 获取管理员密码
    local default_password=$(openssl rand -base64 12)
    read -p "请输入管理员密码 [$default_password]: " admin_password
    admin_password=${admin_password:-$default_password}

    # 确认
    echo ""
    echo "配置信息:"
    echo "  域名/IP: $domain"
    echo "  管理员密码: $admin_password"
    echo ""
    read -p "确认开始安装? [Y/n]: " confirm
    confirm=${confirm:-Y}

    if [[ ! $confirm =~ ^[Yy]$ ]]; then
        log_warn "安装已取消"
        exit 0
    fi

    check_root
    detect_os
    install_docker
    install_docker_compose
    configure_firewall
    quick_install "$domain" "$admin_password"
}

# 显示帮助
show_help() {
    echo "远程调试 SaaS 平台 - 一键部署脚本"
    echo ""
    echo "用法: $0 [命令] [选项]"
    echo ""
    echo "命令:"
    echo "  install          交互式安装（默认）"
    echo "  quick [域名]     快速安装"
    echo "  start            启动服务"
    echo "  stop             停止服务"
    echo "  restart          重启服务"
    echo "  logs             查看日志"
    echo "  status           查看状态"
    echo "  uninstall        卸载"
    echo "  help             显示帮助"
    echo ""
    echo "示例:"
    echo "  $0 install                  # 交互式安装"
    echo "  $0 quick example.com        # 快速安装到指定域名"
    echo "  $0 logs                     # 查看日志"
}

# 服务管理
service_start() {
    cd "$(dirname "$0")"
    docker-compose up -d
    log_success "服务已启动"
}

service_stop() {
    cd "$(dirname "$0")"
    docker-compose down
    log_success "服务已停止"
}

service_restart() {
    cd "$(dirname "$0")"
    docker-compose restart
    log_success "服务已重启"
}

service_logs() {
    cd "$(dirname "$0")"
    docker-compose logs -f
}

service_status() {
    cd "$(dirname "$0")"
    docker-compose ps
}

service_uninstall() {
    read -p "确定要卸载吗？数据将被保留 [y/N]: " confirm
    if [[ $confirm =~ ^[Yy]$ ]]; then
        cd "$(dirname "$0")"
        docker-compose down --rmi all
        log_success "已卸载，数据保留在 deploy/data 目录"
    fi
}

# 主入口
main() {
    case "${1:-install}" in
        install)
            interactive_install
            ;;
        quick)
            check_root
            detect_os
            install_docker
            install_docker_compose
            configure_firewall
            quick_install "$2" "$3"
            ;;
        start)
            service_start
            ;;
        stop)
            service_stop
            ;;
        restart)
            service_restart
            ;;
        logs)
            service_logs
            ;;
        status)
            service_status
            ;;
        uninstall)
            service_uninstall
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            log_error "未知命令: $1"
            show_help
            exit 1
            ;;
    esac
}

main "$@"
