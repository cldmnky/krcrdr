package base

type (
	Service interface{}

	service struct{}
)

func NewService() Service {
	return &service{}
}
