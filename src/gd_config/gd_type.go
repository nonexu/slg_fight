package gd_config

type GDCardConfig struct {
	Id       int16
	CardId   int16
	Name     string
	Level    int16
	AtkLower int64
	AtkUpper int64
	AtkDis   int16
	Speed    int16
	SkillId  int16
	Hp       int64
}

type GDSkillConfig struct {
	Id          int16
	SkillId     int16
	Name        string
	AtkType     int16  //攻击类型
	Action      int16  //效果
	Level       int16
	Pro         int   //触发概率
	DamageLower int64 //技能伤害
	DamageUpper int64 //技能伤害
	AtkDis      int16
	TargetNum   int   //攻击数量
}

type GDActionConfig struct {
	Id int16
	Name string
}
