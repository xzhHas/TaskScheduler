package conf

var Routes map[string]Route

type Route struct {
	Host         string `json:"host"`
	Scheme       string `json:"scheme"`
	Uri          string `json:"uri"`
	Auth         string `json:"auth"`
	LimitTimeout int    `json:"limit_timeout"`
	LimitRate    int    `json:"limit_rate"`
	RetryTime    int    `json:"retry_time"`
	Remarks      string `json:"remarks"`
}
