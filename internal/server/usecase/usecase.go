// Package usecase содержит реализацию логики взаимодействия с сервисами/хранилищем.
// А также интерфейсы для взаимодействия с этим слоем.
package usecase

// Usecase обеспечивает логику взаимодействия с сервисами/хранилищем.
//
//goland:noinspection SpellCheckingInspection
type Usecase struct {
	IAuthorizationService
	IPairsService
	IBankService
	ITextService
}

// New создаёт объект Usecase.
func New(auth IAuthorizationService,
	pairs IPairsService,
	cards IBankService,
	notes ITextService,
) (*Usecase, error) {
	return &Usecase{
		auth,
		pairs,
		cards,
		notes,
	}, nil
}
