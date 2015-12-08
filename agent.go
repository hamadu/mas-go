package mas

import (
	"time"
	"math/rand"
)

type Agent struct {
	Name  		string
	Ability		int
	E           *Environment
	Request     chan Message

	RequestRel  map[string]float64
	ResponseRel map[string]float64
}

const RATE = 0.0

func (a *Agent) ProcessLeaderTask() {
	for {
		res := a.tryFetchTask()
		if res == 1 {
			return
		}
	}
}

func (a *Agent) ProcessWorkerTask() {
	for {
		a.tryReceiveTask()
	}
}

func (a *Agent) pickRequestAgent() Agent {
	sum := 0.0
	for _, v := range a.E.WorkerAgents {
		val, ok := a.RequestRel[v.Name]
		if !ok {
			a.RequestRel[v.Name] = 0.5
			val = 0.5
		}
		sum += val
	}
	r := rand.Float64() * sum

	for _, v := range a.E.WorkerAgents {
		val := a.RequestRel[v.Name]
		if r < val {
			return v
		}
		r -= val
	}
	panic("worker agent doesnt seem to exsit")
}

func (a *Agent) pickResponseAgent() Agent {
	sum := 0.0
	for _, v := range a.E.LeaderAgents {
		val, ok := a.ResponseRel[v.Name]
		if !ok {
			a.ResponseRel[v.Name] = 0.5
			val = 0.5
		}
		sum += val
	}
	r := rand.Float64() * sum

	for _, v := range a.E.LeaderAgents {
		val := a.ResponseRel[v.Name]
		if r < val {
			return v
		}
		r -= val
	}
	panic("leader agent doesnt seem to exsit")
}

func (a *Agent) doitWithoutOther(t Task, b Agent) {
	wait := t.processTime(a)
	time.Sleep(time.Duration(wait))
	a.E.Score += wait / 1000 / 1000

	a.RequestRel[b.Name] = (1.0 - RATE) * a.RequestRel[b.Name]
}

func (a *Agent) doitWithOther(t Task, b Agent) {
	wait := t.processTime(a, &b)
	time.Sleep(time.Duration(wait))
	a.E.Score += wait / 1000 / 1000

	a.RequestRel[b.Name] = (1.0 - RATE) * a.RequestRel[b.Name] + RATE
}

func (a *Agent) doitWithOtherRequest(t Task, b Agent) {
	wait := t.processTime(a, &b)
	time.Sleep(time.Duration(wait))
	a.E.Score += wait / 1000 / 1000

	a.ResponseRel[b.Name] = (1.0 - RATE) * a.ResponseRel[b.Name] + RATE
}

func (a *Agent) doitWithoutOtherRequest(b Agent) {
	a.ResponseRel[b.Name] = (1.0 - RATE) * a.ResponseRel[b.Name]
}

func (a *Agent) tryFetchTask() int {
	var t Task
	select {
	case t = <- a.E.TaskQueue:
	default:
		a.E.FinishQueue <- 1
		return 1
	}

	// 応援呼ぶ
	response := make(chan int, 100)
	response2 := make(chan int, 100)
	target := a.pickRequestAgent()

	target.Request <- Message{*a, t, response, response2}
	select {
	case rec := <- response:
		if rec == 1 {
			response2 <- 1
			// いっしょにやる
			a.doitWithOther(t, target)
		} else {
			// 断られた。自処理
			a.doitWithoutOther(t, target)
		}
	case <-time.After(time.Microsecond * 1000):
		// giveup -> 自処理
		a.doitWithoutOther(t, target)
	}
	close(response2)

	return 0
}

func (a *Agent) tryReceiveTask() {
	<- time.After(time.Microsecond * 250)

	agent := a.pickResponseAgent()

	messages := make([]Message, 0)
	lp: for {
		select {
		case request := <- a.Request:
			messages = append(messages, request)
		default:
			break lp
		}
	}

	found := false
	for _, request := range messages {
		if agent.Name == request.From.Name {
			found = true
			request.BoxA <- 1
			_, ok := <- request.BoxB
			if ok {
				a.doitWithOtherRequest(request.Task, request.From)
			} else {
				a.doitWithoutOtherRequest(request.From)
			}
		} else {
			request.BoxA <- 0
		}
	}
	if !found {
		a.doitWithoutOtherRequest(agent)
	}
}
