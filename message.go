package main

import (
	"time"
	"math/rand"
	"log"
)

type Message struct {
	From        Agent
	Task        Task
	Box         chan int
}
