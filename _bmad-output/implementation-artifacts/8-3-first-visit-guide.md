# Story 8.3: 首次访问引导

Status: ready-for-dev

## Story

As a 新用户,
I want 在首次访问时看到引导说明,
so that 快速了解产品功能和使用流程。

## Acceptance Criteria

1. **AC1: 首次登录检测**
   - **Given** 用户首次登录系统
   - **When** 进入系统首页
   - **Then** 系统检测到用户首次访问
   - **And** 查询用户偏好设置中是否有引导记录

2. **AC2: 欢迎引导弹窗显示**
   - **Given** 用户首次登录系统
   - **When** 进入系统首页
   - **Then** 显示欢迎引导弹窗
   - **And** 弹窗居中显示，带有遮罩层
   - **And** 弹窗可关闭

3. **AC3: 引导内容完整性**
   - **Given** 欢迎引导弹窗已显示
   - **When** 查看引导内容
   - **Then** 引导包含：产品介绍、核心功能、使用流程
   - **And** 使用分步轮播或卡片形式展示
   - **And** 每部分有清晰的标题和说明

4. **AC4: 操作按钮功能**
   - **Given** 用户在引导弹窗中
   - **When** 点击"开始使用"按钮
   - **Then** 关闭弹窗并跳转到配置页
   - **When** 点击"稍后再看"按钮
   - **Then** 关闭引导弹窗

5. **AC5: 不再显示选项**
   - **Given** 用户在引导弹窗中
   - **When** 勾选"不再显示"选项
   - **And** 关闭弹窗
   - **Then** 用户偏好中记录引导已完成
   - **And** 后续登录不再显示引导

6. **AC6: 引导记录存储**
   - **Given** 用户完成或跳过引导
   - **When** 系统保存引导状态
   - **Then** 引导记录存储在用户偏好中
   - **And** API 保存用户偏好设置
   - **And** 前端缓存用户偏好状态

7. **AC7: 无障碍访问**
   - **Given** 用户使用屏幕阅读器
   - **When** 引导弹窗显示
   - **Then** 弹窗内容可被屏幕阅读器读取
   - **And** 支持键盘导航（Tab、Enter、Escape）
   - **And** 焦点正确管理

## Tasks / Subtasks

- [ ] Task 1: 创建用户偏好数据模型和 API (AC: 1, 6)
  - [ ] 1.1 创建 `backend/internal/domain/user/preference.go` 用户偏好模型
  - [ ] 1.2 创建数据库迁移脚本添加 `user_preferences` 表
  - [ ] 1.3 创建 `backend/internal/repository/user_preference.go`
  - [ ] 1.4 创建 API 端点 `GET/PUT /v1/user/preferences`
  - [ ] 1.5 创建前端 API 客户端方法

- [ ] Task 2: 创建引导弹窗组件 (AC: 2, 3)
  - [ ] 2.1 创建 `frontend/src/components/guide/WelcomeGuide.vue`
  - [ ] 2.2 创建 `frontend/src/components/guide/GuideStep.vue` 步骤组件
  - [ ] 2.3 实现弹窗样式和动画效果
  - [ ] 2.4 配置弹窗遮罩层

- [ ] Task 3: 实现引导内容 (AC: 3)
  - [ ] 3.1 设计引导内容数据结构
  - [ ] 3.2 编写产品介绍内容
  - [ ] 3.3 编写核心功能介绍
  - [ ] 3.4 编写使用流程步骤
  - [ ] 3.5 准备引导图片/图标资源

- [ ] Task 4: 实现操作按钮 (AC: 4, 5)
  - [ ] 4.1 实现"开始使用"按钮跳转逻辑
  - [ ] 4.2 实现"稍后再看"按钮关闭逻辑
  - [ ] 4.3 创建"不再显示"复选框组件
  - [ ] 4.4 实现偏好设置保存逻辑

- [ ] Task 5: 创建引导状态管理 (AC: 1, 6)
  - [ ] 5.1 创建 `frontend/src/stores/user-preferences.ts`
  - [ ] 5.2 实现首次访问检测逻辑
  - [ ] 5.3 实现引导状态缓存

- [ ] Task 6: 无障碍优化 (AC: 7)
  - [ ] 6.1 添加弹窗 ARIA 属性
  - [ ] 6.2 实现焦点陷阱（Focus Trap）
  - [ ] 6.3 实现键盘导航支持
  - [ ] 6.4 测试屏幕阅读器兼容性

- [ ] Task 7: 编写测试 (AC: 1-7)
  - [ ] 7.1 编写后端用户偏好 API 测试
  - [ ] 7.2 编写前端引导弹窗组件测试
  - [ ] 7.3 编写引导状态管理测试
  - [ ] 7.4 编写首次访问检测集成测试

## Dev Notes

### 架构模式与约束

**必须遵循的 UX 设计规范 [Source: ux-design-specification.md]:**

1. **布局策略**: 向导式聚焦布局（分步引导）
2. **无障碍要求**: WCAG 2.1 AA 级合规
3. **最小触控区域**: 44x44px
4. **色彩对比度**: 4.5:1

**UX 模式 [Source: ux-design-specification.md]:**
- 主要按钮：蓝色填充 #1677FF
- 成功反馈：Message.success 3秒自动消失

### 现有项目状态

**前端项目结构 [Source: 1-2-frontend-project-init.md]:**

```
frontend/
├── src/
│   ├── views/              # 页面组件
│   ├── components/         # 公共组件
│   ├── stores/             # Pinia 状态管理
│   ├── api/                # API 客户端
│   └── router/             # 路由配置
```

**后端项目结构 [Source: 1-1-backend-project-init.md]:**

```
backend/
├── cmd/server/main.go
├── internal/
│   ├── api/               # REST API 处理器
│   ├── domain/            # 领域模型
│   └── repository/        # 数据访问层
```

### 引导内容设计

**引导步骤设计:**

```typescript
interface GuideStep {
  id: string;
  title: string;
  description: string;
  icon: string;
  image?: string;
}

const guideSteps: GuideStep[] = [
  {
    id: 'welcome',
    title: '欢迎使用 OpenClaw',
    description: 'OpenClaw 是一个智能 AI 助手平台，让您可以快速部署和使用 AI Agent，与飞书无缝集成。',
    icon: 'rocket'
  },
  {
    id: 'features',
    title: '核心功能',
    description: '一键部署 AI 助手实例，配置飞书机器人，实现智能对话。无需编程知识，3分钟即可完成。',
    icon: 'star',
    features: [
      { text: '一键部署 OpenClaw 实例', icon: 'cloud' },
      { text: '飞书机器人集成', icon: 'chat' },
      { text: '智能对话能力', icon: 'brain' }
    ]
  },
  {
    id: 'workflow',
    title: '使用流程',
    description: '只需三步，即可开始使用您的专属 AI 助手。',
    icon: 'flow',
    steps: [
      '配置飞书应用',
      '启动 OpenClaw 实例',
      '在飞书中开始对话'
    ]
  }
];
```

### 数据库设计

**user_preferences 表:**

```sql
CREATE TABLE user_preferences (
    "ID" VARCHAR(36) PRIMARY KEY,
    "UserID" VARCHAR(36) NOT NULL REFERENCES tenant_users("ID") ON DELETE CASCADE,
    "GuideCompleted" BOOLEAN NOT NULL DEFAULT FALSE,
    "GuideCompletedAt" TIMESTAMP,
    "Theme" VARCHAR(50) DEFAULT 'light',
    "Language" VARCHAR(10) DEFAULT 'zh-CN',
    "NotificationEnabled" BOOLEAN DEFAULT TRUE,
    "CreatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "UpdatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_user_preferences_user_id ON user_preferences ("UserID");
```

### API 设计

**获取用户偏好:**

```
GET /v1/user/preferences

Response:
{
  "data": {
    "guideCompleted": false,
    "guideCompletedAt": null,
    "theme": "light",
    "language": "zh-CN",
    "notificationEnabled": true
  },
  "error": null,
  "meta": {}
}
```

**更新用户偏好:**

```
PUT /v1/user/preferences

Request:
{
  "guideCompleted": true
}

Response:
{
  "data": {
    "guideCompleted": true,
    "guideCompletedAt": "2026-03-05T10:30:00Z",
    ...
  },
  "error": null,
  "meta": {}
}
```

### 组件设计

**WelcomeGuide 组件:**

```vue
<!-- frontend/src/components/guide/WelcomeGuide.vue -->
<template>
  <n-modal
    v-model:show="visible"
    :mask-closable="false"
    :closable="true"
    preset="card"
    style="width: 600px; max-width: 90vw;"
    title="欢迎使用 OpenClaw"
    role="dialog"
    aria-modal="true"
    aria-labelledby="guide-title"
  >
    <div class="welcome-guide">
      <!-- 步骤指示器 -->
      <div class="step-indicators">
        <span
          v-for="(_, index) in steps"
          :key="index"
          :class="['indicator', { active: currentStep === index }]"
        />
      </div>

      <!-- 当前步骤内容 -->
      <GuideStep
        :step="steps[currentStep]"
        :step-number="currentStep + 1"
        :total-steps="steps.length"
      />

      <!-- 底部操作区 -->
      <div class="guide-actions">
        <n-checkbox v-model:checked="dontShowAgain">
          不再显示
        </n-checkbox>

        <div class="action-buttons">
          <n-button
            v-if="currentStep > 0"
            @click="prevStep"
          >
            上一步
          </n-button>
          <n-button
            v-if="currentStep < steps.length - 1"
            type="primary"
            @click="nextStep"
          >
            下一步
          </n-button>
          <n-button
            v-else
            type="primary"
            @click="handleStart"
          >
            开始使用
          </n-button>
        </div>
      </div>
    </div>

    <template #footer>
      <n-button text @click="handleSkip">稍后再看</n-button>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useUserPreferencesStore } from '@/stores/user-preferences';
import GuideStep from './GuideStep.vue';

const router = useRouter();
const preferencesStore = useUserPreferencesStore();

const visible = ref(true);
const currentStep = ref(0);
const dontShowAgain = ref(false);

const steps = [
  { id: 'welcome', title: '欢迎使用 OpenClaw', ... },
  { id: 'features', title: '核心功能', ... },
  { id: 'workflow', title: '使用流程', ... }
];

function nextStep() {
  if (currentStep.value < steps.length - 1) {
    currentStep.value++;
  }
}

function prevStep() {
  if (currentStep.value > 0) {
    currentStep.value--;
  }
}

async function handleStart() {
  if (dontShowAgain.value) {
    await preferencesStore.markGuideCompleted();
  }
  visible.value = false;
  router.push('/feishu-config');
}

async function handleSkip() {
  if (dontShowAgain.value) {
    await preferencesStore.markGuideCompleted();
  }
  visible.value = false;
}
</script>
```

### 状态管理设计

**用户偏好 Store:**

```typescript
// frontend/src/stores/user-preferences.ts
import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { userApi } from '@/api/user';

export const useUserPreferencesStore = defineStore('userPreferences', () => {
  const preferences = ref<UserPreferences | null>(null);
  const loading = ref(false);

  const shouldShowGuide = computed(() => {
    return preferences.value?.guideCompleted === false;
  });

  async function fetchPreferences() {
    loading.value = true;
    try {
      const response = await userApi.getPreferences();
      preferences.value = response.data;
    } catch (error) {
      console.error('Failed to fetch preferences:', error);
    } finally {
      loading.value = false;
    }
  }

  async function markGuideCompleted() {
    try {
      const response = await userApi.updatePreferences({
        guideCompleted: true
      });
      preferences.value = response.data;
    } catch (error) {
      console.error('Failed to update preferences:', error);
    }
  }

  return {
    preferences,
    loading,
    shouldShowGuide,
    fetchPreferences,
    markGuideCompleted
  };
});
```

### 项目结构规范

**新增文件位置:**

```
frontend/
├── src/
│   ├── components/
│   │   └── guide/
│   │       ├── WelcomeGuide.vue            # 欢迎引导弹窗（新增）
│   │       └── GuideStep.vue               # 引导步骤组件（新增）
│   ├── stores/
│   │   └── user-preferences.ts             # 用户偏好状态（新增）
│   ├── api/
│   │   └── user.ts                         # 用户 API（扩展）
│   └── data/
│       └── guide/
│           └── guide-steps.ts              # 引导内容数据（新增）

backend/
├── internal/
│   ├── domain/
│   │   └── user/
│   │       └── preference.go               # 用户偏好模型（新增）
│   ├── repository/
│   │   └── user_preference.go             # 用户偏好仓库（新增）
│   └── api/
│       └── handlers/
│           └── user_preferences.go        # 用户偏好 API（新增）
├── migrations/
│   └── 004_create_user_preferences.up.sql # 数据库迁移（新增）
```

### 无障碍实现要点

**焦点管理:**

```typescript
// 焦点陷阱实现
import { useFocusTrap } from '@vueuse/integrations/useFocusTrap';

const modalRef = ref<HTMLElement>();
const { activate, deactivate } = useFocusTrap(modalRef);

watch(visible, (newValue) => {
  if (newValue) {
    activate();
  } else {
    deactivate();
  }
});
```

**键盘导航:**

```typescript
// 键盘事件处理
function handleKeyDown(event: KeyboardEvent) {
  switch (event.key) {
    case 'Escape':
      handleSkip();
      break;
    case 'ArrowRight':
      nextStep();
      break;
    case 'ArrowLeft':
      prevStep();
      break;
  }
}
```

### 测试标准

**测试要求:**
- 前端测试框架: Vitest + Vue Test Utils
- 后端测试框架: Go testing + testify
- 测试覆盖率目标: ≥ 70%

**测试用例设计:**

| 测试场景 | 测试文件 | 测试方法 |
|---------|---------|---------|
| 用户偏好 API | `user_preferences_test.go` | 集成测试 |
| 引导弹窗渲染 | `WelcomeGuide.spec.ts` | 组件测试 |
| 步骤导航 | `WelcomeGuide.spec.ts` | 组件测试 |
| 偏好保存 | `user-preferences.spec.ts` | Store 测试 |
| 首次访问检测 | `App.spec.ts` | 集成测试 |
| 键盘导航 | `WelcomeGuide.spec.ts` | 组件测试 |

### Project Structure Notes

**与 Epic 8 其他 Story 的关系:**

1. **Story 8.4 (帮助中心)**: 引导中的"开始使用"按钮可跳转到配置页或帮助中心

### 前序 Story 的学习经验

**从 Story 1.2 (前端项目初始化) 获得的经验:**

1. **状态管理**: 使用 Pinia 按模块组织
2. **组件规范**: Vue 3 Composition API + TypeScript

**从 Story 2.1 (用户数据模型) 获得的经验:**

1. **数据库迁移**: 使用 golang-migrate 工具
2. **Repository 模式**: Clean Architecture 分层

### 常见问题与解决方案

**问题 1: 引导弹窗重复显示**
- **原因**: 用户偏好未正确保存或加载
- **解决**: 确保前端和后端状态同步，使用缓存减少 API 调用

**问题 2: 焦点管理问题**
- **原因**: 弹窗关闭后焦点未正确恢复
- **解决**: 使用 useFocusTrap 或手动管理焦点

### 无障碍注意事项

1. **弹窗角色**: 使用 role="dialog" 和 aria-modal="true"
2. **标题关联**: 使用 aria-labelledby 关联标题
3. **焦点陷阱**: 弹窗打开时焦点限制在弹窗内
4. **键盘支持**: Escape 关闭、Tab 导航
5. **屏幕阅读器**: 确保内容可被正确读取

### References

- [Source: ux-design-specification.md] - UX 设计规范
- [Source: architecture.md#Frontend] - 前端技术栈
- [Source: architecture.md#Backend] - 后端技术栈
- [Source: epics.md#Story 8.3] - 原始故事定义
- [Source: prd.md#FR37] - 首次访问引导需求
- [Source: 1-1-backend-project-init.md] - 后端项目实现经验
- [Source: 1-2-frontend-project-init.md] - 前端项目实现经验
- [Source: 2-1-user-data-model.md] - 用户数据模型实现经验

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
