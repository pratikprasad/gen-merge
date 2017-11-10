package teststructs

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

	p2.Friends = []Person{p1, merged}
	merged = p1.Merge(p2)
	assert.Equal(t, 2, len(merged.Friends))
}

func TestPerson_MergeOverride(t *testing.T) {
	city := "Atlantis"
	p1 := Person{
		Name:           "Alice",
		Age:            12345,
		favoriteCity:   &city,
		Height:         NullFloat{Valid: true, Float: 45.3},
		secretIdentity: "Athena",
	}
	p2 := Person{
		secretIdentity: "Cheshire Cat",
	}
	merged := p1.MergeOverride(p2)
	assert.Equal(t, "Alice", merged.Name)
	assert.Equal(t, int64(12345), merged.Age)
	assert.Equal(t, "Atlantis", *merged.favoriteCity)
	assert.Equal(t, 45.3, merged.Height.Float)
	assert.Equal(t, true, merged.Height.Valid)
	assert.Equal(t, "Cheshire Cat", merged.secretIdentity)
}
