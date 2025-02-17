package protocol

type TriggerSrc = int64

const (
	TriggerSrcProduct TriggerSrc = 1
	TriggerSrcDevice  TriggerSrc = 2
)

type TriggerTimer = int64

const (
	TriggerTimerBefore TriggerTimer = 1
	TriggerTimerAfter  TriggerTimer = 2
)

type TriggerDir = int64

const (
	TriggerDirUp   TriggerDir = 1
	TriggerDirDown TriggerDir = 2
)

const All = "__all__"
