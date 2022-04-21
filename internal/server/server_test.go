package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jeandeducla/api-plant/internal/models"
	"github.com/jeandeducla/api-plant/internal/plants"
	"github.com/stretchr/testify/suite"
)

type MainTestSuite struct {
    suite.Suite
    server *Server
}

func TestServer(t *testing.T) {
    suite.Run(t, new(MainTestSuite))
}

type ServiceMocker struct {
    db map[uint]models.EnergyManager
}

var em1 = models.EnergyManager{
    Name: "Gerard",
    Surname: "Depardieu",
    Plants: nil,
}

var em2 = models.EnergyManager{
    Name: "Patrick",
    Surname: "Sebastien",
    Plants: nil,
}

var dbID uint = 3


func (sm *ServiceMocker) GetAllEnergyManagers() ([]models.EnergyManager, error) {
    slice := make([]models.EnergyManager, 0, len(sm.db))
    for _, em := range sm.db {
        slice = append(slice, em)
    }
    return slice, nil
}

func (sm *ServiceMocker) CreateEnergyManager(input plants.CreateEnergyManagerInput) error {
    em := models.EnergyManager{
        Name: input.Name,
        Surname: input.Surname,
    }
    sm.db[dbID] = em
    dbID++
    return nil
}

func (sm *ServiceMocker) GetEnergyManager(id uint) (*models.EnergyManager, error) {
    em, ok := sm.db[id]
    if !ok {
        return nil, plants.ErrEmptyResult
    }
    return &em, nil
}

func (sm *ServiceMocker) DeleteEnergyManager(id uint) error {
    _, ok := sm.db[id]
    if !ok {
        return plants.ErrEmptyResult
    }
    delete(sm.db, id)
    return nil
}

func (sm *ServiceMocker) UpdateEnergyManager(id uint, input plants.UpdateEnergyManagerInput) error {
    em, ok := sm.db[id]
    if !ok {
        return plants.ErrEmptyResult
    }
    em.Name = input.Name
    em.Surname = input.Surname
    sm.db[id] = em
    return nil
}

func (sm *ServiceMocker) GetEnergyManagerPlants(id uint) ([]models.Plant, error) {
    return nil, nil
}


func (sm *ServiceMocker) GetAllPlants() ([]models.Plant, error) {
    return nil, nil
}

func (sm *ServiceMocker) CreatePlant(input plants.CreatePlantInput) error {
    return nil
}

func (sm *ServiceMocker) GetPlant(id uint) (*models.Plant, error) {
    return nil, nil
}

func (sm *ServiceMocker) DeletePlant(id uint) error {
    return nil
}

func (sm *ServiceMocker) UpdatePlant(id uint, input plants.UpdatePlantInput) error {
    return nil
}


func (t *MainTestSuite) SetupTest() {
    // re-initialize mock db
    server, err := NewServer(&ServiceMocker{db: map[uint]models.EnergyManager{
        1: em1,
        2: em2,
    }})
    t.Require().NoError(err)
    t.server = server
}

func (t *MainTestSuite) TestGetAllEnergyManagers() {
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/ems", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
}

func (t *MainTestSuite) TestPostEnergyManager() {
    // empty body
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/ems", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(400, w.Code)

    // body is missing parts
    body := []byte(`
        {
            "name": "coucou",
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/ems", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(400, w.Code)
    // check it has not affected the DB
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/ems", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    res := []models.EnergyManager{}
    _ = json.Unmarshal(w.Body.Bytes(), &res)
    t.Equal(2, len(res))

    // body is correct
    body = []byte(`
        {
            "name": "coucou",
            "surname": "salut"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/ems", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    // check DB has changed
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/ems", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    res = []models.EnergyManager{}
    _ = json.Unmarshal(w.Body.Bytes(), &res)
    t.Equal(3, len(res))

    // body has more fields than expected
    body = []byte(`
        {
            "name": "coucou",
            "surname": "salut",
            "useless": 1
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/ems", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    // check DB has changed
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/ems", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    res = []models.EnergyManager{}
    _ = json.Unmarshal(w.Body.Bytes(), &res)
    t.Equal(4, len(res))

    // TODO: test 500
}

func (t *MainTestSuite) TestGetEnergyManager() {
    // id is not parsable to uint
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/ems/hjki", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)

    // no id
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/ems/", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(301, w.Code)

    // id does not exist
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/ems/1234", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)

    // id does exist
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/ems/1", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)

    // bad http verb
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("OPTIONS", "/ems/2", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)

    // TODO: test 500
}

func (t *MainTestSuite) TestDeleteEnergyManager() {
    // id is not parsable to uint
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("DELETE", "/ems/jemangedupoulet", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)

    // id does not exist
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("DELETE", "/ems/4", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)
    // check it has not affected the DB
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/ems", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    res := []models.EnergyManager{}
    _ = json.Unmarshal(w.Body.Bytes(), &res)
    t.Equal(2, len(res))

    // id does exist
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("DELETE", "/ems/1", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    // check DB has changed
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/ems", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    res = []models.EnergyManager{}
    _ = json.Unmarshal(w.Body.Bytes(), &res)
    t.Equal(1, len(res))
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/ems/1", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)

    // TODO: test 500
}

func (t *MainTestSuite) TestPutEnergyManager() {
    // id is not parsable to uint
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("PUT", "/ems/^jlliu", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)

    // id does not exist and request body is valid
    body := []byte(`
        {
            "name": "coucou",
            "surname": "salut"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("PUT", "/ems/4", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)
    // check it has not affected the DB
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/ems", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    res := []models.EnergyManager{}
    _ = json.Unmarshal(w.Body.Bytes(), &res)
    t.Equal(2, len(res))
    t.Equal(em1.Name, res[0].Name)
    t.Equal(em1.Surname, res[0].Surname)
    t.Equal(em2.Name, res[1].Name)
    t.Equal(em2.Surname, res[1].Surname)

    // request body is empty
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("PUT", "/ems/1", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(400, w.Code)

    // request body is not valid
    body = []byte(`
        {
            "name": "coucou",
            "surname": "salut",
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("PUT", "/ems/1", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(400, w.Code)

    // request body has useless fields
    body = []byte(`
        {
            "name": "coucou",
            "surname": "salut",
            "roger": "moore"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("PUT", "/ems/1", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    // check it has affected the DB
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/ems/1", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    result := models.EnergyManager{}
    _ = json.Unmarshal(w.Body.Bytes(), &result)
    t.Equal("coucou", result.Name)
    t.Equal("salut", result.Surname)

    // TODO: test 500
}
