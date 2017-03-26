package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"simple-blockchain/blockchain"
)

var bc = blockchain.GetBlockchain()
var peers = map[string]string{}

func main() {
	portFlag := flag.String("port", "8080", "port to run server")
	flag.Parse()
	port := ":" + *portFlag
	http.HandleFunc("/", home)
	http.HandleFunc("/create", createBlock)
	http.HandleFunc("/blocks", blocks)
	http.HandleFunc("/addBlock", addBlock)
	http.HandleFunc("/peer", peer)
	log.Println("Listening on port", port)
	log.Println(http.ListenAndServe(port, nil))
}

func home(w http.ResponseWriter, r *http.Request) {
	blockData := []string{}
	for _, b := range bc.Blocks() {
		blockData = append(blockData, b.Data)
	}
	peerURLs := make([]string, len(peers))

	i := 0
	for k := range peers {
		peerURLs[i] = k
		i++
	}

	t, _ := template.ParseFiles("home.gtpl")
	t.Execute(w, map[string][]string{
		"Blocks": blockData,
		"Peers":  peerURLs,
	})
}

func createBlock(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	data := r.Form.Get("BlockData")
	if data == "" {
		http.Error(w, "must supply BlockData", http.StatusBadRequest)
		return
	}
	bc.AddNewBlock(data)
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(bc.Head)
	for _, peer := range peers {
		resp, err := http.Post(peer+"/addBlock", "application/json", buf)
		if err != nil {
			log.Println("unable to POST new block to peer", peer, bc.Head, resp, err)
		}
	}
	http.Redirect(w, r, "/home", 302)
}

func blocks(w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(bc.Blocks())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(json))
}

func addBlock(w http.ResponseWriter, r *http.Request) {
	var b blockchain.Block
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		log.Println("unable to decode json into block", err)
	}
	if !bc.IsValidBlock(b, bc.Head) {
		log.Println("need to sync with peer", r.Referer())
		synchronizeWithPeer(r.Referer())
	} else {
		log.Println("adding new valid block from peer", r.Referer())
		bc.InsertBlock(b)
	}
}

func peer(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	peerURL := r.Form.Get("PeerURL")
	if peerURL == "" {
		http.Error(w, "must supply PeerURL", http.StatusBadRequest)
		return
	}
	peers[peerURL] = peerURL
	synchronizeWithPeer(peerURL)
	http.Redirect(w, r, "/home", 302)
}

func synchronizeWithPeer(peer string) {
	peerResponse, err := http.Get(peer + "/blocks")
	if err != nil {
		log.Println("unable to fetch blocks from peer", peer, err)
		return
	}
	body, err := ioutil.ReadAll(peerResponse.Body)
	if err != nil {
		log.Println("unable to read body response from peer", peer, err)
	}
	peerBc := blockchain.BlockchainFromJSON(body)
	if peerBc.Head.Index > bc.Head.Index {
		if peerBc.Head.Index == bc.Head.Index+1 &&
			bc.IsValidBlock(peerBc.Head, bc.Head) {
			log.Println("added one new block from peer", peer)
			bc.InsertBlock(peerBc.Head)
		} else if bc.ShouldReplaceWithChain(peerBc) {
			log.Println("replacing blockchain with that of peer", peer)
			bc = peerBc
		}
	}
}
