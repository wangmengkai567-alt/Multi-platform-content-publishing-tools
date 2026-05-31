# 多平台内容发布工具

一个前后端分离的内容分发工具原型，面向公众号、知乎、B站、小红书等平台的内容同步发布场景。用户输入一份标准化内容后，系统会根据平台特征自动生成格式化预览，并支持一键模拟发布。

## 技术栈

- 后端：Go
- 前端：Vue 3 + Vite
- 架构：前后端分离，平台适配器插件化

## 目录结构

```text
backend/
  cmd/server                HTTP 服务入口
  internal/adapter          平台适配器接口与内置平台实现
  internal/domain           请求、响应与领域模型
  internal/service          内容编排与发布服务
  internal/transport        HTTP 路由层
frontend/
  src/components            平台卡片、预览面板
  src/services              前端 API 封装
  src/styles                样式文件
```

## 核心能力

- 标准化内容输入：标题、摘要、正文、标签、封面图、语气
- 平台自动适配：针对不同平台生成标题压缩、正文结构重组、标签建议、风险提示
- 一键发布：当前实现为模拟发布流程，已预留真实 API 接入点
- 扩展架构：新增平台只需实现统一 `PlatformAdapter` 接口

## 后端 API

- `GET /health`
- `GET /api/platforms`
- `POST /api/previews`
- `POST /api/publish`

### 预览请求示例

```json
{
  "content": {
    "title": "如何把一篇内容高效同步到多个平台",
    "summary": "同一份内容在不同平台要改标题、改段落、改语气。",
    "body": "正文内容",
    "tags": ["内容创作", "效率工具"],
    "coverImage": "",
    "tone": "professional"
  },
  "platforms": ["wechat", "zhihu", "bilibili", "xiaohongshu"]
}
```

## 平台扩展设计

新增平台时，实现以下三个方法即可：

1. `Descriptor()`：声明平台名称、能力和风格信息，供前端展示。
2. `BuildPreview(content)`：根据平台规则转换内容结构，如标题长度、Markdown 结构、口语化程度、标签数量限制等。
3. `Publish(payload)`：封装真实平台 API 调用。当前项目返回模拟结果，接入正式发布时可在这一层处理 OAuth、媒体上传、限流和错误重试。

建议继续演进的模块：

- 账号授权中心：统一管理各平台 Access Token / Cookie / OAuth 凭证
- 媒体资产中心：图片、封面、视频统一上传并按平台规格转换
- 调度任务中心：支持定时发布、失败重试和幂等控制
- 数据回流中心：采集阅读量、点赞、收藏、评论等指标，形成跨平台运营看板

## 本地运行

### 启动后端

```powershell
cd backend
go run ./cmd/server
```

默认监听 `http://localhost:8080`。

### 启动前端

```powershell
cd frontend
npm install
npm run dev
```

默认访问 `http://localhost:5173`。

demo视频链接：`https://www.bilibili.com/video/BV1y6VU6jEnK/?vd_source=9890724d574f35a16b6d22529571bdf4`
