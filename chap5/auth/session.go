package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

var (
	// map的两种初始化方式：make 或字面量 {} 初始化
	users     = make(map[string]string) // 存储用户名和加密后的密码。make 方式初始化：创建空 map，预分配空间、性能稍高，便于后续添加元素
	sessions  = map[string]string{}     // 存储会话令牌和用户名。字面量 {} 初始化：简洁表达，适合初始化时有固定值
	mu        sync.Mutex
	usersFile = "users.txt" // 用户信息存储文件
)

// 加载用户信息
func LoadUsers() error {
	mu.Lock()
	defer mu.Unlock()

	// 读取用户信息文件
	data, err := os.ReadFile(usersFile)
	if err != nil {
		return fmt.Errorf("failed to read user file: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		// 跳过空行
		if line == "" {
			continue
		}
		// 每行是用户名:加密后的密码
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		users[parts[0]] = parts[1]
	}

	return nil
}

// SaveUsers 将内存中的用户凭证数据持久化到文件
// 函数特性:
//   - 线程安全：使用互斥锁保证并发访问的安全性
//   - 数据格式：将数据序列化为 "username:passwordHash" 的文本格式
//
// 返回值说明:
//   - error: 当文件写入失败时返回错误信息，成功时返回nil
func SaveUsers() error {
	// 同步控制：在数据序列化和写入期间保持互斥锁
	mu.Lock()
	defer mu.Unlock()

	// 数据序列化：将map结构转换为文本行格式
	var data []string
	for username, passwordHash := range users {
		data = append(data, fmt.Sprintf("%s:%s", username, passwordHash))
	}

	// 文件持久化：将数据写入用户文件，设置文件权限为0644 (-rw-r--r--)
	return os.WriteFile(usersFile, []byte(strings.Join(data, "\n")), 0644)
}

/*
HashPassword 使用bcrypt算法对明文密码进行安全哈希处理

参数:

	password string - 需要加密的明文密码字符串

返回值:

	string - 加密后的哈希字符串(包含算法版本、cost值和盐值)
	error  - 加密过程中遇到的错误，包括无效密码或cost值超出范围等情况

重要实现细节:
 1. 采用bcrypt.DefaultCost(当前默认值为10)平衡安全性与计算性能
 2. 自动生成随机盐值，确保相同密码每次加密结果不同
 3. 返回的哈希字符串包含完整的加密参数，可直接用于密码验证
*/
func HashPassword(password string) (string, error) {
	// 使用bcrypt标准方法生成密码哈希，自动处理盐值生成和cost设置
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// 处理潜在错误：包括空密码、超过72字节的密码或无效的cost参数
		return "", fmt.Errorf("failed to hash password: %v", err)
	}

	// 将字节数组转换为标准哈希字符串格式返回
	return string(hash), nil
}

// ValidateUser 验证用户提供的用户名和密码是否有效
// 参数:
//   - username: 待验证的用户名字符串
//   - password: 待验证的明文密码字符串
//
// 返回值:
//   - bool: 验证结果，true表示验证通过，false表示验证失败
func ValidateUser(username, password string) bool {
	// 并发控制：加互斥锁保证线程安全
	mu.Lock()
	defer mu.Unlock()

	// 在用户数据集中查找指定用户
	// 若用户名不存在则立即返回验证失败
	passHash, ok := users[username]
	if !ok {
		return false
	}

	// 使用bcrypt算法进行安全密码比对
	// 比较传入的明文密码与存储的哈希值是否匹配
	err := bcrypt.CompareHashAndPassword([]byte(passHash), []byte(password))
	return err == nil
}

// CreateSession 生成新的会话令牌并关联用户名
// 参数:
//   - username: 需要关联会话的用户名
//
// 返回值:
//   - string: 生成的16字节随机十六进制会话令牌
func CreateSession(username string) string {
	// 生成16字节加密安全随机数据作为令牌基础值
	b := make([]byte, 16)
	rand.Read(b)
	token := hex.EncodeToString(b)

	// 使用互斥锁保证并发安全的会话存储
	mu.Lock()
	defer mu.Unlock()
	sessions[token] = username
	return token
}

// GetUsernameByToken 根据会话令牌查找对应的用户名
// 参数:
//
//	token: 会话令牌字符串，用于唯一标识用户会话
//
// 返回值:
//
//	string: 找到的用户名，若不存在则为空字符串
//	bool:  指示是否成功找到对应用户的布尔值
//
// 注意：该函数使用互斥锁保证并发安全访问会话数据
func GetUsernameByToken(token string) (string, bool) {
	// 使用互斥锁保证对共享会话数据的原子操作
	mu.Lock()
	defer mu.Unlock()

	// 从会话存储中直接查询并返回结果
	u, ok := sessions[token]
	return u, ok
}

// RegisterUser 注册新用户并保存到内存和文件
// 参数:
//   - username: 用户名，必须唯一
//   - password: 用户输入的明文密码
//
// 返回值:
//   - error: 注册失败时返回错误原因，包含用户名冲突或保存失败等情况
func RegisterUser(username, password string) error {
	// 检查用户名是否已被注册
	// 通过查询内存中的用户映射判断重复注册
	if _, exists := users[username]; exists {
		return fmt.Errorf("username already exists")
	}

	// 对原始密码进行不可逆加密处理
	// 使用bcrypt等安全哈希算法保护用户密码
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	// 将加密后的密码存入内存缓存
	// 使用全局map实现临时存储，等待持久化操作
	users[username] = hashedPassword

	// 将内存中的用户数据同步到持久化存储
	// 通过文件IO操作确保数据不丢失
	return SaveUsers()
}

// Logout 删除指定 token 对应的会话
func Logout(token string) {
	mu.Lock()
	defer mu.Unlock()
	delete(sessions, token)
}
