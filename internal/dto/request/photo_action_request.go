package request

// PhotoActionRequest is reserved for future action payload extension.
// Current behavior APIs do not require a request body.
type PhotoActionRequest struct {
	Source string `json:"source" validate:"omitempty,max=60"`
}
