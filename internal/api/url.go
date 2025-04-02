package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sunvdong/urlshortener/internal/model"
)

type URLService interface {
	CreateURL(ctx context.Context, req model.CreateURLRequest) (*model.CreateURLResponse, error)
	GetURLByCode(ctx context.Context, code string) (string, error)
}

type URLHandler struct {
	urlService URLService
}

func NewURLHandler(s URLService) *URLHandler {
	return &URLHandler{
		urlService: s,
	}
}

// POST /api/url original_url, custom_code, duration ==> shortURL, 过期时间
func (h *URLHandler) CreateURL(c echo.Context) error {
	// 提取数据
	var req model.CreateURLRequest
	if e := c.Bind(&req); e != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	}
	// 验证数据格式
	if e := c.Validate(req); e != nil {
		return echo.NewHTTPError(http.StatusBadRequest, e.Error())
	}
	// 调用业务函数
	resp, err := h.urlService.CreateURL(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 返回响应
	return c.JSON(http.StatusCreated, resp)
}

// GET /:code ==> 把短URL重定向到长URL
func (h *URLHandler) RedirectURL(c echo.Context) error {
	// 取出来 code
	shortCode := c.Param("code")

	// shortcode ==> url 调取业务函数
	originalUrl, err := h.urlService.GetURLByCode(c.Request().Context(), shortCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusPermanentRedirect, originalUrl)
}
