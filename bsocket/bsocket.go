package bsocket

import (
	"log"
	"time"

	"github.com/chaosknight/gowebsocket"
)

const timestep = 20 * time.Second

type Bsocket struct {
	socket    gowebsocket.Socket
	timer     *time.Timer
	isconnect bool
	isTimeOut bool
}

func NewBsocket(wss string) *Bsocket {
	return &Bsocket{
		socket: gowebsocket.New(wss),
	}
}

func (sys *Bsocket) Connect() {
	sys.socket.Connect()
}

func (sys *Bsocket) Init(funs func(message string, socket gowebsocket.Socket), onconnected func(socket gowebsocket.Socket)) {
	sys.socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server")
		sys.isconnect = true
		sys.resettimer()
		if onconnected != nil {
			onconnected(socket)
		}
	}
	sys.socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Recieved connect error ", err)
	}
	sys.socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from server ")
		sys.isconnect = false
	}

	sys.socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {

		if message == "pong" {
			return
		}
		if funs != nil {
			funs(message, socket)
		}

	}
}

func (sys *Bsocket) SendText(msg string) {
	if sys.isconnect {
		sys.socket.SendText(msg)
	}
}

func (sys *Bsocket) resettimer() {
	if sys.timer == nil {
		sys.timer = time.NewTimer(timestep)
		go func() {
			for !sys.isTimeOut {
				<-sys.timer.C
				sys.timer.Reset(timestep)
				if !sys.isconnect {
					sys.socket.Connect()
				} else {
					go sendping(sys)
				}

			}

		}()
	} else {
		sys.timer.Reset(timestep)
	}
}

func sendping(sys *Bsocket) {
	defer func() {
		if x := recover(); x != nil {
			log.Println("send ping error", x)
		}
	}()
	sys.socket.SendText("ping")
}
