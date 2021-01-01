package postgres

// StorePSQL ...
type StorePSQL struct {
	ConnString   string
	storageLevel int
}

// New - инициализирует Store
func New(connString string) *StorePSQL {
	return &StorePSQL{
		ConnString: connString,
	}
}
