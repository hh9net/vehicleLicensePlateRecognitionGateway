package gatewayWeb

import "vehicleLicensePlateRecognitionGateway/service"

//查询网关基本数据
func QueryGatewayBasicData() (int, error, *GatewayBasicDataResp) {
	//获取数据
	GatewayBasicData := new(GatewayBasicDataResp)
	GatewayBasicData.Gatewayid = service.Deviceid      //1、网关id
	GatewayBasicData.Version = service.MainVersion     //2、版本号
	GatewayBasicData.StartTime = service.MainStartTime //3、启动时间
	GatewayBasicData.CameraCnt = service.CameraCount   //4、摄像头数量
	switch Gatewaylocation {
	case "1":
		GatewayBasicData.GatewayType = "门架"
	case "2":
		GatewayBasicData.GatewayType = "服务区"
	case "3":
		GatewayBasicData.GatewayType = "收费站"

	}
	//GatewayBasicData.GatewayType = Gatewaylocation + "|" //5、网关类型  1门架、2、服务区[默认] 3、收费站 +站点
	//返回数据
	return StatusSuccessfully, nil, GatewayBasicData
}

//查询网关动态数据
func QueryGatewayDynamicData() (int, error, *GatewayDynamicDataResp) {

	//获取数据
	GatewayDynamicData := new(GatewayDynamicDataResp)
	GatewayDynamicData.IMgCnt = service.CapCnt               //1、网关启动后共抓拍照片数量
	GatewayDynamicData.CameraNormalCnt = service.CameraCount //2、正常摄像头数量
	GatewayDynamicData.CameraAbnormalCnt = 0                 //3、异常摄像头数量
	GatewayDynamicData.OfflineCameraCnt = 0                  //4、离线摄像头数量

	//返回数据
	return StatusSuccessfully, nil, GatewayDynamicData
}

//查询摄像头基本信息列表查询
func QueryCameraInfoData() (int, error, *[]CameraInfo) {
	Resp := make([]CameraInfo, service.CameraCount)
	for id, cmeraid := range service.Cmeraid {
		i := 0
		cameraInfo := new(CameraInfo)
		cameraInfo.CameraId = cmeraid                                  //1、摄像头id
		cameraInfo.BrandName = service.EngineId[id]                    //2、品牌名称
		cameraInfo.CameraIMGCnt = service.CmeraCapCnt[id]              //3、抓拍统计图片的数量
		cameraInfo.LatestSnapshotTime = service.LatestSnapshotTime[id] //4、最近抓拍时间

		cameraInfo.MainrestartCnt = 0 //5、进程重启次数
		l := ""
		switch Gatewaylocation {
		case "1":
			l = "门架"
		case "2":
			l = "服务区"
		case "3":
			l = "收费站"
		}
		cameraInfo.Location = l + service.Name[id] //6、所在位置 1门架、2、服务区 3、收费站
		Resp[i] = *cameraInfo
		i = i + 1
	}

	//返回数据
	return StatusSuccessfully, nil, &Resp
}
