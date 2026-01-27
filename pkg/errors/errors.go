package errors

import "errors"

// 预定义的业务错误
var (
	// 通用错误
	ErrInternal      = errors.New("服务器内部错误")
	ErrInvalidParams = errors.New("参数错误")
	ErrNotFound      = errors.New("记录不存在")

	// 认证相关
	ErrUnauthorized     = errors.New("未登录")
	ErrInvalidToken     = errors.New("无效的 Token")
	ErrTokenExpired     = errors.New("Token 已过期")
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrAccountDisabled  = errors.New("账号已被禁用")

	// 权限相关
	ErrForbidden       = errors.New("无权限访问")
	ErrPermissionDenied = errors.New("权限不足")

	// 用户相关
	ErrUserNotFound     = errors.New("用户不存在")
	ErrUsernameExists   = errors.New("用户名已存在")
	ErrPhoneExists      = errors.New("手机号已存在")

	// 数据相关
	ErrRecordExists     = errors.New("记录已存在")
	ErrRecordNotFound   = errors.New("记录不存在")
	ErrRecordInUse      = errors.New("记录正在使用中，无法删除")

	// 业务相关
	ErrInvalidStatus    = errors.New("无效的状态")
	ErrOperationFailed  = errors.New("操作失败")
)

// AppError 应用错误
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// NewAppError 创建应用错误
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// 常用错误构造函数
func BadRequest(message string) *AppError {
	return &AppError{Code: 400, Message: message}
}

func Unauthorized(message string) *AppError {
	return &AppError{Code: 401, Message: message}
}

func Forbidden(message string) *AppError {
	return &AppError{Code: 403, Message: message}
}

func NotFound(message string) *AppError {
	return &AppError{Code: 404, Message: message}
}

func InternalServer(message string) *AppError {
	return &AppError{Code: 500, Message: message}
}
