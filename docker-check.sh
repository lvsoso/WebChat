#!/bin/bash

# 颜色定义
GREEN="\033[0;32m"
YELLOW="\033[1;33m"
RED="\033[0;31m"
BLUE="\033[0;34m"
NC="\033[0m" # No Color

echo -e "${BLUE}=== Docker 环境检查工具 ===${NC}\n"

# 检查Docker是否安装
echo -e "${BLUE}[1/4] 检查Docker安装状态${NC}"
if ! command -v docker &> /dev/null; then
    echo -e "${RED}✗ Docker未安装${NC}"
    echo -e "${YELLOW}请先安装Docker: https://docs.docker.com/get-docker/${NC}"
    exit 1
else
    echo -e "${GREEN}✓ Docker已安装${NC}"
    docker --version
fi

# 检查Docker Compose是否安装
echo -e "\n${BLUE}[2/4] 检查Docker Compose安装状态${NC}"
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}✗ Docker Compose未安装${NC}"
    echo -e "${YELLOW}请先安装Docker Compose: https://docs.docker.com/compose/install/${NC}"
    exit 1
else
    echo -e "${GREEN}✓ Docker Compose已安装${NC}"
    docker-compose --version
fi

# 检查Docker守护进程是否运行
echo -e "\n${BLUE}[3/4] 检查Docker守护进程状态${NC}"
if ! docker info &> /dev/null; then
    echo -e "${RED}✗ Docker守护进程未运行${NC}"
    echo -e "${YELLOW}请启动Docker守护进程:${NC}"
    echo -e "${YELLOW}  - macOS: 启动Docker Desktop应用${NC}"
    echo -e "${YELLOW}  - Linux: sudo systemctl start docker${NC}"
    echo -e "${YELLOW}  - Windows: 启动Docker Desktop应用${NC}"
    
    # 检查OrbStack (macOS特有)
    if [[ "$(uname)" == "Darwin" ]] && [[ -d "$HOME/.orbstack" ]]; then
        echo -e "\n${YELLOW}检测到OrbStack环境:${NC}"
        echo -e "${YELLOW}  - 请确保OrbStack应用已启动${NC}"
        echo -e "${YELLOW}  - 或者在终端运行: open -a OrbStack${NC}"
    fi
    
    exit 1
else
    echo -e "${GREEN}✓ Docker守护进程正在运行${NC}"
fi

# 尝试拉取测试镜像
echo -e "\n${BLUE}[4/4] 测试Docker网络连接${NC}"
echo -e "正在拉取测试镜像..."
if ! docker pull hello-world &> /dev/null; then
    echo -e "${RED}✗ 无法拉取测试镜像${NC}"
    echo -e "${YELLOW}可能的原因:${NC}"
    echo -e "${YELLOW}  - 网络连接问题${NC}"
    echo -e "${YELLOW}  - Docker Hub访问受限${NC}"
    echo -e "${YELLOW}  - 代理或防火墙设置${NC}"
else
    echo -e "${GREEN}✓ 成功拉取测试镜像${NC}"
fi

echo -e "\n${BLUE}=== 检查完成 ===${NC}"

# 如果所有检查都通过，提供启动命令
if docker info &> /dev/null; then
    echo -e "\n${GREEN}Docker环境正常，可以使用以下命令启动服务:${NC}"
    echo -e "${GREEN}cd $(dirname "$0") && docker-compose up -d${NC}"
fi