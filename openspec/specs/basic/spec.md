# Specification: Basic Data Management

基础数据管理模块，提供性别、年级、学科等字典数据的查询功能。

## Requirements

### Requirement: 获取性别列表
系统 SHALL 提供性别字典数据供前端下拉框使用。

#### Scenario: 获取所有性别
- **WHEN** 已认证用户请求性别列表
- **THEN** 系统返回所有性别选项（男、女）
- **AND** 响应格式为 `{"code": 0, "message": "success", "data": {"sexes": [...]}}`

#### Scenario: 性别按ID排序
- **WHEN** 返回性别列表
- **THEN** 性别按ID升序排列

#### Scenario: 未认证用户访问
- **WHEN** 未提供有效JWT token的用户请求性别列表
- **THEN** 系统返回401错误
- **AND** 错误消息为 "未授权"

### Requirement: 获取启用的年级列表
系统 SHALL 提供启用状态的年级列表供前端下拉框使用。

#### Scenario: 获取启用年级成功
- **WHEN** 已认证用户请求年级列表
- **THEN** 系统返回所有status=0（启用）的年级
- **AND** 响应格式为 `{"code": 0, "message": "success", "data": {"grades": [...]}}`
- **AND** 每个年级包含id和grade字段

#### Scenario: 年级按ID排序
- **WHEN** 返回年级列表
- **THEN** 年级按ID升序排列（一年级、二年级...高三）

#### Scenario: 过滤禁用年级
- **WHEN** 存在status=1（禁用）的年级
- **THEN** 该年级不出现在返回列表中

#### Scenario: 年级字段名称
- **WHEN** 返回年级数据
- **THEN** 年级名称字段为"grade"而非"name"
- **AND** 与数据库列名保持一致

### Requirement: 获取启用的学科列表
系统 SHALL 提供启用状态的学科列表供前端下拉框使用。

#### Scenario: 获取启用学科成功
- **WHEN** 已认证用户请求学科列表
- **THEN** 系统返回所有status=0（启用）的学科
- **AND** 响应格式为 `{"code": 0, "message": "success", "data": {"subjects": [...]}}`
- **AND** 每个学科包含id和subject字段

#### Scenario: 学科按ID排序
- **WHEN** 返回学科列表
- **THEN** 学科按ID升序排列

#### Scenario: 过滤禁用学科
- **WHEN** 存在status=1（禁用）的学科
- **THEN** 该学科不出现在返回列表中

#### Scenario: 学科字段名称
- **WHEN** 返回学科数据
- **THEN** 学科名称字段为"subject"而非"name"
- **AND** 与数据库列名保持一致

### Requirement: 基础数据认证保护
系统 SHALL 要求JWT认证才能访问基础数据接口。

#### Scenario: 有效Token访问成功
- **WHEN** 使用有效JWT token请求基础数据
- **THEN** 系统允许访问并返回数据

#### Scenario: Token过期
- **WHEN** 使用过期JWT token请求基础数据
- **THEN** 系统返回401错误
- **AND** 错误消息为 "token无效或已过期"

#### Scenario: 未提供Token
- **WHEN** 未提供JWT token请求基础数据
- **THEN** 系统返回401错误
- **AND** 错误消息为 "未提供token"

## API Endpoints

### GET /api/sexes
获取所有性别（需要认证）

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200):**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "sexes": [
      {
        "id": 1,
        "name": "男"
      },
      {
        "id": 2,
        "name": "女"
      }
    ]
  }
}
```

### GET /api/grades/active
获取启用的年级（需要认证）

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200):**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "grades": [
      {"id": 1, "grade": "一年级"},
      {"id": 2, "grade": "二年级"},
      {"id": 3, "grade": "三年级"},
      {"id": 4, "grade": "四年级"},
      {"id": 5, "grade": "五年级"},
      {"id": 6, "grade": "六年级"},
      {"id": 7, "grade": "初一"},
      {"id": 8, "grade": "初二"},
      {"id": 9, "grade": "初三"},
      {"id": 10, "grade": "高一"},
      {"id": 11, "grade": "高二"},
      {"id": 12, "grade": "高三"}
    ]
  }
}
```

### GET /api/subjects/active
获取启用的学科（需要认证）

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200):**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "subjects": [
      {"id": 1, "subject": "语文"},
      {"id": 2, "subject": "数学"},
      {"id": 3, "subject": "英语"},
      {"id": 4, "subject": "物理"},
      {"id": 5, "subject": "化学"},
      {"id": 6, "subject": "生物"},
      {"id": 7, "subject": "历史"},
      {"id": 8, "subject": "地理"},
      {"id": 9, "subject": "政治"},
      {"id": 10, "subject": "音乐"},
      {"id": 11, "subject": "美术"},
      {"id": 12, "subject": "体育"}
    ]
  }
}
```

## Data Models

### Sex
```go
type Sex struct {
    ID   uint   `gorm:"column:id;primaryKey"`
    Name string `gorm:"column:name;size:10;not null"`
}
```

**表名**: `sex`

**字段说明**:
- `id`: 主键，自增
- `name`: 性别名称（男、女）

### Grade
```go
type Grade struct {
    ID         uint      `gorm:"column:id;primaryKey"`
    Name       string    `gorm:"column:name;size:50;not null"`
    Status     int       `gorm:"column:status;default:0"`
    CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
    UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime"`
}
```

**表名**: `grade`

**字段说明**:
- `id`: 主键，自增
- `name`: 年级名称
- `status`: 状态（0=启用，1=禁用）
- `create_time`: 创建时间
- `update_time`: 更新时间

**JSON字段映射**:
- 返回JSON时`name`字段映射为`grade`

### Subject
```go
type Subject struct {
    ID         uint      `gorm:"column:id;primaryKey"`
    Subject    string    `gorm:"column:subject;size:50;not null"`
    Status     int       `gorm:"column:status;default:0"`
    CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
    UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime"`
}
```

**表名**: `subject`

**字段说明**:
- `id`: 主键，自增
- `subject`: 学科名称
- `status`: 状态（0=启用，1=禁用）
- `create_time`: 创建时间
- `update_time`: 更新时间

## Business Rules

1. **性别数据**: 固定为"男"和"女"两个选项
2. **年级范围**: 一年级至高三，共12个年级
3. **学科范围**: 包含语文、数学等12个学科
4. **状态过滤**: 只返回status=0（启用）的年级和学科
5. **排序规则**: 所有列表按ID升序排列
6. **字段命名**: 必须与数据库列名保持一致
7. **认证要求**: 所有接口都需要JWT认证

## Database Initialization

初始化脚本位于: `scripts/init_basic_data.sql`

### 性别数据
```sql
INSERT INTO `sex` (`id`, `name`) VALUES
(1, '男'),
(2, '女');
```

### 年级数据
```sql
INSERT INTO `grade` (`id`, `name`, `status`) VALUES
(1, '一年级', 0),
(2, '二年级', 0),
(3, '三年级', 0),
(4, '四年级', 0),
(5, '五年级', 0),
(6, '六年级', 0),
(7, '初一', 0),
(8, '初二', 0),
(9, '初三', 0),
(10, '高一', 0),
(11, '高二', 0),
(12, '高三', 0);
```

### 学科数据
```sql
INSERT INTO `subject` (`id`, `subject`, `status`) VALUES
(1, '语文', 0),
(2, '数学', 0),
(3, '英语', 0),
(4, '物理', 0),
(5, '化学', 0),
(6, '生物', 0),
(7, '历史', 0),
(8, '地理', 0),
(9, '政治', 0),
(10, '音乐', 0),
(11, '美术', 0),
(12, '体育', 0);
```

## Frontend Integration

### 前端使用场景

1. **学生管理**:
   - 性别下拉框：选择学生性别
   - 年级下拉框：选择学生所在年级

2. **教练管理**:
   - 性别下拉框：选择教练性别
   - 学科下拉框：选择教练授课学科

3. **数据筛选**:
   - 按性别筛选学生/教练
   - 按年级筛选学生
   - 按学科筛选教练

### 前端调用示例
```javascript
// 获取性别列表
const response = await axios.get('/api/sexes', { withCredentials: true });
this.sexes = response.data.data.sexes;

// 获取年级列表
const response = await axios.get('/api/grades/active', { withCredentials: true });
this.grades = response.data.data.grades;

// 获取学科列表
const response = await axios.get('/api/subjects/active', { withCredentials: true });
this.subjects = response.data.data.subjects;
```

## Performance Considerations

1. **数据量**: 基础数据量极小（性别2条、年级12条、学科12条）
2. **查询性能**: 全表扫描性能开销可忽略
3. **缓存策略**: 当前未实现，可考虑添加内存缓存（数据基本不变）
4. **数据库索引**: 主键索引已足够，无需额外索引

## Future Enhancements

1. **管理接口**: 添加年级和学科的CRUD管理接口
2. **缓存优化**: 添加内存缓存减少数据库查询
3. **动态配置**: 支持在管理界面动态添加年级和学科
4. **国际化**: 支持多语言的年级和学科名称
5. **扩展字段**: 为年级和学科添加更多属性（如描述、图标等）
