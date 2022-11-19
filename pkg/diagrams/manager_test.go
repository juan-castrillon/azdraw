package diagrams

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDiagramManager(t *testing.T) {
	dm, err := NewDiagramManager()
	require.NoError(t, err)
	require.NotNil(t, dm)
	require.NotNil(t, dm.registeredDiagrams)
}

type registerTestCase struct {
	Name     string
	Diagram  AzDiagram
	Metadata *DiagramMetadata
	Error    bool
}

// DummyDiagram implements the AzDiagram interface just for the tests
type dummyDiagram struct {
	name string
}

func (d *dummyDiagram) GetDescription() string {
	return "dummy diagram"
}

func (d *dummyDiagram) GetType() string {
	return "dummy diagram"
}

func (d *dummyDiagram) Render(m *DiagramMetadata) error {
	d.name = m.Name
	return nil
}

func TestRegister(t *testing.T) {
	diag := &dummyDiagram{}
	badDM := &DiagramManager{}
	dm, err := NewDiagramManager()
	require.NoError(t, err)
	testCases := []registerTestCase{
		{
			Name:     "Complete metadata",
			Diagram:  diag,
			Metadata: &DiagramMetadata{Name: "name", Filename: "dir/name", FileFormat: "dot"},
			Error:    false,
		},
		{
			Name:     "Missing metadata",
			Diagram:  diag,
			Metadata: nil,
			Error:    true,
		},
		{
			Name:     "Missing name",
			Diagram:  diag,
			Metadata: &DiagramMetadata{Filename: "something", FileFormat: "dot"},
			Error:    true,
		},
		{
			Name:     "Missing filename",
			Diagram:  diag,
			Metadata: &DiagramMetadata{Name: "name", FileFormat: "dot"},
			Error:    false,
		},
		{
			Name:     "Missing fileformat",
			Diagram:  diag,
			Metadata: &DiagramMetadata{Name: "name", Filename: "some/file"},
			Error:    false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			require.Error(t, badDM.Register(tc.Diagram, tc.Metadata))
			err := dm.Register(tc.Diagram, tc.Metadata)
			if tc.Error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				a, ok := dm.registeredDiagrams[tc.Metadata.Name]
				require.True(t, ok)
				require.Equal(t, tc.Diagram, a.diag)
				require.Equal(t, tc.Metadata.Name, a.metadata.Name)
				if tc.Metadata.Filename == "" {
					exp := filepath.Join(".", tc.Metadata.Name)
					require.Equal(t, exp, a.metadata.Filename)
				} else {
					require.Equal(t, tc.Metadata.Filename, a.metadata.Filename)
				}
				if tc.Metadata.FileFormat == "" {
					require.Equal(t, "dot", a.metadata.FileFormat)
				} else {
					require.Equal(t, tc.Metadata.FileFormat, a.metadata.FileFormat)
				}
			}
		})
	}

}

type renderTestCase struct {
	Name        string
	DiagramName string
	Error       bool
}

func TestRender(t *testing.T) {
	diag := &dummyDiagram{}
	badDM := &DiagramManager{}
	dm, err := NewDiagramManager()
	require.NoError(t, err)
	testCases := []renderTestCase{
		{
			Name:        "Registered diagram",
			DiagramName: "dummy",
			Error:       false,
		},
		{
			Name:        "Unregistered diagram",
			DiagramName: "foo",
			Error:       true,
		},
	}
	require.NoError(t, dm.Register(diag, &DiagramMetadata{Name: "dummy"}))
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			require.Error(t, badDM.Render(tc.DiagramName))
			err := dm.Render(tc.DiagramName)
			if tc.Error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				// The dummy implementation changes the name when render is called
				require.Equal(t, "dummy", diag.name)
			}
		})
	}
}

func TestRenderAll(t *testing.T) {
	diag := &dummyDiagram{}
	diag2 := &dummyDiagram{}
	badDM := &DiagramManager{}
	dm, err := NewDiagramManager()
	require.NoError(t, err)
	dm.Register(diag, &DiagramMetadata{Name: "dummy"})
	dm.Register(diag2, &DiagramMetadata{Name: "dummy2"})
	require.Empty(t, diag.name)
	require.Empty(t, diag2.name)
	require.Error(t, badDM.RenderAll())
	require.NoError(t, dm.RenderAll())
	require.Equal(t, "dummy", diag.name)
	require.Equal(t, "dummy2", diag2.name)
}
