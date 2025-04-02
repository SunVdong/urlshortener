package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/sunvdong/urlshortener/internal/model"
	"github.com/sunvdong/urlshortener/internal/repo"
)

type ShortCodeGenerator interface {
	GenerateShortCode() string
}

type Cacher interface {
	SetURL(ctx context.Context, url repo.Url) error
	GetURL(ctx context.Context, shortCode string) (*repo.Url, error)
}

type URLService struct {
	querier            repo.Querier
	shortCodeGenerator ShortCodeGenerator
	defaultDuration    time.Duration
	baseURL            string
	cache              Cacher
}

func NewURLService(db *sql.DB, shortCodeGenerator ShortCodeGenerator,
	defaultDuration time.Duration, cache Cacher, baseURL string) *URLService {
	return &URLService{
		querier:            repo.New(db),
		shortCodeGenerator: shortCodeGenerator,
		defaultDuration:    defaultDuration,
		baseURL:            baseURL,
		cache:              cache,
	}
}

func (s *URLService) CreateURL(ctx context.Context, req model.CreateURLRequest) (*model.CreateURLResponse, error) {
	var shortCode string
	var isCustom bool
	var expiredAt time.Time

	if req.CustomCode != "" {
		isAvailable, err := s.querier.IsShortCodeAvailable(ctx, req.CustomCode)
		if err != nil {
			return nil, err
		}
		if !isAvailable {
			return nil, fmt.Errorf("短链接已经存在")
		}
		shortCode = req.CustomCode
		isCustom = true
	} else {
		code, err := s.getShortCode(ctx, 0)
		if err != nil {
			return nil, err
		}
		shortCode = code
	}

	if req.Duration == nil {
		expiredAt = time.Now().Add(s.defaultDuration)
	} else {
		expiredAt = time.Now().Add(time.Hour * time.Duration(*req.Duration))
	}

	// 插入数据库
	url, err := s.querier.CreateURL(ctx, repo.CreateURLParams{
		OriginalUrl: req.OriginalUrl,
		ShortCode:   shortCode,
		IsCustom:    isCustom,
		ExpiredAt:   expiredAt,
	})
	if err != nil {
		return nil, err
	}

	// 插入缓存
	if err := s.cache.SetURL(ctx, url); err != nil {
		return nil, err
	}

	return &model.CreateURLResponse{
		ShortURL:  s.baseURL + "/" + url.ShortCode,
		ExpiredAt: url.ExpiredAt,
	}, nil
}

func (s *URLService) GetURLByCode(ctx context.Context, code string) (string, error) {
	// 访问缓存
	url, err := s.cache.GetURL(ctx, code)
	if err != nil {
		return "", err
	}
	if url != nil {
		return url.OriginalUrl, nil
	}
	// 不存在 ==> 访问数据库
	url2, err := s.querier.GetUrlByShortCode(ctx, code)
	if err != nil {
		return "", err
	}
	// 存入缓存
	if err := s.cache.SetURL(ctx, url2); err != nil {
		return "", err
	}

	return url2.OriginalUrl, nil
}

func (s *URLService) DeleteURL(ctx context.Context) error {
	return s.querier.DeleUrlExpired(ctx)
}

func (s *URLService) getShortCode(ctx context.Context, n int) (string, error) {
	if n > 5 {
		return "", errors.New("重试过多")
	}
	shortUrl := s.shortCodeGenerator.GenerateShortCode()

	isAvailable, err := s.querier.IsShortCodeAvailable(ctx, shortUrl)
	if err != nil {
		return "", err
	}

	if !isAvailable {
		return s.getShortCode(ctx, n+1)
	}

	return shortUrl, nil
}
