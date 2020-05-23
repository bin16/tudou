package forms

// NoteForm structure for JSON input to create todo
type NoteForm struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Attachments string `json:"attachments"`
}

// Valid method return if form is valid and error message
func (nf NoteForm) Valid() (bool, string) {
	if nf.Content == "" {
		return false, "Content is required"
	}

	return true, ""
}
