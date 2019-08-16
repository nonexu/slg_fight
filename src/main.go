package main 
import(
	"fmt"
	"time"
)

func main() {
	start := time.Now().UnixNano()
	fight := InitFight()
	fight.StartFight()
	if fight.Result() {
		fmt.Println("Atk win")
	} else {
		fmt.Println("Atk lose")
	}
	fmt.Println("fight start")
	fight.DebugProcess()
	fmt.Println("fight end")
	end := time.Now().UnixNano()
	fmt.Println("it takes:", end - start)
}