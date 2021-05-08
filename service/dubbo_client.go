package service

import (
	"context"
	"time"

	"github.com/mk1010/industry_adaptor/nclink"

	"github.com/apache/dubbo-go/config"
	"google.golang.org/grpc"
)

func dubboInit() error {
	grpcNCLinkServiceClientImpl := new(GrpcNCLinkServiceImpl)
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
		grpcNCLinkServiceClientImpl.NCLinkSendBasicData,
		grpcNCLinkServiceClientImpl.NCLinkGetMeta,
	}
	return nil
}

type GrpcNCLinkServiceImpl struct {
	NCLinkAuth          func(ctx context.Context, in *nclink.NCLinkAuthReq, out *nclink.NCLinkAuthResp) error
	NCLinkSubscribe     func(ctx context.Context, in *nclink.NCLinkTopicSub) (nclink.NCLinkService_NCLinkSubscribeClient, error)
	NCLinkSendData      func(ctx context.Context, in *nclink.NCLinkDataMessage, out *nclink.NCLinkBaseResp) error
	NCLinkSendBasicData func(ctx context.Context, in *nclink.NCLinkTopicMessage, out *nclink.NCLinkBaseResp) error
	NCLinkGetMeta       func(ctx context.Context, in *nclink.NCLinkMetaDataReq, out *nclink.NCLinkMetaDataResp) error
}

func (u *GrpcNCLinkServiceImpl) Reference() string {
	return "GrpcNCLinkServiceClientImpl"
}

func (u *GrpcNCLinkServiceImpl) GetDubboStub(cc *grpc.ClientConn) nclink.NCLinkServiceClient {
	return nclink.NewNCLinkServiceClient(cc)
}
