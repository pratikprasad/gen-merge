package teststructs

type NullFloat struct {
	Valid bool
	Float float64
}

type Person struct {
	Name           string
	Age            int64
	secretIdentity string

	favoriteCity *string

	Height NullFloat

	Friends []Person
}
