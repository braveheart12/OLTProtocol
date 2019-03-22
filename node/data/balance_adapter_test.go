package data

import "testing"

func TestCreateAdapter(t *testing.T) {

	a := NewBalanceFromString("10", "OLT")
	a.MakeDataAdapter()
}
