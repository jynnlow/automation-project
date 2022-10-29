package dto

type AuditLogReq struct {
	Timestamp      int
	Map_id         uint
	Key            string
	Original_value int
	New_value      int
	Action         string
	Is_latest      bool
	User_id        int
}

type RevertReq struct {
	Timestamp int
}
