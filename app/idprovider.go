package app

import (
	"context"
	"log"
	"math/rand"

	"go.mongodb.org/mongo-driver/mongo"
)

type IDProvider struct {
	idDAO  *IDDAO
	idChan chan string
}

func NewIDProvider(ctx context.Context, idDAO *IDDAO) *IDProvider {
	provider := &IDProvider{
		idChan: make(chan string, 10),
		idDAO:  idDAO,
	}
	go provider.generate(ctx)
	return provider
}

var symbols = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func (p *IDProvider) generate(ctx context.Context) {
	p.preloadIds(ctx)
	for {
		idCandidate := generateRandomID()
		err := p.idDAO.Insert(ctx, &UrlID{
			ID: idCandidate,
		})
		if err != nil && !mongo.IsDuplicateKeyError(err) {
			log.Println(err)
		}
		p.idChan <- idCandidate
	}
}

func (p *IDProvider) preloadIds(ctx context.Context) {
	ids, err := p.idDAO.FindUnused(ctx)
	if err != nil {
		log.Println(err)
	}
	for _, id := range ids {
		p.idChan <- id.ID
	}
}

func (p *IDProvider) GetID(ctx context.Context) (string, error) {
	for {
		id := <-p.idChan
		err := p.idDAO.ReserveID(ctx, id)
		if err == nil {
			return id, nil
		}
		if err != ErrIDNotFound {
			return "", err
		}
	}
}

func generateRandomID() string {
	const idLength = 6
	id := make([]rune, idLength)
	for i := range id {
		id[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(id)
}
