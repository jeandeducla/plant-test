package server

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/mock"

	"github.com/jeandeducla/api-plant/internal/models"
	"github.com/jeandeducla/api-plant/internal/plants"
)

type MainTestSuite struct {
    suite.Suite
    service *plants.Service
    server *Server
}

func TestServer(t *testing.T) {
    suite.Run(t, new(MainTestSuite))
}

var mockDb = &plants.DbMock{}

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

func (t *MainTestSuite) SetupSuite() {
    service := plants.NewPlantsService(mockDb)
    t.service = service

    server, err := NewServer(service)
    t.Require().NoError(err)
    t.server = server
}

func (t *MainTestSuite) TestGetAllEnergyManagers() {
    mockDb.On("GetAllEnergyManagers").Return([]models.EnergyManager{em1, em2}, nil)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/ems", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
}

func (t *MainTestSuite) TestPostEnergyManager() {
    mockDb.On("CreateEnergyManager", mock.Anything).Return(nil)

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
}

func (t *MainTestSuite) TestGetEnergyManager() {
    mockDb.On("GetEnergyManagerById", uint(1)).Return(&em1, nil)
    mockDb.On("GetEnergyManagerById", uint(2)).Return(&em2, nil)
    mockDb.On("GetEnergyManagerById", uint(1234)).Return(&models.EnergyManager{}, plants.ErrEmptyResult)
    mockDb.On("GetEnergyManagerById", uint(42)).Return(&models.EnergyManager{}, errors.New("DB error"))

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

    // test 500
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/ems/42", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(500, w.Code)
}

func (t *MainTestSuite) TestDeleteEnergyManager() {
    mockDb.On("DeleteEnergyManagerById", uint(1)).Return(nil)
    mockDb.On("DeleteEnergyManagerById", uint(4)).Return(plants.ErrEmptyResult)

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

    // id does exist
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("DELETE", "/ems/1", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
}

func (t *MainTestSuite) TestPutEnergyManager() {
    mockDb.On("GetEnergyManagerById", uint(1)).Return(&em1, nil)
    mockDb.On("UpdateEnergyManager", mock.Anything).Return(nil)
    mockDb.On("GetEnergyManagerById", uint(4)).Return(&models.EnergyManager{}, plants.ErrEmptyResult)

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
}
