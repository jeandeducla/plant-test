package plants


import (
	"github.com/stretchr/testify/mock"

	"github.com/jeandeducla/api-plant/internal/models"
)

type DbMock struct {
    mock.Mock
}

func (dbm *DbMock) GetAllEnergyManagers() ([]models.EnergyManager, error) {
    args := dbm.Called()
    return args.Get(0).([]models.EnergyManager), args.Error(1)
}

func (dbm *DbMock) CreateEnergyManager(em *models.EnergyManager) error {
    args := dbm.Called(em)
    return args.Error(0)
}

func (dbm *DbMock) GetEnergyManagerById(id uint) (*models.EnergyManager, error) {
    args := dbm.Called(id)
    return args.Get(0).(*models.EnergyManager), args.Error(1)
}

func (dbm *DbMock) DeleteEnergyManagerById(id uint) error {
    args := dbm.Called(id)
    return args.Error(0)
}

func (dbm *DbMock) UpdateEnergyManager(em *models.EnergyManager) error {
    args := dbm.Called(em)
    return args.Error(0)
}

func (dbm *DbMock) GetAllPlants() ([]models.Plant, error) {
    args := dbm.Called()
    return args.Get(0).([]models.Plant), args.Error(1)
}

func (dbm *DbMock) CreatePlant(plant *models.Plant) error {
    args := dbm.Called(plant)
    return args.Error(0)
}

func (dbm *DbMock) GetPlantById(id uint) (*models.Plant, error) {
    args := dbm.Called(id)
    return args.Get(0).(*models.Plant), args.Error(1)
}

func (dbm *DbMock) DeletePlantById(id uint) error {
    args := dbm.Called(id)
    return args.Error(0)
}

func (dbm *DbMock) UpdatePlant(plant *models.Plant) error {
    args := dbm.Called(plant)
    return args.Error(0)
}

func (dbm *DbMock) GetPlantsByEnergyManagerId(id uint) ([]models.Plant, error) {
    args := dbm.Called(id)
    return args.Get(0).([]models.Plant), args.Error(1)
}
