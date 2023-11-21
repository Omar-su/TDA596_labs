package mr

//TODO :
// IMPLEMENT SLEEP AND CHECK FOR DEAD WORKERS
// IMPLEMENT REDUCE TASKS
// IMPLEMENT LOCKS ON COORDINATOR
// IMPLEMENT DISTRIBUTED

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type mapTask struct {
	nReduce                   int
	status                    int
	file                      string
	workerID                  int
	intermediateFileLocations []string
}

type reduceTask struct {
	status   int
	file     string
	workerID int
}

type Coordinator struct {
	nReduce              int
	files                []string
	completeMapTasks     int
	completedReduceTasks int
	mapTasks             []mapTask
	reduceTasks          []reduceTask
	// Your definitions here.

}

// Your code here -- RPC handlers for the worker to call.

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

func (c *Coordinator) GiveWork(args *WorkArgs, reply *WorkReply) error {

	// print the first file name in C

	givenTaskType := ""
	givenTaskNum := -1

	fmt.Printf(c.files[0])

	if c.completeMapTasks < len(c.files) {

		for i := 0; i < len(c.files); i++ {

			if c.mapTasks[i].status == 0 {
				reply.TaskType = "map"
				reply.File = c.files[i]
				reply.TaskNum = i
				reply.NReduce = c.nReduce
				c.mapTasks[i].status = 1
				fmt.Printf("task given")
				break
			}
		}
	} /* else if c.completeMapTasks == len(c.files) {

		if len(c.reduceTasks) < c.nReduce {
			for i := 0; i < len(c.reduceTasks); i++ {
				if c.reduceTasks[i].status == 0 {
					reply.File = ""
					reply.taskType = "reduce"
					reply.TaskNum = i
					reply.nReduce = c.nReduce
					c.reduceTasks = append(c.reduceTasks, reduceTask{0, reply.File, 0})
					break
				}
			}
		}
	} */
	return nil
}

//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	ret := false

	if c.completedReduceTasks == c.nReduce {
		ret = true
	}
	return ret

}

func (c *Coordinator) ReportWork(args *ReportArgs, reply *ReportReply) error {

	if args.TaskType == "map" {
		c.mapTasks[args.TaskNum].status = 2
		c.mapTasks[args.TaskNum].workerID = args.WorkerID
		c.completeMapTasks++

	} else if args.TaskType == "reduce" {
		c.reduceTasks[args.TaskNum].status = 2
		c.reduceTasks[args.TaskNum].workerID = args.WorkerID
		c.completedReduceTasks++
	}
	for i := 0; i < len(c.mapTasks); i++ {
		if c.mapTasks[i].status == 0 {
			reply.MoreWork = true
			return nil
		}
	}
	for i := 0; i < len(c.reduceTasks); i++ {
		if c.reduceTasks[i].status == 0 {
			reply.MoreWork = true
			return nil
		}
	}
	reply.MoreWork = false
	return nil
}

//
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}
	c.files = files
	c.nReduce = nReduce
	c.completeMapTasks = 0
	c.completedReduceTasks = 0
	//make empty map tasks equal to number of files

	for i := 0; i < len(files); i++ {
		c.mapTasks = append(c.mapTasks, mapTask{nReduce, 0, files[i], 0, []string{}})
		c.reduceTasks = append(c.reduceTasks, reduceTask{0, "", 0})
	}

	//TODO Your code here.

	c.server()
	return &c
}
