package ops

type WorkOrderStatus = int64

const (
	WorkOrderStatusWait     WorkOrderStatus = 1 //待处理
	WorkOrderStatusHandling WorkOrderStatus = 2 //处理中
	WorkOrderStatusFinished WorkOrderStatus = 3 //处理完成
)
