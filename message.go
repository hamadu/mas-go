package mas

type Message struct {
	From        Agent
	Task        Task
	BoxA        chan int
	BoxB        chan int
}
