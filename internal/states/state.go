package states

import "math"

const (
	Playing  = "Playing"  // 1
	Boring   = "Boring"   // 2
	Studying = "Studying" // 3
)

func GetStateByArgs(arg string) string {
	switch arg {
	case "1", "p":
		return Playing
	case "2", "b":
		return Boring
	case "3", "s":
		return Studying
	}
	return Studying
}

type State interface {
	GetState() string
	AddInfluenceFromState(state string, inf float64)
	GetNextState() string
	GetTransCoef() float64
}

type StateMath struct {
	infCoef   float64
	transCoef float64

	transStateCoef map[string]float64
}

func (s StateMath) GetState() string {
	return ""
}

func (s StateMath) AddInfluenceFromState(state string, inf float64) {
	s.transStateCoef[state] += inf
}

func (s StateMath) GetTransCoef() float64 {
	return s.transCoef
}

func (s StateMath) GetNextState() string {
	max := math.Inf(-1)
	nextStr := ""
	for stateStr, coef := range s.transStateCoef {
		if coef > max {
			max = coef
			nextStr = stateStr
		}
	}
	return nextStr
}
