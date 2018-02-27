package etcd_client

import (
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

type EtcdClient struct {
	Client client.KeysAPI
}

func NewClient(config []string) (*EtcdClient, error) {
	cfg := client.Config{
		Endpoints: config,
		Transport: client.DefaultTransport,
	}
	c, err := client.New(cfg)

	api := client.NewKeysAPI(c)

	kapi := &EtcdClient{Client: api}

	return kapi, err
}

func (kapi *EtcdClient) Set(key, value string) error {
	_, err := kapi.Client.Set(context.Background(), key, value, nil)

	return err
}

func (kapi *EtcdClient) Get(key string) (string, error) {
	resp, err := kapi.Client.Get(context.Background(), key, nil)
	if err != nil {
		return "", err
	}

	return resp.Node.Value, err
}

func (kapi *EtcdClient) Delete(key string) error {
	_, err := kapi.Client.Delete(context.Background(), key, nil)

	return err
}

func (kapi *EtcdClient) List(dir string) ([]string, error) {
	resp, err := kapi.Client.Get(context.Background(), dir, nil)

	if err != nil {
		return nil, err
	}

	var list []string
	for _, value := range resp.Node.Nodes {
		list = append(list, value.Key)
	}

	return list, err
}

func (kapi *EtcdClient) Update(key, value string) error {
	_, err := kapi.Client.Update(context.Background(), key, value)

	return err
}
