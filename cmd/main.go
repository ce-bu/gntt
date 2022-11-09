package main

import (
	"gntt/app/tcp_client"
	"gntt/app/tcp_server"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use: "gntt",
}

var optLogLvl string

var tcpServerCmd = &cobra.Command{
	Use: "tcp-server",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Flags().Visit(func(f *pflag.Flag) {

			if f.Name == "mtu-discover" {
				v, _ := cmd.Flags().GetInt("mtu-discover")
				tcpServerConfig.MtuDiscover.Set(v)
			}
		})
		tcp_server.New(tcpServerConfig).Run()
	},
}

var tcpServerConfig = &tcp_server.Config{
	Address:    "localhost",
	Port:       58822,
	MaxClients: 5,
	BufferSize: 65535,
}

var tcpClientCmd = &cobra.Command{
	Use: "tcp-client",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Flags().Visit(func(f *pflag.Flag) {

			switch f.Name {
			case "mtu-discover":
				v, _ := cmd.Flags().GetInt("mtu-discover")
				tcpClientConfig.MtuDiscover.Set(v)
			case "num-conn":
				v, _ := cmd.Flags().GetInt("num-conn")
				tcpClientConfig.NumConnections.Set(v)
			case "num-bytes":
				v, _ := cmd.Flags().GetInt64("num-bytes")
				tcpClientConfig.NumBytes.Set(v)
			}
		})
		tcp_client.New(tcpClientConfig).Run()
	},
}

var tcpClientConfig = &tcp_client.Config{
	Address:        "localhost",
	Port:           58822,
	MaxClients:     5,
	BufferSize:     65535,
	ConnTimeoutSec: 1,
}

func init() {
	tcpServerCmd.Flags().StringVarP(&tcpServerConfig.Address, "address", "a", tcpServerConfig.Address, "address")
	tcpServerCmd.Flags().IntVarP(&tcpServerConfig.Port, "port", "p", tcpServerConfig.Port, "port")
	tcpServerCmd.Flags().IntVarP(&tcpServerConfig.MaxClients, "max-clients", "c", tcpServerConfig.MaxClients, "max clients")
	tcpServerCmd.Flags().Int("mtu-discover", 1, "mtu discover")

	tcpClientCmd.Flags().StringVarP(&tcpClientConfig.Address, "address", "a", tcpClientConfig.Address, "address")
	tcpClientCmd.Flags().IntVarP(&tcpClientConfig.Port, "port", "p", tcpClientConfig.Port, "port")
	tcpClientCmd.Flags().IntVarP(&tcpClientConfig.MaxClients, "max-clients", "c", tcpClientConfig.MaxClients, "max clients")
	tcpClientCmd.Flags().IntVar(&tcpClientConfig.ConnTimeoutSec, "conn-timeout", tcpClientConfig.ConnTimeoutSec, "connection timeout")
	tcpClientCmd.Flags().Int("mtu-discover", 1, "mtu discover")
	tcpClientCmd.Flags().Int("num-conn", 0, "num connections")
	tcpClientCmd.Flags().Int64("num-bytes", 0, "num bytes")

	rootCmd.Flags().StringVar(&optLogLvl, "log", "trace", "log level")
	rootCmd.AddCommand(tcpServerCmd)
	rootCmd.AddCommand(tcpClientCmd)
}

func main() {
	logLvl, err := logrus.ParseLevel(optLogLvl)
	if err != nil {
		logrus.Panic("invalid log level")
	}
	logrus.SetLevel(logLvl)
	rootCmd.Execute()
	panic("done")

}
