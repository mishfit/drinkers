package main

import (
  "drinkers/bottle"
  "drinkers/philosopher"
  "flag"
  "fmt"
  "math/rand"
  "os"
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
  fmt.Println("Press Ctrl+C to exit...")

  bottles := make([]bottle.Bottle, *countOfBottles)
  philosophers := make([]philosopher.Philosopher, *countOfDrinkers)

  // create a specified number of bottles
  for i := 0; i < len(bottles); i++ {
    b := bottle.New(i)
    bottles[i] = b
  }

  // create a specified number of philosophers
  for i := 0; i < len(philosophers); i++ {
    // assign each philosopher a random subset of bottles
    countOfRequiredBottles := rand.Intn(len(bottles)) + 1
    secondsForThinking := rand.Intn(20)
    secondsForDrinking := rand.Intn(30)
    requiredBottles := make([]philosopher.Bottle, 0)

    var uniqueBottles = make(map[int]bool)

    for len(requiredBottles) < countOfRequiredBottles {
      indexOfBottleToAdd := rand.Intn(len(bottles))
      _, exists := uniqueBottles[indexOfBottleToAdd]
      if (!exists) {
        uniqueBottles[indexOfBottleToAdd] = true
        requiredBottles = append(requiredBottles, &bottles[indexOfBottleToAdd])
      }
    }

    philosophers[i] = philosopher.New(i, requiredBottles, secondsForThinking, secondsForDrinking)
  }


  signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

  go func() {
    sig := <-sigs
    fmt.Println()
    fmt.Println(sig)
    done <- true
  }()

  <-done
}
