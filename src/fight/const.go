package fight

//攻击类型
const (
	ATK_TYPE_NORMAL_ATK        = iota + 1 //直接攻击
	ATK_TYPE_SKILL_DIRECT                 //技能直接攻击
	ATK_TYPE_SKILL_DIRECT_CURE            //直接治愈
	ATK_TYPE_SKILL_STATUS_ADD             //状态加成
)

//skill， 主动触发
const (
	NORMAL_ATK        = iota + 1 //普通攻击
	FIGHT_BACK_ATK               //反击
	SKILL_DIRECT_ATK             //野火燎原
	SKILL_DIRECT_CURE            //治疗
)

//action效果
const (
	ACTION            = iota + 1
	ACTION_FIGHT_BACK //反击
	ACTION_MISS       //闪避
	ACTION_NO_TARGET  //没有目标
	ACTION_POISONED   //中毒
)

const (
	KILL_NUM  = iota + 1 //杀敌数
	CURE_NUM             //治疗数
	SKILL_NUM            //技能释放数
)

//error code
const (
	OK = iota + 1
	ERR_INVALID_SKILL
	ERR_CFG_MISS
)

var id2Name = map[int64]string{
	10000: "四叶草",
	20000: "如花",
}
