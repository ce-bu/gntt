//go:build windows
// +build windows

package tcp_client

import "net"

func configureConn(app *App, conn *net.Conn) {

}

func configureRawConn(app *App, fd uintptr) {
}
