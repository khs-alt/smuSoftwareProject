package app

import "github.com/gorilla/mux"

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", mainHandler)
	r.HandleFunc("/get", requestHandler)
	//r.HandleFunc("/receviedVideo", receviedVideoHandler)
	r.HandleFunc("/postvideo/{id:[1-50]+}", serveSegmentVideoHandler)
	r.HandleFunc("/start", startJobHandler)
	r.HandleFunc("/video_summary", requestVideoSummaryHandler)
	r.HandleFunc("/check_video_summary", checkvideoSummaryHandler)
	r.HandleFunc("/api/uploadVideo", postVideoHandler)
	r.HandleFunc("/summerized_video", serveSummarizeVideoHandler)
	//r.HandleFunc("serveSegmentVideo", serveSegmentVideoHandler)

	return r
}
