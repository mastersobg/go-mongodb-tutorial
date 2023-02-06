package app

import (
	"context"
	"time"
)

type Service struct {
	urlDAO     *UrlDAO
	idProvider *IDProvider
}

func NewService(urlDAO *UrlDAO, idProvider *IDProvider) *Service {
	return &Service{
		urlDAO:     urlDAO,
		idProvider: idProvider,
	}
}

func (s *Service) Shorten(ctx context.Context, url string, ttlDays int) (*ShortURL, error) {
	shortID, err := s.idProvider.GetID(ctx)
	if err != nil {
		return nil, err
	}
	shortURL := &ShortURL{
		ID:       shortID,
		URL:      url,
		ExpireAt: getExpirationTime(ttlDays),
	}
	err = s.urlDAO.Insert(ctx, shortURL)
	if err != nil {
		return nil, err
	}
	return shortURL, nil
}

func (s *Service) Update(ctx context.Context, id string, url string, ttlDays int) (*ShortURL, error) {
	sURL, err := s.urlDAO.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	sURL.URL = url
	sURL.ExpireAt = getExpirationTime(ttlDays)

	return sURL, s.urlDAO.Update(ctx, sURL)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.urlDAO.DeleteByID(ctx, id)
}

func (s *Service) GetFullURL(ctx context.Context, shortURL string) (string, error) {
	sURL, err := s.urlDAO.FindByID(ctx, shortURL)
	if err != nil {
		return "", err
	}
	return sURL.URL, nil
}

func getExpirationTime(ttlDays int) time.Time {
	if ttlDays <= 0 {
		return time.Time{}
	}
	return time.Now().Add(time.Hour * 24 * time.Duration(ttlDays))
}
