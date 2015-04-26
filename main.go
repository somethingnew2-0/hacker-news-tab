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
		MaxActive:   10,
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

	http.HandleFunc("/", screenshot)
	log.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func phantom(jobs <-chan string) {
	for job := range jobs {
		cmd := exec.Command("phantomjs", "rasterize.js", job, "300px*300px", "0.25")
		cmd.Stderr = os.Stderr
		out, err := cmd.Output()
		conn := pool.Get()
		if err != nil {
			log.Println("Error rasterizing: ", err)
		} else {
			conn.Do("HSET", "screenshot", job, string(out))
		}
		conn.Close()
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

		exists, err := redis.Bool(conn.Do("HEXISTS", "screenshot", url))
		if err != nil || !exists {
			jobs <- url
		}
		for err != nil || !exists {
			if err != nil {
				log.Println(err)
				return
			}
			time.Sleep(500 * time.Millisecond)
			exists, err = redis.Bool(conn.Do("HEXISTS", "screenshot", url))
		}
		screenshot, err := redis.String(conn.Do("HGET", "screenshot", url))
		if err != nil {
			log.Println(err)
			return
		}
		res.Header().Set("Content-Type", "image/png")
		decode, _ := base64.StdEncoding.DecodeString(screenshot)
		res.Write(decode)
	}
}
