package request

type PhotoActionRequest struct {
	Source string `json:"source" validate:"omitempty,max=60"`
}
