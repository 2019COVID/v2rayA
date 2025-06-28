# v2rayA Docker 运行指南

## 构建镜像

```bash
docker build -t v2raya:latest .
```

## 运行容器

### 基本运行命令

```bash
docker run -d \
  -p 2017:2017 \
  --restart=always \
  --name v2raya \
  -v /etc/v2raya:/etc/v2raya \
  v2raya:latest
```

### 完整配置运行命令

```bash
docker run -d \
  -p 3017:2017 \
  -p 30170-30172:20170-20172 \
  --restart=always \
  --name v2raya-web \
  -e V2RAYA_V2RAY_BIN=/usr/bin/v2ray \
  -e V2RAYA_LOG_FILE=/tmp/v2raya.log \
  -v /etc/v2raya:/etc/v2raya \
  v2raya:latest
```

## 环境变量说明

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `V2RAYA_V2RAY_BIN` | `/usr/bin/v2ray` | v2ray二进制文件路径 |
| `V2RAYA_LOG_FILE` | - | 日志文件路径 |

## 端口说明

| 端口 | 说明 |
|------|------|
| `2017` | Web管理界面端口 |
| `20170-20172` | v2ray代理端口范围 |

## 数据卷说明

| 路径 | 说明 |
|------|------|
| `/etc/v2raya` | v2rayA配置文件存储目录 |

## 访问方式

启动容器后，通过浏览器访问：
- 默认: http://localhost:2017
- 自定义端口: http://localhost:3017 (如使用上述完整配置)

## 容器管理

### 查看日志
```bash
docker logs v2raya-web
```

### 停止容器
```bash
docker stop v2raya-web
```

### 重新启动
```bash
docker restart v2raya-web
```

### 删除容器
```bash
docker stop v2raya-web && docker rm v2raya-web
```