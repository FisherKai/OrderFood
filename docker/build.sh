#!/bin/bash
# OrderFood Docker 构建脚本

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}   OrderFood Docker 构建脚本${NC}"
echo -e "${GREEN}========================================${NC}"

# 切换到项目根目录
cd "$(dirname "$0")/.."

# 镜像名称和标签
IMAGE_NAME="${IMAGE_NAME:-orderfood}"
IMAGE_TAG="${IMAGE_TAG:-latest}"

echo -e "\n${YELLOW}[1/3] 检查必要文件...${NC}"
required_files=("Dockerfile" "docker-compose.yml" "init.sql" "backend/go.mod" "frontend-admin/package.json" "frontend-mobile/package.json")
for file in "${required_files[@]}"; do
    if [ ! -f "$file" ]; then
        echo -e "${RED}错误: 缺少必要文件 $file${NC}"
        exit 1
    fi
done
echo -e "${GREEN}✓ 所有必要文件存在${NC}"

echo -e "\n${YELLOW}[2/3] 构建Docker镜像...${NC}"
echo "镜像名称: ${IMAGE_NAME}:${IMAGE_TAG}"

docker build \
    --progress=plain \
    -t "${IMAGE_NAME}:${IMAGE_TAG}" \
    .

echo -e "${GREEN}✓ 镜像构建完成${NC}"

echo -e "\n${YELLOW}[3/3] 镜像信息${NC}"
docker images "${IMAGE_NAME}:${IMAGE_TAG}"

echo -e "\n${GREEN}========================================${NC}"
echo -e "${GREEN}   构建完成！${NC}"
echo -e "${GREEN}========================================${NC}"

echo -e "\n使用方式:"
echo -e "  1. 使用docker-compose启动（推荐）:"
echo -e "     ${YELLOW}cp .env.example .env${NC}"
echo -e "     ${YELLOW}# 编辑 .env 修改密码和密钥${NC}"
echo -e "     ${YELLOW}docker-compose up -d${NC}"
echo -e ""
echo -e "  2. 单独运行镜像（需要外部MySQL）:"
echo -e "     ${YELLOW}docker run -d \\${NC}"
echo -e "     ${YELLOW}  -p 80:80 -p 3000:3000 -p 3001:3001 -p 8080:8080 \\${NC}"
echo -e "     ${YELLOW}  -e DB_HOST=your-mysql-host \\${NC}"
echo -e "     ${YELLOW}  -e DB_PASSWORD=your-password \\${NC}"
echo -e "     ${YELLOW}  -e JWT_SECRET=your-jwt-secret \\${NC}"
echo -e "     ${YELLOW}  ${IMAGE_NAME}:${IMAGE_TAG}${NC}"
echo -e ""
echo -e "访问地址:"
echo -e "  - 移动端:     http://localhost:3000"
echo -e "  - 管理后台:   http://localhost:3001"
echo -e "  - 统一入口:   http://localhost (移动端) / http://localhost/admin (管理后台)"
echo -e "  - 后端API:    http://localhost:8080"
