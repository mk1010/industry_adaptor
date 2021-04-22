package service

import (
	"context"
	"sync"

	"github.com/mk1010/industry_adaptor/config"
	"github.com/mk1010/industry_adaptor/nclink"

	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/protocol/dubbo"
	_ "github.com/apache/dubbo-go/protocol/grpc"
	_ "github.com/apache/dubbo-go/registry/protocol"

	// todo etcd?
	_ "github.com/apache/dubbo-go/registry/zookeeper"
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

// impl类型 提供service mesh 屏蔽底层细节
var NCLinkClient NcLinkService

func Init() (err error) {
	serviceInitOnce.Do(func() {
		switch config.ConfInstance.ConnectMethod {
		case "dubbo":
			err = dubboInit()
		case "http":
			// todo
		default:
			panic("Error ConnectMethod")
		}
	})
	return
}
