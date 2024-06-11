package data

// MsgType
const (
	Login       = "Login"
	Snapshot    = "Snapshot"
	Trade       = "Trade"
	Entrust     = "Entrust"
	SnapshotOpt = "Snapshot_Opt" //期权快照
	SnapshotHk  = "Snapshot_Hk"  //港股快照
	SnapshotCtp = "Snapshot_Ctp" //期货快照
	TimeLine    = "TimeLine"     // 分时
	Bargain     = "Bargain"      // 分笔明细
)

type Conf struct {
	User     string
	Password string
	Hiscenter
	Rtcenter
}

type Hiscenter struct {
	Addr string
}

type Rtcenter struct {
	Level1
	Level2
}

type Level1 struct {
	Addr string
}

type Level2 struct {
	Addr string
}
