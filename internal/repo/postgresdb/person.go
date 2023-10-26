package postgresdb

import (
	"context"
	"fmt"

	"github.com/realPointer/EnrichInfo/internal/entity"
	"github.com/realPointer/EnrichInfo/pkg/postgres"
)

type PersonRepo struct {
	*postgres.Postgres
}

func NewPersonRepo(pg *postgres.Postgres) *PersonRepo {
	return &PersonRepo{
		Postgres: pg,
	}
}

func (r *PersonRepo) CreatePerson(ctx context.Context, person *entity.EnrichedPerson) error {
	sql, args, _ := r.Builder.
		Insert("people").
		Columns("name", "surname", "patronymic", "age", "gender", "nationality").
		Values(person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality).
		ToSql()

	_, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PersonRepo - CreatePerson - tx.Exec: %v", err)
	}

	return nil
}

func (r *PersonRepo) UpdatePerson(ctx context.Context, id int, updatedPerson *entity.EnrichedPerson) error {
	sql, args, _ := r.Builder.
		Update("people").
		Set("name", updatedPerson.Name).
		Set("surname", updatedPerson.Surname).
		Set("patronymic", updatedPerson.Patronymic).
		Set("age", updatedPerson.Age).
		Set("gender", updatedPerson.Gender).
		Set("nationality", updatedPerson.Nationality).
		Where("id = ?", id).
		ToSql()

	_, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PersonRepo - UpdatePerson - tx.Exec: %v", err)
	}

	return nil
}

func (r *PersonRepo) DeletePerson(ctx context.Context, id int) error {
	sql, args, _ := r.Builder.
		Delete("people").
		Where("id = ?", id).
		ToSql()

	_, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PersonRepo - DeletePerson - tx.Exec: %v", err)
	}

	return nil
}

func (r *PersonRepo) GetPerson(ctx context.Context, id int) (*entity.EnrichedPerson, error) {
	sql, args, _ := r.Builder.
		Select("*").
		From("people").
		Where("id = ?", id).
		ToSql()

	row := r.Pool.QueryRow(ctx, sql, args...)
	person := &entity.EnrichedPerson{}
	err := row.Scan(&person.ID, &person.Name, &person.Surname, &person.Patronymic, &person.Age, &person.Gender, &person.Nationality)
	if err != nil {
		return nil, fmt.Errorf("PersonRepo - GetPerson - row.Scan: %v", err)
	}

	return person, nil
}

func (r *PersonRepo) SearchPeople(ctx context.Context, filters map[string]string, page, perPage uint64) ([]*entity.EnrichedPerson, error) {
	builder := r.Builder.Select("*").From("people")

	// Add filters to the query
	for key, value := range filters {
		if value != "" {
			builder = builder.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}

	builder = builder.Limit(perPage)
	if page > 1 {
		builder = builder.Offset((page - 1) * perPage)
	}

	sql, args, _ := builder.ToSql()
	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("PersonRepo - SearchPeople - p.Pool.Query: %v", err)
	}
	defer rows.Close()

	// Parse the query results into a slice of EnrichedPerson structs
	people := []*entity.EnrichedPerson{}
	for rows.Next() {
		person := &entity.EnrichedPerson{}
		err := rows.Scan(&person.ID, &person.Name, &person.Surname, &person.Patronymic, &person.Age, &person.Gender, &person.Nationality)
		if err != nil {
			return nil, fmt.Errorf("PersonRepo - SearchPeople - rows.Scan: %v", err)
		}
		people = append(people, person)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("PersonRepo - SearchPeople - rows.Err: %v", err)
	}

	return people, nil
}
