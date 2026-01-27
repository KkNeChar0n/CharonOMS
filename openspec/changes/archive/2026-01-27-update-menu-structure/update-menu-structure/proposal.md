# Change: 统一菜单结构与原项目保持完全一致

## Why

当前Go项目的菜单结构与原Python Flask项目存在多处不一致，导致：

1. **Route命名规范冲突**：Go项目使用kebab-case（如`menu-management`），原项目使用snake_case（如`menu_management`），导致前端调用API时路由不匹配
2. **缺失重要菜单功能**：原项目有5个退款相关菜单和详细的分账管理功能，Go项目仅有1个退款菜单
3. **审批管理模块重设计**：Go项目采用三层审批（我发起的/待我审批/已审批），原项目采用更详细的审批流程管理（审批流类型/模板/管理）
4. **前端兼容性问题**：前端代码基于原项目菜单结构开发，当前菜单不匹配导致页面无法正常显示

这些不一致严重影响了系统的前后端集成和用户体验。

## What Changes

### 1. 修正Route命名规范（6处修改）

将所有kebab-case格式的route改为snake_case，与原项目保持一致：

| 当前Route | 修改为 | 影响菜单 |
|----------|--------|---------|
| `menu-management` | `menu_management` | 菜单管理 |
| `activity-templates` | `activity_template` | 活动模板 |
| `activities` | `activity_management` | 活动列表 |
| `contracts` | `contract_management` | 合同管理 |
| `payment-collections` | `payment_collection` | 收款管理 |
| `refund-orders` | `refund_orders` | 退款订单 |
| `approval-initiated` | `approval_flow_type` | **BREAKING** 结构性变更 |
| `approval-pending` | `approval_flow_template` | **BREAKING** 结构性变更 |
| `approval-completed` | `approval_flow_management` | **BREAKING** 结构性变更 |

### 2. 补充缺失的财务管理菜单（4个新增）

在财务管理一级菜单下补充原项目的完整退款和分账功能：

- **退款子订单** (`refund_childorders`) - 退款订单明细
- **退款管理** (`refund_management`) - 常规退费和淘宝退费管理
- **退款支付详情** (`refund_payment_detail`) - 退款支付明细
- **分账明细** (`separate_account`) - 分账管理

### 3. 修正审批管理结构（**BREAKING** 重大变更）

将当前的三层审批改为原项目的设计：

**当前结构（删除）：**
- 我发起的 (approval-initiated)
- 待我审批 (approval-pending)
- 已审批 (approval-completed)

**新结构（采用原项目）：**
- 审批流类型 (approval_flow_type) - 管理审批流程类型定义
- 审批流模板 (approval_flow_template) - 管理审批流程模板配置
- 审批流管理 (approval_flow_management) - 查看我发起的、待我审批、已审批的流程

### 4. 更新菜单ID编号规则

调整菜单ID以与原项目保持一致的编号区间：

| 一级菜单 | ID编号区间 | 说明 |
|---------|----------|------|
| 系统管理 | 1-14 | 系统管理相关菜单 |
| 学生管理 | 2, 5 | 学生相关菜单 |
| 教练管理 | 2, 6 | 教练相关菜单 |
| 商品管理 | 15-19 | 商品体系菜单 |
| 订单管理 | 3, 7, 20, 29, 37 | 订单和退款订单 |
| 活动管理 | 21-23 | 营销活动菜单 |
| 财务管理 | 24-28, 38, 39 | 合同、收款、退款、分账 |
| 审批管理 | 30-34 | 审批流程管理 |

### 5. 修正中文名称细微差异

| 当前名称 | 修改为 | 原因 |
|---------|--------|------|
| 学生列表 | 学生管理 | 与原项目保持一致 |
| 教练列表 | 教练管理 | 与原项目保持一致 |
| 订单列表 | 订单管理 | 与原项目保持一致 |
| 活动列表 | 活动管理 | 与原项目保持一致 |

## Impact

### 影响的规范
- `specs/rbac/spec.md` - 菜单管理需求需要更新

### 影响的代码

**数据库初始化脚本：**
- `scripts/init_data_bcrypt.sql` - 需要完整重写菜单定义（第28-74行）
- `scripts/init_data_final.sql` - 同步更新

**路由配置：**
- `internal/interfaces/http/router/router.go` - 需要更新所有受影响的路由注册

**前端代码：**
- `frontend/app.js` - 前端菜单route引用已基于原项目，无需修改
- `frontend/index.html` - 前端页面组件引用已基于原项目，无需修改

### 破坏性变更

**⚠️ BREAKING CHANGES：**

1. **所有Menu Route变更**：
   - 影响范围：前端路由跳转、API调用、权限检查
   - 迁移方案：更新所有硬编码的route引用

2. **审批管理模块完全重构**：
   - 删除当前的3个审批菜单
   - 新增原项目的3个审批菜单
   - 需要同步实现后端API和前端页面

3. **菜单ID重新编号**：
   - 影响范围：权限表中的menu_id外键关联
   - 迁移方案：提供数据迁移SQL脚本

### 回滚计划

如果变更导致问题，提供回滚SQL脚本：
```sql
-- 备份当前菜单
CREATE TABLE menu_backup AS SELECT * FROM menu;

-- 回滚时恢复
TRUNCATE TABLE menu;
INSERT INTO menu SELECT * FROM menu_backup;
```

## Migration Path

1. **Phase 1**: 更新数据库初始化脚本（非破坏性）
2. **Phase 2**: 运行数据迁移脚本，更新现有数据库菜单
3. **Phase 3**: 更新路由配置，支持新的route命名
4. **Phase 4**: 验证前端页面能否正常访问所有菜单
5. **Phase 5**: 删除旧的审批管理菜单，实现新的审批流程功能

## Validation

### 验证标准

- [ ] 所有菜单的route与原项目完全一致
- [ ] 前端能够正常访问所有菜单页面
- [ ] 数据库中的菜单数量和原项目一致（24个二级菜单）
- [ ] 权限ID和菜单ID的映射关系正确
- [ ] 超级管理员能看到所有菜单
- [ ] 普通角色根据权限只能看到授权的菜单

### 测试场景

1. 使用admin账号登录，验证能看到完整的菜单树
2. 点击每个二级菜单，验证route跳转正确
3. 使用非超级管理员账号登录，验证权限过滤生效
4. 测试新增的4个退款/分账菜单功能
5. 测试重构后的审批管理菜单功能
