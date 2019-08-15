package main

//skill
const (
	NORMAL_ATK = iota + 1 //普通攻击
	FIGHT_BACK_ATK
	SKILL_ATK1
)

//action效果
const (
	ACTION = iota
	ACTION_FIGHT_BACK //反击
	ACTION_DODGE   //闪避
)


//error code
const (
	OK = iota + 1
	ERR_INVALID_SKILL
)



