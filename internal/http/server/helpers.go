package server

import (
	"path"
	"strings"
)

func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

// func writeJSON(w http.ResponseWriter, value interface{}) {
// 	data, err := json.Marshal(&value)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		log.Println("failed to marshal:", err)
// 		fmt.Fprintf(w, "failed to marshal: %v", err)
// 		return
// 	}

// 	w.Write(data)
// }
