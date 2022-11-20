package diagrams

// DiagramMetadata represents metadata of a diagram output file.
type DiagramMetadata struct {
	// Name of the diagram, used to register
	Name string
	// Path on where to output the file
	SaveDir string
	// Format for the output file. SUpported are dot and png
	FileFormat string
}

// AzDiagram represents a generic diagram of azure resources
// The specific implementations determine the types of diagrams we support
type AzDiagram interface {
	// Description of what the diagram shows
	GetDescription() string
	// Return the diagram's type
	GetType() string
	// Create and Render the diagram or return an error
	Render(*DiagramMetadata) error
}
