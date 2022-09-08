package entity

import (
	"errors"
)

var ErrNotEnoughMoney = errors.New("not enough money")
var ErrTransactionAlreadyReverted = errors.New("transaction already reverted")
var ErrNotEnoughFreerounds = errors.New("not enough freerounds")
