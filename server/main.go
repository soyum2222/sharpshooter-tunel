package main

import (
	"io"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"sharpshooter"
	"sharpshooterTunnel/server/config"
)

func main() {

	go http.ListenAndServe(":8888", nil)

	log.SetFlags(log.Llongfile | log.LstdFlags)
	l, err := sharpshooter.Listen(&net.UDPAddr{
		IP:   nil,
		Port: config.CFG.ListenPort,
		Zone: "",
	})

	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go func() {

			local_conn, err := net.Dial("tcp", config.CFG.LocalAddr)
			if err != nil {
				log.Println(err)
				conn.Close()
				return
			}

			go func() {
				_, err = io.Copy(conn, local_conn)
				if err != nil {
					log.Println(err)
				}

				conn.Close()
				local_conn.Close()
				return
			}()

			go func() {
				_, err = io.Copy(local_conn, conn)
				if err != nil {
					log.Println(err)
				}

				conn.Close()
				local_conn.Close()
				return
			}()

		}()
	}

}
