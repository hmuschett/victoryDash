package test

import (
	"testing"
	"victorydash/configs"
	"victorydash/handlers"
)

func TestConnection(t *testing.T) {
	configs.CreateConnection()
	configs.Ping()
	defer configs.CloseConnection()
}

/*func TestUpdateOrders(t *testing.T) {
	configs.CreateConnection()
	handlers.UpDateOrders()
	defer configs.CloseConnection()

}*/
func TestUpdateOrders(t *testing.T) {
	configs.CreateConnection()

	handlers.UpdateStatusOrder("2653824680096", "WERM", "received")
	defer configs.CloseConnection()
}
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
