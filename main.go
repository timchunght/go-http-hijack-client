package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("usage:", os.Args[0], "<url> <token>")

	}
	rawURL := os.Args[1]
	token := os.Args[2]
	url, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalln(rawURL, "is not a valid URL", err)
	}

	req, err := http.NewRequest("CONNECT", rawURL, nil)
	if err != nil {
		log.Fatalln("Fail to create request CONNECT", rawURL)
	}

	req.SetBasicAuth("", token)

	dial, err := net.Dial("tcp", url.Host)
	if err != nil {
		log.Fatalln("Fail to connect to", url.Host, err)
	}

	var conn *httputil.ClientConn
	if url.Scheme == "https" {
		host := strings.Split(url.Host, ":")[0]
		tls_conn := tls.Client(dial, &tls.Config{ServerName: host})
		conn = httputil.NewClientConn(tls_conn, nil)
	} else if url.Scheme == "http" {
		conn = httputil.NewClientConn(dial, nil)
	} else {
		log.Println("Scheme format should be 'http' or 'https' not", url.Scheme)
	}

	_, err = conn.Do(req)
	if err != httputil.ErrPersistEOF && err != nil {
		log.Fatalln("Fail to execute http request", err)
	}

	log.Println("Hijack connection")
	connection, _ := conn.Hijack()

	fmt.Println("Pipe stdin to socket, and socket to stdout")
	go io.Copy(connection, os.Stdin)
	_, err = io.Copy(&ServerWriter{os.Stdout}, connection)
	if err != nil {
		log.Println("Erreur when copying socket to stdout:", err)
	}
	fmt.Println("End of hijacking")
}

type ServerWriter struct{ io.Writer }

func (w *ServerWriter) Write(b []byte) (int, error) {
	w.Writer.Write([]byte("Server: "))
	return w.Writer.Write(b)
}
