package utils_test

// A living tour of Gomega: frequent scenarios, known gotchas, and advanced
// patterns. Everything here RUNS and PASSES — traps are demonstrated by
// asserting their true (sometimes surprising) behaviour, so `go test ./...`
// is the documentation.
//
// House rules, learned the hard way:
//   - One `g := NewWithT(t)` per test (sub)test. A gomega instance is bound to
//     the `t` it was created with; reuse across parallel subtests attributes
//     failures to the wrong node.
//   - The dot-import is Gomega's documented style. If a linter ever objects,
//     quarantine this file rather than changing call sites elsewhere.
//   - gotcha #0: this v1.43 surface has NO BeCloseTo, Send, MatchFile or
//     ExpectAll. Float deltas go through BeNumerically("~"), channel "can
//     receive a send" is BeSent, and file checks are stdlib + matchers.

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"
	. "github.com/onsi/gomega"
)

// Point is a tiny all-exported struct used across examples.
type Point struct{ X, Y int }

// --- Frequent scenarios ----------------------------------------------------

func TestAssertions101(t *testing.T) {
	g := NewWithT(t)

	g.Expect(1+1).To(Equal(2), "messages are the LAST argument, printf-style: %s", "optional")
	g.Expect(true).To(BeTrue())
	g.Expect([]string{"b", "a"}).NotTo(Equal([]string{"a", "b"}), "Equal respects order")
	g.Expect([]string{"b", "a"}).To(ConsistOf("a", "b"), "ConsistOf ignores order")
	g.Expect([]int{1, 2, 3}).To(ContainElement(2))
	g.Expect([]int{1, 2, 3}).To(HaveEach(BeNumerically(">", 0)), "every element")
	g.Expect(map[string]int{"id": 7}).To(HaveKey("id"))
	g.Expect(map[string]int{"id": 7}).To(HaveKeyWithValue("id", 7))
	g.Expect("deadbeef").To(MatchRegexp(`^[0-9a-f]{8}$`))
	g.Expect("stingray").To(ContainSubstring("ing"))
	g.Expect([]string{"a"}).To(HaveLen(1))
}

func TestErrors(t *testing.T) {
	g := NewWithT(t)

	notFound := errors.New("not found")
	wrapped := fmt.Errorf("loading state: %w", notFound)
	var noErr error

	g.Expect(wrapped).To(HaveOccurred())  // non-nil
	g.Expect(noErr).NotTo(HaveOccurred()) // nil error
	g.Expect(func() error { return nil }()).To(Succeed(), "Succeed == err == nil")

	// MatchError uses errors.Is semantics — Equal would compare structurally and
	// miss through wrapping:
	g.Expect(wrapped).To(MatchError(notFound))
	// matcher arguments give substring matching through the chain:
	g.Expect(wrapped).To(MatchError(ContainSubstring("state")))

	// errors.As for typed errors — plain Go, then assert the extracted value:
	_, openErr := os.Open("/nonexistent/zz")
	var target *os.PathError
	g.Expect(errors.As(openErr, &target)).To(BeTrue())
	g.Expect(target.Path).To(Equal("/nonexistent/zz"))
}

func TestNumbersAndTimes(t *testing.T) {
	g := NewWithT(t)

	// Float delta idiom (there is no BeCloseTo in this version):
	g.Expect(3.14159).To(BeNumerically("~", 22.0/7.0, 2e-3), "22/7 is only this good an approximation")

	// ~ also bridges float32/float64, the JSON round-trip case:
	var f32 float32 = 3.14159 // float32 loses precision vs the float64 literal
	g.Expect(f32).To(BeNumerically("~", 3.14159, 1e-4))
	g.Expect(f32).NotTo(Equal(3.14159), "exact compare across widths bites")

	// time: Equal/reflect.DeepEqual looks at zone & monotonic internals —
	// same instant, different *representation*, fails:
	now := time.Now()
	g.Expect(now).NotTo(Equal(now.UTC()))
	g.Expect(now).To(BeTemporally("==", now.UTC()), "BeTemporally uses time.Equal")
	g.Expect(now).To(BeTemporally("~", now.Add(500*time.Millisecond), time.Second))
}

// --- Gotchas ----------------------------------------------------------------

func TestDeepEqualityThreeEngines(t *testing.T) {
	g := NewWithT(t)

	// Engine A: Equal == reflect.DeepEqual, TYPE-STRICT. int32 vs int never match:
	g.Expect(int32(42)).NotTo(Equal(42))
	g.Expect(Point{1, 2}).To(Equal(Point{1, 2}))

	// BeEquivalentTo converts actual to expected's type first:
	g.Expect(int32(42)).To(BeEquivalentTo(42))

	// nil vs empty slices are DIFFERENT for DeepEqual...
	var decoded []string // e.g. missing field after json.Unmarshal
	g.Expect(decoded).NotTo(Equal([]string{}))
	// ...and go-cmp can opt in to equivalence:
	g.Expect(decoded).To(BeComparableTo([]string{}, cmpopts.EquateEmpty()))

	// Engine B: BeComparableTo == go-cmp: options + a path diff on failure.
	got := Point{X: 1, Y: 3}
	g.Expect(got).To(BeComparableTo(Point{X: 1, Y: 3}))
	// gotcha: unexported fields make go-cmp PANIC unless you pass
	// cmp.AllowUnexported(T{}) or cmpopts.IgnoreUnexported(T{}) — loud failure
	// by design. DeepEqual-based matchers compare them silently instead.

	// Field surgery beats full-struct equality when only part is stable:
	g.Expect(got).To(HaveField("Y", BeNumerically(">", 2)))
}

func TestNilAndEmptiness(t *testing.T) {
	g := NewWithT(t)

	// gotcha: Equal(nil) is an ERROR in gomega ("both actual and expected are
	// nil"). Use BeNil:
	var noErr error
	g.Expect(noErr).To(BeNil())
	var ptr *Point
	g.Expect(ptr).To(BeNil())

	// BeEmpty covers nil AND zero-length AND zero-value; BeZero is stricter:
	g.Expect([]string(nil)).To(BeEmpty())
	g.Expect(map[string]int{}).To(BeEmpty())
	g.Expect(Point{}).To(BeZero())
	g.Expect(Point{X: 1}).NotTo(BeZero())
}

func TestChannelsAndAsync(t *testing.T) {
	g := NewWithT(t)

	// Eventually polls; defaults are 1s timeout / 10ms polling — override per
	// call (or globally via SetDefaultEventuallyTimeout) for flake-free CI.
	msgs := make(chan string, 1)
	go func() {
		time.Sleep(20 * time.Millisecond)
		msgs <- "hello"
		close(msgs)
	}()

	var got string
	// Receive CONSUMES the value; the pointer argument BINDS it:
	g.Eventually(msgs).
		WithTimeout(2 * time.Second).
		WithPolling(10 * time.Millisecond).
		Should(Receive(&got))
	g.Expect(got).To(Equal("hello"))

	// After the close, BeClosed holds:
	g.Eventually(msgs).Should(BeClosed())

	// "Nothing must arrive" costs the full window — Consistently can't fail
	// fast. Keep its timeout tiny:
	quiet := make(chan string, 1)
	g.Consistently(quiet).WithTimeout(50 * time.Millisecond).ShouldNot(Receive())

	// BeSent asserts a NON-BLOCKING send is possible — beware: it PERFORMS the
	// send if it succeeds, so the value is really in the buffer afterwards:
	buf := make(chan int, 1)
	g.Expect(buf).To(BeSent(42))
	g.Expect(<-buf).To(Equal(42))
}

// --- Advanced ---------------------------------------------------------------

func TestCombinatorsAndTransforms(t *testing.T) {
	g := NewWithT(t)

	g.Expect(7).To(And(BeNumerically(">", 0), BeNumerically("<=", 100)), "And() is a function, not a matcher method")
	g.Expect("stingray").To(Or(HavePrefix("sting"), HavePrefix("foo")))
	g.Expect([]int{2, 4}).To(SatisfyAll(HaveLen(2), ContainElement(4)), "SatisfyAll == And")
	g.Expect(13).To(SatisfyAny(BeNumerically("<", 0), BeNumerically(">", 10)), "SatisfyAny == Or")

	// Satisfy: ad-hoc predicate when no matcher fits:
	g.Expect(11).To(Satisfy(func(n int) bool { return n%2 == 1 && n > 10 }))

	// WithTransform: map the actual into something matchable first:
	type record struct{ Created time.Time }
	g.Expect(record{Created: time.Now()}).
		To(WithTransform(func(r record) int64 { return r.Created.Unix() }, BeNumerically(">", 0)))
}

func TestJSONAndFiles(t *testing.T) {
	g := NewWithT(t)

	// MatchJSON compares semantically — key order & indentation don't matter
	// (state-file assertion of the future):
	b, err := json.Marshal(map[string]any{"b": 1, "a": []int{1, 2}})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(string(b)).To(MatchJSON(`{"a":[1,2],"b":1}`))

	// No MatchFileContents in this version — stdlib + matchers compose fine:
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")
	g.Expect(os.WriteFile(path, b, 0o644)).To(Succeed())

	g.Expect(path).To(BeARegularFile())
	contents, err := os.ReadFile(path)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(contents).To(MatchJSON(`{"b":1,"a":[1,2]}`), "MatchJSON takes []byte too")
	g.Expect(string(contents)).To(ContainSubstring(`"a"`))
}

// --- Escape hatch -------------------------------------------------------------

func TestOldschool(t *testing.T) {
	// Gomega never demands full conversion. Plain testing remains the baseline,
	// and mixing both in one file is fine.
	if got := 2 + 2; got != 4 {
		t.Errorf("got %d, want 4", got)
	}
}
