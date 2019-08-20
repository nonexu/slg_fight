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
	Type int16
	Num  int64
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

func (card *CardInfo) LoseHp(loseHp int64) {
	card.Hp = card.Hp - loseHp
	if card.Hp < 0 {
		card.Hp = 0
	}
}

func (card *CardInfo) TriggerDodge() bool {
	return RandomHappen(50)
}

func (card *CardInfo) TriggerFightBack() bool {
	if card.Dead() {
		return false
	}

	return RandomHappen(50)
}

func (card *CardInfo) GetSkillDamage(skillId int16) int64 {
	return int64(Random(10))
}

func (card *CardInfo) AddTotalData(typ int16, num int64) {
	info, ok := card.TotalData[typ] //map[int16]*TotalInfo
	if !ok {
		info := &TotalInfo{typ, 0}
		card.TotalData[typ] = info
	}
	info.Num += num
}

//添加卡牌状态，后期添加状态逻辑
func (card *CardInfo) AddStatus(info *StatusInfo) bool {
	card.Status[info.Type] = info
	return true
}

func (card *CardInfo) NormalAtkDis() int16 {
	cardCfg := gd_config.GetCardLevelCfg(card.CardId, card.CardLevel)
	if cardCfg == nil {
		return 0
	}
	return cardCfg.AtkDis
}
