package task

import (
	"context"
	"fmt"
	"sync"

	"github.com/mk1010/industry_adaptor/nclink"
	"github.com/mk1010/industry_adaptor/service"

	"github.com/apache/dubbo-go/common/logger"
)

func deviceInit()

type TC130 struct {
	DeviceId   string
	DeviceName string
	DeviceMeta *nclink.NCLinkDevice
	done       *chan struct{}
	mutex      sync.Mutex
}

func (t *TC130) Start(ctx context.Context) (err error) {
	if t.DeviceId == "" && (t.DeviceMeta == nil || t.DeviceMeta.DeviceId == "") {
		err = fmt.Errorf("未被标识的设备")
		logger.Error(err)
		return
	}
	if t.DeviceMeta == nil && t.DeviceId != "" {
		metaResp := new(nclink.NCLinkMetaDataResp)
		err = service.NCLinkClient.NCLinkGetMeta(ctx, &nclink.NCLinkMetaDataReq{
			DeviceId: []string{t.DeviceId},
		}, metaResp)
		if err != nil || metaResp.BaseResp == nil || metaResp.BaseResp.StatusCode != nclink.StatusOk {
			logger.Errorf("获取设备%s元数据失败-%v", t.DeviceId, err)
			return
		}
		for _, device := range metaResp.Devices {
			if t.DeviceId == device.DeviceId {
				t.DeviceMeta = device
			}
		}
		if t.DeviceMeta == nil {
			err = fmt.Errorf("未查询到设备元数据")
			logger.Error(err)
			return
		}
	}
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.done != nil && (*t.done) != nil {
		close(*t.done)
	}
	NCLinkDeviceMeta.Store(t.DeviceId, t.DeviceMeta)
	ch := (make(chan struct{}))
	t.done = &ch
	return nil
}
