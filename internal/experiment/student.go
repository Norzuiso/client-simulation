package experiment

import (
	"math/rand"
	"strconv"

	"github.com/Norzuiso/client/internal/states"
)

type Student struct {
	Name         string
	ExtInfluence float32
	IntChange    float32
	Seed         int64

	CurrentState states.State
	States       map[string]states.State

	ExtEvent    chan string
	EventOutput chan string
	InitState   string
}

func (s *Student) Initialize() {

	StateStudying := states.NewStateStudying()
	StatePlaying := states.NewStatePlaying()
	StateBoring := states.NewStateBoring()
	s.States = make(map[string]states.State)

	s.States[states.Studying] = StateStudying
	s.States[states.Boring] = StateBoring
	s.States[states.Playing] = StatePlaying

	s.ExtEvent = make(chan string, 50)
	s.EventOutput = make(chan string, 1)
	s.CurrentState = s.States[s.InitState]
}

func (s *Student) DeltInt() {
	nextState := s.CurrentState.GetNextState()
	if nextState != "" {
		s.CurrentState = s.States[nextState]
	}

	if nextState == s.CurrentState.GetState() {
		prob := rand.NewSource(s.Seed)
		r := rand.New(prob)
		randomProb := r.Float64()

		if randomProb > s.CurrentState.GetTransCoef() {
			randOpt := r.Intn(3)
			randOptStr := strconv.FormatInt(int64(randOpt), 10)
			s.CurrentState = s.States[states.GetStateByArgs(randOptStr)]
		}

	}
}

func (s *Student) DeltExt(e float64) {
	for event := range s.ExtEvent {
		s.CurrentState.AddInfluenceFromState(event, 0.1)
	}
}

func (s *Student) DeltCon(e float64) {
	s.DeltInt()
	s.DeltExt(0)
}

func (s *Student) Lambda() string {
	return s.CurrentState.GetState()
}
