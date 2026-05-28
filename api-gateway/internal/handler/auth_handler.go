package handler

import (
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/labstack/echo/v5"
)

type AuthHandler struct {
	proxy *httputil.ReverseProxy
}

func NewAuthHandler() *AuthHandler {
	target, _ := url.Parse(os.Getenv("AUTH_SERVICE_URL"))
	return &AuthHandler{proxy: httputil.NewSingleHostReverseProxy(target)}
}

func (h *AuthHandler) Proxy(c *echo.Context) error {
	h.proxy.ServeHTTP(c.Response(), c.Request())
	return nil
}