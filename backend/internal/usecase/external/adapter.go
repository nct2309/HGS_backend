package usecase

type ExternalServiceAdapter struct {
	adaptee ExternalService
}

func NewExternalServiceAdapter(adaptee ExternalService) *ExternalServiceAdapter {
	return &ExternalServiceAdapter{
		adaptee: adaptee,
	}
}

func (a *ExternalServiceAdapter) Execute(strategy Strategy, des any) error {
	return a.adaptee.Execute(strategy, des)
}
