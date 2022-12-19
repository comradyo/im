package main

import (
	"fmt"
	"math/rand"
	"time"
)

var computersBreakAt = []int64{
	15, 300,
}

const (
	sleepDuration = time.Millisecond * 1000
	//////////////////
	simulationDuration = 100000
	queueLenLimit      = 30
	//////////////////
	numOfCores            = 3
	numOfServers          = 5
	numOfUsers            = 50
	numOfRepeatedRequests = 10
	//////////////////
	maxProcessQueueLen = 10
	//////////////////
	handleDuration     = 10 //100
	timeOut            = handleDuration * 5
	checkServersPeriod = 30
	userRequestPeriod  = 10
)

type Request struct {
	userID  string
	reqTime int64
	resTime int64
}

type User struct {
	id              string
	createRequestAt int64
	currentRequest  *Request
	website         *Website
}

func (u *User) Calculate(t int64) {
	if u.createRequestAt == 0 {
		u.createRequestAt = rand.Int63n(userRequestPeriod) + userRequestPeriod/2
	}
	if t == u.createRequestAt {
		u.currentRequest = u.createRequest(t)
		u.website.HandleRequest(u.currentRequest)
	}
	if u.currentRequest != nil {
		if t-u.currentRequest.reqTime >= timeOut {
			repeatedRequests := u.repeatRequest(t)
			for i := range repeatedRequests {
				u.website.HandleRequest(repeatedRequests[i])
			}
		} else if t == u.currentRequest.resTime {
			u.currentRequest = nil
			u.createRequestAt = 0
		}
	}
}

func (u *User) Break(t int64) {
	repeatedRequests := u.repeatRequest(t)
	for i := range repeatedRequests {
		u.website.HandleRequest(repeatedRequests[i])
	}
	u.createRequestAt = t
}

func (u *User) createRequest(t int64) *Request {
	return &Request{
		userID:  u.id,
		reqTime: t,
	}
}

func (u *User) repeatRequest(t int64) []*Request {
	numOfRequests := rand.Intn(numOfRepeatedRequests-numOfRepeatedRequests/2) + numOfRepeatedRequests/2
	requests := make([]*Request, numOfRequests)
	for i := range requests {
		requests[i] = &Request{
			userID:  u.id,
			reqTime: t,
		}
	}
	u.currentRequest = requests[len(requests)-1]
	return requests
}

type Core struct {
	jobStartsAt int64
	jobEndsAt   int64
}

func (c *Core) Calculate(t int64) {
	if t >= c.jobEndsAt {
		c.jobStartsAt = 0
		c.jobEndsAt = 0
	}
}

func (c *Core) SetJob(startsAt int64) {
	jobDuration := rand.Int63n(handleDuration) + handleDuration/2
	c.jobStartsAt = startsAt
	c.jobEndsAt = c.jobStartsAt + jobDuration
}

type Server struct {
	cores         []Core
	requestsQueue []*Request
}

func (s *Server) Calculate(t int64) {
	for i := range s.cores {
		s.cores[i].Calculate(t)
	}
	s.processRequestsQueue(t)
}

func (s *Server) HandleRequest(req *Request) {
	s.requestsQueue = append(s.requestsQueue, req)
}

func (s *Server) processRequestsQueue(t int64) {
	for i := 0; i < len(s.requestsQueue); i++ {
		availableCoreID, found := s.availableCore()
		if found {
			s.cores[availableCoreID].SetJob(t)
			s.requestsQueue[i].resTime = s.cores[availableCoreID].jobEndsAt
			s.requestsQueue = s.requestsQueue[1:]
		} else {
			i = len(s.requestsQueue)
		}
	}
}

func (s *Server) availableCore() (int, bool) {
	for i := range s.cores {
		if s.cores[i].jobEndsAt == 0 {
			return i, true
		}
	}
	return 0, false
}

type Website struct {
	userIDtoServerID map[string]int
	servers          []Server
}

func NewWebsite() *Website {
	servers := make([]Server, numOfServers)
	for i := range servers {
		cores := make([]Core, numOfCores)
		var queue []*Request
		servers[i] = Server{
			cores:         cores,
			requestsQueue: queue,
		}
	}
	mp := make(map[string]int, numOfUsers)
	return &Website{
		servers:          servers,
		userIDtoServerID: mp,
	}
}

func (w *Website) Calculate(t int64) {
	for i := range w.servers {
		w.servers[i].Calculate(t)
	}
	if t%checkServersPeriod == 0 {
		maxQueueLen, serverID := w.longestProcessQuery()
		if maxQueueLen > maxProcessQueueLen {
			w.rewriteUserIDs(serverID)
		}
	}
}

func (w *Website) HandleRequest(req *Request) {
	serverID, ok := w.userIDtoServerID[req.userID]
	if !ok {
		panic("FUCK")
	}
	w.servers[serverID].HandleRequest(req)
}

func (w *Website) longestProcessQuery() (int, int) {
	maxQueueLen := len(w.servers[0].requestsQueue)
	serverID := 0
	for i := 1; i < len(w.servers); i++ {
		if len(w.servers[i].requestsQueue) > maxQueueLen {
			maxQueueLen = len(w.servers[i].requestsQueue)
			serverID = i
		}
	}
	return maxQueueLen, serverID
}

func (w *Website) rewriteUserIDs(downServerID int) {
	minQueueLen := len(w.servers[0].requestsQueue)
	newServerID := 0
	for i := 1; i < len(w.servers); i++ {
		if len(w.servers[i].requestsQueue) < minQueueLen {
			minQueueLen = len(w.servers[i].requestsQueue)
			newServerID = i
		}
	}
	if minQueueLen > queueLenLimit {
		panic("Очень большая очередь")
	}
	for userId, serverID := range w.userIDtoServerID {
		if serverID == downServerID {
			w.userIDtoServerID[userId] = newServerID
		}
	}
}

func (w *Website) RegisterUsers(users []User) {
	k := 0
	for i := range users {
		w.userIDtoServerID[users[i].id] = k
		k++
		if k == len(w.servers) {
			k = 0
		}
	}
}

type Clock struct {
	ws   *Website
	usrs []User
	//////////
	serverUsers []string
}

func NewClock(ws *Website, usrs []User) *Clock {
	serverUsers := make([]string, numOfServers)
	return &Clock{
		ws:          ws,
		usrs:        usrs,
		serverUsers: serverUsers,
	}
}

var (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

func (c *Clock) Run() {
	for t := int64(0); t < simulationDuration; t++ {
		time.Sleep(sleepDuration)
		for i := range computersBreakAt {
			if t == computersBreakAt[i] {
				userIDToBreak := rand.Intn(numOfUsers)
				fmt.Printf(
					"\n\n%sСкоро пользователь %d отправит много запросов%s\n\n",
					colorRed, userIDToBreak, colorReset,
				)
				c.usrs[userIDToBreak].Break(t)
				time.Sleep(time.Second * 3)
			}
		}
		c.ws.Calculate(t)
		for i := range c.usrs {
			c.usrs[i].Calculate(t)
		}
		c.print(t)
	}
}

func (c *Clock) print(t int64) {
	fmt.Print("\033[H\033[2J")
	fmt.Printf("\n\n")
	fmt.Printf("Время от начала симуляции: %s%d%s\n", colorYellow, t, colorReset)
	for i := range c.serverUsers {
		c.serverUsers[i] = ""
	}
	for u, s := range c.ws.userIDtoServerID {
		c.serverUsers[s] += " " + u
	}
	fmt.Printf("Распределение пользователей по серверам:\n")
	for i := range c.serverUsers {
		fmt.Printf("\t Server %s%d%s: %s%s%s\n", colorCyan, i, colorReset, colorGreen, c.serverUsers[i], colorReset)
	}
	for i := range c.ws.servers {
		fmt.Printf("server [%d]: \n", i)
		fmt.Printf("\tqueue: %s%v%s\n", colorRed, len(c.ws.servers[i].requestsQueue), colorReset)
		for j := range c.ws.servers[i].cores {
			fmt.Printf("\tcore [%d]: %+v\n", j, c.ws.servers[i].cores[j])
		}
	}
	fmt.Printf("\n\n")
}

func main() {
	rand.Seed(time.Now().UnixNano())
	website := NewWebsite()
	users := make([]User, numOfUsers)
	for i := range users {
		users[i].id = fmt.Sprintf("%d", i)
		users[i].website = website
	}
	website.RegisterUsers(users)
	clock := NewClock(website, users)
	clock.Run()
}
