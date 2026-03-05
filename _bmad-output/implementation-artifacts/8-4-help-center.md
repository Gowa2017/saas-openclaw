# Story 8.4: 帮助中心页面

Status: ready-for-dev

## Story

As a 用户,
I want 访问帮助中心获取各类帮助,
so that 遇到问题时可以自助解决。

## Acceptance Criteria

1. **AC1: 帮助入口可见性**
   - **Given** 用户在任意页面
   - **When** 查看页面布局
   - **Then** 显示"帮助"入口（导航栏或固定按钮）
   - **And** 入口位置明显且易于访问

2. **AC2: 帮助中心页面展示**
   - **Given** 用户点击"帮助"入口
   - **When** 页面加载完成
   - **Then** 显示帮助中心页面 `/help`
   - **And** 页面布局清晰，内容分类展示

3. **AC3: 快速入门指南**
   - **Given** 用户在帮助中心页面
   - **When** 查看"快速入门"区域
   - **Then** 显示快速入门指南内容
   - **And** 包含关键操作步骤概述
   - **And** 提供跳转到详细教程的链接

4. **AC4: 常见问题（FAQ）**
   - **Given** 用户在帮助中心页面
   - **When** 查看"常见问题"区域
   - **Then** 显示 FAQ 列表
   - **And** FAQ 按分类组织
   - **And** 支持展开/收起答案

5. **AC5: 视频教程链接**
   - **Given** 用户在帮助中心页面
   - **When** 查看"视频教程"区域
   - **Then** 显示视频教程卡片列表
   - **And** 每个卡片显示视频标题、时长、封面
   - **And** 点击卡片跳转到视频播放页面

6. **AC6: 联系客服入口**
   - **Given** 用户在帮助中心页面
   - **When** 点击"联系客服"按钮
   - **Then** 显示客服联系方式或在线咨询窗口
   - **And** 支持多种联系方式（在线、邮件、电话）

7. **AC7: 关键词搜索功能**
   - **Given** 用户在帮助中心页面
   - **When** 在搜索框输入关键词
   - **Then** 实时搜索匹配的帮助内容
   - **And** 显示搜索结果列表
   - **And** 高亮匹配的关键词

8. **AC8: 无障碍访问**
   - **Given** 用户使用屏幕阅读器
   - **When** 访问帮助中心页面
   - **Then** 所有内容可被正确读取
   - **And** 支持键盘导航
   - **And** 搜索功能可通过键盘操作

## Tasks / Subtasks

- [ ] Task 1: 创建帮助中心页面路由 (AC: 1, 2)
  - [ ] 1.1 在 `frontend/src/router/` 添加帮助中心路由
  - [ ] 1.2 创建 `frontend/src/views/help/HelpCenter.vue`
  - [ ] 1.3 在导航栏添加"帮助"入口

- [ ] Task 2: 实现快速入门区域 (AC: 3)
  - [ ] 2.1 创建 `frontend/src/components/help/QuickStart.vue`
  - [ ] 2.2 编写快速入门内容
  - [ ] 2.3 添加跳转链接到详细教程

- [ ] Task 3: 实现 FAQ 区域 (AC: 4)
  - [ ] 3.1 创建 `frontend/src/components/help/HelpFaq.vue`
  - [ ] 3.2 定义 FAQ 数据结构（可复用 Story 8.2 FaqAccordion）
  - [ ] 3.3 编写 FAQ 内容数据
  - [ ] 3.4 实现分类筛选功能

- [ ] Task 4: 实现视频教程区域 (AC: 5)
  - [ ] 4.1 创建 `frontend/src/components/help/VideoTutorials.vue`
  - [ ] 4.2 创建 `frontend/src/components/help/VideoCard.vue`
  - [ ] 4.3 定义视频教程数据结构
  - [ ] 4.4 实现视频链接跳转

- [ ] Task 5: 实现联系客服入口 (AC: 6)
  - [ ] 5.1 创建 `frontend/src/components/help/ContactSupport.vue`
  - [ ] 5.2 实现客服联系方式展示
  - [ ] 5.3 实现在线咨询弹窗（可选）

- [ ] Task 6: 实现搜索功能 (AC: 7)
  - [ ] 6.1 创建 `frontend/src/components/help/HelpSearch.vue`
  - [ ] 6.2 创建搜索结果组件
  - [ ] 6.3 实现前端搜索逻辑
  - [ ] 6.4 实现关键词高亮

- [ ] Task 7: 无障碍优化 (AC: 8)
  - [ ] 7.1 检查标题层级结构
  - [ ] 7.2 添加 ARIA 属性
  - [ ] 7.3 实现键盘导航
  - [ ] 7.4 测试屏幕阅读器兼容性

- [ ] Task 8: 编写测试 (AC: 1-8)
  - [ ] 8.1 编写页面渲染测试
  - [ ] 8.2 编写 FAQ 组件测试
  - [ ] 8.3 编写搜索功能测试
  - [ ] 8.4 编写无障碍测试

## Dev Notes

### 架构模式与约束

**必须遵循的 UX 设计规范 [Source: ux-design-specification.md]:**

1. **布局策略**: 卡片式仪表盘布局
2. **无障碍要求**: WCAG 2.1 AA 级合规
3. **最小触控区域**: 44x44px
4. **色彩对比度**: 4.5:1

**技术栈要求 [Source: architecture.md]:**
- 前端框架: Vite + Vue 3 + TypeScript
- UI 组件库: Naive UI
- 状态管理: Pinia

### 现有项目状态

**前端项目结构 [Source: 1-2-frontend-project-init.md]:**

```
frontend/
├── src/
│   ├── views/              # 页面组件
│   ├── components/         # 公共组件
│   ├── router/             # 路由配置
│   └── assets/             # 静态资源
```

**可复用组件 (Story 8.1, 8.2):**
- `TutorialStep.vue` - 步骤展示组件
- `TutorialImage.vue` - 图片展示组件
- `FaqAccordion.vue` - FAQ 折叠组件
- `FeedbackButton.vue` - 反馈按钮组件

### 页面布局设计

**帮助中心页面布局:**

```
+--------------------------------------------------+
|                    导航栏                         |
|  [Logo] [首页] [实例] [配置] [帮助]      [用户]  |
+--------------------------------------------------+
|                                                  |
|  [搜索框: 搜索帮助内容...]                        |
|                                                  |
+--------------------------------------------------+
|                                                  |
|  快速入门                                         |
|  +------------+  +------------+  +------------+ |
|  | 配置飞书   |  | 部署实例   |  | 开始使用   | |
|  +------------+  +------------+  +------------+ |
|                                                  |
+--------------------------------------------------+
|                                                  |
|  常见问题 (FAQ)                                  |
|  +------------------------------------------+   |
|  | > 如何配置飞书应用？                      |   |
|  | > 机器人没有回复怎么办？                  |   |
|  | > 如何重置配置？                          |   |
|  | ...                                       |   |
|  +------------------------------------------+   |
|                                                  |
+--------------------------------------------------+
|                                                  |
|  视频教程                                        |
|  +--------+  +--------+  +--------+             |
|  | 视频1  |  | 视频2  |  | 视频3  |             |
|  +--------+  +--------+  +--------+             |
|                                                  |
+--------------------------------------------------+
|                                                  |
|  需要更多帮助？                                  |
|  [联系客服]  [提交反馈]                          |
|                                                  |
+--------------------------------------------------+
```

### 内容数据设计

**快速入门内容:**

```typescript
interface QuickStartItem {
  id: string;
  title: string;
  description: string;
  icon: string;
  link: string;
}

const quickStartItems: QuickStartItem[] = [
  {
    id: 'config-feishu',
    title: '配置飞书应用',
    description: '了解如何在飞书开放平台创建应用并获取配置信息',
    icon: 'settings',
    link: '/help/feishu-platform-tutorial'
  },
  {
    id: 'deploy-instance',
    title: '部署 OpenClaw 实例',
    description: '一键部署您的专属 AI 助手实例',
    icon: 'cloud-upload',
    link: '/instances'
  },
  {
    id: 'use-bot',
    title: '开始使用',
    description: '在飞书中添加机器人并开始对话',
    icon: 'chat-bubbles',
    link: '/help/feishu-bot-tutorial'
  }
];
```

**FAQ 内容:**

```typescript
interface FaqCategory {
  id: string;
  name: string;
  items: FaqItem[];
}

const faqCategories: FaqCategory[] = [
  {
    id: 'getting-started',
    name: '快速入门',
    items: [
      {
        id: 'faq-1',
        question: '如何开始使用 OpenClaw？',
        answer: '1. 配置飞书应用（获取 App ID 和 Secret）\n2. 点击"启动实例"按钮\n3. 在飞书中添加机器人\n4. 开始对话！'
      },
      {
        id: 'faq-2',
        question: '部署需要多长时间？',
        answer: '通常在 3 分钟内完成部署。部署过程中请勿关闭页面。'
      }
    ]
  },
  {
    id: 'feishu-config',
    name: '飞书配置',
    items: [
      {
        id: 'faq-3',
        question: 'App ID 格式是什么？',
        answer: 'App ID 以 "cli_" 开头，例如：cli_a1b2c3d4e5f6g7h8'
      },
      {
        id: 'faq-4',
        question: '配置验证失败怎么办？',
        answer: '请检查：1) App ID 和 Secret 是否正确复制\n2) 应用是否已发布\n3) 权限是否已配置'
      }
    ]
  },
  {
    id: 'instance',
    name: '实例管理',
    items: [
      {
        id: 'faq-5',
        question: '实例状态有哪些？',
        answer: '运行中：实例正常运行\n部署中：正在部署实例\n已停止：实例已停止\n错误：部署或运行出错'
      },
      {
        id: 'faq-6',
        question: '如何重启实例？',
        answer: '在实例详情页点击"重启"按钮即可。重启不会丢失数据。'
      }
    ]
  },
  {
    id: 'bot-usage',
    name: '机器人使用',
    items: [
      {
        id: 'faq-7',
        question: '机器人没有回复怎么办？',
        answer: '请检查：1) 实例状态是否为"运行中"\n2) 消息是否发送到正确的机器人\n3) 查看实例日志排查问题'
      },
      {
        id: 'faq-8',
        question: '支持哪些消息类型？',
        answer: '目前支持文本消息。图片、文件等类型暂不支持。'
      }
    ]
  }
];
```

**视频教程内容:**

```typescript
interface VideoTutorial {
  id: string;
  title: string;
  description: string;
  duration: string;
  thumbnail: string;
  url: string;
}

const videoTutorials: VideoTutorial[] = [
  {
    id: 'video-1',
    title: 'OpenClaw 快速入门',
    description: '5 分钟了解 OpenClaw 的核心功能',
    duration: '5:30',
    thumbnail: '/assets/images/tutorials/video-1-thumb.png',
    url: 'https://example.com/videos/quick-start'
  },
  {
    id: 'video-2',
    title: '飞书应用配置详解',
    description: '手把手教你配置飞书应用',
    duration: '8:15',
    thumbnail: '/assets/images/tutorials/video-2-thumb.png',
    url: 'https://example.com/videos/feishu-config'
  },
  {
    id: 'video-3',
    title: 'OpenClaw 实例部署',
    description: '一键部署 AI 助手实例',
    duration: '3:45',
    thumbnail: '/assets/images/tutorials/video-3-thumb.png',
    url: 'https://example.com/videos/deploy'
  }
];
```

### 组件设计

**HelpCenter 页面:**

```vue
<!-- frontend/src/views/help/HelpCenter.vue -->
<template>
  <div class="help-center">
    <!-- 搜索区域 -->
    <HelpSearch @search="handleSearch" />

    <!-- 搜索结果 -->
    <SearchResults
      v-if="searchQuery"
      :query="searchQuery"
      :results="searchResults"
    />

    <!-- 主内容区域 -->
    <template v-else>
      <!-- 快速入门 -->
      <section class="section quick-start" aria-labelledby="quick-start-title">
        <h2 id="quick-start-title">快速入门</h2>
        <QuickStart :items="quickStartItems" />
      </section>

      <!-- 常见问题 -->
      <section class="section faq" aria-labelledby="faq-title">
        <h2 id="faq-title">常见问题</h2>
        <HelpFaq :categories="faqCategories" />
      </section>

      <!-- 视频教程 -->
      <section class="section videos" aria-labelledby="videos-title">
        <h2 id="videos-title">视频教程</h2>
        <VideoTutorials :videos="videoTutorials" />
      </section>

      <!-- 联系客服 -->
      <section class="section contact" aria-labelledby="contact-title">
        <h2 id="contact-title">需要更多帮助？</h2>
        <ContactSupport />
      </section>
    </template>
  </div>
</template>
```

**HelpSearch 组件:**

```vue
<!-- frontend/src/components/help/HelpSearch.vue -->
<template>
  <div class="help-search">
    <n-input
      v-model:value="searchQuery"
      placeholder="搜索帮助内容..."
      size="large"
      clearable
      @update:value="debouncedSearch"
    >
      <template #prefix>
        <n-icon :component="SearchOutline" />
      </template>
    </n-input>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useDebounceFn } from '@vueuse/core';
import { SearchOutline } from '@vicons/ionicons5';

const emit = defineEmits<{
  (e: 'search', query: string): void;
}>();

const searchQuery = ref('');

const debouncedSearch = useDebounceFn((value: string) => {
  emit('search', value);
}, 300);
</script>
```

**QuickStart 组件:**

```vue
<!-- frontend/src/components/help/QuickStart.vue -->
<template>
  <div class="quick-start-grid">
    <n-card
      v-for="item in items"
      :key="item.id"
      hoverable
      class="quick-start-card"
      @click="handleClick(item)"
    >
      <div class="card-icon">
        <n-icon :component="getIcon(item.icon)" size="48" />
      </div>
      <h3>{{ item.title }}</h3>
      <p>{{ item.description }}</p>
    </n-card>
  </div>
</template>
```

**ContactSupport 组件:**

```vue
<!-- frontend/src/components/help/ContactSupport.vue -->
<template>
  <div class="contact-support">
    <n-space>
      <n-button type="primary" size="large" @click="openOnlineChat">
        <template #icon>
          <n-icon :component="ChatbubbleOutline" />
        </template>
        在线客服
      </n-button>
      <n-button size="large" @click="showEmailModal = true">
        <template #icon>
          <n-icon :component="MailOutline" />
        </template>
        发送邮件
      </n-button>
      <n-button size="large" @click="showPhoneModal = true">
        <template #icon>
          <n-icon :component="CallOutline" />
        </template>
        电话咨询
      </n-button>
    </n-space>

    <n-modal v-model:show="showEmailModal" preset="card" title="发送邮件">
      <p>客服邮箱：support@openclaw.com</p>
      <p>工作时间：周一至周五 9:00-18:00</p>
    </n-modal>

    <n-modal v-model:show="showPhoneModal" preset="card" title="电话咨询">
      <p>客服电话：400-XXX-XXXX</p>
      <p>工作时间：周一至周五 9:00-18:00</p>
    </n-modal>
  </div>
</template>
```

### 搜索功能实现

**前端搜索逻辑:**

```typescript
// frontend/src/composables/useHelpSearch.ts
import { computed, ref } from 'vue';
import { quickStartItems, faqCategories, videoTutorials } from '@/data/help';

interface SearchResult {
  type: 'quickstart' | 'faq' | 'video';
  title: string;
  description: string;
  link: string;
}

export function useHelpSearch() {
  const query = ref('');

  const results = computed<SearchResult[]>(() => {
    if (!query.value.trim()) return [];

    const lowerQuery = query.value.toLowerCase();
    const results: SearchResult[] = [];

    // 搜索快速入门
    quickStartItems.forEach(item => {
      if (
        item.title.toLowerCase().includes(lowerQuery) ||
        item.description.toLowerCase().includes(lowerQuery)
      ) {
        results.push({
          type: 'quickstart',
          title: item.title,
          description: item.description,
          link: item.link
        });
      }
    });

    // 搜索 FAQ
    faqCategories.forEach(category => {
      category.items.forEach(item => {
        if (
          item.question.toLowerCase().includes(lowerQuery) ||
          item.answer.toLowerCase().includes(lowerQuery)
        ) {
          results.push({
            type: 'faq',
            title: item.question,
            description: item.answer.slice(0, 100) + '...',
            link: `/help#faq-${item.id}`
          });
        }
      });
    });

    // 搜索视频教程
    videoTutorials.forEach(video => {
      if (
        video.title.toLowerCase().includes(lowerQuery) ||
        video.description.toLowerCase().includes(lowerQuery)
      ) {
        results.push({
          type: 'video',
          title: video.title,
          description: video.description,
          link: video.url
        });
      }
    });

    return results;
  });

  return {
    query,
    results
  };
}
```

### 项目结构规范

**新增文件位置:**

```
frontend/
├── src/
│   ├── views/
│   │   └── help/
│   │       ├── HelpCenter.vue               # 帮助中心页面（新增）
│   │       ├── FeishuPlatformTutorial.vue   # 飞书平台教程（Story 8.1）
│   │       └── FeishuBotTutorial.vue        # 飞书机器人教程（Story 8.2）
│   ├── components/
│   │   └── help/
│   │       ├── QuickStart.vue               # 快速入门组件（新增）
│   │       ├── HelpFaq.vue                  # FAQ 组件（新增）
│   │       ├── VideoTutorials.vue           # 视频教程组件（新增）
│   │       ├── VideoCard.vue                # 视频卡片组件（新增）
│   │       ├── ContactSupport.vue           # 联系客服组件（新增）
│   │       ├── HelpSearch.vue               # 搜索组件（新增）
│   │       ├── SearchResults.vue            # 搜索结果组件（新增）
│   │       ├── FaqAccordion.vue             # FAQ 折叠组件（Story 8.2）
│   │       └── FeedbackButton.vue           # 反馈按钮（Story 8.2）
│   ├── composables/
│   │   └── useHelpSearch.ts                 # 搜索逻辑（新增）
│   └── data/
│       └── help/
│           ├── quick-start.ts               # 快速入门数据（新增）
│           ├── faq.ts                       # FAQ 数据（新增）
│           └── videos.ts                    # 视频教程数据（新增）
```

### 测试标准

**测试要求:**
- 测试框架: Vitest + Vue Test Utils
- 测试覆盖率目标: ≥ 70%

**测试用例设计:**

| 测试场景 | 测试文件 | 测试方法 |
|---------|---------|---------|
| 页面渲染 | `HelpCenter.spec.ts` | 组件测试 |
| 快速入门卡片点击 | `QuickStart.spec.ts` | 组件测试 |
| FAQ 展开/收起 | `HelpFaq.spec.ts` | 组件测试 |
| 搜索功能 | `useHelpSearch.spec.ts` | 单元测试 |
| 搜索结果高亮 | `SearchResults.spec.ts` | 组件测试 |
| 联系客服弹窗 | `ContactSupport.spec.ts` | 组件测试 |

### Project Structure Notes

**与 Epic 8 其他 Story 的关系:**

1. **Story 8.1 (飞书开放平台配置教程)**: 帮助中心提供链接跳转
2. **Story 8.2 (飞书机器人添加教程)**: 帮助中心提供链接跳转
3. **Story 8.3 (首次访问引导)**: 引导中的"开始使用"可跳转到帮助中心

### 前序 Story 的学习经验

**从 Story 8.1, 8.2 获得的经验:**

1. **组件复用**: FaqAccordion 组件可直接复用
2. **内容管理**: 帮助内容使用数据驱动，便于维护
3. **无障碍实现**: 标题层级、ARIA 属性等已实现

### 常见问题与解决方案

**问题 1: 搜索性能**
- **原因**: 大量内容时前端搜索可能卡顿
- **解决**: 使用 debounce 防抖，后续可对接后端搜索 API

**问题 2: 视频加载慢**
- **原因**: 视频封面图较大
- **解决**: 使用懒加载，优化图片大小

**问题 3: FAQ 内容维护**
- **原因**: FAQ 内容可能频繁更新
- **解决**: 使用数据文件或 CMS 管理，支持热更新

### 无障碍注意事项

1. **标题层级**: 使用正确的 h1-h6 层级
2. **区域标识**: 使用 section + aria-labelledby 标识内容区域
3. **卡片组件**: 确保卡片可通过键盘操作
4. **搜索功能**: 搜索框关联 label，结果区域使用 aria-live

### References

- [Source: ux-design-specification.md] - UX 设计规范
- [Source: architecture.md#Frontend] - 前端技术栈
- [Source: epics.md#Story 8.4] - 原始故事定义
- [Source: prd.md#FR35-FR37] - 用户引导相关需求
- [Source: 8-1-feishu-platform-tutorial.md] - 教程页面参考
- [Source: 8-2-feishu-bot-add-tutorial.md] - FAQ 组件参考
- [Source: 8-3-first-visit-guide.md] - 引导内容参考

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
