package philosopher

import (
  "errors"
  "fmt"
  "sync"
  "time"
)

type bottle interface {
  Id() int
  GetDrinker() Philosopher
  SetDrinker(Philosopher) error
}

type State int

const (
  THINKING State = iota
  THIRSTING = iota
  DRINKING = iota
)

type Philosopher struct {
  id int
  CurrentState State
  requiredBottles []*bottle
  bottles []bottle
  bottleQueues map[int][]*Philosopher
  stateChannel chan State
  sync.RWMutex
}

func (p *Philosopher) ReleaseBottleTo (b bottle, to *Philosopher) (error) {
  p.Lock()
  if (len(p.bottles) == 0) {
    return errors.New("philosopher doesn't have any bottles")
  }
  // do I have this bottle?
  found := false
  for _, x := range p.bottles {
    if (x == b) {
      queue, exists := p.bottleQueues[b.Id()]
      p.bottleQueues[b.Id()] = append(queue, to)
    }
  }

  if (!found) {
    return errors.New("philosopher doesn't have specified bottle")
  }

  p.Unlock()
  return nil
}

func (p *Philosopher) ReceiveBottle (b bottle, queue []*Philosopher) (error) {
  // do I have this bottle?
  p.Lock()
  for _, x := range p.bottles {
    if (x == b) {
      return errors.New("philosopher already has this bottle")
    }
  }
  // am I thirsty?
  if (p.CurrentState != THIRSTING) {
    return errors.New("philosopher isn't thirsty")
  }

  p.bottles = append(p.bottles, b)
  p.bottleQueues[b.Id()] = queue

  p.Unlock()

  return nil
}

func (p *Philosopher) think(secondsForThinking int) {
  for state := range p.stateChannel {
    if (state == THINKING) {
      p.CurrentState = state
      fmt.Printf("philosopher %03d is thinking at %#v\n", p.id, time.Now())
      // send off bottles
      p.Lock()
      for i, b := range p.bottles {
        queue, exists := p.bottleQueues[b.Id()]
        if (exists) {
          for len(queue) > 0 {
            to, queue := queue[0], queue[1:]
            e := to.ReceiveBottle(b, queue)

            if (e == nil) {
              break;
            }
          }

          delete(p.bottleQueues, b.Id())
        }
      }

      p.bottles = make([]bottle, len(p.requiredBottles))

      p.Unlock()

      time.Sleep(secondsForThinking * time.Second)
      p.stateChannel <- THIRSTING
    }
  }
}

func (p *Philosopher) thirst () {
  for state := range p.stateChannel {
    if (state == THIRSTING) {
      p.CurrentState = state
      fmt.Printf("philosopher %03d is thirsty at %#v\n", p.id, time.Now())
      // should need to ask for necessary bottles only once
      for _, b := range p.requiredBottles {
        otherDrinker := b.GetDrinker()
        otherDrinker.ReleaseBottleTo(&p)
      }
    }
  }
}

func (p *Philosopher) drink (secondsForDrinking int) {
  for state := range p.stateChannel {
    if (state == DRINKING) {
      p.CurrentState = state
      fmt.Printf("philosopher %03d is drinking at %#v\n", p.id, time.Now())
      time.Sleep(secondsForDrinking * time.Second)
      p.stateChannel <- THINKING
    }
  }
}

func New (id int, bottles []bottle, secondsForThinking int, secondsForDrinking int) (Philosopher) {

  p := Philosopher { id: id, requiredBottles: bottles, stateChannel: make(chan State, 10) }

  // thinking thread
  go p.think(secondsForThinking)

  // thirsting thread
  go p.thirst()

  // drinking thread
  go p.drink(secondsForDrinking)

  p.stateChannel <- THINKING

  return p
}
