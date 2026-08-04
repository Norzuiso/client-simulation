package states

// =============StatePlaying===================
type StatePlaying struct {
	StateMath
}

// GetState implements [State].
func (s StatePlaying) GetState() string {
	return Playing
}

func NewStatePlaying() *StatePlaying {
	s := &StatePlaying{}
	s.infCoef = 0.1
	s.transStateCoef = map[string]float64{}
	s.transStateCoef[Studying] = 0
	s.transStateCoef[Boring] = 0
	s.transStateCoef[Playing] = 0
	return s
}

// ========================================
