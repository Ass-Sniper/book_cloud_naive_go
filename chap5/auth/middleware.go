package auth

import (
	"net/http"

	"github.com/kay/kvstore/logger"
)

// AuthMiddleware HTTP身份验证中间件，用于验证请求中的会话Cookie
// 参数说明:
//   - next: 后续处理函数，类型为http.Handler
//
// 返回值:
//   - http.Handler: 返回新的处理器，执行身份验证后传递给后续处理
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查会话Cookie是否存在
		cookie, err := r.Cookie("session")
		if err != nil {
			logger.Log.WithFields(map[string]interface{}{
				"module":    "auth",
				"goroutine": logger.GetGID(),
				"reason":    "no session cookie",
				"remote":    r.RemoteAddr,
				"path":      r.URL.Path,
			}).Warn("unauthenticated access")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// 验证会话Token有效性
		if _, ok := GetUsernameByToken(cookie.Value); !ok {
			logger.Log.WithFields(map[string]interface{}{
				"module":    "auth",
				"goroutine": logger.GetGID(),
				"reason":    "invalid session token",
				"remote":    r.RemoteAddr,
				"path":      r.URL.Path,
			}).Warn("unauthenticated access")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// 验证通过，继续后续处理流程
		next.ServeHTTP(w, r)
	})
}
