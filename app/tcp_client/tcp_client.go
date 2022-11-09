package tcp_client

import (
	"fmt"
	"gntt/pkg/gntt_math"
	"gntt/pkg/gntt_optional"
	"gntt/pkg/gntt_worker"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"

	"syscall"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Address        string
	Port           int
	MaxClients     int
	BufferSize     int
	MtuDiscover    gntt_optional.Optional[int]
	NumConnections gntt_optional.Optional[int]
	NumBytes       gntt_optional.Optional[int64]
	ConnTimeoutSec int
}

type App struct {
	config       *Config
	actvConn     sync.Map
	jobsFinished chan bool
	bytesSent    uint64
}

func New(c *Config) *App {
	return &App{
		config: c,
	}
}

func (app *App) NumJobs() gntt_optional.Optional[int] {
	return app.config.NumConnections
}

func (app *App) MaxClients() int {
	return app.config.MaxClients
}

func (app *App) handleConn(conn *net.Conn) {

	defer (*conn).Close()
	tcpc := (*conn).(*net.TCPConn)
	fd, err := tcpc.File()
	if err != nil {
		log.Errorf("cannot get socket descriptor error=%s", err.Error())
		return
	}
	if app.config.MtuDiscover.HasValue() {
		err = syscall.SetsockoptInt(int(fd.Fd()), syscall.IPPROTO_IP, syscall.IP_MTU_DISCOVER, app.config.MtuDiscover.Get())
		if err != nil {
			log.Errorf("cannot set IP_MTU_DISCOVER error=%s", err.Error())
			return
		}
	}

	var buf []byte = make([]byte, app.config.BufferSize)

	unlimited := !app.config.NumBytes.HasValue()
	var total int64 = int64(app.config.BufferSize)
	if !unlimited {
		total = app.config.NumBytes.Get()
	}

	for unlimited || total > 0 {

		n, err := (*conn).Write(buf[0:gntt_math.Max(total, int64(app.config.BufferSize))])
		if err != nil {
			log.Tracef("conn %s write error=%s", (*conn).RemoteAddr().String(), err.Error())
			break
		}

		atomic.AddUint64(&app.bytesSent, uint64(n))

		if !unlimited {
			total = total - int64(n)
		}
	}
}

func (app *App) Perform() {
	address := fmt.Sprintf("%s:%d", app.config.Address, app.config.Port)
	d := net.Dialer{
		Timeout: time.Duration(app.config.ConnTimeoutSec) * time.Second,
	}
	newc, err := d.Dial("tcp", address)
	if err != nil {
		log.Errorf("connection error=%s", err.Error())
	} else {
		app.actvConn.Store(newc, 1)
		app.handleConn(&newc)
		app.actvConn.Delete(newc)
	}
}

func (app *App) CancelAll() {
	app.actvConn.Range(func(conn any, value any) bool {
		conn.(net.Conn).Close()
		return true
	})

}

func (app *App) JobsFinished() {
	app.jobsFinished <- true
}

func (app *App) Run() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	app.jobsFinished = make(chan bool)

	endWork := gntt_worker.ClientWorker(app)

	sampler := func() chan bool {
		stop := make(chan bool)
		go func() {
			pt := time.Now()
			pb := atomic.LoadUint64(&app.bytesSent)
			for {
				select {
				case <-stop:
					goto end
				case <-time.After(100 * time.Millisecond):
					ct := time.Now()
					cb := atomic.LoadUint64(&app.bytesSent)
					rate := float64(cb-pb) / float64(ct.Sub(pt).Seconds())
					pt = ct
					pb = cb
					fmt.Printf("%.2fGb\n", rate/(100.0*1000.0*1000.0))
				}
			}
		end:
			stop <- true
		}()
		return stop
	}()

	select {
	case <-sigs:
	case <-app.jobsFinished:
	}
	sampler <- true
	<-sampler
	endWork <- true
	<-endWork

}
