# 商品管理模块实现文档

## 模块概述

商品管理模块是 CharonOMS 系统中最复杂的模块之一，支持单品和组合商品（套餐）管理。该模块严格遵循 DDD（领域驱动设计）四层架构。

## 架构层次

### 1. Domain 层（领域层）

#### 实体 (Entity)
- `internal/domain/goods/entity/goods.go`
  - **Goods**: 商品实体
  - **GoodsAttributeValue**: 商品与属性值关联实体
  - **GoodsGoods**: 商品组合关系实体

#### 仓储接口 (Repository)
- `internal/domain/goods/repository/goods_repository.go`
  - 定义了 9 个方法接口：
    - `GetList()`: 获取商品列表（含品牌、分类、属性）
    - `GetByID()`: 获取商品详情（含属性值ID数组和子商品ID数组）
    - `GetActiveForOrder()`: 获取可用于下单的商品
    - `GetAvailableForCombo()`: 获取可用于组合的单品
    - `GetIncludedGoods()`: 获取子商品列表
    - `GetTotalPrice()`: 计算商品总价
    - `Create()`: 创建商品
    - `Update()`: 更新商品
    - `UpdateStatus()`: 更新状态

### 2. Infrastructure 层（基础设施层）

#### 仓储实现
- `internal/infrastructure/persistence/mysql/goods/goods_repository_impl.go`
  - 实现所有 Repository 接口方法
  - **关键实现**：
    - `GetList`: LEFT JOIN 三张表（brand, classify, attribute），使用 `GROUP_CONCAT` 构建 attributes 字符串
    - `GetByID`: 返回 `attributevalue_ids` 数组和 `included_goods_ids` 数组
    - `GetActiveForOrder`: 使用 `CASE WHEN` 计算 total_price（单品=price，组合=SUM子商品price）
    - `GetTotalPrice`: 使用子查询计算总价
    - `Create/Update`: 使用事务处理，全量替换关联表数据（goods_attributevalue 和 goods_goods）

### 3. Application 层（应用层）

#### DTOs
- `internal/application/service/goods/dto.go`
  - `CreateGoodsRequest`: 数字字段使用 `interface{}` 接收（brandid, classifyid, isgroup, price）
  - `UpdateGoodsRequest`: 同上
  - `UpdateStatusRequest`: 状态字段使用 `interface{}`

#### 业务服务
- `internal/application/service/goods/goods_service.go`
  - 包含辅助函数：`toInt()`, `toFloat64()`（参考合同模块）
  - **业务验证**：
    - 验证必填字段
    - 验证组合商品至少包含一个子商品
    - 错误消息统一使用中文

### 4. Interface 层（接口层）

#### HTTP 处理器
- `internal/interfaces/http/handler/goods/goods_handler.go`
  - 实现 9 个 HTTP 方法
  - **错误处理**：
    - "商品信息不完整" (400)
    - "组合商品必须至少包含一个子商品" (400)
    - "商品不存在" (404)
    - "状态不能为空" (400)

## 数据库表结构

### goods 表
```sql
CREATE TABLE `goods` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(200) NOT NULL,
  `brandid` int NOT NULL,
  `classifyid` int NOT NULL,
  `isgroup` tinyint NOT NULL DEFAULT '1',  -- 0=套餐，1=单品
  `price` decimal(10,2) NOT NULL,
  `status` tinyint NOT NULL DEFAULT '0',   -- 0=启用，1=禁用
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
)
```

### goods_attributevalue 表（商品属性关联）
```sql
CREATE TABLE `goods_attributevalue` (
  `id` int NOT NULL AUTO_INCREMENT,
  `goodsid` int NOT NULL,
  `attributevalueid` int NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_goods_attrvalue` (`goodsid`,`attributevalueid`)
)
```

### goods_goods 表（商品组合关联）
```sql
CREATE TABLE `goods_goods` (
  `id` int NOT NULL AUTO_INCREMENT,
  `goodsid` int NOT NULL,      -- 子商品ID
  `parentsid` int NOT NULL,    -- 父商品ID（组合商品）
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_goods_parents` (`goodsid`,`parentsid`)
)
```

## API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/goods | 获取商品列表 |
| GET | /api/goods/active-for-order | 获取可用于下单的商品 |
| GET | /api/goods/available-for-combo | 获取可用于组合的单品 |
| GET | /api/goods/:id | 获取商品详情 |
| GET | /api/goods/:id/included | 获取子商品列表 |
| GET | /api/goods/:id/total-price | 计算商品总价 |
| POST | /api/goods | 创建商品 |
| PUT | /api/goods/:id | 更新商品 |
| PUT | /api/goods/:id/status | 更新商品状态 |

## 关键特性

### 1. 组合商品逻辑
- **IsGroup = 0**：套餐（组合商品），可包含多个子商品
- **IsGroup = 1**：单品
- 组合商品的总价 = 所有子商品价格之和
- 组合商品必须至少包含一个子商品

### 2. 属性管理
- 商品可关联多个属性值
- 列表查询时自动构建 attributes 字符串（格式："属性名:值1/值2,属性名2:值3"）
- 详情查询返回 attributevalue_ids 数组

### 3. 事务处理
- 创建/更新商品使用事务
- 采用全量替换策略处理关联表（先删除旧关联，再创建新关联）

### 4. 状态管理
- **Status = 0**: 启用（可用于下单）
- **Status = 1**: 禁用

## 待办事项

- [ ] 路由注册（参考 ROUTES.md）
- [ ] 依赖注入配置
- [ ] 单元测试
- [ ] 集成测试

## 注意事项

1. **路由顺序**：静态路径（/active-for-order）必须在参数路径（/:id）之前注册
2. **数据类型转换**：前端可能发送字符串形式的数字，Service 层使用 toInt/toFloat64 转换
3. **空数组处理**：确保返回空数组而非 nil
4. **错误消息**：统一使用中文错误消息
5. **事务安全**：关联表更新失败时需要回滚整个事务
