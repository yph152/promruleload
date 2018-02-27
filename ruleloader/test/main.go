package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	inputFile := "./config"
	buf, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("File Error: %s\n", err)
		return
	}

	str := string(buf) + "alertname hello\n  for = " + `"` + "1m" + `"` + "\n"
	err = ioutil.WriteFile(inputFile, []byte(str), 0x644)
}

/*type helloHandler struct{}

func (h *helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!\n"))
}

func main() {
	http.Handle("/", &helloHandler{})
	err := http.ListenAndServe(":8899", nil)

	if err != nil {
		fmt.Println("Quit test...")
	}
}*/
