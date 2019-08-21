package fight

import (
	"fmt"
	"gd_config"
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

func (skill *SkillInfo) AtkDis() int16 {
	fmt.Println("skillid", skill.SkillId)
	skcfg := gd_config.GetSkillCfg(skill.SkillId, skill.Level)
	if skcfg == nil {
		return 0
	}
	return skcfg.AtkDis
}
