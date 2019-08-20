package fight

//攻击类型
const (
	ATK_TYPE_NORMAL_ATK       = iota + 1 //直接攻击
	ATK_TYPE_SKILL_DIRECT                //技能直接攻击
	ATK_TYPE_SKILL_STATUS_ADD            //状态加成
)

//skill， 主动触发
const (
	NORMAL_ATK       = iota + 1 //普通攻击
	FIGHT_BACK_ATK              //反击
	SKILL_DIRECT_ATK            //技能直接攻击
)

//action效果
const (
	ACTION            = iota + 1
	ACTION_FIGHT_BACK //反击
	ACTION_DODGE      //闪避
)

//error code
const (
	OK = iota + 1
	ERR_INVALID_SKILL
)

var id2Name = map[int64]string{
	10000: "四叶草",
	20000: "如花",
}
