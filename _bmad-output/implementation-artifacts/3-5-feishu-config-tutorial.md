# Story 3.5: 飞书开放平台配置教程组件

Status: ready-for-dev

## Story

As a 用户,
I want 查看飞书开放平台配置教程,
so that 可以了解如何创建飞书应用并获取 App ID 和 Secret。

## Acceptance Criteria

1. **AC1: 教程组件位置**
   - **Given** 用户在飞书配置页面
   - **When** 查看"如何获取配置"区域
   - **Then** 教程组件显示在配置表单下方
   - **And** 默认收起状态
   - **And** 显示标题"如何获取飞书应用配置"

2. **AC2: 分步教程结构**
   - **Given** 用户展开教程
   - **When** 查看教程内容
   - **Then** 显示 4 个步骤
   - **And** 每个步骤有编号、标题和说明
   - **And** 当前步骤高亮显示

3. **AC3: 步骤 1 - 跳转到飞书开放平台**
   - **Given** 用户查看步骤 1
   - **When** 阅读步骤内容
   - **Then** 显示说明"访问飞书开放平台"
   - **And** 提供跳转链接按钮（https://open.feishu.cn）
   - **And** 点击链接在新窗口打开

4. **AC4: 步骤 2 - 创建企业自建应用截图指引**
   - **Given** 用户查看步骤 2
   - **When** 阅读步骤内容
   - **Then** 显示操作说明
   - **And** 显示创建应用界面截图
   - **And** 关键位置有红色边框标注
   - **And** 支持点击放大查看

5. **AC5: 步骤 3 - 获取 App ID/Secret 位置截图**
   - **Given** 用户查看步骤 3
   - **When** 阅读步骤内容
   - **Then** 显示凭证位置截图
   - **And** App ID 和 App Secret 位置高亮标注
   - **And** 提供复制按钮（点击复制示例值）
   - **And** 提示保密 App Secret

6. **AC6: 步骤 4 - 返回平台填写指引**
   - **Given** 用户查看步骤 4
   - **When** 阅读步骤内容
   - **Then** 显示"返回本页面填写配置"说明
   - **And** 提供"返回顶部"快捷按钮
   - **And** 显示完成提示

7. **AC7: 展开/收起功能**
   - **Given** 教练组件已渲染
   - **When** 点击展开/收起按钮
   - **Then** 切换教程内容显示状态
   - **And** 使用动画过渡效果
   - **And** 展开时图标旋转 180 度

8. **AC8: 截图展开/收起功能**
   - **Given** 用户查看带截图的步骤
   - **When** 点击截图区域
   - **Then** 切换截图显示/隐藏状态
   - **And** 显示时加载高清截图
   - **And** 隐藏时显示缩略图或占位符

## Tasks / Subtasks

- [ ] Task 1: 创建教程组件基础结构 (AC: 1, 7)
  - [ ] 1.1 创建 `src/components/config/FeishuConfigTutorial.vue` 组件
  - [ ] 1.2 实现展开/收起状态管理
  - [ ] 1.3 实现动画过渡效果
  - [ ] 1.4 添加标题和展开按钮

- [ ] Task 2: 创建步骤组件 (AC: 2)
  - [ ] 2.1 创建 `src/components/config/TutorialStep.vue` 步骤组件
  - [ ] 2.2 实现步骤编号、标题、内容布局
  - [ ] 2.3 实现当前步骤高亮样式
  - [ ] 2.4 支持插槽显示自定义内容

- [ ] Task 3: 实现步骤 1 - 飞书开放平台链接 (AC: 3)
  - [ ] 3.1 创建步骤 1 内容
  - [ ] 3.2 添加飞书开放平台链接按钮
  - [ ] 3.3 设置链接在新窗口打开

- [ ] Task 4: 实现步骤 2 - 创建应用截图 (AC: 4)
  - [ ] 4.1 准备创建应用截图素材
  - [ ] 4.2 创建截图标注组件
  - [ ] 4.3 实现点击放大功能
  - [ ] 4.4 添加操作说明文字

- [ ] Task 5: 实现步骤 3 - 获取凭证截图 (AC: 5)
  - [ ] 5.1 准备凭证位置截图素材
  - [ ] 5.2 添加 App ID/Secret 高亮标注
  - [ ] 5.3 实现复制示例值功能
  - [ ] 5.4 添加保密提示

- [ ] Task 6: 实现步骤 4 - 返回填写指引 (AC: 6)
  - [ ] 6.1 创建步骤 4 内容
  - [ ] 6.2 实现"返回顶部"功能
  - [ ] 6.3 添加完成提示

- [ ] Task 7: 实现截图展开/收起功能 (AC: 8)
  - [ ] 7.1 创建截图容器组件
  - [ ] 7.2 实现懒加载高清截图
  - [ ] 7.3 实现缩略图/占位符显示
  - [ ] 7.4 添加展开/收起动画

- [ ] Task 8: 准备教程截图素材 (AC: 4, 5)
  - [ ] 8.1 截取飞书开放平台首页截图
  - [ ] 8.2 截取创建应用界面截图
  - [ ] 8.3 截取凭证位置截图
  - [ ] 8.4 添加标注和说明

- [ ] Task 9: 编写组件测试 (AC: 1-8)
  - [ ] 9.1 编写教程组件测试
  - [ ] 9.2 编写步骤组件测试
  - [ ] 9.3 测试展开/收起功能
  - [ ] 9.4 测试链接跳转功能

## Dev Notes

### 架构模式与约束

**必须遵循的前端架构 [Source: architecture.md]:**

1. **组件组织**:
   - 按功能组织组件（`src/components/config/`）
   - 组件可复用、可测试

2. **命名约定 [Source: architecture.md#Naming Patterns]:**
   - 组件名: PascalCase (例: `FeishuConfigTutorial`)
   - 文件名: kebab-case (例: `feishu-config-tutorial.vue`)

3. **UI 组件库 [Source: architecture.md]:**
   - 使用 Naive UI 组件库
   - 遵循 Naive UI 设计规范

4. **UX 模式 [Source: ux-design-specification.md]:**
   - 飞书配置引导组件：分步图文教程
   - 支持展开/收起截图

### 现有项目状态

**Story 3.4 已完成的页面结构:**

```
frontend/
├── src/
│   ├── components/
│   │   └── config/
│   │       ├── FeishuConfigForm.vue   # ✅ 已创建
│   │       ├── ConfigStatusCard.vue   # ✅ 已创建
│   │       └── FeishuConfigTutorial.vue # 待创建
│   └── pages/
│       └── feishu-config/
│           └── index.vue              # ✅ 已创建
```

### 技术栈要求

**核心依赖:**

| 依赖 | 用途 | 版本 |
|------|------|------|
| Vue | 前端框架 | 3.x |
| TypeScript | 类型支持 | 5.x |
| Naive UI | 组件库 | 2.x |

### 教程组件实现

**飞书配置教程组件:**

```vue
<!-- src/components/config/FeishuConfigTutorial.vue -->
<template>
  <n-card class="tutorial-card">
    <template #header>
      <div class="tutorial-header" @click="toggleExpand">
        <n-space align="center">
          <n-icon :component="HelpCircleOutline" />
          <span class="title">如何获取飞书应用配置</span>
        </n-space>
        <n-icon :component="ChevronDownOutline" :class="{ 'expanded': isExpanded }" />
      </div>
    </template>

    <n-collapse-transition :show="isExpanded">
      <n-space vertical size="large">
        <!-- 步骤 1: 访问飞书开放平台 -->
        <TutorialStep
          :step="1"
          title="访问飞书开放平台"
        >
          <p>打开浏览器，访问飞书开放平台官网</p>
          <n-button
            type="primary"
            tag="a"
            href="https://open.feishu.cn"
            target="_blank"
          >
            <template #icon>
              <n-icon :component="OpenOutline" />
            </template>
            打开飞书开放平台
          </n-button>
        </TutorialStep>

        <!-- 步骤 2: 创建企业自建应用 -->
        <TutorialStep
          :step="2"
          title="创建企业自建应用"
        >
          <p>登录后，点击"创建企业自建应用"按钮</p>
          <ScreenshotViewer
            src="/images/tutorial/feishu-create-app.png"
            alt="创建应用界面截图"
            annotation="点击红色框标注的按钮创建应用"
          />
        </TutorialStep>

        <!-- 步骤 3: 获取 App ID 和 Secret -->
        <TutorialStep
          :step="3"
          title="获取 App ID 和 App Secret"
        >
          <p>在应用详情页，找到"凭证与基础信息"区域</p>
          <ScreenshotViewer
            src="/images/tutorial/feishu-credentials.png"
            alt="凭证位置截图"
            :annotations="[
              { x: 20, y: 30, text: 'App ID 在此处', highlight: true },
              { x: 20, y: 60, text: 'App Secret 在此处', highlight: true }
            ]"
          />
          <n-alert type="warning" title="安全提示">
            App Secret 是敏感信息，请勿泄露给他人
          </n-alert>
        </TutorialStep>

        <!-- 步骤 4: 返回平台填写 -->
        <TutorialStep
          :step="4"
          title="返回平台填写配置"
        >
          <p>复制 App ID 和 App Secret，返回本页面填写配置</p>
          <n-space>
            <n-button @click="scrollToTop">
              <template #icon>
                <n-icon :component="ArrowUpOutline" />
              </template>
              返回顶部
            </n-button>
          </n-space>
        </TutorialStep>
      </n-space>
    </n-collapse-transition>
  </n-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import {
  NCard, NSpace, NButton, NIcon, NAlert, NCollapseTransition
} from 'naive-ui'
import {
  HelpCircleOutline,
  ChevronDownOutline,
  OpenOutline,
  ArrowUpOutline
} from '@vicons/ionicons5'
import TutorialStep from './TutorialStep.vue'
import ScreenshotViewer from './ScreenshotViewer.vue'

const isExpanded = ref(false)

const toggleExpand = () => {
  isExpanded.value = !isExpanded.value
}

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}
</script>

<style scoped>
.tutorial-card {
  margin-top: 16px;
}

.tutorial-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
  user-select: none;
}

.tutorial-header .title {
  font-weight: 500;
  font-size: 16px;
}

.tutorial-header .expanded {
  transform: rotate(180deg);
  transition: transform 0.3s ease;
}
</style>
```

### 步骤组件实现

**教程步骤组件:**

```vue
<!-- src/components/config/TutorialStep.vue -->
<template>
  <div class="tutorial-step">
    <div class="step-header">
      <n-badge :value="step" :max="99" color="#1677FF" />
      <h4 class="step-title">{{ title }}</h4>
    </div>
    <div class="step-content">
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { NBadge } from 'naive-ui'

interface Props {
  step: number
  title: string
}

defineProps<Props>()
</script>

<style scoped>
.tutorial-step {
  padding: 16px;
  border-radius: 8px;
  background-color: #fafafa;
}

.step-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.step-title {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
}

.step-content {
  padding-left: 32px;
}

.step-content p {
  margin-bottom: 12px;
  color: #666;
}
</style>
```

### 截图查看器组件

**截图展开/收起组件:**

```vue
<!-- src/components/config/ScreenshotViewer.vue -->
<template>
  <div class="screenshot-viewer">
    <div class="screenshot-header" @click="toggleScreenshot">
      <n-space align="center">
        <n-icon :component="ImageOutline" />
        <span>查看截图指引</span>
      </n-space>
      <n-icon :component="ChevronDownOutline" :class="{ 'expanded': showScreenshot }" />
    </div>

    <n-collapse-transition :show="showScreenshot">
      <div class="screenshot-content">
        <div class="screenshot-wrapper" @click="openModal = true">
          <img
            :src="src"
            :alt="alt"
            loading="lazy"
            class="screenshot-image"
          />
          <div v-if="annotation" class="annotation">
            {{ annotation }}
          </div>
        </div>
      </div>
    </n-collapse-transition>

    <!-- 图片放大模态框 -->
    <n-modal v-model:show="openModal" preset="card" style="width: 90%; max-width: 1200px;">
      <img :src="src" :alt="alt" style="width: 100%;" />
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import {
  NSpace, NIcon, NCollapseTransition, NModal
} from 'naive-ui'
import { ImageOutline, ChevronDownOutline } from '@vicons/ionicons5'

interface Props {
  src: string
  alt: string
  annotation?: string
  annotations?: Array<{ x: number; y: number; text: string; highlight?: boolean }>
}

defineProps<Props>()

const showScreenshot = ref(false)
const openModal = ref(false)

const toggleScreenshot = () => {
  showScreenshot.value = !showScreenshot.value
}
</script>

<style scoped>
.screenshot-viewer {
  margin-top: 12px;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  overflow: hidden;
}

.screenshot-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background-color: #f5f5f5;
  cursor: pointer;
  user-select: none;
}

.screenshot-header:hover {
  background-color: #ebebeb;
}

.screenshot-content {
  padding: 16px;
}

.screenshot-wrapper {
  position: relative;
  cursor: pointer;
}

.screenshot-image {
  width: 100%;
  border-radius: 4px;
  border: 1px solid #e8e8e8;
}

.screenshot-wrapper:hover .screenshot-image {
  border-color: #1677FF;
}

.annotation {
  position: absolute;
  bottom: 16px;
  left: 16px;
  right: 16px;
  padding: 8px 12px;
  background-color: rgba(0, 0, 0, 0.7);
  color: white;
  border-radius: 4px;
  font-size: 14px;
}
</style>
```

### 教程截图素材

**截图清单:**

| 文件名 | 用途 | 说明 |
|-------|------|------|
| `feishu-homepage.png` | 飞书开放平台首页 | 展示入口位置 |
| `feishu-create-app.png` | 创建应用界面 | 标注"创建企业自建应用"按钮 |
| `feishu-credentials.png` | 凭证位置截图 | 标注 App ID 和 App Secret 位置 |

**截图要求:**
- 分辨率: 至少 1280x720
- 格式: PNG（保持清晰度）
- 标注: 红色边框，2px 宽度
- 存储: `public/images/tutorial/`

### UX 设计规范

**来自 UX 设计文档的要求:**

| 元素 | 规范 |
|------|------|
| 步骤编号 | 蓝色徽章 #1677FF |
| 步骤标题 | 16px，中等字重 |
| 展开按钮 | 右侧对齐，点击区域 44x44px |
| 截图边框 | 1px 灰色边框，悬停蓝色 |
| 安全提示 | 黄色警告框 |
| 动画时长 | 0.3s |

### 测试标准

**测试要求:**
- 测试框架: Vitest + Vue Test Utils
- 测试覆盖率目标: >= 70%

**测试用例设计:**

| 测试场景 | 测试方法 | 预期结果 |
|---------|---------|---------|
| 组件渲染 | mount 组件 | 显示标题和展开按钮 |
| 默认收起 | 初始状态 | 内容不显示 |
| 展开教程 | 点击标题 | 内容展开显示 |
| 收起教程 | 再次点击 | 内容收起隐藏 |
| 步骤渲染 | 展开后 | 显示 4 个步骤 |
| 链接跳转 | 点击链接按钮 | 新窗口打开链接 |
| 截图展开 | 点击截图区域 | 截图显示/隐藏 |
| 图片放大 | 点击截图 | 模态框显示大图 |

### 无障碍要求

**WCAG 2.1 AA 级合规:**

1. **键盘导航**:
   - 展开/收起按钮支持键盘操作
   - Tab 键可以在步骤间导航

2. **屏幕阅读器**:
   - 步骤编号和标题使用语义化标签
   - 图片提供 alt 属性

3. **对比度**:
   - 文字与背景对比度 >= 4.5:1
   - 标注颜色醒目

### 前序 Story 的学习经验

**从 Story 3.4 (飞书配置前端页面) 获得的经验:**

1. **组件组织**: 按功能模块组织组件
2. **状态管理**: 使用 ref 管理组件内部状态
3. **样式规范**: 使用 scoped CSS
4. **Naive UI**: 使用 CollapseTransition 实现动画

### References

- [Source: architecture.md#Frontend Architecture] - 前端架构
- [Source: ux-design-specification.md#Custom Components] - 飞书配置引导组件
- [Source: prd.md#FR35] - 飞书开放平台配置教程
- [Source: epics.md#Story 3.5] - 原始故事定义
- [Source: 3-4-feishu-config-frontend.md] - 前端页面实现参考

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
