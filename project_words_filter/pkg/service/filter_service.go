package filter_service_v1

import (
	"context"
	"github.com/sjmshsh/grpc-gin-admin/project_words_filter/pkg/service/ac"
)

var AcMachine *ac.AcAutoMachine

type FilterService struct {
	UnimplementedFilterServiceServer
}

func InitAcMachine() {
	acMachine := ac.New("./dict.txt")
	AcMachine = acMachine
}

func New() *FilterService {
	return &FilterService{}
}

func (FilterService) mustEmbedUnimplementedLoginServiceServer() {}

func (f *FilterService) Filter(ctx context.Context, msg *ContentMessage) (*ContentResponse, error) {
	// 1. 获取参数
	content := msg.Content
	afterContent := AcMachine.Filter(content)
	// 处理业务逻辑
	return &ContentResponse{
		AfterContent: afterContent,
	}, nil
}
