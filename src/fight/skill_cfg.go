package fight

import(
	"gd_config"
)

func GetCardBaseAtk(cardId int16, level int16) int64 {
	cardCfg := gd_config.GetCardLevelCfg(cardId, level)
	if cardCfg == nil {
		return 0
	}
	return RandomBetween2Num(cardCfg.AtkLower , cardCfg.AtkUpper)
}
