# 智慧社区 AI 服务平台

基于 **uni-app + Go + Python + MySQL** 的社区管理系统，内置 AI 医疗诊断、知识问答、地图识图、社区周报等智能功能。

---

## 项目架构

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  微信小程序  │────▶│  Go + Gin   │────▶│  MySQL 8.0  │
│  (前端)     │     │  主服务:8080 │     │  数据库      │
└─────────────┘     └──────┬──────┘     └─────────────┘
                           │
                           ▼
                    ┌─────────────┐     ┌─────────────┐
                    │ Python +    │────▶│  GLM AI     │
                    │ FastAPI:8001│     │  大模型      │
                    └─────────────┘     └─────────────┘
```

| 模块 | 技术栈 | 端口 | 目录 |
|---|---|---|---|
| 前端 | uni-app（微信小程序） | - | `frontend/community/` |
| 主后端 | Go + Gin | 8080 | `backend/Service/` |
| AI 服务 | Python + FastAPI | 8001 | `backend/AIend/integration/` |
| 数据库 | MySQL 8.0 | 3306 | `db/community.sql` |

---

## 功能模块

- **用户系统**：注册、登录、JWT 鉴权
- **社区管理**：创建/解散社区、成员管理、社区申请
- **消息聊天**：WebSocket 实时群聊
- **公告通知**：发布、查看最新公告
- **便民服务**：社区申请、失物招领分享
- **AI 医疗助手**：症状诊断、图片识别、医疗日志
- **知识问答**：AI 百科问答（支持图片 + 文字）
- **地图识图**：上传地图图片自动识别
- **社区周报**：AI 自动生成本周社区周报

---

## 快速启动（Docker 一键部署）

### 前置要求

- 已安装 [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- Docker 已启动

### 启动步骤

```bash
# 1. 克隆项目后，进入项目根目录
cd Community_Management_System

# 2. 一键构建并启动所有服务（MySQL + Go + Python）
docker compose up -d --build
```

首次启动会拉取镜像 + 构建，大约需要 3-5 分钟。

### 验证启动

```bash
# 查看容器状态，三个都显示 Up 即成功
docker compose ps

# 输出应该类似：
# NAME                                      STATUS
# community_management_system-go-server-1   Up
# community_management_system-mysql-1       Up
# community_management_system-python-ai-1   Up
```

### 访问服务

| 服务 | 地址 | 说明 |
|---|---|---|
| Go 主后端 | http://127.0.0.1:8080 | API 接口 |
| Python AI 服务 | http://127.0.0.1:8001 | AI 接口 |
| MySQL | 127.0.0.1:3306 | root / 123456 |

---

## 前端启动（微信开发者工具）

1. 打开 **微信开发者工具**
2. 导入项目，选择目录：`frontend/community/`
3. 修改接口地址：`util/configuration.js`

```javascript
const BASEURL_NAME = "http://127.0.0.1:8080"
const WS_BASEURL_NAME = "ws://127.0.0.1:8080"
```

4. 点击编译即可运行

---

## 开发模式（本地调试，不用 Docker）

### 1. 启动 MySQL

本地安装 MySQL 8.0，创建数据库 `community`，导入 `db/community.sql`：

```bash
mysql -u root -p < db/community.sql
```

### 2. 启动 Go 后端

```bash
cd backend/Service
go mod download
go run main.go
```

确保 `.env` 配置正确：

```env
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=123456
DB_NAME=community
SERVER_PORT=8080
CLIENT_URL=http://127.0.0.1:8001
```

### 3. 启动 Python AI 服务

```bash
cd backend/AIend/integration
pip install -r requirements.txt
uvicorn main:app --host 0.0.0.0 --port 8001
```

确保 `.env` 配置正确（填入你的 GLM API Key）：

```env
GLM_BASEURL=https://open.bigmodel.cn/api/paas/v4/
GLM_46V_KEY=你的key
GLM_45V_KEY=你的key
GIN_BASEURL=http://127.0.0.1:8080
```

---

## 常用 Docker 命令

| 命令 | 作用 |
|---|---|
| `docker compose up -d` | 后台启动所有服务 |
| `docker compose up -d --build` | 重新构建并启动 |
| `docker compose down` | 停止并删除容器 |
| `docker compose ps` | 查看运行状态 |
| `docker compose logs -f go-server` | 实时查看 Go 后端日志 |
| `docker compose restart go-server` | 重启某个服务 |

---

## 项目结构

```
Community_Management_System/
├── backend/
│   ├── Service/              # Go 主后端
│   │   ├── main.go           # 入口文件
│   │   ├── server/           # 路由处理（user/community/chat/...）
│   │   ├── Agent/             # AI 模块对接
│   │   ├── Tools/             # 工具库（JWT/加密/SQL）
│   │   ├── .env              # 配置文件
│   │   ├── Dockerfile        # Go 镜像构建
│   │   └── go.mod
│   │
│   └── AIend/integration/    # Python AI 服务
│       ├── main.py          # FastAPI 入口
│       ├── service/          # LLM 调用、业务逻辑
│       ├── .env              # GLM API 配置
│       ├── Dockerfile        # Python 镜像构建
│       └── requirements.txt   # Python 依赖
│
├── db/
│   └── community.sql         # 数据库建表 SQL
│
├── frontend/community/       # uni-app 微信小程序
│   ├── pages/                # 登录页、首页
│   ├── component/            # 各功能页面组件
│   ├── util/                 # 请求封装、API 接口
│   ├── static/               # 静态资源
│   └── pages.json            # 路由配置
│
├── docker-compose.yml        # Docker 编排配置
└── README.md                 # 项目说明
```

---

## 常见问题

### 1. 前端白屏

检查微信开发者工具的详情 → 本地设置 → 勾选「不校验合法域名」。

### 2. 接口 500 错误

- 先看后端日志：`docker compose logs go-server`
- 确认 MySQL 容器是否正常：`docker compose ps`
- 确认数据库表是否导入：进 MySQL 执行 `USE community; SHOW TABLES;`

### 3. AI 问答失败

- 检查 Python 容器是否在运行：`docker compose ps -a`
- 看 Python 日志：`docker logs community_management_system-python-ai-1`
- 确认 `.env` 里的 GLM API Key 是否有效

### 4. 端口冲突

如果本地已占用 3306/8080/8001，修改 `docker-compose.yml` 里的 `ports` 映射，例如 `"3307:3306"`。
