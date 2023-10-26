package entity

import (
	"errors"
	"net/http"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type PersonInput struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}

func (p *PersonInput) Bind(r *http.Request) error {
	if p.Name == "" {
		return errors.New("missing required name fields")
	}

	if p.Surname == "" {
		return errors.New("missing required surname fields")
	}

	p.Name = strings.TrimSpace(p.Name)
	p.Surname = strings.TrimSpace(p.Surname)
	p.Patronymic = strings.TrimSpace(p.Patronymic)

	p.Name = cases.Title(language.English).String(p.Name)
	p.Surname = cases.Title(language.English).String(p.Surname)
	p.Patronymic = cases.Title(language.English).String(p.Patronymic)

	return nil
}

type EnrichedPerson struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic,omitempty"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

func (p *EnrichedPerson) Bind(r *http.Request) error {
	p.Name = strings.TrimSpace(p.Name)
	p.Surname = strings.TrimSpace(p.Surname)
	p.Patronymic = strings.TrimSpace(p.Patronymic)

	p.Name = cases.Title(language.English).String(p.Name)
	p.Surname = cases.Title(language.English).String(p.Surname)
	p.Patronymic = cases.Title(language.English).String(p.Patronymic)

	return nil
}
