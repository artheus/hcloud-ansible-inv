package hcloudinventory

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"regexp"
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud"
	yaml "gopkg.in/yaml.v2"
)

type ExampleHostGrouper struct {
	groupMappings map[string][]string
}

func newExampleHostGrouper() *ExampleHostGrouper {
	hg := new(ExampleHostGrouper)
	hg.groupMappings = make(map[string][]string)
	return hg
}

func (s *ExampleHostGrouper) LoadYaml(fileName string) {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, s.groupMappings)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func (s *ExampleHostGrouper) GroupsForHost(hostname string) (groupNames []string) {
	re := regexp.MustCompile("[a-z]-([a-z0-9]+)[0-9]{2}[0-9]*-[a-z0-9]+")
	submatches := re.FindStringSubmatch(hostname)
	groups := []string{}
	simplehostname := submatches[1]

	for group, hosts := range s.groupMappings {
		if SliceContains(hosts, simplehostname) {
			groups = append(groups, group)
		}
	}

	return groups
}

type MockServerClient struct{}

func (s *MockServerClient) All(ctx context.Context) ([]*hcloud.Server, error) {
	var srvList = []*hcloud.Server{}
	var srvNameList = []string{
		"p-servername01-dc1",
		"p-servername02-dc1",
		"p-otherserver01-dc1",
		"p-tools01-dc1",
		"p-swarmmanager01-dc1",
		"p-swarmmanager02-dc1",
		"p-swarmworker01-dc1",
		"p-swarmworker02-dc1",
	}

	for _, srvname := range srvNameList {
		server := new(hcloud.Server)
		server.Image = new(hcloud.Image)
		server.Image.Name = "imagename"
		server.Name = srvname
		server.PublicNet = hcloud.ServerPublicNet{}
		server.PublicNet.IPv4 = hcloud.ServerPublicNetIPv4{}
		server.PublicNet.IPv4.DNSPtr = "dns-fqdn"
		server.PublicNet.IPv4.IP = net.IPv4(byte(192), byte(168), byte(0), byte(1))
		server.Datacenter = new(hcloud.Datacenter)
		server.Datacenter.Name = "dc1"
		server.Datacenter.Location = new(hcloud.Location)
		server.Datacenter.Location.Name = "rack1"
		srvList = append(srvList, server)
	}

	return srvList, nil
}

func TestWithDefaultGrouper(t *testing.T) {
	var srvClient ServerClient
	mockClient := new(MockServerClient)
	srvClient = mockClient
	grouper := newDefaultGrouper()
	grouper.LoadYaml("test/default_grouper_groups.yml")
	inventory := GetInventoryFromAPI(srvClient, grouper)

	inventory.UpdateAllGroup()

	//fmt.Println(inventory.Json())
}

func TestWithCustomGrouper(t *testing.T) {
	var srvClient ServerClient
	mockClient := new(MockServerClient)
	srvClient = mockClient
	grouper := newExampleHostGrouper()
	grouper.LoadYaml("test/example_grouper_groups.yml")
	inventory := GetInventoryFromAPI(srvClient, grouper)

	inventory.UpdateAllGroup()
	inventory.UpdateAllGroup() // Twice only for the coverage..
	inventory.Json()
	//fmt.Println(inventory.Json())

	// TODO: Actually make some assertions to make sure that we are getting the right results.
}
