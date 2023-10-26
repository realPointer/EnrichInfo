package services

import (
	"context"

	"github.com/realPointer/EnrichInfo/internal/entity"
	"github.com/realPointer/EnrichInfo/internal/repo"
)

type PersonService struct {
	personRepo repo.Person
}

func NewPersonService(personRepo repo.Person) *PersonService {
	return &PersonService{
		personRepo: personRepo,
	}
}

func (s *PersonService) CreatePerson(ctx context.Context, person *entity.EnrichedPerson) error {
	return s.personRepo.CreatePerson(ctx, person)
}

func (s *PersonService) UpdatePerson(ctx context.Context, id int, updatedPerson *entity.EnrichedPerson) error {
	return s.personRepo.UpdatePerson(ctx, id, updatedPerson)
}

func (s *PersonService) DeletePerson(ctx context.Context, id int) error {
	return s.personRepo.DeletePerson(ctx, id)
}

func (s *PersonService) GetPerson(ctx context.Context, id int) (*entity.EnrichedPerson, error) {
	return s.personRepo.GetPerson(ctx, id)
}

func (s *PersonService) SearchPeople(ctx context.Context, filters map[string]string, page, perPage uint64) ([]*entity.EnrichedPerson, error) {
	return s.personRepo.SearchPeople(ctx, filters, page, perPage)
}
