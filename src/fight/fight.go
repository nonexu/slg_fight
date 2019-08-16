package fight

import (
	"fmt"
	"sort"
	"gd_config"
)

var (
	FIGHT_MAX_ROUND = int16(8)
)

type Fighter struct {
	UserId int64
	Cards  []*CardInfo
}

func (fighter *Fighter) AllDead() bool {
	for _, card := range fighter.Cards {
		if !card.Dead() {
			return false
		}
	}
	return true
}

func (fighter *Fighter) AddCard(userId int64, pos int16, cardId int16, cardLevel int16) {
	fighter.UserId = userId
	card := &CardInfo{
		UserId:    userId,
		Pos:       pos,
		CardId:    cardId,
		CardLevel: cardLevel,
		Skills:    make(map[int16]*SkillInfo),
	}

	cardCfg := gd_config.GetCardLevelCfg(cardId, cardLevel)
	if cardCfg != nil {
		card.Speed = cardCfg.Speed
		card.Hp = cardCfg.Hp
		card.InitHp = card.Hp
	}

	if cardCfg.SkillId != 0 {
		card.Skills[1] = &SkillInfo{3, 1}
	}
	fighter.Cards = append(fighter.Cards, card)
}

func (fighter *Fighter) DebugFighterInfo() {
	for _, card := range fighter.Cards {
		cardCfg := gd_config.GetCardLevelCfg(card.CardId, card.CardLevel)
		if cardCfg == nil {
			continue
		}

		str := fmt.Sprintf("Pos[%v] card[%v] level[%v] hp[%v] speed[%v]", card.Pos, cardCfg.Name, card.CardLevel, cardCfg.Hp, cardCfg.Speed)
		fmt.Println(str)
	}
}

func (fight *FightBattle) DebugInitInfo() {
	fmt.Println("atk:", fight.GetAtkName())
	fight.AtkInfo.DebugFighterInfo()
	fmt.Println("def:", fight.GetDefName())
	fight.DefInfo.DebugFighterInfo()
}

type AtkDetail struct {
	AtkUserId int64
	DefUserId int64

	AtkCard int16
	DefCard int16
	//AtkType   int16 //攻击类型，普通攻击，技能攻击，反击
	SkillType  int16 //特效类型
	ActionType []int16 //是否触发闪避类
	LoseHp     int64
	FinalHp    int64
	Trigger    []*AtkDetail //触发的技能
}

type AtkSkill struct {
	AtkType   int16 //攻击类型，普通攻击，技能攻击，反击
	SkillType int16 ////特效类型 技能
}

type FightProcess struct {
	Round   int16
	AtkTree []*AtkDetail
}

func (atkProcess *AtkDetail) DebugAtk() {
	atkCardCfg := gd_config.GetCardCfg(atkProcess.AtkCard)
	defCardCfg := gd_config.GetCardCfg(atkProcess.DefCard)

	atkName := id2Name[atkProcess.AtkUserId]
	defName := id2Name[atkProcess.DefUserId]

	skillCfg := gd_config.GetSkillCfg(atkProcess.SkillType)
	if atkCardCfg == nil  || defCardCfg == nil || skillCfg == nil{
		return
	} 

	actionStr := ""
	for _, action := range atkProcess.ActionType{
		if actionCfg :=gd_config.GetActionCfg(action); actionCfg != nil {
			if actionStr == ""{
				actionStr +=actionCfg.Name
			}else {
				actionStr = actionStr + "," + actionCfg.Name
			}
		}
	}

	str1 := fmt.Sprintf("    atkUser[%v]'card[%v] 使用  [%v] 攻击 defUser[%v]'card[%v]", atkName, atkCardCfg.Name, skillCfg.Name,
		defName, defCardCfg.Name)
	str2 := fmt.Sprintf("    defUser[%v]'card[%v] 触发 特效[%v] and 失去 hp[%v] newhp[%v]", defName, defCardCfg.Name,
		actionStr, atkProcess.LoseHp, atkProcess.FinalHp)
	//str2 := fmt.Sprintf("    defUser[%v]'card[%v] 触发 特效[%v] and 失去 hp[%v] newhp[%v] and Trigger skill[%v]", defName, defCardCfg.Name,
	//	actionStr, atkProcess.LoseHp, atkProcess.FinalHp, len(atkProcess.Trigger))
	fmt.Println(str1)
	fmt.Println(str2)	

	if atkProcess.FinalHp == 0 {
		str3 := fmt.Sprintf("    defUser[%v]'card[%v] dead", defName, defCardCfg.Name)
		fmt.Println(str3)	
	}

	//触发技能
	for _, triger := range atkProcess.Trigger {
		triger.DebugAtk()
	}
}

func (process *FightProcess) Debug() {
	fmt.Println("this is round:", process.Round)

	for _, atkDetail := range process.AtkTree {
		atkDetail.DebugAtk()
	}
}

type FightBattle struct {
	Round   int16
	AtkInfo *Fighter
	DefInfo *Fighter
	AtkSort []*CardInfo
	Process []*FightProcess
}

func (fight *FightBattle) InitAtkSort() {

	for _, card := range fight.AtkInfo.Cards {
		fight.AtkSort = append(fight.AtkSort, card)
	}

	for _, card := range fight.DefInfo.Cards {
		fight.AtkSort = append(fight.AtkSort, card)
	}
}

func (fight *FightBattle) StartFight() {
	fight.InitAtkSort()
	for i := int16(0); i <= FIGHT_MAX_ROUND; i++ {
		if fight.Done() {
			break
		}
		fight.RoundFight()
	}
}

func (fight *FightBattle) Done() bool {
	if fight.Round >= FIGHT_MAX_ROUND {
		return true
	}

	if fight.AtkInfo.AllDead() || fight.DefInfo.AllDead() {
		return true
	}

	return false
}

func (fight *FightBattle) RoundFight() {
	fight.Round++
	sort.Sort(SortCards(fight.AtkSort))
	fight.Process = append(fight.Process, &FightProcess{
		Round:   fight.Round,
		AtkTree: make([]*AtkDetail, 0),
	})

	for _, atk := range fight.AtkSort {
		if atk.Dead() {
			continue
		}
		//fmt.Println(*atk)
		fight.AtkFight(atk)
	}
}

func (fight *FightBattle) GetAtkTarget(atkCard *CardInfo, num int) []*CardInfo {
	cards := make([]*CardInfo, 0)
	ids := make([]int, 0)

	for pos, card := range fight.AtkSort {
		if card.UserId == atkCard.UserId {
			continue
		}
		if card.Dead() {
			continue
		}
		ids = append(ids, pos)
	}

	ids = RandomIds(ids, num)
	for _, pos := range ids {
		cards = append(cards, fight.AtkSort[pos])
	}
	return cards
}

//开始攻击
func (fight *FightBattle) AtkFight(atkCard *CardInfo) {
	//普通攻击
	fight.NormalAtk(atkCard)
	fight.SkillAtk(atkCard)
}

//普通攻击
func (fight *FightBattle) NormalAtk(atkCard *CardInfo) {
	targets := fight.GetAtkTarget(atkCard, 1)
	for _, card := range targets {
		ret, atkInfo := fight.DoAtk(atkCard, card, &AtkSkill{AtkType: ATK_TYPE_NORMAL_ATK, SkillType: NORMAL_ATK})
		if ret == OK {
			fight.Process[fight.Round-1].AtkTree = append(fight.Process[fight.Round-1].AtkTree, atkInfo)
		}
	}
}

//卡牌技能攻击
func (fight *FightBattle) SkillAtk(atkCard *CardInfo) {
	for i := int16(0); i <= 3; i++ {
		skill, ok := atkCard.Skills[i]
		if !ok {
			continue
		}
		fight.CardSkillAction(atkCard, skill)
	}
}

func (fight *FightBattle) CardSkillAction(atkCard *CardInfo, skill *SkillInfo) {
	if !skill.Trigger() {
		return
	}

	targets := fight.GetAtkTarget(atkCard, skill.TargetNum())
	for _, card := range targets {
		ret, atkInfo := fight.DoAtk(atkCard, card, &AtkSkill{AtkType: ATK_TYPE_SKILL_DIRECT, SkillType: SKILL_DIRECT_ATK})
		if ret == OK {
			fight.Process[fight.Round-1].AtkTree = append(fight.Process[fight.Round-1].AtkTree, atkInfo)
		}
	}
}

func (fight *FightBattle) DoAtk(atkCard *CardInfo, defCard *CardInfo, atkSkill *AtkSkill) (int16, *AtkDetail) {
	return Invoke(atkSkill.AtkType, atkCard, defCard, atkSkill)
}

func (fight *FightBattle) Result() bool {
	if fight.AtkInfo.AllDead() {
		return false
	}

	if !fight.AtkInfo.AllDead() && fight.DefInfo.AllDead() {
		return true
	}
	return false
}

func (fight *FightBattle) GetAtkName() string {
	return id2Name[fight.AtkInfo.UserId]
}

func (fight *FightBattle) GetDefName() string {
	return id2Name[fight.DefInfo.UserId]
}


func (fight *FightBattle) DebugProcess() {
	fight.DebugInitInfo()
	fmt.Println("fight total round:", fight.Round)
	for _, process := range fight.Process {
		//fmt.Println("round:", round)
		process.Debug()
	}
}

func (fight *FightBattle) AddCard(atk bool, userId int64, pos int16, cardId int16, cardLevel int16) {
	var fighter *Fighter
	if atk {
		fighter = fight.AtkInfo
	} else {
		fighter = fight.DefInfo
	}
	fighter.AddCard(userId, pos, cardId, cardLevel)
}

func InitFight() *FightBattle {
	fightBattle := &FightBattle{
		Round:   0,
		AtkInfo: &Fighter{UserId: 1, Cards: make([]*CardInfo, 0)},
		DefInfo: &Fighter{UserId: 2, Cards: make([]*CardInfo, 0)},
		AtkSort: make([]*CardInfo, 0),
		Process: make([]*FightProcess, 0),
	}

	fightBattle.AddCard(true, 10000, 1, 1, 1)
	fightBattle.AddCard(true, 10000, 2, 3, 1)
	//fightBattle.AddCard(true, 10000, 2, 2, 1, 3, 100)
	//fightBattle.AddCard(true, 10000, 3, 3, 1, 2, 100)
	fightBattle.AddCard(false, 20000, 1, 2, 1)
	fightBattle.AddCard(false, 20000, 2, 4, 1)
	//fightBattle.AddCard(false, 20000, 2, 2, 1, 3, 100)
	//fightBattle.AddCard(false, 20000, 3, 3, 1, 2, 100)
	return fightBattle
}