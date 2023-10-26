package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/realPointer/EnrichInfo/internal/entity"
	"github.com/realPointer/EnrichInfo/internal/service"
	"github.com/realPointer/EnrichInfo/pkg/logger"
)

type peopleRoutes struct {
	peopleService service.Person
	l             logger.Interface
}

func NewPeopleRouter(peopleService service.Person, l logger.Interface) http.Handler {
	p := peopleRoutes{
		peopleService: peopleService,
		l:             l,
	}
	r := chi.NewRouter()

	r.Get("/", p.searchPeople)
	r.Post("/", p.createPerson)
	r.Put("/{id}", p.updatePerson)
	r.Delete("/{id}", p.deletePerson)

	return r
}

// @Summary Create person
// @Description Get name, surname, patronymic, enrich with age, gender and nationality, and save to database
// @Tags People
// @Accept json
// @Success 201
// @Failure 400
// @Router /people [post]
func (p *peopleRoutes) createPerson(w http.ResponseWriter, r *http.Request) {
	// Bind the request body to a PersonInput struct
	person := &entity.PersonInput{}
	if err := render.Bind(r, person); err != nil {
		p.l.Debug("Error binding request body: %v", err)
		render.Render(w, r, ErrorInvalidRequest(err))
		return
	}

	// Enrich the person
	enrichedPerson, err := enrichPerson(person, p.l)
	if err != nil {
		p.l.Debug("Error enriching person: %v", err)
		render.Render(w, r, ErrorInvalidRequest(err))
		return
	}

	// Create the person using the peopleService
	err = p.peopleService.CreatePerson(r.Context(), enrichedPerson)
	if err != nil {
		p.l.Debug("Error creating person: %v", err)
		render.Render(w, r, ErrorInvalidRequest(err))
		return
	}

	// Log the success and set the response status to 201 Created
	p.l.Info("Person created successfully: %v", enrichedPerson)
	render.Status(r, http.StatusCreated)
}

// @Summary Update person
// @Description Update person by id
// @Tags People
// @Accept json
// @Param id path int true "Person ID"
// @Success 200
// @Failure 400
// @Router /people/{id} [put]
func (p *peopleRoutes) updatePerson(w http.ResponseWriter, r *http.Request) {
	personId, err := getIdFromRequest(r)
	if err != nil {
		p.l.Debug("Error getting person ID from request: %v", err)
		render.Render(w, r, ErrorNotFound(err))
		return
	}

	// Check and get info if person with the given ID exists.
	previousPerson, err := p.peopleService.GetPerson(r.Context(), personId)
	if err != nil {
		p.l.Debug("Error getting person with ID %d: %v", personId, err)
		render.Render(w, r, ErrorInvalidRequest(err))
		return
	}

	// Bind the request body to a person struct.
	person := &entity.EnrichedPerson{}
	if err := render.Bind(r, person); err != nil {
		p.l.Debug("Error binding request body: %v", err)
		render.Render(w, r, ErrorInvalidRequest(err))
		return
	}

	// Update the enriched person with the entered data only.
	// Как-то упростить эти конструкции. Какой-нибудь reflect?
	if person.Name != "" {
		previousPerson.Name = person.Name
	}
	if person.Surname != "" {
		previousPerson.Surname = person.Surname
	}
	if person.Patronymic != "" {
		previousPerson.Patronymic = person.Patronymic
	}
	if person.Age != 0 {
		previousPerson.Age = person.Age
	}
	if person.Gender != "" {
		previousPerson.Gender = person.Gender
	}
	if person.Nationality != "" {
		previousPerson.Nationality = person.Nationality
	}

	// If the person's name is provided, re-enrich the person's information.
	if person.Name != "" {
		p.l.Debug("Re-enriching person with name: %s", person.Name)

		// Enrich the person's information with the provided name.
		reEnrichedPerson, err := enrichPerson(&entity.PersonInput{Name: previousPerson.Name, Surname: previousPerson.Surname, Patronymic: previousPerson.Patronymic}, p.l)
		if err != nil {
			p.l.Debug("Error re-enriching person: %v", err)
			render.Render(w, r, ErrorInvalidRequest(err))
			return
		}

		// Update the person's information with the re-enriched data.
		err = p.peopleService.UpdatePerson(r.Context(), personId, reEnrichedPerson)
		if err != nil {
			p.l.Debug("Error updating person: %v", err)
			render.Render(w, r, ErrorInvalidRequest(err))
			return
		}
	} else {
		p.l.Debug("Updating person without re-enriching")

		// Update the person's information without re-enriching.
		err = p.peopleService.UpdatePerson(r.Context(), personId, previousPerson)
		if err != nil {
			p.l.Debug("Error updating person: %v", err)
			render.Render(w, r, ErrorInvalidRequest(err))
			return
		}
	}

	// Return a 200 status code.
	render.Status(r, http.StatusOK)
}

// enrichPerson enriches a person with additional information such as age, gender, and nationality.
func enrichPerson(person *entity.PersonInput, l logger.Interface) (*entity.EnrichedPerson, error) {
	// Create a new EnrichedPerson struct with the provided name, surname, and patronymic.
	enrichedPerson := &entity.EnrichedPerson{
		Name:       person.Name,
		Surname:    person.Surname,
		Patronymic: person.Patronymic,
	}

	// Get the age of the person using the agify.io API.
	ageURL := fmt.Sprintf("https://api.agify.io/?name=%s", person.Name)
	l.Debug("ageURL: ", ageURL)
	ageResp, err := http.Get(ageURL)
	if err != nil {
		l.Error("error while getting age: ", err)
		return nil, err
	}
	defer ageResp.Body.Close()

	var ageData struct {
		Age int `json:"age"`
	}
	if err := json.NewDecoder(ageResp.Body).Decode(&ageData); err != nil {
		l.Error("error while decoding age response: ", err)
		return nil, err
	}

	// Add the age to the enrichedPerson struct.
	enrichedPerson.Age = ageData.Age
	l.Debug("age: ", enrichedPerson.Age)

	// Get the gender of the person using the genderize.io API.
	genderURL := fmt.Sprintf("https://api.genderize.io/?name=%s", person.Name)
	l.Debug("genderURL: ", genderURL)
	genderResp, err := http.Get(genderURL)
	if err != nil {
		l.Error("error while getting gender: ", err)
		return nil, err
	}
	defer genderResp.Body.Close()

	var genderData struct {
		Gender string `json:"gender"`
	}
	if err := json.NewDecoder(genderResp.Body).Decode(&genderData); err != nil {
		l.Error("error while decoding gender response: ", err)
		return nil, err
	}

	// Add the gender to the enrichedPerson struct.
	enrichedPerson.Gender = genderData.Gender
	l.Debug("gender: ", enrichedPerson.Gender)

	// Get the nationality of the person using the nationalize.io API.
	nationalityURL := fmt.Sprintf("https://api.nationalize.io/?name=%s", person.Name)
	l.Debug("nationalityURL: ", nationalityURL)
	nationalityResp, err := http.Get(nationalityURL)
	if err != nil {
		l.Error("error while getting nationality: ", err)
		return nil, err
	}
	defer nationalityResp.Body.Close()

	var nationalityData struct {
		Country []struct {
			Code string `json:"country_id"`
		}
	}
	err = json.NewDecoder(nationalityResp.Body).Decode(&nationalityData)
	if err != nil {
		l.Error("error while decoding nationality response: ", err)
		return nil, err
	}

	// Add the nationality to the enrichedPerson struct.
	if len(nationalityData.Country) != 0 {
		personNationality := nationalityData.Country[0].Code
		enrichedPerson.Nationality = personNationality
	} else {
		l.Error("no nationality found for ", enrichedPerson.Name)
		return nil, fmt.Errorf("no nationality found for %s", enrichedPerson.Name)
	}
	l.Debug("nationality: ", enrichedPerson.Nationality)
	l.Debug("enrichedPerson: ", enrichedPerson)

	// Return the enrichedPerson struct and nil error.
	return enrichedPerson, nil
}

// @Summary Delete person
// @Description Delete person by id
// @Tags People
// @Param id path int true "Person ID"
// @Success 200
// @Failure 400
// @Router /people/{id} [delete]
func (p *peopleRoutes) deletePerson(w http.ResponseWriter, r *http.Request) {
	personId, err := getIdFromRequest(r)
	if err != nil {
		p.l.Debug("Error getting person ID from request: %v", err)
		render.Render(w, r, ErrorNotFound(err))
		return
	}

	// Log the deletion of the person with the given ID.
	p.l.Info("Deleting person with ID %d", personId)

	// Delete the person with the given ID from the database.
	err = p.peopleService.DeletePerson(r.Context(), personId)
	if err != nil {
		p.l.Debug("Error deleting person with ID %d: %v", personId, err)
		render.Render(w, r, ErrorInvalidRequest(err))
		return
	}

	// Log the successful deletion of the person with the given ID.
	p.l.Info("Person with ID %d deleted successfully", personId)

	// Set the response status to 200 OK.
	render.Status(r, http.StatusOK)
}

func getIdFromRequest(r *http.Request) (int, error) {
	personIdStr := chi.URLParam(r, "id")

	personId, err := strconv.Atoi(personIdStr)
	if err != nil {
		return 0, err
	}

	return personId, nil
}

// @Summary Search people
// @Description Returns a list of people matching the specified search criteria
// @Tags People
// @Tags People
// @Param name query string false "Name"
// @Param surname query string false "Surname"
// @Param patronymic query string false "Patronymic"
// @Param age query int false "Age"
// @Param gender query string false "Gender"
// @Param nationality query string false "Nationality"
// @Param page query int false "Page"
// @Param perPage query int false "Persons per page"
// @Success 200
// @Failure 400
// @Router /people [get]
func (p *peopleRoutes) searchPeople(w http.ResponseWriter, r *http.Request) {
	// Get filters from query parameters
	filters := map[string]string{}
	queryParams := []string{"name", "surname", "patronymic", "age", "gender", "nationality"}
	for _, param := range queryParams {
		filters[param] = r.URL.Query().Get(param)
	}

	// Get page number from query parameters
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	// Get number of results per page from query parameters
	perPage, err := strconv.Atoi(r.URL.Query().Get("perPage"))
	if err != nil {
		perPage = 10
	}

	// Convert page and perPage to uint64
	pageUint := uint64(page)
	perPageUint := uint64(perPage)

	// Log search parameters
	p.l.Debug(fmt.Sprintf("searchPeople: filters=%v, page=%d, perPage=%d", filters, pageUint, perPageUint))

	// Search people based on filters, page number, and number of results per page
	people, err := p.peopleService.SearchPeople(r.Context(), filters, pageUint, perPageUint)
	if err != nil {
		// Return error response if there is an error while searching
		p.l.Error(fmt.Sprintf("searchPeople: error=%v", err))
		render.Render(w, r, ErrorInvalidRequest(err))
		return
	}

	// Log search results
	p.l.Debug(fmt.Sprintf("searchPeople: people=%v", people))

	// Return JSON response with search results
	render.JSON(w, r, people)
}
