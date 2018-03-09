package middleware

// 中间件
type Middleware interface {
	Use(string, interface{})
	Exec(string, func(v interface{}, next func()))
}
