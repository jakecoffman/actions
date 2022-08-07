package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func main() {
	if ip := os.Getenv("PROXY_ADDR"); ip != "" {
		url := fmt.Sprintf("http://%v:8000", ip)
		buf := bytes.NewBufferString("Hello")
		resp, err := http.Post(url, "text/html", buf)
		if err != nil {
			log.Fatal(err)
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(data))
		return
	}

	var err error
	for i := 0; i < 60*5; i++ {
		err = tailscaleCall()
		if err == nil {
			return
		}
		log.Println(err)
		time.Sleep(time.Second)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func tailscaleCall() error {
	ip := GetOtherIP()
	log.Println("Other IP:", ip)
	url := fmt.Sprintf("http://%v:8000", ip)
	buf := bytes.NewBufferString("Hello")
	resp, err := http.Post(url, "text/html", buf)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println(string(data))
	return nil
}

func GetOtherIP() string {
	buf := &bytes.Buffer{}

	cmd := exec.Command("tailscale", "status", "--json")
	cmd.Stdout = buf
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	var v tail
	if err := json.NewDecoder(buf).Decode(&v); err != nil {
		log.Fatal(err)
	}

	var peers []Peer

	for _, peer := range v.Peer {
		if strings.HasPrefix(peer.HostName, "github") {
			peers = append(peers, peer)
		}
	}

	sort.Slice(peers, func(i, j int) bool {
		return peers[i].Created.After(peers[j].Created)
	})

	return peers[0].TailscaleIPs[0]
}

type tail struct {
	Peer map[string]Peer
}

type Peer struct {
	HostName         string
	TailscaleIPs     []string
	Active           bool
	Online           bool
	RxBytes, TxBytes int
	Created          time.Time
}
