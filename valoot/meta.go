package valoot

type Meta struct {
	Pagination Pagination    `json:"pagination"`
	Include    []interface{} `json:"include"`
	Custom     []interface{} `json:"custom"`
}
