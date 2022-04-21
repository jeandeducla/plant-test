package plants


import (
	"testing"

	"github.com/stretchr/testify/suite"
    "gorm.io/gorm"

	"github.com/jeandeducla/api-plant/internal/models"
)

type MainTestSuite struct {
    suite.Suite
    service *Service
    db *gorm.DB
}

func TestMainTestSuite(t *testing.T) {
    suite.Run(t, new(MainTestSuite))
}

func (t *MainTestSuite) SetupTest() {
    db, err := models.NewDB("postgres://postgres:postgres@postgresql/metron")
    t.db = db
    t.Require().NoError(err)

    service := NewPlantsService(NewPlantsDB(db))
    t.service = service
}

func (t *MainTestSuite) TearDownTest() {
    t.db.Migrator().DropTable(&models.Asset{})
    t.db.Migrator().DropTable(&models.Plant{})
    t.db.Migrator().DropTable(&models.EnergyManager{})
}

func (t *MainTestSuite) TestCreateEnergyManager() {
    err := t.service.CreateEnergyManager(CreateEnergyManagerInput{
        Name: "Gerard",
        Surname: "Depardieu",
    })
    t.Require().NoError(err)

    em, err := t.service.GetEnergyManager(uint(1))
    t.Require().NoError(err)
    t.Equal(em.Name, "Gerard")
    t.Equal(em.Surname, "Depardieu")

    // You can create another EM with the same data
    err = t.service.CreateEnergyManager(CreateEnergyManagerInput{
        Name: "Gerard",
        Surname: "Depardieu",
    })
    t.Require().NoError(err)
    em, err = t.service.GetEnergyManager(uint(2))
    t.Require().NoError(err)
    t.Equal(em.Name, "Gerard")
    t.Equal(em.Surname, "Depardieu")

    // We've created 2 ems so far
    ems, err := t.service.GetAllEnergyManagers()
    t.Require().NoError(err)
    t.Equal(len(ems), 2)
}

func (t *MainTestSuite) TestDeleteEnergyManager() {
    // should not be able to delete a em that does not exist
    err := t.service.DeleteEnergyManager(uint(1))
    t.Require().Error(err)

    // should be able to delete one that exists
    err = t.service.CreateEnergyManager(CreateEnergyManagerInput{
        Name: "Gerard",
        Surname: "Depardieu",
    })
    t.Require().NoError(err)
    err = t.service.DeleteEnergyManager(uint(1))
    t.Require().NoError(err)
    em, err := t.service.GetEnergyManager(uint(1))
    t.Require().Error(err)
    t.Nil(em)

    // Deleting a em with plants attached to him should not delete the plants
    err = t.service.CreateEnergyManager(CreateEnergyManagerInput{
        Name: "Gerard",
        Surname: "Depardieu",
    })
    err = t.service.CreatePlant(CreatePlantInput{
        Name: "plant1",
        Address: "17 rue truc",
        MaxPower: 100,
        EnergyManagerID: 2,
    })
    t.Require().NoError(err)
    err = t.service.DeleteEnergyManager(uint(2))
    t.Require().NoError(err)
    plant, err := t.service.GetAllPlants()
    t.Require().NoError(err)
    t.NotNil(plant)
    t.Equal(len(plant), 1)
    t.Equal(plant[0].Name, "plant1")
}

func (t *MainTestSuite) TestUpdateEnergyManager() {
    // Updating an em that does not exist should return an error
    err := t.service.UpdateEnergyManager(uint(123), UpdateEnergyManagerInput{"Jacques", "Chirac"})
    t.Require().Error(err)

    // You can update an existing em
    err = t.service.CreateEnergyManager(CreateEnergyManagerInput{
        Name: "Gerard",
        Surname: "Depardieu",
    })
    t.Require().NoError(err)
    err = t.service.UpdateEnergyManager(uint(1), UpdateEnergyManagerInput{"Jacques", "Chirac"})
    t.Require().NoError(err)
    em, err := t.service.GetEnergyManager(uint(1))
    t.Require().NoError(err)
    t.Equal(em.Name, "Jacques")
    t.Equal(em.Surname, "Chirac")
}

func (t *MainTestSuite) TestGetEnergyManagerPlants() {
    // If an em has no plant assigned yet you should get an empty slice
    err := t.service.CreateEnergyManager(CreateEnergyManagerInput{
        Name: "Gerard",
        Surname: "Depardieu",
    })
    t.Require().NoError(err)
    plants, err := t.service.GetEnergyManagerPlants(uint(1))
    t.Require().NoError(err)
    t.NotNil(plants)
    t.Equal(len(plants), 0)

    // Getting plants of an unexisting em should return an error
    plants, err = t.service.GetEnergyManagerPlants(uint(123))
    t.Require().Error(err)
    t.Equal(len(plants), 0)
}

func (t *MainTestSuite) TestGetAllPlants() {
    // If there are no plants we should get no error, just an empty slice
    plants, err := t.service.GetAllPlants()
    t.Require().NoError(err)
    t.NotNil(plants)
    t.Equal(len(plants), 0)
}

func (t *MainTestSuite) TestCreatePlant() {
    // Creating a plant with a emid that does not exist should return an error
    err := t.service.CreatePlant(CreatePlantInput{
        Name: "plant1",
        Address: "17 rue truc",
        MaxPower: 100,
        EnergyManagerID: 123,
    })
    t.Require().Error(err)

    // Creating a plant with an existing emid should be ok!
    err = t.service.CreateEnergyManager(CreateEnergyManagerInput{
        Name: "Gerard",
        Surname: "Depardieu",
    })
    t.Require().NoError(err)
    err = t.service.CreatePlant(CreatePlantInput{
        Name: "plant1",
        Address: "17 rue truc",
        MaxPower: 100,
        EnergyManagerID: 1,
    })
    t.Require().NoError(err)
    plants, err := t.service.GetAllPlants()
    t.Require().NoError(err)
    t.NotNil(plants)
    t.Equal(len(plants), 1)

    // Adding more plants
    err = t.service.CreatePlant(CreatePlantInput{
        Name: "plant2",
        Address: "13 rue truc",
        MaxPower: 200,
        EnergyManagerID: 1,
    })
    t.Require().NoError(err)
    err = t.service.CreateEnergyManager(CreateEnergyManagerInput{
        Name: "Jacques",
        Surname: "Chirac",
    })
    t.Require().NoError(err)
    err = t.service.CreatePlant(CreatePlantInput{
        Name: "plant4",
        Address: "67 rue de la paix",
        MaxPower: 1001,
        EnergyManagerID: 2,
    })
    t.Require().NoError(err)
    plants, err = t.service.GetAllPlants()
    t.Require().NoError(err)
    t.NotNil(plants)
    t.Equal(len(plants), 3)
}

func (t *MainTestSuite) GetPlant() {
    // Getting a plant that does not exist should return an error
    plant, err := t.service.GetPlant(uint(123))
    t.Require().Error(err)
    t.Nil(plant)

    // Getting a plant with an existing emid should be ok!
    err = t.service.CreateEnergyManager(CreateEnergyManagerInput{
        Name: "Gerard",
        Surname: "Depardieu",
    })
    t.Require().NoError(err)
    err = t.service.CreatePlant(CreatePlantInput{
        Name: "plant1",
        Address: "17 rue truc",
        MaxPower: 100,
        EnergyManagerID: 1,
    })
    t.Require().NoError(err)
    plant, err = t.service.GetPlant(uint(1))
    t.Require().NoError(err)
    t.NotNil(plant)
    t.Equal(plant.Name, "plant1")
}

func (t *MainTestSuite) TestDeletePlant() {
    // Deleting a plant that does not exist should return an error
    err := t.service.DeletePlant(uint(123))
    t.Require().Error(err)

    // Deleting a plant with an existing emid should be ok!
    err = t.service.CreateEnergyManager(CreateEnergyManagerInput{
        Name: "Gerard",
        Surname: "Depardieu",
    })
    t.Require().NoError(err)
    err = t.service.CreatePlant(CreatePlantInput{
        Name: "plant1",
        Address: "17 rue truc",
        MaxPower: 100,
        EnergyManagerID: 1,
    })
    t.Require().NoError(err)
    err = t.service.DeletePlant(uint(1))
    t.Require().NoError(err)
    plant, err := t.service.GetPlant(uint(1))
    t.Require().Error(err)
    t.Nil(plant)
    plants, err := t.service.GetEnergyManagerPlants(uint(1))
    t.Require().NoError(err)
    t.Equal(len(plants), 0)

    // Deleting a plant should delete all its assets
    err = t.service.CreatePlant(CreatePlantInput{
        Name: "plant1",
        Address: "17 rue truc",
        MaxPower: 100,
        EnergyManagerID: 1,
    })
    t.Require().NoError(err)
    err = t.service.CreateAsset(uint(2), CreateAssetInput{
        Name: "asset1",
        MaxPower: 10,
        Type: "furnace",
    })
    t.Require().NoError(err)
    err = t.service.CreateAsset(uint(2), CreateAssetInput{
        Name: "asset2",
        MaxPower: 10,
        Type: "furnace",
    })
    t.Require().NoError(err)
    err = t.service.DeletePlant(uint(2))
    t.Require().NoError(err)
    asset, err := t.service.DB.GetAssetById(uint(1))
    t.Require().Error(err)
    t.Nil(asset)
    asset, err = t.service.DB.GetAssetById(uint(2))
    t.Require().Error(err)
    t.Nil(asset)
    plant, err = t.service.GetPlant(uint(1))
    t.Require().Error(err)
    t.Nil(plant)
}

func (t *MainTestSuite) TestUpdatePlant() {
    // Updating a plant that does not exist should return an error
    err := t.service.UpdatePlant(uint(123), UpdatePlantInput{
        Name: "plantDeOuf",
        Address: "Mars a droite",
        MaxPower: 1234,
        EnergyManagerID: 2,
    })
    t.Require().Error(err)

    // Assigning a plant to a em that does not exist should return an error
    err = t.service.CreateEnergyManager(CreateEnergyManagerInput{
        Name: "Gerard",
        Surname: "Depardieu",
    })
    t.Require().NoError(err)
    err = t.service.CreatePlant(CreatePlantInput{
        Name: "plant1",
        Address: "17 rue truc",
        MaxPower: 100,
        EnergyManagerID: 1,
    })
    t.Require().NoError(err)
    err = t.service.UpdatePlant(uint(1), UpdatePlantInput{
        Name: "plantDeOuf",
        Address: "Mars a droite",
        MaxPower: 1234,
        EnergyManagerID: 22,
    })
    t.Require().Error(err)

    // Assigning a plant to a em that does exist should not return an error
    err = t.service.CreateEnergyManager(CreateEnergyManagerInput{
        Name: "Jacques",
        Surname: "Chirac",
    })
    err = t.service.UpdatePlant(uint(1), UpdatePlantInput{
        Name: "plantDeOuf",
        Address: "Mars a droite",
        MaxPower: 1234,
        EnergyManagerID: 2,
    })
    t.Require().NoError(err)
    plants, err := t.service.GetEnergyManagerPlants(uint(2))
    t.Require().NoError(err)
    t.Equal(len(plants), 1)
    plants, err = t.service.GetEnergyManagerPlants(uint(1))
    t.Require().NoError(err)
    t.Equal(len(plants), 0)

    // Updating max power of a plant that breaks its power constraint should return an error
    err = t.service.CreateAsset(uint(1), CreateAssetInput{
        Name: "asset1",
        MaxPower: 10,
        Type: "furnace",
    })
    t.Require().NoError(err)
    err = t.service.UpdatePlant(uint(1), UpdatePlantInput{
        Name: "plantDeOuf",
        Address: "Mars a droite",
        MaxPower: 1,
        EnergyManagerID: 2,
    })
    t.Require().Error(err)
}

func (t *MainTestSuite) TestGetPlantAssets() {
    // Getting assets from a plant that does not exist should return an error
    assets, err := t.service.GetPlantAssets(uint(112))
    t.Require().Error(err)
    t.Nil(assets)
    t.Equal(len(assets), 0)

    // A plant with no asset should return an empty slice
    err = t.service.CreateEnergyManager(CreateEnergyManagerInput{
        Name: "Gerard",
        Surname: "Depardieu",
    })
    t.Require().NoError(err)
    err = t.service.CreatePlant(CreatePlantInput{
        Name: "plant1",
        Address: "17 rue truc",
        MaxPower: 100,
        EnergyManagerID: 1,
    })
    t.Require().NoError(err)
    assets, err = t.service.GetPlantAssets(uint(1))
    t.Require().NoError(err)
    t.NotNil(assets)
    t.Equal(len(assets), 0)
}

func (t *MainTestSuite) TestCreateAsset() {
    // Creating an asset that is not the right type should raise an error
    err := t.service.CreateAsset(uint(1), CreateAssetInput{
        Name: "asset1",
        MaxPower: 10,
        Type: "eau",
    })
    t.Require().Error(err)

    // Creating an asset to an unexisting plant should retiurn an error
    err = t.service.CreateAsset(uint(1), CreateAssetInput{
        Name: "asset1",
        MaxPower: 10,
        Type: "furnace",
    })
    t.Require().Error(err)
    err = t.service.CreateEnergyManager(CreateEnergyManagerInput{
        Name: "Gerard",
        Surname: "Depardieu",
    })
    t.Require().NoError(err)
    err = t.service.CreatePlant(CreatePlantInput{
        Name: "plant1",
        Address: "17 rue truc",
        MaxPower: 100,
        EnergyManagerID: 1,
    })
    err = t.service.CreateAsset(uint(123), CreateAssetInput{
        Name: "asset1",
        MaxPower: 10,
        Type: "furnace",
    })
    t.Require().Error(err)

    // Creating an asset that overpass the plant max power should return an error
    err = t.service.CreateAsset(uint(1), CreateAssetInput{
        Name: "asset1",
        MaxPower: 101,
        Type: "furnace",
    })
    t.Require().Error(err)
    err = t.service.CreateAsset(uint(1), CreateAssetInput{
        Name: "asset1",
        MaxPower: 99,
        Type: "furnace",
    })
    t.Require().NoError(err)
    asset, err := t.service.GetPlantAsset(uint(1), uint(1))
    t.Require().NoError(err)
    t.Equal(asset.Name, "asset1")
    assets, err := t.service.GetPlantAssets(uint(1))
    t.Require().NoError(err)
    t.Equal(len(assets), 1)
    err = t.service.CreateAsset(uint(1), CreateAssetInput{
        Name: "asset2",
        MaxPower: 2,
        Type: "compressor",
    })
    t.Require().Error(err)
    assets, err = t.service.GetPlantAssets(uint(1))
    t.Require().NoError(err)
    t.Equal(len(assets), 1)
}

func (t *MainTestSuite) TestGetPlantAsset() {
    // Getting asset from a plant that does not exist should return an error
    asset, err := t.service.GetPlantAsset(uint(1), uint(1))
    t.Require().Error(err)
    t.Nil(asset)

    err = t.service.CreateEnergyManager(CreateEnergyManagerInput{
        Name: "Gerard",
        Surname: "Depardieu",
    })
    t.Require().NoError(err)
    err = t.service.CreatePlant(CreatePlantInput{
        Name: "plant1",
        Address: "17 rue truc",
        MaxPower: 100,
        EnergyManagerID: 1,
    })
    err = t.service.CreateAsset(uint(1), CreateAssetInput{
        Name: "asset1",
        MaxPower: 10,
        Type: "furnace",
    })
    t.Require().NoError(err)
    asset, err = t.service.GetPlantAsset(uint(2), uint(1))
    t.Require().Error(err)
    t.Nil(asset)

    // Getting asset that does not exist from a plant that does
    asset, err = t.service.GetPlantAsset(uint(1), uint(111))
    t.Require().Error(err)
    t.Nil(asset)
}

func (t *MainTestSuite) TestDeletePlantAsset() {
    // Deleting asset from a plant that does not exist should return an error
    err := t.service.DeletePlantAsset(uint(1), uint(1))
    t.Require().Error(err)

    err = t.service.CreateEnergyManager(CreateEnergyManagerInput{
        Name: "Gerard",
        Surname: "Depardieu",
    })
    t.Require().NoError(err)
    err = t.service.CreatePlant(CreatePlantInput{
        Name: "plant1",
        Address: "17 rue truc",
        MaxPower: 100,
        EnergyManagerID: 1,
    })
    t.Require().NoError(err)
    err = t.service.CreateAsset(uint(1), CreateAssetInput{
        Name: "asset1",
        MaxPower: 10,
        Type: "furnace",
    })
    t.Require().NoError(err)
    err = t.service.DeletePlantAsset(uint(2), uint(1))
    t.Require().Error(err)

    // Deleting asset that does not exist from a plant that does
    err = t.service.DeletePlantAsset(uint(1), uint(111))
    t.Require().Error(err)

    // Deleting asset that does exist but that does not belong to the plant
    err = t.service.CreateEnergyManager(CreateEnergyManagerInput{
        Name: "Jacques",
        Surname: "Check",
    })
    t.Require().NoError(err)
    err = t.service.CreatePlant(CreatePlantInput{
        Name: "plant2",
        Address: "Le lune",
        MaxPower: 101,
        EnergyManagerID: 2,
    })
    t.Require().NoError(err)
    err = t.service.CreateAsset(uint(2), CreateAssetInput{
        Name: "asset1",
        MaxPower: 10,
        Type: "furnace",
    })
    t.Require().NoError(err)
    err = t.service.DeletePlantAsset(uint(1), uint(2))
    t.Require().Error(err)
}

func (t *MainTestSuite) TestUpdatePlantAsset() {
    // Updating asset from a plant that does not exist should return an error
    err := t.service.UpdatePlantAsset(uint(1), uint(1), UpdateAssetInput{
        Name: "asset123",
        MaxPower: 67,
        Type: "chiller",
    })
    t.Require().Error(err)

    // Updating asset with a new invalid type should return an error
    err = t.service.UpdatePlantAsset(uint(1), uint(1), UpdateAssetInput{
        Name: "asset123",
        MaxPower: 67,
        Type: "chill",
    })
    t.Require().Error(err)

    err = t.service.CreateEnergyManager(CreateEnergyManagerInput{
        Name: "Gerard",
        Surname: "Depardieu",
    })
    t.Require().NoError(err)
    err = t.service.CreatePlant(CreatePlantInput{
        Name: "plant1",
        Address: "17 rue truc",
        MaxPower: 100,
        EnergyManagerID: 1,
    })
    t.Require().NoError(err)
    err = t.service.CreateAsset(uint(1), CreateAssetInput{
        Name: "asset1",
        MaxPower: 10,
        Type: "furnace",
    })
    t.Require().NoError(err)
    err = t.service.UpdatePlantAsset(uint(1), uint(1), UpdateAssetInput{
        Name: "asset123",
        MaxPower: 67,
        Type: "chill",
    })
    t.Require().Error(err)

    // Updating asset from a plant that does exist but the asset does not exist
    err = t.service.UpdatePlantAsset(uint(1), uint(112), UpdateAssetInput{
        Name: "asset123",
        MaxPower: 67,
        Type: "chiller",
    })
    t.Require().Error(err)

    // Updating asset that breaks the plant power limit should return an error
    err = t.service.UpdatePlantAsset(uint(1), uint(1), UpdateAssetInput{
        Name: "asset123",
        MaxPower: 1000,
        Type: "chiller",
    })
    t.Require().Error(err)
}
