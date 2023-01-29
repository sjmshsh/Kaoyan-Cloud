package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sjmshsh/grpc-gin-admin/project_api/api/user/common"
	pb "github.com/sjmshsh/grpc-gin-admin/project_api/api/user/protoc"

	"github.com/sjmshsh/grpc-gin-admin/project_api/pkg/dao"
	"github.com/sjmshsh/grpc-gin-admin/project_api/pkg/util"
	"github.com/sjmshsh/grpc-gin-admin/project_common"
	"github.com/sjmshsh/grpc-gin-admin/project_common/errs"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type HandlerUser struct {
}

func New() *HandlerUser {
	return &HandlerUser{}
}

type Params struct {
	ctx    *gin.Context
	chunks int64
	chunk  int64
	size   int64
	name   string
	md5    string
}

func (h *HandlerUser) UploadFile(ctx *gin.Context) {
	result := &project_common.Result{}
	file, head, err := ctx.Request.FormFile("file")
	defer file.Close()
	size := head.Size
	name := head.Filename
	// 用户ID从JWT里面解析
	token := ctx.Request.Header.Get("token")
	parseToken, err := util.ParseToken(token)
	if err != nil {
		log.Println(err)
	}
	userId := parseToken.Uid
	// 任务ID使用雪花算法生成
	worker, err := util.NewWorker(0)
	if err != nil {
		log.Println(err)
	}
	// 任务ID
	id := worker.GetId()
	// md5值也是一样的
	md5 := ctx.Query("md5")
	Schunk := ctx.Query("chunk")
	Schunks := ctx.Query("chunks")
	chunk, _ := strconv.ParseInt(Schunk, 10, 64)
	chunks, _ := strconv.ParseInt(Schunks, 10, 64)
	tempDirPath := common.FilePath + md5
	tempFileName := name + "_tmp"
	os.MkdirAll(tempDirPath, os.ModePerm)
	newFile, err := os.Create(tempDirPath + tempFileName)
	if err != nil {
		log.Println(err)
	}
	defer newFile.Close()
	// 用来计算offset的初始位置
	// 0 = 文件开始位置
	// 1 = 当前位置
	// 2 = 文件结尾处
	var whence = 0
	offset := chunk * common.ChunkSize
	path := tempDirPath + tempFileName

	if err != nil {
		log.Println(err)
	}
	newFile.Seek(offset, whence)
	io.Copy(newFile, file)
	pararm := &Params{
		ctx:    ctx,
		chunks: chunks,
		chunk:  chunk,
		size:   size,
		name:   name,
		md5:    md5,
	}
	isOk := checkAndSetUploadProgress(pararm)
	if isOk {
		renameFile(path, tempDirPath+name)
		// 文件重命名之后马上向服务发送rpc请求，录入数据库
		resp, err := UserServiceClient.UploadFile(context.Background(), &pb.UploadFileRequest{
			UserId: int64(userId),
			Id:     id,
			Size:   size,
			Name:   name,
			Md5:    md5,
			Path:   path,
		})
		if err != nil {
			code, msg := errs.ParseGrpcError(err)
			ctx.JSON(http.StatusOK, result.Fail(code, msg))
			return
		}
		ctx.JSON(http.StatusOK, resp)
	}
}

func renameFile(oldPath string, newPath string) {
	os.Rename(oldPath, newPath)
}

func checkAndSetUploadProgress(param *Params) bool {
	// 把该分段标记成true表示完成
	chunk := param.chunk
	chunks := param.chunks
	md5 := param.md5
	key := common.FileProcessingStatus + md5
	dao.Rdb.SetBit(dao.RCtx, key, chunk, 1)
	flag := true
	// 检查是否全部完成
	for i := 0; i < int(chunks); i++ {
		result, err := dao.Rdb.GetBit(dao.RCtx, key, int64(i)).Result()
		if err != nil {
			log.Println(err)
		}
		if result == 0 {
			flag = false
			break
		}
	}
	if flag == true {
		// 说明文件已经上传完毕了
		dao.Rdb.HSet(dao.RCtx, common.FileUploadStatus, md5, "true")
		return true
	} else {
		// 说明文件还没有上传完毕
		// 如果你不存在这个哈希的key就顺便创建，如果存在的话就直接跳过
		n, err := dao.Rdb.Exists(dao.RCtx, common.FileUploadStatus, md5).Result()
		if err != nil {
			log.Println(err)
		}
		if n <= 0 {
			dao.Rdb.HSet(dao.RCtx, common.FileUploadStatus, md5, "false")
		}
		return false
	}
}

// 文件分片
type filePart struct {
	Index int    // 文件分片的序号
	From  int    // 开始的byte
	To    int    // 结束的byte
	Data  []byte // http 下载得到的文件内容
}

// FileDownloader 文件下载器
type FileDownloader struct {
	fileSize       int
	url            string
	outputFileName string
	totalPart      int // 下载线程
	outputDir      string
	doneFilePart   []filePart
}

func NewFileDownloader(url, outputFileName, outputDir string, totalPart int) *FileDownloader {
	return &FileDownloader{
		fileSize:       0,
		url:            url,
		outputFileName: outputFileName,
		totalPart:      totalPart,
		outputDir:      outputDir,
		doneFilePart:   make([]filePart, totalPart),
	}
}

func (h *HandlerUser) DownLoadFile(ctx *gin.Context) {
	outputFileName := ctx.Query("filename")
	outputDir := ctx.Query("dir")
	// 我上传在考研云的文件也需要生成上传和下载的URL
	url := ctx.Query("url")
	downloader := NewFileDownloader(url, outputFileName, outputDir, common.DownloadChunks)
	if err := downloader.Run(); err != nil {
		log.Println(err)
	}
	ctx.JSON(http.StatusOK, "文件下载完成")
}

// 创建一个request
func (d *FileDownloader) getNewRequest(method string) (*http.Request, error) {
	r, err := http.NewRequest(method, d.url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("User-Agent", "lxy")
	return r, nil
}

// 获取要下载的文件的基本信息
func (d *FileDownloader) head() (int, error) {
	r, err := d.getNewRequest("HEAD")
	if err != nil {
		return 0, err
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode > 299 {
		return 0, errors.New(fmt.Sprintf("Can't process, response is %v", resp.StatusCode))
	}
	// 检查是否支持断点续传
	if resp.Header.Get("Accept-Ranges") != "bytes" {
		return 0, errors.New("服务器不支持断点续传")
	}
	return strconv.Atoi(resp.Header.Get("Content-Length"))
}

// Run 开始下载任务
func (d *FileDownloader) Run() error {
	fileTotalSize, err := d.head()
	if err != nil {
		log.Println(err)
	}
	d.fileSize = fileTotalSize

	jobs := make([]filePart, d.totalPart)
	eachSize := fileTotalSize / d.totalPart

	for i := range jobs {
		jobs[i].Index = i
		if i == 0 {
			jobs[i].From = 0
		} else {
			jobs[i].From = jobs[i-1].To + 1
		}
		if i < d.totalPart-1 {
			jobs[i].To = jobs[i].From + eachSize
		} else {
			// the last filePart
			jobs[i].To = fileTotalSize - 1
		}
	}

	var wg sync.WaitGroup
	for _, j := range jobs {
		wg.Add(1)
		go func(job filePart) {
			defer wg.Done()
			err := d.downloadPart(job)
			if err != nil {
				log.Println("下载文件失败: ", err, job)
			}
		}(j)
	}
	wg.Wait()
	return d.mergeFileParts()
}

// 下载分片
func (d *FileDownloader) downloadPart(c filePart) error {
	r, err := d.getNewRequest("GET")
	if err != nil {
		log.Println(err)
	}
	log.Printf("开销下载[%d]下载from:%d to: %d\n", c.Index, c.From, c.To)
	r.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", c.From, c.To))
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode > 299 {
		return errors.New(fmt.Sprintf("服务器错误状态码: %v", resp.StatusCode))
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	if len(bs) != (c.To - c.From + 1) {
		return errors.New("下载文件分片长度错误")
	}
	c.Data = bs
	d.doneFilePart[c.Index] = c
	return nil
}

func (d *FileDownloader) mergeFileParts() error {
	log.Println("开始合并文件")
	path := d.outputDir + d.outputFileName
	mergedFile, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}
	defer mergedFile.Close()
	totalSize := 0
	for _, s := range d.doneFilePart {
		mergedFile.Write(s.Data)
		totalSize += len(s.Data)
	}
	if totalSize != d.fileSize {
		return errors.New("文件不完善")
	}
	return nil
}

func (h *HandlerUser) CheckFileMd5(ctx *gin.Context) {
	result := &project_common.Result{}
	md5 := ctx.Query("md5")
	SChunks := ctx.Query("chunks")
	chunks, err := strconv.ParseInt(SChunks, 10, 32)
	if err != nil {
		log.Println(err)
	}
	resp, err := UserServiceClient.CheckFileMd5(context.Background(), &pb.CheckFileMd5Request{
		Md5:    md5,
		Chunks: int32(chunks),
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *HandlerUser) CheckIn(ctx *gin.Context) {
	result := &project_common.Result{}
	token := ctx.Request.Header.Get("token")
	parseToken, err := util.ParseToken(token)
	if err != nil {
		log.Println(err)
	}
	userId := parseToken.Uid
	// 获取目前的年份和月份还有天
	year := time.Now().Format("2006")
	month := time.Now().Format("1")
	day := time.Now().Format("02")
	// control层把东西发给service层进行业务逻辑开发
	resp, err := UserServiceClient.CheckIn(context.Background(), &pb.CheckSignRequest{
		UserId: int64(userId),
		Year:   year,
		Month:  month,
		Day:    day,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(int(resp.Status), resp.Msg)
}

func (h *HandlerUser) GetSign(ctx *gin.Context) {
	result := &project_common.Result{}
	token := ctx.Request.Header.Get("token")
	parseToken, err := util.ParseToken(token)
	if err != nil {
		log.Println(err)
	}
	userId := parseToken.Uid
	year := ctx.Query("year")
	month := ctx.Query("month")
	// control层把东西发给service层进行业务逻辑开发
	resp, err := UserServiceClient.GetSign(context.Background(), &pb.GetSignRequest{
		UserId: int64(userId),
		Year:   year,
		Month:  month,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(int(resp.Status), resp.Msg)
}

func (h *HandlerUser) WatchUv(ctx *gin.Context) {
	result := &project_common.Result{}
	resp, err := UserServiceClient.WatchUv(context.Background(), &pb.WatchUvRequest{})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, resp.Uv)
}

func (h *HandlerUser) Location(ctx *gin.Context) {
	result := &project_common.Result{}
	longitude := ctx.Query("longitude")
	latitude := ctx.Query("latitude")
	location := ctx.Query("location")
	header := ctx.GetHeader("token")
	token, err := util.ParseToken(header)
	if err != nil {
		log.Println(err)
	}
	userId := token.Uid
	resp, err := UserServiceClient.Location(context.Background(), &pb.LocationRequest{
		Longitude: longitude,
		Latitude:  latitude,
		UserId:    int64(userId),
		Location:  location,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(int(resp.Status), resp.Msg)
}

func (h *HandlerUser) FindFriend(ctx *gin.Context) {
	result := &project_common.Result{}
	longitude := ctx.Query("longitude")
	latitude := ctx.Query("latitude")
	resp, err := UserServiceClient.FindFriend(context.Background(), &pb.FindFriendRequest{
		Longitude: longitude,
		Latitude:  latitude,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(&struct {
		Name []string
		Dist []float32
	}{
		Name: resp.Name,
		Dist: resp.Dist,
	}))
}

func (h *HandlerUser) Watch(ctx *gin.Context) {
	// 获取用户ID
	//header := ctx.GetHeader("token")
	//token, err := util.ParseToken(header)
	//if err != nil {
	//	log.Println(err)
	//}
	//userId := token.Uid
	userId := 1
	id := ctx.Query("id")
	attentionId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err)
	}
	resp, err := UserServiceClient.Watch(context.Background(), &pb.WatchRequest{
		UserId:          int64(userId),
		AttentionUserId: attentionId,
	})
	if err != nil {
		log.Println(err)
	}
	ctx.JSON(http.StatusOK, resp.Msg)
}

func (h *HandlerUser) PostBlog(ctx *gin.Context) {
	// 获取用户ID
	//header := ctx.GetHeader("token")
	//token, err := util.ParseToken(header)
	//if err != nil {
	//	log.Println(err)
	//}
	//userId := token.Uid
	// 获取文章内容
	content := ctx.PostForm("content")
	resp, err := UserServiceClient.PostBlog(context.Background(), &pb.PostBlogRequest{
		Content: content,
		UserId:  1,
	})
	if err != nil {
		log.Println(err)
	}
	ctx.JSON(http.StatusOK, resp.Msg)
}

func (h *HandlerUser) OList(ctx *gin.Context) {
	header := ctx.GetHeader("token")
	token, err := util.ParseToken(header)
	if err != nil {
		log.Println(err)
	}
	id := token.Uid
	t := ctx.Query("type")
	ty, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		log.Println(err)
	}
	start := ctx.Query("start")
	s, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		log.Println(err)
	}
	offset := ctx.Query("offset")
	o, err := strconv.ParseInt(offset, 10, 64)
	if err != nil {
		log.Println(err)
	}
	resp, err := UserServiceClient.OList(context.Background(), &pb.ListRequest{
		UserId: int64(id),
		Start:  int64(int(s)),
		Offset: int64(int(o)),
		Type:   int64(int(ty)),
	})
	ctx.JSON(http.StatusOK, resp.UserList)
}

func (h *HandlerUser) List(ctx *gin.Context) {
	t := ctx.Query("type")
	ty, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		log.Println(err)
	}
	start := ctx.Query("start")
	s, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		log.Println(err)
	}
	offset := ctx.Query("offset")
	o, err := strconv.ParseInt(offset, 10, 64)
	if err != nil {
		log.Println(err)
	}
	resp, err := UserServiceClient.List(context.Background(), &pb.ListRequest{
		UserId: 1, // 这里应该用JWT获取用户ID，为了测试方便我们暂时这么去写
		Start:  int64(int(s)),
		Offset: int64(int(o)),
		Type:   int64(int(ty)),
	})
	ctx.JSON(http.StatusOK, resp.UserList)
}

func (h *HandlerUser) Comment(ctx *gin.Context) {
	value := ctx.PostForm("content")
	id := ctx.PostForm("id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err)
	}
	resp, err := UserServiceClient.Comment(context.Background(), &pb.CommentRequest{
		Content: value,
		Id:      i,
	})
	if err != nil {
		log.Println(err)
	}
	ctx.JSON(http.StatusOK, resp.Msg)
}

func (h *HandlerUser) GetFeedList(ctx *gin.Context) {
	//header := ctx.GetHeader("token")
	//token, err := util.ParseToken(header)
	s := ctx.Query("start")
	o := ctx.Query("offset")
	start, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Println(err)
	}
	offset, err := strconv.ParseInt(o, 10, 64)
	if err != nil {
		log.Println(err)
	}
	// userId := token.Uid
	resp, err := UserServiceClient.GetFeedList(context.Background(), &pb.GetFeedListRequest{
		UserId: 1,
		Start:  start,
		Offset: offset,
	})
	if err != nil {
		log.Println(err)
	}
	ctx.JSON(http.StatusOK, resp.Feed)
}
