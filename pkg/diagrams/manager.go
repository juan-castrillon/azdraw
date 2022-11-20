package diagrams

import (
	"errors"
	"fmt"
	"path/filepath"
)

// DiagramManager is the main interaction point with the package
// It allows to register diagrams ,together with metadata
// It is also the interface to render the diagrams
type DiagramManager struct {
	// registeredDiagrams holds the diagrams and their metadata
	registeredDiagrams map[string]*managedDiagram
}

// managedDiagram is an internal structure used to group diagrams with their metadata
type managedDiagram struct {
	// Diagram
	diag AzDiagram
	// Metadata of the diagram
	metadata *DiagramMetadata
}

// NewDiagramManager returns an initialized manager.
// This method must be called to get a manager, as a manager obtained in any other way
// will be uninitialized and therefore raise errors
func NewDiagramManager() (*DiagramManager, error) {
	return &DiagramManager{
		registeredDiagrams: make(map[string]*managedDiagram),
	}, nil
}

// Register a diagram and its metadata
// It returns an error if called on an uninitialized manager
// It returns an error if metadata is nil or if the name field is missing
// For other missing data in the metadata, it sets up the follwing default values:
// - Filename: Name (./$NAME)
// - FileFormat: dot
func (dm *DiagramManager) Register(diag AzDiagram, metadata *DiagramMetadata) error {
	if dm.registeredDiagrams == nil {
		return errors.New("manager has not been initialized")
	}
	if metadata == nil || metadata.Name == "" {
		return errors.New("invalid metadata")
	}
	if metadata.SaveDir == "" {
		metadata.SaveDir = filepath.Join(".", metadata.Name)
	}
	if metadata.FileFormat == "" {
		metadata.FileFormat = "dot"
	}
	d := &managedDiagram{
		diag:     diag,
		metadata: metadata,
	}
	dm.registeredDiagrams[metadata.Name] = d
	return nil
}

// Render the given diagram
// It returns an error if called on an uninitialized manager
// It recieves the name of the diagram as parameter
// It will attempt to look for the diagram based on the name, so it will return error for
// diagrams that have not been registered first.
func (dm *DiagramManager) Render(name string) error {
	if dm.registeredDiagrams == nil {
		return errors.New("manager has not been initialized")
	}
	d, ok := dm.registeredDiagrams[name]
	if !ok {
		return fmt.Errorf("diagram %s was not found", name)
	}
	return d.diag.Render(d.metadata)
}

// RenderAll registered diagrams
func (dm *DiagramManager) RenderAll() error {
	if dm.registeredDiagrams == nil {
		return errors.New("manager has not been initialized")
	}
	for n := range dm.registeredDiagrams {
		err := dm.Render(n)
		if err != nil {
			return err
		}
	}
	return nil
}
