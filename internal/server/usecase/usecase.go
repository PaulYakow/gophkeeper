// Package usecase содержит реализацию логики взаимодействия с сервисами/хранилищем.
// А также интерфейсы для взаимодействия с этим слоем.
package usecase

// Usecase обеспечивает логику взаимодействия с сервисами/хранилищем.
//
//goland:noinspection SpellCheckingInspection
type Usecase struct {
	IAuthorizationService
	IPairsService
}

// New создаёт объект Usecase.
func New(auth IAuthorizationService, pairs IPairsService) (*Usecase, error) {
	return &Usecase{
		auth,
		pairs,
	}, nil
}
