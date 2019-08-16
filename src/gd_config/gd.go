package gd_config

import (
	//"fmt"
)

type GDDataInfo struct {
	GD_CARD_CFG  map[int16]*GDCardConfig
	GD_SKILL_CFG map[int16]*GDSkillConfig
	GD_ACTION_CFG map[int16]*GDActionConfig
}

var GDData GDDataInfo


func loadData() {
	loadCardConfig()
	loadSkillConfig()
	loadActionConfig()
	//fmt.Println(GDData)
}

func loadCardConfig() {
	GDData.GD_CARD_CFG = make(map[int16]*GDCardConfig)
	recs := GetAll("card_config")
	for _, rec := range recs {
		fields := rec.Fields
		info := &GDCardConfig{
			Id:       Sto16(fields["Id"]),
			CardId:   Sto16(fields["CardId"]),
			Name:     fields["Name"],
			Level:    Sto16(fields["Level"]),
			AtkLower: Sto64(fields["AtkLower"]),
			AtkUpper: Sto64(fields["AtkUpper"]),
			Speed:    Sto16(fields["Speed"]),
			SkillId:  Sto16(fields["SkillId"]),
			Hp: Sto64(fields["Hp"]),
		}
		GDData.GD_CARD_CFG[info.Id] = info
	}
}

func GetCardLevelCfg(cardId int16, level int16) *GDCardConfig {
	for _, rec := range GDData.GD_CARD_CFG{
		if rec.CardId == cardId && rec.Level == level{
			return rec
		}
	}
	return nil
}

func GetCardCfg(cardId int16) *GDCardConfig {
	for _, rec := range GDData.GD_CARD_CFG{
		if rec.CardId == cardId{
			return rec
		}
	}
	return nil
}

func loadSkillConfig() {
	GDData.GD_SKILL_CFG = make(map[int16]*GDSkillConfig)
	recs := GetAll("skill_config")
	for _, rec := range recs {
		fields := rec.Fields
		info := &GDSkillConfig{
			Id:          Sto16(fields["Id"]),
			SkillId:     Sto16(fields["SkillId"]),
			Name:        fields["Name"],
			AtkType:     Sto16(fields["AtkType"]),
			Action:      Sto16(fields["Action"]),
			Level:       Sto16(fields["Level"]),
			Pro:         Stoi(fields["Pro"]),
			DamageLower: Sto64(fields["DamageLower"]),
			DamageUpper: Sto64(fields["DamageUpper"]),
			TargetNum:   Stoi(fields["TargetNum"]),
		}
		GDData.GD_SKILL_CFG[info.Id] = info
	}
}

func GetSkillCfg( skillId int16)*GDSkillConfig {
	for _, rec := range GDData.GD_SKILL_CFG {
		if rec.SkillId == skillId {
			return rec
		}
	}
	return nil
}

func loadActionConfig() {
	GDData.GD_ACTION_CFG = make(map[int16]*GDActionConfig)
	recs := GetAll("action_config")
	for _, rec := range recs {
		fields := rec.Fields
		info := &GDActionConfig{
			Id:          Sto16(fields["Id"]),
			Name:        fields["Name"],
		}
		GDData.GD_ACTION_CFG[info.Id] = info
	}
}

func GetActionCfg(id int16)*GDActionConfig {
	return GDData.GD_ACTION_CFG[id]
}




