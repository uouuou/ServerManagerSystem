package middleware

// ResultBody 默认数据返回结构体
type ResultBody struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// Pages 分页设置
type Pages struct {
	Page        int   `json:"page"`         //当前页码
	PageSize    int   `json:"page_size"`    //每页条数
	TotalAmount int64 `json:"total_amount"` //总条数
}

// ResultPageBody 默认带页码的数据返回结构体
type ResultPageBody struct {
	Code    int         `json:"code"`
	Pages   Pages       `json:"pages"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// ResultTokenBody 默认token的数据返回结构体
type ResultTokenBody struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Token   string      `json:"token"`
	Message string      `json:"message"`
}
