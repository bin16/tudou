package forms

// TodoForm structure for JSON input to create todo
type TodoForm struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Valid method return if form is valid and error message
func (tf TodoForm) Valid() (bool, string) {
	if tf.Title == "" {
		return false, "Title is required"
	}

	return true, ""
}
