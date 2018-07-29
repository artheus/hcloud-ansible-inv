package hcloudinventory

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/hetznercloud/hcloud-go/hcloud"
	yaml "gopkg.in/yaml.v2"
)

// HostVars is host specific variables to be stored in meta
type HostVars struct {
	IP         string `json:"ansible_host"`
	DNS        string `json:"hcloud_dns"`
	Location   string `json:"hcloud_location"`
	Datacenter string `json:"hcloud_datacenter"`
	Image      string `json:"hcloud_image"`
}

// Meta is the _meta tag in the resulted Inventory
type Meta struct {
	Hostvars map[string]*HostVars `json:"hostvars"`
}

// GroupDefinition is used to define a group of hosts
type GroupDefinition struct {
	Hosts    []string               `json:"hosts"`
	Vars     map[string]interface{} `json:"vars,omitempty"`
	Children []string               `json:"children,omitempty"`
}

// Host is a tool struct to let Inventory keep track of groups
type Host struct {
	Name   string
	Groups []string
}

// Inventory is the main struct which will be returned as a result of GetInventoryFromAPI
type Inventory struct {
	inventory map[string]interface{}
	allHosts  []*Host
}

func newHostVars(ip string, dns string, location string, datacenter string, image string) *HostVars {
	hv := new(HostVars)
	hv.IP = ip
	hv.DNS = dns
	hv.Location = location
	hv.Datacenter = datacenter
	hv.Image = image
	return hv
}

func newMeta() *Meta {
	m := new(Meta)
	m.Hostvars = make(map[string]*HostVars)
	return m
}

func (s *Meta) addHostvar(name string, hostVar *HostVars) {
	s.Hostvars[name] = hostVar
}

func newInventory() *Inventory {
	i := new(Inventory)
	i.inventory = make(map[string]interface{})
	return i
}

// SetMeta will set the meta field of the inventory
func (s *Inventory) SetMeta(meta *Meta) {
	s.inventory["_meta"] = meta
}

// AddHost will add a Host to the inventory
func (s *Inventory) AddHost(host *Host) {
	s.allHosts = append(s.allHosts, host)
}

// Group is used for getting and/or creating a group for storing hosts
func (s *Inventory) Group(name string) *GroupDefinition {
	if _, ok := s.inventory[name]; !ok {
		s.inventory[name] = newGroupDefinition()
	}

	return s.inventory[name].(*GroupDefinition)
}

// UpdateAllGroup will refresh the Inventory group named all
func (s *Inventory) UpdateAllGroup() {
	if _, ok := s.inventory["all"]; ok {
		s.inventory["all"] = nil
	}

	allGroup := newGroupDefinition()
	for group := range s.inventory {
		if strings.HasPrefix(group, "_") {
			continue
		}
		allGroup.addChild(group)
	}

	for _, host := range s.allHosts {
		if len(host.Groups) == 0 {
			allGroup.addHost(host.Name)
		}
	}

	s.inventory["all"] = allGroup
}

// JSON will output a json representation of the Inventory
func (s *Inventory) JSON() (jsonString string) {
	output, _ := json.Marshal(s.inventory)
	return string(output)
}

func newHost(name string, groups []string) *Host {
	h := new(Host)
	h.Name = name
	h.Groups = groups
	return h
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

// ServerClient is an interface for allowing mocking in tests
type ServerClient interface {
	All(ctx context.Context) ([]*hcloud.Server, error)
}

// Grouper is an interface for the grouping of hosts
type Grouper interface {
	// GroupsForHost should return any groups that the host is a part of
	GroupsForHost(hostname string) (groupNames []string)
}

// DefaultGrouper is the default implementation of Grouper
type DefaultGrouper struct {
	// GroupMappings is a mapping of hosts, like group => host1, host2 ...
	GroupMappings map[string][]string
}

func newDefaultGrouper() *DefaultGrouper {
	hg := new(DefaultGrouper)
	hg.GroupMappings = make(map[string][]string)
	return hg
}

// LoadYaml loads a YAML file and tries to Unmarshal it into the GroupMappings type
func (s *DefaultGrouper) LoadYaml(fileName string) {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, s.GroupMappings)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

// GroupsForHost will return any groups that the host is a part of
func (s *DefaultGrouper) GroupsForHost(hostname string) (groupNames []string) {
	groups := []string{}

	for group, hosts := range s.GroupMappings {
		if sliceContains(hosts, hostname) {
			groups = append(groups, group)
		}
	}

	return groups
}

func sliceContains(slice []string, item string) (result bool) {
	for _, sliceItem := range slice {
		if sliceItem == item {
			return true
		}
	}
	return false
}

// GetInventoryFromAPI returns a JSON-formatted and Ansible-compatible representation of all virtual servers that are listed under the specified Hetzner Cloud API account.
func GetInventoryFromAPI(client ServerClient, grouper Grouper) *Inventory {
	if grouper == nil {
		defaultGrouper := newDefaultGrouper()
		defaultGrouper.LoadYaml("groups.yml")
		grouper = defaultGrouper
	}

	// Fetch servers from Hetzner Cloud API using it's official golang API client
	serverList, err := client.All(context.Background())
	if err != nil {
		fmt.Printf("%e", err)
	}

	inventory := newInventory()
	meta := newMeta()

	for _, hostDef := range serverList {
		hostVars := newHostVars(hostDef.PublicNet.IPv4.IP.String(), hostDef.PublicNet.IPv4.DNSPtr, hostDef.Datacenter.Location.Name, hostDef.Datacenter.Name, hostDef.Image.Name)
		meta.addHostvar(hostDef.Name, hostVars)
		groups := grouper.GroupsForHost(hostDef.Name)
		host := newHost(hostDef.Name, groups)
		inventory.AddHost(host)

		for _, groupName := range groups {
			inventory.Group(groupName).addHost(hostDef.Name)
		}
	}

	inventory.SetMeta(meta)

	return inventory
}
