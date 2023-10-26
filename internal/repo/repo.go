package repo

import (
	"context"

	"github.com/realPointer/EnrichInfo/internal/entity"
	"github.com/realPointer/EnrichInfo/internal/repo/postgresdb"
	"github.com/realPointer/EnrichInfo/pkg/postgres"
)

type Person interface {
	CreatePerson(ctx context.Context, person *entity.EnrichedPerson) error
	UpdatePerson(ctx context.Context, id int, updatedPerson *entity.EnrichedPerson) error
	DeletePerson(ctx context.Context, id int) error
	GetPerson(ctx context.Context, id int) (*entity.EnrichedPerson, error)
	SearchPeople(ctx context.Context, filters map[string]string, page, perPage uint64) ([]*entity.EnrichedPerson, error)
}

type Repositories struct {
	Person
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Person: postgresdb.NewPersonRepo(pg),
	}
}
