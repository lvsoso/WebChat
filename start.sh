#!/bin/bash

# 颜色定义
GREEN="\033[0;32m"
YELLOW="\033[1;33m"
BLUE="\033[0;34m"
NC="\033[0m" # No Color

echo -e "${BLUE}=== Web Chat Service 启动脚本 ===${NC}\n"

# 检查docker是否安装
if ! command -v docker &> /dev/null; then
    echo -e "${YELLOW}警告: Docker未安装，请先安装Docker${NC}"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo -e "${YELLOW}警告: Docker Compose未安装，请先安装Docker Compose${NC}"
    exit 1
fi

# 启动数据库服务
echo -e "${GREEN}[1/4] 启动数据库服务...${NC}"
docker-compose up -d postgres redis
echo -e "${GREEN}✓ 数据库服务启动成功${NC}\n"

# 等待PostgreSQL就绪
echo -e "${GREEN}[2/4] 等待PostgreSQL就绪...${NC}"
sleep 5
echo -e "${GREEN}✓ PostgreSQL就绪${NC}\n"

# 启动后端服务
echo -e "${GREEN}[3/4] 启动后端服务...${NC}"
cd backend
go mod tidy
go run main.go &
BACKEND_PID=$!
cd ..
echo -e "${GREEN}✓ 后端服务启动成功 (PID: $BACKEND_PID)${NC}\n"

# 启动前端服务
echo -e "${GREEN}[4/4] 启动前端服务...${NC}"
cd frontend
npm install
npm run dev &
FRONTEND_PID=$!
cd ..
echo -e "${GREEN}✓ 前端服务启动成功 (PID: $FRONTEND_PID)${NC}\n"

echo -e "${BLUE}=== 所有服务已启动 ===${NC}"
echo -e "${GREEN}前端访问地址: http://localhost:5173${NC}"
echo -e "${GREEN}后端API地址: http://localhost:8080${NC}\n"

echo -e "按Ctrl+C停止所有服务\n"

# 捕获SIGINT信号(Ctrl+C)，优雅关闭所有服务
trap "echo -e '\n${YELLOW}正在停止服务...${NC}'; kill $FRONTEND_PID $BACKEND_PID; docker-compose down; echo -e '${GREEN}所有服务已停止${NC}'; exit 0" INT

# 保持脚本运行
wait