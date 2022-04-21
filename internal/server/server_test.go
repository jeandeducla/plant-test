package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/jeandeducla/api-plant/internal/models"
	"github.com/jeandeducla/api-plant/internal/plants"
)

type MainTestSuite struct {
    suite.Suite
    service *plants.Service
    server *Server
    db *gorm.DB
}

func TestServer(t *testing.T) {
    suite.Run(t, new(MainTestSuite))
}

func (t *MainTestSuite) SetupTest() {
    db, err := models.NewDB("postgres://postgres:postgres@postgresql/metron")
    t.Require().NoError(err)

    t.db = db

    service := plants.NewPlantsService(plants.NewPlantsDB(db))
    t.service = service

    server, err := NewServer(service)
    t.Require().NoError(err)
    t.server = server
}

func (t *MainTestSuite) TearDownTest() {
    t.db.Migrator().DropTable(&models.Asset{})
    t.db.Migrator().DropTable(&models.Plant{})
    t.db.Migrator().DropTable(&models.EnergyManager{})
}

func (t *MainTestSuite) TestGetAllEnergyManagers() {
    // No ems should not return an error and res should be an empty slice
    {
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("GET", "/ems", nil)
        t.server.Router().ServeHTTP(w, req)
        t.Equal(200, w.Code)
        var res []models.EnergyManager
        t.Require().NoError(json.NewDecoder(w.Result().Body).Decode(&res))
        t.Equal(len(res), 0)
    }

    // There are some ems
    {
        body := []byte(`
             {
                 "name": "coucou",
                 "surname": "salut"
             }
         `)
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("POST", "/ems", bytes.NewReader(body))
        t.server.Router().ServeHTTP(w, req)
        t.Equal(200, w.Code)
        req, _ = http.NewRequest("GET", "/ems", nil)
        t.server.Router().ServeHTTP(w, req)
        var res []models.EnergyManager
        t.Require().NoError(json.NewDecoder(w.Result().Body).Decode(&res))
        t.Equal(len(res), 1)
        t.Equal(res[0].Name, "coucou")
    }
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
            "name": "Jacques",
            "surname": "salut",
            "useless": 1
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/ems", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)

    {
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("GET", "/ems", nil)
        t.server.Router().ServeHTTP(w, req)
        t.Equal(200, w.Code)
        var res []models.EnergyManager
        t.Require().NoError(json.NewDecoder(w.Result().Body).Decode(&res))
        t.Equal(len(res), 2)
    }

    {
        w := httptest.NewRecorder()
        req, _ = http.NewRequest("GET", "/ems/2", nil)
        t.server.Router().ServeHTTP(w, req)
        t.Equal(200, w.Code)
        var res models.EnergyManager
        t.Require().NoError(json.NewDecoder(w.Result().Body).Decode(&res))
        t.Equal(res.Name, "Jacques")
    }
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
    body := []byte(`
        {
            "name": "coucou",
            "surname": "salut"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/ems", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/ems/1", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)

    // bad http verb
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("OPTIONS", "/ems/2", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)
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

    // id does exist
    body := []byte(`
        {
            "name": "coucou",
            "surname": "salut"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/ems", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("DELETE", "/ems/1", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/ems/1", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)
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
            "surname": "salut"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/ems", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    w = httptest.NewRecorder()
    body = []byte(`
        {
            "name": "Jacques",
            "surname": "salut",
            "roger": "moore"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("PUT", "/ems/1", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/ems/1", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    var res models.EnergyManager
    t.Require().NoError(json.NewDecoder(w.Result().Body).Decode(&res))
    t.Equal(res.Name, "Jacques")
}

func (t *MainTestSuite) TestGetEnergyManagerPlants() {
    // id is not parsable to uint
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("PUT", "/ems/^jlliu/plants", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)

    // id does not exist
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("PUT", "/ems/1/plants", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)

    // no plants should return empty slice of plants
    body := []byte(`
        {
            "name": "Gerard",
            "surname": "Depardieu"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/ems", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/ems/1/plants", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    var res []models.Plant
    t.Require().NoError(json.NewDecoder(w.Result().Body).Decode(&res))
    t.Equal(len(res), 0)

    // Some plants 
    body = []byte(`
        {
            "name": "Gerard",
            "address": "187 rue triuy",
            "max_power": 189,
            "energy_manager_id": 1
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/plants", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    req, _ = http.NewRequest("GET", "/ems/1/plants", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    {
        var res []models.Plant
        t.Require().NoError(json.NewDecoder(w.Result().Body).Decode(&res))
        t.Equal(len(res), 1)
        t.Equal(res[0].Name, "Gerard")
    }
}

func (t *MainTestSuite) TestGetPlants() {
    // no plants should return empty slice
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/plants", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    var res []models.Plant
    t.Require().NoError(json.NewDecoder(w.Result().Body).Decode(&res))
    t.Equal(len(res), 0)

    // Some plant
    body := []byte(`
        {
            "name": "Gerard",
            "surname": "Depardieu"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/ems", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    body = []byte(`
        {
            "name": "Gerard",
            "address": "187 rue triuy",
            "max_power": 189,
            "energy_manager_id": 1
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/plants", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)


    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/plants", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    {
        var res []models.Plant
        t.Require().NoError(json.NewDecoder(w.Result().Body).Decode(&res))
        t.Equal(len(res), 1)
        t.Equal(res[0].Name, "Gerard")
    }
}

func (t *MainTestSuite) TestPostPlants() {
    // empty body
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/plants", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(400, w.Code)

    // bad body
    body := []byte(`
        {
            "name": "Gerard",
            "surname": "Depardieu"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/plants", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(400, w.Code)

    // bad body
    body = []byte(`
        {
            "name": "Gerard",
            "Address": 1,
            "max)Power": 
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/plants", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(400, w.Code)

    // can't create plant if em does not exist
    body = []byte(`
        {
            "name": "Gerard",
            "address": "187 rue triuy",
            "max_power": 189,
            "energy_manager_id": 1
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/plants", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(400, w.Code)

    // works
    body = []byte(`
        {
            "name": "Gerard",
            "surname": "Bond"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/ems", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    body = []byte(`
        {
            "name": "Plant",
            "address": "187 rue triuy",
            "max_power": 189,
            "energy_manager_id": 1
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/plants", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/plants/1", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    var res models.Plant
    t.Require().NoError(json.NewDecoder(w.Result().Body).Decode(&res))
    t.Equal(res.Name, "Plant")
}

func (t *MainTestSuite) TestDeletePlant() {
    // id is not parsable
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("DELETE", "/plants/joel", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)

    // id does not exist
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("DELETE", "/plants/1", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)

    // works
    body := []byte(`
        {
            "name": "Gerard",
            "surname": "Depardieu"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/ems", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    body = []byte(`
        {
            "name": "Plant",
            "address": "187 rue triuy",
            "max_power": 189,
            "energy_manager_id": 1
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/plants", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("DELETE", "/plants/1", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/plants/1", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/ems/1/plants", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    var res []models.Plant
    t.Require().NoError(json.NewDecoder(w.Result().Body).Decode(&res))
    t.Equal(len(res), 0)
}

func (t *MainTestSuite) TestPutPlant() {
    // id is not parsable
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("PUT", "/plants/joel", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)

    // id does not exist
    body := []byte(`
        {
            "name": "Plant",
            "address": "187 rue triuy",
            "max_power": 189,
            "energy_manager_id": 1
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("PUT", "/plants/1", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)

    // invalid body
    body = []byte(`
        {
            "name": "Plant",
            "address": "187 rue triuy",
            "max_Power": 189
            "energy_manager_id": 1
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("PUT", "/plants/1", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(400, w.Code)

    // update to em that does not exist
    body = []byte(`
        {
            "name": "Eric",
            "surname": "judor"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/ems", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    body = []byte(`
        {
            "name": "Plant",
            "address": "187 rue triuy",
            "max_power": 189,
            "energy_manager_id": 1
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/plants", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    body = []byte(`
        {
            "name": "Plant",
            "address": "187 rue triuy",
            "max_power": 189,
            "energy_manager_id": 1123
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("PUT", "/plants/1", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(400, w.Code)
}

func (t *MainTestSuite) TestGetPlantAssets() {
    // id is not parsable
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/plants/joel/assets", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)

    // id does not exist
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/plants/1/assets", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)

    // No assets = no errors and empty list
    body := []byte(`
        {
            "name": "Eric",
            "surname": "judor"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/ems", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    body = []byte(`
        {
            "name": "Plant",
            "address": "187 rue triuy",
            "max_power": 189,
            "energy_manager_id": 1
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/plants", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/plants/1/assets", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    {
        var res []models.Asset
        t.Require().NoError(json.NewDecoder(w.Result().Body).Decode(&res))
        t.Equal(len(res), 0)
    }

    // Some assets
    body = []byte(`
        {
            "name": "asset",
            "max_power": 189,
            "type": "furnace",
            "plant_id": 1
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/plants/1/assets", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/plants/1/assets", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    {
        var res []models.Asset
        t.Require().NoError(json.NewDecoder(w.Result().Body).Decode(&res))
        t.Equal(len(res), 1)
        t.Equal(res[0].Name, "asset")
    }
}

func (t *MainTestSuite) TestPostPlantAssets() {
    // id is not parsable
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/plants/joel/assets", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)

    // id does not exist
    body := []byte(`
        {
            "name": "asset",
            "max_power": 10,
            "type": "furnace"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/plants/1/assets", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)

    // bad body
    body = []byte(`
        {
            "name": "Eric",
            "surname": "judor"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/ems", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    body = []byte(`
        {
            "name": "Plant",
            "address": "187 rue triuy",
            "max_power": 189,
            "energy_manager_id": 1
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/plants", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    body = []byte(`
        {
            "name": "asset",
            "max_power": 10,
            "coucou": "furnace"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/plants/1/assets", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(400, w.Code)

    // wrong type
    body = []byte(`
        {
            "name": "asset",
            "max_power": 10,
            "type": "JACK CHIRAK"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/plants/1/assets", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(400, w.Code)

    // plant id does not exist
    body = []byte(`
        {
            "name": "asset",
            "max_power": 10,
            "type": "chiller"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/plants/1123/assets", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)
    
    // power limit is broken
    body = []byte(`
        {
            "name": "asset",
            "max_power": 1000000,
            "type": "chiller"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/plants/1/assets", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(400, w.Code)
}

func (t *MainTestSuite) TestDeletePlantAsset() {
    // id is not parsable
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("DELETE", "/plants/joel/assets", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)

    // id does not exist
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("DELETE", "/plants/1/assets/1", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)
    // id does not exist 2
    body := []byte(`
        {
            "name": "Eric",
            "surname": "judor"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/ems", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    body = []byte(`
        {
            "name": "Plant",
            "address": "187 rue triuy",
            "max_power": 189,
            "energy_manager_id": 1
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/plants", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("DELETE", "/plants/1/assets/1", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)

    // id does exist
    body = []byte(`
        {
            "name": "asset",
            "max_power": 10,
            "type": "furnace"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/plants/1/assets", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("DELETE", "/plants/1/assets/1", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/plants/1/assets/1", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)
}

func (t *MainTestSuite) TestPutPlantAsset() {
    // id is not parsable
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("PUT", "/plants/joel/assets", nil)
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)

    // id does not exist
    body := []byte(`
        {
            "name": "asset",
            "max_power": 10,
            "type": "furnace"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("PUT", "/plants/1/assets/1", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)
    // id does not exist 2
    body = []byte(`
        {
            "name": "Eric",
            "surname": "judor"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/ems", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    body = []byte(`
        {
            "name": "Plant",
            "address": "187 rue triuy",
            "max_power": 189,
            "energy_manager_id": 1
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/plants", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    body = []byte(`
        {
            "name": "asset",
            "max_power": 10,
            "type": "furnace"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("PUT", "/plants/1/assets/1", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(404, w.Code)

    // bad body
    body = []byte(`
        {
            "name": "asset",
            "max_power": 10,
            "type": "furnace"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/plants/1/assets", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(200, w.Code)
    body = []byte(`
        {
            "name": "asset",
            "max_power": 10,
            "coucou": "furnace"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("PUT", "/plants/1/assets/1", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(400, w.Code)

    // update to bad type
    body = []byte(`
        {
            "name": "asset",
            "max_power": 10,
            "type": "JACK CHIRAK"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("PUT", "/plants/1/assets/1", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(400, w.Code)

    // update breaks power limit
    body = []byte(`
        {
            "name": "asset",
            "max_power": 100000,
            "type": "furnace"
        }
    `)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("PUT", "/plants/1/assets/1", bytes.NewReader(body))
    t.server.Router().ServeHTTP(w, req)
    t.Equal(400, w.Code)
}
