package main

import (
	"flag"
	"fmt"
	"github.com/manifoldco/promptui"
	"keep/components"
)

type pepper struct {
	Name     string
	HeatUnit int
	Peppers  int
}

func main() {

	defer fmt.Println("\ngoodbye!")

	var err error
	endpoints := args()

	etcd, err := components.NewEtcd(endpoints)
	if err != nil {
		fmt.Println(err)
		return
	}

	pd, err := components.NewPlacementDriver(endpoints[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	s := components.Server{
		Etcd:            etcd,
		PlacementDriver: pd,
	}

	prompt := promptui.Select{
		Label: " what can i do for you ",
		Items: []string{
			"cdc",
		},
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	switch result {
	case "cdc":
		err = s.Cdc()
	default:
	}
	if err != nil {
		fmt.Println(err)
		return
	}
}

func args() []string {
	var pd string
	flag.StringVar(&pd, "pd", "", "127.0.0.1:2379")
	flag.Parse()

	if pd == "" {
		pd = "172.16.5.133:2379"
	}
	var endpoints []string
	endpoints = append(endpoints, pd)
	return endpoints
}
