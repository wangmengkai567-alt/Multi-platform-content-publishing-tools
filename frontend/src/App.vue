<template>
  <div class="page-shell">
    <header class="hero">
      <div>
        <p class="eyebrow">Creator Workflow Console</p>
        <h1>多平台内容发布工具</h1>
        <p class="hero__summary">
          一次写作，自动适配公众号、知乎、B站、小红书的内容风格，支持模拟发布和后续平台扩展。
        </p>
      </div>
      <aside class="hero__aside">
        <div>
          <span>架构</span>
          <strong>Go API + Vue 控制台</strong>
        </div>
        <div>
          <span>模式</span>
          <strong>前后端分离 / 平台插件化</strong>
        </div>
      </aside>
    </header>

    <main class="workspace">
      <section class="panel panel--editor">
        <div class="panel__heading">
          <h2>内容输入</h2>
          <button class="ghost-button" @click="fillDemo">填充示例</button>
        </div>

        <div class="form-grid">
          <label>
            标题
            <input v-model="form.title" placeholder="输入内容标题" />
          </label>
          <label>
            语气
            <select v-model="form.tone">
              <option value="professional">专业</option>
              <option value="casual">轻松</option>
              <option value="storytelling">叙事</option>
            </select>
          </label>
          <label class="full">
            摘要
            <textarea v-model="form.summary" rows="3" placeholder="输入内容摘要，用于平台导语和推荐语"></textarea>
          </label>
          <label class="full">
            正文
            <textarea
              v-model="form.body"
              rows="12"
              placeholder="输入正文。可用空行分段，系统会按平台重组结构。"
            ></textarea>
          </label>
          <label>
            标签
            <input v-model="tagInput" placeholder="如：AI工具, 内容创作, 自媒体" />
          </label>
          <label>
            封面图 URL
            <input v-model="form.coverImage" placeholder="可选，用于偏视觉的平台" />
          </label>
        </div>
      </section>

      <section class="panel">
        <div class="panel__heading">
          <h2>平台选择</h2>
          <button class="ghost-button" @click="toggleAll">
            {{ selectedPlatforms.length === platforms.length ? '清空选择' : '全选平台' }}
          </button>
        </div>

        <div v-if="platformLoadError" class="connection-banner">
          <span>平台列表加载失败：{{ platformLoadError }}</span>
          <button class="ghost-button" @click="loadPlatforms">重试加载</button>
        </div>

        <div v-else-if="!platforms.length" class="connection-banner">
          <span>正在连接后端并加载平台列表...</span>
          <button class="ghost-button" @click="loadPlatforms">重新加载</button>
        </div>

        <div class="platform-grid">
          <PlatformCard
            v-for="platform in platforms"
            :key="platform.id"
            :platform="platform"
            :checked="selectedPlatforms.includes(platform.id)"
            @toggle="togglePlatform"
          />
        </div>

        <div class="action-row">
          <label class="switch">
            <input v-model="simulate" type="checkbox" />
            <span>使用模拟发布</span>
          </label>
          <label class="switch">
            <input v-model="enableAnalytics" type="checkbox" />
            <span>生成发布追踪标识</span>
          </label>
          <button class="primary-button" :disabled="loading" @click="buildPreview">
            {{ loading ? '生成中...' : '生成适配预览' }}
          </button>
          <button class="secondary-button" :disabled="loading || !previews.length" @click="publish">
            {{ loading ? '处理中...' : simulate ? '模拟发布' : '发起发布请求' }}
          </button>
        </div>

        <p class="mode-hint">
          {{
            simulate
              ? '当前为模拟发布：仅验证适配和分发流程，不会调用外部平台接口。'
              : '当前为演示发布：界面会发起“真实发布”请求，但后端暂未接入平台官方 API，因此仍不会真正发到公众号、知乎、B站或小红书。'
          }}
        </p>

        <p v-if="errorMessage" class="feedback feedback--error">{{ errorMessage }}</p>
        <p v-if="successMessage" class="feedback feedback--success">{{ successMessage }}</p>
      </section>

      <section class="panel">
        <div class="panel__heading">
          <h2>适配预览</h2>
          <span>{{ previews.length }} 个平台</span>
        </div>

        <div v-if="previews.length" class="preview-grid">
          <PreviewPanel v-for="preview in previews" :key="preview.platformId" :preview="preview" />
        </div>
        <div v-else class="empty-state">
          选择平台并生成预览后，这里会展示不同平台的标题、正文结构、标签和注意事项。
        </div>
      </section>

      <section class="panel">
        <div class="panel__heading">
          <h2>发布结果</h2>
          <span v-if="publishResults.length">请求 {{ publishRequestId }}</span>
        </div>

        <div v-if="publishResults.length" class="result-list">
          <article v-for="result in publishResults" :key="result.externalRef" class="result-card">
            <div>
              <strong>{{ result.platformName }}</strong>
              <span>{{ statusText(result.status) }}</span>
            </div>
            <p>{{ result.message }}</p>
            <small>{{ result.externalRef }} · {{ formatTime(result.publishedAt) }}</small>
          </article>
        </div>
        <div v-else class="empty-state">
          发布后会显示每个平台的处理状态。当前版本只实现了演示流程和模拟发布，尚未接入真实平台发布接口。
        </div>
      </section>

      <section class="panel architecture">
        <div class="panel__heading">
          <h2>扩展架构</h2>
        </div>
        <div class="architecture-grid">
          <div>
            <strong>Adapter 接口</strong>
            <p>每个平台实现 `Descriptor / BuildPreview / Publish` 三个能力，独立管理格式规则和 API 调用。</p>
          </div>
          <div>
            <strong>Service 编排</strong>
            <p>统一做内容校验、平台分发、模拟发布和发布请求跟踪，避免前端直接耦合平台细节。</p>
          </div>
          <div>
            <strong>可扩展方向</strong>
            <p>后续可增加 OAuth 凭证管理、媒体资源上传队列、定时任务、发布日志和回流数据看板。</p>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import PlatformCard from './components/PlatformCard.vue'
import PreviewPanel from './components/PreviewPanel.vue'
import { fetchPlatforms, generatePreviews, publishContent } from './services/api'

const platforms = ref([])
const selectedPlatforms = ref([])
const previews = ref([])
const publishResults = ref([])
const publishRequestId = ref('')
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const platformLoadError = ref('')
const simulate = ref(true)
const enableAnalytics = ref(true)

const form = ref({
  title: '',
  summary: '',
  body: '',
  tags: [],
  coverImage: '',
  tone: 'professional',
})

const tagInput = computed({
  get: () => form.value.tags.join(', '),
  set: (value) => {
    form.value.tags = value
      .split(',')
      .map((item) => item.trim())
      .filter(Boolean)
  },
})

onMounted(async () => {
  await loadPlatforms()
})

async function loadPlatforms() {
  platformLoadError.value = ''
  try {
    platforms.value = await fetchPlatforms()
    if (!selectedPlatforms.value.length) {
      selectedPlatforms.value = platforms.value.map((item) => item.id)
    } else {
      selectedPlatforms.value = selectedPlatforms.value.filter((id) =>
        platforms.value.some((item) => item.id === id),
      )
    }
  } catch (error) {
    platforms.value = []
    selectedPlatforms.value = []
    platformLoadError.value = error.message || '无法连接后端服务'
  }
}

function togglePlatform(id) {
  if (selectedPlatforms.value.includes(id)) {
    selectedPlatforms.value = selectedPlatforms.value.filter((item) => item !== id)
    return
  }
  selectedPlatforms.value = [...selectedPlatforms.value, id]
}

function toggleAll() {
  if (selectedPlatforms.value.length === platforms.value.length) {
    selectedPlatforms.value = []
    return
  }
  selectedPlatforms.value = platforms.value.map((item) => item.id)
}

function fillDemo() {
  form.value = {
    title: '如何把一篇内容高效同步到多个平台',
    summary: '同一份内容在不同平台往往要改标题、改段落、改语气，这正是创作者重复劳动最多的环节。',
    body:
      '很多创作者每天都要在多个平台更新内容，但每个平台的受众和排版习惯都不同。\n\n公众号更适合结构完整的长文，知乎偏好观点清晰、论证充分的表达，B站动态要求节奏紧凑，而小红书更强调口语化和标签曝光。\n\n如果能把内容抽象成统一输入，再通过平台适配器自动重组结构、压缩标题、补充标签和生成发布建议，创作者就能把精力放回创作本身。',
    tags: ['内容创作', '效率工具', '多平台分发', '自媒体运营'],
    coverImage: 'https://images.unsplash.com/photo-1516321318423-f06f85e504b3?auto=format&fit=crop&w=900&q=80',
    tone: 'professional',
  }
}

async function buildPreview() {
  errorMessage.value = ''
  successMessage.value = ''
  loading.value = true

  try {
    if (!platforms.value.length) {
      await loadPlatforms()
    }
    if (!form.value.title.trim()) {
      throw new Error('请先输入标题。')
    }
    if (!form.value.body.trim()) {
      throw new Error('请先输入正文。')
    }

    const response = await generatePreviews({
      content: form.value,
      platforms: selectedPlatforms.value,
    })
    previews.value = response.previews
    successMessage.value = `已生成 ${response.previews.length} 个平台预览。`
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    loading.value = false
  }
}

async function publish() {
  errorMessage.value = ''
  successMessage.value = ''
  loading.value = true

  try {
    const response = await publishContent({
      content: form.value,
      platforms: selectedPlatforms.value,
      simulate: simulate.value,
      enableAnalytics: enableAnalytics.value,
      scheduledAt: '',
    })
    publishResults.value = response.results
    publishRequestId.value = response.requestId
    successMessage.value = simulate.value
      ? `已完成 ${response.results.length} 个平台的模拟发布校验。`
      : `已完成 ${response.results.length} 个平台的演示发布流程，但尚未真正发布到外部平台。`
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    loading.value = false
  }
}

function formatTime(value) {
  return new Date(value).toLocaleString('zh-CN')
}

function statusText(status) {
  if (status === 'simulated') {
    return '模拟完成'
  }
  if (status === 'demo_ready') {
    return '演示完成'
  }
  if (status === 'published') {
    return '已发布'
  }
  return status
}
</script>
