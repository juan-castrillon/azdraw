package diagrams

import (
	"fmt"

	dg "github.com/blushft/go-diagrams/diagram"
	"github.com/blushft/go-diagrams/nodes/azure"
)

type VMNetworkDiagram struct {
	VMList []VMInfo
}

type SubnetInfo struct {
	Name         string
	Vnet         string
	AddressRange string
	NSG          string
}

type VMInfo struct {
	Name      string
	PrivateIP string
	PublicIP  string
	NSG       string
	Subnet    SubnetInfo
}

// Description of what the diagram shows
func (d *VMNetworkDiagram) GetDescription() string {
	return "Diagram that shows how VMs are placed in the network architecture"
}

// Return the diagram's type
func (d *VMNetworkDiagram) GetType() string {
	return "VM_Network"
}

// Create and Render the diagram or return an error
func (d *VMNetworkDiagram) Render(metadata *DiagramMetadata) error {
	diag, err := d.createDiagram(metadata)
	if err != nil {
		return err
	}
	r := d.generateMaps()
	for vnetName, subnetMap := range r.vnet2subnet {
		g := dg.NewGroup(vnetName).Label(vnetName)
		for subnetName := range subnetMap {
			sg := g.NewGroup(subnetName).Label(subnetName)
			for vmName := range r.subnet2vm[subnetName] {
				vm := r.vm[vmName]
				sg.Group(processVM(&vm))
			}
		}
		diag.Group(g)
	}
	return diag.Render()
}

type OrganizeResult struct {
	vnet2subnet map[string]map[string]bool
	subnet2vm   map[string]map[string]bool
	vm          map[string]VMInfo
	subnet      map[string]SubnetInfo
}

func (d *VMNetworkDiagram) generateMaps() *OrganizeResult {
	vnetSet := make(map[string]bool)
	m1 := make(map[string]map[string]bool)
	m2 := make(map[string]map[string]bool)
	m3 := make(map[string]VMInfo)
	m4 := make(map[string]SubnetInfo)
	for _, vm := range d.VMList {
		vnetName := vm.Subnet.Vnet
		subnetName := vm.Subnet.Name
		if !vnetSet[vnetName] {
			vnetSet[vnetName] = true
			m1[vnetName] = make(map[string]bool)
		}
		if !m1[vnetName][vm.Subnet.Name] {
			m1[vnetName][vm.Subnet.Name] = true
			m4[subnetName] = vm.Subnet
			m2[subnetName] = make(map[string]bool)
		}
		if !m2[subnetName][vm.Name] {
			m2[subnetName][vm.Name] = true
		}
		m3[vm.Name] = vm
	}
	return &OrganizeResult{
		vnet2subnet: m1,
		subnet2vm:   m2,
		vm:          m3,
		subnet:      m4,
	}
}

func processVM(i *VMInfo) *dg.Group {
	label := fmt.Sprintf("%s [%s]", i.Name, i.PrivateIP)
	g := dg.NewGroup(i.Name, dg.GroupLabel(label))
	vm := azure.Compute.Vm(dg.NodeLabel(i.Name), dg.LabelLocation("b"))
	g.Add(vm)
	if i.PublicIP != "" {
		pip := azure.Network.PublicIpAddresses(dg.NodeLabel(i.PublicIP))
		g.Add(pip)
		g.Connect(vm, pip)
	}
	if i.NSG != "" {
		nsg := azure.Network.NetworkSecurityGroupsClassic(dg.NodeLabel(i.NSG))
		g.Add(nsg)
		g.Connect(vm, nsg)
	}
	return g
}

func (d *VMNetworkDiagram) createDiagram(metadata *DiagramMetadata) (*dg.Diagram, error) {
	name := metadata.Name
	baseDir := metadata.SaveDir
	diag, err := dg.New(SavePath(baseDir, name), dg.Label(name), dg.Direction("TB"))
	if err != nil {
		return nil, err
	}
	return diag, nil
}
