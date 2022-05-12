package main

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"

	"github.com/l-ross/miio/client"
	"github.com/l-ross/miio/vacuum"
)

func main() {
	ip := os.Getenv("MIIO_IP")
	token := os.Getenv("MIIO_TOKEN")

	c, err := client.New(token, client.SetIP(ip))
	if err != nil {
		log.Fatal(err)
	}

	v, err := vacuum.New(c)
	if err != nil {
		log.Fatal(err)
	}

	info, err := v.Info()
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(info)

	//m, err := v.Map()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(m)
}
