package fight

import (
	"fmt"
)

type statusFunc func(atkCard *CardInfo, status *StatusInfo, fight *FightBattle) int16

//新的技能map
var statusBehavior map[int16]statusFunc

func registerStatusFunc(atkType int16, fn statusFunc) {
	statusBehavior[atkType] = fn
}

func init() {
	statusBehavior = make(map[int16]statusFunc)
	registerStatusFunc(ACTION_POISONED, PrisonHandle)
}

func InvokeStatus(typ int16, atkCard *CardInfo, status *StatusInfo, fight *FightBattle) int16 {
	fun, ok := statusBehavior[typ]
	if !ok {
		fmt.Println("can't find status:", typ)
		return ERR_INVALID_SKILL
	}
	return fun(atkCard, status, fight)
}

func PrisonHandle(card *CardInfo, info *StatusInfo, fight *FightBattle) int16 {
	action := make([]int16, 0)
	damage := info.Damage
	damage = card.LoseHp(damage)
	atkCard := fight.GetCard(info.AtkUserId, info.CardId)
	if atkCard == nil {
		return OK
	}
	atkCard.AddTotalData(CURE_NUM, damage)

	atkDetail := pktAtkDetail(atkCard, card, &AtkSkill{0, info.SkillId}, pktAtkAction(action, damage))
	fight.Process[fight.Round-1].AtkTree = append(fight.Process[fight.Round-1].AtkTree, atkDetail)

	delete(card.Status, info.Status)
	return OK
}
