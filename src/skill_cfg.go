package main

type GDCardCfg struct {
	Id int32
	CardId int16
	Level int16
	Attack int64
	SkillId1 int16
	SkillId2 int16
	SkillId3 int16
}

type GDSkillCfg struct {
	Id int32	
	SkillId int16
	Level int16
	TriggerPer int32
	TargetNum int16
}

var GD_CARDS map[int32]*GDCardCfg
var GD_SKILLS map[int32]*GDSkillCfg

func init() {
	GD_CARDS = map[int32]*GDCardCfg{
		1: &GDCardCfg{1,1,1,3,0,0,0},
		2: &GDCardCfg{2,2,1,2,0,0,0},
		3: &GDCardCfg{3,3,1,3,0,0,0},
	}


	GD_SKILLS = map[int32]*GDSkillCfg{
		1: &GDSkillCfg{1,1,1,50,1},
		2: &GDSkillCfg{2,2,1,50,2},
		3: &GDSkillCfg{3,2,1,50,3},
	}
}

func GetCardBaseAtk(cardId int16, level int16) int64 {
	for _,card := range GD_CARDS {
		if card.CardId == cardId && card.Level ==level {
			return card.Attack
		}
	}
	return 0
}





