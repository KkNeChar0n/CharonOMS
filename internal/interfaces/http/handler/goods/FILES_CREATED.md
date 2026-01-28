# 商品管理模块 - 创建文件清单

## 创建日期
2026-01-28

## 文件列表

### Domain 层（领域层）
1. **internal/domain/goods/entity/goods.go**
   - Goods 实体
   - GoodsAttributeValue 实体
   - GoodsGoods 实体
   - 表名映射

2. **internal/domain/goods/repository/goods_repository.go**
   - GoodsRepository 接口
   - 9个方法定义

### Infrastructure 层（基础设施层）
3. **internal/infrastructure/persistence/mysql/goods/goods_repository_impl.go**
   - GoodsRepositoryImpl 实现
   - 9个方法实现（含复杂SQL查询）
   - 事务处理逻辑

### Application 层（应用层）
4. **internal/application/service/goods/dto.go**
   - CreateGoodsRequest
   - UpdateGoodsRequest
   - UpdateStatusRequest
   - CreateGoodsResponse
   - MessageResponse
   - GoodsListResponse

5. **internal/application/service/goods/goods_service.go**
   - GoodsService 结构
   - 9个业务方法
   - toInt/toFloat64 辅助函数
   - 业务验证逻辑

### Interface 层（接口层）
6. **internal/interfaces/http/handler/goods/goods_handler.go**
   - GoodsHandler 结构
   - 9个HTTP处理方法
   - 错误处理和状态码映射

### 文档文件
7. **internal/interfaces/http/handler/goods/README.md**
   - 模块概述
   - 架构层次说明
   - 数据库表结构
   - API端点列表
   - 关键特性说明

8. **internal/interfaces/http/handler/goods/ROUTES.md**
   - 路由注册说明
   - 路由顺序重要性说明
   - 正确示例和错误示例
   - API端点表格

9. **internal/interfaces/http/handler/goods/IMPLEMENTATION_SUMMARY.md**
   - 实现完成清单
   - 核心实现要点
   - 待注册路由示例
   - 编译验证状态

10. **internal/interfaces/http/handler/goods/API_EXAMPLES.md**
    - 9个API的请求响应示例
    - 数据类型说明
    - HTTP状态码说明
    - 前端注意事项

11. **internal/interfaces/http/handler/goods/FILES_CREATED.md**
    - 本文件（文件清单）

## 统计

- **Go源代码文件**: 6个
- **文档文件**: 5个
- **总文件数**: 11个
- **代码行数**: 约800行（不含注释和空行）
- **文档行数**: 约600行

## 代码质量检查

- ✅ 所有代码编译通过
- ✅ Go代码格式化完成（gofmt）
- ✅ 遵循DDD四层架构
- ✅ 错误处理完善
- ✅ 中文错误消息
- ✅ 完整的注释说明

## 待办事项

- [ ] 在 router.go 中注册路由
- [ ] 在依赖注入中配置商品模块
- [ ] 编写单元测试
- [ ] 编写集成测试
- [ ] 测试组合商品创建和价格计算
- [ ] 前后端联调测试

## 特殊注意事项

### 路由注册顺序
必须按以下顺序注册路由：
1. 静态路径（/active-for-order, /available-for-combo）
2. 参数路径（/:id, /:id/included, /:id/total-price）
3. 其他路径

### 数据类型转换
- 前端可能发送字符串形式的数字
- Service层使用 toInt/toFloat64 自动转换
- 支持 float64, string, int 三种输入类型

### 组合商品验证
- isgroup = 0（套餐）时，必须至少包含一个子商品
- 组合商品总价 = 子商品价格之和

### 事务处理
- 创建/更新商品时使用事务
- 关联表采用全量替换策略（先删除，再创建）
- 任何步骤失败都会回滚整个事务
