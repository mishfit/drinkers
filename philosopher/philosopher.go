package philosopher

import (
  "errors"
  "fmt"
  "sync"
  "time"
)

type Bottle interface {
  Id() int
  GetDrinker() *Philosopher
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
  requiredBottles []Bottle
  bottles []Bottle
  bottleQueues map[int][]*Philosopher
  stateChannel chan State
  sync.RWMutex
}

func (p *Philosopher) ReleaseBottleTo (b Bottle, to *Philosopher) (error) {
  p.Lock()
  if (len(p.bottles) == 0) {
    return errors.New("philosopher doesn't have any bottles")
  }
  // do I have this bottle?
  found := false
  for _, x := range p.bottles {
    if (x == b) {
      queue, exists := p.bottleQueues[b.Id()]
      if (exists) {
        p.bottleQueues[b.Id()] = append(queue, to)
      } else {
        p.bottleQueues[b.Id()] = []*Philosopher{ to }
      }
    }
  }

  if (!found) {
    return errors.New("philosopher doesn't have specified bottle")
  }

  p.Unlock()
  return nil
}

func (p *Philosopher) ReceiveBottle (b Bottle, queue []*Philosopher) (error) {
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
      fmt.Printf("philosopher %03d is thinking at %s\n", p.id, time.Now().Format("20060102150405"))
      // send off bottles
      p.Lock()
      fmt.Printf("about to send bottles %v\n", p.bottles)
      for _, b := range p.bottles {
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

      p.bottles = make([]Bottle, len(p.requiredBottles))

      p.Unlock()

      time.Sleep(time.Duration(secondsForThinking) * time.Second)
      p.stateChannel <- THIRSTING
    }
  }
}

func (p *Philosopher) thirst () {
  for state := range p.stateChannel {
    if (state == THIRSTING) {
      p.CurrentState = state
      fmt.Printf("philosopher %03d is thirsty at %s\n", p.id, time.Now().Format("20060102150405"))
      // should need to ask for necessary bottles only once
      for _, b := range p.requiredBottles {
        bottle := b
        fmt.Printf("checking bottle availability %#v\n", bottle)
        otherDrinker := bottle.GetDrinker()

        if (otherDrinker == nil) {
          fmt.Printf("no one has bottle %02d, taking it immediately\n", b.Id())
          b.SetDrinker(*p)
          p.bottles = append(p.bottles, b)
        } else {
          otherDrinker.ReleaseBottleTo(b, p)
        }
      }

      if (len(p.bottles) == len(p.requiredBottles)) {
        p.stateChannel <- DRINKING
      }
    }
  }
}

func (p *Philosopher) drink (secondsForDrinking int) {
  for state := range p.stateChannel {
    if (state == DRINKING) {
      p.CurrentState = state
      fmt.Printf("philosopher %03d is drinking at %s with %#v bottles\n", p.id, time.Now().Format("20060102150405"), p.bottles)
      time.Sleep(time.Duration(secondsForDrinking) * time.Second)
      p.stateChannel <- THINKING
    }
  }
}

func New (id int, bottles []Bottle, secondsForThinking int, secondsForDrinking int) (Philosopher) {

  p := Philosopher { id: id, requiredBottles: bottles, stateChannel: make(chan State, 10) }
  //fmt.Printf("created philosopher %03d: %+v\n", id, p)

  // thinking thread
  go p.think(secondsForThinking)

  // thirsting thread
  go p.thirst()

  // drinking thread
  go p.drink(secondsForDrinking)

  p.stateChannel <- THINKING

  return p
}
