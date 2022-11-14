package tcp_client

import (
	"fmt"
	"gntt/pkg/gntt_math"
	"gntt/pkg/gntt_optional"
	"gntt/pkg/gntt_utils"
	"gntt/pkg/gntt_worker"
	"net"
	"os"
	"os/signal"
	"sync"
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
	ConnTimeSec    gntt_optional.Optional[int]
	ConnTimeoutSec int
}

type App struct {
	config       *Config
	actvConn     sync.Map
	jobsFinished chan bool
	bytesSent    uint64
	sampler      *gntt_utils.RateSampler
}

func New(c *Config) *App {
	return &App{
		config: c,
		sampler: gntt_utils.NewSampler(100*time.Millisecond, func(rate float64) {
			fmt.Printf("%.2f Gb/s\n", rate/(1000.0*1000.0*1000.0))
		}),
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

	configureConn(app, conn)

	var buf []byte = make([]byte, app.config.BufferSize)

	unlimited := !app.config.NumBytes.HasValue()
	var total int64 = int64(app.config.BufferSize)
	if !unlimited {
		total = app.config.NumBytes.Get()
	}

	app.sampler.Start()

	connStartTime := time.Now()

	for unlimited || total > 0 {

		n, err := (*conn).Write(buf[0:gntt_math.Max(total, int64(app.config.BufferSize))])
		if err != nil {
			log.Tracef("conn %s write error=%s", (*conn).RemoteAddr().String(), err.Error())
			break
		}

		app.sampler.AddSample(uint64(n))

		if !unlimited {
			total = total - int64(n)
		}

		if app.config.ConnTimeSec.HasValue() {
			delta := time.Now().Sub(connStartTime).Seconds()
			if delta > float64(app.config.ConnTimeSec.Get()) {
				log.Tracef("conn %s expired %f s %d", (*conn).RemoteAddr().String(), delta, app.config.ConnTimeSec.Get())
				break
			}
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

	select {
	case <-sigs:
	case <-app.jobsFinished:
	}

	app.sampler.Stop()
	endWork <- true
	<-endWork

}
