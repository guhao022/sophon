package sophon

import (
	"testing"
	"fmt"
)

func TestChem_Generate(t *testing.T) {
	c := NewChemical()

	c.Generate(water(), "水")

	fmt.Println(c.Name())
}

func water() map[ChemicalNum]int {
	return map[ChemicalNum]int{
		H: 2,
		O: 1,
	}
}
