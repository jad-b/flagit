package flagit

import "time"

// Test Structs
type ChipotleOrder struct {
	Rice             string
	Beans            string
	FajitaVegetables boolean
	Meat             string
	Salsa            []string
	Corn             boolean
	SourCream        boolean
	Cheese           boolean
	Guacamole        boolean
}

type Address struct {
	Street  string
	City    string
	State   int
	LongLat [2]float64
}

type PlacedOrder struct {
	ChipotleOrder
	Address
	TimeReady time.Time
}
