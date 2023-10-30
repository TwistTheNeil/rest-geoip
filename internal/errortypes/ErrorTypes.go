package errortypes

type ErrorDatabaseNotFound struct {
	Message    string `json:"message"`
	FileName   string `json:"-"`
	StatusCode int    `json:"-"`
	Err        error  `json:"-"`
}

func (e *ErrorDatabaseNotFound) Error() string {
	return "Database file not found"
}

func (e *ErrorDatabaseNotFound) Unwrap() error {
	return e.Err
}

func NewErrorDatabaseNotFound(err error, filename string) *ErrorDatabaseNotFound {
	return &ErrorDatabaseNotFound{
		FileName:   filename,
		Message:    "We couldn't access the database, please contact the administrator",
		StatusCode: 500,
		Err:        err,
	}
}
