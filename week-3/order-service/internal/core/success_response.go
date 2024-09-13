package core

type success struct {
	Data interface{} `json:"data"`
}

func ResponseData(data interface{}) *success {
	return &success{Data: data}
}
