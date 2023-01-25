package password

//go:generate mockgen -source=interface.go -destination=../../server/mocks/mocks_password.go -package=mocks

type IPasswordHash interface {
	Hash(password string) (string, error)
	Check(password, hashedPassword string) error
}
