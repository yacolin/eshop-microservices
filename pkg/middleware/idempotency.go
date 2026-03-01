package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"

	"eshop-microservices/pkg/database"
)

// Idempotency 幂等性检查中间件
// 适用于所有需要幂等性保障的接口
func Idempotency() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 RequestID
		requestID := getRequestID(c)
		if requestID == "" {
			// 没有 RequestID，继续处理
			c.Next()
			return
		}

		// 构建 Redis 键
		key := "idempotency:" + requestID
		ctx := database.GetCtx()
		client := database.GetClient()

		// 检查是否存在重复请求
		existingResponse, err := client.Get(ctx, key).Result()
		if err == nil {
			// 找到重复请求，返回缓存的响应
			var cachedResp cachedResponse
			if json.Unmarshal([]byte(existingResponse), &cachedResp) == nil {
				// 设置状态码
				c.Writer.WriteHeader(cachedResp.StatusCode)
				// 设置响应头
				for key, value := range cachedResp.Headers {
					c.Writer.Header().Set(key, value)
				}
				// 写入响应体
				c.Writer.Write([]byte(cachedResp.Body))
				// 终止后续处理
				c.Abort()
				return
			}
		}

		// 捕获响应
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 继续处理请求
		c.Next()

		// 缓存响应
		if c.Writer.Status() < 400 { // 只缓存成功的响应
			resp := cachedResponse{
				StatusCode: c.Writer.Status(),
				Headers:    make(map[string]string),
				Body:       blw.body.String(),
			}

			// 复制响应头
			for key, values := range c.Writer.Header() {
				if len(values) > 0 {
					resp.Headers[key] = values[0]
				}
			}

			// 序列化响应
			respBytes, err := json.Marshal(resp)
			if err == nil {
				// 缓存响应，过期时间 24 小时
				client.Set(ctx, key, respBytes, 24*time.Hour)
			}
		}
	}
}

// getRequestID 从请求中获取 RequestID
func getRequestID(c *gin.Context) string {
	// 优先从请求头获取
	requestID := c.GetHeader("X-Request-ID")
	if requestID != "" {
		return requestID
	}

	// 从请求体获取（适用于 POST/PUT 请求）
	if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
		// 读取请求体
		body, err := io.ReadAll(c.Request.Body)
		if err == nil {
			// 重置请求体，以便后续处理
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			// 尝试解析 JSON
			var req map[string]interface{}
			if json.Unmarshal(body, &req) == nil {
				if rid, ok := req["request_id"].(string); ok {
					return rid
				}
			}
		}
	}

	// 从查询参数获取
	requestID = c.Query("request_id")
	return requestID
}

// cachedResponse 缓存的响应结构
type cachedResponse struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

// bodyLogWriter 用于捕获响应体
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
