package hcloudinventory

import (
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func TestInvalidToken(t *testing.T) {
	hetznerClient := hcloud.NewClient(hcloud.WithToken("oiehfoewifhwofihwefoih"))
	GetInventoryFromAPI(hetznerClient)
}
