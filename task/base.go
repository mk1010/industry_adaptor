package task

import (
	"context"

	"github.com/apache/dubbo-go/common/logger"
	"github.com/mk1010/industry_adaptor/config"
	"github.com/mk1010/industry_adaptor/nclink"
	"github.com/mk1010/industry_adaptor/service"
)

func Init() {
	ctx := context.Background()
	metaResp := new(nclink.NCLinkMetaDataResp)
	err := service.NCLinkClient.NCLinkGetMeta(ctx, &nclink.NCLinkMetaDataReq{
		AdaptorId: []string{config.ConfInstance.AdaptorID},
	}, metaResp)
	if err != nil || len(metaResp.Adaptors) <= 0 {
		logger.Error("获取适配器元数据失败", err)
		panic(err)
	}
	adaptorMeta := metaResp.Adaptors[0]
	if err := adaptorInit(adaptorMeta); err != nil {
		panic(err)
	}
}
