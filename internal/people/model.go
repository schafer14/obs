package people

// NewPerson is a person that has not yet been validated and is not ready to persist to the database.
type NewPerson struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

// Person is a person who has been validated and is ready to persist to the database.
type Person struct {
	ID    string `json:"id"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}
