# 任务清单：统一菜单结构

## Phase 1: 数据库脚本更新（非破坏性）

- [ ] 备份当前的初始化脚本
  - [ ] 复制 `scripts/init_data_bcrypt.sql` 到 `scripts/init_data_bcrypt.sql.backup`
  - [ ] 复制 `scripts/init_data_final.sql` 到 `scripts/init_data_final.sql.backup`

- [ ] 更新 `scripts/init_data_bcrypt.sql` 中的菜单定义（第28-74行）
  - [ ] 修正6处route命名（kebab-case → snake_case）
    - [ ] `menu-management` → `menu_management`
    - [ ] `activity-templates` → `activity_template`
    - [ ] `activities` → `activity_management`
    - [ ] `contracts` → `contract_management`
    - [ ] `payment-collections` → `payment_collection`
    - [ ] `refund-orders` → `refund_orders`
  - [ ] 补充4个缺失的财务管理菜单
    - [ ] 退款子订单 (`refund_childorders`)
    - [ ] 退款管理 (`refund_management`)
    - [ ] 退款支付详情 (`refund_payment_detail`)
    - [ ] 分账明细 (`separate_account`)
  - [ ] 重构审批管理模块（3个菜单）
    - [ ] 删除：我发起的 (`approval-initiated`)
    - [ ] 删除：待我审批 (`approval-pending`)
    - [ ] 删除：已审批 (`approval-completed`)
    - [ ] 新增：审批流类型 (`approval_flow_type`)
    - [ ] 新增：审批流模板 (`approval_flow_template`)
    - [ ] 新增：审批流管理 (`approval_flow_management`)
  - [ ] 调整菜单ID编号与原项目一致
  - [ ] 修正中文名称细微差异
    - [ ] 学生列表 → 学生管理
    - [ ] 教练列表 → 教练管理
    - [ ] 订单列表 → 订单管理
    - [ ] 活动列表 → 活动管理

- [ ] 同步更新 `scripts/init_data_final.sql`

- [ ] 验证SQL语法正确性
  - [ ] 使用MySQL语法检查工具验证
  - [ ] 确认所有外键引用正确

## Phase 2: 数据迁移脚本

- [ ] 创建数据迁移脚本 `scripts/migrate_menu_structure.sql`
  - [ ] 备份现有菜单数据
    ```sql
    CREATE TABLE IF NOT EXISTS menu_backup_20260127 AS SELECT * FROM menu;
    ```
  - [ ] 更新现有菜单的route字段
  - [ ] 插入4个新增的财务管理菜单
  - [ ] 删除旧的3个审批管理菜单
  - [ ] 插入新的3个审批管理菜单
  - [ ] 更新菜单ID（如果需要）
  - [ ] 更新菜单中文名称

- [ ] 创建回滚脚本 `scripts/rollback_menu_structure.sql`
  ```sql
  TRUNCATE TABLE menu;
  INSERT INTO menu SELECT * FROM menu_backup_20260127;
  ```

- [ ] 测试迁移脚本
  - [ ] 在测试数据库上运行迁移
  - [ ] 验证迁移后的菜单数量和结构
  - [ ] 测试回滚脚本是否正常工作

## Phase 3: 路由配置更新

- [ ] 更新 `internal/interfaces/http/router/router.go`
  - [ ] 修改菜单管理路由：`/menu-management` → `/menu_management`
  - [ ] 修改活动模板路由：`/activity-templates` → `/activity_template`
  - [ ] 修改活动管理路由：`/activities` → `/activity_management`
  - [ ] 修改合同管理路由：`/contracts` → `/contract_management`
  - [ ] 修改收款管理路由：`/payment-collections` → `/payment_collection`
  - [ ] 修改退款订单路由：`/refund-orders` → `/refund_orders`
  - [ ] 添加4个新的财务管理路由（占位符处理器）
    - [ ] `/refund_childorders` → 返回占位符响应
    - [ ] `/refund_management` → 返回占位符响应
    - [ ] `/refund_payment_detail` → 返回占位符响应
    - [ ] `/separate_account` → 返回占位符响应
  - [ ] 添加3个新的审批管理路由（占位符处理器）
    - [ ] `/approval_flow_type` → 返回占位符响应
    - [ ] `/approval_flow_template` → 返回占位符响应
    - [ ] `/approval_flow_management` → 返回占位符响应
  - [ ] 删除旧的审批管理路由
    - [ ] `/approval-initiated`
    - [ ] `/approval-pending`
    - [ ] `/approval-completed`

- [ ] 创建占位符处理器 `internal/interfaces/http/handler/placeholder/placeholder_handler.go`
  - [ ] 实现通用占位符响应函数
  - [ ] 返回 `{message: "功能开发中", status: "pending"}`

## Phase 4: 前端验证

- [ ] 启动开发服务器
  - [ ] 运行数据迁移脚本
  - [ ] 启动Go后端服务
  - [ ] 使用admin账号登录

- [ ] 菜单显示验证
  - [ ] 验证左侧菜单树显示完整（8个一级菜单，24个二级菜单）
  - [ ] 验证菜单中文名称正确
  - [ ] 验证菜单图标正常显示

- [ ] 路由跳转验证
  - [ ] 点击每个二级菜单，验证route跳转正确
  - [ ] 验证浏览器地址栏URL与菜单route一致
  - [ ] 验证已实现功能的页面正常加载
  - [ ] 验证占位符功能返回开发中提示

- [ ] 权限验证
  - [ ] 使用超级管理员账号，验证能看到所有菜单
  - [ ] 创建测试角色，分配部分权限
  - [ ] 使用测试账号登录，验证只能看到授权的菜单
  - [ ] 验证无权限菜单不显示在菜单树中

## Phase 5: 文档更新

- [ ] 更新 `REFACTORING_STATUS.md`
  - [ ] 记录菜单结构统一完成
  - [ ] 更新进度百分比
  - [ ] 标记占位符功能为待实现

- [ ] 更新 `README.md`
  - [ ] 添加菜单结构说明章节
  - [ ] 说明占位符功能的存在
  - [ ] 提供数据迁移指南

- [ ] 创建 `docs/MENU_MIGRATION_GUIDE.md`
  - [ ] 详细记录迁移步骤
  - [ ] 提供troubleshooting指南
  - [ ] 记录已知问题和限制

## Phase 6: 规范更新

- [ ] 更新 `specs/rbac/spec.md`
  - [ ] 添加菜单结构变更说明
  - [ ] 更新菜单管理需求
  - [ ] 说明route命名规范（使用snake_case）
  - [ ] 补充财务管理和审批管理的详细需求

## 验收标准

- [ ] 所有菜单route与原Python项目完全一致
- [ ] 数据库中共有24个二级菜单
- [ ] 前端能正常显示完整菜单树
- [ ] 点击菜单route跳转正确
- [ ] 权限过滤功能正常工作
- [ ] 已实现功能页面正常显示
- [ ] 占位符功能返回友好提示
- [ ] 文档完整更新

## 注意事项

1. **破坏性变更警告**：此变更会修改所有菜单route，影响前端路由跳转
2. **数据库备份**：执行迁移前必须备份menu表数据
3. **占位符处理**：新增菜单暂时返回占位符响应，不影响已有功能
4. **测试优先**：在测试环境完整验证后再应用到生产环境
5. **分阶段执行**：按照Phase顺序执行，每个Phase完成后验证再继续
