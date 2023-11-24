package mr

//TODO :
//  WORKARGS we need to add the files that the computer has
// WORKARGS the ip of the worker
import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"
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
	lock                 sync.Mutex

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
	c.lock.Lock()
	defer c.lock.Unlock()

	givenTaskType := "wait"
	givenTaskNum := -1
	reply.TaskType = "wait"
	reply.TaskNum = -1

	if c.completeMapTasks < len(c.files) {

		for i := 0; i < len(c.files); i++ {

			if c.mapTasks[i].status == 0 {
				reply.TaskType = "map"
				reply.File = c.files[i]
				reply.TaskNum = i
				reply.NReduce = c.nReduce
				c.mapTasks[i].status = 1
				fmt.Printf("task given")
				givenTaskType = "map"
				givenTaskNum = i
				break
			}
		}
	} else if c.completeMapTasks == len(c.files) {
		if c.completedReduceTasks < c.nReduce {
			for i := 0; i < len(c.reduceTasks); i++ {
				if c.reduceTasks[i].status == 0 {
					reply.File = ""
					reply.TaskType = "reduce"
					reply.TaskNum = i
					reply.NReduce = c.nReduce
					reply.Nfiles = len(c.files)
					c.reduceTasks[i].status = 1
					givenTaskType = "reduce"
					givenTaskNum = i
					fmt.Println(os.Stderr, "reduce task given:", givenTaskNum)
					break

				}
			}
		}
	}
	if reply.TaskType != "wait" {
		go checkIfTaskComplete(givenTaskNum, givenTaskType, c)
	}
	return nil
}

func checkIfTaskComplete(taskNum int, taskType string, c *Coordinator) {

	time.Sleep(10 * time.Second)
	c.lock.Lock()
	defer c.lock.Unlock()

	// Check if the task is still in progress after 10 seconds.
	if taskType == "map" && c.mapTasks[taskNum].status == 1 {
		c.mapTasks[taskNum].status = 0
		fmt.Printf("Map task %d still in progress after 10 seconds\n", taskNum)

	} else if taskType == "reduce" && c.reduceTasks[taskNum].status == 1 {

		c.reduceTasks[taskNum].status = 0
		fmt.Printf("Reduce task %d still in progress after 10 seconds\n", taskNum)
	}
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
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := false

	if c.completedReduceTasks == c.nReduce {
		ret = true
	}
	return ret

}

func (c *Coordinator) ReportWork(args *ReportArgs, reply *ReportReply) error {

	c.lock.Lock()
	defer c.lock.Unlock()

	fmt.Fprintln(os.Stderr, "work done: ", args.TaskSucess, "task type: ", args.TaskType, "task num: ", args.TaskNum)

	if args.TaskSucess == true {

		if args.TaskType == "map" && c.mapTasks[args.TaskNum].status != 2 {
			c.mapTasks[args.TaskNum].status = 2
			c.mapTasks[args.TaskNum].workerID = args.WorkerID
			c.completeMapTasks++

		} else if args.TaskType == "reduce" && c.reduceTasks[args.TaskNum].status != 2 {
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
	}
	for i := 0; i < nReduce; i++ {
		c.reduceTasks = append(c.reduceTasks, reduceTask{0, "", 0})
	}

	//TODO Your code here.

	c.server()
	return &c
}
