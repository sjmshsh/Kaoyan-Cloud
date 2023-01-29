package common

const (
	// FileMd5Key 保存文件所在的路径 eg:FILE_MD5:468s4df6s4a
	FileMd5Key = "FILE_MD5:"
	// FileUploadStatus 保存上传文件的状态
	FileUploadStatus = "FILE_UPLOAD_STATUS"
	// FileProcessingStatus 文件上传过程，位图，判断哪些文件已经上传过了，哪些没有上传过
	FileProcessingStatus = "FILE_PROCESSING_STATUS:"

	// UserCheckIn 签到的key
	UserCheckIn = "usercheckin:"

	WebUv = "uv"

	UserLocation = "locations:"

	EXIST = "existence:"

	// ATTENTION 已经关注
	ATTENTION = "attention"

	// LIKE 已经点赞
	LIKE = "like"

	StarType = 1

	AttentionType = 4

	RepostType = 2

	CommentType = 3

	// AttentionList 用户关注列表
	AttentionList = "attention:"

	// FollowerList 用户粉丝列表
	FollowerList = "follower:"

	// BlogComment 博客的评论
	BlogComment = 1

	Attention = 1

	Follower = 2

	Like = 3

	// User 用户的缓存信息, 暂时有用户姓名，但是后期可以加上性别，头像，年龄等等
	User = "user:"

	// Feed 后面填追随者的id，为追随者构建feed流
	Feed = "feed:"

	// Content 博客内容
	Content = "content:"

	LockListKey = "locklistkey:"
)
