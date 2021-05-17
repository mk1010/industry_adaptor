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
	NCLinkClient = NcLinkService{
		grpcNCLinkServiceClientImpl.NCLinkAuth,
		grpcNCLinkServiceClientImpl.NCLinkSubscribe,
		grpcNCLinkServiceClientImpl.NCLinkSendData,
		grpcNCLinkServiceClientImpl.NCLinkSendBasicData,
		grpcNCLinkServiceClientImpl.NCLinkGetMeta,
	}
	return nil
}

type GrpcNCLinkServiceImpl struct {
	NCLinkAuth          func(ctx context.Context, in *nclink.NCLinkAuthReq, out *nclink.NCLinkAuthResp) error
	NCLinkSubscribe     func(ctx context.Context) (nclink.NCLinkService_NCLinkSubscribeClient, error)
	NCLinkSendData      func(ctx context.Context, in *nclink.NCLinkDataMessage, out *nclink.NCLinkBaseResp) error
	NCLinkSendBasicData func(ctx context.Context, in *nclink.NCLinkTopicMessage, out *nclink.NCLinkBaseResp) error
	NCLinkGetMeta       func(ctx context.Context, in *nclink.NCLinkMetaDataReq, out *nclink.NCLinkMetaDataResp) error
}

func (u *GrpcNCLinkServiceImpl) Reference() string {
	return "nCLinkServiceImpl"
}

func (u *GrpcNCLinkServiceImpl) GetDubboStub(cc *grpc.ClientConn) nclink.NCLinkServiceClient {
	return nclink.NewNCLinkServiceClient(cc)
}
