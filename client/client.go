package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"os/exec"
	"strings"
)

type rpcConn struct {
	client *rpc.Client
}

type Args struct {
	InputFile, OutputFile string
}

func printHelp() {
	fmt.Printf("\n---- Client Usage:" +
		"\n * ping (no arguments) - to check connection with server *" +
		"\n * echo [OPTION]... TEXT TO PRINT - to send text print command to server *" +
		"\n * process [INPUT FILE] [OUTPUT FILE] - to process selection sort on server side *" +
		"\n * help (no arguments) - to show help *" +
		"\n * clear, press [ENTER](no arguments) - to clear screen *\n\n")
}

func (c *rpcConn) handleEchoRequest(args string) error {
	var result string
	err := c.client.Call("Server.Echo", args, &result)
	if err != nil {
		return err
	} else {
		fmt.Printf("\n---- Echo() = %s ----\n", result)
		return nil
	}
}

func (c *rpcConn) handlePingRequest() error {
	var result string
	err := c.client.Call("Server.Ping", 0, &result)
	if err != nil {
		return err
	} else {
		fmt.Printf("\n---- Ping() = %s ----\n", result)
		return nil
	}

}

func (c *rpcConn) handleProcessRequest(args string) error {
	parts := strings.SplitN(args, " ", 2)

	if len(parts) != 2 {
		return fmt.Errorf("\nWrong number of arguments, see help for more info\n")
	}

	request := &Args {
		parts[0],
		parts[1],
	}

	var result string
	err := c.client.Call("Server.Process", request, &result)
	if err != nil {
		return err
	} else {
		fmt.Printf("\n---- Process() = %s ----\n", result)
		return nil
	}
}

func (c *rpcConn) parseInputCmd(input string) error {
	parts := strings.SplitN(input, " ", 2)

	switch parts[0] {
	case "", "clear":
		c := exec.Command("clear")
		c.Stdout = os.Stdout
		c.Run()
	case "echo", "Echo", "ECHO":
		if len(parts) > 1 {
			return c.handleEchoRequest(parts[1])
		} else {
			fmt.Printf("\n---- Not enough parameters for echo cmd, see help ----\n")
			printHelp()
			return nil
		}
	case "ping", "Ping", "PING":
		return c.handlePingRequest()
	case "process", "Process", "PROCESS":
		if len(parts) > 1 {
			return c.handleProcessRequest(parts[1])
		} else {
			fmt.Printf("\n---- Not enough parameters for process cmd, see help ----\n")
			printHelp()
			return nil
		}
	case "help", "h":
		printHelp()
	default:
		fmt.Printf("\n---- Received wrong command ----\n")
		printHelp()
		return nil
	}

	return nil
}

func (c *rpcConn) StartCli() error {
	scanner := bufio.NewScanner(os.Stdin)
	printHelp()
	fmt.Printf("\n---- Waiting for user input: ----\n")
	for scanner.Scan() {
		err := c.parseInputCmd(scanner.Text())
		if err != nil {
			return err
		}
	}

	return nil
}

func client() error {
	client, err := rpc.Dial("tcp", "127.0.0.1:8333")
	if err != nil {
		return err
	}

	fmt.Printf("\n---- Client connected to .8333 ----\n")

	conn := rpcConn{
		client,
	}

	err = conn.StartCli()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := client(); err != nil {
		fmt.Printf(err.Error())
	}
}