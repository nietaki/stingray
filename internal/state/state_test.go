package state_test

import (
	"github.com/google/uuid"
	"github.com/mvrahden/go-test/pkg/gotest"
	"github.com/nietaki/stingray/internal/state"
)

type StateManagerTestSuite struct {
	name string
}

type SampleState struct {
	String string
	Int    int32
	Float  float32
	Bool   bool
	Arr    [3]string
}

func (s *StateManagerTestSuite) BeforeEach(t *gotest.T) {
	s.name = uuid.NewString()
}

func (s *StateManagerTestSuite) TestHashing(t *gotest.T) {
	t.It("spawns unchanged", func(it *gotest.T) {
		sm := state.NewStateManager[SampleState]()
		gotest.False(it, sm.CheckChanged())
		sm.Data.String = "foo"
		gotest.True(it, sm.CheckChanged())
		gotest.False(it, sm.CheckChanged())
	})
}

func (s *StateManagerTestSuite) TestPersistance(t *gotest.T) {
	t.It("can figure out its name", func(it *gotest.T) {
		sm := state.NewStateManager[SampleState]()
		gotest.Equal(it, "SampleState", sm.GetName())
		sm.SetName(s.name)
		gotest.Equal(it, s.name, sm.GetName())
	})

	t.It("loads the default values if marshalled state is present", func(it *gotest.T) {
		ss := SampleState{
			Int: 42,
		}
		sm := state.NewStateManager[SampleState]()
		gotest.Equal(it, 0, sm.Data.Int)
		err := sm.Load()
		gotest.Error(it, err)
		gotest.Equal(it, 0, sm.Data.Int)
		sm.LoadOrDefault(ss)
		gotest.Equal(it, 42, sm.Data.Int)
	})
}
