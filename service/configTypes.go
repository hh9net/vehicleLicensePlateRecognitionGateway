package service

import "encoding/xml"

//车牌抓取结果
type PlateInfo struct {
	XMLName        xml.Name `xml:"PlateInfo"`
	Id             string   `xml:"id"`
	PassId         string   `xml:"passId"`
	CamId          string   `xml:"camId"`
	PassTime       string   `xml:"passTime"`
	VehicleImgPath string   `xml:"vehicleImgPath"`
	Channel        string   `xml:"channel"`
	PlateNo        string   `xml:"plateNo"`
	PlateColor     string   `xml:"plateColor"`
}

/*
<PlateInfo>
    <passId>sxjgl_yzjtd_320200_G2_K1071_2_0_001_20201113_000075</passId>
    <camId>sxjgl_yzjtd_320200_G2_K1071_2_0_001</camId>
    <passTime>2020-11-13 17:05:03</passTime>
    <vehicleImgPath>C:\Users\WIN10\Desktop\HSJDEBUG\images\20201113\sxjgl_yzjtd_320200_G2_K1071_2_0_001_20201113_000075.jpg</vehicleImgPath>
    <channel>0</channel>
    <plateNo>苏BE736M</plateNo>
    <plateColor>蓝</plateColor>
</PlateInfo>
*/

//一对一启动相机进程的配置xml
type OneToOneConfig struct {
	XMLName     xml.Name                  `xml:"config"`
	DevCompId   string                    `xml:"devCompId"`
	Devlist     OneToOneConfigDevlist     `xml:"devlist"`
	Channellist OneToOneConfigChannellist `xml:"channellist"`
}

type OneToOneConfigDevlist struct {
	XMLName xml.Name          `xml:"devlist"`
	Dev     OneToOneConfigDev `xml:"dev"`
}
type OneToOneConfigDev struct {
	XMLName  xml.Name `xml:"dev"`
	DevIp    string   `xml:"devIp"`
	Port     string   `xml:"port"`
	UserName string   `xml:"userName"`
	Password string   `xml:"password"`
	Id       string   `xml:"id"`
}

type OneToOneConfigChannellist struct {
	XMLName xml.Name              `xml:"channellist"`
	Channel OneToOneConfigChannel `xml:"channel"`
}
type OneToOneConfigChannel struct {
	XMLName xml.Name `xml:"channel"`
	Id      string   `xml:"id"`
	Index   string   `xml:"index"`
}

/*

<config>
  <devCompId>HIK</devCompId>
  <devlist>
    <dev>
      <devIp>10.25.51.21</devIp>
      <port>8000</port>
      <userName>admin</userName>
      <password>Wx@12345+</password>
      <id/>
    </dev>
  </devlist>
  <channellist>
    <channel>
      <id>sxjgl_yzjtd_320200_G2_K1071_2_0_001</id>
      <index>0</index>
    </channel>
  </channellist>
</config>
*/

//一对多启动相机进程的配置xml  海康    一个海康存储设备 挂多个相机   一个进程  一个文件   多个"通道"   相机id+通道号 可以关联车道号
type OneToMoreConfig struct {
	XMLName     xml.Name                   `xml:"config"`
	DevCompId   string                     `xml:"devCompId"`
	Devlist     OneToMoreConfigDevlist     `xml:"devlist"`
	Channellist OneToMoreConfigChannellist `xml:"channellist"`
}

type OneToMoreConfigDevlist struct {
	XMLName xml.Name           `xml:"devlist"`
	Dev     OneToMoreConfigDev `xml:"dev"`
}
type OneToMoreConfigDev struct {
	XMLName  xml.Name `xml:"dev"`
	DevIp    string   `xml:"devIp"`
	Port     string   `xml:"port"`
	UserName string   `xml:"userName"`
	Password string   `xml:"password"`
	Id       string   `xml:"id"`
}

type OneToMoreConfigChannellist struct {
	XMLName xml.Name                 `xml:"channellist"`
	Channel []OneToMoreConfigChannel `xml:"channel"`
}
type OneToMoreConfigChannel struct {
	XMLName xml.Name `xml:"channel"`
	Id      string   `xml:"id"`
	Index   string   `xml:"index"`
}

//多对多启动相机进程的配置xml
type MoreToMoreConfig struct {
	XMLName     xml.Name                    `xml:"config"`
	DevCompId   string                      `xml:"devCompId"`
	Devlist     MoreToMoreConfigDevlist     `xml:"devlist"`
	Channellist MoreToMoreConfigChannellist `xml:"channellist"`
}

type MoreToMoreConfigDevlist struct {
	XMLName xml.Name              `xml:"devlist"`
	Dev     []MoreToMoreConfigDev `xml:"dev"`
}
type MoreToMoreConfigDev struct {
	XMLName  xml.Name `xml:"dev"`
	DevIp    string   `xml:"devIp"`
	Port     string   `xml:"port"`
	UserName string   `xml:"userName"`
	Password string   `xml:"password"`
	Id       string   `xml:"id"`
}

type MoreToMoreConfigChannellist struct {
	XMLName xml.Name                  `xml:"channellist"`
	Channel []MoreToMoreConfigChannel `xml:"channel"`
}
type MoreToMoreConfigChannel struct {
	XMLName xml.Name `xml:"channel"`
	Id      string   `xml:"id"`
	Index   string   `xml:"index"`
}
