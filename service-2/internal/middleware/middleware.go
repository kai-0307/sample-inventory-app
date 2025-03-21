package middleware

import "net/http"

// JSONMiddleware はすべてのレスポンスにJSONヘッダーを設定する
func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
