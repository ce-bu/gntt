//go:build linux
// +build linux

package tcp_server

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
