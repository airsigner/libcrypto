package eth

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	addrRegex = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
)

func IsValidAddress(address string) bool {
	return addrRegex.MatchString(address)
}

func IsSmartContract(address string, client *ethclient.Client) (bool, error) {
	return IsSmartContractCtx(context.Background(), address, client)
}

func IsSmartContractCtx(ctx context.Context, address string, client *ethclient.Client) (bool, error) {
	if !IsValidAddress(address) {
		return false, errors.New("invalid address")
	}

	addr := common.HexToAddress(address)
	byteCode, err := client.CodeAt(ctx, addr, nil)
	if err != nil {
		err = fmt.Errorf("failed to get bytecode: %w", err)
		return false, err
	}
	return len(byteCode) > 0, nil
}
