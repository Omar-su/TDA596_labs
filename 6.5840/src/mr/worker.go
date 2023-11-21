package mr

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
)

//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

//
// main/mrworker.go calls this function.
//
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	// Your worker implementation here.

	// uncomment to send the Example RPC to the coordinator.
	CallAskForWork(mapf, reducef)

}

//
// example function to show how to make an RPC call to the coordinator.
//
// the RPC argument and reply types are defined in rpc.go.
//
func CallAskForWork(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {
	workArgs := WorkArgs{WorkerdID: 1}
	workReply := WorkReply{}

	ok := call("Coordinator.GiveWork", &workArgs, &workReply)
	if ok {
		fmt.Printf(workReply.File)
	} else {
		fmt.Printf("does not work")
	}

	fmt.Printf("reply.tasktype %s\n", workReply.TaskType)
	fmt.Printf("reply.tasknum %v\n", workReply.TaskNum)
	fmt.Printf("reply.file %s\n", workReply.File)
	fmt.Printf("reply.nreduce %v\n", workReply.NReduce)

	if workReply.TaskType == "map" {
		processMapTask(mapf, reducef, workReply.File, workReply.NReduce, workReply.TaskNum)
	} else if workReply.TaskType == "reduce" {
		processReduceTask(reducef, workReply.File, workReply.TaskNum)

	}
	

}

func processMapTask(mapf func(string, string) []KeyValue, reducef func(string, []string) string, filename string, nReduce int, taskNum int) {

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("cannot open %v", filename)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read %v", filename)
	}
	file.Close()
	kva := mapf(filename, string(content))

	for i := range kva {
		bucketNum := ihash(kva[i].Key) % nReduce

		//check if file named mr-mapTaskNum-bucketNum exists if not create it  and append key value pair to it , otherwise append key value pair to it
		filename := fmt.Sprintf("mr-%d-%d", taskNum, bucketNum)
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("cannot open %v", filename)
		}
		file.WriteString(kva[i].Key + " " + kva[i].Value + "\n")
		file.Close()
	}

	fmt.Printf("map task done")
	ReportWork("map", taskNum, true, mapf, reducef)
}

func processReduceTask(reducef func(string, []string) string, file string, taskNum int) {

}

func CallReportWork() {

}

func ReportWork(taskType string, taskNum int, success bool, mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	args := ReportArgs{TaskType: taskType, TaskNum: taskNum, TaskSucess: success}
	reply := ReportReply{}

	ok := call("Coordinator.ReportWork", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("reply.Y")
	} else {
		fmt.Printf("call failed!\n")
	}

	if reply.MoreWork {
		CallAskForWork(mapf, reducef)
	}

}

func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	// the "Coordinator.Example" tells the
	// receiving server that we'd like to call
	// the Example() method of struct Coordinator.
	ok := call("Coordinator.Example", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("reply.Y %v\n", reply.Y)
	} else {
		fmt.Printf("call failed!\n")
	}
}

//
// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
