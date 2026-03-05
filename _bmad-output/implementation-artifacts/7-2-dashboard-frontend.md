# Story 7.2: 管理仪表盘前端页面

Status: ready-for-dev

## Story

As a 平台管理员,
I want 在仪表盘查看平台运营概览,
so that 快速了解平台运行状况。

## Acceptance Criteria

1. **AC1: 用户总数卡片**
   - **Given** 管理员已登录管理后台
   - **When** 访问仪表盘页面
   - **Then** 显示用户总数卡片
   - **And** 显示数字和环比变化百分比
   - **And** 环比上升显示绿色箭头，下降显示红色箭头

2. **AC2: 活跃实例数卡片**
   - **Given** 管理员已登录管理后台
   - **When** 访问仪表盘页面
   - **Then** 显示活跃实例数卡片
   - **And** 显示活跃实例数/总实例数

3. **AC3: 部署成功率卡片**
   - **Given** 管理员已登录管理后台
   - **When** 访问仪表盘页面
   - **Then** 显示部署成功率卡片
   - **And** 显示进度条可视化
   - **And** 显示最近 7 天成功率百分比

4. **AC4: 系统可用性卡片**
   - **Given** 管理员已登录管理后台
   - **When** 访问仪表盘页面
   - **Then** 显示系统可用性卡片
   - **And** 显示最近 7 天可用性百分比
   - **And** 可用性 >= 99% 显示绿色，< 99% 显示黄色

5. **AC5: 最近告警列表**
   - **Given** 管理员已登录管理后台
   - **When** 访问仪表盘页面
   - **Then** 显示最近告警列表
   - **And** 每条告警显示类型图标、级别、时间、状态
   - **And** 点击告警可查看详情

6. **AC6: 页面性能要求**
   - **Given** 仪表盘页面已实现
   - **When** 访问仪表盘页面
   - **Then** 页面加载时间 < 2 秒
   - **And** 支持数据自动刷新（可选）

## Tasks / Subtasks

- [ ] Task 1: 创建仪表盘页面路由和布局 (AC: 1-6)
  - [ ] 1.1 创建 `src/views/admin/Dashboard.vue` 页面组件
  - [ ] 1.2 配置路由 `/admin/dashboard`
  - [ ] 1.3 创建页面布局（卡片网格布局）
  - [ ] 1.4 添加页面标题和刷新按钮

- [ ] Task 2: 创建统计卡片组件 (AC: 1-4)
  - [ ] 2.1 创建 `src/components/admin/StatCard.vue` 通用卡片组件
  - [ ] 2.2 实现 UserStatsCard 用户统计卡片
  - [ ] 2.3 实现 InstanceStatsCard 实例统计卡片
  - [ ] 2.4 实现 DeployStatsCard 部署成功率卡片（含进度条）
  - [ ] 2.5 实现 SystemStatsCard 系统可用性卡片

- [ ] Task 3: 创建告警列表组件 (AC: 5)
  - [ ] 3.1 创建 `src/components/admin/AlertList.vue` 告警列表组件
  - [ ] 3.2 实现告警类型图标（CPU/内存/磁盘）
  - [ ] 3.3 实现告警级别样式（warning/critical）
  - [ ] 3.4 实现告警状态标签（pending/resolved）
  - [ ] 3.5 实现点击查看详情功能

- [ ] Task 4: 创建 API 服务和数据获取 (AC: 1-6)
  - [ ] 4.1 创建 `src/api/admin/dashboard.ts` API 服务
  - [ ] 4.2 实现 fetchDashboardStats() 方法
  - [ ] 4.3 创建 `src/stores/admin/dashboard.ts` Pinia store
  - [ ] 4.4 实现数据加载状态管理
  - [ ] 4.5 实现错误处理和重试机制

- [ ] Task 5: 实现页面响应式布局 (AC: 1-6)
  - [ ] 5.1 使用 Naive UI Grid 实现响应式布局
  - [ ] 5.2 桌面端 4 列布局
  - [ ] 5.3 平板端 2 列布局
  - [ ] 5.4 移动端单列布局

- [ ] Task 6: 添加单元测试 (AC: 1-6)
  - [ ] 6.1 创建 `Dashboard.spec.ts` 页面测试
  - [ ] 6.2 创建 `StatCard.spec.ts` 组件测试
  - [ ] 6.3 创建 `AlertList.spec.ts` 组件测试
  - [ ] 6.4 确保测试覆盖率 >= 70%

## Dev Notes

### 架构模式与约束

**必须遵循的前端架构原则：**
1. **组件化**: 可复用的卡片组件
2. **状态管理**: Pinia store 管理仪表盘数据
3. **响应式设计**: 支持多设备访问

**关键架构决策 [Source: architecture.md]:**
- 前端框架: Vite + Vue 3 + TypeScript + Naive UI
- 状态管理: Pinia（按模块组织）
- 布局策略: Web 桌面优先，平板 P1

### 组件设计

**StatCard 通用卡片组件:**

```vue
<template>
  <n-card class="stat-card">
    <div class="stat-card__header">
      <span class="stat-card__title">{{ title }}</span>
      <n-icon :component="icon" />
    </div>
    <div class="stat-card__value">{{ formattedValue }}</div>
    <div class="stat-card__footer">
      <span :class="trendClass">
        <n-icon :component="trendIcon" />
        {{ trendText }}
      </span>
    </div>
  </n-card>
</template>
```

**Props:**
| 属性 | 类型 | 描述 |
|------|------|------|
| title | string | 卡片标题 |
| value | number | 主要数值 |
| icon | Component | 图标组件 |
| trend | number | 环比变化百分比 |
| format | string | 格式化类型 (number/percent) |

**AlertList 告警列表组件:**

```vue
<template>
  <n-card title="最近告警" class="alert-list">
    <n-list>
      <n-list-item v-for="alert in alerts" :key="alert.id">
        <div class="alert-item" @click="showDetail(alert)">
          <n-tag :type="levelType">{{ alert.level }}</n-tag>
          <span class="alert-type">{{ alert.type }}</span>
          <span class="alert-time">{{ formatTime(alert.createdAt) }}</span>
        </div>
      </n-list-item>
    </n-list>
  </n-card>
</template>
```

### 页面布局设计

**仪表盘页面布局:**

```
+------------------------------------------+
|  管理仪表盘                    [刷新按钮] |
+------------------------------------------+
|  +--------+  +--------+  +--------+  +--------+
|  | 用户   |  | 实例   |  | 部署   |  | 系统   |
|  | 总数   |  | 活跃   |  | 成功率 |  | 可用性 |
|  +--------+  +--------+  +--------+  +--------+
+------------------------------------------+
|  最近告警列表                              |
|  +------------------------------------+  |
|  | [warning] CPU告警  10:30 [查看]    |  |
|  | [critical] 内存告警 09:15 [查看]   |  |
|  +------------------------------------+  |
+------------------------------------------+
```

**响应式断点:**

| 断点 | 列数 | 描述 |
|------|------|------|
| >= 1200px | 4 列 | 大屏桌面 |
| 768-1199px | 2 列 | 平板/小桌面 |
| < 768px | 1 列 | 移动端 |

### API 集成

**API 服务 (src/api/admin/dashboard.ts):**

```typescript
import { request } from '@/utils/request';

export interface DashboardStats {
  userStats: UserStats;
  instanceStats: InstanceStats;
  deployStats: DeployStats;
  systemStats: SystemStats;
  recentAlerts: AlertSummary[];
}

export const dashboardApi = {
  getStats(): Promise<DashboardStats> {
    return request.get('/v1/admin/dashboard');
  }
};
```

**Pinia Store (src/stores/admin/dashboard.ts):**

```typescript
import { defineStore } from 'pinia';
import { dashboardApi } from '@/api/admin/dashboard';

export const useDashboardStore = defineStore('admin-dashboard', {
  state: () => ({
    stats: null as DashboardStats | null,
    loading: false,
    error: null as Error | null,
  }),

  actions: {
    async fetchStats() {
      this.loading = true;
      try {
        this.stats = await dashboardApi.getStats();
      } catch (error) {
        this.error = error;
      } finally {
        this.loading = false;
      }
    }
  }
});
```

### 样式规范

**配色方案 [Source: ux-design-specification.md]:**

| 状态 | 颜色 | 用途 |
|------|------|------|
| 成功/正常 | #18A058 | 高可用性、上升趋势 |
| 警告 | #F0A020 | 中等告警 |
| 错误/严重 | #D03050 | 严重告警、下降趋势 |
| 信息 | #2080F0 | 主要按钮 |

**卡片样式:**

```css
.stat-card {
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.stat-card__value {
  font-size: 32px;
  font-weight: 600;
  color: #333;
}

.stat-card__trend--up {
  color: #18A058;
}

.stat-card__trend--down {
  color: #D03050;
}
```

### 性能优化

1. **数据缓存**:
   - 使用 Pinia 缓存仪表盘数据
   - 设置合理的刷新间隔（可选 30 秒）

2. **懒加载**:
   - 告警列表支持分页懒加载
   - 使用 v-show 替代 v-if 优化切换性能

3. **骨架屏**:
   - 数据加载时显示骨架屏
   - 提升用户感知速度

### 与其他 Story 的依赖关系

**前序依赖:**
- Story 1.2: 前端项目初始化 - 需要前端项目结构
- Story 2.6: 管理员登录页面 - 需要管理员认证
- Story 7.1: 管理仪表盘数据 API - 需要后端 API

**后续依赖:**
- Story 7.3: 实例日志查看功能 - 告警详情可跳转

### 测试标准

**组件测试要求:**
- 测试所有 props 变化
- 测试事件触发
- 测试样式变化

**E2E 测试要求:**
- 页面正常渲染
- 数据正确显示
- 告警列表交互正常

### References

- [Source: architecture.md#Frontend] - 前端技术栈
- [Source: ux-design-specification.md] - UX 设计规范
- [Source: epics.md#Story 7.2] - 原始故事定义
- [Source: prd.md#FR28-FR31] - 功能需求定义
- [Source: prd.md#NFR-P3] - 页面加载性能要求

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
