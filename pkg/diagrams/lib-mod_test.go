package diagrams

import (
	"os"
	"path/filepath"
	"testing"

	d "github.com/blushft/go-diagrams/diagram"
	"github.com/blushft/go-diagrams/nodes/azure"
	"github.com/stretchr/testify/require"
)

func TestSavePath(t *testing.T) {
	baseDir := "testdir"
	fileName := "testfile"
	// Render method will always create the dir and fail if exists
	defer os.RemoveAll(baseDir)
	completePath := filepath.Join(baseDir, fileName+".dot")
	require.NoFileExists(t, completePath)
	// Create a dummy diagram to test that all resources are saved in the path
	diag, err := d.New(SavePath(baseDir, fileName))
	f1 := azure.Compute.Vm(d.NodeLabel("vm1"))
	f2 := azure.Network.NetworkInterfaces(d.NodeLabel("nic1"))
	diag.Connect(f1, f2)
	require.NoError(t, err)
	require.NoError(t, diag.Render())
	require.FileExists(t, completePath)
	require.DirExists(t, filepath.Join(baseDir, "assets"))
}
