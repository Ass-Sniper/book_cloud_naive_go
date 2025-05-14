package errors

const (
	// 成功码高位段统一为 0x80000000 开头
	SUCCESS      = 0x80000000 // 通用成功
	AUTH_SUCCESS = 0x81000000 // Auth 成功（子模块）
)

// Module: Authentication
// Authentication相关的错误码，前缀 "AUTH"
const (
	AUTH_ERR_USER_ALREADY_EXISTS = 0x01010001 // 用户名已存在
	AUTH_ERR_PASSWORD_MISMATCH   = 0x01010002 // 密码不匹配
	AUTH_ERR_INVALID_CREDENTIALS = 0x01010003 // 无效用户名或密码
	AUTH_ERR_USER_NOT_FOUND      = 0x01010004 // 用户不存在
	AUTH_ERR_SESSION_INVALID     = 0x01010005 // 用户会话无效
)

// Module: Key-Value Store
// KV存储相关的错误码，前缀 "KV"
const (
	KV_ERR_KEY_NOT_FOUND    = 0x02020001 // 键值存储不存在或已过期
	KV_ERR_INVALID_TTL      = 0x02020002 // 无效的TTL
	KV_ERR_SET_VALUE_FAILED = 0x02020003 // 设置键值失败
	KV_ERR_DELETE_FAILED    = 0x02020004 // 删除键值失败
	KV_ERR_LIST_KEYS_FAILED = 0x02020005 // 列出键值失败
	KV_ERR_INVALID_JSON     = 0x02020006 // 无效JSON数据
)

// Module: Session Management
// 会话管理相关的错误码，前缀 "SESSION"
const (
	SESSION_ERR_INVALID_TOKEN           = 0x03030001 // 会话令牌无效
	SESSION_ERR_SESSION_CREATION_FAILED = 0x03030002 // 创建会话失败
)

// Module: Template Management
// 模板管理相关的错误码，前缀 "TEMPLATE"
const (
	TEMPLATE_ERR_LOAD_FAILED       = 0x04040001 // 模板加载失败
	TEMPLATE_ERR_FILE_READ_FAILED  = 0x04040002 // 文件读取失败
	TEMPLATE_ERR_FILE_WRITE_FAILED = 0x04040003 // 文件写入失败
)

// Module: Unknown Errors
// 其它未知错误的通用错误码
const (
	UNKNOWN_ERR_GENERAL = 0x09090001 // 未知错误
)
