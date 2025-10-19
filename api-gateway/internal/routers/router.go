package routers

import (
	"api-gateway/internal/config"
	"api-gateway/internal/middleware"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, cfg *config.Config) {
	usersProxy := setupProxy(cfg.UsersServiceURL)
	r.Any("/v1/users/*path", middleware.JWTAuth(), middleware.RequestID(), middleware.Logging(), middleware.CORS(), usersProxy)

	ordersProxy := setupProxy(cfg.OrdersServiceURL)
	r.Any("/v1/orders/*path", middleware.JWTAuth(), middleware.RequestID(), middleware.Logging(), middleware.CORS(), ordersProxy)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "API Gateway is running"})
	})
}

func setupProxy(target string) gin.HandlerFunc {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.Director = func(req *http.Request) {
		req.Host = url.Host
		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host
		req.URL.Path = req.URL.Path[1:]
		req.Header.Set("X-Request-ID", req.Header.Get("X-Request-ID"))
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "")
		}
	}

	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {
		writer.WriteHeader(http.StatusBadGateway)
		writer.Write([]byte("Service unavailable"))
	}

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
