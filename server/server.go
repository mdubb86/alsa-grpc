package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
	pb "server/grpc_gen"
)

const (
	port = ":50051"
)

type AlsaServer struct {
	card string
	device string
	monitor Monitor
	ctrlIdMap map[string]int32
	updates  chan *ControlInfo
	requests chan *pb.Request
	clients []*pb.Alsamixer_CommunicateServer
	clientMux sync.Mutex
	grpcServer *grpc.Server
}

func (s *AlsaServer) Start() error {

	cards, err := aplayList()
	if err != nil {
		return err
	}

	if cardNum, found := cards[s.card]; found {
		s.device = fmt.Sprintf("hw:%d", cardNum)
	} else {
		return fmt.Errorf("unable to identify card %s", s.card)
	}

	fmt.Printf("Card %s is present as %s\n", s.card, s.device)

	// Get control states
	controls, err := amixerContents(s.device)
	s.ctrlIdMap = make(map[string]int32)
	for _, ctrl := range controls {
		s.ctrlIdMap[ctrl.name] = ctrl.id
	}
	fmt.Printf("Map", s.ctrlIdMap)

	// Monitor alsa controls
	s.updates = make(chan *ControlInfo)
	s.clients = []*pb.Alsamixer_CommunicateServer{}

	s.monitor = Monitor{
		device: s.device,
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
			if id, found := s.ctrlIdMap[req.Control]; found {
				fmt.Printf("Setting volume level for %s to %d\n", req.Control, req.Volume)
				cset(s.device, id, req.Volume)
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
	// Build update
	update := pb.Response{Card: s.card, Controls: []*pb.Response_Control{createResponseControl(ctrl)}}

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

	// Get current states
	current, err := amixerContents(s.device)

	// Send current states to client that just connected
	if err == nil {
		fmt.Println("Gathering updates", current)
		updates := make([]*pb.Response_Control, len(current))
		for i, ctrl := range current {
			updates[i] = createResponseControl(&ctrl)
		}
		fmt.Println("Sending initial states", updates)
		stream.Send(&pb.Response{Card: s.card, Controls: updates})
	} else {
		fmt.Println(err)
	}

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

	server := AlsaServer{
		card: "snd_rpi_hifiberry_amp",
	}
	server.Start()
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
