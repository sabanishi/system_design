package main

import (
	"fmt"
	"log"
	"net" // standard network package
	"strings"
)

func main() {
	// config
	port := 8000
	protocol := "tcp"

	// resolve TCP address
	addr, err := net.ResolveTCPAddr(protocol, fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalln(err)
	}

	// get TCP socket
	socket, err := net.ListenTCP(protocol, addr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Listen: ", socket.Addr().String())

	// keep listening
	for {
		// wait for connection
		conn, err := socket.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("Connected by ", conn.RemoteAddr().String())

		// yield connection to concurrent process
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// close connection when this function ends
	defer conn.Close()

	buf := make([]byte, 1024)
	conn.Read(buf)

	//HTTPリクエストを行毎に分割
	splits := strings.Split(string(buf), "\n")
	if len(splits) < 1 {
		return
	}

	//アクセスURLが格納された1番目の要素を取り出す
	resource := splits[0]
	resources := strings.Split(resource, " ")

	if len(resources) < 2 {
		buf = createNotFoundRespose()
	} else {
		switch resources[1] {
		case "/hello":
			buf = createHttpResponse([]byte("Hello world."))
		case "/bye":
			buf = createHttpResponse([]byte("Good bye."))
		case "/hello.jp":
			buf = createHttpResponse([]byte("こんにちは"))
		default:
			buf = createNotFoundRespose()
		}
	}

	conn.Write(buf)
}

/*引数のbyte列をBodyとするHTTPリクエストを生成する関数*/
func createHttpResponse(body []byte) []byte {
	buf := []byte(fmt.Sprintf(
		"HTTP/1.1 200 OK\n"+
			"Content-Type: text/plain; charset=utf-8\n"+
			"Content-Length: %d\n"+
			"\n"+
			"%s\n",
		len(body),
		string(body)))
	return buf
}

/*ページが見つからないことを示すHTTPリクエストを生成する関数*/
func createNotFoundRespose() []byte {
	buf := []byte(fmt.Sprintf(
		"HTTP/1.1 404 Not Found\n"+
			"%d\n"+
			"\n",
		0))
	return buf
}
