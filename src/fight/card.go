package fight

import(
	"gd_config"
	"fmt"
)

type SkillInfo struct {
	SkillId int16
	Level   int16
}

func (skill *SkillInfo) Trigger() bool {
	return RandomHappen(50)
}

func (skill *SkillInfo) TargetNum() int {
	return Random(3)
}

type CardInfo struct {
	UserId    int64
	Pos       int16
	CardId    int16
	CardLevel int16
	Speed     int16
	Hp        int64
	InitHp    int64
	Skills    map[int16]*SkillInfo //pos对应点的技能
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
	return GetCardBaseAtk(card.CardId, card.CardLevel)
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

func Test() {
	fmt.Println(gd_config.GDData)
}
