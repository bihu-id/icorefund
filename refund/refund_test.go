package refund

import "testing"

func Test_GetHolders(t *testing.T) {
	//holder := make(map[string]int64)
	//GetSenderAtBlock(start, contractAddress)

	holders := GetHolders()
	//ShowHolder(holders)
	SendEther(holders)
}

func Test_Refund(t *testing.T) {
	//holder := make(map[string]int64)
	//GetSenderAtBlock(start, contractAddress)

	// holders := GetHolders()
	// ShowHolder(holders)

	ICORefund1()

	//SendEther(holders)
}
