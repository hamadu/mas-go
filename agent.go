package main

import (
	"time"
	"math/rand"
	"log"
)

type Agent struct {
	Name  		string
	Ability		int
	Score 		int
	E           Environment
	Request     chan Message
}

func (a Agent) TryProcessTask() {
	for {
		if rand.Intn(2) <= 0 {
			res := a.tryFetchTask()
			if res == 1 {
				return
			}
		} else {
			a.tryReceiveTask()
		}
	}
}

func (a Agent) pickAgent() Agent {
	k := rand.Intn(len(a.E.Agents) - 1)

	for _, v := range a.E.Agents {
		if v.Name == a.Name {
			continue
		}
		if k <= 0 {
			return v
		}
		k--
	}
	panic("another agent doesnt seem to exsit")
}

func (a Agent) doitSelf(t Task) {
	time.Sleep(time.Duration(100 * t.Difficulty * int(time.Millisecond) / a.Ability / a.Ability))
}

func (a Agent) doitWithOther(t Task, b Agent) {
	ability := a.Ability + b.Ability
	time.Sleep(time.Duration(100 * t.Difficulty * int(time.Millisecond) / ability / ability))
}

func (a Agent) tryFetchTask() int {
	t := <- a.E.TaskQueue
	if t.Difficulty == -1 {
		a.E.FinishQueue <- a.Score
		log.Println("finished",a.Name,a.Score)
		return 1
	}

	if rand.Intn(2) <= 0 {
		// 自処理
		a.Score += t.Difficulty
		a.doitSelf(t)
	} else {
		// 応援呼ぶ
		response := make(chan int, 100)
		target := a.pickAgent()

		target.Request <- Message{a, t, response}
		select {
		case rec := <- response:
			if rec == 1 {
				// いっしょにやる
				a.Score += t.Difficulty
				a.doitWithOther(t, target)
			} else {
				// 断られた。自処理
				a.Score += t.Difficulty
				a.doitSelf(t)
			}
		default:
		// giveup -> 自処理
			a.Score += t.Difficulty
			a.doitSelf(t)
		}
	}
	return 0
}

func (a Agent) tryReceiveTask() {
	select {
	case request := <- a.Request:
		if rand.Intn(2) == 1 {
			request.Box <- 1
			a.doitWithOther(request.Task, request.From)
		} else {
			request.Box <- 0
		}
	default:
	}
}
