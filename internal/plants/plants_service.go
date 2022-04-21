package plants

import (
	"errors"

	"github.com/jeandeducla/api-plant/internal/models"
)

var (
    ErrAssetPower = errors.New("Asset MaxPower is too big for the plant")
    ErrAssetType = errors.New("Asset Type must be one of 'furnace', 'compressor', 'chiller' or 'rolling mill'")
    ErrNewEmDoesNotExist = errors.New("The EM you want to change to does not exist")
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
    plants, err := s.DB.GetPlantsByEnergyManagerId(id)
    if err != nil && err != ErrEmptyResult {
        return nil, err
    }
    return plants, nil
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

    // checking new max power is ok with existing assets
    existing_assets, err := s.DB.GetAssetsByPlantId(id)
    if err != nil && err != ErrEmptyResult {
        return err
    }
    if sumAssetPower(existing_assets) > input.MaxPower {
        return ErrAssetPower
    }
    plant.MaxPower = input.MaxPower

    // checking em exists
    _, err = s.DB.GetEnergyManagerById(input.EnergyManagerID)
    if err != nil {
        return ErrNewEmDoesNotExist
    }
    plant.EnergyManagerID = input.EnergyManagerID

    plant.Name = input.Name
    plant.Address = input.Address
    return s.DB.UpdatePlant(plant)
}

func (s *Service) GetPlantAssets(id uint) ([]models.Asset, error) {
    if _, err := s.DB.GetPlantById(id); err != nil {
        return nil, err
    }
    assests, err := s.DB.GetAssetsByPlantId(id) 
    if err != nil && err != ErrEmptyResult {
        return nil, err
    }
    return assests, nil
}

type CreateAssetInput struct {
    Name     string `json:"name"      binding:"required"`
    MaxPower uint   `json:"max_power" binding:"required"`
    Type     string `json:"type"      binding:"required"`
}

func (s *Service) CreateAsset(id uint, input CreateAssetInput) error  {
    if input.Type != "furnace" && input.Type != "compressor" && input.Type != "chiller" && input.Type != "rolling mill" {
        return ErrAssetType
    }

    plant, err := s.DB.GetPlantById(id)
    if err != nil {
        return err
    }

    existing_assets, err := s.DB.GetAssetsByPlantId(id)
    if err != nil && err != ErrEmptyResult {
        return err
    }

    if sumAssetPower(existing_assets) + input.MaxPower > plant.MaxPower {
        return ErrAssetPower
    }

    asset := models.Asset{
        Name: input.Name,
        MaxPower: input.MaxPower,
        Type:  input.Type,
        PlantID: id,
    }
    return s.DB.CreateAsset(&asset)
}

func (s *Service) GetPlantAsset(plant_id uint, asset_id uint) (*models.Asset, error) {
    _, err := s.DB.GetPlantById(plant_id)
    if err != nil {
        return nil, err
    }
    return s.DB.GetAssetByPlantId(plant_id, asset_id)
}

func (s *Service) DeletePlantAsset(plant_id uint, asset_id uint) error {
    _, err := s.GetPlantAsset(plant_id, asset_id)
    if err != nil {
        return err
    }
    return s.DB.DeleteAssetById(asset_id)
}

type UpdateAssetInput struct {
    Name     string `json:"name"      binding:"required"`
    MaxPower uint   `json:"max_power" binding:"required"`
    Type     string `json:"type"      binding:"required"`
}

func (s *Service) UpdatePlantAsset(plant_id uint, asset_id uint, input UpdateAssetInput) error {
    if input.Type != "furnace" && input.Type != "compressor" && input.Type != "chiller" && input.Type != "rolling mill" {
        return ErrAssetType
    }

    // checks the asset belongs to the plant
    asset_to_change, err := s.GetPlantAsset(plant_id, asset_id)
    if err != nil {
        return err
    }

    // business rule enforcement
    plant, err := s.DB.GetPlantById(plant_id)
    if err != nil {
        return err
    }
    existing_assets, err := s.DB.GetAssetsByPlantId(plant_id)
    if err != nil && err != ErrEmptyResult {
        return err
    }
    if sumAssetPower(existing_assets) - asset_to_change.MaxPower + input.MaxPower > plant.MaxPower {
        return ErrAssetPower
    }

    asset_to_change.Name = input.Name
    asset_to_change.MaxPower = input.MaxPower
    asset_to_change.Type = input.Type
    return s.DB.UpdateAsset(asset_to_change)
}
