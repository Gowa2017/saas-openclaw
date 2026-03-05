# Story 7.5: 系统通知功能

Status: ready-for-dev

## Story

As a 平台管理员,
I want 向用户发送系统通知,
so that 可以告知用户重要信息。

## Acceptance Criteria

1. **AC1: 通知对象选择**
   - **Given** 管理员在通知管理页面
   - **When** 创建通知
   - **Then** 可以选择通知对象为全部用户
   - **And** 可以选择通知对象为指定用户
   - **And** 支持按用户名/邮箱搜索用户

2. **AC2: 通知内容编辑**
   - **Given** 管理员在创建通知页面
   - **When** 编辑通知内容
   - **Then** 可以编辑通知标题（必填，最长 100 字符）
   - **And** 可以编辑通知内容（必填，最长 2000 字符）
   - **And** 支持富文本编辑（可选）

3. **AC3: 用户端通知显示**
   - **Given** 管理员已发送通知
   - **When** 用户登录业务平台
   - **Then** 用户看到新通知提示（未读标记）
   - **And** 点击可查看通知详情
   - **And** 支持通知列表查看

4. **AC4: 通知已读/未读状态**
   - **Given** 用户收到通知
   - **When** 用户查看通知
   - **Then** 通知状态更新为已读
   - **And** 未读通知数量显示在导航栏
   - **And** 支持批量标记已读

5. **AC5: 通知发送日志**
   - **Given** 通知已发送
   - **When** 管理员查看通知日志
   - **Then** 显示发送时间、发送人、接收人数
   - **And** 显示已读/未读统计
   - **And** 支持按时间范围筛选

## Tasks / Subtasks

- [ ] Task 1: 创建通知数据模型 (AC: 1-5)
  - [ ] 1.1 创建 `internal/domain/notification/` 目录
  - [ ] 1.2 定义 Notification 结构体
  - [ ] 1.3 定义 NotificationRecipient 结构体
  - [ ] 1.4 定义 NotificationStatus 枚举
  - [ ] 1.5 创建数据库表和索引

- [ ] Task 2: 创建通知 Repository (AC: 1-5)
  - [ ] 2.1 创建 `internal/repository/notification_repository.go`
  - [ ] 2.2 实现 CreateNotification() 方法
  - [ ] 2.3 实现 GetNotificationByID() 方法
  - [ ] 2.4 实现 ListNotifications() 方法（管理员视角）
  - [ ] 2.5 实现 ListUserNotifications() 方法（用户视角）
  - [ ] 2.6 实现 MarkAsRead() 方法
  - [ ] 2.7 实现 MarkAllAsRead() 方法
  - [ ] 2.8 实现 GetUnreadCount() 方法
  - [ ] 2.9 实现 GetNotificationStats() 方法（已读/未读统计）

- [ ] Task 3: 创建通知服务 (AC: 1-5)
  - [ ] 3.1 创建 `internal/domain/notification/service.go`
  - [ ] 3.2 实现 CreateAndSendNotification() 方法
  - [ ] 3.3 实现用户筛选逻辑
  - [ ] 3.4 实现批量通知创建
  - [ ] 3.5 实现通知统计逻辑

- [ ] Task 4: 创建管理员通知 API (AC: 1, 2, 5)
  - [ ] 4.1 创建 `internal/api/admin/notification_handler.go`
  - [ ] 4.2 实现 POST /v1/admin/notifications 端点（创建通知）
  - [ ] 4.3 实现 GET /v1/admin/notifications 端点（通知列表）
  - [ ] 4.4 实现 GET /v1/admin/notifications/:id 端点（通知详情）
  - [ ] 4.5 实现 GET /v1/admin/notifications/:id/stats 端点（通知统计）
  - [ ] 4.6 实现 GET /v1/admin/users/search 端点（用户搜索）

- [ ] Task 5: 创建用户端通知 API (AC: 3, 4)
  - [ ] 5.1 创建 `internal/api/notification_handler.go`
  - [ ] 5.2 实现 GET /v1/notifications 端点（用户通知列表）
  - [ ] 5.3 实现 GET /v1/notifications/unread-count 端点（未读数量）
  - [ ] 5.4 实现 PUT /v1/notifications/:id/read 端点（标记已读）
  - [ ] 5.5 实现 PUT /v1/notifications/read-all 端点（全部已读）

- [ ] Task 6: 创建管理后台通知页面 (AC: 1, 2, 5)
  - [ ] 6.1 创建 `src/views/admin/NotificationManage.vue` 页面
  - [ ] 6.2 创建 `src/components/admin/NotificationForm.vue` 表单组件
  - [ ] 6.3 创建 `src/components/admin/NotificationList.vue` 列表组件
  - [ ] 6.4 创建 `src/components/admin/UserSelector.vue` 用户选择器
  - [ ] 6.5 创建 `src/components/admin/NotificationStats.vue` 统计组件

- [ ] Task 7: 创建用户端通知组件 (AC: 3, 4)
  - [ ] 7.1 创建 `src/components/common/NotificationBell.vue` 通知铃铛
  - [ ] 7.2 创建 `src/components/common/NotificationList.vue` 通知列表
  - [ ] 7.3 创建 `src/views/Notifications.vue` 通知页面
  - [ ] 7.4 实现未读数量徽章显示
  - [ ] 7.5 实现通知弹窗/抽屉

- [ ] Task 8: 添加单元测试 (AC: 1-5)
  - [ ] 8.1 创建 `notification_repository_test.go`
  - [ ] 8.2 创建 `notification_service_test.go`
  - [ ] 8.3 创建 `notification_handler_test.go`
  - [ ] 8.4 创建前端组件测试
  - [ ] 8.5 确保测试覆盖率 >= 70%

## Dev Notes

### 架构模式与约束

**必须遵循的架构原则：**
1. **数据隔离**: 通知按租户隔离
2. **批量处理**: 批量通知使用异步处理
3. **性能优化**: 大量用户时分页处理

**关键架构决策:**
- 通知类型: 系统通知（管理员发送）
- 通知存储: PostgreSQL
- 实时更新: 可选 WebSocket 推送

### 数据模型设计

**Notification 结构体:**

```go
type Notification struct {
    ID          string               `json:"id" db:"ID"`
    Title       string               `json:"title" db:"Title"`
    Content     string               `json:"content" db:"Content"`
    Type        NotificationType     `json:"type" db:"Type"`
    TargetType  NotificationTarget   `json:"targetType" db:"TargetType"`
    CreatedBy   string               `json:"createdBy" db:"CreatedBy"`
    CreatedAt   time.Time            `json:"createdAt" db:"CreatedAt"`
    Recipients  []NotificationRecipient `json:"recipients,omitempty"`
}

type NotificationRecipient struct {
    ID             string `json:"id" db:"ID"`
    NotificationID string `json:"notificationId" db:"NotificationID"`
    TenantID       string `json:"tenantId" db:"TenantID"`
    UserID         string `json:"userId" db:"UserID"`
    IsRead         bool   `json:"isRead" db:"IsRead"`
    ReadAt         *time.Time `json:"readAt,omitempty" db:"ReadAt"`
}

type NotificationType string

const (
    NotificationTypeSystem NotificationType = "system"
)

type NotificationTarget string

const (
    NotificationTargetAll      NotificationTarget = "all"
    NotificationTargetSpecific NotificationTarget = "specific"
)
```

### 数据库表设计

```sql
-- 通知主表
CREATE TABLE notifications (
    ID VARCHAR(36) PRIMARY KEY,
    Title VARCHAR(100) NOT NULL,
    Content TEXT NOT NULL,
    Type VARCHAR(20) NOT NULL DEFAULT 'system',
    TargetType VARCHAR(20) NOT NULL,
    CreatedBy VARCHAR(100) NOT NULL,
    CreatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 通知接收者表
CREATE TABLE notification_recipients (
    ID VARCHAR(36) PRIMARY KEY,
    NotificationID VARCHAR(36) NOT NULL REFERENCES notifications(ID) ON DELETE CASCADE,
    TenantID VARCHAR(36) NOT NULL,
    UserID VARCHAR(36) NOT NULL,
    IsRead BOOLEAN NOT NULL DEFAULT FALSE,
    ReadAt TIMESTAMP,
    CreatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(NotificationID, UserID)
);

-- 索引
CREATE INDEX idx_notifications_created_at ON notifications(CreatedAt);
CREATE INDEX idx_notification_recipients_tenant ON notification_recipients(TenantID);
CREATE INDEX idx_notification_recipients_user ON notification_recipients(UserID);
CREATE INDEX idx_notification_recipients_read ON notification_recipients(IsRead);
```

### API 端点设计

**管理员 API:**

| 端点 | 方法 | 描述 |
|------|------|------|
| /v1/admin/notifications | GET | 获取通知列表 |
| /v1/admin/notifications | POST | 创建通知 |
| /v1/admin/notifications/:id | GET | 获取通知详情 |
| /v1/admin/notifications/:id/stats | GET | 获取通知统计 |
| /v1/admin/users/search | GET | 搜索用户 |

**用户 API:**

| 端点 | 方法 | 描述 |
|------|------|------|
| /v1/notifications | GET | 获取用户通知列表 |
| /v1/notifications/unread-count | GET | 获取未读数量 |
| /v1/notifications/:id/read | PUT | 标记已读 |
| /v1/notifications/read-all | PUT | 全部标记已读 |

**创建通知请求体:**

```json
{
  "title": "系统维护通知",
  "content": "系统将于今晚 22:00 进行维护...",
  "targetType": "all"
}
```

或指定用户:

```json
{
  "title": "账户安全提醒",
  "content": "您的账户存在安全风险...",
  "targetType": "specific",
  "userIds": ["user-1", "user-2"]
}
```

**通知统计响应:**

```json
{
  "data": {
    "totalRecipients": 100,
    "readCount": 45,
    "unreadCount": 55,
    "readRate": 45.0
  }
}
```

### 前端组件设计

**用户选择器组件:**

```vue
<template>
  <div class="user-selector">
    <n-radio-group v-model:value="targetType">
      <n-radio value="all">全部用户</n-radio>
      <n-radio value="specific">指定用户</n-radio>
    </n-radio-group>

    <n-select
      v-if="targetType === 'specific'"
      v-model:value="selectedUsers"
      multiple
      filterable
      remote
      :loading="loading"
      :options="userOptions"
      @search="handleSearch"
    />
  </div>
</template>
```

**通知铃铛组件:**

```vue
<template>
  <n-badge :value="unreadCount" :max="99">
    <n-button quaternary circle @click="showNotifications = true">
      <template #icon>
        <n-icon :component="BellIcon" />
      </template>
    </n-button>
  </n-badge>

  <n-drawer v-model:show="showNotifications" width="400">
    <NotificationList @read="handleRead" />
  </n-drawer>
</template>
```

**通知列表组件:**

```vue
<template>
  <div class="notification-list">
    <div class="notification-list__header">
      <span>系统通知</span>
      <n-button text @click="markAllRead">全部已读</n-button>
    </div>

    <n-list>
      <n-list-item
        v-for="notification in notifications"
        :key="notification.id"
        :class="{ 'is-unread': !notification.isRead }"
        @click="viewDetail(notification)"
      >
        <div class="notification-item">
          <div class="notification-item__title">{{ notification.title }}</div>
          <div class="notification-item__time">{{ formatTime(notification.createdAt) }}</div>
        </div>
      </n-list-item>
    </n-list>
  </div>
</template>
```

### 性能优化

1. **批量插入优化**:
   - 发送给全部用户时，使用批量插入
   - 分批处理，每批 1000 条

2. **未读数量缓存**:
   - 使用 Redis 缓存未读数量
   - 用户标记已读时更新缓存

3. **分页查询**:
   - 通知列表分页加载
   - 默认每页 20 条

### 与其他 Story 的依赖关系

**前序依赖:**
- Story 2.1: 用户数据模型 - 需要用户表
- Story 2.4: 平台管理员独立认证系统 - 需要管理员信息
- Story 2.5: 前端登录页面与状态管理 - 需要用户认证状态

**后续依赖:**
- 无直接后续依赖

### 测试标准

**单元测试要求:**
- 通知创建逻辑测试
- 批量插入测试
- 已读/未读状态测试

**集成测试要求:**
- 完整通知流程测试
- 用户搜索功能测试

### 安全考虑

1. **权限控制**:
   - 只有管理员可以发送通知
   - 用户只能查看自己的通知

2. **数据隔离**:
   - 通知按租户隔离
   - 用户无法查看其他租户通知

3. **内容过滤**:
   - 标题和内容长度限制
   - 可选: 敏感词过滤

### References

- [Source: architecture.md#API Design] - RESTful API 设计规范
- [Source: epics.md#Story 7.5] - 原始故事定义
- [Source: prd.md#FR34] - 功能需求定义

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
