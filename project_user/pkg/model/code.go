package model

import (
	"github.com/sjmshsh/grpc-gin-admin/project_common/errs"
)

var (
	NoLegalMobile = errs.NewError(2001, "手机号不合法")
)