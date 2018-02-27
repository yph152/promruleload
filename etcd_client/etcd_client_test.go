package etcd_client

import (
	"fmt"
	"testing"
)

func TestEtcd(t *testing.T) {
	config := []string{"http://127.0.0.1:2379"}
	kapi, err := NewClient(config)

	if err != nil {
		t.Error("NewClient error")
	}
	err = kapi.Set("key", "value")

	if err != nil {
		t.Error("Set error")
	}

	_, err = kapi.Get("key")

	if err != nil {
		t.Error("Get error")
	} else {
		t.Log("Get success")
	}

	err = kapi.Delete("key")

	if err != nil {
		t.Error("Delete Error")
	}

	list, err := kapi.List("/api/v1/autoscalers/default")

	if err != nil {
		t.Error("List Error")
	}

	fmt.Println(list)

}
