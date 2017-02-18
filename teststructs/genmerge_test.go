package teststructs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPerson_Merge(t *testing.T) {
	p1 := Person{}
	p2 := Person{
		Name: "Alice",
	}
	merged := p1.Merge(p2)
	assert.Equal(t, "Alice", merged.Name)

	p1.secretIdentity = "Alpaca"
	merged = p1.Merge(p2)
	assert.Equal(t, "Alpaca", merged.secretIdentity)

	p1.Height = NullFloat{
		Valid: true,
		Float: 12.4,
	}

	merged = p1.Merge(p2)
	assert.Equal(t, 12.4, merged.Height.Float)
	assert.Equal(t, true, merged.Height.Valid)

	city := "NYC"
	p2.favoriteCity = &city
	merged = p1.Merge(p2)
	assert.Equal(t, "NYC", *merged.favoriteCity)
}
