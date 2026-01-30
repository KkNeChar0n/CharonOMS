# 营销管理模块测试文档

## 概述

本文档描述营销管理模块的完整测试流程,包括活动模板和活动管理的所有API端点。

## 测试环境

- **服务器地址**: http://localhost:8080
- **数据库**: charonoms (MySQL)
- **认证方式**: JWT Bearer Token
- **测试用户**: admin / admin123

## 测试数据准备

测试数据已通过以下SQL脚本导入:

```bash
mysql -h localhost -u root -pqweasd123Q! charonoms < scripts/test_data_marketing.sql
```

测试数据包括:
- 3个活动模板 (ID: 1, 2, 3)
  - 模板1: 春季满减促销 (按分类选择, 启用)
  - 模板2: 夏季满折促销 (按商品选择, 启用)
  - 模板3: 秋季满赠促销 (按分类选择, 禁用)
- 3个活动 (ID: 1, 2, 3)
  - 活动1: 春季大促销 (模板1, 启用, 2026-03-01 至 2026-03-31)
  - 活动2: 夏季满折活动 (模板2, 启用, 2026-06-01 至 2026-06-30)
  - 活动3: 测试已禁用活动 (模板1, 禁用, 2026-01-01 至 2026-01-31)

## 测试脚本

提供两个测试脚本:

### Bash脚本 (Linux/Mac/Git Bash)

```bash
cd scripts
chmod +x test_marketing_api.sh
./test_marketing_api.sh
```

### PowerShell脚本 (Windows)

```powershell
cd scripts
.\test_marketing_api.ps1
```

## API测试用例

### 1. 认证

#### POST /api/login

**请求:**
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**响应:**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

### 2. 活动模板管理

#### GET /api/activity-templates

获取所有活动模板列表

**请求头:**
```
Authorization: Bearer {token}
```

**响应:**
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "id": 1,
      "name": "春季满减促销",
      "type": 1,
      "select_type": 1,
      "status": 0,
      "create_time": "2026-01-30T14:30:40Z",
      "update_time": "2026-01-30T14:30:40Z"
    }
  ]
}
```

#### GET /api/activity-templates/active

获取启用的活动模板列表 (status=0)

**响应:** 同上,仅返回status=0的模板

#### GET /api/activity-templates/:id

获取活动模板详情,包含关联的商品或分类信息

**响应:**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "name": "春季满减促销",
    "type": 1,
    "select_type": 1,
    "status": 0,
    "create_time": "2026-01-30T14:30:40Z",
    "update_time": "2026-01-30T14:30:40Z",
    "classify_list": [
      {
        "classify_id": 15,
        "classify_name": "电子设备"
      }
    ]
  }
}
```

#### POST /api/activity-templates

创建新的活动模板

**请求:**
```json
{
  "name": "测试满减活动模板",
  "type": 1,
  "select_type": 1,
  "classify_ids": [15, 16],
  "status": 0
}
```

**响应:**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 4
  }
}
```

**验证规则:**
- name: 必填
- type: 必填,取值 1(满减)/2(满折)/3(满赠)
- select_type: 必填,取值 1(按分类)/2(按商品)
- select_type=1时,classify_ids不能为空
- select_type=2时,goods_ids不能为空

#### PUT /api/activity-templates/:id

更新活动模板

**请求:** 同POST

**响应:**
```json
{
  "code": 0,
  "msg": "success"
}
```

**验证规则:**
- 启用中的模板不能编辑 (status=0)

#### PUT /api/activity-templates/:id/status

更新活动模板状态

**请求:**
```json
{
  "status": 1
}
```

**响应:**
```json
{
  "code": 0,
  "msg": "success"
}
```

#### DELETE /api/activity-templates/:id

删除活动模板

**响应:**
```json
{
  "code": 0,
  "msg": "success"
}
```

**验证规则:**
- 有关联活动的模板不能删除

### 3. 活动管理

#### GET /api/activities

获取所有活动列表

**响应:**
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "id": 1,
      "name": "春季大促销",
      "template_id": 1,
      "template_name": "春季满减促销",
      "template_type": 1,
      "start_time": "2026-03-01T00:00:00Z",
      "end_time": "2026-03-31T23:59:59Z",
      "status": 0,
      "create_time": "2026-01-30T14:30:40Z"
    }
  ]
}
```

#### GET /api/activities/:id

获取活动详情,包含折扣规则

**响应:**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 2,
    "name": "夏季满折活动",
    "template_id": 2,
    "template_name": "夏季满折促销",
    "template_type": 2,
    "select_type": 2,
    "start_time": "2026-06-01T00:00:00Z",
    "end_time": "2026-06-30T23:59:59Z",
    "status": 0,
    "details": [
      {
        "id": 1,
        "activity_id": 2,
        "threshold_amount": 100.00,
        "discount_value": 0.95
      },
      {
        "id": 2,
        "activity_id": 2,
        "threshold_amount": 200.00,
        "discount_value": 0.90
      }
    ]
  }
}
```

#### POST /api/activities

创建新活动

**请求:**
```json
{
  "name": "测试满减活动",
  "template_id": 1,
  "start_time": "2026-02-01T00:00:00Z",
  "end_time": "2026-02-28T23:59:59Z",
  "status": 0,
  "details": [
    {
      "threshold_amount": 100.0,
      "discount_value": 10.0
    },
    {
      "threshold_amount": 200.0,
      "discount_value": 25.0
    }
  ]
}
```

**响应:**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 4
  }
}
```

**验证规则:**
- name, template_id, start_time, end_time: 必填
- 关联的模板必须启用 (status=0)
- start_time 必须早于 end_time
- details中的discount_value对于满折类型(type=2)必须在0-1之间

#### PUT /api/activities/:id

更新活动

**请求:** 同POST

**响应:**
```json
{
  "code": 0,
  "msg": "success"
}
```

**验证规则:**
- 启用中的活动不能编辑 (status=0)
- 其他验证规则同POST

#### PUT /api/activities/:id/status

更新活动状态

**请求:**
```json
{
  "status": 1
}
```

**响应:**
```json
{
  "code": 0,
  "msg": "success"
}
```

#### DELETE /api/activities/:id

删除活动 (关联的活动详情会级联删除)

**响应:**
```json
{
  "code": 0,
  "msg": "success"
}
```

#### GET /api/activities/by-date-range

根据付款时间查询有效活动,并检测冲突

**查询参数:**
- payment_time: 必填,日期时间格式 (支持RFC3339, "2006-01-02 15:04:05", "2006-01-02")

**示例:**
```
GET /api/activities/by-date-range?payment_time=2026-02-15T12:00:00Z
```

**响应:**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "has_duplicate": false,
    "duplicate_type": null,
    "type_name": "",
    "activities": [
      {
        "id": 1,
        "name": "春季大促销",
        "template_id": 1,
        "template_name": "春季满减促销",
        "template_type": 1,
        "template_select_type": 1,
        "start_time": "2026-03-01T00:00:00Z",
        "end_time": "2026-03-31T23:59:59Z",
        "status": 0,
        "create_time": "2026-01-30T14:30:40Z",
        "details": []
      }
    ]
  }
}
```

**冲突检测:**
- 如果同一时间段内存在多个相同类型的活动,has_duplicate=true
- duplicate_type 指示冲突的活动类型
- type_name 显示冲突类型的中文名称 ("满减"/"满折"/"满赠")

## 测试场景

### 场景1: 创建和管理活动模板

1. 创建按分类选择的满减模板
2. 创建按商品选择的满折模板
3. 查询模板列表,验证创建成功
4. 获取模板详情,验证关联数据正确
5. 更新模板信息
6. 更新模板状态为禁用
7. 尝试删除有关联活动的模板 (应失败)
8. 删除无关联活动的模板 (应成功)

### 场景2: 创建和管理活动

1. 选择启用的模板创建活动
2. 添加活动详情 (折扣规则)
3. 查询活动列表,验证创建成功
4. 获取活动详情,验证规则正确
5. 更新活动信息和规则
6. 更新活动状态为禁用
7. 删除活动

### 场景3: 日期范围查询和冲突检测

1. 创建两个不同类型的活动在同一时间段
2. 查询该时间段,验证无冲突
3. 创建两个相同类型的活动在同一时间段
4. 查询该时间段,验证检测到冲突

### 场景4: 业务规则验证

1. 尝试用禁用的模板创建活动 (应失败)
2. 尝试编辑启用中的模板 (应失败)
3. 尝试编辑启用中的活动 (应失败)
4. 验证时间范围 (start_time < end_time)
5. 验证满折类型的折扣值在0-1之间

## 错误处理

所有API都应正确处理以下错误:

1. **认证错误** (401)
   - 缺少或无效的JWT token

2. **验证错误** (400)
   - 缺少必填字段
   - 字段值不符合规则
   - 业务规则违反

3. **资源不存在** (404)
   - 模板或活动ID不存在

4. **服务器错误** (500)
   - 数据库连接失败
   - 其他内部错误

## 性能测试

建议测试:
1. 并发创建多个模板
2. 查询大量活动时的响应时间
3. 复杂日期范围查询的性能

## 总结

营销管理模块包含:
- 7个活动模板API端点
- 7个活动管理API端点
- 完整的CRUD操作
- 业务规则验证
- 冲突检测机制

所有实现遵循:
- DDD领域驱动设计
- RESTful API规范
- 与Python原项目API兼容
- 完整的事务支持
- 详细的错误处理
