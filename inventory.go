package hcloudinventory

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

type HostVars struct {
	Ip         string `json:"ip"`
	Location   string `json:"location"`
	Datacenter string `json:"datacenter"`
	Image      string `json:"image"`
}

func newHostVars(Ip string, Location string, Datacenter string, Image string) *HostVars {
	hv := new(HostVars)
	hv.Ip = Ip
	hv.Location = Location
	hv.Datacenter = Datacenter
	hv.Image = Image
	return hv
}

type Meta struct {
	Hostvars map[string]*HostVars `json:"hostvars"`
}

func newMeta() *Meta {
	m := new(Meta)
	m.Hostvars = make(map[string]*HostVars)
	return m
}

func (s *Meta) addHostvar(name string, hostVar *HostVars) {
	s.Hostvars[name] = hostVar
}

type GroupDefinition struct {
	Hosts    []string               `json:"hosts"`
	Vars     map[string]interface{} `json:"vars,omitempty"`
	Children []string               `json:"children,omitempty"`
}

func newGroupDefinition() *GroupDefinition {
	gd := new(GroupDefinition)
	gd.Hosts = []string{}
	gd.Vars = make(map[string]interface{})
	gd.Children = []string{}
	return gd
}

func (s *GroupDefinition) addHost(hostname string) {
	s.Hosts = append(s.Hosts, hostname)
}

func (s *GroupDefinition) addVar(name string, obj interface{}) {
	s.Vars[name] = obj
}

func (s *GroupDefinition) addChild(name string) {
	s.Children = append(s.Children, name)
}

// GetInventoryFromAPI returns a JSON-formatted and Ansible-compatible representation of all virtual servers that are listed under the specified Hetzner Cloud API account.
func GetInventoryFromAPI(client *hcloud.Client) (jsonString string) {
	// Fetch servers from Hetzner Cloud API using it's official golang API client
	serverList, err := client.Server.All(context.Background())
	if err != nil {
		fmt.Errorf("%e", err)
	}

	inventory := make(map[string]interface{})
	meta := newMeta()
	gd := newGroupDefinition()

	for _, hostDef := range serverList {
		hostVars := newHostVars(hostDef.PublicNet.IPv4.IP.String(), hostDef.Datacenter.Location.Name, hostDef.Datacenter.Name, hostDef.Image.Name)
		gd.addHost(hostDef.Name)
		meta.addHostvar(hostDef.Name, hostVars)
		//meta.addHostvar("", newHostVars(serv, Ip, Location, Datacenter, Image))
	}

	inventory["all"] = gd
	inventory["_meta"] = meta

	svlst, _ := json.Marshal(inventory)

	return string(svlst)
}
