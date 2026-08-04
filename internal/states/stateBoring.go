package states

// ==============StateBoring==================

type StateBoring struct {
	StateMath
}

func (s *StateBoring) GetState() string {
	return Boring
}

func NewStateBoring() *StateBoring {
	s := &StateBoring{}
	s.infCoef = 0.1
	s.transStateCoef = map[string]float64{}
	s.transStateCoef[Studying] = 0
	s.transStateCoef[Boring] = 0
	s.transStateCoef[Playing] = 0
	return s
}

// ========================================
