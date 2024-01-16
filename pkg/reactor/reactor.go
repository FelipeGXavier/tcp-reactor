package reactor

import (
	"eventloop/pkg/data"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Reactor struct {
	incomingEvents data.Queue
	handlers map[string](func(net.Conn, []byte)data.Result)
	workersCount int 
	incomingEventsChannel chan net.Conn
	processedEventsChannel chan data.Result
	waitGroup *sync.WaitGroup
}

func MakeReactor(workersCount int) Reactor {
	incomingEventsQueue := data.MakeLinkedQueue()
	incomingEventsChannel := make(chan net.Conn)
	processedEventsChannel := make(chan data.Result)
	handlersMap := make(map[string]func(net.Conn, []byte)data.Result, 0)
	var waitGroup sync.WaitGroup

	reactor := Reactor{
		incomingEventsChannel: incomingEventsChannel,
		processedEventsChannel: processedEventsChannel,
		incomingEvents: incomingEventsQueue,
		handlers: handlersMap,
		workersCount: workersCount,
		waitGroup: &waitGroup,
	}

	addDefaultHandlers(&reactor)
	reactor.startHandleTasks()
	go reactor.responseHandler()
	go reactor.pollTasks()

	return reactor
}


func (reactor *Reactor) AddHandler(key string, execution func(net.Conn, []byte)data.Result) {
	reactor.handlers[key] = execution
}

func (reactor *Reactor) Stop() {
	log.Default().Print("Stopping the server...")
	reactor.waitGroup.Wait()
	log.Default().Print("Server exit 0")
	os.Exit(0)
}

func (reactor *Reactor) Run() {
	server, err := net.Listen("tcp", ":3090")
	if err != nil {
		log.Fatal("Error while start listening on :3090", err)
		panic(err)
	}
	log.Default().Print("Listening on port :3090")
	for {
		connection, err := server.Accept()
		if err != nil {
			log.Fatal("Error while accepting connection", err)
			continue
		}
		reactor.incomingEvents.Add(connection)
	}
}

func (reactor *Reactor) responseHandler() {
	for result := range reactor.processedEventsChannel {
		reactor.waitGroup.Done()
		connection := result.Socket 
		connection.Write([]byte(result.Data))
	}
}

func (reactor *Reactor) pollTasks() {
	for {
		task := reactor.incomingEvents.Pop()
		if task != nil {
			reactor.incomingEventsChannel <- task.(net.Conn)
		}
	}
}

func (reactor *Reactor) startHandleTasks() {
	for i := 0; i < reactor.workersCount; i++ {
		go func ()  {
			for task := range reactor.incomingEventsChannel {
				reactor.waitGroup.Add(1)
				sleepTime := rand.Intn(3) + 1
				time.Sleep(time.Second * time.Duration(sleepTime))
				resultResponse := reactor.parsePacket(task)
				log.Default().Print("Processing TCP packet")
				reactor.processedEventsChannel <- resultResponse
			}
		}()
	}
}

func addDefaultHandlers(reactor *Reactor) {
	reactor.AddHandler("SUM", SumHandler)
	reactor.AddHandler("MULT", MultiplyHandler)
}

func (reactor Reactor) parsePacket(connection net.Conn) (data.Result) {
	buffer := make([]byte, 4096)
	connection.Read(buffer)
	packetData := string(buffer)
	parts := strings.Split(packetData, " ")

	if len(parts) <= 1 {
		return data.Result{Socket: connection, Data: "Invalid packet data"}
	}

	command := strings.ToUpper(parts[0])

	handler := reactor.handlers[command]

	if handler == nil {
		return data.Result{Socket: connection, Data: "Invalid packet data"}
	}

	result := handler(connection, buffer)

	return result
}
