package main

import (
	"encoding/binary"
	"github.com/soyum2222/sharpshooter"
	"io"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"sharpshooterTunnel/client/config"
	"sharpshooterTunnel/crypto"
)

func main() {

	go http.ListenAndServe(":9999", nil)

	log.SetFlags(log.Llongfile | log.LstdFlags)
	aes := crypto.AesCbc{
		Key:    config.CFG.Key,
		KenLen: 16,
	}

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
				log.Println(err)
				local_conn.Close()
				return
			}

			// local to remote
			go func() {

				b := make([]byte, 1<<10)
				head := make([]byte, 4)

				for {

					n, err := local_conn.Read(b)
					if err != nil {
						log.Println(err)
						local_conn.Close()
						remote_conn.Close()
						return
					}

					if n == 0 {
						continue
					}

					data, err := aes.Encrypt(b[:n])
					if err != nil {
						log.Println(err)
						local_conn.Close()
						remote_conn.Close()
						return
					}

					binary.BigEndian.PutUint32(head, uint32(len(data)))

					_, err = remote_conn.Write(append(head, data...))
					if err != nil {
						log.Println(err)
						local_conn.Close()
						remote_conn.Close()
						return
					}

				}

			}()

			// remote to local
			go func() {
				head := make([]byte, 4)

				for {

					_, err := io.ReadFull(remote_conn, head)
					if err != nil {
						log.Println(err)
						local_conn.Close()
						remote_conn.Close()
						return
					}

					var length uint32

					length = binary.BigEndian.Uint32(head)

					data := make([]byte, length)

					_, err = io.ReadFull(remote_conn, data)
					if err != nil {
						log.Println(err)
						local_conn.Close()
						remote_conn.Close()
						return
					}

					realdata, err := aes.Decrypt(data)
					if err != nil {
						log.Println(err)
						local_conn.Close()
						remote_conn.Close()
						return
					}

					_, err = local_conn.Write(realdata)
					if err != nil {
						log.Println(err)
						local_conn.Close()
						remote_conn.Close()
						return
					}
				}
			}()
		}()
	}
}
