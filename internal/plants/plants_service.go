package plants

import (
	"errors"

	"github.com/jeandeducla/api-plant/internal/models"
)

var (
    ErrAssetPower = errors.New("Asset MaxPower is too big for the plant")
    ErrAssetType = errors.New("Asset Type must be one of 'furnace', 'compressor', 'chiller' or 'rolling mill'")
)

func sumAssetPower(assets []models.Asset) uint {
    var sum uint
    for _, asset := range assets {
        sum += asset.MaxPower
    }
    return sum
}

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

func (s *Service) GetPlantAssets(id uint) ([]models.Asset, error) {
    if _, err := s.DB.GetPlantById(id); err != nil {
        return nil, err
    }
    return s.DB.GetAssetsByPlantId(id)
}

type CreateAssetInput struct {
    Name     string `json:"name"      binding:"required"`
    MaxPower uint   `json:"max_power" binding:"required"`
    Type     string `json:"type"      binding:"required"`
}

func (s *Service) CreateAsset(id uint, input CreateAssetInput) error  {
    plant, err := s.DB.GetPlantById(id)
    if err != nil {
        return err
    }

    existing_assets, err := s.DB.GetAssetsByPlantId(id)
    if err != nil {
        return err
    }

    if sumAssetPower(existing_assets) + input.MaxPower > plant.MaxPower {
        return ErrAssetPower
    }

    if input.Type != "furnace" && input.Type != "compressor" && input.Type != "chiller" && input.Type != "rolling mill" {
        return ErrAssetType
    }

    asset := models.Asset{
        Name: input.Name,
        MaxPower: input.MaxPower,
        Type:  input.Type,
        PlantID: id,
    }
    return s.DB.CreateAsset(&asset)
}

