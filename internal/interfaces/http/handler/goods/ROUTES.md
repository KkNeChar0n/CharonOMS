# 商品管理模块路由注册说明

## 重要提示

**路由顺序非常重要**：静态路径必须在参数路径（如 `:id`）之前注册，否则静态路径会被参数路径匹配。

## 正确的路由注册顺序

```go
// 在 router.go 中注册商品路由时，按以下顺序：

goodsGroup := apiGroup.Group("/goods")
{
    // 1. 静态路径（无参数）- 必须在前面
    goodsGroup.GET("/active-for-order", goodsHandler.GetActiveForOrder)
    goodsGroup.GET("/available-for-combo", goodsHandler.GetAvailableForCombo)

    // 2. 带参数的路径 - 必须在后面
    goodsGroup.GET("/:id", goodsHandler.GetGoodsByID)
    goodsGroup.GET("/:id/included", goodsHandler.GetIncludedGoods)
    goodsGroup.GET("/:id/total-price", goodsHandler.GetTotalPrice)
    goodsGroup.PUT("/:id", goodsHandler.UpdateGoods)
    goodsGroup.PUT("/:id/status", goodsHandler.UpdateStatus)

    // 3. 其他路径
    goodsGroup.GET("", goodsHandler.GetGoods)
    goodsGroup.POST("", goodsHandler.CreateGoods)
}
```

## 错误示例（不要这样做）

```go
// ❌ 错误：参数路径在前，会导致 /active-for-order 被误匹配为 /:id
goodsGroup.GET("/:id", goodsHandler.GetGoodsByID)
goodsGroup.GET("/active-for-order", goodsHandler.GetActiveForOrder)  // 这个路由永远不会被匹配到
```

## API 端点列表

| 方法 | 路径 | 处理器 | 说明 |
|------|------|--------|------|
| GET | /api/goods | GetGoods | 获取商品列表 |
| GET | /api/goods/active-for-order | GetActiveForOrder | 获取可用于下单的商品列表 |
| GET | /api/goods/available-for-combo | GetAvailableForCombo | 获取可用于组合的单品商品列表 |
| GET | /api/goods/:id | GetGoodsByID | 获取商品详情 |
| GET | /api/goods/:id/included | GetIncludedGoods | 获取包含的子商品列表 |
| GET | /api/goods/:id/total-price | GetTotalPrice | 计算商品总价 |
| POST | /api/goods | CreateGoods | 创建商品 |
| PUT | /api/goods/:id | UpdateGoods | 更新商品 |
| PUT | /api/goods/:id/status | UpdateStatus | 更新商品状态 |

## 路由注册位置

路由应在 `internal/interfaces/http/router/router.go` 文件中注册。
