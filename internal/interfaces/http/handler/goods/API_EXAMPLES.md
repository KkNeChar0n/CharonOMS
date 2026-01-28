# 商品管理 API 请求响应示例

## 1. 获取商品列表

### 请求
```http
GET /api/goods
```

### 响应
```json
{
  "goods": [
    {
      "id": 1,
      "name": "钢琴基础课程",
      "brandid": 1,
      "brand_name": "雅马哈",
      "classifyid": 1,
      "classify_name": "钢琴课程",
      "isgroup": 1,
      "price": 5000.00,
      "status": 0,
      "create_time": "2026-01-20T10:00:00Z",
      "update_time": "2026-01-20T10:00:00Z",
      "attributes": "年龄段:儿童,课时:48节"
    },
    {
      "id": 2,
      "name": "音乐启蒙套餐",
      "brandid": 2,
      "brand_name": "珠江",
      "classifyid": 2,
      "classify_name": "组合课程",
      "isgroup": 0,
      "price": 8000.00,
      "status": 0,
      "create_time": "2026-01-21T14:30:00Z",
      "update_time": "2026-01-21T14:30:00Z",
      "attributes": "年龄段:幼儿"
    }
  ]
}
```

## 2. 获取商品详情

### 请求
```http
GET /api/goods/1
```

### 响应
```json
{
  "goods": {
    "id": 1,
    "name": "钢琴基础课程",
    "brandid": 1,
    "brand_name": "雅马哈",
    "classifyid": 1,
    "classify_name": "钢琴课程",
    "isgroup": 1,
    "price": 5000.00,
    "status": 0,
    "create_time": "2026-01-20T10:00:00Z",
    "update_time": "2026-01-20T10:00:00Z",
    "attributevalue_ids": [1, 2, 5],
    "included_goods_ids": []
  }
}
```

### 组合商品详情响应示例
```json
{
  "goods": {
    "id": 2,
    "name": "音乐启蒙套餐",
    "brandid": 2,
    "brand_name": "珠江",
    "classifyid": 2,
    "classify_name": "组合课程",
    "isgroup": 0,
    "price": 8000.00,
    "status": 0,
    "create_time": "2026-01-21T14:30:00Z",
    "update_time": "2026-01-21T14:30:00Z",
    "attributevalue_ids": [3],
    "included_goods_ids": [1, 3, 5]
  }
}
```

## 3. 获取可用于下单的商品列表

### 请求
```http
GET /api/goods/active-for-order
```

### 响应
```json
{
  "goods": [
    {
      "id": 1,
      "name": "钢琴基础课程",
      "brandid": 1,
      "brand_name": "雅马哈",
      "classifyid": 1,
      "classify_name": "钢琴课程",
      "isgroup": 1,
      "price": 5000.00,
      "total_price": 5000.00,
      "status": 0,
      "create_time": "2026-01-20T10:00:00Z"
    },
    {
      "id": 2,
      "name": "音乐启蒙套餐",
      "brandid": 2,
      "brand_name": "珠江",
      "classifyid": 2,
      "classify_name": "组合课程",
      "isgroup": 0,
      "price": 8000.00,
      "total_price": 12000.00,
      "status": 0,
      "create_time": "2026-01-21T14:30:00Z"
    }
  ]
}
```

## 4. 获取可用于组合的单品列表

### 请求
```http
GET /api/goods/available-for-combo
```

### 响应
```json
{
  "goods": [
    {
      "id": 1,
      "name": "钢琴基础课程",
      "brandid": 1,
      "brand_name": "雅马哈",
      "classifyid": 1,
      "classify_name": "钢琴课程",
      "price": 5000.00,
      "status": 0
    },
    {
      "id": 3,
      "name": "乐理课程",
      "brandid": 3,
      "brand_name": "施坦威",
      "classifyid": 3,
      "classify_name": "理论课程",
      "price": 3000.00,
      "status": 0
    }
  ]
}
```

## 5. 获取包含的子商品列表

### 请求
```http
GET /api/goods/2/included
```

### 响应
```json
{
  "goods": [
    {
      "id": 1,
      "name": "钢琴基础课程",
      "brandid": 1,
      "brand_name": "雅马哈",
      "classifyid": 1,
      "classify_name": "钢琴课程",
      "price": 5000.00,
      "status": 0
    },
    {
      "id": 3,
      "name": "乐理课程",
      "brandid": 3,
      "brand_name": "施坦威",
      "classifyid": 3,
      "classify_name": "理论课程",
      "price": 3000.00,
      "status": 0
    },
    {
      "id": 5,
      "name": "视唱练耳课程",
      "brandid": 1,
      "brand_name": "雅马哈",
      "classifyid": 4,
      "classify_name": "技能课程",
      "price": 4000.00,
      "status": 0
    }
  ]
}
```

## 6. 计算商品总价

### 请求
```http
GET /api/goods/2/total-price
```

### 响应
```json
{
  "total_price": 12000.00
}
```

## 7. 创建商品

### 单品商品请求
```http
POST /api/goods
Content-Type: application/json

{
  "name": "钢琴进阶课程",
  "brandid": 1,
  "classifyid": 1,
  "isgroup": 1,
  "price": 6000.00,
  "attributevalue_ids": [1, 2, 6],
  "included_goods_ids": []
}
```

### 组合商品请求
```http
POST /api/goods
Content-Type: application/json

{
  "name": "音乐全科套餐",
  "brandid": 2,
  "classifyid": 2,
  "isgroup": 0,
  "price": 15000.00,
  "attributevalue_ids": [3, 7],
  "included_goods_ids": [1, 3, 5, 8]
}
```

### 成功响应
```json
{
  "id": 10,
  "message": "商品新增成功"
}
```

### 错误响应 - 商品信息不完整
```json
{
  "error": "商品信息不完整"
}
```

### 错误响应 - 组合商品验证失败
```json
{
  "error": "组合商品必须至少包含一个子商品"
}
```

## 8. 更新商品

### 请求
```http
PUT /api/goods/1
Content-Type: application/json

{
  "name": "钢琴基础课程（升级版）",
  "brandid": 1,
  "classifyid": 1,
  "isgroup": 1,
  "price": 5500.00,
  "attributevalue_ids": [1, 2, 5, 9],
  "included_goods_ids": []
}
```

### 成功响应
```json
{
  "message": "操作成功"
}
```

### 错误响应 - 商品不存在
```json
{
  "error": "商品不存在"
}
```

## 9. 更新商品状态

### 请求 - 禁用商品
```http
PUT /api/goods/1/status
Content-Type: application/json

{
  "status": 1
}
```

### 请求 - 启用商品
```http
PUT /api/goods/1/status
Content-Type: application/json

{
  "status": 0
}
```

### 成功响应
```json
{
  "message": "操作成功"
}
```

### 错误响应 - 状态为空
```json
{
  "error": "状态不能为空"
}
```

### 错误响应 - 商品不存在
```json
{
  "error": "商品不存在"
}
```

## 数据类型说明

### 字段类型
- `id`: 整数
- `name`: 字符串
- `brandid`, `classifyid`: 整数（可以是字符串形式的数字）
- `isgroup`: 整数（0=套餐，1=单品，可以是字符串形式的数字）
- `price`, `total_price`: 浮点数（可以是字符串形式的数字）
- `status`: 整数（0=启用，1=禁用）
- `attributevalue_ids`: 整数数组
- `included_goods_ids`: 整数数组
- `create_time`, `update_time`: ISO 8601 时间格式

### 前端注意事项
1. **数字字段**：前端可以发送字符串或数字，后端会自动转换
2. **空数组**：如果没有属性或子商品，发送空数组 `[]`，不要发送 `null`
3. **组合商品**：`isgroup=0` 时，`included_goods_ids` 不能为空
4. **价格计算**：组合商品的 `total_price` 由后端自动计算（子商品价格之和）

## HTTP 状态码

- `200 OK`: 查询成功
- `201 Created`: 创建成功
- `400 Bad Request`: 请求参数错误或业务验证失败
- `404 Not Found`: 商品不存在
- `500 Internal Server Error`: 服务器内部错误
