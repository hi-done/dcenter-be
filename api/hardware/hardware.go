// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package hardware

import (
	"context"

	"dcenter/api/hardware/v1"
)

type IHardwareV1 interface {
	GetCpuTemp(ctx context.Context, req *v1.GetCpuTempReq) (res *v1.GetCpuTempRes, err error)
}
