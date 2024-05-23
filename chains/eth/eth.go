package eth

import (
	"math/big"

	"github.com/airsigner/libcrypto/types"
	"github.com/shopspring/decimal"
)

type ethDefinition struct{}

func (ethDefinition) CoinName() string { return "ETH" }
func (ethDefinition) UnitExp() int32   { return 18 }

type Eth struct {
	*types.CoinValue[ethDefinition]
}

func NewEth(ether decimal.Decimal) *Eth {
	return &Eth{
		types.NewCoinValueFromCoins[ethDefinition](ether),
	}
}

func NewEthFromWei(wei *big.Int) *Eth {
	return &Eth{
		types.NewCoinValue[ethDefinition](wei),
	}
}

func NewEthFromKWei(kwei decimal.Decimal) *Eth {
	return &Eth{
		types.NewCoinValueFromScaled[ethDefinition](kwei, 3),
	}
}

func NewEthFromMWei(mwei decimal.Decimal) *Eth {
	return &Eth{
		types.NewCoinValueFromScaled[ethDefinition](mwei, 6),
	}
}

func NewEthFromGWeil(gwei decimal.Decimal) *Eth {
	return &Eth{
		types.NewCoinValueFromScaled[ethDefinition](gwei, 9),
	}
}

// Wei returns the value of the Eth type in Wei.
func (e Eth) Wei() *big.Int {
	return e.Units()
}

// KWei returns the value of the Eth type in KWei.
func (e Eth) KWei() decimal.Decimal {
	return e.ScaledValue(3)
}

// MWei returns the value of the Eth type in MWei.
func (e Eth) MWei() decimal.Decimal {
	return e.ScaledValue(6)
}

// GWei returns the value of the Eth type in GWei.
func (e Eth) GWei() decimal.Decimal {
	return e.ScaledValue(9)
}

// Eth returns the value of the Eth type in Ether
func (e Eth) Eth() decimal.Decimal {
	return e.Coins()
}
