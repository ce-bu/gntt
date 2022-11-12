package tcp_server

import (
	"fmt"
	"gntt/pkg/gntt_optional"
	"gntt/pkg/gntt_worker"
	"io"
	"net"
	"os"
	"os/signal"

	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Address     string
	Port        int
	MaxClients  int
	BufferSize  int
	MtuDiscover gntt_optional.Optional[int]
}

type App struct {
	config   *Config
	listener net.Listener
	actvConn sync.Map
}

func New(c *Config) *App {
	return &App{
		config: c,
	}
}

func (app *App) handleConn(conn *net.Conn) {

	defer (*conn).Close()
	configureConn(app, conn)

	var buf []byte = make([]byte, app.config.BufferSize)
	for {
		_, err := (*conn).Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Tracef("conn %s read error=%s", (*conn).RemoteAddr().String(), err.Error())
			}
			break
		}
	}
}

func (app *App) Setup() error {
	address := fmt.Sprintf("%s:%d", app.config.Address, app.config.Port)
	log.Infof("server listen at %s", address)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("accept error=%s", err.Error())
	} else {
		app.listener = listener
	}
	return err
}

func (app *App) Teardown() {
	app.listener.Close()
}

func (app *App) MaxConcurrentJobs() int {
	return app.config.MaxClients
}

func (app *App) Accept() (net.Conn, error) {
	newc, err := app.listener.Accept()
	if err != nil {
		log.Tracef("accept error=%s", err.Error())
	} else {
		log.Infof("new connection remote=%s", newc.RemoteAddr().String())
	}
	return newc, err
}

func (app *App) Perform(newc net.Conn) {
	app.actvConn.Store(newc, 0)
	app.handleConn(&newc)
	app.actvConn.Delete(newc)

}

func (app *App) CancelAll() {
	app.actvConn.Range(func(conn any, value any) bool {
		conn.(net.Conn).Close()
		return true
	})
}

func (app *App) Run() {
	endServer, err := gntt_worker.ServerWorker[net.Conn](app)
	if err != nil {
		return
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	endServer <- true
}
