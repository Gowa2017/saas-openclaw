# Story 8.2: 飞书机器人添加教程页面

Status: ready-for-dev

## Story

As a 用户,
I want 查看如何在飞书中添加 OpenClaw 机器人的教程,
so that 可以在飞书中开始使用。

## Acceptance Criteria

1. **AC1: 教程页面路由和入口**
   - **Given** 用户访问帮助中心
   - **When** 点击"飞书机器人添加教程"
   - **Then** 跳转到教程页面 `/help/feishu-bot-tutorial`
   - **And** 页面标题显示"飞书机器人添加教程"

2. **AC2: 教程步骤展示**
   - **Given** 用户进入教程页面
   - **When** 页面加载完成
   - **Then** 显示完整的图文教程
   - **And** 教程包含 3 个步骤
   - **And** 每步包含截图和文字说明

3. **AC3: 步骤内容完整性**
   - **Given** 教程页面已加载
   - **When** 查看教程内容
   - **Then** 步骤 1：在飞书搜索机器人
   - **And** 步骤 2：添加到单聊/群聊
   - **And** 步骤 3：发送消息测试

4. **AC4: 常见问题解答（FAQ）**
   - **Given** 用户在教程页面
   - **When** 查看 FAQ 区域
   - **Then** 显示常见问题列表
   - **And** 点击问题可展开/收起答案
   - **And** FAQ 内容覆盖常见使用场景

5. **AC5: 问题反馈入口**
   - **Given** 用户遇到问题
   - **When** 点击"遇到问题"按钮
   - **Then** 显示问题反馈表单或跳转客服
   - **And** 表单包含问题描述输入框
   - **And** 提供提交功能

6. **AC6: 无障碍访问**
   - **Given** 用户使用屏幕阅读器
   - **When** 访问教程页面
   - **Then** 图片有 alt 文字描述
   - **And** FAQ 可通过键盘操作
   - **And** 反馈表单可通过键盘填写

## Tasks / Subtasks

- [ ] Task 1: 创建教程页面路由 (AC: 1)
  - [ ] 1.1 在 `frontend/src/router/modules/help.ts` 添加路由
  - [ ] 1.2 创建 `frontend/src/views/help/FeishuBotTutorial.vue`
  - [ ] 1.3 配置页面 meta 信息（标题、描述）

- [ ] Task 2: 复用并扩展教程组件 (AC: 2, 3)
  - [ ] 2.1 复用 Story 8.1 的 TutorialStep 组件
  - [ ] 2.2 复用 Story 8.1 的 TutorialImage 组件
  - [ ] 2.3 创建教程内容数据

- [ ] Task 3: 实现教程内容 (AC: 3)
  - [ ] 3.1 准备飞书机器人相关截图资源
  - [ ] 3.2 编写步骤 1：搜索机器人的内容和截图
  - [ ] 3.3 编写步骤 2：添加到会话的内容和截图
  - [ ] 3.4 编写步骤 3：发送测试的内容和截图

- [ ] Task 4: 实现 FAQ 组件 (AC: 4)
  - [ ] 4.1 创建 `frontend/src/components/help/FaqAccordion.vue`
  - [ ] 4.2 定义 FAQ 数据结构
  - [ ] 4.3 实现展开/收起动画效果
  - [ ] 4.4 编写 FAQ 内容

- [ ] Task 5: 实现问题反馈入口 (AC: 5)
  - [ ] 5.1 创建 `frontend/src/components/help/FeedbackButton.vue`
  - [ ] 5.2 创建反馈表单弹窗组件
  - [ ] 5.3 实现表单提交功能
  - [ ] 5.4 添加成功/失败提示

- [ ] Task 6: 无障碍优化 (AC: 6)
  - [ ] 6.1 为所有图片添加 alt 属性
  - [ ] 6.2 确保 FAQ 组件键盘可访问
  - [ ] 6.3 确保反馈表单键盘可操作
  - [ ] 6.4 添加 ARIA 属性支持

- [ ] Task 7: 编写测试 (AC: 1-6)
  - [ ] 7.1 编写页面渲染测试
  - [ ] 7.2 编写步骤展示测试
  - [ ] 7.3 编写 FAQ 组件测试
  - [ ] 7.4 编写反馈功能测试

## Dev Notes

### 架构模式与约束

**必须遵循的 UX 设计规范 [Source: ux-design-specification.md]:**

1. **布局策略**: 向导式聚焦布局（分步引导）
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

**可复用组件 (Story 8.1):**
- `TutorialStep.vue` - 步骤展示组件
- `TutorialImage.vue` - 图片展示组件

### 教程内容设计

**步骤 1: 在飞书搜索机器人**

```
标题: 在飞书搜索机器人
内容:
1. 打开飞书客户端
2. 点击左上角搜索框
3. 输入机器人名称（如：OpenClaw 助手）
4. 在搜索结果中找到机器人

截图:
- 飞书客户端首页截图
- 搜索框位置截图
- 搜索结果截图
```

**步骤 2: 添加到单聊/群聊**

```
标题: 添加到单聊/群聊
内容:
单聊方式:
1. 点击机器人头像
2. 点击"添加到通讯录"
3. 发送消息开始对话

群聊方式:
1. 进入目标群聊
2. 点击群设置 -> 群机器人
3. 点击"添加机器人"
4. 选择 OpenClaw 助手

截图:
- 机器人详情页截图
- 添加到通讯录截图
- 群设置入口截图
- 添加机器人截图
```

**步骤 3: 发送消息测试**

```
标题: 发送消息测试
内容:
1. 在对话中输入任意消息
2. 等待机器人回复
3. 确认机器人正常工作
4. 如果没有回复，请检查:
   - 应用配置是否正确
   - 机器人服务是否启动

截图:
- 发送消息截图
- 收到回复截图
```

### FAQ 内容设计

**FAQ 数据结构:**

```typescript
interface FaqItem {
  id: string;
  question: string;
  answer: string;
  category?: string;
}

const faqItems: FaqItem[] = [
  {
    id: 'faq-1',
    question: '搜索不到机器人怎么办？',
    answer: '请确认：1) 管理员已发布应用；2) 你在同一企业内；3) 搜索名称正确。',
    category: '搜索问题'
  },
  {
    id: 'faq-2',
    question: '机器人没有回复怎么办？',
    answer: '请检查：1) 应用配置是否正确；2) OpenClaw 实例是否运行中；3) 查看实例日志排查问题。',
    category: '消息问题'
  },
  {
    id: 'faq-3',
    question: '如何在群里使用机器人？',
    answer: '进入群设置 -> 群机器人 -> 添加机器人 -> 选择 OpenClaw 助手即可。',
    category: '群聊使用'
  },
  {
    id: 'faq-4',
    question: '机器人权限不足怎么办？',
    answer: '请联系应用管理员在飞书开放平台配置相应权限，然后重新发布应用。',
    category: '权限问题'
  },
  {
    id: 'faq-5',
    question: '支持哪些消息类型？',
    answer: '目前支持文本消息、富文本消息和卡片消息。图片和文件消息暂不支持。',
    category: '消息类型'
  }
];
```

### 组件设计

**FaqAccordion 组件:**

```vue
<!-- frontend/src/components/help/FaqAccordion.vue -->
<template>
  <div class="faq-accordion">
    <div
      v-for="item in items"
      :key="item.id"
      class="faq-item"
    >
      <button
        class="faq-question"
        :aria-expanded="expandedItems.includes(item.id)"
        @click="toggleItem(item.id)"
      >
        <span class="question-text">{{ item.question }}</span>
        <n-icon :component="expandedItems.includes(item.id) ? ChevronUp : ChevronDown" />
      </button>
      <Transition name="collapse">
        <div
          v-show="expandedItems.includes(item.id)"
          class="faq-answer"
        >
          {{ item.answer }}
        </div>
      </Transition>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { ChevronUp, ChevronDown } from '@vicons/ionicons5';

interface FaqItem {
  id: string;
  question: string;
  answer: string;
}

defineProps<{
  items: FaqItem[];
}>();

const expandedItems = ref<string[]>([]);

function toggleItem(id: string) {
  const index = expandedItems.value.indexOf(id);
  if (index === -1) {
    expandedItems.value.push(id);
  } else {
    expandedItems.value.splice(index, 1);
  }
}
</script>
```

**FeedbackButton 组件:**

```vue
<!-- frontend/src/components/help/FeedbackButton.vue -->
<template>
  <div class="feedback-button">
    <n-button
      type="primary"
      ghost
      @click="showModal = true"
    >
      <template #icon>
        <n-icon :component="HelpCircleOutline" />
      </template>
      遇到问题？
    </n-button>

    <n-modal
      v-model:show="showModal"
      preset="card"
      title="问题反馈"
      style="width: 500px;"
    >
      <n-form ref="formRef" :model="formValue" :rules="rules">
        <n-form-item label="问题描述" path="description">
          <n-input
            v-model:value="formValue.description"
            type="textarea"
            placeholder="请描述你遇到的问题..."
            :rows="4"
          />
        </n-form-item>
        <n-form-item label="联系方式" path="contact">
          <n-input
            v-model:value="formValue.contact"
            placeholder="手机号或邮箱（选填）"
          />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" @click="submitFeedback">提交</n-button>
      </template>
    </n-modal>
  </div>
</template>
```

### 项目结构规范

**新增文件位置:**

```
frontend/
├── src/
│   ├── views/
│   │   └── help/
│   │       └── FeishuBotTutorial.vue        # 教程页面（新增）
│   ├── components/
│   │   └── help/
│   │       ├── FaqAccordion.vue              # FAQ 组件（新增）
│   │       └── FeedbackButton.vue            # 反馈按钮（新增）
│   ├── assets/
│   │   └── images/
│   │       └── tutorials/
│   │           └── feishu-bot/               # 教程截图（新增）
│   │               ├── step1-search.png
│   │               ├── step2-add.png
│   │               ├── step3-test.png
│   │               └── ...
│   └── data/
│       └── faq/
│           └── feishu-bot-faq.ts             # FAQ 数据（新增）
```

### 测试标准

**测试要求:**
- 测试框架: Vitest + Vue Test Utils
- 测试覆盖率目标: ≥ 70%

**测试用例设计:**

| 测试场景 | 测试文件 | 测试方法 |
|---------|---------|---------|
| 页面渲染 | `FeishuBotTutorial.spec.ts` | 组件测试 |
| 步骤展示 | 复用 Story 8.1 测试 | - |
| FAQ 展开/收起 | `FaqAccordion.spec.ts` | 组件测试 |
| 反馈表单提交 | `FeedbackButton.spec.ts` | 组件测试 |
| 键盘导航 | `FaqAccordion.spec.ts` | 集成测试 |

### Project Structure Notes

**与 Epic 8 其他 Story 的关系:**

1. **Story 8.1 (飞书开放平台配置教程)**: 复用 TutorialStep、TutorialImage 组件
2. **Story 8.4 (帮助中心)**: 帮助中心页面链接到此教程页面

### 前序 Story 的学习经验

**从 Story 8.1 (飞书开放平台配置教程) 获得的经验:**

1. **组件复用**: TutorialStep、TutorialImage 组件可直接复用
2. **图片管理**: 按步骤组织图片资源
3. **无障碍实现**: 图片 alt、键盘导航等已实现

### 常见问题与解决方案

**问题 1: FAQ 展开动画性能**
- **原因**: Vue Transition 需要正确处理高度变化
- **解决**: 使用 CSS max-height 过渡或 JS 动画库

**问题 2: 反馈提交接口**
- **原因**: 需要后端支持反馈提交 API
- **解决**: 可以先实现前端 Mock，后续对接真实 API

### 无障碍注意事项

1. **FAQ 折叠组件**:
   - 使用 button 元素触发展开/收起
   - 添加 aria-expanded 属性
   - 支持键盘 Enter/Space 操作

2. **反馈表单**:
   - 表单字段关联 label
   - 错误提示可被屏幕阅读器识别
   - 提交按钮有明确的加载状态

### References

- [Source: ux-design-specification.md] - UX 设计规范
- [Source: architecture.md#Frontend] - 前端技术栈
- [Source: epics.md#Story 8.2] - 原始故事定义
- [Source: prd.md#FR36] - 飞书机器人添加教程需求
- [Source: prd.md#FR23] - 系统提供飞书机器人添加教程
- [Source: 8-1-feishu-platform-tutorial.md] - 复用组件参考

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
