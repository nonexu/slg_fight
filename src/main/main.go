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


	str := fmt.Sprintf("%v 获胜!!!", winner)
	fmt.Printf("%c[%d;%d;%dm%s%c[0m", 0x1B, 0, 0, 32, str, 0x1B)

	fmt.Println("\n战斗开始")
	fight.DebugProcess()
	fmt.Println("战斗结束")
	end := time.Now().UnixNano()
	fmt.Printf("战斗耗时: %v 微秒\n", (end-start)/1000)
}
