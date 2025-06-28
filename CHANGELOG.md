# v2rayA 修复日志

## 修复时间: 2025-06-28

### 1. 防抖问题修复

**问题描述:**
- 文件: `gui/src/components/modalSharing.vue`
- 问题: 点击 copyLink 按钮时弹窗会多次弹出
- 影响: 用户体验差，快速点击时出现重复的成功/失败提示

**修复方案:**
- 在 `mounted()` 方法中添加防抖机制
- 使用时间戳检查，防止1秒内重复显示 toast 通知
- 修改位置: `gui/src/components/modalSharing.vue:95-120`

**修复代码:**
```javascript
// 添加防抖变量
this.lastToastTime = 0;

// 在成功和错误回调中添加时间检查
const now = Date.now();
if (now - this.lastToastTime > 1000) {
    this.lastToastTime = now;
    // 显示 toast
}
```

### 2. Docker 镜像基础镜像优化

**问题描述:**
- 原 Dockerfile 使用 `mzz2017/git:alpine` 个人镜像
- 缺乏官方支持和安全保障

**修复方案:**
- 替换为官方 `alpine:3.22` 镜像
- 手动安装 git 工具: `RUN apk add --no-cache git`
- 备份原始文件为 `Dockerfile.backup`

### 3. Docker 镜像版本固定

**问题描述:**
- 使用 `latest` 标签导致构建不可重现
- 版本漂移风险

**修复方案:**
- 使用具体版本号和 SHA256 哈希值
- 固定版本:
  - `alpine:3.22@sha256:8a1f59ff...`
  - `node:22-alpine@sha256:5340cbfc...`
  - `golang:1.23-alpine@sha256:68932fa6...`
  - `v2fly/v2fly-core:v5.21.0@sha256:9068d1e6...`

### 4. v2ray 二进制路径问题修复

**问题描述:**
- 容器中 v2ray 二进制文件位于 `/usr/bin/v2ray`
- 默认环境变量指向 `/usr/local/bin/v2ray`
- 导致 v2rayA 无法找到 v2ray 可执行文件

**修复方案:**
- 在 Dockerfile 中添加环境变量: `ENV V2RAYA_V2RAY_BIN=/usr/bin/v2ray`
- 更新运行命令中的环境变量设置

**验证结果:**
- v2ray 版本: V2Ray 5.30.0 (V2Fly)
- geosite.dat: 2.3MB ✅
- geoip.dat: 21MB ✅
- LoyalsoldierSite.dat: 9.5MB ✅

### 5. 文档创建

**新增文件:**
- `docker-run.md`: Docker 运行指南
- `CHANGELOG.md`: 本修复日志

## 最终状态

### Docker 镜像特性
- ✅ 基于官方通用镜像
- ✅ 版本完全固定，构建可重现
- ✅ v2ray 路径问题已修复
- ✅ 防抖问题已解决

### 验证结果
- ✅ Web 界面正常访问: http://localhost:3017
- ✅ v2ray-core 正常运行
- ✅ 所有必需的 .dat 文件存在
- ✅ 容器启动无错误

### 推荐运行命令
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