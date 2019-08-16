package main

import (
	"fmt"
	"time"
	"fight"
)

func main() {
	start := time.Now().UnixNano()
	fight := fight.InitFight()
	fight.StartFight()

	winner := ""
	if fight.Result() {
		winner = fight.GetAtkName()
	} else {
		winner = fight.GetDefName()
	}

	str := fmt.Sprintf("%v win!!!", winner)
	fmt.Println(str)

	fmt.Println("\nfight start")
	fight.DebugProcess()
	fmt.Println("fight end")
	end := time.Now().UnixNano()
	fmt.Println("it takes:", end-start)
}
