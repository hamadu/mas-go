package main

import (
	"fmt"
	"time"
	"math/rand"
	"log"
)

type Environment struct {
	Agents      []Agent
	TaskQueue   chan Task
	FinishQueue chan int
}

func (e Environment) GenerateTask(d int) {
	e.TaskQueue <- Task{d}
}
