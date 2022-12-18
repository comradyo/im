package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Event interface {
	Calculate(clock int64)
	IsDone() bool
}

type Simulation struct {
	clock int64
	events []Event
}

func (s *Simulation) Iterate() {
	s.clock++
	eventsLen := len(s.events)
	for i := 0; i < eventsLen; i++ {
		if s.events[i].IsDone() {
			eventsLen--
			if i != eventsLen {
				s.events = append(slice[:i], slice[i+1:]...)
			}
			i--
		} else {
			s.events[i].Calculate(s.clock)
		}
	}
}

func (s *Simulation) AddEvent(e Event) {
	s.events = append(s.events, e)
}

func (s *Simulation) Run() {
	for s.clock != 10000 {
		s.Iterate()
	}
}

type computer struct {
	website Balancer
}

type User struct {
	userID int64
}

func (u *User) Calculate(clock int64) {
	if 
}






const (
	// количество пользователей
	numberOfUsers = 10
	// количество серверов
	numberOfServers = 5
	// среднее количество повторных запросов, которые отправляет пользователь, если ему приходит таймаут
	numberOfResends = 10
	// количество потоков - количество запросов, которые сервер может одновременно обработать
	numberOfThreads = 50
	// время, за которое сервер обрабатывает запрос
	handleTime = 5000
	// время (в единицах времени симуляции), при превышении которого клиенту отдается таймаут
	timeout = 50000
	// максимальное число таймаутов, при превышении которого сервер перезагружается
	maxTimeoutTimes = 5
	// время, характеризующее частоту запросов от пользователей
	userRequestTime = 1000000
)

// структура, отвечающая за симуляцию программного времени
type Clock struct {
	mu            *sync.RWMutex
	t             int64
}

func NewProgramTime() *ProgramTime {
	return &ProgramTime{
		mu:            &sync.RWMutex{},
		t:             0,
	}
}

func (p *ProgramTime) SetTime(t int64) SimulationTime {

}

func (p *ProgramTime) GetTime() SimulationTime {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.t
}

// Тип, в котором хранятся id пользователей
type ID string

// Запрос пользователя
type Request struct {
	userID ID
	// время, когда пришел запрос от пользователя
	requestTime SimulationTime
}

// Статус ответа
type Status int64

const (
	StatusOK Status = iota
	StatusTimeout
)

// Ответ от сервера
type Response struct {
	status Status
	// время, когда сервер отдал ответ
	responseTime SimulationTime
}

// Пользователь
type User struct {
	id ID
}

// Создать запрос
func (u *User) GenerateRequest(rTime SimulationTime) Request {
	return Request{
		userID:      u.id,
		requestTime: rTime,
	}
}

// Повторная генерация запросов пользователем
func (u *User) RegenerateRequests(rTime SimulationTime) []Request {
	// будет сгенерировано повторных запросов от 0 до resends * 2
	rsnds := rand.Intn(numberOfResends*2) + 1
	requests := make([]Request, rsnds)
	for i := 0; i < rsnds; i++ {
		requests[i] = Request{
			userID:      u.id,
			requestTime: rTime,
		}
	}
	return requests
}

func (u *U)

type Server struct {
	mu *sync.RWMutex
	// число потоков на сервере
	threads int64
}

func (s *Server) HandleRequest(req *Request) Response {
	responseTime := req.requestTime
	s.mu.Lock()
	if s.threads <= 0 {
		for s.threads <= 0 {
		}
		s.mu.Lock()
		responseTime = progTime.GetTime()
	}
	for s.threads <= 0 {
		responseTime += 1
	}
	/*
		если все серверы заняты, то встаем в очередь
		while threads <= 0 {
			wait
		}
		handleTime := rand(averageHandleTime) + deltaHandleTime
		return ok
	*/
	return Response{}
}

// Балансировщик нагрузки
type Balancer struct {
	usersToServers SyncMap[ID, int]
	servers        []*Server
	// количество таймаутов, произошедших на серверах
	serversTimeoutsCount []int64
}

func (b *Balancer) getServerID(userID ID) int {
	serverID, ok := b.usersToServers.Load(userID)
	if !ok {
		panic("user to server not found")
	}
	return serverID
}

func (b *Balancer) HandleRequest(req *Request) Status {
	// получение номера сервера, к которому надо обращаться
	serverID, ok := b.usersToServers.Load(req.userID)
	if !ok {
		panic("user to server not found")
	}
	// получение самого сервера
	b.mu.RLock()
	server := b.servers[serverID]
	b.mu.RUnlock()
	res := server.HandleRequest(req)
	processingTime := res.responseTime - req.requestTime
	if processingTime >= timeout {
		b.serversTimeoutsCount[serverID]++
		return StatusTimeout
	}
	b.serversTimeoutsCount[serverID] = 0
	return StatusOK
}

func (b *Balancer) CheckAndBalance() bool {
	serversToRestart := make([]*Server, 0, len(b.servers))
	serversThatAreOk := make([]*Server, 0, len(b.servers))
	for i := range b.serversTimeoutsCount {
		if b.serversTimeoutsCount[i] > maxTimeoutTimes {
			serversToRestart = append(serversToRestart, b.servers[i])
		} else {
			serversToRestart = append(serversThatAreOk, b.servers[i])
		}
	}
	if len(serversToRestart) > len(serversThatAreOk) {
		return false
	}
	return false
}

func (b *Balancer) 

func test(i int) {
	go func() {
		for {
			t := time.Now()
			val := progTime.GetTime()
			fmt.Printf("test %d val = %v, time = %v\n", i, val, t)
		}
	}()
}

func main() {
	progTime = NewProgramTime(tickDuration)
	go progTime.Run()

	for i := 0; i < 5; i++ {
		go test(i, mu)
	}
	time.Sleep(time.Second * 100)
}

var progTime *ProgramTime
