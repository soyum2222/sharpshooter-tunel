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
	"runtime/debug"
	"sharpshooterTunnel/client/config"
	"sharpshooterTunnel/crypto"
	"sharpshooterTunnel/prof"
	"strconv"
	"time"
)

func createConn() (*smux.Session, error) {

	conn, err := sharpshooter.Dial(config.CFG.RemoteAddr)
	if err != nil {
		return nil, err
	}

	sniper := conn.(*sharpshooter.Sniper)

	sniper.SetSendWin(int32(config.CFG.SendWin))

	sniper.SetInterval(config.CFG.Interval)

	sniper.SetPackageSize(int64(config.CFG.MTU))

	if config.CFG.FEC {
		sniper.OpenFec(10, 3)
	}

	if config.CFG.Debug {
		sniper.OpenStaTraffic()
	}

	sharpPool = append(sharpPool, sniper)

	remote, err := smux.Client(conn, nil)
	if err != nil {
		return nil, err
	}

	return remote, nil
}

var cPool []*smux.Session

var sharpPool []*sharpshooter.Sniper

func main() {

	defer func() {
		err := recover()
		if err != nil {
			stack := debug.Stack()
			fmt.Println(string(stack))
			fmt.Println(err)
		}
	}()

	log.SetFlags(log.Llongfile | log.LstdFlags)
	aes := crypto.AesCbc{
		Key:    config.CFG.Key,
		KenLen: 16,
	}

	cPool = make([]*smux.Session, config.CFG.ConNum)

	l, err := net.Listen("tcp", config.CFG.LocalAddr)
	if err != nil {
		panic(err)
	}

	for i := 0; i < config.CFG.ConNum; i++ {

	loop:
		cPool[i], err = createConn()
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second)
			goto loop
		}
	}

	if config.CFG.Debug {
		http.HandleFunc("/statistics", func(writer http.ResponseWriter, request *http.Request) {

			sta := map[string]sharpshooter.Statistics{}
			for _, sniper := range sharpPool {
				s := sniper.TrafficStatistics()
				sta[sniper.LocalAddr().String()] = s
			}

			data, _ := json.Marshal(sta)

			_, _ = writer.Write(data)
		})

		go func() { _ = http.ListenAndServe(":"+strconv.Itoa(config.CFG.PPort), nil) }()
		prof.Monitor(time.Second * 30)
	}

	var i uint32

	for {

		local_conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		i++

	loop2:
		var session *smux.Session
		session = cPool[i%uint32(config.CFG.ConNum)]

		var remoteStream *smux.Stream
		if session != nil {
			remoteStream, err = session.OpenStream()
		}

		if err != nil || session == nil {
			index := i % uint32(config.CFG.ConNum)
			cPool[index], err = createConn()
			if err != nil {
				log.Println(err)
				local_conn.Close()
				continue
			}
			goto loop2
		}

		go func() {

			// local to remote
			go func() {

				b := make([]byte, 1<<10)
				head := make([]byte, 4)

				for {

					n, err := local_conn.Read(b)
					if err != nil {
						log.Println(err)
						local_conn.Close()
						remoteStream.Close()
						return
					}

					if n == 0 {
						continue
					}

					data, err := aes.Encrypt(b[:n])
					if err != nil {
						log.Println(err)
						local_conn.Close()
						remoteStream.Close()
						return
					}

					binary.BigEndian.PutUint32(head, uint32(len(data)))

					_, err = remoteStream.Write(append(head, data...))
					if err != nil {
						log.Println(err)
						local_conn.Close()
						remoteStream.Close()
						return
					}

				}

			}()

			// remote to local
			go func() {
				head := make([]byte, 4)

				for {

					_, err := io.ReadFull(remoteStream, head)
					if err != nil {
						log.Println(err)
						local_conn.Close()
						remoteStream.Close()
						return
					}

					var length uint32

					length = binary.BigEndian.Uint32(head)
					if length > 1<<11 {
						continue
					}

					data := make([]byte, length)

					_, err = io.ReadFull(remoteStream, data)
					if err != nil {
						log.Println(err)
						local_conn.Close()
						remoteStream.Close()
						return
					}

					realdata, err := aes.Decrypt(data)
					if err != nil {
						log.Println(err)
						local_conn.Close()
						remoteStream.Close()
						return
					}

					_, err = local_conn.Write(realdata)
					if err != nil {
						log.Println(err)
						local_conn.Close()
						remoteStream.Close()
						return
					}
				}
			}()
		}()
	}
}
