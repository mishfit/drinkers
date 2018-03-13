package main

import (
  "drinkers/bottle"
  "drinkers/philosopher"
  "flag"
  "fmt"
  "math/rand"
  "os/signal"
  "syscall"
)

func main () {
  countOfDrinkers := flag.Int("drinkers",  5, "number of drinkers to spawn")
  countOfBottles  := flag.Int("bottles", 3, "number of bottles to spawn")
  sigs := make(chan os.Signal, 1)
  done := make(chan bool, 1)

  flag.Parse()

  fmt.Println("drinkers", *countOfDrinkers)
  fmt.Println("bottles", *countOfBottles)

  var bottles [*countOfBottles]bottle.Bottle
  var philosophers [*countOfDrinkers]*philosopher.Philosopher

  // create a specified number of bottles
  for i := 0; i < len(bottles); i++ {
    bottle[i] = bottle.New(i)
  }

  // create a specified number of philosophers
  for i := 0; i < len(philosopher); i++ {
    // assign each philosopher a random subset of bottles
    countOfRequiredBottles := rand.Intn(len(bottles)) + 1
    secondsForThinking := rand.Intn(20)
    secondsForDrinking := rand.Intn(30)
    requiredBottles := make([]bottle.Bottle, countOfRequiredBottles)
    var uniqueBottles = map[int]bottle.Bottle
    indexOfBottle := 0

    for {
      indexOfBottleToAdd := rand.Intn(len(bottles))
      _, exists := uniqueBottles[indexOfBottle]
      if (!exists) {
        uniqueBottles[indexOfBottleToAdd] = bottles[indexOfBottleToAdd]
        requiredBottles[indexOfBottle] = bottles[indexOfBottleToAdd]
      }

      if (len(uniqueBottles) == countOfRequiredBottles) {
        break
      }

      indexOfBottle++
    }

    philosopher[i] = philosopher.New(i, requiredBottles, secondsForThinking, secondsForDrinking)
  }


  signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

  go func() {
    sig := <-sigs
    fmt.Println()
    fmt.Println(sig)
    done <- true
  }()

  fmt.Println("Press Ctrl+C to exit...")
  <-done
}
