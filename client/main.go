package main

import (
	"io"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"sharpshooter"
	"sharpshooterTunnel/client/config"
)

func main() {

	go http.ListenAndServe(":9999", nil)

	log.SetFlags(log.Llongfile | log.LstdFlags)

	l, err := net.Listen("tcp", config.CFG.LocalAddr)
	if err != nil {
		panic(err)
	}
	for {

		local_conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go func() {

			remote_conn, err := sharpshooter.Dial(config.CFG.RemoteAddr)
			if err != nil {
				local_conn.Close()
				return
			}

			go func() {

				_, err := io.Copy(remote_conn, local_conn)
				if err != nil {
					log.Println(err)
				}
				_ = local_conn.Close()
				remote_conn.Close()
				return

			}()

			go func() {
				_, err := io.Copy(local_conn, remote_conn)
				if err != nil {
					log.Println(err)
				}
				_ = local_conn.Close()
				remote_conn.Close()
				return

			}()

		}()

	}
}
