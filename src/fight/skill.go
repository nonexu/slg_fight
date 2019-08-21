package fight

import (
	"fmt"
)

type AtkAction struct {
	Action []int16 //效果，伤害，闪避， 反击等
	Damage int64
}

type behaviorFunc func(atkCard *CardInfo, atkSkill *AtkSkill, fight *FightBattle) int16

//新的技能map
var SkillBehavior map[int16]behaviorFunc

func init() {
	SkillBehavior = make(map[int16]behaviorFunc)
	registerSkillFunc(ATK_TYPE_NORMAL_ATK, NormalAtkBehavior)
	registerSkillFunc(ATK_TYPE_SKILL_DIRECT, SkillDirectAtk)
	registerSkillFunc(ATK_TYPE_SKILL_DIRECT_CURE, SkillCure)
}

//同一种效果的攻击， 可以用同样的回调函数处理
func registerSkillFunc(atkType int16, fn behaviorFunc) {
	SkillBehavior[atkType] = fn
}

func InvokeSkill(atkType int16, atkCard *CardInfo, atkSkill *AtkSkill, fight *FightBattle) int16 {
	fun, ok := SkillBehavior[atkType]
	if !ok {
		fmt.Println("can't find atkType:", atkType)
		return ERR_INVALID_SKILL
	}
	return fun(atkCard, atkSkill, fight)
}

func pktNoTargetAtkDetail(atkCard *CardInfo, atkSkill *AtkSkill, status int16) *AtkDetail {
	info := &AtkDetail{
		AtkUserId: atkCard.UserId,
		AtkCard:   atkCard.CardId,
		//AtkType   int16 //攻击类型，普通攻击，技能攻击，反击
		AtkStatus:  make([]int16, 0),
		SkillId:    atkSkill.SkillId,      //特效类型
		ActionType: make([]int16, 0),      //是否触发闪避类
		Trigger:    make([]*AtkDetail, 0), //触发的技能
	}
	info.AtkStatus = append(info.AtkStatus, status)
	return info
}

func pktAtkDetail(atkCard *CardInfo, defCard *CardInfo, atkSkill *AtkSkill, action *AtkAction) *AtkDetail {
	return &AtkDetail{
		AtkUserId: atkCard.UserId,
		DefUserId: defCard.UserId,

		AtkCard: atkCard.CardId,
		DefCard: defCard.CardId,
		//AtkType  : , //攻击类型，普通攻击，技能攻击，反击
		SkillId:    atkSkill.SkillId, //特效类型
		ActionType: action.Action,
		LoseHp:     action.Damage,
		FinalHp:    defCard.Hp,
		Trigger:    make([]*AtkDetail, 0), //触发的技能
	}
}

func pktAtkAction(action []int16, damage int64) *AtkAction {
	return &AtkAction{
		Action: action,
		Damage: damage,
	}
}

func NormalAtkBehavior(atkCard *CardInfo, atkSkill *AtkSkill, fight *FightBattle) int16 {
	if atkCard.Dead() {
		return OK
	}
	targets := fight.GetAtkTarget(atkCard, atkCard.NormalAtkDis(), 1) //攻击目标数量为1
	//攻击范围类没有攻击目标
	if len(targets) == 0 {
		fight.Process[fight.Round-1].AtkTree = append(fight.Process[fight.Round-1].AtkTree,
			pktNoTargetAtkDetail(atkCard, atkSkill, ACTION_NO_TARGET))
	}

	for _, card := range targets {
		ret, atkInfo := doNormalAtk(atkCard, card, atkSkill)
		if ret == OK {
			fight.Process[fight.Round-1].AtkTree = append(fight.Process[fight.Round-1].AtkTree, atkInfo)
		}
	}
	return OK
}

func doNormalAtk(atkCard *CardInfo, defCard *CardInfo, atkSkill *AtkSkill) (int16, *AtkDetail) {
	damage := atkCard.NormalDamage()
	action := make([]int16, 0)
	if defCard.TriggerDodge() {
		action = append(action, int16(ACTION_MISS))
		damage = 0
	}

	defCard.LoseHp(damage)
	atkDetail := pktAtkDetail(atkCard, defCard, atkSkill, pktAtkAction(action, damage))

	if defCard.TriggerFightBack() {
		ret, fightBack := FightBackAtk(defCard, atkCard, &AtkSkill{ATK_TYPE_NORMAL_ATK, FIGHT_BACK_ATK})
		if ret == OK {
			action = append(action, int16(ACTION_FIGHT_BACK))
			atkDetail.Trigger = append(atkDetail.Trigger, fightBack)
		}
	}
	atkDetail.ActionType = action
	return OK, atkDetail
}

//反击效果
func FightBackAtk(atkCard *CardInfo, defCard *CardInfo, atkSkill *AtkSkill) (int16, *AtkDetail) {
	damage := atkCard.NormalDamage()
	action := make([]int16, 0)
	if defCard.TriggerDodge() {
		action = append(action, int16(ACTION_MISS))
		damage = 0
	}

	defCard.LoseHp(damage)
	atkDetail := pktAtkDetail(atkCard, defCard, atkSkill, pktAtkAction(action, damage))
	return OK, atkDetail
}

//技能直接攻击
func SkillDirectAtk(atkCard *CardInfo, atkSkill *AtkSkill, fight *FightBattle) int16 {
	if !atkCard.SkillTrigger(atkSkill.SkillId) {
		return OK
	}

	targets := fight.GetAtkTarget(atkCard, atkCard.SkillAtkDis(atkSkill.SkillId), atkCard.SkillTargetNum(atkSkill.SkillId))
	//fmt.Println("SKILL NUM", atkCard.SkillAtkDis(atkSkill.SkillId), atkCard.SkillTargetNum(atkSkill.SkillId), targets)
	//攻击范围类没有攻击目标
	if len(targets) == 0 {
		fight.Process[fight.Round-1].AtkTree = append(fight.Process[fight.Round-1].AtkTree,
			pktNoTargetAtkDetail(atkCard, atkSkill, ACTION_NO_TARGET))
	}

	for _, card := range targets {
		ret, atkInfo := doSkillDirectAttack(atkCard, card, atkSkill)
		if ret == OK {
			fight.Process[fight.Round-1].AtkTree = append(fight.Process[fight.Round-1].AtkTree, atkInfo)
		}
	}
	return OK
}

func doSkillDirectAttack(atkCard *CardInfo, defCard *CardInfo, atkSkill *AtkSkill) (int16, *AtkDetail) {
	action := make([]int16, 0)
	damage := atkCard.GetSkillDamage(atkSkill.SkillId)
	if defCard.TriggerDodge() {
		action = append(action, int16(ACTION_MISS))
		damage = 0
	}
	defCard.LoseHp(damage)
	atkDetail := pktAtkDetail(atkCard, defCard, atkSkill, pktAtkAction(action, damage))
	return OK, atkDetail
}

func SkillCure(atkCard *CardInfo, atkSkill *AtkSkill, fight *FightBattle) int16 {
	if !atkCard.SkillTrigger(atkSkill.SkillId) {
		return OK
	}

	targets := fight.GetCureTarget(atkCard, atkCard.SkillAtkDis(atkSkill.SkillId), atkCard.SkillTargetNum(atkSkill.SkillId))

	if len(targets) == 0 {
		fight.Process[fight.Round-1].AtkTree = append(fight.Process[fight.Round-1].AtkTree,
			pktNoTargetAtkDetail(atkCard, atkSkill, ACTION_NO_TARGET))
	}

	for _, card := range targets {
		ret, atkInfo := doSkillDirectCure(atkCard, card, atkSkill)
		if ret == OK {
			fight.Process[fight.Round-1].AtkTree = append(fight.Process[fight.Round-1].AtkTree, atkInfo)
		}
	}
	return OK
}

//治愈
func doSkillDirectCure(atkCard *CardInfo, defCard *CardInfo, atkSkill *AtkSkill) (int16, *AtkDetail) {
	action := make([]int16, 0)
	damage := atkCard.GetSkillDamage(atkSkill.SkillId)
	defCard.LoseHp(-damage)
	atkDetail := pktAtkDetail(atkCard, defCard, atkSkill, pktAtkAction(action, -damage))
	return OK, atkDetail
}
