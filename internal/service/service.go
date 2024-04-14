package service

import (
	"avito-banner/internal/domain"
	"avito-banner/internal/repository/cache"
	"avito-banner/internal/repository/db"
	"context"
	"errors"
)

type Service struct {
	repo  *db.Repo
	cache *cache.Cache
}

func NewService(repo *db.Repo, cache *cache.Cache) *Service {
	return &Service{repo: repo, cache: cache}

}

func (s *Service) GetUserBanner(ctx context.Context, token string, tagId, featureId int, useLastRevision bool) (domain.ResponseUserBanner, error) {
	if !useLastRevision {
		data, err := s.cache.Get(ctx, tagId, featureId)
		if err == nil {
			if data.IsActive || token == "admin_token" {
				return domain.ResponseUserBanner{
					Title: data.Title,
					Text:  data.Text,
					Url:   data.Url,
				}, nil
			}
			return domain.ResponseUserBanner{}, errors.New("banner not active")
		}
	}
	data, err := s.repo.GetUserBanner(ctx, tagId, featureId)
	if err == nil {
		s.cache.Set(ctx, tagId, featureId, data.Title, data.Text, data.Url, data.IsActive)
		if data.IsActive || token == "admin_token" {
			return domain.ResponseUserBanner{
				Title: data.Title,
				Text:  data.Text,
				Url:   data.Url,
			}, nil
		}
		return domain.ResponseUserBanner{}, errors.New("banner not active")
	}

	return domain.ResponseUserBanner{}, err
}

func (s *Service) GetBanner(ctx context.Context, tagId, featureId, limit, offset int) ([]domain.ResponseBanner, error) {
	return s.repo.GetBanner(ctx, tagId, featureId, limit, offset)
}

func (s *Service) PostBanner(ctx context.Context, TagIds []int, featureId int, content domain.ResponseUserBanner, isActive bool) (int, error) {
	return s.repo.PostBanner(ctx, TagIds, featureId, content, isActive)
}

func (s *Service) PatchBanner(ctx context.Context, bannerId int, TagIds []int, featureId int, content domain.ResponseUserBanner, isActive bool) error {
	return s.repo.PatchBanner(ctx, bannerId, TagIds, featureId, content, isActive)
}

func (s *Service) DeleteBanner(ctx context.Context, bannerId int) error {
	return s.repo.DeleteBanner(ctx, bannerId)
}
