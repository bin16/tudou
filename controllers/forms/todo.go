package forms

// TodoForm structure for JSON input to create todo
type TodoForm struct {
	ID int `json:"id"`
	// Title       string `json:"title"`
	// Description string `json:"description"`
	Content string `json:"content"`
}

// Valid method return if form is valid and error message
func (tf TodoForm) Valid() (bool, string) {
	if tf.Content == "" {
		return false, "Title is required"
	}

	return true, ""
}
