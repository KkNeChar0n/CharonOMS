# 任务列表：系统设置迁移与权限种植

## 任务

- [x] **[T1] 修复新建角色默认状态**
  - 修改角色entity或Service层，新建角色status默认为1（禁用）
  - 验证：创建角色后查看状态为禁用
  - 依赖：无

- [x] **[T2] 添加角色编辑状态守卫**
  - Service层UpdateRole和UpdateRolePermissions前校验角色status
  - 启用中的角色返回400："角色启用中，无法编辑"
  - 依赖：无

- [x] **[T3] 添加菜单编辑状态守卫**
  - Service层UpdateMenu前校验菜单status
  - 启用中的菜单返回400："菜单启用中，无法编辑"
  - 依赖：无

- [x] **[T4] 添加菜单sort_order唯一性校验**
  - Service层UpdateMenu中校验同parent_id下sort_order不重复
  - 冲突时返回400："同级菜单中排序已存在"
  - 依赖：T3

- [x] **[T5] 创建数据库种植脚本**
  - 从backup提取menu和permissions的INSERT语句
  - 使用INSERT IGNORE保证幂等性
  - 文件位置：scripts/seed_permissions.sql
  - 依赖：无

- [x] **[T6] 执行种植脚本并验证**
  - 运行种植脚本填入数据
  - 验证前端菜单正常显示
  - 验证权限树结构正确
  - 依赖：T5

## 并行关系

- T1、T2、T3、T5 互相独立，可并行实施
- T4 依赖 T3 完成
- T6 依赖 T5 完成

## 预估工作项

| 编号 | 描述 | 类型 |
|------|------|------|
| T1 | 角色默认状态修复 | 代码修改 |
| T2 | 角色编辑守卫 | 代码修改 |
| T3 | 菜单编辑守卫 | 代码修改 |
| T4 | sort_order校验 | 代码修改 |
| T5 | 种植脚本 | SQL脚本 |
| T6 | 执行验证 | 验证 |
