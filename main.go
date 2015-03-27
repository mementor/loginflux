package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/influxdb/influxdb/client"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type CommandLine struct {
	Host      string
	Port      int
	User      string
	Pass      string
	Values    string
	Separator string
	Name      string
}

const (
	default_host = "localhost"
	default_port = 8086
	default_user = "root"
	default_pass = "root"
)

func main() {

	c := CommandLine{}

	fs := flag.NewFlagSet("default", flag.ExitOnError)
	fs.StringVar(&c.Host, "host", default_host, "host to connect to")
	fs.IntVar(&c.Port, "port", default_port, "port to connect to")
	fs.StringVar(&c.User, "user", default_user, "user to auth on db")
	fs.StringVar(&c.Pass, "pass", default_host, "pass to auth on db")
	fs.StringVar(&c.Values, "vals", "value", "values to set coma separated")
	fs.StringVar(&c.Separator, "sep", ",", "values separator")
	fs.StringVar(&c.Name, "name", "", "collection")
	fs.Parse(os.Args[1:])

	u := url.URL{
		Scheme: "http",
	}
	u.Host = fmt.Sprintf("%s:%d", c.Host, c.Port)

	u.User = url.UserPassword(c.User, c.Pass)

	cl, err := client.NewClient(
		client.Config{
			URL:       u,
			Username:  c.User,
			Password:  c.Pass,
			UserAgent: "loginflux",
		})
	if err != nil {
		fmt.Println("FAIL")
	} else {
		//fmt.Println("OK")
		//dur, str, err := cl.Ping()
		//if err != nil {
		//	fmt.Printf("FAIL: '%v', '%v'\n", dur, str)
		//} else {
		//	fmt.Printf("OK: '%v', '%v'\n", dur, str)
		//}
	}
	buf := make([]byte, 1024)
	br := bufio.NewReader(os.Stdin)
	batchpoints := client.BatchPoints{}
	batchpoints.Database = "webui"
	for {
		// read a chunk
		n, err := br.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}

		// write a chunk
		input := fmt.Sprintf("%s", buf[:n-1])
		//line = strings.Replace(line, "\n", "", -1)
		for _, line := range strings.Split(input, "\n") {
			fields := map[string]interface{}{}
			var timestamp string
			for num, str := range strings.Split(line, c.Separator) {
				fmt.Printf("num = '%d', str = '%s'\n", num, str)
				if num == 0 {
					timestamp = str
					continue
				}
				fields["test"] = str
			}
			ts, _ := strconv.ParseInt(timestamp, 10, 64)
			point := client.Point{}
			point.Name = c.Name
			point.Fields = fields
			point.Timestamp = time.Unix(ts, 0)
			batchpoints.Points = append(batchpoints.Points, point)
		}
		res, err := cl.Write(batchpoints)
		if err != nil {
			fmt.Printf("FAIL: '%v'", res)
		} else {
			fmt.Printf("OK: '%v'", res)
		}
		batchpoints.Points = []client.Point{}
	}
}
