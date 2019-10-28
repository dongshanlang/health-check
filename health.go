// There is a standard HTTP interface to report information. Adding
// the following line will install handlers under the /status
// usage:
//
//    import _ "net/http/pprof"
//

//Another case, when there is not existing a http server, the usage should be like this
//
//    import _ "net/http/pprof"
//    health.Start(":8080")
//
package health

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

var h = new(Health)

func init() {
	http.HandleFunc("/status", NewJSONHandlerFunc(*h, nil))
}

type Health struct {
	//TODO: get a logger from outside.
	//Logger     log.Logger
	addr string
}

// State is a struct that contains the results of the latest
// run of a particular check.
type State struct {
	// App information
	AppInfo AppInformation `json:"app_info"`
	//MemStatus
	MemStatus MemoryStatus `json:"mem_status"`
	// Err is the error returned from a failed health check
	Err string `json:"error,omitempty"`
	// ServerTime is the time of the last health check
	ServerTime time.Time `json:"server_time"`
}

//MemoryStatus is a struct for the runtime memory information
type MemoryStatus struct {
	TotalAlloc   uint64 `json:"TotalAlloc"`
	Alloc        uint64 `json:"Alloc"`
	Mallocs      uint64 `json:"Mallocs"`
	HeapAlloc    uint64 `json:"HeapAlloc"`
	LastGc       uint64 `json:"LastGC"`
	NextGc       uint64 `json:"NextGC"`
	PauseTotalNs uint64 `json:"PauseTotalNs"`
	NumGC        uint32 `json:"NumGC"`
}

//AppInformation is a struct used to telling  the service build info
type AppInformation struct {
	Version  string `json:"version"`
	Debug    bool   `json:"debug"`
	Build    string `json:"build"`
	HostName string `json:"host_name"`
}

// State will return a map of all current healthcheck states,
// The returned structs can be used for figuring out additional analytics or
// used for building your own status handler (todo)
func (h *Health) State() (State, error) {
	return h.safeGetStates(), nil
}

func (h *Health) Start(addr string) error {
	if len(addr) == 0 {
		h.addr = ":8080"
	} else {
		h.addr = addr
	}

	go http.ListenAndServe(h.addr, nil)
	return nil
}

//Start will start a http server listen at addr, the addr should be like "192.168.3.12:8090" or ":8090".
func Start(addr string) error {
	return h.Start(addr)
}

// get all states
func (h *Health) safeGetStates() State {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	return State{
		ServerTime: time.Now(),
		//AppInfo:    nil,
		MemStatus: MemoryStatus{
			TotalAlloc:   m.TotalAlloc,
			Alloc:        m.Alloc,
			Mallocs:      m.Mallocs,
			HeapAlloc:    m.HeapAlloc,
			LastGc:       m.LastGC,
			NextGc:       m.NextGC,
			PauseTotalNs: m.PauseTotalNs,
			NumGC:        m.NumGC,
		},
		Err: "",
	}
}

func (h *Health) setAddr(addr string) error {
	h.addr = addr
	return nil
}

type jsonStatus struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// NewJSONHandlerFunc will return an `http.HandlerFunc` that will marshal and
// write the contents of `h.State()` to `rw` and set status code to
//  `http.StatusOK` .
// todo: It also accepts a set of optional custom fields to be added to the final JSON body
func NewJSONHandlerFunc(h Health, custom map[string]interface{}) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		sts, err := h.State()
		if err != nil {
			writeJSONStatus(rw, "error", fmt.Sprintf("Unable to fetch states: %v", err), http.StatusOK)
			return
		}
		statusCode := http.StatusOK
		data, err := json.Marshal(sts)
		if err != nil {
			writeJSONStatus(rw, "error", fmt.Sprintf("Failed to marshal state data: %v", err), http.StatusOK)
			return
		}

		writeJSONResponse(rw, statusCode, data)
	})
}

func writeJSONStatus(rw http.ResponseWriter, status, message string, statusCode int) {
	jsonData, _ := json.Marshal(&jsonStatus{
		Message: message,
		Status:  status,
	})

	writeJSONResponse(rw, statusCode, jsonData)
}

func writeJSONResponse(rw http.ResponseWriter, statusCode int, content []byte) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
	rw.WriteHeader(statusCode)
	rw.Write(content)
}
