# 任务清单：实现商品管理模块

> 总计 27 个接口，涉及 brand/classify/attribute/goods 四个子模块。
> 品牌、分类、属性为商品的依赖模块，优先实现。

---

## Phase 1: 数据库Schema验证和准备

- [ ] 验证所有 7 张表已存在于数据库
  - [ ] brand: id, name, status, create_time, update_time
  - [ ] classify: id, name, level, parentid, status, created_at, updated_at
  - [ ] attribute: id, name, classify, status, create_time, update_time
  - [ ] attribute_value: id, name, attributeid
  - [ ] goods: id, name, brandid, classifyid, isgroup, price, status, create_time, update_time
  - [ ] goods_attributevalue: id, goodsid, attributevalueid, create_time
  - [ ] goods_goods: id, goodsid, parentsid, create_time
- [ ] 插入测试数据（可选）
  - [ ] 品牌：2-3条
  - [ ] 分类：一级2条 + 二级3条
  - [ ] 属性：2条属性 + 2条规格，每条含2-3个值
  - [ ] 商品：3条单商品 + 1条组合商品

---

## Phase 2: Brand（品牌）模块实现

- [ ] Domain层
  - [ ] 创建 `internal/domain/brand/entity/brand.go`
  - [ ] 创建 `internal/domain/brand/repository/brand_repository.go`（接口）
- [ ] Infrastructure层
  - [ ] 创建 `internal/infrastructure/persistence/mysql/brand/brand_repository_impl.go`
  - [ ] 实现 GetAll、GetActive、GetByID、GetByName、Create、Update、UpdateStatus
- [ ] Application层
  - [ ] 创建 `internal/application/service/brand/dto.go`
  - [ ] 创建 `internal/application/service/brand/brand_service.go`
  - [ ] 实现名称唯一性验证逻辑
- [ ] Interface层
  - [ ] 创建 `internal/interfaces/http/handler/brand/brand_handler.go`
  - [ ] 实现 5 个 HTTP 处理方法
- [ ] 路由注册

---

## Phase 3: Classify（分类）模块实现

- [ ] Domain层
  - [ ] 创建 `internal/domain/classify/entity/classify.go`
  - [ ] 创建 `internal/domain/classify/repository/classify_repository.go`（接口）
- [ ] Infrastructure层
  - [ ] 创建 `internal/infrastructure/persistence/mysql/classify/classify_repository_impl.go`
  - [ ] 实现 GetAll（含LEFT JOIN parent_name）、GetParents、GetActive
  - [ ] 实现 Create、Update、UpdateStatus
  - [ ] 实现名称唯一性检查（按级别和父级）
- [ ] Application层
  - [ ] 创建 `internal/application/service/classify/dto.go`
  - [ ] 创建 `internal/application/service/classify/classify_service.go`
  - [ ] 实现两级分类名称唯一性验证
  - [ ] 实现级别值验证(0或1)
  - [ ] 实现二级分类必须有parent_id验证
- [ ] Interface层
  - [ ] 创建 `internal/interfaces/http/handler/classify/classify_handler.go`
  - [ ] 实现 6 个 HTTP 处理方法
- [ ] 路由注册

---

## Phase 4: Attribute（属性）模块实现

- [ ] Domain层
  - [ ] 创建 `internal/domain/attribute/entity/attribute.go`
  - [ ] 创建 `internal/domain/attribute/entity/attribute_value.go`
  - [ ] 创建 `internal/domain/attribute/repository/attribute_repository.go`（接口）
- [ ] Infrastructure层
  - [ ] 创建 `internal/infrastructure/persistence/mysql/attribute/attribute_repository_impl.go`
  - [ ] 实现 GetAll（含LEFT JOIN value_count统计）
  - [ ] 实现 GetActive（含嵌套values数组）
  - [ ] 实现 Create、Update、UpdateStatus
  - [ ] 实现 GetValues（按属性ID查询值列表）
  - [ ] 实现 SaveValues（DELETE + INSERT 全量替换）
- [ ] Application层
  - [ ] 创建 `internal/application/service/attribute/dto.go`
  - [ ] 创建 `internal/application/service/attribute/attribute_service.go`
  - [ ] 实现classify值验证(0或1)
  - [ ] 实现属性值非空验证
- [ ] Interface层
  - [ ] 创建 `internal/interfaces/http/handler/attribute/attribute_handler.go`
  - [ ] 实现 7 个 HTTP 处理方法
- [ ] 路由注册

---

## Phase 5: Goods（商品）模块实现

- [ ] Domain层
  - [ ] 创建 `internal/domain/goods/entity/goods.go`
  - [ ] 创建 `internal/domain/goods/repository/goods_repository.go`（接口）
- [ ] Infrastructure层
  - [ ] 创建 `internal/infrastructure/persistence/mysql/goods/goods_repository_impl.go`
  - [ ] 实现 GetList（LEFT JOIN brand+classify，含attributes字符串构建）
  - [ ] 实现 GetByID（含attributevalue_ids和included_goods_ids）
  - [ ] 实现 GetActiveForOrder（含total_price子查询计算）
  - [ ] 实现 GetAvailableForCombo（启用+单商品筛选）
  - [ ] 实现 GetIncludedGoods（组合商品子商品列表）
  - [ ] 实现 GetTotalPrice（单商品=price，组合=子商品价格之和）
  - [ ] 实现 Create（插入goods + goods_attributevalue + goods_goods）
  - [ ] 实现 Update（全量替换属性关联和子商品关联）
  - [ ] 实现 UpdateStatus
- [ ] Application层
  - [ ] 创建 `internal/application/service/goods/dto.go`（interface{} 接收数字字段）
  - [ ] 创建 `internal/application/service/goods/goods_service.go`
  - [ ] 实现必填字段验证
  - [ ] 实现组合商品至少含一个子商品验证
- [ ] Interface层
  - [ ] 创建 `internal/interfaces/http/handler/goods/goods_handler.go`
  - [ ] 实现 9 个 HTTP 处理方法
  - [ ] 注意路由顺序：静态路径在参数路径之前
- [ ] 路由注册（替换占位符）

---

## Phase 6: 测试验证

- [ ] 单商品 CRUD 流程测试
  - [ ] 新增单商品（绑定属性）
  - [ ] 获取列表（验证attributes字符串格式）
  - [ ] 获取详情（验证attributevalue_ids数组）
  - [ ] 编辑商品（替换属性绑定）
  - [ ] 状态切换
- [ ] 组合商品流程测试
  - [ ] 新增组合商品（含子商品）
  - [ ] 获取详情（验证included_goods_ids）
  - [ ] 编辑组合商品（替换子商品）
  - [ ] 获取总价（验证子商品价格之和）
  - [ ] 获取组合子商品列表
- [ ] 辅助接口测试
  - [ ] active-for-order（含total_price）
  - [ ] available-for-combo（仅启用单商品）
- [ ] 品牌/分类/属性 CRUD 测试
  - [ ] 名称唯一性验证
  - [ ] 分类层级和父级约束
  - [ ] 属性值全量替换
- [ ] 错误情况测试
  - [ ] 必填字段缺失
  - [ ] 组合商品无子商品
  - [ ] ID不存在
  - [ ] 名称重复

---

## Phase 7: 文档更新

- [ ] 更新 REFACTORING_STATUS.md
  - [ ] 品牌管理 100%（5接口）
  - [ ] 分类管理 100%（6接口）
  - [ ] 属性管理 100%（7接口）
  - [ ] 商品管理 100%（9接口）

---

## 依赖关系

```
Phase 1 (DB验证) → Phase 2 (Brand) ─┐
                  → Phase 3 (Classify) ─┤→ Phase 5 (Goods)
                  → Phase 4 (Attribute) ─┘
Phase 5 (Goods) → Phase 6 (测试)
```

Phase 2、3、4 可并行实现（互不依赖），但 Phase 5 依赖 2、3、4 完成。

---

## 注意事项

1. **路由顺序**：Goods 模块中 `/goods/active-for-order` 和 `/goods/available-for-combo` 必须注册在 `/goods/:id` 之前
2. **Vue 字段类型**：商品 price 和 brandid 等数字字段可能以字符串形式从前端发送，DTO 使用 `interface{}`
3. **attributes 字符串格式**：`"属性名:值1/值2,属性名2:值3"` — 按属性分组，多值用 `/` 分隔，不同属性用 `,` 分隔
4. **组合商品总价**：需要子查询 `SUM(sub_goods.price)`，组合商品自身的 price 字段存储标准售价但显示时用总价
5. **classify 表的时间字段**：使用 `created_at`/`updated_at`（与其他表的 `create_time`/`update_time` 不同）
