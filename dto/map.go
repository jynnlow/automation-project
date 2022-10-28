package dto

type MapModel struct {
	Key        string `json:"key"`
	Value      int    `json:"value"`
	Created_at int    `json:"created_at"`
	Updated_at int    `json:"updated_at"`
}

type MapReq struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}
