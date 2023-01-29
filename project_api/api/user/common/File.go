package common

const (
	// FilePath 文件路径
	FilePath = "D:\\图床"

	// ChunkSize 分片大小, 这个是前端和后端安排好的
	ChunkSize = 1024 * 1024 * 50 // 50MB

	// DownloadChunks 分片下载的个数，这里限制为10个，不能太多了
	DownloadChunks = 10
)
