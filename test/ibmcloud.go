package test

import (
	"errors"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/codeready-toolchain/devcluster/pkg/ibmcloud"
)

type MockIBMCloudClient struct {
	mux            sync.RWMutex
	clustersByID   map[string]*ibmcloud.Cluster
	clustersByName map[string]*ibmcloud.Cluster
}

func NewMockIBMCloudClient() *MockIBMCloudClient {
	return &MockIBMCloudClient{
		clustersByName: make(map[string]*ibmcloud.Cluster),
		clustersByID:   make(map[string]*ibmcloud.Cluster),
	}
}

func (c *MockIBMCloudClient) GetZones() ([]string, error) {
	return []string{
		"ams03",
		"che01",
		"dal10",
		"dal12",
		"dal13",
		"fra02",
		"fra04",
		"fra05",
		"hkg02",
		"lon02",
		"lon04",
		"lon05",
		"lon06",
		"mel01",
		"mex01",
		"mil01",
		"mon01",
		"osl01",
		"par01",
		"sao01",
		"seo01",
		"sjc03",
		"sjc04",
		"sng01",
		"syd01",
		"syd04",
		"syd05",
		"tok02",
		"tok04",
		"tok05",
		"tor01",
		"wdc04",
		"wdc06",
		"wdc07",
	}, nil
}

func (c *MockIBMCloudClient) GetVlans(zone string) ([]ibmcloud.Vlan, error) {
	return []ibmcloud.Vlan{}, nil
}

func (c *MockIBMCloudClient) CreateCluster(name, zone string) (string, error) {
	defer c.mux.Unlock()
	c.mux.Lock()
	if c.clustersByName[name] != nil {
		return "", errors.New("cluster already exist")
	}
	newCluster := &ibmcloud.Cluster{
		ID:          uuid.NewV4().String(),
		Name:        name,
		Region:      zone,
		CreatedDate: time.Now().String(),
		State:       "deploying",
		Ingress:     ibmcloud.Ingress{},
	}
	c.clustersByName[name] = newCluster
	c.clustersByID[newCluster.ID] = newCluster
	return newCluster.ID, nil
}

func (c *MockIBMCloudClient) GetCluster(id string) (*ibmcloud.Cluster, error) {
	defer c.mux.RUnlock()
	c.mux.RLock()
	return c.clustersByID[id], nil
}

func (c *MockIBMCloudClient) UpdateCluster(cluster ibmcloud.Cluster) error {
	defer c.mux.Unlock()
	c.mux.Lock()
	found := c.clustersByID[cluster.ID]
	if found == nil {
		return errors.New("cluster not found")
	}
	found.State = cluster.State
	found.Ingress = cluster.Ingress
	return nil
}
