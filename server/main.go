package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/soyum2222/sharpshooter"
	"github.com/xtaci/smux"
	"io"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"sharpshooterTunnel/crypto"
	"sharpshooterTunnel/prof"
	"sharpshooterTunnel/server/config"
	"strconv"
	"sync"
	"time"
)

var currentPool sync.Map

func main() {

	if config.CFG.Debug {
		http.HandleFunc("/statistics", func(writer http.ResponseWriter, request *http.Request) {

			sta := map[string]sharpshooter.Statistics{}
			currentPool.Range(func(key, value interface{}) bool {
				sta[value.(*sharpshooter.Sniper).RemoteAddr().String()] = value.(*sharpshooter.Sniper).TrafficStatistics()
				return true
			})

			data, _ := json.Marshal(sta)
			_, _ = writer.Write(data)
		})

		go func() { fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", config.CFG.PPort), nil)) }()
		prof.Monitor(time.Second * 30)
	}

	log.SetFlags(log.Llongfile | log.LstdFlags)

	l, err := sharpshooter.Listen(":" + strconv.Itoa(config.CFG.ListenPort))
	if err != nil {
		panic(err)
	}

	aes := crypto.AesCbc{
		Key:    config.CFG.Key,
		KenLen: 16,
	}

	for {

		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		rawconn := conn.(*sharpshooter.Sniper)

		rawconn.SetSendWin(int32(config.CFG.SendWin))

		rawconn.SetInterval(config.CFG.Interval)

		rawconn.SetPackageSize(int64(config.CFG.MTU))

		currentPool.Store(conn.RemoteAddr(), rawconn)

		if config.CFG.Debug {
			rawconn.OpenStaTraffic()
		}

		if config.CFG.FEC {
			rawconn.OpenFec(10, 3)
		}

		serconn, err := smux.Server(rawconn, nil)
		if err != nil {
			log.Println(err)
			return
		}

		go func() {
			defer func() {
				currentPool.Delete(conn.RemoteAddr())
				_ = conn.Close()
				_ = serconn.Close()
			}()

			for {

				conn, err := serconn.AcceptStream()
				if err != nil {
					log.Println(err)
					return
				}

				go func() {

					local_conn, err := net.Dial("tcp", config.CFG.LocalAddr)
					if err != nil {
						log.Println(err)
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
								conn.Close()
								return
							}

							if n == 0 {
								continue
							}

							data, err := aes.Encrypt(b[:n])
							if err != nil {
								log.Println(err)
								local_conn.Close()
								conn.Close()
								return
							}

							binary.BigEndian.PutUint32(head, uint32(len(data)))

							_, err = conn.Write(append(head, data...))
							if err != nil {
								log.Println(err)
								local_conn.Close()
								conn.Close()
								return
							}
						}
					}()

					// remote to local
					go func() {

						head := make([]byte, 4)
						for {

							_, err := io.ReadFull(conn, head)
							if err != nil {
								log.Println(err)
								local_conn.Close()
								conn.Close()
								return
							}

							var length uint32

							length = binary.BigEndian.Uint32(head)

							data := make([]byte, length)

							_, err = io.ReadFull(conn, data)
							if err != nil {
								log.Println(err)
								local_conn.Close()
								conn.Close()
								return
							}

							realdata, err := aes.Decrypt(data)
							if err != nil {
								log.Println(err)
								local_conn.Close()
								conn.Close()
								return
							}

							_, err = local_conn.Write(realdata)
							if err != nil {
								log.Println(err)
								local_conn.Close()
								conn.Close()
								return
							}
						}
					}()
				}()
			}
		}()
	}
}
