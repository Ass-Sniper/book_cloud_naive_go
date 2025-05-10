package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func init() {
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	Log.SetLevel(logrus.DebugLevel)
	Log.AddHook(&GIDHook{})
}

type GIDHook struct{}

func (h *GIDHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *GIDHook) Fire(entry *logrus.Entry) error {
	gid := GetGID()
	pc, file, _, ok := runtime.Caller(8)
	module := "unknown"
	if ok {
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			module = extractModuleName(fn.Name(), file)
		}
	}
	entry.Data["goroutine"] = gid
	entry.Data["module"] = module
	return nil
}

// 辅助函数：提取模块名
func extractModuleName(funcName, file string) string {
	parts := strings.Split(funcName, "/")
	last := parts[len(parts)-1]
	if idx := strings.Index(last, "."); idx != -1 {
		return last[:idx]
	}
	fileParts := strings.Split(file, "/")
	if len(fileParts) > 1 {
		return fileParts[len(fileParts)-2]
	}
	return "unknown"
}

// 在 Go 中，将 标识符从私有变为导出 的方式非常简单：
// 只需将函数、变量、常量、类型等的首字母从小写改为大写。
// GetGID 获取当前 goroutine 的 ID
//
// 返回值:
//
//	int - 当前 goroutine 的唯一标识符。注意：
//	      1. 该实现依赖 runtime 包的栈信息格式，不同 Go 版本可能存在兼容性问题
//	      2. 获取 goroutine ID 在 Go 编程规范中不被推荐，仅建议用于调试等特殊场景
//	      3. 该方法有性能损耗，不宜高频调用
func GetGID() int {
	// 通过截取运行时栈信息的前64字节来解析 GID
	// runtime.Stack 第二个参数 false 表示不打印所有 goroutine 的栈信息
	var buf [64]byte
	n := runtime.Stack(buf[:], false)

	// 使用格式化扫描从栈信息中提取 GID
	// 栈信息格式示例："goroutine 1 [running]:..."
	var gid int
	_, _ = fmt.Sscanf(string(buf[:n]), "goroutine %d ", &gid)
	return gid
}
