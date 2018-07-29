package hcloudinventory

import (
	"context"
	"net"
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

type MockServerClient struct{}

func (s *MockServerClient) All(ctx context.Context) ([]*hcloud.Server, error) {
	var srvList = []*hcloud.Server{}

	server := new(hcloud.Server)
	server.Image = new(hcloud.Image)
	server.Image.Name = "imagename"
	server.Name = "name"
	server.PublicNet = hcloud.ServerPublicNet{}
	server.PublicNet.IPv4 = hcloud.ServerPublicNetIPv4{}
	server.PublicNet.IPv4.DNSPtr = "dns-fqdn"
	server.PublicNet.IPv4.IP = net.IPv4(byte(192), byte(168), byte(0), byte(1))
	server.Datacenter = new(hcloud.Datacenter)
	server.Datacenter.Name = "datacenter"
	server.Datacenter.Location = new(hcloud.Location)
	server.Datacenter.Location.Name = "location"
	srvList = append(srvList, server)

	return srvList, nil
}

func TestInvalidToken(t *testing.T) {
	var srvClient ServerClient
	mockClient := new(MockServerClient)
	srvClient = mockClient
	GetInventoryFromAPI(srvClient)
}
