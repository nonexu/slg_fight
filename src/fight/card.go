package fight

import (
	"gd_config"
	//"fmt"
)

type TotalInfo struct {
	Type int16
	Num  int64
}

type StatusInfo struct {
	AtkUserId int64
	CardId    int16
	SkillId   int16
	Status    int16
	Damage    int64
}

type CardInfo struct {
	UserId    int64
	Pos       int16
	CardId    int16
	CardLevel int16
	Speed     int16
	Hp        int64
	InitHp    int64
	Skills    map[int16]*SkillInfo  //pos对应点的技能
	TotalData map[int16]*TotalInfo  //统计数据，比如杀敌数
	Status    map[int16]*StatusInfo //状态记录
}

type SortCards []*CardInfo

func (p SortCards) Len() int           { return len(p) }
func (p SortCards) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p SortCards) Less(i, j int) bool { return p[i].Speed > p[j].Speed }

func (card *CardInfo) Dead() bool {
	if card.Hp <= 0 {
		return true
	}
	return false
}

//卡牌伤害
func (card *CardInfo) NormalDamage() int64 {
	cardCfg := gd_config.GetCardLevelCfg(card.CardId, card.CardLevel)
	if cardCfg == nil {
		return 0
	}
	return RandomBetween2Num(cardCfg.AtkLower, cardCfg.AtkUpper)
}

func (card *CardInfo) LoseHp(loseHp int64) int64 {
	if card.Hp >= loseHp {
		card.Hp = card.Hp - loseHp
	} else {
		loseHp = card.Hp
		card.Hp = 0
	}
	return loseHp
}

func (card *CardInfo) TriggerMiss() bool {
	return RandomHappen(1)
}

func (card *CardInfo) TriggerFightBack() bool {
	if card.Dead() {
		return false
	}

	return RandomHappen(1)
}

func (card *CardInfo) GetSkillDamage(skillId int16) int64 {
	skillCfg := gd_config.GetSkillCfg(skillId, 1)
	if skillCfg == nil {
		return 0
	}
	return RandomBetween2Num(skillCfg.DamageLower, skillCfg.DamageUpper)
}
func (card *CardInfo) GetSkillStatus(skillId int16) int16 {
	skillCfg := gd_config.GetSkillCfg(skillId, 1)
	if skillCfg == nil {
		return 0
	}
	return skillCfg.Action
}

func (card *CardInfo) AddTotalData(typ int16, num int64) {
	info, ok := card.TotalData[typ] //map[int16]*TotalInfo
	if !ok {
		info = &TotalInfo{typ, 0}
		card.TotalData[typ] = info
	}
	info.Num += num
}

//添加卡牌状态，后期添加状态逻辑
func (card *CardInfo) AddStatus(info *StatusInfo) bool {
	card.Status[info.Status] = info
	return true
}

func (card *CardInfo) NormalAtkDis() int16 {
	cardCfg := gd_config.GetCardLevelCfg(card.CardId, card.CardLevel)
	if cardCfg == nil {
		return 0
	}
	return cardCfg.AtkDis
}

func (card *CardInfo) SkillTrigger(skillId int16) bool {
	skillCfg := gd_config.GetSkillCfg(skillId, 1)
	if skillCfg == nil {
		return false
	}
	return RandomHappen(skillCfg.Pro)
}

func (card *CardInfo) SkillAtkDis(skillId int16) int16 {
	skillCfg := gd_config.GetSkillCfg(skillId, 1)
	if skillCfg == nil {
		return 0
	}
	return skillCfg.AtkDis
}

func (card *CardInfo) SkillTargetNum(skillId int16) int {
	skillCfg := gd_config.GetSkillCfg(skillId, 1)
	if skillCfg == nil {
		return 0
	}
	return skillCfg.TargetNum
}

func (card *CardInfo) GetDataNum(typ int16) int64 {
	data, ok := card.TotalData[typ]
	if !ok {
		return 0
	}
	return data.Num
}
