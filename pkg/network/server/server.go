package server

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/van-pelt/netcalc/pkg/protolite"
	"log"
	"net"
	"strconv"
	"strings"
)

type CalcServer struct {
	protocol string
	addr     string
	listener net.Listener
	conn     Conn
}

type Conn map[string]*net.Conn

func NewCalcServer(protocol, address string) *CalcServer {
	log.Print("Init Server....OK")
	return &CalcServer{
		protocol: protocol,
		addr:     address,
	}
}

func (c *CalcServer) StartServer() (err error) {

	proto := protolite.NewCalcProtocol()
	c.listener, err = net.Listen(c.protocol, c.addr)
	log.Print("Start Listen [", c.protocol, ":", c.addr, "]...")
	if err != nil {
		return fmt.Errorf("Listen():%w", err)
	}
	log.Print("Listen is OK")
	c.conn = make(map[string]*net.Conn)
	for {
		newConn, err := c.listener.Accept()

		if err != nil {
			log.Print("Warning:", fmt.Errorf("Listen.Accept():%w", err))
			continue
		}

		go c.handleCli(newConn, proto)
	}
	log.Print("Shutdown server")
	return
}

func (c *CalcServer) handleCli(conn net.Conn, protocolMessage *protolite.CalcProtocol) {

	newId := uuid.New().String()
	c.conn[newId] = &conn
	log.Print("New connection ID:", newId, ".Count connection:", len(c.conn))

	defer func() {
		conn.Close()
		delete(c.conn, newId)
		log.Print("Close connection ID:", newId)
	}()

	buf := make([]byte, 64)
	_, err := conn.Write([]byte(newId + "\n"))
	if err != nil {
		log.Print("Warning:", fmt.Errorf("conn.Write():%w", err))
		return
	}

	for {
		readLen, err := conn.Read(buf)
		if err != nil {
			log.Print("Warning:", fmt.Errorf("conn.Read(%v):%w", buf, err))
			break
		}
		cmd := string(buf[:readLen])
		newPack := c.ParseCommand(cmd)
		_, err = conn.Write([]byte(newPack))
		if err != nil {
			log.Print("Warning:", fmt.Errorf("newPack:conn.Write():%w", err))
			return
		}
	}

}

func (c *CalcServer) ParseCommand(cmd string) string {

	newPack := protolite.NetPackage{
		Command:    "",
		Body:       protolite.CmdBody{},
		ErrorsData: protolite.ErrBody{},
		Response:   protolite.CmdResponse{},
	}

	line := strings.Split(cmd, ":")
	if len(line) != 2 {
		newPack.Command = "ERR"
		newPack.ErrorsData.ErrMessage = "BAD FORMAT"
		return newPack.Pack()
	}
	p := protolite.NewCalcProtocol()
	_, ok := (*p)[line[0]]
	if !ok {
		newPack.Command = "ERR"
		newPack.ErrorsData.ErrMessage = ":COMMAND" + line[0] + " NOT FOUND"
		return newPack.Pack()
	}
	param := strings.Split(line[1], ";")
	if len(param) != 5 {
		newPack.Command = "ERR"
		newPack.ErrorsData.ErrMessage = "BAD BODY FORMAT"
		return newPack.Pack()
	}
	arg1, err := strconv.ParseFloat(param[0], 64)
	if err != nil {
		newPack.Command = "ERR"
		newPack.ErrorsData.ErrMessage = fmt.Errorf(":ERR ARGS1 %w", err).Error()
		return newPack.Pack()
	}
	arg2, err := strconv.ParseFloat(param[1], 64)
	if err != nil {
		newPack.Command = "ERR"
		newPack.ErrorsData.ErrMessage = fmt.Errorf(":ERR ARGS2 %w", err).Error()
		return newPack.Pack()
	}
	newPack.Body.Arg1 = arg1
	newPack.Body.Arg2 = arg2
	if line[0] == protolite.CMD_SUM {
		newPack.Command = protolite.CMD_SUM
		newPack.Response.Code = protolite.CMD_SUM
		newPack.Response.Result = newPack.Body.Arg1 + newPack.Body.Arg2
	} else if line[0] == protolite.CMD_MULT {
		newPack.Command = protolite.CMD_MULT
		newPack.Response.Code = protolite.CMD_MULT
		newPack.Response.Result = newPack.Body.Arg1 * newPack.Body.Arg2
	} else if line[0] == protolite.CMD_DIV {
		if newPack.Body.Arg2 == 0 {
			newPack.Command = "ERR"
			newPack.ErrorsData.ErrMessage = "ZERO DIVISION"
			return newPack.Pack()
		}
		newPack.Command = protolite.CMD_DIV
		newPack.Response.Code = protolite.CMD_DIV
		newPack.Response.Result = newPack.Body.Arg1 / newPack.Body.Arg2
	} else if line[0] == protolite.CMD_SUB {
		newPack.Command = protolite.CMD_SUB
		newPack.Response.Code = protolite.CMD_SUB
		newPack.Response.Result = newPack.Body.Arg1 - newPack.Body.Arg2
	} else {
		newPack.Command = "ERR"
		newPack.ErrorsData.ErrMessage = "UNKNOWN ERR"
		return newPack.Pack()
	}
	return newPack.Pack()
}
