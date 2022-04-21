package plants

import (
	"github.com/jeandeducla/api-plant/internal/models"
)

type Service struct {
    DB DB
}

func NewPlantsService(plantsDB DB) *Service {
    return &Service{DB: plantsDB}
}

func (s *Service) GetAllEnergyManagers() ([]models.EnergyManager, error) {
    return s.DB.GetAllEnergyManagers()
}

type CreateEnergyManagerInput struct {
    Name    string `json:"name" binding:"required"`
    Surname string `json:"surname" binding:"required"`
}

func (s *Service) CreateEnergyManager(input CreateEnergyManagerInput) error {
    em := models.EnergyManager{
        Name: input.Name,
        Surname: input.Surname,
    }
    return s.DB.CreateEnergyManager(&em)
}

func (s *Service) GetEnergyManager(id uint) (*models.EnergyManager, error) {
    return s.DB.GetEnergyManagerById(id)
}

func (s *Service) DeleteEnergyManager(id uint) error {
    return s.DB.DeleteEnergyManagerById(id)
}

type UpdateEnergyManagerInput struct {
    Name    string `json:"name"    binding:"required"`
    Surname string `json:"surname" binding:"required"`
}

func (s *Service) UpdateEnergyManager(id uint, input UpdateEnergyManagerInput) error {
    em, err := s.DB.GetEnergyManagerById(id)
    if err != nil {
        return err
    }
    em.Name = input.Name
    em.Surname = input.Surname
    return s.DB.UpdateEnergyManager(em)
}

func (s *Service) GetEnergyManagerPlants(id uint) ([]models.Plant, error) {
    if _, err := s.DB.GetEnergyManagerById(id); err != nil {
        return nil, err
    }
    return s.DB.GetPlantsByEnergyManagerId(id)
}

func (s *Service) GetAllPlants() ([]models.Plant, error) {
    return s.DB.GetAllPlants()
}

type CreatePlantInput struct {
    Name            string `json:"name"              binding:"required"`
    Address         string `json:"address"           binding:"required"`
    MaxPower        uint   `json:"max_power"         binding:"required"`
    EnergyManagerID uint   `json:"energy_manager_id" binding:"required"`
}

func (s *Service) CreatePlant(input CreatePlantInput) error {
    if _, err := s.DB.GetEnergyManagerById(input.EnergyManagerID); err != nil {
        return err
    }
    return s.DB.CreatePlant(&models.Plant{
        Name: input.Name,
        Address: input.Address,
        MaxPower: input.MaxPower,
        EnergyManagerID: input.EnergyManagerID,
    })
}

func (s *Service) GetPlant(id uint) (*models.Plant, error) {
    return s.DB.GetPlantById(id)
}

func (s *Service) DeletePlant(id uint) error {
    return s.DB.DeletePlantById(id)
}

type UpdatePlantInput struct {
    Name            string `json:"name"              binding:"required"`
    Address         string `json:"address"           binding:"required"`
    MaxPower        uint   `json:"max_power"         binding:"required"`
    EnergyManagerID uint   `json:"energy_manager_id" binding:"required"`
}

func (s *Service) UpdatePlant(id uint, input UpdatePlantInput) error {
    plant, err := s.DB.GetPlantById(id)
    if err != nil {
        return err
    }
    plant.Name = input.Name
    plant.Address = input.Address
    // TODO: validate that the sum of asset power still is lower than the new MaxPower, otherwise error
    plant.MaxPower = input.MaxPower
    // TODO: validate that the new EnergyManagerID is valid? (maybe we want to allow to set it to a null value, like detaching a plant from a EM?)
    plant.EnergyManagerID = input.EnergyManagerID
    return s.DB.UpdatePlant(plant)
}
