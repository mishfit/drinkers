package bottle

import (
  "fmt"
  "drinkers/philosopher"
  "testing"
)

func Test (t *testing.T) {
  b := New()
  if b.drinker != nil {
    t.Errorf("drinker should be nil")
  }
}

func TestSetDrinker (t *testing.T) {
  b := New()
  p := philosopher.New(nil, 0, 0)
  s := 0

  fmt.Printf("%v\n", b)
  fmt.Printf("%v\n", p)
  fmt.Printf("%p\n", &p)
  fmt.Printf("%T\n", philosopher.THINKING)
  fmt.Printf("%t\n", philosopher.THINKING == philosopher.State(s))

  b.SetDrinker(&p)

  fmt.Printf("%v\n", b)
  fmt.B

  if (&p != b.GetDrinker()) {
    t.Errorf("Expected GetDrinker() equal %p got %p", &p, b.GetDrinker())
  }
}
