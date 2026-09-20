package state_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/nietaki/stingray/internal/state"
	. "github.com/onsi/gomega"
)

type SampleState struct {
	String string
	Int    int32
	Float  float32
	Bool   bool
	Arr    [3]string
}

// setup replaces the old BeforeEach/AfterEach hooks: it moves to the parent of
// the package directory (so the manager's relative "state/<name>.json" paths
// resolve), stages a unique name for the test, and removes any state files
// written during the test. It returns the generated name.
//
// go test runs with the working directory set to the package directory, so
// t.Chdir("..") is equivalent to the old runtime.Caller dance. As with the old
// BeforeEach, the name is only used by tests that opt into it.
func setup(t *testing.T) string {
	t.Helper()
	t.Chdir("..")

	name := uuid.NewString()
	t.Cleanup(func() {
		// The default-named file is only removed if some test actually saved
		// state under it (i.e. without calling SetName first); os.Remove
		// ignores missing files.
		os.Remove(filepath.Join("state", "state_test-SampleState.json"))
		os.Remove(filepath.Join("state", name+".json"))
	})
	return name
}

func TestHashing(t *testing.T) {
	setup(t)

	t.Run("spawns unchanged", func(t *testing.T) {
		g := NewWithT(t)

		sm := state.NewStateManager[SampleState]()
		g.Expect(sm.CheckChanged()).To(BeFalse())

		sm.Data.String = "foo"
		g.Expect(sm.CheckChanged()).To(BeTrue())
		g.Expect(sm.CheckChanged()).To(BeFalse())
	})
}

func TestPersistance(t *testing.T) {
	name := setup(t)

	t.Run("can figure out its name", func(t *testing.T) {
		g := NewWithT(t)

		sm := state.NewStateManager[SampleState]()
		g.Expect(sm.GetName()).To(Equal("state_test-SampleState"))

		sm.SetName(name)
		g.Expect(sm.GetName()).To(Equal(name))
	})

	t.Run("loads the default values if marshalled state is present", func(t *testing.T) {
		g := NewWithT(t)

		ss := SampleState{
			Int: 42,
		}
		sm := state.NewStateManager[SampleState]()
		g.Expect(sm.Data.Int).To(BeEquivalentTo(0))

		err := sm.Load()
		g.Expect(err).To(HaveOccurred())
		g.Expect(sm.Data.Int).To(BeEquivalentTo(0))

		sm.LoadOrDefault(ss)
		g.Expect(sm.Data.Int).To(BeEquivalentTo(42))
	})

	t.Run("saves its data right", func(t *testing.T) {
		g := NewWithT(t)

		ss := SampleState{
			Int: 42,
		}
		sm := state.NewStateManager[SampleState]()
		sm.LoadOrDefault(ss)
		g.Expect(sm.CheckChanged()).To(BeTrue())

		g.Expect(sm.Save()).To(Succeed())

		sm2 := state.NewStateManager[SampleState]()
		g.Expect(sm2.Data.Int).To(BeEquivalentTo(0))

		g.Expect(sm2.Load()).To(Succeed())
		g.Expect(sm2.Data.Int).To(BeEquivalentTo(42))
	})
}
