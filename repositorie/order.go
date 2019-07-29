package repositorie

import (
	"order-backend/common"
	"order-backend/model"
	"order-backend/rpcserver/project"
	"fmt"
	"github.com/pkg/errors"
)

type orderRepository struct {
	orderTableName string
	orderDetailTableName string
	orderIndexTableName string
}

func NewDefaultOrderRepositories(orderId string) *orderRepository {
	suffix := ""
	if orderId != "" {
		suffix = "00"
	}

	return &orderRepository{
		orderTableName: "m_order_" + suffix,
		orderDetailTableName: "m_order_detail_" + suffix,
		orderIndexTableName: "m_order_index",
	}
}

func (o *orderRepository) initTableByOrderId(orderId string) {
	suffix := ""
	if orderId != "" {
		suffix = "00"
	}

	o.orderTableName = "m_order_" + suffix
	o.orderDetailTableName = "m_order_detail_" + suffix
	o.orderIndexTableName = "m_order_index"
}

//通过用户id,平台id,订单号获取订单表名
func (o *orderRepository) getOrderTableByOrderId(order model.Order) string {
	return o.orderTableName
}
func (o *orderRepository) getOrderTableBySellerId(id int) string {
	return o.orderTableName
}

func (o *orderRepository) OrderUpdate(orderModel model.Order) (int64, error) {
	tableName := o.getOrderTableByOrderId(orderModel)
	var whereSql string
	if common.APP.Id > 0 {
		whereSql = fmt.Sprintf("order_id=%s AND platform_id = %d", orderModel.OrderId, common.APP.Id)
	} else {
		whereSql = fmt.Sprintf("order_id=%s ", orderModel.OrderId)
	}
	affected, err := common.DB.Table(tableName).Where(whereSql).MustCols().Update(orderModel)
	if err != nil {
		return 0, nil
	}
	if affected == 0 {
		return 0, nil
	}
	return affected, nil
}

//获取订单数据
func (o *orderRepository) OrderGetOneByOid(orderId string) (model.Order, error) {
	var order model.Order
	o.initTableByOrderId(orderId)
	session := common.DB.Table(o.orderTableName)
	where := map[string]interface{}{
		"order_id" : orderId,
	}
	has, err := GetWhere(session, where).Get(&order)
	if has == false {
		fmt.Println("===", err)
		return order, errors.New("订单不存在")
	}
	return order, nil
}
func (o *orderRepository) OrderSearchByBuyerId(id int, filter string) (int64, error) {
	tableName := o.getOrderTableBySellerId(id)
	query := common.DB.Table(tableName).Where("buyer_id=? AND platform_id= ?", id, common.APP.Id)
	if filter != "" {
		query = query.And(filter)
	}
	return query.Count()
}
func (o *orderRepository) OrderSearchBySellerId(id int, filter string) (int64, error) {
	tableName := o.getOrderTableBySellerId(id)
	query := common.DB.Table(tableName).Where("seller_id=? AND platform_id= ?", id, common.APP.Id)
	if filter != "" {
		query = query.And(filter)
	}
	return query.Count()
}

//分页查询订单
func (o *orderRepository) OrderSearch(filter string, orderby string, pageParam common.PageParam) ([]model.Order, *common.Pagination, error) {
	parkRecordList := make([]model.Order, 0)
	tableName := "m_order_00"
	query := common.DB.Table(tableName).Where("platform_id= ?", common.APP.Id)
	if filter != "" {
		query = query.Where(filter)
	}
	if orderby != "" {
		query = query.OrderBy(orderby)
	}
	p, err := common.FindPaginationData(query, &parkRecordList, pageParam, new(model.Order))
	return parkRecordList, p, err
}

func (o *orderRepository) CreateOrder(
	order model.Order,
	orderDetail []model.OrderDetail,
	orderIndex model.OrderIndex) (err error) {
	//事务开始
	session := common.DB.NewSession()
	session.Begin()
	_, errOrder := session.Table(o.orderTableName).InsertOne(order)
	//----------写入订单索引表结束
	if errOrder != nil {
		fmt.Println("insert mysql order error === ", errOrder)
		return common.NewError(common.ErrorUpdateDB)
	}

	//写入order detail表
	_, errDetail := session.Table(o.orderDetailTableName).Insert(orderDetail)
	if errDetail != nil {
		fmt.Println("insert mysql order detail error === ", errDetail)
		session.Rollback()
		return common.NewError(common.ErrorUpdateDB)
	}
	//写入orderIndex表
	_, errIoi := session.Table(o.orderIndexTableName).InsertOne(orderIndex)
	if errIoi != nil {
		fmt.Println("insert mysql order index error === ", errIoi)
		session.Rollback()
		return common.NewError(common.ErrorUpdateDB)
	}

	//优惠券领奖
	if order.CouponUserCode != "" {
		errCoupon := new(project.RpcCouponService).DoUseCoupon(order.BuyerId, order.CouponUserCode)
		if errCoupon != nil {
			fmt.Println("exchange coupon error === ", errCoupon)
			session.Rollback()
			return errCoupon
		}
	}

	err = session.Commit()
	if err != nil {
		fmt.Println("commit mysql error === ", err)
		return common.NewError(common.ErrorUpdateDB)
	}

	return nil
}

func (o *orderRepository) UpdateOrder(
	where map[string]interface{},
	order model.Order,
	orderColumns []string,
	orderIndex model.OrderIndex,
	orderIndexColumns []string) error {
	//事务开始
	session := common.DB.NewSession()
	session.Begin()

	session.Table(o.orderTableName)
	i, err := GetWhere(session, where).MustCols(orderColumns...).Update(order)
	//----------写入订单索引表结束
	if err != nil || i == 0 {
		fmt.Println("update mysql order error === ", err)
		return common.NewError(common.ErrorUpdateDB)
	}

	//写入orderIndex表
	session.Table(o.orderIndexTableName)
	i, err = GetWhere(session, where).MustCols(orderIndexColumns...).Update(orderIndex)
	if err != nil || i == 0 {
		fmt.Println("update mysql order index error === ", err)
		session.Rollback()
		return common.NewError(common.ErrorUpdateDB)
	}

	err = session.Commit()
	if err != nil {
		fmt.Println("commit mysql error === ", err)
		return common.NewError(common.ErrorUpdateDB)
	}

	return nil
}

