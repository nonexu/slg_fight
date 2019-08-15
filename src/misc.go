package main 

import(
	"math/rand"
	"time"
)
var (
	rd             *rand.Rand
)

func init() {
	rd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomIds(ids []int, num int) []int{
	if len(ids) <= num {
		return ids
	}
	idMap := make(map[int]bool)
	newIds := make([]int, 0)
	leng := len(ids)
	for {
		pos := rd.Intn(leng)
		idMap[ids[pos]] = true
		if len(idMap) >= num {
			break
		}
	}

	for id, _:= range idMap{
		newIds = append(newIds, id)
	}

	return newIds
}

func RandomHappen(percent int) bool {
	return rd.Intn(100) > percent	
}



