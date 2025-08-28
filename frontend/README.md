# sunflower

Vue3 + Vite + TypeScript sunflower frontend

## Recommended IDE Setup

[VSCode](https://code.visualstudio.com/) + [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur).

## Type Support for `.vue` Imports in TS

TypeScript cannot handle type information for `.vue` imports by default, so we replace the `tsc` CLI with `vue-tsc` for type checking. In editors, we need [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) to make the TypeScript language service aware of `.vue` types.

## Customize configuration

See [Vite Configuration Reference](https://vite.dev/config/).

## 视频教程

[B站视频教程](https://www.bilibili.com/cheese/play/ss167267536)


## Project Setup

前置依赖，需要安装 node 和 npm .

具体参考：[https://nodejs.org/zh-cn/download](https://nodejs.org/zh-cn/download)

**重要：执行以下命令前，确保处于 `sunflower/frontend` 目录下。**

### 安装依赖

```sh
npm install
```

### 本地启动

```sh
npm run dev
```

### 编译发布

```sh
npm run build
```

### Lint with [ESLint](https://eslint.org/)

```sh
npm run lint
```

## API 配置

本项目支持灵活的 API 配置，可以根据不同环境（开发、测试、生产）设置不同的后端 API 地址。

### 环境变量配置

项目使用 Vite 的环境变量机制进行 API 配置：

#### 配置文件

- `.env` - 默认环境配置（所有环境共享）
- `.env.development` - 开发环境配置
- `.env.production` - 生产环境配置

#### 可用的环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `VITE_API_BASE_URL` | API 基础地址 | `http://127.0.0.1:8000/api/v1` |
| `VITE_API_TIMEOUT` | 请求超时时间（毫秒） | `10000` |
| `VITE_ENABLE_REQUEST_LOG` | 是否启用请求日志 | `true` |

### 使用示例

#### 开发环境 (.env.development)

```bash
# 开发环境配置
VITE_API_BASE_URL=http://127.0.0.1:8000/api/v1
VITE_API_TIMEOUT=10000
VITE_ENABLE_REQUEST_LOG=true
```

#### 生产环境 (.env.production)

```bash
# 生产环境配置
VITE_API_BASE_URL=https://api.yourdomain.com/api/v1
VITE_API_TIMEOUT=15000
VITE_ENABLE_REQUEST_LOG=false
```

### 部署配置

#### 方式一：构建时配置（当前方案）

在不同环境构建时使用对应的环境变量：

```bash
# 开发环境构建
npm run dev

# 生产环境构建
npm run build
```

#### 方式二：CI/CD 环境变量覆盖

在 CI/CD 流程中通过环境变量覆盖：

```bash
# GitHub Actions / GitLab CI 示例
VITE_API_BASE_URL=https://api.production.com/api/v1 npm run build
```

#### 方式三：Docker 部署

```dockerfile
# Dockerfile 示例
FROM node:18-alpine as builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .

# 通过构建参数传入 API 地址
ARG VITE_API_BASE_URL
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
```

```bash
# 构建不同环境的镜像
docker build --build-arg VITE_API_BASE_URL=https://api.staging.com/api/v1 -t app:staging .
docker build --build-arg VITE_API_BASE_URL=https://api.production.com/api/v1 -t app:prod .
```

### 开发代理配置

开发环境使用 Vite 代理解决跨域问题，配置在 `vite.config.ts`：

```typescript
server: {
  proxy: {
    '/api': {
      target: env.VITE_API_BASE_URL?.replace('/api/v1', '') || 'http://127.0.0.1:8000/',
      changeOrigin: true,
    },
  },
}
```

### 最佳实践

1. **环境隔离**：不同环境使用不同的配置文件，避免配置混乱
2. **安全性**：生产环境关闭请求日志，避免敏感信息泄露
3. **超时设置**：生产环境适当增加超时时间，应对网络延迟
4. **构建优化**：在 CI/CD 中通过环境变量动态设置 API 地址
5. **容器化部署**：使用构建参数传入不同环境的配置

### 配置验证

启动应用后，可以在浏览器控制台查看当前使用的 API 配置：

```javascript
// 在浏览器控制台执行
console.log('当前 API 配置：', window.API_CONFIG || '未找到配置');
```

### 故障排查

如果 API 请求失败，请检查：

1. **网络连接**：确保能够访问配置的 API 地址
2. **环境变量**：检查 `.env` 文件中的配置是否正确
3. **代理设置**：开发环境检查 `vite.config.ts` 中的代理配置
4. **CORS 策略**：确保后端 API 允许前端域名的跨域请求

更多详细的 API 使用说明，请参考 `src/api/README.md` 文件。
