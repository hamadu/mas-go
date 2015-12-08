package mas

type Environment struct {
	LeaderAgents []Agent
	WorkerAgents []Agent
	TaskQueue    chan Task
	FinishQueue  chan int
	Score        int
}

func (e *Environment) GenerateTask(d int) {
	e.TaskQueue <- Task{d}
}

func (e *Environment) GenerateLeaderAgent(name string, a int) {
	e.LeaderAgents = append(e.LeaderAgents, Agent{
		name,
		a,
		e,
		make(chan Message, 100000),
		map[string]float64{},
		map[string]float64{},
	}, )
}

func (e *Environment) GenerateWorkerAgent(name string, a int) {
	e.WorkerAgents = append(e.WorkerAgents, Agent{
		name,
		a,
		e,
		make(chan Message, 100000),
		map[string]float64{},
		map[string]float64{},
	}, )
}

func (e *Environment) Run() int {
	for i := 0 ; i < len(e.LeaderAgents); i++ {
		go e.LeaderAgents[i].ProcessLeaderTask()
	}
	for i := 0 ; i < len(e.WorkerAgents); i++ {
		go e.WorkerAgents[i].ProcessWorkerTask()
	}
	score := 0
	for i := 0 ; i < len(e.LeaderAgents) ; i++ {
		score += <- e.FinishQueue
	}
	return score
}

func SetupEnvironment() Environment {
	return Environment{
		make([]Agent, 0),
		make([]Agent, 0),
		make(chan Task, 100000),
		make(chan int),
		0,
	}
}
