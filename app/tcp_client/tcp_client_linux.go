//go:build linux
// +build linux

package tcp_client

import (
	"net"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func configureConn(app *App, conn *net.Conn) {
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

}

func configureRawConn(app *App, fd uintptr) {
	if app.config.TcpFastOpen.HasValue() {
		err := syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, 30, app.config.TcpFastOpen.Get())
		if err != nil {
			log.Errorf("cannot set tcp_fast_open error=%s", err.Error())
			return
		}
	}

}
