package task

import (
	"context"
	"fmt"

	"github.com/apache/dubbo-go/common/logger"
	"github.com/mk1010/industry_adaptor/config"
	"github.com/mk1010/industry_adaptor/nclink"
	"github.com/mk1010/industry_adaptor/service"
	"github.com/mk1010/industry_adaptor/task/adaptor"
)

func Init(ctx context.Context) error {
	metaResp := new(nclink.NCLinkMetaDataResp)
	err := service.NCLinkClient.NCLinkGetMeta(ctx, &nclink.NCLinkMetaDataReq{
		AdaptorId: []string{config.ConfInstance.AdaptorID},
	}, metaResp)
	if err != nil || len(metaResp.Adaptors) <= 0 {
		logger.Error("获取适配器元数据失败", err)
		if err == nil {
			err = fmt.Errorf("获取适配器元数据失败")
		}
		return err
	}
	adaptorMeta := metaResp.Adaptors[0]
	if err := adaptor.AdaptorInit(ctx, adaptorMeta); err != nil {
		return err
	}
	return nil
}
