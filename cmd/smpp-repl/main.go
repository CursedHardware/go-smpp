package main

import (
	"context"
	"crypto/tls"
	"encoding/hex"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/NiceLabs/go-smpp"
	"github.com/NiceLabs/go-smpp/coding"
	"github.com/NiceLabs/go-smpp/pdu"
	"github.com/abiosoft/ishell"
	"github.com/davecgh/go-spew/spew"
)

var conn *smpp.Conn

func main() {
	shell := ishell.New()
	shell.Println("Short Message Peer-to-Peer interactive shell")
	shell.AutoHelp(true)
	shell.SetHistoryPath(".smpp_repl_history")
	shell.AddCmd(&ishell.Cmd{Name: "connect", Help: "connect to server", Func: onConnect})
	shell.AddCmd(&ishell.Cmd{Name: "disconnect", Help: "disconnect", Func: onDisconnect})
	shell.AddCmd(&ishell.Cmd{Name: "send-message", Help: "send message", Func: onSendMessage})
	shell.AddCmd(&ishell.Cmd{Name: "send-ussd", Help: "send ussd", Func: onUSSDCommand})
	shell.AddCmd(&ishell.Cmd{Name: "query", Help: "query status", Func: onQueryCommand})
	go onHandler()
	shell.Run()
}

func onHandler() {
	var err error
	for conn != nil {
		for {
			packet := <-conn.PDU()
			spew.Dump(packet)
			switch p := packet.(type) {
			case *pdu.DeliverSM:
				resp := p.Resp().(*pdu.DeliverSMResp)
				messageId := make([]byte, 32)
				rand.Read(messageId)
				resp.MessageID = hex.EncodeToString(messageId)
				spew.Dump(resp)
				err = conn.Send(resp)
				if err != nil {
					fmt.Println(err)
				}
			case pdu.Responsable:
				err = conn.Send(p.Resp())
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func onConnect(c *ishell.Context) {
	c.ShowPrompt(false)
	defer c.ShowPrompt(true)
	if conn != nil {
		fmt.Println("connected")
		fmt.Println("use `disconnect` command, disconnect")
		return
	}
	var host, port, systemId, password, systemType string
	var enableTLS bool
	flags := makeFlags(func(flags *flag.FlagSet) {
		flags.StringVar(&host, "host", "", "Host")
		flags.StringVar(&port, "port", "2775", "Port")
		flags.StringVar(&systemId, "system-id", "", "System ID")
		flags.StringVar(&password, "password", "", "Password")
		flags.StringVar(&systemType, "system-type", "", "System Type")
		flags.BoolVar(&enableTLS, "tls", false, "Use TLS Mode")
	})
	if err := flags.Parse(c.Args); err != nil {
		fmt.Println("Error:", err.Error())
		return
	} else if flags.NFlag() < 3 {
		flags.Usage()
		return
	}
	address := net.JoinHostPort(host, port)
	var parent net.Conn
	var err error
	if enableTLS {
		parent, err = tls.Dial("tcp", address, &tls.Config{InsecureSkipVerify: true})
	} else {
		parent, err = net.Dial("tcp", address)
	}
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	conn = smpp.NewConn(context.Background(), parent)
	conn.WriteTimeout = time.Minute
	conn.ReadTimeout = time.Minute
	go conn.Watch()
	fmt.Printf("Connect %q successfully\n", address)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	resp, err := conn.Submit(ctx, &pdu.BindTransmitter{
		SystemID:   systemId,
		Password:   password,
		SystemType: systemType,
		Version:    pdu.SMPPVersion50,
	})
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	spew.Dump(resp)
	if status := pdu.ReadCommandStatus(resp); status == 0 {
		go conn.EnquireLink(time.Minute, time.Minute)
		fmt.Println("Bind successfully")
	}
}

func onDisconnect(c *ishell.Context) {
	c.ShowPrompt(false)
	defer c.ShowPrompt(true)
	if conn == nil {
		return
	}
	if err := conn.Close(); err != nil {
		fmt.Println(err)
	}
	conn = nil
}

func onSendMessage(c *ishell.Context) {
	c.ShowPrompt(false)
	defer c.ShowPrompt(true)
	if conn == nil {
		fmt.Println("You are not connected to the server")
		return
	}
	var source, dest, message string
	flags := makeFlags(func(flags *flag.FlagSet) {
		flags.StringVar(&source, "source", "", "Source address")
		flags.StringVar(&dest, "dest", "", "Destination address")
		flags.StringVar(&message, "message", "Test", "Message Content")
	})
	if err := flags.Parse(c.Args); err != nil {
		fmt.Println("Error:", err.Error())
		return
	} else if flags.NFlag() < 1 {
		flags.Usage()
		return
	}
	reference := uint16(rand.Intn(0xFFFF))
	parts, err := pdu.ComposeMultipartShortMessage(message, coding.BestCoding(message), reference)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	for _, message := range parts {
		packet := &pdu.SubmitSM{
			SourceAddr: pdu.Address{TON: 1, NPI: 1, No: source},
			DestAddr:   pdu.Address{TON: 1, NPI: 1, No: dest},
			ESMClass:   pdu.ESMClass{UDHIndicator: true},
			Message:    message,
		}
		spew.Dump(packet)
		resp, err := conn.Submit(context.Background(), packet)
		if err != nil {
			fmt.Println("Error:", err.Error())
			break
		}
		spew.Dump(resp)
	}
}

func onUSSDCommand(c *ishell.Context) {
	c.ShowPrompt(false)
	defer c.ShowPrompt(true)
	if conn == nil {
		fmt.Println("You are not connected to the server")
		return
	}
	var source, dest, message string
	flags := makeFlags(func(flags *flag.FlagSet) {
		flags.StringVar(&source, "source", "", "Source address")
		flags.StringVar(&dest, "dest", "", "Destination address")
		flags.StringVar(&message, "ussd", "*100#", "USSD command")
	})
	if err := flags.Parse(c.Args); err != nil {
		fmt.Println("Error:", err.Error())
		return
	} else if flags.NFlag() < 1 {
		flags.Usage()
		return
	}
	packet := &pdu.SubmitSM{
		ServiceType: "USSD",
		SourceAddr:  pdu.Address{TON: 1, NPI: 1, No: source},
		DestAddr:    pdu.Address{TON: 1, NPI: 1, No: dest},
		Tags:        pdu.Tags{0x5010: []byte{0x02}},
	}
	err := packet.Message.Compose(message)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	spew.Dump(packet)
	resp, err := conn.Submit(context.Background(), packet)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	spew.Dump(resp)
}

func onQueryCommand(c *ishell.Context) {
	if conn == nil {
		fmt.Println("You are not connected to the server")
		return
	}
	var id, source string
	var broadcast bool
	flags := makeFlags(func(flags *flag.FlagSet) {
		flags.StringVar(&id, "id", "", "Message ID")
		flags.StringVar(&source, "source", "", "Source address")
		flags.BoolVar(&broadcast, "broadcast", false, "Query Broadcast")
	})
	if err := flags.Parse(c.Args); err != nil {
		fmt.Println("Error:", err.Error())
		return
	} else if flags.NFlag() < 2 {
		flags.Usage()
		return
	}
	var packet pdu.Responsable
	if !broadcast {
		packet = &pdu.QuerySM{
			MessageID:  id,
			SourceAddr: pdu.Address{TON: 1, NPI: 1, No: source},
		}
	} else {
		packet = &pdu.QueryBroadcastSM{
			MessageID:  id,
			SourceAddr: pdu.Address{TON: 1, NPI: 1, No: source},
		}
	}
	spew.Dump(packet)
	resp, err := conn.Submit(context.Background(), packet)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	spew.Dump(resp)
}
