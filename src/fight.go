package main

import (
	"fmt"
	"sort"
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

func (fighter *Fighter) AddCard(userId int64, pos int16, cardId int16, cardLevel int16, speed int16, hp int64) {
	fighter.UserId = userId
	card := &CardInfo{
		UserId:    userId,
		Pos:       pos,
		CardId:    cardId,
		CardLevel: cardLevel,
		Speed:     speed,
		Hp:        hp,
		InitHp:    hp,
		Skills:   make(map[int16]*SkillInfo),
	}
	card.Skills[1] = &SkillInfo{3,1}
	fighter.Cards = append(fighter.Cards, card)
}

func (fighter *Fighter) DebugFighterInfo(){
	for _, card := range fighter.Cards {
		str := fmt.Sprintf("Pos[%v] cardId[%v] level[%v] hp[%v] speed[%v]", card.Pos, card.CardId, card.CardLevel, card.InitHp, card.Speed)
		fmt.Println(str)
	}			
}

func (fight *FightBattle) DebugInitInfo() {
	fmt.Println("atk:")
	fight.AtkInfo.DebugFighterInfo()
	fmt.Println("def:")
	fight.DefInfo.DebugFighterInfo()
}

type AtkDetail struct {
	AtkUserId int64
	DefUserId int64

	AtkCard int16
	DefCard int16
	//AtkType   int16 //攻击类型，普通攻击，技能攻击，反击
	SkillType  int16 //特效类型
	ActionType int16 //是否触发闪避类
	LoseHp     int64
	FinalHp    int64
	Trigger    []*AtkDetail //触发的技能
}

type AtkSkill struct {
	AtkType       int16 //攻击类型，普通攻击，技能攻击，反击
	SkillType     int16 ////特效类型 技能
}

type FightProcess struct {
	Round   int16
	AtkTree []*AtkDetail
}

func (atkProcess *AtkDetail) DebugAtk() {
	str1 := fmt.Sprintf("    atkUser[%v]'card[%v] use skill [%v] atk defUser[%v]'card[%v]", atkProcess.AtkUserId, atkProcess.AtkCard, atkProcess.SkillType,
		atkProcess.DefUserId, atkProcess.DefCard)
	str2 := fmt.Sprintf("    defUser[%v]'card[%v] use action[%v] and lose hp[%v] newhp[%v] and Trigger skill[%v]", atkProcess.DefUserId, atkProcess.DefCard,
		atkProcess.ActionType, atkProcess.LoseHp, atkProcess.FinalHp, len(atkProcess.Trigger))
	fmt.Println(str1)
	fmt.Println(str2)
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

func (fight *FightBattle) GetAtkTarget(atkCard *CardInfo, num int)[]*CardInfo {
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
		ret, atkInfo := fight.DoAtk(atkCard, card, &AtkSkill{AtkType:ATK_TYPE_NORMAL_ATK, SkillType: NORMAL_ATK })
		if ret == OK {
			fight.Process[fight.Round-1].AtkTree = append(fight.Process[fight.Round-1].AtkTree, atkInfo)
		}
	}
}

//卡牌技能攻击
func (fight *FightBattle) SkillAtk(atkCard *CardInfo) {
	for i:= int16(0);i<=3;i++{
		skill, ok := atkCard.Skills[i] 
		if !ok {
			continue
		}
		fight.CardSkillAction(atkCard, skill)
	}
}

func (fight *FightBattle) CardSkillAction (atkCard *CardInfo, skill *SkillInfo){
	if !skill.Trigger() {
		return
	}

	targets := fight.GetAtkTarget(atkCard, skill.TargetNum())
	for _, card := range targets {
		ret, atkInfo := fight.DoAtk(atkCard, card, &AtkSkill{AtkType:ATK_TYPE_SKILL_DIRECT, SkillType: SKILL_DIRECT_ATK })
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

func (fight *FightBattle) DebugProcess() {
	fight.DebugInitInfo()
	fmt.Println("fight total round:", fight.Round)
	for _, process := range fight.Process {
		//fmt.Println("round:", round)
		process.Debug()
	}
}

func (fight *FightBattle) AddCard(atk bool, userId int64, pos int16, cardId int16, cardLevel int16, speed int16, hp int64) {
	var fighter *Fighter
	if atk {
		fighter = fight.AtkInfo
	} else {
		fighter = fight.DefInfo
	}
	fighter.AddCard(userId, pos, cardId, cardLevel, speed, hp)
}

func InitFight() *FightBattle {
	fightBattle := &FightBattle{
		Round:   0,
		AtkInfo: &Fighter{UserId: 1, Cards: make([]*CardInfo, 0)},
		DefInfo: &Fighter{UserId: 2, Cards: make([]*CardInfo, 0)},
		AtkSort: make([]*CardInfo, 0),
		Process: make([]*FightProcess, 0),
	}

	fightBattle.AddCard(true, 10000, 1, 1, 1, 1, 20)
	//fightBattle.AddCard(true, 10000, 2, 2, 1, 3, 100)
	//fightBattle.AddCard(true, 10000, 3, 3, 1, 2, 100)
	fightBattle.AddCard(false, 20000, 1, 1, 1, 2, 10)
	//fightBattle.AddCard(false, 20000, 2, 2, 1, 3, 100)
	//fightBattle.AddCard(false, 20000, 3, 3, 1, 2, 100)
	return fightBattle
}

