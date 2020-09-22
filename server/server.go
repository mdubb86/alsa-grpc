package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"os"
	"reflect"
	pb "server/grpc_gen"
	"strings"
	"sync"
)

const (
	port = ":50051"
)

type AlsaServer struct {
	cardName   string
	cardNum    int32
	monitor    Monitor
	controls   map[string]*ControlInfo
	updates    chan *ControlInfo
	requests   chan *pb.Request
	clients    []*pb.Alsamixer_CommunicateServer
	clientMux  sync.Mutex
	grpcServer *grpc.Server
}

func (s *AlsaServer) Start(cardName string, controlNames []string) error {
	s.cardName = cardName

	cards, err := aplayList()
	if err != nil {
		return err
	}

	if cardNum, found := cards[s.cardName]; found {
		s.cardNum = cardNum
	} else {
		return fmt.Errorf("unable to identify card %s in %v", s.cardName, cards)
	}

	fmt.Printf("Card %s is present as card number %d\n", s.cardName, s.cardNum)

	// Prepare controlMap
	s.controls = make(map[string]*ControlInfo, len(controlNames))
	for _, controlName := range controlNames {
		s.controls[controlName] = nil
	}

	// Get control IDs from amixer
	allControls, err := amixerContents(s.cardNum)
	for i := 0; i < len(allControls); i++ {
		control := allControls[i]
		fmt.Printf("Discovered Control %s (%p -- %p): %v\n", control.name, &allControls[i], &control, control)
		if _, found := s.controls[control.name]; found {
			s.controls[control.name] = &control
		} else {
			fmt.Println("Ignoring", control.name)
		}
	}

	// Check for missing all_controls
	for controlName, control := range s.controls {
		if control == nil {
			return fmt.Errorf("unable to find control %s on card %s (%d): %v", controlName, s.cardName,
				s.cardNum, allControls)
		} else {
			fmt.Printf("Control %s: %v\n", controlName, control)
		}
	}

	// Monitor alsa all_controls
	s.updates = make(chan *ControlInfo)
	s.clients = []*pb.Alsamixer_CommunicateServer{}

	s.monitor = Monitor{
		cardNum: s.cardNum,
		updates:  s.updates,
	}
	s.monitor.Start()

	go func() {
		for {
			s.processUpdateFromMonitor(<-s.updates)
		}
	}()

	s.requests = make(chan *pb.Request)
	go func() {
		for {
			req := <-s.requests
			if ctrl, found := s.controls[req.Control]; found {
				fmt.Printf("Setting volume level for %s to %d\n", req.Control, req.Volume)
				cset(s.cardNum, ctrl.id, req.Volume)
			} else {
				fmt.Printf("Unrecognized control %s for\n", req.Control)
			}
		}
	}()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s.grpcServer = grpc.NewServer()
	pb.RegisterAlsamixerService(s.grpcServer, &pb.AlsamixerService{
		Communicate: s.communicate,
	})

	if err := s.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	 return nil
}

func createResponseControl(ctrl *ControlInfo) *pb.Response_Control {
	return &pb.Response_Control{
		Name:   ctrl.name,
		Volume: ctrl.percents,
	}
}

func (s *AlsaServer) processUpdateFromMonitor(ctrl *ControlInfo) {
	// Store update
	if _, found := s.controls[ctrl.name]; found {
		fmt.Println("Applying update for", ctrl.name)
		s.controls[ctrl.name] = ctrl
	} else{
		fmt.Println("Ignoring update to", ctrl.name)
	}

	// Build update
	update := pb.Response{Card: s.cardName, Controls: []*pb.Response_Control{createResponseControl(ctrl)}}

	// Broadcast update
	fmt.Printf("Broadcasting update to %s to %d clients (%v)\n", ctrl.name, len(s.clients), &s.clients)
	s.clientMux.Lock()
	defer s.clientMux.Unlock()
	var unavailableClients []*pb.Alsamixer_CommunicateServer
	for _, client := range s.clients {
		if err := (*client).Send(&update); err != nil {
			if status.Code(err) == codes.Unavailable {
				unavailableClients = append(unavailableClients, client)
			} else {
				fmt.Println("Client send error", err, reflect.TypeOf(err))
			}
		}
	}

	if len(unavailableClients) > 0 {
		fmt.Printf("Processing unavailable clients Clients %d unvailable %d\n", len(s.clients), len(unavailableClients))
		s.removeUnavailable(unavailableClients)
		fmt.Printf("Num remaining clients %d\n", len(s.clients))
	}

	fmt.Println("Done broadcasting update")
}

func (s *AlsaServer) addClient(client *pb.Alsamixer_CommunicateServer) {
	s.clientMux.Lock()
	defer s.clientMux.Unlock()
	fmt.Printf("Clients before: %d (%v)\n", len(s.clients), &s.clients)
	s.clients = append(s.clients, client)
	fmt.Printf("Clients after: %d (%v)\n", len(s.clients), &s.clients)
}

func (s *AlsaServer) removeUnavailable(unavailables []*pb.Alsamixer_CommunicateServer) {
	numRemainingClients := len(s.clients) - len(unavailables)
	if numRemainingClients > 0 {
		isAvailable := func(chk *pb.Alsamixer_CommunicateServer) bool {
			for _, unavailable := range unavailables {
				if chk == unavailable {
					return false
				}
			}
			return true
		}
		var stillAvailable []*pb.Alsamixer_CommunicateServer
		for _, client := range s.clients {
			if isAvailable(client) {
				stillAvailable = append(stillAvailable, client)
			}
		}
		s.clients = stillAvailable
	} else {
		s.clients = nil
	}
}

func (s *AlsaServer) communicate(stream pb.Alsamixer_CommunicateServer) error {
	s.addClient(&stream)

	// Send current states to client that just connected
	fmt.Println("Sending initial", s.controls)
	var updates []*pb.Response_Control
	for _, ctrl := range s.controls {
		updates = append(updates, createResponseControl(ctrl))
	}
	stream.Send(&pb.Response{Card: s.cardName, Controls: updates})

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("Client is disconnected")
			return nil
		}
		if err != nil {
			fmt.Println("Client is disconnected by error")
			return err
		}
		fmt.Println("Received req", req, &s, s.requests, &s.requests)
		s.requests <- req
	}
}

func main() {
	// Parse controls argument
	if len(os.Args) != 3 {
		log.Fatal("Must provide 2 arguments 'card name(s)' 'control name(s)'")
	}

	cardName := os.Args[1]
	controlNames := strings.Split(os.Args[2], ",")

	server := AlsaServer{}
	err := server.Start(cardName, controlNames)
	if err != nil {
		log.Panic(err)
	}
}

//go func() {
//	for {
//		err := stream.Send(&pb.Update{Volume: 0.5})
//		if err != nil {
//			fmt.Println(err)
//		}
//		time.Sleep(1 * time.Second)
//	}
//}()
