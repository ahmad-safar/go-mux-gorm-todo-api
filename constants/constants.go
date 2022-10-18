package constants

type AppError struct {
	Err        error
	Message    string
	StatusCode int
}

type DbKey struct{}
