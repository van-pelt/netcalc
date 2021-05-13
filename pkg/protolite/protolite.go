package protolite

import "fmt"

const (
	CMD_SUM  = "SUM" //+
	CMD_DIV  = "DIV" // /
	CMD_MULT = "MLT" //*
	CMD_SUB  = "SUB" //-

)

type CalcProtocol map[string]NetPackage

type NetPackage struct {
	Command    string
	Body       CmdBody
	ErrorsData ErrBody
	Response   CmdResponse
}

func (n *NetPackage) Pack() string {
	return fmt.Sprintf("%s:%s;%s;%s", n.Command, n.Body.String(), n.ErrorsData.String(), n.Response.String())
}

type CmdBody struct {
	Arg1 float64
	Arg2 float64
}

func (b *CmdBody) String() string {
	return fmt.Sprintf("%f;%f", b.Arg1, b.Arg2)
}

type CmdResponse struct {
	Code   string
	Result float64
}

func (r *CmdResponse) String() string {
	return fmt.Sprintf("%s;%f", r.Code, r.Result)
}

type ErrBody struct {
	ErrMessage string
}

func (e *ErrBody) String() string {
	return fmt.Sprintf("%s", e.ErrMessage)
}

// CMD_FLAG:arg1;arg2;err_message;CMD_FLAG;result

func NewCalcProtocol() *CalcProtocol {
	return &CalcProtocol{
		CMD_SUM: NetPackage{
			Command: CMD_SUM,
			Body: CmdBody{
				Arg1: 0,
				Arg2: 0,
			},
			ErrorsData: ErrBody{ErrMessage: ""},
			Response: CmdResponse{
				Code:   "",
				Result: 0,
			},
		},
		CMD_DIV: NetPackage{
			Command: CMD_DIV,
			Body: CmdBody{
				Arg1: 0,
				Arg2: 0,
			},
			ErrorsData: ErrBody{ErrMessage: ""},
			Response: CmdResponse{
				Code:   "",
				Result: 0,
			},
		},
		CMD_MULT: NetPackage{
			Command: CMD_MULT,
			Body: CmdBody{
				Arg1: 0,
				Arg2: 0,
			},
			ErrorsData: ErrBody{ErrMessage: ""},
			Response: CmdResponse{
				Code:   "",
				Result: 0,
			},
		},
		CMD_SUB: NetPackage{
			Command: CMD_SUB,
			Body: CmdBody{
				Arg1: 0,
				Arg2: 0,
			},
			ErrorsData: ErrBody{ErrMessage: ""},
			Response: CmdResponse{
				Code:   "",
				Result: 0,
			},
		},
	}
}
