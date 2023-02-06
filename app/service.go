package app

import (
	"context"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	rnd    *rand.Rand
	urlDAO *UrlDAO
}

func NewService(urlDAO *UrlDAO) *Service {
	return &Service{
		urlDAO: urlDAO,
		rnd:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *Service) Shorten(ctx context.Context, url string, ttlDays int) (*ShortURL, error) {
	shortURL := &ShortURL{
		URL:      url,
		ExpireAt: getExpirationTime(ttlDays),
	}

	for it := 0; it < 10; it++ {
		shortURL.ID = s.generateRandomID()
		err := s.urlDAO.Insert(ctx, shortURL)
		if err == nil {
			return shortURL, nil
		}
		if !mongo.IsDuplicateKeyError(err) {
			return nil, err
		}
	}
	return nil, ErrCollision
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

var symbols = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func (s *Service) generateRandomID() string {
	const idLength = 6
	id := make([]rune, idLength)
	for i := range id {
		id[i] = symbols[s.rnd.Intn(len(symbols))]
	}
	return string(id)
}

func getExpirationTime(ttlDays int) *time.Time {
	if ttlDays <= 0 {
		return nil
	}
	t := time.Now().Add(time.Hour * 24 * time.Duration(ttlDays))
	return &t
}
