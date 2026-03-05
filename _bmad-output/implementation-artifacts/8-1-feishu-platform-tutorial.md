# Story 8.1: 飞书开放平台配置教程页面

Status: ready-for-dev

## Story

As a 用户,
I want 查看详细的飞书开放平台配置教程,
so that 可以独立完成飞书应用创建。

## Acceptance Criteria

1. **AC1: 教程页面路由和入口**
   - **Given** 用户访问帮助中心
   - **When** 点击"飞书开放平台配置教程"
   - **Then** 跳转到教程页面 `/help/feishu-platform-tutorial`
   - **And** 页面标题显示"飞书开放平台配置教程"

2. **AC2: 教程步骤展示**
   - **Given** 用户进入教程页面
   - **When** 页面加载完成
   - **Then** 显示完整的图文教程
   - **And** 教程包含 5 个步骤
   - **And** 每步包含截图和文字说明

3. **AC3: 步骤内容完整性**
   - **Given** 教程页面已加载
   - **When** 查看教程内容
   - **Then** 步骤 1：注册飞书开放平台账号
   - **And** 步骤 2：创建企业自建应用
   - **And** 步骤 3：配置应用权限
   - **And** 步骤 4：获取 App ID 和 Secret
   - **And** 步骤 5：返回平台填写配置

4. **AC4: 打印和分享功能**
   - **Given** 用户在教程页面
   - **When** 点击"打印"按钮
   - **Then** 打开浏览器打印对话框
   - **And** 打印样式优化（隐藏导航、按钮等）
   - **When** 点击"分享"按钮
   - **Then** 复制页面链接到剪贴板

5. **AC5: 无障碍访问**
   - **Given** 用户使用屏幕阅读器
   - **When** 访问教程页面
   - **Then** 图片有 alt 文字描述
   - **And** 步骤标题使用正确的标题层级
   - **And** 支持键盘导航

## Tasks / Subtasks

- [ ] Task 1: 创建教程页面路由 (AC: 1)
  - [ ] 1.1 在 `frontend/src/router/` 添加教程页面路由
  - [ ] 1.2 创建 `frontend/src/views/help/FeishuPlatformTutorial.vue`
  - [ ] 1.3 配置页面 meta 信息（标题、描述）

- [ ] Task 2: 创建教程步骤组件 (AC: 2, 3)
  - [ ] 2.1 创建 `frontend/src/components/help/TutorialStep.vue` 步骤组件
  - [ ] 2.2 创建 `frontend/src/components/help/TutorialImage.vue` 图片组件
  - [ ] 2.3 定义教程内容数据结构

- [ ] Task 3: 实现教程内容 (AC: 3)
  - [ ] 3.1 准备飞书开放平台截图资源
  - [ ] 3.2 编写步骤 1：注册账号的内容和截图
  - [ ] 3.3 编写步骤 2：创建应用的内容和截图
  - [ ] 3.4 编写步骤 3：配置权限的内容和截图
  - [ ] 3.5 编写步骤 4：获取凭证的内容和截图
  - [ ] 3.6 编写步骤 5：填写配置的内容和截图

- [ ] Task 4: 实现打印和分享功能 (AC: 4)
  - [ ] 4.1 添加打印按钮和打印样式
  - [ ] 4.2 实现分享链接复制功能
  - [ ] 4.3 添加操作成功提示

- [ ] Task 5: 无障碍优化 (AC: 5)
  - [ ] 5.1 为所有图片添加 alt 属性
  - [ ] 5.2 检查标题层级结构
  - [ ] 5.3 添加键盘导航支持
  - [ ] 5.4 测试屏幕阅读器兼容性

- [ ] Task 6: 编写测试 (AC: 1-5)
  - [ ] 6.1 编写页面渲染测试
  - [ ] 6.2 编写步骤展示测试
  - [ ] 6.3 编写打印功能测试
  - [ ] 6.4 编写分享功能测试

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
├── package.json            # npm 依赖
└── vite.config.ts          # Vite 配置
```

**关键依赖版本:**
- Vue: 3.x
- TypeScript: 5.x
- Naive UI: 最新
- Vue Router: 4.x

### 教程内容设计

**步骤 1: 注册飞书开放平台账号**

```
标题: 注册飞书开放平台账号
内容:
1. 访问飞书开放平台官网: https://open.feishu.cn
2. 点击右上角"登录/注册"按钮
3. 使用企业邮箱或手机号完成注册
4. 完成账号验证

截图:
- 飞书开放平台首页截图
- 注册页面截图
```

**步骤 2: 创建企业自建应用**

```
标题: 创建企业自建应用
内容:
1. 登录后进入"开发者后台"
2. 点击"创建企业自建应用"
3. 填写应用名称（如：OpenClaw 助手）
4. 选择应用类型
5. 点击"创建"完成

截图:
- 开发者后台入口截图
- 创建应用表单截图
- 应用创建成功截图
```

**步骤 3: 配置应用权限**

```
标题: 配置应用权限
内容:
1. 进入应用详情页
2. 点击"权限管理"
3. 搜索并添加以下权限:
   - 获取用户基本信息
   - 接收消息
   - 发送消息
4. 发布版本使权限生效

截图:
- 权限管理页面截图
- 权限搜索截图
- 权限配置完成截图
```

**步骤 4: 获取 App ID 和 Secret**

```
标题: 获取 App ID 和 Secret
内容:
1. 在应用详情页找到"凭证与基础信息"
2. 复制 App ID（以 cli_ 开头）
3. 点击"查看"获取 App Secret
4. 安全保存这两个凭证

截图:
- 凭证位置截图（高亮标注）
- App ID 复制截图
- App Secret 查看截图
```

**步骤 5: 返回平台填写配置**

```
标题: 返回平台填写配置
内容:
1. 返回 OpenClaw 管理平台
2. 进入"飞书配置"页面
3. 填写 App ID
4. 填写 App Secret
5. 点击"验证配置"确认有效性
6. 保存配置

截图:
- 飞书配置页面截图
- 填写表单截图
- 验证成功截图
```

### 组件设计

**TutorialStep 组件:**

```vue
<!-- frontend/src/components/help/TutorialStep.vue -->
<template>
  <div class="tutorial-step">
    <div class="step-header">
      <span class="step-number">步骤 {{ number }}</span>
      <h2 class="step-title">{{ title }}</h2>
    </div>
    <div class="step-content">
      <ol class="step-instructions">
        <li v-for="(instruction, index) in instructions" :key="index">
          {{ instruction }}
        </li>
      </ol>
      <TutorialImage
        v-for="(image, index) in images"
        :key="index"
        :src="image.src"
        :alt="image.alt"
        :caption="image.caption"
      />
    </div>
  </div>
</template>
```

**TutorialImage 组件:**

```vue
<!-- frontend/src/components/help/TutorialImage.vue -->
<template>
  <figure class="tutorial-image">
    <img
      :src="src"
      :alt="alt"
      @click="openPreview"
    />
    <figcaption v-if="caption">{{ caption }}</figcaption>
  </figure>
</template>
```

### 打印样式设计

```css
/* 打印样式 */
@media print {
  /* 隐藏导航 */
  .app-header,
  .app-sidebar,
  .app-footer {
    display: none !important;
  }

  /* 隐藏按钮 */
  .print-hide {
    display: none !important;
  }

  /* 优化打印布局 */
  .tutorial-step {
    page-break-inside: avoid;
  }

  .tutorial-image img {
    max-width: 100%;
    height: auto;
  }
}
```

### 项目结构规范

**新增文件位置:**

```
frontend/
├── src/
│   ├── views/
│   │   └── help/
│   │       └── FeishuPlatformTutorial.vue   # 教程页面（新增）
│   ├── components/
│   │   └── help/
│   │       ├── TutorialStep.vue              # 步骤组件（新增）
│   │       └── TutorialImage.vue             # 图片组件（新增）
│   ├── assets/
│   │   └── images/
│   │       └── tutorials/
│   │           └── feishu-platform/          # 教程截图（新增）
│   │               ├── step1-home.png
│   │               ├── step1-register.png
│   │               ├── step2-console.png
│   │               ├── step2-create.png
│   │               ├── ...
│   └── router/
│       └── modules/
│           └── help.ts                       # 帮助模块路由（新增）
```

### 测试标准

**测试要求:**
- 测试框架: Vitest + Vue Test Utils
- 测试覆盖率目标: ≥ 70%

**测试用例设计:**

| 测试场景 | 测试文件 | 测试方法 |
|---------|---------|---------|
| 页面渲染 | `FeishuPlatformTutorial.spec.ts` | 组件测试 |
| 步骤展示 | `TutorialStep.spec.ts` | 组件测试 |
| 图片加载 | `TutorialImage.spec.ts` | 组件测试 |
| 打印功能 | `FeishuPlatformTutorial.spec.ts` | 集成测试 |
| 分享功能 | `FeishuPlatformTutorial.spec.ts` | 集成测试 |

### Project Structure Notes

**与 Epic 8 其他 Story 的关系:**

1. **Story 8.2 (飞书机器人添加教程)**: 复用 TutorialStep、TutorialImage 组件
2. **Story 8.4 (帮助中心)**: 帮助中心页面链接到此教程页面

### 前序 Story 的学习经验

**从 Story 1.2 (前端项目初始化) 获得的经验:**

1. **组件规范**: 使用 Vue 3 Composition API + TypeScript
2. **样式方案**: 使用 Naive UI 主题变量
3. **路由配置**: 模块化路由配置

### 常见问题与解决方案

**问题 1: 截图资源管理**
- **原因**: 教程截图较多，需要合理组织
- **解决**: 按步骤分类存放，使用 Webpack/Vite 自动导入

**问题 2: 打印样式问题**
- **原因**: 不同浏览器打印效果不一致
- **解决**: 使用 @media print 统一打印样式，测试主流浏览器

**问题 3: 图片加载性能**
- **原因**: 教程图片较多，影响页面加载
- **解决**: 使用图片懒加载，优化图片大小

### 无障碍注意事项

1. **图片替代文本**: 所有截图必须有描述性的 alt 文字
2. **键盘导航**: 确保所有交互元素可通过 Tab 键访问
3. **标题层级**: 使用正确的 h1-h6 层级结构
4. **色彩对比度**: 文字与背景对比度至少 4.5:1
5. **焦点指示器**: 为焦点元素提供可见的焦点样式

### References

- [Source: ux-design-specification.md] - UX 设计规范
- [Source: architecture.md#Frontend] - 前端技术栈
- [Source: epics.md#Story 8.1] - 原始故事定义
- [Source: prd.md#FR35] - 飞书开放平台配置教程需求
- [Source: prd.md#NFR-P3] - 页面加载时间 < 2 秒
- [Source: 1-2-frontend-project-init.md] - 前端项目实现经验

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
