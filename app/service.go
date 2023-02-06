package app

import (
	"context"
	"math/rand"
	"time"
)

type Service struct {
	rnd     *rand.Rand
	urlsMap map[string]*ShortURL
}

func NewService() *Service {
	return &Service{
		urlsMap: make(map[string]*ShortURL),
		rnd:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *Service) Shorten(ctx context.Context, url string, ttlDays int) (*ShortURL, error) {
	shortURL := &ShortURL{
		ID:       s.generateRandomID(),
		URL:      url,
		ExpireAt: getExpirationTime(ttlDays),
	}

	for it := 0; it < 10; it++ {
		shortURL.ID = s.generateRandomID()
		if _, ok := s.urlsMap[shortURL.ID]; !ok {
			s.urlsMap[shortURL.ID] = shortURL
			return shortURL, nil
		}
	}
	return nil, ErrCollision
}

func (s *Service) Update(ctx context.Context, id string, url string, ttlDays int) (*ShortURL, error) {
	sURL, ok := s.urlsMap[id]
	if !ok {
		return nil, ErrNotFound
	}

	sURL.URL = url
	sURL.ExpireAt = getExpirationTime(ttlDays)

	s.urlsMap[id] = sURL
	return sURL, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	delete(s.urlsMap, id)
	return nil
}

func (s *Service) GetFullURL(ctx context.Context, shortURL string) (string, error) {
	sURL, ok := s.urlsMap[shortURL]
	if !ok {
		return "", ErrNotFound
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

func getExpirationTime(ttlDays int) time.Time {
	if ttlDays <= 0 {
		return time.Time{}
	}
	return time.Now().Add(time.Hour * 24 * time.Duration(ttlDays))
}
