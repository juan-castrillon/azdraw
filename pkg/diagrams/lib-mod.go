package diagrams

import (
	d "github.com/blushft/go-diagrams/diagram"
)

// SavePath is a custom diagram.Option method to allow for configuration of the complete path
// It is meant to be passed as argument to the diagram::New method
// The standard library only allows to modifiy the filename, however,
// the assets folder is created in the go-diagrams directory all the time
// This custom option changes that and allows to save all resources (diagram and assets) in a desired path
func SavePath(baseDir, fileName string) d.Option {
	return func(o *d.Options) {
		o.Name = baseDir
		o.FileName = fileName
	}
}
