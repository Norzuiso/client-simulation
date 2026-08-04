package devs

type Atomic interface {
	DeltInt()          // Internal transition function.
	DeltExt(e float64) // External transition function.
	DeltCon(e float64) // Confluent transition function.
	Lambda()           // Output function.
}
