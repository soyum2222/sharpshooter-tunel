package main

import (
	"encoding/binary"
	"github.com/soyum2222/sharpshooter"
	"github.com/xtaci/smux"
	"io"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"sharpshooterTunnel/client/config"
	"sharpshooterTunnel/crypto"
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

	if config.CFG.FEC {
		sniper.OpenFec(10, 3)
	}

	//if config.CFG.Debug {
	//	sniper.OpenStaFlow()
	//}

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
			//
			//var total int64
			//var effective int64
			//for _, sniper := range sharpPool {
			//	t, e := sniper.FlowStatistics()
			//	total += t
			//	effective += e
			//}
			//
			//resp := struct {
			//	Total     int64 `json:"total"`
			//	Effective int64 `json:"effective"`
			//}{}
			//
			//resp.Total = total
			//resp.Effective = effective
			//
			//data, _ := json.Marshal(resp)
			//
			//_, _ = writer.Write(data)
		})

		go func() { _ = http.ListenAndServe(":"+strconv.Itoa(config.CFG.PPort), nil) }()
	}

	var i uint32

	for {

		local_conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		i++

		var session *smux.Session
		for {
			session = cPool[i%uint32(config.CFG.ConNum)]
			if session != nil {
				break
			}
			i++
		}

		remote_streem, err := session.OpenStream()
		if err != nil {

			cPool[i%uint32(config.CFG.ConNum)] = nil

			go func() {
			loop:
				cPool[i%uint32(config.CFG.ConNum)], err = createConn()
				if err != nil {
					log.Println(err)
					goto loop
				}
			}()

			continue
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
						remote_streem.Close()
						return
					}

					if n == 0 {
						continue
					}

					data, err := aes.Encrypt(b[:n])
					if err != nil {
						log.Println(err)
						local_conn.Close()
						remote_streem.Close()
						return
					}

					binary.BigEndian.PutUint32(head, uint32(len(data)))

					_, err = remote_streem.Write(append(head, data...))
					if err != nil {
						log.Println(err)
						local_conn.Close()
						remote_streem.Close()
						return
					}

				}

			}()

			// remote to local
			go func() {
				head := make([]byte, 4)

				for {

					_, err := io.ReadFull(remote_streem, head)
					if err != nil {
						log.Println(err)
						local_conn.Close()
						remote_streem.Close()
						return
					}

					var length uint32

					length = binary.BigEndian.Uint32(head)

					data := make([]byte, length)

					_, err = io.ReadFull(remote_streem, data)
					if err != nil {
						log.Println(err)
						local_conn.Close()
						remote_streem.Close()
						return
					}

					realdata, err := aes.Decrypt(data)
					if err != nil {
						log.Println(err)
						local_conn.Close()
						remote_streem.Close()
						return
					}

					_, err = local_conn.Write(realdata)
					if err != nil {
						log.Println(err)
						local_conn.Close()
						remote_streem.Close()
						return
					}
				}
			}()
		}()
	}
}
