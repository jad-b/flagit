package flagit

import "time"

// Test Structs
type ChipotleOrder struct {
	Rice             string
	Beans            string
	FajitaVegetables bool
	Meat             string
	Salsa            []string
	Corn             bool
	SourCream        bool
	Cheese           bool
	Guacamole        bool
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
