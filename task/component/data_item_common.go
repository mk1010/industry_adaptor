package component

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"

	"github.com/apache/dubbo-go/common/logger"
	"github.com/mk1010/industry_adaptor/nclink"
	"github.com/mk1010/industry_adaptor/task/common"
)

func DataInfoInit(ctx context.Context, adaptorID, deviceID, componentID string, dataInfoMeta *nclink.NCLinkDataInfo) (common.NCLinkDataInfoAPI, error) {
	if dataInfoMeta == nil || dataInfoMeta.DataItem == nil || dataInfoMeta.SampleInfo == nil {
		err := fmt.Errorf("组件元数据为空")
		logger.Error(err)
		return nil, err
	}
	switch dataInfoMeta.DataItem.DataItemType {
	case "nclink_common_data_item":
		{
			dataItem := &NCLinkCommonDataInfo{
				AdaptorID:   adaptorID,
				DeviceID:    deviceID,
				ComponentID: componentID,
				DataInfo:    dataInfoMeta,
				DataChan:    make(chan []byte, 10),
				done:        make(chan struct{}),
			}
			err := dataItem.Start(ctx)
			if err != nil {
				return nil, err
			}
			return dataItem, nil
		}
		// 花里胡哨的还是算了
	/*case "MethodRegisterExample":
	methReg := reflect.ValueOf(NCLinkComponetMethods{})
	callFunc := methReg.MethodByName(dataInfo.DataItem.DataItemType)
	if callFunc.Kind() == reflect.Func {
		msgChan := make(chan *nclink.NCLinkTopicMessage, 3)
		util.GoSafely(func() {
			callFunc.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(dataInfo), reflect.ValueOf(t.done), reflect.ValueOf(msgChan)})
		}, nil)
	} else {
		err = fmt.Errorf("组件数据项%s未找到可执行方法", dataInfo.DataItem.DataItemType)
		logger.Error(err)
		return
	}*/
	default:
		{
			err := fmt.Errorf("未知的数据项类型%s", dataInfoMeta.DataItem.DataItemType)
			logger.Error(err)
			return nil, err
		}
	}
}

var kindDealFuncMap = make(map[nclink.DataKind]func(data *bytes.Buffer, b binary.ByteOrder) (interface{}, error))

func init() {
	kindDealFuncMap[nclink.DataKind_Int32] = getInt32
	kindDealFuncMap[nclink.DataKind_Float64] = getFloat64
	kindDealFuncMap[nclink.DataKind_String] = getString
}

func GetFiledValueLittleEnd(kind nclink.DataKind, data *bytes.Buffer) (interface{}, error) {
	if f, ok := kindDealFuncMap[kind]; ok {
		return f(data, binary.LittleEndian)
	}
	err := fmt.Errorf("未找到执行方法 kind:%s", kind.String())
	logger.Error(err)
	return nil, err
}

func GetFiledValueBigEnd(kind nclink.DataKind, data *bytes.Buffer) (interface{}, error) {
	if f, ok := kindDealFuncMap[kind]; ok {
		return f(data, binary.BigEndian)
	}
	err := fmt.Errorf("未找到执行方法 kind:%s", kind.String())
	logger.Error(err)
	return nil, err
}

func getInt32(data *bytes.Buffer, b binary.ByteOrder) (interface{}, error) {
	var val int32
	err := binary.Read(data, b, &val)
	return val, err
}

func getFloat64(data *bytes.Buffer, b binary.ByteOrder) (interface{}, error) {
	var val float64
	err := binary.Read(data, b, &val)
	return val, err
}

func getString(data *bytes.Buffer, b binary.ByteOrder) (interface{}, error) {
	val, err := data.ReadString('\n')
	return val, err
}
