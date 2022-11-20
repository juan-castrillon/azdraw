package diagrams

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateMaps(t *testing.T) {
	vmi := []VMInfo{
		{
			Name: "vm1",
			Subnet: SubnetInfo{
				Name: "subnet11",
				Vnet: "vnet1",
			},
			NSG:       "nsg1",
			PrivateIP: "10_0_0_1",
			PublicIP:  "120_34_23_54",
		},
		{
			Name: "vm2",
			Subnet: SubnetInfo{
				Name: "subnet22",
				Vnet: "vnet2",
			},
		},
		{
			Name: "vm3",
			Subnet: SubnetInfo{
				Name: "subnet11",
				Vnet: "vnet1",
			},
		},
		{
			Name: "vm4",
			Subnet: SubnetInfo{
				Name: "subnet12",
				Vnet: "vnet2",
			},
		},
		{
			Name: "vm5",
			Subnet: SubnetInfo{
				Name: "subnet21",
				Vnet: "vnet1",
			},
		},
	}
	d := &VMNetworkDiagram{
		VMList: vmi,
	}
	//r := d.generateMaps()
	//oe, err := json.MarshalIndent(r, "", "\t")
	//require.NoError(t, err)
	//t.Log(string(oe))
	//t.Logf("%+v", r)
	require.NoError(t, d.Render(&DiagramMetadata{
		Name:    "test",
		SaveDir: "testjp",
	}))
	defer os.RemoveAll("testjp")
	imagePath := "../../test.png"
	cmd := exec.Command("dot", "-Tpng", "test.dot")
	cmd.Dir = "testjp"
	var b bytes.Buffer
	cmd.Stderr = &b
	bytes, err := cmd.Output()
	require.NoError(t, err, "output: %s", b.String())
	require.NoError(t, os.WriteFile(imagePath, bytes, 0777))
	cmd2 := exec.Command("eog", "../"+imagePath)
	cmd2.Dir = "testjp"
	require.NoError(t, cmd2.Run())
	//defer os.Remove(imagePath)
}
