# activity-template 规范增量

## 新增需求

### 需求： 活动模板删除

系统MUST支持删除活动模板。在删除前，系统MUST检查是否有关联的活动。如果存在关联活动，系统MUST禁止删除并返回明确的错误信息。

- **标识符**: `activity-template-delete`
- **优先级**: P0

#### 场景：删除有关联活动的模板

**前置条件**:
- 活动模板存在
- 存在关联该模板的活动记录
- 用户拥有 `delete_activity_template` 权限

**操作**:
```http
DELETE /api/activity-templates/1
Authorization: Bearer {token}
```

**预期结果**:
- 返回 HTTP 400 Bad Request
- 响应格式：
  ```json
  {
    "error": "该模板有关联活动，无法删除"
  }
  ```
- 活动模板未被删除
- 数据库中的数据保持不变

**验证规则**:
- 删除前MUST查询 `activity` 表中 `template_id` 匹配的记录数
- 如果记录数 > 0，MUST拒绝删除
- 错误消息MUST为："该模板有关联活动，无法删除"

---


