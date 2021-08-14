package instruction

type InputDataContainer struct {
	Data DataBody `json:"data"`
}

type DataBody struct {
	Id         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	Message string `json:"message"`
	Width   int16  `json:"width"`
	Height  int16  `json:"height"`
}
