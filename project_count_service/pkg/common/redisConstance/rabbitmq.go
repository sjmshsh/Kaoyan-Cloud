package redisConstance

const (
	// RetryMessageDeleteFailedQueueName 专门用来重试缓存删除失败的队列
	RetryMessageDeleteFailedQueueName = "retryqueue"

	RetryMessageDeleteFailedExchangeName = "retryexchange"

	RetryMessageDeleteFailedRoute = "retry"
)
