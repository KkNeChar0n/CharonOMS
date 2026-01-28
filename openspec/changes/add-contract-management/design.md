# Design: 合同管理模块技术设计

## Context

合同管理模块需要从原Python Flask项目迁移到Go项目。该模块包含合同的完整生命周期管理，特别需要处理"中止合作"功能中的文件上传（termination_agreement）。

### 背景
- 原项目使用Flask处理文件上传
- 需要在Go中实现相同的文件上传功能
- 需要保持与原项目的API兼容性
- 需要考虑文件存储位置、安全性、命名规则

### 约束
- 前端代码不修改，必须兼容现有API
- 文件上传需要在HTTP multipart/form-data中处理
- 需要支持多种文件格式（pdf, doc, docx, jpg, png）
- 文件大小限制在10MB以内

## Goals / Non-Goals

### Goals
- 实现与原项目完全兼容的合同管理API
- 实现安全的文件上传功能
- 支持合同状态的正确流转
- 保持DDD架构的清晰分层

### Non-Goals
- 不实现付款状态的自动更新（由未来的订单收款模块负责）
- 不实现合同审批流程（由未来的审批流程模块负责）
- 不实现合同编辑和删除功能（原项目也没有）
- 不实现文件的在线预览功能（仅存储路径）

## Decisions

### Decision 1: 文件上传实现方式

**选择**: 使用Gin的`c.FormFile()`处理multipart/form-data文件上传

**理由**:
- Gin提供了内置的文件上传支持，简单易用
- 与原Flask的`request.files`功能类似
- 无需引入额外的第三方库
- 性能良好，适合中小规模文件

**实现细节**:
```go
file, err := c.FormFile("termination_agreement")
if err != nil {
    return err
}

// 验证文件类型
ext := filepath.Ext(file.Filename)
allowedExts := []string{".pdf", ".doc", ".docx", ".jpg", ".png"}
if !contains(allowedExts, ext) {
    return errors.New("不支持的文件类型")
}

// 验证文件大小
if file.Size > 10*1024*1024 { // 10MB
    return errors.New("文件大小超过限制")
}

// 生成唯一文件名
filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), generateRandomString(8), ext)
filepath := filepath.Join("uploads", "contracts", filename)

// 保存文件
if err := c.SaveUploadedFile(file, filepath); err != nil {
    return err
}
```

### Decision 2: 文件存储位置

**选择**: 存储在项目根目录下的 `uploads/contracts/` 目录

**理由**:
- 与前端静态文件分离，避免混淆
- 便于备份和管理
- 可以通过Nginx配置单独的访问路径
- 与原项目的存储方式保持一致

**目录结构**:
```
CharonOMS/
├── uploads/
│   └── contracts/
│       ├── 1706428800_abc12345.pdf
│       ├── 1706428801_def67890.docx
│       └── ...
```

**注意事项**:
- 需要在应用启动时检查并创建该目录
- 需要配置Nginx/Apache以允许访问该目录
- 生产环境可考虑使用对象存储（OSS/S3）

### Decision 3: 文件命名规则

**选择**: 使用 `时间戳_随机字符串.扩展名` 格式

**理由**:
- 避免文件名冲突
- 时间戳便于排序和定位
- 随机字符串增加安全性
- 保留原文件扩展名

**示例**:
- `1706428800_abc12345.pdf`
- `1706428801_def67890.docx`

### Decision 4: 文件路径存储格式

**选择**: 存储相对路径（从项目根目录开始）

**理由**:
- 便于项目迁移和部署
- 不依赖绝对路径
- 与前端兼容

**示例**:
- 数据库存储: `uploads/contracts/1706428800_abc12345.pdf`
- 前端访问: `http://localhost:5001/uploads/contracts/1706428800_abc12345.pdf`

### Decision 5: 中止合作接口设计

**选择**: 使用单独的PUT接口 `/api/contracts/:id/terminate`

**理由**:
- 与原项目保持一致
- 语义清晰，表明这是一个特殊的状态变更操作
- 与撤销操作（revoke）分离，便于权限控制
- 支持文件上传和状态更新在一个请求中完成

**请求格式**:
```
PUT /api/contracts/:id/terminate
Content-Type: multipart/form-data

termination_agreement: <file>
```

**响应格式**:
```json
{
  "message": "操作成功"
}
```

### Decision 6: 不实现编辑和删除功能

**选择**: 不实现 `PUT /api/contracts/:id` 和 `DELETE /api/contracts/:id` 接口

**理由**:
- 原Python项目中也没有这两个接口
- 合同通常不允许直接编辑（业务规则）
- 合同不允许删除（数据完整性要求）
- 状态变更通过专门的操作（撤销、中止）完成

**替代方案**:
- 如需修改合同信息，通过创建新合同实现
- 如需作废合同，通过撤销或中止操作实现

## Alternatives Considered

### Alternative 1: 使用第三方文件上传库

**考虑**: 使用如 `github.com/disintegration/imaging` 等专门的文件处理库

**拒绝理由**:
- Gin内置功能已经足够
- 避免增加不必要的依赖
- 当前需求简单，不需要复杂的图片处理功能

### Alternative 2: 使用对象存储（OSS/S3）

**考虑**: 直接将文件上传到云对象存储服务

**拒绝理由**:
- 增加了系统复杂度和外部依赖
- 需要配置额外的云服务账号和权限
- 当前文件量不大，本地存储已足够
- 可以在未来需要时再迁移

**保留选项**: 在`infrastructure/storage/`层抽象文件存储接口，便于未来切换到对象存储

### Alternative 3: 将文件上传和状态更新分离为两个接口

**考虑**:
- `/api/contracts/:id/upload` - 仅上传文件
- `/api/contracts/:id/terminate` - 仅更新状态

**拒绝理由**:
- 与原项目API不一致
- 增加了前端调用的复杂度
- 需要处理上传成功但状态更新失败的情况
- 用户体验较差（需要两步操作）

## Risks / Trade-offs

### Risk 1: 文件存储空间

**风险**: 随着合同数量增加，文件占用的磁盘空间会持续增长

**缓解措施**:
- 实施文件大小限制（10MB）
- 定期清理已作废合同的文件（需要新的后台任务）
- 监控磁盘使用情况
- 未来可迁移到对象存储

### Risk 2: 文件安全性

**风险**: 上传的文件可能包含恶意内容

**缓解措施**:
- 限制允许的文件类型（仅pdf, doc, docx, jpg, png）
- 限制文件大小
- 使用随机文件名，避免路径遍历攻击
- 配置Nginx禁止执行上传目录中的脚本
- 考虑添加病毒扫描（未来）

### Risk 3: 并发上传冲突

**风险**: 多个用户同时中止合同可能导致文件名冲突

**缓解措施**:
- 使用时间戳+随机字符串生成唯一文件名
- 冲突概率极低（纳秒级时间戳 + 8位随机字符）
- 如果发生冲突，重试生成新文件名

### Risk 4: 文件路径泄露

**风险**: 文件路径可能暴露服务器信息

**缓解措施**:
- 使用相对路径，不暴露绝对路径
- 文件名不包含敏感信息
- 通过Nginx配置访问权限（仅认证用户可访问）

### Trade-off 1: 简单性 vs 可扩展性

**选择**: 优先选择简单的本地文件存储

**权衡**:
- ✅ 实现简单，无外部依赖
- ✅ 开发和调试容易
- ❌ 不适合大规模部署
- ❌ 不支持CDN加速

**决定**: 当前选择简单性，在需要时可迁移到对象存储

### Trade-off 2: 状态验证严格性 vs 灵活性

**选择**: 严格验证状态流转规则

**权衡**:
- ✅ 确保数据一致性
- ✅ 防止非法操作
- ❌ 可能在某些特殊情况下不够灵活

**决定**: 选择严格验证，特殊情况通过数据库直接修改

## Migration Plan

### Step 1: 实现基础合同管理（不含文件上传）
1. 创建Domain层实体和仓储接口
2. 实现Infrastructure层仓储
3. 实现Application层服务
4. 实现Interface层handler
5. 测试基础CRUD功能

### Step 2: 实现文件上传功能
1. 创建`infrastructure/storage/file_storage.go`
2. 实现文件保存逻辑
3. 集成到中止合作接口
4. 测试文件上传功能

### Step 3: 前端联调测试
1. 使用原前端测试所有接口
2. 验证响应格式兼容性
3. 验证文件上传和下载

### Step 4: 部署和配置
1. 创建uploads目录并设置权限
2. 配置Nginx静态文件访问
3. 测试生产环境文件上传

## Open Questions

1. **文件访问权限**: 是否需要检查用户权限才能访问上传的文件？
   - 建议: 是，通过Nginx配置或Go中间件验证JWT token

2. **文件备份策略**: 如何备份上传的文件？
   - 建议: 与数据库备份一起，定期备份uploads目录

3. **文件删除策略**: 合同删除后是否删除关联的文件？
   - 建议: 暂不删除，由管理员手动清理（合同本身也不支持删除）

4. **文件命名冲突**: 如何处理文件名冲突的极端情况？
   - 建议: 在保存前检查文件是否存在，如存在则重新生成文件名

5. **大文件上传**: 10MB限制是否合适？
   - 建议: 先使用10MB，根据实际使用情况调整
