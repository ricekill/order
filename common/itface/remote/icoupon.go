package remote

type RpcCouponService struct {
	Coupon_checkIsOnCoupon  func(int, string, float64, int) string
	Coupon_returnCouponCode func(int, string) string
	Coupon_doUseCoupon      func(string, int) string
}

type CouponService interface {
	CheckIsOnCoupon(int, string, float32, int) (string, string, string, error)
	ReturnCouponCode(int, string) error
	DoUseCoupon(int, string) error
}
