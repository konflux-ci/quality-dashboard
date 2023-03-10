package v1alpha1

// A User is a GitHub user
type User struct {
	ID      string
	Login   string
	Name    string
	Company string
	Email   string
	URL     string
}
