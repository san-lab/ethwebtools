package create

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const NETWORK = "https://ethereum-holesky.publicnode.com/"

// using go-ethereum client dial and get the given address nonce
// if there is an error - return it as a string
func GetNonce(address string) string {
	client, err := ethclient.Dial(NETWORK)
	if err != nil {
		return err.Error()
	}
	defer client.Close()
	adr := common.HexToAddress(address)
	if err != nil {
		return err.Error()
	}
	nonce, err := client.PendingNonceAt(context.Background(), adr)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%v", nonce)
}
