package bottle

import (
  "drinkers/philosopher"
  "errors"
  "sync"
)

type Bottle struct {
  id int
  drinker *philosopher.Philosopher
  sync.Mutex
}

func (bottle *Bottle) Id() (int) {
  return bottle.id
}

func (bottle *Bottle) GetDrinker () (*philosopher.Philosopher) {
  return bottle.drinker
}

func (bottle *Bottle) SetDrinker (drinker *philosopher.Philosopher) (error) {
  bottle.Lock()
  if (bottle.drinker == nil || bottle.drinker.CurrentState == philosopher.THINKING) {
    bottle.drinker = drinker
  } else {
    return errors.New("bottle drinker is not in a 'THINKING' state")
  }
  bottle.Unlock()
  return nil
}

func New (id int) (Bottle) {
  return Bottle { id: id }
}

