package postgres

import "time"

// Option применяет заданную настройку к репозиторию (Postgres).
type Option func(*Postgres)

// MaxOpenConn задаёт максимальное количество подключений к БД
func MaxOpenConn(size int) Option {
	return func(p *Postgres) {
		p.maxOpenConn = size
	}
}

// MaxIdleConn задаёт максимальное количество бездействующих подключений к БД
func MaxIdleConn(size int) Option {
	return func(p *Postgres) {
		p.maxIdleConn = size
	}
}

// MaxConnIdleTime задаёт время, после которого бездействующее соединение будет закрыто.
func MaxConnIdleTime(duration time.Duration) Option {
	return func(p *Postgres) {
		p.maxConnIdleTime = duration
	}
}

// MaxConnLifeTime задаёт время с момента создания, после которого соединение будет закрыто.
func MaxConnLifeTime(duration time.Duration) Option {
	return func(p *Postgres) {
		p.maxConnLifetime = duration
	}
}

// ConnAttempts задаёт количество попыток подключения.
func ConnAttempts(attempts int) Option {
	return func(p *Postgres) {
		p.connAttempts = attempts
	}
}

// ConnTimeout задаёт таймаут между попытками подключения.
func ConnTimeout(timeout time.Duration) Option {
	return func(p *Postgres) {
		p.connTimeout = timeout
	}
}
