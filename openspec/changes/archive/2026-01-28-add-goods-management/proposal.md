# 提案：迁移商品管理功能

## 变更ID
`add-goods-management`

## 背景与目的

从 ZhixinStudentSaaS (Python Flask) 迁移商品体系到 CharonOMS (Go)。商品管理是订单系统的核心依赖，涉及品牌、分类、属性等多个实体及其管理功能。

商品体系包含单商品和组合商品（套餐），组合商品由多个单商品组成，总价为子商品价格之和。

## 迁移范围

### 品牌管理 (5接口)
- `GET /api/brands` — 获取品牌列表（全部）
- `POST /api/brands` — 新增品牌
- `PUT /api/brands/:id` — 编辑品牌
- `PUT /api/brands/:id/status` — 更新品牌状态
- `GET /api/brands/active` — 获取启用品牌（用于商品表单下拉）

### 分类管理 (6接口)
- `GET /api/classifies` — 获取分类列表（含父级信息）
- `GET /api/classifies/parents` — 获取启用的一级分类（用于二级分类下拉）
- `POST /api/classifies` — 新增分类
- `PUT /api/classifies/:id` — 编辑分类
- `PUT /api/classifies/:id/status` — 更新分类状态
- `GET /api/classifies/active` — 获取启用的二级分类（用于商品表单下拉）

### 属性管理 (7接口)
- `GET /api/attributes` — 获取属性列表（含值数量统计）
- `POST /api/attributes` — 新增属性
- `PUT /api/attributes/:id` — 编辑属性
- `PUT /api/attributes/:id/status` — 更新属性状态
- `GET /api/attributes/:id/values` — 获取指定属性的值列表
- `POST /api/attributes/:id/values` — 保存属性值（全量替换）
- `GET /api/attributes/active` — 获取启用属性及其值（用于商品表单）

### 商品管理 (9接口)
- `GET /api/goods` — 获取商品列表（支持分类和状态筛选，含品牌/分类名称/属性字符串）
- `POST /api/goods` — 新增商品（含属性绑定和组合商品子商品）
- `GET /api/goods/:id` — 获取商品详情（含attributevalue_ids和included_goods_ids数组）
- `PUT /api/goods/:id` — 编辑商品（全量替换属性和子商品关联）
- `PUT /api/goods/:id/status` — 更新商品状态
- `GET /api/goods/active-for-order` — 获取订单可用商品（状态启用，含总价计算）
- `GET /api/goods/:id/included-goods` — 获取组合商品的子商品列表
- `GET /api/goods/:id/total-price` — 计算商品总价（单商品=price，组合商品=子商品价格之和）
- `GET /api/goods/available-for-combo` — 获取可用于组合的商品（启用的单商品）

## 数据库表（已存在）

| 表名 | 说明 | 关键字段 |
|------|------|----------|
| `brand` | 品牌 | id, name, status, create_time, update_time |
| `classify` | 分类（两级） | id, name, level(0=一级,1=二级), parentid, status |
| `attribute` | 属性/规格定义 | id, name, classify(0=属性,1=规格), status |
| `attribute_value` | 属性值 | id, name, attributeid |
| `goods` | 商品 | id, name, brandid, classifyid, isgroup(0=套餐,1=单品), price, status |
| `goods_attributevalue` | 商品-属性值关联 | id, goodsid, attributevalueid |
| `goods_goods` | 组合商品-子商品关联 | id, goodsid(子商品), parentsid(父商品) |

## 关键业务规则

1. **品牌名称全局唯一**
2. **分类层级**：一级分类(level=0) 名称全局唯一；二级分类(level=1) 名称在同父级下唯一
3. **属性分类**：classify=0 为属性，classify=1 为规格
4. **属性值保存**：POST /attributes/:id/values 全量替换该属性下所有值
5. **组合商品**：isgroup=0 时必须包含至少一个子商品；isgroup 字段创建后不可修改
6. **组合商品总价**：等于所有子商品 price 之和
7. **商品属性绑定**：编辑时先删除旧关联再插入新关联（全量替换）
8. **状态值**：0=启用，1=禁用

## 验收标准

- 全部 27 个接口实现并通过测试
- 响应格式与原 Python 项目完全一致
- 品牌/分类/属性/商品的 CRUD 流程正常
- 组合商品的子商品管理和总价计算正确
- 属性值全量替换逻辑正确
- DDD 四层架构实现
- 路由注册替换占位符
