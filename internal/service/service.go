package service

import (
	"context"

	"github.com/realPointer/EnrichInfo/internal/entity"
	"github.com/realPointer/EnrichInfo/internal/repo"
	"github.com/realPointer/EnrichInfo/internal/service/services"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Person interface {
	CreatePerson(ctx context.Context, person *entity.EnrichedPerson) error
	UpdatePerson(ctx context.Context, id int, updatedPerson *entity.EnrichedPerson) error
	DeletePerson(ctx context.Context, id int) error
	GetPerson(ctx context.Context, id int) (*entity.EnrichedPerson, error)
	SearchPeople(ctx context.Context, filters map[string]string, page, perPage uint64) ([]*entity.EnrichedPerson, error)
}

type Services struct {
	Person
}

type ServicesDependencies struct {
	Repos *repo.Repositories
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Person: services.NewPersonService(deps.Repos.Person),
	}
}
