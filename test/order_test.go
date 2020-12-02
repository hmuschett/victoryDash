package test

import (
	"testing"
	"victorydash/configs"
	"victorydash/handlers"
	"victorydash/models"
)

func TestConnection(t *testing.T) {
	configs.CreateConnection()
	configs.Ping()
	defer configs.CloseConnection()
}
func TestLogin(t *testing.T) {
	configs.CreateConnection()
	defer configs.CloseConnection()

	_, err := models.Login("henry2", "1234567")
	if err != nil {
		t.Error()
	}

}
func TestGetUserByUsername(t *testing.T) {
	configs.CreateConnection()
	defer configs.CloseConnection()
	nameFile, err := models.GetUserByUsername("henry2")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(nameFile)
	}
}

/*func TestUpdateOrders(t *testing.T) {
	configs.CreateConnection()
	handlers.UpDateOrders()
	defer configs.CloseConnection()

}*/

func TestCreateCsvOrderByProvider(t *testing.T) {
	configs.CreateConnection()
	defer configs.CloseConnection()

	arr := []string{"2537442967712", "2508171051168"}
	nameFile, err := handlers.CreateCsvOrderByProvider(arr, "WERM")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(nameFile)
	}
}
