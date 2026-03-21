package command_status

type GetStatusRequest struct {
	OrganizationID string `json:"organizationId"`
	CorrelationID  string `json:"correlationId"`
}

type StatusResponse struct {
	State       string `json:"state"`
	ErrorReason string `json:"errorReason,omitempty"`
}
