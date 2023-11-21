package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import (
	"os"
	"strconv"
)

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

type WorkArgs struct {
	WorkerdID int
	// filesinmemory []string   maybe used  for distributed system , tell cordinator which files are in storage so he gives the worker a file that is  in storage
}

type WorkReply struct {
	TaskType string
	TaskNum  int
	File     string
	NReduce  int //amount of reduce tasks to create from bucket
}

type ReportArgs struct {
	TaskType   string
	TaskNum    int
	TaskSucess bool
	WorkerID   int
}

type ReportReply struct {
	MoreWork bool
}

// Add your RPC definitions here.

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/5840-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
