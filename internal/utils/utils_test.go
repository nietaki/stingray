package utils_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/mvrahden/go-test/pkg/gotest"
)

type ExampleTestSuite struct {
	name string
}

func (s *ExampleTestSuite) BeforeEach(t *gotest.T) {
	s.name = uuid.NewString()
}

func (s *ExampleTestSuite) TestExample(t *gotest.T) {
	t.It("fails on a broken assertion", func(it *gotest.T) {
		it.Skipf("example failing test skipped")
		gotest.Less(it, 2, 1, "this is surely incorrect")
	})
}

func (s *ExampleTestSuite) TestNonBDD(t *gotest.T) {
	t.Skipf("skipping the gotest way")
	gotest.Less(t, 2, 1, "this is surely incorrect")
}

func (s *ExampleTestSuite) TestOldschool(t *testing.T) {
	t.Skipf("skipping in an oldschool way")
	gotest.Less(t, 2, 1, "this is surely incorrect")
}

func TestEvenMoreOldschool(t *testing.T) {
	t.Skipf("skipping in an oldschool way")
	if 2 < 1 {
		t.Errorf("maths is wrong")
	}
}
