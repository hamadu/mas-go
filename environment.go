package mas
import (
	"log"
)

type Environment struct {
	Agents      []Agent
	TaskQueue   chan Task
	FinishQueue chan int
	Score       int
}

func (e *Environment) GenerateTask(d int) {
	e.TaskQueue <- Task{d}
}

func (e *Environment) GenerateAgent(name string, a int) {
	e.Agents = append(e.Agents, Agent{
		name,
		a,
		0,
		e,
		make(chan Message, 100000),
		map[string]float64{},
		map[string]float64{},
	}, )
}

func (e *Environment) Run() int {
	for i := 0 ; i < len(e.Agents); i++ {
		go e.Agents[i].TryProcessTask()
	}
	score := 0
	for i := 0 ; i < len(e.Agents)/2 ; i++ {
		score += <- e.FinishQueue
		log.Println(i, score)
	}
	return score
}

func SetupEnvironment() Environment {
	return Environment{make([]Agent, 0), make(chan Task, 100000), make(chan int), 0}
}
