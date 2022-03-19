package njson

import (
	"fmt"
	"testing"
)

type BaseOrderCURLGetOrderByOrderNumberResponseEntity struct {
	OrderID         string  `njson:"_data._data.basic.orderId" `
	OrderNumber     string  `njson:"_data._data.basic.orderNum" `
	ArrivalType     string  `njson:"_data._data.basic.arrivalType" `
	AccessBusiness  string  `njson:"_data._data.basic.accessBusiness" `
	AppointmentTime string  `njson:"fullname:_data._data.basic.customData.appointmentTime;format:json;formatPath:_data._data.basic.customData"`
	CityID          string  `njson:"_data._data.basic.cityId"`
	City            string  `njson:"_data._data.basic.city"`
	Latitude        float64 `njson:"_data._data.basic.customData.latitude"`
	Longitude       float64 `njson:"_data._data.basic.customData.longitude"`
	HttpStatus      int     `json:"httpstatus"`
	BusinessCode    string  `njson:"_data._retcode"`
	BusinessMessage string  `njson:"_data._retinfo"`
	*ResponseStatus
}

func (o *BaseOrderCURLGetOrderByOrderNumberResponseEntity) SetCityID(cityId string) {
	o.CityID = fmt.Sprintf("%s00", cityId)
}

type ResponseStatus struct {
	HttpStatus   int
	BusinessCode string
	Message      string
}

func (r *ResponseStatus) GetHttpStatus() int {
	return r.HttpStatus
}
func (r *ResponseStatus) SetHttpStatus(httpstatus int) {
	r.HttpStatus = httpstatus
}

func (r *ResponseStatus) GetBusinessMessage() string {
	return r.Message
}

func (r *ResponseStatus) GetBusinessCode() string {
	return r.BusinessCode
}

func TestOrder(t *testing.T) {

	baseOrder := `{"_data":{"_data":{"basic":{"accessBusiness":"7","arrivalType":"2","businessMode":"3","businessModeName":"C端回收","businessPartnerId":"5","businessPartnerName":"闲鱼","cancelTime":"","charges":"0","checkFinishTime":"","city":"杭州市","cityId":"3301","costPrice":"0","couponId":"","couponList":"{\\\"activity\\\":[],\\\"coupon\\\":[]}\n","couponOperationList":"","creditFlag":"0","currentPrice":"","cushouQuotate":"125800","customData":"{\\\"address\\\":\\\"浙江省杭州市西湖区蒋村街道 蒋村花园如意苑50号楼1单元7楼\\\",\\\"allowanceRedPacketId\\\":\\\"\\\",\\\"appointmentTime\\\":\\\"2022-03-19 17:00:00\\\",\\\"areaId\\\":\\\"3301\\\",\\\"dktype\\\":\\\"\\\",\\\"evaluateId\\\":\\\"22033704793\\\",\\\"latitude\\\":\\\"30.286060\\\",\\\"longitude\\\":\\\"120.082605\\\",\\\"orderName\\\":\\\"潘健\\\",\\\"orderPhone\\\":\\\"18868009021\\\",\\\"shopId\\\":\\\"50462\\\",\\\"tmallIndoorCashFlag\\\":\\\"0\\\",\\\"tmallPayType\\\":\\\"\\\"}","deliveryMode":"上门","deliveryModeId":"3","destPID":"1426","externalOrderNo":"1520041910041682879","getInTime":"","isTest":"0","masterOrderId":"0","merchTags":"default","mergeFlag":"0","mergePayment":"","needPayPrice":"0","operationCouponValue":"0","operationPrice":"","orderBiz":"0","orderId":"16344076","orderKind":"0","orderNum":"220319201000168","orderProperty":"8","orderTime":"2022-03-19 15:41:58","orderType":"普通订单","payOutPrice":"","payOutTime":"","payOutUser":"","payWays":"2","paymentMode":"现付(当面付费)","paymentModeId":"2","phoneNum":"18868009021","quotation":"125800","quotationCouponValue":"0","recycleType":"2","recycleTypeID":"2","remark":"","remarkStatus":"0","sendTime":"","srcPID":"2010","stateMachineNo":"5","status":"工程师已接单","statusId":"210","supplyPartnerId":"2","supplyPartnerName":"小站(自营)","totalAddPrice":"0","transPrice":"","userId":"655826918","userOpinion":"0","videoState":"0","vivoWelfareCode":"","wechatPayOpenId":"2088422815025892","wechatPayoutType":"0"},"channel":{"businessSourceChannelId":"10000335","channelFlagID":"1","channelId":"10000335","channelName":"闲鱼上门","destChannelId":"10000335","destChannelName":"闲鱼小站","destPID":"1426","srcPID":"2010"},"good":{"IMEI":"","actualWarehouse":"","amount":"1","brandId":"11","brandIdV2":"2","brandName":"苹果","categoryId":"11","channelId":"10000335","classId":"1","className":"手机","consignor":"2","errorModelBackup":"","goodId":"18279008","isAccessory":"0","isPurchaseGood":"0","materialCode":"","newProductId":"","newProductName":"iPhone 7 Plus","plannedWarehouse":"01","priceInputResult":"","productId":"30832","productName":"iPhone 7 Plus","productPic":"30832_20191106172219_783.png","purchaseLabel":"","seriesNumber":"","sourceProductPic":"30832_20191106172219_783.png","useNewProduct":"0","userItem":"[\\\"12\\\",\\\"18\\\",\\\"37\\\",\\\"42\\\",\\\"1689\\\",\\\"83\\\",\\\"10202\\\",\\\"10173\\\",\\\"71\\\",\\\"63\\\",\\\"21\\\",\\\"23\\\",\\\"53\\\",\\\"55\\\",\\\"59\\\",\\\"65\\\",\\\"223\\\",\\\"3246\\\",\\\"2171\\\",\\\"5535\\\",\\\"6931\\\",\\\"7641\\\"]","warehouseId":"6"},"logistics":{"boxCode":"","logisticsCompany":"","logisticsId":"","logisticsType":"","postage":"0","postageType":"","receiverAddress":"","receiverCity":"","receiverCounty":"","receiverName":"","receiverPhone":"","receiverProvince":"","recycleType":"上门回收","recycleTypeID":"2","sender":"","senderAddress":"","senderCity":"","senderCounty":"","senderName":"","senderPhone":"18868009021","senderProvince":"","senderRealName":"","senderRealPhone":"","senderTel":"","trackingNumber":"","visitAddress":"","visitTime":""},"stateTrackInfo":[{"operater":"","operaterId":"655826918","orderStatus":"10","remark":"已下单","statusName":"已下单","updateTime":"2022-03-19 15:41:58"},{"operater":"655826918","operaterId":"655826918","orderStatus":"210","remark":"工程师已接单","statusName":"工程师已接单","updateTime":"2022-03-19 15:41:58"}]},"_ret":"0","_retcode":"0","_retinfo":"success"},"_head":{"_callerServiceId":"110001","_groupNo":"1","_interface":"getOrderInfo","_invokeId":"f1dd7725a35ad4e3ba588e9e6c6e4ddc","_msgType":"response","_remark":"","_timestamps":"1647675720","_version":"0.01"}}`

	actual := BaseOrderCURLGetOrderByOrderNumberResponseEntity{}

	err := Unmarshal([]byte(baseOrder), &actual)
	if err != nil {
		panic(err)
	}
	fmt.Println(actual)
}
