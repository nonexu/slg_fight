package main

import (
	"fmt"
)

type AtkAction struct {
	Action int16 //效果，伤害，闪避， 反击等
	Damage int64
}

type behavior func(atkCard *CardInfo, defCard *CardInfo, atkSkill *AtkSkill) (int16, *AtkDetail)

var Skills map[int16]behavior

func init() {
	Skills = make(map[int16]behavior)
	registerSkill(ATK_TYPE_NORMAL_ATK, NormalAtk)
	registerSkill(ATK_TYPE_SKILL_DIRECT, SkillAttack)

}

func registerSkill(skillId int16, fn behavior) {
	Skills[skillId] = fn
}

func Invoke(skillId int16, atkCard *CardInfo, defCard *CardInfo, atkSkill *AtkSkill) (int16, *AtkDetail) {
	fun, ok := Skills[skillId]
	if !ok {
		fmt.Println("can't find the skill:", skillId)
		return ERR_INVALID_SKILL, nil
	}
	return fun(atkCard, defCard, atkSkill)
}

func pktAtkDetail(atkCard *CardInfo, defCard *CardInfo, atkSkill *AtkSkill, action *AtkAction) *AtkDetail {
	return &AtkDetail{
		AtkUserId: atkCard.UserId,
		DefUserId: defCard.UserId,

		AtkCard: atkCard.CardId,
		DefCard: defCard.CardId,
		//AtkType  : , //攻击类型，普通攻击，技能攻击，反击
		SkillType:  atkSkill.SkillType, //特效类型
		ActionType: action.Action,
		LoseHp:     action.Damage,
		FinalHp:    defCard.Hp,
		Trigger:    make([]*AtkDetail, 0), //触发的技能
	}
}

func pktAtkAction(action int16, damage int64) *AtkAction {
	return &AtkAction{
		Action: action,
		Damage: damage,
	}
}

/*

	ACTION_FIGHT_BACK //反击
	ACTION_DODGE   //闪避

*/
//普通攻击
func NormalAtk(atkCard *CardInfo, defCard *CardInfo, atkSkill *AtkSkill) (int16, *AtkDetail) {
	damage := atkCard.NormalDamage()
	Action := int16(0)
	if defCard.TriggerDodge() {
		Action = ACTION_DODGE
		damage = 0
	}

	defCard.LoseHp(damage)
	atkDetail := pktAtkDetail(atkCard, defCard, atkSkill, pktAtkAction(Action, damage))

	if defCard.TriggerFightBack() {
		ret, fightBack := FightBackAtk(defCard, atkCard, &AtkSkill{ATK_TYPE_NORMAL_ATK, FIGHT_BACK_ATK})
		if ret == OK{
			atkDetail.Trigger = append(atkDetail.Trigger, fightBack)
		}
	}
	return OK, atkDetail
}

//反击效果
func FightBackAtk(atkCard *CardInfo, defCard *CardInfo, atkSkill *AtkSkill) (int16, *AtkDetail) {
	damage := atkCard.NormalDamage()
	Action := int16(0)
	if defCard.TriggerDodge() {
		Action = ACTION_DODGE
		damage = 0
	}

	defCard.LoseHp(damage)
	atkDetail := pktAtkDetail(atkCard, defCard, atkSkill, pktAtkAction(Action, damage))
	return OK, atkDetail
}

//技能直接攻击
func SkillAttack(atkCard *CardInfo, defCard *CardInfo, atkSkill *AtkSkill) (int16, *AtkDetail) {
	Action := int16(0)
	damage := atkCard.GetSkillDamage(atkSkill.SkillType)
	if defCard.TriggerDodge(){
		Action = ACTION_DODGE
		damage = 0
	}
	defCard.LoseHp(damage)
	atkDetail := pktAtkDetail(atkCard, defCard, atkSkill, pktAtkAction(Action, damage))
	return OK, atkDetail
}



