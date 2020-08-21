package dog

// Breed _
type Breed string

func toBreed(val string) Breed {
	return Breed(val)
}

// Dog _
type Dog string

func toDog(val string) Dog {
	return Dog(val)
}
