# 审批流模块接口测试计划

## 测试环境
- 服务地址: http://localhost:5001
- 数据库: charonoms (MySQL)
- 认证方式: JWT Token

## 测试账号准备
需要准备以下测试账号:
- 管理员账号 (用于创建审批流类型和模板)
- 发起人账号 (用于发起审批流)
- 审批人账号1、2、3 (用于测试会签和或签场景)
- 抄送人账号 (用于接收抄送通知)

## API端点清单

### 1. 审批流类型管理 (4个接口)

#### 1.1 GET /api/approval_flow_type/list - 获取类型列表
**请求示例:**
```bash
curl -X GET "http://localhost:5001/api/approval_flow_type/list?status=0" \
  -H "Authorization: Bearer {token}"
```

**测试用例:**
- [ ] TC-001: 无筛选条件，返回所有类型
- [ ] TC-002: 按status=0筛选，返回启用的类型
- [ ] TC-003: 按status=1筛选，返回禁用的类型
- [ ] TC-004: 按name模糊搜索
- [ ] TC-005: 未授权访问，返回401

**预期结果:**
- 返回类型列表，包含id、name、status、create_time、update_time字段
- 状态码200

#### 1.2 POST /api/approval_flow_type/create - 创建类型
**请求示例:**
```bash
curl -X POST "http://localhost:5001/api/approval_flow_type/create" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "请假审批",
    "status": 0
  }'
```

**测试用例:**
- [ ] TC-006: 正常创建，返回成功
- [ ] TC-007: 缺少name字段，返回400
- [ ] TC-008: name为空字符串，返回400
- [ ] TC-009: 重复名称，返回错误
- [ ] TC-010: 未授权访问，返回401

**预期结果:**
- 创建成功，返回类型ID
- 状态码200

#### 1.3 GET /api/approval_flow_type/detail/:id - 获取类型详情
**请求示例:**
```bash
curl -X GET "http://localhost:5001/api/approval_flow_type/detail/1" \
  -H "Authorization: Bearer {token}"
```

**测试用例:**
- [ ] TC-011: 存在的ID，返回详情
- [ ] TC-012: 不存在的ID，返回404
- [ ] TC-013: 无效的ID格式，返回400
- [ ] TC-014: 未授权访问，返回401

**预期结果:**
- 返回类型详情对象
- 状态码200

#### 1.4 POST /api/approval_flow_type/update_status - 更新类型状态
**请求示例:**
```bash
curl -X POST "http://localhost:5001/api/approval_flow_type/update_status" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "status": 1
  }'
```

**测试用例:**
- [ ] TC-015: 启用 → 禁用，返回成功
- [ ] TC-016: 禁用 → 启用，返回成功
- [ ] TC-017: 不存在的ID，返回404
- [ ] TC-018: 无效的status值，返回400
- [ ] TC-019: 未授权访问，返回401

**预期结果:**
- 状态更新成功
- 状态码200

---

### 2. 审批流模板管理 (4个接口)

#### 2.1 GET /api/approval_flow_template/list - 获取模板列表
**请求示例:**
```bash
curl -X GET "http://localhost:5001/api/approval_flow_template/list?approval_flow_type_id=1&status=0" \
  -H "Authorization: Bearer {token}"
```

**测试用例:**
- [ ] TC-020: 无筛选，返回所有模板
- [ ] TC-021: 按类型筛选
- [ ] TC-022: 按状态筛选
- [ ] TC-023: 按名称模糊搜索
- [ ] TC-024: 多条件组合筛选
- [ ] TC-025: 未授权访问，返回401

**预期结果:**
- 返回模板列表，包含关联的flow_type_name
- 状态码200

#### 2.2 POST /api/approval_flow_template/create - 创建模板
**请求示例:**
```bash
curl -X POST "http://localhost:5001/api/approval_flow_template/create" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "approval_flow_type_id": 1,
    "name": "部门主管请假审批流程",
    "status": 0,
    "nodes": [
      {
        "type": 0,
        "sort": 1,
        "approver_ids": [2, 3]
      },
      {
        "type": 1,
        "sort": 2,
        "approver_ids": [4, 5, 6]
      }
    ],
    "copy_user_ids": [7, 8]
  }'
```

**测试用例:**
- [ ] TC-026: 正常创建单节点模板
- [ ] TC-027: 创建多节点模板（会签+或签）
- [ ] TC-028: 缺少nodes，返回400
- [ ] TC-029: nodes为空数组，返回400
- [ ] TC-030: 节点缺少approver_ids，返回400
- [ ] TC-031: approver_ids为空数组，返回400
- [ ] TC-032: 不存在的approval_flow_type_id，返回错误
- [ ] TC-033: 不存在的approver_id，返回错误
- [ ] TC-034: 未授权访问，返回401

**预期结果:**
- 事务创建成功，模板、节点、审批人、抄送人全部入库
- 状态码200

#### 2.3 GET /api/approval_flow_template/detail/:id - 获取模板详情
**请求示例:**
```bash
curl -X GET "http://localhost:5001/api/approval_flow_template/detail/1" \
  -H "Authorization: Bearer {token}"
```

**测试用例:**
- [ ] TC-035: 存在的模板ID，返回完整详情
- [ ] TC-036: 返回数据包含节点列表（按sort排序）
- [ ] TC-037: 每个节点包含审批人列表
- [ ] TC-038: 包含抄送人列表
- [ ] TC-039: 不存在的ID，返回404
- [ ] TC-040: 未授权访问，返回401

**预期结果:**
- 返回嵌套结构: template、nodes(含approvers)、copy_users
- 状态码200

#### 2.4 POST /api/approval_flow_template/update_status - 更新模板状态
**请求示例:**
```bash
curl -X POST "http://localhost:5001/api/approval_flow_template/update_status" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "status": 0
  }'
```

**测试用例:**
- [ ] TC-041: 启用模板，同类型其他模板自动禁用
- [ ] TC-042: 禁用模板
- [ ] TC-043: 启用已启用的模板（幂等性）
- [ ] TC-044: 不存在的ID，返回404
- [ ] TC-045: 未授权访问，返回401

**预期结果:**
- 启用时，同类型其他模板status变为1
- 状态码200

---

### 3. 审批流实例管理 (7个接口)

#### 3.1 POST /api/approval_flow_management/initiate - 发起审批流
**请求示例:**
```bash
curl -X POST "http://localhost:5001/api/approval_flow_management/initiate" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "approval_flow_template_id": 1,
    "title": "张三请假3天",
    "info": "因个人原因请假3天，2024-02-01至2024-02-03"
  }'
```

**测试用例:**
- [ ] TC-046: 正常发起审批流
- [ ] TC-047: 验证自动创建第一个节点实例
- [ ] TC-048: 验证节点审批人员记录创建
- [ ] TC-049: 验证初始状态为0（待审批）
- [ ] TC-050: 验证step为1
- [ ] TC-051: 不存在的模板ID，返回错误
- [ ] TC-052: 禁用的模板，返回错误
- [ ] TC-053: 缺少必填字段，返回400
- [ ] TC-054: 未授权访问，返回401

**预期结果:**
- 创建审批流实例及第一节点
- 返回approval_flow_management_id
- 状态码200

#### 3.2 GET /api/approval_flow_management/initiated_flows - 我发起的审批流
**请求示例:**
```bash
curl -X GET "http://localhost:5001/api/approval_flow_management/initiated_flows?status=0&page=1&size=10" \
  -H "Authorization: Bearer {token}"
```

**测试用例:**
- [ ] TC-055: 获取当前用户发起的所有审批流
- [ ] TC-056: 按状态筛选（待审批）
- [ ] TC-057: 按状态筛选（已通过）
- [ ] TC-058: 按状态筛选（已驳回）
- [ ] TC-059: 分页参数测试
- [ ] TC-060: 返回数据包含模板名称
- [ ] TC-061: 未授权访问，返回401

**预期结果:**
- 返回发起人为当前用户的审批流列表
- 状态码200

#### 3.3 GET /api/approval_flow_management/pending_flows - 待我审批的
**请求示例:**
```bash
curl -X GET "http://localhost:5001/api/approval_flow_management/pending_flows?page=1&size=10" \
  -H "Authorization: Bearer {token}"
```

**测试用例:**
- [ ] TC-062: 获取当前用户待审批的所有任务
- [ ] TC-063: 只返回result为null的记录
- [ ] TC-064: 返回数据包含发起人、模板名称
- [ ] TC-065: 分页参数测试
- [ ] TC-066: 未授权访问，返回401

**预期结果:**
- 返回待审批任务列表，包含node_case_user_id
- 状态码200

#### 3.4 GET /api/approval_flow_management/completed_flows - 我已审批的
**请求示例:**
```bash
curl -X GET "http://localhost:5001/api/approval_flow_management/completed_flows?page=1&size=10" \
  -H "Authorization: Bearer {token}"
```

**测试用例:**
- [ ] TC-067: 获取当前用户已审批的所有记录
- [ ] TC-068: 包含通过和驳回的记录
- [ ] TC-069: 返回数据包含审批结果、审批时间
- [ ] TC-070: 分页参数测试
- [ ] TC-071: 未授权访问，返回401

**预期结果:**
- 返回已审批记录列表，result非null
- 状态码200

#### 3.5 GET /api/approval_flow_management/copied_flows - 抄送给我的
**请求示例:**
```bash
curl -X GET "http://localhost:5001/api/approval_flow_management/copied_flows?page=1&size=10" \
  -H "Authorization: Bearer {token}"
```

**测试用例:**
- [ ] TC-072: 获取抄送给当前用户的所有记录
- [ ] TC-073: 返回数据包含抄送信息、时间
- [ ] TC-074: 分页参数测试
- [ ] TC-075: 未授权访问，返回401

**预期结果:**
- 返回抄送记录列表
- 状态码200

#### 3.6 POST /api/approval_flow_management/approve - 审批通过
**请求示例:**
```bash
curl -X POST "http://localhost:5001/api/approval_flow_management/approve" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "node_case_user_id": 1
  }'
```

**测试用例（会签节点）:**
- [ ] TC-076: 第一人通过，节点未通过，流程继续
- [ ] TC-077: 所有人通过，节点通过，流转下一节点
- [ ] TC-078: 所有节点通过，审批流完成，状态变为10
- [ ] TC-079: 审批流完成后，创建抄送记录

**测试用例（或签节点）:**
- [ ] TC-080: 任意一人通过，节点通过，删除其他待审批人
- [ ] TC-081: 节点通过后，流转下一节点

**测试用例（边界情况）:**
- [ ] TC-082: 重复审批，返回错误"该审批已处理"
- [ ] TC-083: 不存在的node_case_user_id，返回404
- [ ] TC-084: 审批其他用户的任务，返回权限错误
- [ ] TC-085: 未授权访问，返回401

**预期结果:**
- 更新审批人员记录result=0
- 根据节点类型判断是否流转
- 状态码200

#### 3.7 POST /api/approval_flow_management/reject - 审批驳回
**请求示例:**
```bash
curl -X POST "http://localhost:5001/api/approval_flow_management/reject" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "node_case_user_id": 1
  }'
```

**测试用例（会签节点）:**
- [ ] TC-086: 任意一人驳回，节点驳回，删除其他待审批人
- [ ] TC-087: 节点驳回后，审批流状态变为20

**测试用例（或签节点）:**
- [ ] TC-088: 第一人驳回，节点未驳回，流程继续
- [ ] TC-089: 所有人驳回，节点驳回，审批流状态变为20

**测试用例（边界情况）:**
- [ ] TC-090: 重复驳回，返回错误"该审批已处理"
- [ ] TC-091: 不存在的node_case_user_id，返回404
- [ ] TC-092: 驳回其他用户的任务，返回权限错误
- [ ] TC-093: 未授权访问，返回401

**预期结果:**
- 更新审批人员记录result=1
- 根据节点类型判断是否终止流程
- 状态码200

---

## 复杂业务场景测试

### 场景1: 单节点会签审批流（全部通过）
1. 创建模板：1个会签节点，3个审批人
2. 发起审批流
3. 第1个审批人通过 → 验证节点状态未变化
4. 第2个审批人通过 → 验证节点状态未变化
5. 第3个审批人通过 → 验证节点通过，审批流完成（状态10），创建抄送记录

### 场景2: 单节点会签审批流（一人驳回）
1. 创建模板：1个会签节点，3个审批人
2. 发起审批流
3. 第1个审批人通过
4. 第2个审批人驳回 → 验证节点驳回，审批流状态20，删除第3人待审批记录

### 场景3: 单节点或签审批流（一人通过）
1. 创建模板：1个或签节点，3个审批人
2. 发起审批流
3. 第1个审批人通过 → 验证节点通过，审批流完成（状态10），删除其他待审批记录

### 场景4: 单节点或签审批流（全部驳回）
1. 创建模板：1个或签节点，3个审批人
2. 发起审批流
3. 第1个审批人驳回 → 验证节点未驳回
4. 第2个审批人驳回 → 验证节点未驳回
5. 第3个审批人驳回 → 验证节点驳回，审批流状态20

### 场景5: 多节点审批流（会签 → 或签 → 会签）
1. 创建模板：
   - 节点1(会签): 审批人A、B
   - 节点2(或签): 审批人C、D、E
   - 节点3(会签): 审批人F、G
2. 发起审批流
3. A通过 → 验证step=1，未流转
4. B通过 → 验证节点1通过，step=2，创建节点2实例
5. D通过 → 验证节点2通过，step=3，创建节点3实例，删除C、E待审批
6. F通过 → 验证step=3，未流转
7. G通过 → 验证节点3通过，审批流完成（状态10），创建抄送记录

### 场景6: 多节点审批流（第二节点驳回）
1. 创建模板：节点1(会签: A、B)、节点2(会签: C、D)
2. 发起审批流
3. A通过、B通过 → 节点1通过，流转节点2
4. C驳回 → 验证节点2驳回，审批流状态20，删除D待审批

### 场景7: 启用模板自动禁用同类型其他模板
1. 创建类型"请假审批"
2. 创建模板1（请假审批，启用）
3. 创建模板2（请假审批，禁用）
4. 启用模板2 → 验证模板1自动变为禁用

### 场景8: 并发审批冲突测试
1. 创建模板：1个或签节点，2个审批人A、B
2. 发起审批流
3. A、B同时调用approve接口 → 验证只有一人成功，使用行锁避免冲突

---

## 数据库验证检查点

每个测试用例执行后，应验证以下数据库状态:

### 表: approval_flow_management
- [ ] initiator_id 正确
- [ ] approval_flow_template_id 正确
- [ ] status 状态正确 (0/10/20/99)
- [ ] step 步骤正确递增
- [ ] create_time、update_time 正确

### 表: approval_node_case
- [ ] approval_flow_management_id 关联正确
- [ ] node_id 关联正确
- [ ] type 类型正确 (0会签/1或签)
- [ ] result 结果正确 (null/0/1)

### 表: approval_node_case_user
- [ ] approval_node_case_id 关联正确
- [ ] useraccount_id 审批人正确
- [ ] result 结果正确 (null/0/1)
- [ ] approval_time 审批时间正确

### 表: approval_copy_useraccount_case
- [ ] approval_flow_management_id 关联正确
- [ ] useraccount_id 抄送人正确
- [ ] info 抄送信息正确
- [ ] create_time 创建时间正确

---

## 性能测试

- [ ] 并发发起审批流 (100个请求/秒)
- [ ] 并发审批 (50个请求/秒)
- [ ] 大数据量查询 (10000条记录分页查询)
- [ ] 复杂模板创建 (10个节点，每节点10个审批人)

---

## 安全测试

- [ ] 未登录访问所有接口，验证返回401
- [ ] Token过期，验证返回401
- [ ] 审批其他用户的任务，验证权限校验
- [ ] SQL注入测试（name、title、info字段）
- [ ] XSS测试（返回数据HTML转义）
- [ ] CSRF测试（验证token机制）

---

## 测试工具

推荐使用以下工具:
1. **Postman** - API手动测试和Collection管理
2. **curl** - 命令行快速测试
3. **MySQL Workbench** - 数据库状态验证
4. **JMeter** - 性能和并发测试

---

## 测试报告模板

```
测试日期: YYYY-MM-DD
测试人员: XXX
测试环境: 开发环境

| 模块 | 通过 | 失败 | 跳过 | 通过率 |
|------|------|------|------|--------|
| 审批流类型 | 0/19 | 0 | 0 | 0% |
| 审批流模板 | 0/26 | 0 | 0 | 0% |
| 审批流实例 | 0/48 | 0 | 0 | 0% |
| **总计** | **0/93** | **0** | **0** | **0%** |

失败用例详情:
- TC-XXX: 描述...
  - 预期: ...
  - 实际: ...
  - 问题: ...
```
