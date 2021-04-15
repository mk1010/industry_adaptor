package service

import (
	"context"
	"github/mk1010/industry_adaptor/nclink"
	"time"

	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	"github.com/apache/dubbo-go/config"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/protocol/dubbo"
	_ "github.com/apache/dubbo-go/protocol/grpc"
	_ "github.com/apache/dubbo-go/registry/protocol"

	// todo etcd?
	_ "github.com/apache/dubbo-go/registry/zookeeper"
	"google.golang.org/grpc"
)

func dubboInit() error {
	grpcNCLinkServiceClientImpl := new(GrpcNCLinkServiceClientImpl)
	config.SetConsumerService(grpcNCLinkServiceClientImpl)
	config.Load()
	time.Sleep(3 * time.Second)
	funcSub := grpcNCLinkServiceClientImpl.NCLinkSubscribe
	NCLinkClient = NcLinkService{
		grpcNCLinkServiceClientImpl.NCLinkAuth,
		func(ctx context.Context, in *nclink.NCLinkTopicSub) (NcLinkSubClient, error) {
			return funcSub(ctx, in)
		},
		grpcNCLinkServiceClientImpl.NCLinkSendData,
		grpcNCLinkServiceClientImpl.NCLinkGetMeta,
	}
	return nil
}

type GrpcNCLinkServiceClientImpl struct {
	NCLinkAuth      func(ctx context.Context, in *nclink.NCLinkAuthReq, out *nclink.NCLinkAuthResp) error
	NCLinkSubscribe func(ctx context.Context, in *nclink.NCLinkTopicSub) (nclink.NCLinkService_NCLinkSubscribeClient, error)
	NCLinkSendData  func(ctx context.Context, in *nclink.NCLinkDataMessage, out *nclink.NCLinkBaseResp) error
	NCLinkGetMeta   func(ctx context.Context, in *nclink.NCLinkMetaDataReq, out *nclink.NCLinkMetaDataResp) error
}

func (u *GrpcNCLinkServiceClientImpl) Reference() string {
	return "GrpcNCLinkServiceClientImpl"
}

func (u *GrpcNCLinkServiceClientImpl) GetDubboStub(cc *grpc.ClientConn) nclink.NCLinkServiceClient {
	return nclink.NewNCLinkServiceClient(cc)
}
