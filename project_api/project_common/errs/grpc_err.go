package errs

import (
	"github.com/sjmshsh/grpc-gin-admin/project_common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcError(err *BError) error {
	return status.Error(codes.Code(err.Code), err.Msg)
}

func ParseGrpcError(err error) (project_common.BusinessCode, string) {
	fromError, _ := status.FromError(err)
	return project_common.BusinessCode(fromError.Code()), fromError.Message()
}
