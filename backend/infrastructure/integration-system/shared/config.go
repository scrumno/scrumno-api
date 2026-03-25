package shared

type ProviderType string

const (
	ProviderIiko    ProviderType = "iiko"
	ProviderSyrve   ProviderType = "syrve"
	ProviderRkeeper ProviderType = "rkeeper"
	ProviderManual  ProviderType = "manual"
)

type IntegrationSystemConfig struct {
	Provider ProviderType
}
