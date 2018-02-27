package cache

import (
	"fmt"
	"testing"
)

func TestCache(t *testing.T) {
	cc := NewCache()

	status := cc.Set("hello", "world")
	if !status {
		fmt.Println("Set failed...")
		return
	}
	str := cc.Get("hello")
	fmt.Println(str)

	status = cc.Update("hello", "world1")
	if !status {
		fmt.Println("Update failed...")
		return
	}
	str = cc.Get("hello")
	fmt.Println(str)

	stri := "hello"
	for _, value := range cc.Rulemap {
		fmt.Println("key")
		fmt.Println(value)
	}
	cc.Delete(stri)

	/*if !status {
		fmt.Println("Delete failed...")
		return
	}*/
	str = cc.Get("hello")
	fmt.Println(str)

}
