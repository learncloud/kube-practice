package main

import(
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
)

var version = os.Getenv("version")
var connection int32

func main(){
	log.Printf("%s / starting process on %v",version, os.Getpid())

	var status int

	if version == "v1"{
		status = 201
	} else if version =="v2"{
		status = 202
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Println(version, req.URL.Path)

		defer func() {
			atomic.AddInt32(&connection, -1)
		}()
		atomic.AddInt32(&connection, 1)
		// /sleep/N 요청에는 N초간 슬립모드
		if strings.HasPrefix(req.URL.Path, "/sleep") {
			id := strings.TrimPrefix(req.URL.Path, "/sleep/")
			i, _ := strconv.Atoi(id)
			time.Sleep(time.Second * time.Duration(i))
		}
		w.WriteHeader(status)
	})

	//SIGTERM, SIGINT 무시
	signalChannel :=make(chan os.Signal)
	signal.Notify(signalChannel,syscall.SIGTERM, syscall.SIGINT)
	go func(){
		for{
			log.Println(version, "/ connection", atomic.LoadInt32(&connection))
			time.Sleep(time.Second)
		}		
	}()

	http.ListenAndServe(":5000",nil)

}

/* 
웹서버 기능

	1. version 변수값에 따라 v1은 http 상태코드 201,v2는 202로 변환
	2. sleep요청에 n초간 슬립
	3. sigterm 신호 무시
	-> 이미지 교체시 트래픽 손실이 있는지 확인하는 용도
*/