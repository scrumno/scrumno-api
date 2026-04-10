package shared

type PaymentType string

const (
	PaymentSystemYookassa = "yookassa"
)

type PaymentProviderConfig struct {
	PaymentType PaymentType
}
