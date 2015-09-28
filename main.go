package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/soveran/redisurl"
)

var (
	jobs = make(chan string, 100)
	pool redis.Pool
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	pool = redis.Pool{
		MaxIdle:     3,
		MaxActive:   30,
		IdleTimeout: 240 * time.Second,
		Dial: func() (conn redis.Conn, err error) {
			if len(os.Getenv("REDISCLOUD_URL")) > 0 {
				conn, err = redisurl.ConnectToURL(os.Getenv("REDISCLOUD_URL"))
			} else {
				conn, err = redis.Dial("tcp", ":6379")
			}
			return conn, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	conn := pool.Get()
	_, err := conn.Do("FLUSHALL")
	if err != nil {
		panic(err)
	}
	conn.Close()
}

func main() {
	for w := 1; w <= 3; w++ {
		go phantom(jobs)
	}

	fs := http.FileServer(http.Dir("dist"))
	http.Handle("/", fs)

	http.HandleFunc("/screenshot", screenshot)
	log.Println("listening...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

func phantom(jobs <-chan string) {
	for job := range jobs {
		func(job string) {
			cmd := exec.Command("phantomjs", "rasterize.js", job, "300px*300px", "0.25")
			cmd.Stderr = os.Stderr
			out, err := cmd.Output()
			conn := pool.Get()
			defer conn.Close()
			if err != nil {
				log.Println("Error rasterizing: ", err)
			} else {
				conn.Do("SETEX", job, 21600, string(out))
			}
		}(job)
	}
}

func screenshot(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	req.ParseForm()
	if len(req.Form.Get("url")) > 0 {
		conn := pool.Get()
		defer conn.Close()

		url := req.Form.Get("url")
		log.Println(url)

		exists, err := redis.Bool(conn.Do("EXISTS", url))
		if err != nil || !exists {
			jobs <- url
		}
		for err != nil || !exists {
			if err != nil {
				log.Println(err)
				return
			}
			time.Sleep(500 * time.Millisecond)
			exists, err = redis.Bool(conn.Do("EXISTS", url))
		}
		screenshot, err := redis.String(conn.Do("GET", url))
		if err != nil {
			log.Println(err)
			return
		}
		res.Header().Set("Content-Type", "image/png")
		decode, _ := base64.StdEncoding.DecodeString(screenshot)
		res.Write(decode)
	}
}
