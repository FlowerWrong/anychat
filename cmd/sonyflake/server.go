package main

import (
	"encoding/json"
	"net/http"

	"github.com/sony/sonyflake"
)

// doc https://chai2010.cn/advanced-go-programming-book/ch6-cloud/ch6-01-dist-id.html

var sf *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings
	// st.MachineID = awsutil.AmazonEC2MachineID
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sonyflake not created")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	id, err := sf.NextID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(sonyflake.Decompose(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header()["Content-Type"] = []string{"application/json; charset=utf-8"}
	w.Write(body)
}

// curl 127.0.0.1:8090
func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8090", nil)
}
