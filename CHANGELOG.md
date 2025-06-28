# v2rayA 修复日志

## 修复时间: 2025-06-28

### 1. 防抖问题修复 - 完全解决

**问题描述:**
- 主要文件: `gui/src/components/modalSharing.vue`, `gui/src/node.vue`
- 问题: 各种操作(复制、导入、删除、连接等)快速点击时弹窗会多次弹出
- 影响: 用户体验差，快速点击时出现重复的成功/失败提示，toast消息堆叠

**修复范围:**
- ✅ `modalSharing.vue`: 复制操作防抖 (已完成)
- ✅ `node.vue`: 全面防抖修复 (新增)
  - 复制操作 (14个toast调用)
  - QR码导入错误处理 
  - 服务器/订阅导入成功/失败
  - 服务器删除操作
  - 网络连接/断开操作
  - 延迟测试提示和错误
  - 订阅更新成功通知
  - 服务器创建/修改成功
  - 订阅修改成功

**技术实现:**
- 在组件data中添加 `lastToastTime: 0` 防抖变量
- 每次toast调用前检查时间间隔
- 1秒内重复操作将被忽略，防止toast堆叠

**修复代码模式:**
```javascript
// 防抖变量 (在data中)
lastToastTime: 0,

// 防抖检查 (在每个toast调用前)
const now = Date.now();
if (now - this.lastToastTime > 1000) {
    this.lastToastTime = now;
    this.$buefy.toast.open({...});
}
```

**修改位置:**
- `gui/src/components/modalSharing.vue:95-120` (原有)
- `gui/src/node.vue:781, 858-1604` (新增14处修复)

### 关键修复: 重复弹窗根本原因解决

**发现的核心问题:**
- `modalSharing.vue` 中有两个 `sharingAddressTag` 元素 (第10行和第21行)
- 两个元素都绑定了相同的 ClipboardJS 事件
- 点击任意一个复制按钮会触发两次 success 事件

**冲突问题:**
- `node.vue` 和 `modalSharing.vue` 都初始化了独立的 ClipboardJS 实例
- 造成事件绑定冲突和重复处理

**彻底解决方案:**
1. **区分复制元素**: 改用 `sharingAddressTag-short` 和 `sharingAddressTag-full` 类名
2. **移除重复绑定**: 完全移除 `node.vue` 中的 ClipboardJS 代码
3. **统一管理**: 让 `modalSharing.vue` 组件完全管理自己的复制功能

**最终修改:**
- `modalSharing.vue`: 使用独立的类名，防止重复绑定
- `node.vue`: 移除 ClipboardJS 导入和初始化代码
- 确保单点复制事件管理，彻底解决重复弹窗问题

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