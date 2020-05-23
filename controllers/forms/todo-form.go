package forms

// TodoForm structure for JSON input to create todo
type TodoForm struct {
	ID       int    `json:"id"`
	Content  string `json:"content"`
	Time     string `json:"time"`
	Duration int    `json:"duration"`
}

// Valid method return if form is valid and error message
func (tf TodoForm) Valid() (bool, string) {
	if tf.Content == "" {
		return false, "Title is required"
	}

	return true, ""
}

// TimeStart , event started time
func (tf TodoForm) TimeStart() string {
	if tf.Time == "" {
		return "23:59"
	}

	return tf.Time
}
