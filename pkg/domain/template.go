package domain

type TemplateI interface {
	// GetName returns the name of the template
	GetName() string
	// GetShortDescription returns the short description of the template
	GetShortDescription() string

	// Setup is the code to run when the template is selected
	Setup(path string) error
}
