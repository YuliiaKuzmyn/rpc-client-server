package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"time"
)

type Server struct {}

type Args struct {
	InputFile, OutputFile string
}

func (this *Server) Ping(i int64, reply *string) error {
	fmt.Printf("\n---- PING Transmitted ----\n")
	*reply = "PING 127.0.0.1:8333 successful"
	return nil
}

func (this *Server) Echo(cmd string, reply *string) error {
	fmt.Printf("\n---- ECHO Request received ----\n")
	*reply = cmd
	return nil
}

func (this *Server) Process(args *Args, reply *string) error {
	fmt.Printf("\n---- Process Request received ----\n")
	num, err := readFile(args.InputFile)
	if err != nil {
		*reply = "Process failed with error " + err.Error()
		return err
	}

	selectionSort(num)

	err = writeFile(num, args.OutputFile)
	if err != nil {
		*reply = "Process failed with error " + err.Error()
		return err
	}

	*reply = "Process completed: data written to " + args.OutputFile

	return nil
}

/*writeFile write the data into file*/
func writeFile(p []int, path string) error {

	// open file using READ & WRITE permission
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// write into file
	_, err = file.WriteString(fmt.Sprintln(p))
	if err != nil {
		return err
	}

	// save changes
	err = file.Sync()
	if err != nil {
		return err
	}

	fmt.Printf("\n---- Process completed: data written to %q ----\n", path)
	return nil
}

func readFile(fname string) (nums []int, err error) {
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(b), "\n")
	// Assign cap to avoid resize on every append.
	nums = make([]int, 0, len(lines))

	for _, line := range lines {
		// Empty line occurs at the end of the file when we use Split.
		if len(line) == 0 {
			continue
		}

		strLine := strings.Split(line, " ")

		for _, line := range strLine {
			// Atoi better suits the job when we know exactly what we're dealing
			// with. Scanf is the more general option.
			n, err := strconv.Atoi(line)

			if err != nil {
				return nil, err
			}
			nums = append(nums, n)
		}
	}

	return nums, nil
}

func selectionSort(items []int) {

	start := time.Now()

	var arrayLength = len(items)
	for i := 0; i < arrayLength; i++ {
		var minIdx = i
		for j := i; j < arrayLength; j++ {
			if items[j] < items[minIdx] {
				minIdx = j
			}
		}
		items[i], items[minIdx] = items[minIdx], items[i]
	}

	fmt.Printf("\n---- Selection sort of %d random values took server: %s ----\n", arrayLength, time.Since(start))
}

func server() error {
	rpc.Register(new(Server))
	ln, err := net.Listen("tcp", ":8333")
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("\n---- Server listening on .8333 ----\n")

	for {
		c, err := ln.Accept()
		if err != nil {
			continue
		}

		fmt.Printf("\n---- Connected new peer %s ----\n", c.LocalAddr().String())

		go rpc.ServeConn(c)
	}
}

func main() {
	if err := server(); err != nil {
		fmt.Println(err)
	}
}