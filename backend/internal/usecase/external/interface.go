package usecase

type Strategy interface {
	Execute(des any) error
}

type ExternalService interface {
	Execute(strategy Strategy, des any) error
}
