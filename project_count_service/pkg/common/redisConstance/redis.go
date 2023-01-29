package redisConstance

const (
	// CounterFuzzyService 模糊计数服务
	CounterFuzzyService = "counterfuzzy:"

	// CountPreciseService 精准计数服务
	CountPreciseService = "counterprecise:"

	// CounterServiceTTL TTL设置成一个星期
	CounterServiceTTL = 604800

	StarType = 1

	AttentionType = 4

	RepostType = 2

	CommentType = 3

	WatchType = 5

	// StarNum 点赞数量
	StarNum = "star_num"

	// WatchNum 浏览次数
	WatchNum = "watch_num"

	// Attention 关注数量
	Attention = "attention_num"

	// RepostNum 转发数量
	RepostNum = "repost_num"

	// CommentNum 评论数量
	CommentNum = "comment_num"
	/*
	* 这些各种数量你想的话可以随便添加的，可拓展性非常好
	 */

	// Id 计数的ID
	Id = "id"

	// VUser 大V用户的判断标准
	VUser = 100000

	// ATTENTIONS 粉丝关注列表
	ATTENTIONS = "attentions:"
)
