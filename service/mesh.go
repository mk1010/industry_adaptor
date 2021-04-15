package service

import (
	"context"
	"github/mk1010/industry_adaptor/config"
	"github/mk1010/industry_adaptor/nclink"
	"sync"
)

var serviceInitOnce sync.Once

type NcLinkService struct {
	NCLinkAuth      func(ctx context.Context, in *nclink.NCLinkAuthReq, out *nclink.NCLinkAuthResp) error
	NCLinkSubscribe func(ctx context.Context, in *nclink.NCLinkTopicSub) (NcLinkSubClient, error)
	NCLinkSendData  func(ctx context.Context, in *nclink.NCLinkDataMessage, out *nclink.NCLinkBaseResp) error
	NCLinkGetMeta   func(ctx context.Context, in *nclink.NCLinkMetaDataReq, out *nclink.NCLinkMetaDataResp) error
}

type NcLinkSubClient interface {
	Recv() (*nclink.NCLinkTopicMessage, error)
}

// 接口类型，将dubbo封装成MQTT的方法
var NCLinkClient NcLinkService

func Init() (err error) {
	serviceInitOnce.Do(func() {
		switch config.ConfInstance.ConnectMethod {
		case "DUBBO":
			err = dubboInit()
		case "HTTP":
			// todo
		default:
			panic("Error ConnectMethod")
		}
	})
	return
}
