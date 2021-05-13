package client

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/van-pelt/netcalc/pkg/protolite"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type CalcClient struct {
	protocol           string
	addr               string
	ServerConnectionId string
	conn               net.Conn
}

func NewClient(protocol, address string) *CalcClient {
	return &CalcClient{
		protocol:           protocol,
		addr:               address,
		ServerConnectionId: "",
		conn:               nil,
	}
}

func (l *CalcClient) StartCalcTerminal() (err error) {
	conn, err := net.Dial(l.protocol, l.addr)

	if err != nil {
		return fmt.Errorf("Can`t create connection:%w", err)
	}
	defer conn.Close()
	buf := make([]byte, 64)
	readLen, err := conn.Read(buf)
	if err != nil {
		log.Print("Warning:", fmt.Errorf("conn.Read(%v):%w", buf, err))
	}
	l.ServerConnectionId = string(buf[:readLen])
	log.Print("SERVER:Ok,you ID:", l.ServerConnectionId)
	for {
		str, err := l.parseCommand(os.Stdin)

		if err != nil {
			return err
		}

		if str == "BYE" {
			return nil
		}

		_, err = conn.Write([]byte(str))
		if err != nil {
			log.Print("Warning:", fmt.Errorf("conn.Write():%w", err))
			continue
		}
		buf := make([]byte, 64)
		_, err = conn.Read(buf)
		if err != nil {
			log.Print("Warning:", fmt.Errorf("conn.Read(%v):%w", buf, err))
			continue
		}
		result, err := l.parseResponse(string(buf))
		if err != nil {
			log.Print("Warning:", fmt.Errorf("ERROR RESULT COMMAND:%w", err))
			continue
		}
		log.Print("SERVER:Ok,Result:", result)
	}

	return nil
}

func (l *CalcClient) parseCommand(src io.Reader) (string, error) {
	reader := bufio.NewReader(src)
	fmt.Print(": ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("reader.ReadString():%w", err)
	}
	if text == "BYE" {
		return "BYE", nil
	}
	cmd := strings.Split(text, ":")
	if len(cmd) != 2 {
		return "", fmt.Errorf(":ERR COMMAND\n")
	}
	p := protolite.NewCalcProtocol()
	body, ok := (*p)[cmd[0]]
	if !ok {
		return "", fmt.Errorf(":COMMAND " + cmd[0] + " NOT FOUND\n")
	}
	args := strings.Split(cmd[1], ";")
	if len(cmd) != 2 {
		return "", fmt.Errorf(":ERR ARGS\n")
	}
	body.Body.Arg1, err = strconv.ParseFloat(args[0], 64)
	if err != nil {
		return "", fmt.Errorf(":ERR ARGS1 %w", err)
	}
	body.Body.Arg2, err = strconv.ParseFloat(strings.TrimSuffix(args[1], "\r\n"), 64)
	if err != nil {
		return "", fmt.Errorf(":ERR ARGS2 %w", err)
	}
	return body.Pack(), nil
}

func (l *CalcClient) parseResponse(resp string) (float64, error) {

	pack := protolite.NetPackage{
		Command:    "",
		Body:       protolite.CmdBody{},
		ErrorsData: protolite.ErrBody{},
		Response:   protolite.CmdResponse{},
	}
	cmd := strings.Split(resp, ":")
	if len(cmd) != 2 {
		return 0.0, fmt.Errorf("Responce bad format")
	}
	body := strings.Split(cmd[1], ";")
	if len(body) != 5 {
		return 0.0, fmt.Errorf("Responce bad format")
	}
	pack.Command = cmd[0]
	arg1, err := strconv.ParseFloat(strings.TrimSuffix(body[0], "\r\n"), 64)
	if err != nil {
		return 0.0, fmt.Errorf(":ERR ARGS1 %w", err)
	}
	arg2, err := strconv.ParseFloat(strings.TrimSuffix(body[1], "\r\n"), 64)
	if err != nil {
		return 0.0, fmt.Errorf(":ERR ARGS2 %w", err)
	}
	bs := bytes.Trim([]byte(body[4]), "\x00")
	result, err := strconv.ParseFloat(strings.TrimSuffix(string(bs), "\r\n"), 64)
	if err != nil {
		return 0.0, fmt.Errorf(":ERR RESULT %w", err)
	}
	pack.Body.Arg1 = arg1
	pack.Body.Arg2 = arg2
	pack.ErrorsData.ErrMessage = body[2]
	pack.Response.Code = body[3]
	pack.Response.Result = result

	if cmd[0] == "ERR" {
		return 0.0, fmt.Errorf(":ERR DATA:CODE=ERR,MESS=%v", pack.ErrorsData.ErrMessage)
	}
	return pack.Response.Result, nil
}

func (l *CalcClient) copyTo(dst io.Writer, src io.Reader) {

	fmt.Println(dst, src)

	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

func (l *CalcClient) PrintHelp() {
	fmt.Println("	Commands:")
	fmt.Println("	Syntax:CMD:arg1;arg2")
	fmt.Println("	--------------------")
	fmt.Println("	(+)->SUM:arg1;arg2")
	fmt.Println("	(/)->DIV:arg1;arg2")
	fmt.Println("	(*)->MLT:arg1;arg2")
	fmt.Println("	(-)->SUB:arg1;arg2")
	fmt.Println("		 BYE - exit")
}
