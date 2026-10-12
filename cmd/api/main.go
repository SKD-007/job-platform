package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"vendor/golang.org/x/net/idna"

	"golang.org/x/text/date"
)

type JobsQueuetype struct{
	Queue  map[int]Job
	mutex sync.RWMutex
}

type Job struct{
	name string
	status string
	id int
	jobtype int
}

type IDGenerator interface{
	NextID() int
}

func (i * GlbIndx) NextID() int{
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.gbIndex++
	return i.gbIndex
}

type GlbIndx struct{
	gbIndex int
	mutex sync.RWMutex
}

var glbIndx GlbIndx
var JobsQueue JobsQueuetype

func ServerCheck(w http.ResponseWriter, r *http.Request ){
	fmt.Printf("Server Working....\n")
	w.Write([]byte("Working"))
}


func Addjobs(w http.ResponseWriter, r *http.Request ){
	
	data, err := io.ReadAll(r)
	defer r.Body.Close()

	if( err != nil){
		http.Error(w, "Failed to read the request Body\n", http.StatusInternalServerError)
		return
	}
	
	
	var obj Job  
	err = json.Unmarshal(data, &obj)
	
	if( err != nil ){
		http.Error(w, "Failed to parse the Json data\n", http.StatusInternalServerError)
		return
	}

	indx := glbIndx.NextID()
	obj.id = indx

	JobsQueue.Queue[indx] = obj

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	str := fmt.Sprint("Sucessfully sent the job with id = %d", indx)
	w.Write([]byte(str))
}


func DisplayJobs(w http.ResponseWriter, r *http.Request ){
	for k, v := range JobsQueue.Queue {
		fmt.Printf("job id = %d, job_name = %v, jpb")
	}
}

func main()  {
	mux := http.NewServeMux()

	http.HandleFunc("/hello", ServerCheck)
	http.HandleFunc("/job", Addjobs)
	http.HandleFunc("/jobs", DisplayJobs)

	http.ListenAndServe(":8080", mux)

}