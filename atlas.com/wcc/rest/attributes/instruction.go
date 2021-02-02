package attributes

type InstructionInputDataContainer struct {
	Data InstructionData `json:"data"`
}

type InstructionData struct {
	Id         string                `json:"id"`
	Type       string                `json:"type"`
	Attributes InstructionAttributes `json:"attributes"`
}

type InstructionAttributes struct {
	Message string `json:"message"`
	Width   int16 `json:"width"`
	Height  int16 `json:"height"`
}