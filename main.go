package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// форма, в которой производится отсчет времени в программе
type SimulationTime int64

type ProgramTime struct {
	mu            *sync.RWMutex
	t             SimulationTime
	sleepDuration time.Duration
}

func initProgramTime(sleepDuration time.Duration) {
	progTime = ProgramTime{
		mu:            &sync.RWMutex{},
		t:             0,
		sleepDuration: sleepDuration,
	}
}

func (p *ProgramTime) Run() {
	for {
		time.Sleep(time.Microsecond)
		p.t++
	}
}

func (p *ProgramTime) GetTime() SimulationTime {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.t
}

var progTime ProgramTime

const (
	// среднее количество повторных запросов, которые отправляет пользователь, если ему приходит ошибка
	resends = 10
	// максимальное количество запросов, которое сервер может одновременно обрабатывать
	capability = 50
)

// Тип, который из себя представляют id пользователей
type ID string

// Запрос пользователя
type Request struct {
	userID ID
	// время, когда пришел запрос от пользователя
	startTime SimulationTime
}

// статус ответа
type Status int64

const (
	StatusOK Status = iota
	StatusTimeout
)

// Ответ от сервера
type Response struct {
	status Status
	// время, когда сервер отдал ответ
	endTime SimulationTime
}

// Пользователь
type User struct {
	id ID
}

// Создать запрос
func (u *User) GenerateRequest() Request {
	return Request{
		userID:    u.id,
		startTime: progTime.GetTime(),
	}
}

// Повторная генерация запросов пользователем
func (u *User) RegenerateRequests() []Request {
	// будет сгенерировано повторных запросов от 0 до resends * 2
	rsnds := rand.Intn(resends*2) + 1
	requests := make([]Request, rsnds)
	startTime := progTime.GetTime()
	for i := 0; i < rsnds; i++ {
		requests[i] = Request{
			userID:    u.id,
			startTime: startTime,
		}
	}
	return requests
}

type Server struct {
	userIDs []ID
	// число потоков на сервере
	mu      *sync.Mutex
	threads int64
}

func (s *Server) HandleRequest(req *Request) Response {
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
	mu             *sync.Mutex
	usersToServers map[ID]int
	servers        []Server
}

func (b *Balancer) HandleRequest(req *Request) Status {
	serverID, ok := b.usersToServers[req.userID]
	if !ok {
		panic("user to server not found")
	}
	res := b.servers[serverID].HandleRequest(req)
	res.endTime
}

func main() {
	fmt.Printf("Hello, World!")

	initProgramTime(time.Microsecond)
	go func() {
		progTime.Run()
	}()
	for {
		fmt.Println(progTime.GetTime())
	}
}
