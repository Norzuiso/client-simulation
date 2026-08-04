package states

// ===============StateStudying=================

type StateStudying struct {
	StateMath
}

// GetState implements [State].
func (s StateStudying) GetState() string {
	return Studying
}

func NewStateStudying() *StateStudying {
	s := &StateStudying{}
	s.infCoef = 0.1
	s.transStateCoef = map[string]float64{}
	s.transStateCoef[Studying] = 0
	s.transStateCoef[Boring] = 0
	s.transStateCoef[Playing] = 0
	return s
}

// ========================================
