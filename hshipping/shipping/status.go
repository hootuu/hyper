package shipping

// 物流状态常量定义（兼容快递100/京东等常见状态）
const (
	_               Status = iota
	StatusCreated          // 物流单已创建（待揽收）
	StatusPickedUp         // 已揽收（货物已被快递员取走）
	StatusInTransit        // 运输中
	StatusOutFor           // 出库/开始配送（最后一公里）
	StatusDelivered        // 已送达（签收成功）
	StatusFailed           // 配送失败（如：拒收、无法联系到收件人等）
	StatusException        // 运输异常（如：包裹丢失、损坏等）
)
