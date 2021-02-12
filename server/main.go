package main

import (
	"context"
	"fmt"
	"net"
	"strings"

	prsn "github.com/MURA199918/gRPC_example/person"
	"google.golang.org/grpc"
)

const (
	port = ":3333"
)

//Person struct
type Person struct {
	savedPersons []*prsn.PersonRequest
}

//CreatePerson function
func (p *Person) CreatePerson(c context.Context, input *prsn.PersonRequest) (*prsn.PersonResponse, error) {
	p.savedPersons = append(p.savedPersons, input)
	return &prsn.PersonResponse{Id: input.Id, Success: true}, nil
}

//GetPersons function
func (p *Person) GetPersons(fltr *prsn.PersonFilter, stream prsn.Person_GetPersonServer) error {
	for _, person := range p.savedPersons {
		if fltr.Keyword != "" {
			if !strings.Contains(person.Name, fltr.Keyword) {
				continue
			}
		}
		err := stream.Send(person)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Failed to listen: ", err)
		return
	}
	s := grpc.NewServer()
	prsn.RegisterPersonServer(s, &Person{})
	s.Serve(lis)
}
