package types

import (
	"math/big"

	"github.com/shopspring/decimal"
)

type Value interface {
	Units() *big.Int
	Coins() decimal.Decimal
	ScaledValue(exp int32) decimal.Decimal

	CoinName() string

	Same(other Value) bool

	Add(other Value) Value
	Sub(other Value) Value
	Mul(other Value) Value
	Div(other Value) Value

	MulScalar(scalar *big.Int) Value
	DivScalar(scalar *big.Int) Value
}

type ValueDefinition interface {
	// returns the name of the coin
	// To be used to determine if we are talking about the same coin
	CoinName() string

	// returns the exp of the number of units in one coint
	// e.g. for Ethereum, 10^18 wei per Eth, so this should return 18
	UnitExp() int32
}

type CoinValue[D ValueDefinition] struct {
	def   D
	value *big.Int
}

func NewCoinValue[D ValueDefinition](value *big.Int) *CoinValue[D] {
	return &CoinValue[D]{
		value: func() *big.Int {
			if value == nil {
				return big.NewInt(0)
			}
			return value
		}(),
	}
}

func NewCoinValueFromCoins[D ValueDefinition](value decimal.Decimal) *CoinValue[D] {
	cv := NewCoinValue[D](nil)
	cv.value = value.Mul(decimal.New(1, cv.def.UnitExp())).BigInt()
	return cv
}

func NewCoinValueFromScaled[D ValueDefinition](value decimal.Decimal, exp int32) *CoinValue[D] {
	cv := NewCoinValue[D](nil)
	cv.value = value.Mul(decimal.New(1, cv.def.UnitExp()-exp)).BigInt()
	return cv
}

// Units returns the value of the CoinValue in the smallest unit.
//
// For example for Ethereum this would return the value denominated in wei.
//
// Returns:
// - *big.Int: the value of the CoinValue in the smallest unit..
func (v CoinValue[D]) Units() *big.Int {
	return v.value
}

// Coins returns the value of the CoinValue in whole coin units.
//
// For example for Ethereum this would return the value denomitated in Ether.
//
// Returns:
// - decimal.Decimal: The value of the CoinValue in whole coind units.
func (v CoinValue[D]) Coins() decimal.Decimal {
	return decimal.NewFromBigInt(v.value, 0).DivRound(decimal.New(1, v.def.UnitExp()), v.def.UnitExp())
}

// ScaledValue returns the value of the CoinValue in decimal form, scaled by the given exponent.
//
// For example for Ethereum the exponent value 9 would return the value denomitated in Gwei.
//
// Parameters:
// - exp: the exponent to scale the value by.
//
// Returns:
// - decimal.Decimal: the scaled value of the CoinValue.
func (v CoinValue[D]) ScaledValue(exp int32) decimal.Decimal {
	return decimal.NewFromBigInt(v.value, 0).DivRound(decimal.New(1, exp), v.def.UnitExp())
}

// CoinName returns the name of the coin associated with the CoinValue.
func (v CoinValue[D]) CoinName() string {
	return v.def.CoinName()
}

// Same checks if the CoinValue is the same as another Value by comparing their coin names.
//
// Parameters:
// - other: the Value to compare with.
//
// Returns:
// - bool: true if the coin names are the same, false otherwise.
func (v CoinValue[D]) Same(other Value) bool {
	return v.CoinName() == other.CoinName()
}

// Add adds the value of another CoinValue to the current CoinValue.
//
// It takes a Value as a parameter and returns a Value.
// The function checks if the current CoinValue and the other Value have the same coin name.
// If they don't, it panics with the message "cannot add values of different coins".
// If they are the same, it creates a new CoinValue with the same definition and adds the units of the other Value to the current CoinValue's value.
// The function returns the new CoinValue.
//
// Parameters:
// - other: the Value to add with.
//
// Returns:
// - Value: the new CoinValue after the addition.
func (v *CoinValue[D]) Add(other Value) Value {
	if !v.Same(other) {
		panic("cannot add values of different coins")
	}

	return &CoinValue[D]{
		def:   v.def,
		value: new(big.Int).Add(v.value, other.Units()),
	}
}

// Sub subtracts the value of another CoinValue from the current CoinValue.
//
// It takes a Value as a parameter and returns a Value.
// The function checks if the current CoinValue and the other Value have the same coin name.
// If they don't, it panics with the message "cannot subtract values of different coins".
// If they are the same, it creates a new CoinValue with the same definition and subtracts the units of the other Value from the current CoinValue's value.
// The function returns the new CoinValue.
//
// Parameters:
// - other: the Value to substract with.
//
// Returns:
// - Value: the new CoinValue after the subtraction.
func (v *CoinValue[D]) Sub(other Value) Value {
	if !v.Same(other) {
		panic("cannot subtract values of different coins")
	}

	return &CoinValue[D]{
		def:   v.def,
		value: new(big.Int).Sub(v.value, other.Units()),
	}
}

// Mul multiplies the value of another CoinValue with the current CoinValue.
//
// It takes a Value as a parameter and returns a Value.
// The function checks if the current CoinValue and the other Value have the same coin name.
// If they don't, it panics with the message "cannot multiply values of different coins".
// If they are the same, it creates a new CoinValue with the same definition and multiplies the units of the other Value with the current CoinValue's value.
// The function returns the new CoinValue.
//
// Parameters:
// - other: the Value to multiply with.
//
// Returns:
// - Value: the new CoinValue after the multiplication.
func (v *CoinValue[D]) Mul(other Value) Value {
	if !v.Same(other) {
		panic("cannot multiply values of different coins")
	}

	return &CoinValue[D]{
		def:   v.def,
		value: new(big.Int).Mul(v.value, other.Units()),
	}
}

// Div divides the value of a CoinValue by another Value.
//
// It takes a Value as a parameter and returns a Value.
// The function checks if the current CoinValue and the other Value have the same coin name.
// If they don't, it panics with the message "cannot divide values of different coins".
// If they are the same, it creates a new CoinValue with the same definition and divides the units of the current CoinValue's value by the units of the other Value.
// The function returns the new CoinValue.
//
// Parameters:
// - other: the Value to divide with.
//
// Returns:
// - Value: the new CoinValue after the division.
func (v *CoinValue[D]) Div(other Value) Value {
	if !v.Same(other) {
		panic("cannot divide values of different coins")
	}

	return &CoinValue[D]{
		def:   v.def,
		value: new(big.Int).Div(v.value, other.Units()),
	}
}

// MulScalar multiplies the value of a CoinValue by a scalar value.
//
// It takes a pointer to a big.Int as a parameter and returns a Value.
// The function creates a new CoinValue with the same definition and multiplies the units of the current CoinValue's value by the scalar value.
// The function returns the new CoinValue.
//
// Parameters:
// - scalar: a pointer to a big.Int representing the scalar value to multiply with.
//
// Returns:
// - Value: the new CoinValue after the multiplication.
func (v *CoinValue[D]) MulScalar(scalar *big.Int) Value {
	return &CoinValue[D]{
		def:   v.def,
		value: new(big.Int).Mul(v.value, scalar),
	}
}

// DivScalar divides the value of a CoinValue by a scalar value.
//
// It takes a pointer to a big.Int as a parameter and returns a Value.
// The function creates a new CoinValue with the same definition and divides the units of the current CoinValue's value by the scalar value.
// The function returns the new CoinValue.
//
// Parameters:
// - scalar: a pointer to a big.Int representing the scalar value to divide with.
//
// Returns:
// - Value: the new CoinValue after the division.
func (v *CoinValue[D]) DivScalar(scalar *big.Int) Value {
	return &CoinValue[D]{
		def:   v.def,
		value: new(big.Int).Div(v.value, scalar),
	}
}
