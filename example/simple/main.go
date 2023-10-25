package main

import (
	"fmt"

	helper "github.com/snowlyg/win-helper"
)

func (p *program) Start(s helper.Service) error {
	// do some work
	return nil
}

func (p *program) Stop(s helper.Service) error {
	//stop
	return nil
}

type program struct{}

func main() {
	// new windows service
	s, err := helper.NewService(&program{}, &helper.Config{Name: "service-name"})
	if err != nil {
		fmt.Printf("new service get error %v \n", err)
	}
	s.Run()
}
