package app

//Now: developed asynchronous communication method
//Next: developed long polling or socket programing method
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sfprojectBack/util"
	"strings"

	"github.com/gorilla/mux"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)
	data := map[string]interface{}{
		"message": "Hello Flask server!",
		"value":   48,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("json Marshal error", err)
		return
	}

	url := "http://localhost:8070/receive_data"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("HTTP Requset error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			fmt.Println("응답 디코딩 오류:", err)
			return
		}
		// 응답 데이터 출력
		fmt.Println("Python 서버 응답:", response)
	} else {
		fmt.Println("Python 서버 응답 실패. 상태 코드:", resp.StatusCode)
	}

}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://localhost:8070")
	if err != nil {
		http.Error(w, "Failed to connect to Flask server", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(body)
}

// 프론트에서 날라온 Original Video를 저장함.
func receviedVideoHandler(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)
	originalVideo, originalvideoHandler, err := r.FormFile("video")

	s := strings.Split(originalvideoHandler.Filename, ".")
	lastindex := len(s) - 1
	fileForm := s[lastindex]
	originalVideoPath := "./originalVideo/"
	originalVideoName := "originalVideo"
	originalFilePath := originalVideoPath + originalVideoName + "." + fileForm

	originalOutputFile, err := os.Create(originalFilePath)
	if err != nil {
		fmt.Println("orignial image fail")
		http.Error(w, "Unable to create the file for writing", http.StatusInternalServerError)
		return
	}
	defer originalOutputFile.Close()

	_, err = io.Copy(originalOutputFile, originalVideo) //originFile: multipart.File
	if err != nil {
		http.Error(w, "Unable to write the file", http.StatusInternalServerError)
		return
	}
}

func serveSummarizeVideoHandler(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)
	if r.Method == http.MethodGet {
		videoFilePath := fmt.Sprint("./output_final.mp4")
		//각 URL에 알맞는 비디오 지정
		http.ServeFile(w, r, videoFilePath)
	}
}

func serveSegmentVideoHandler(w http.ResponseWriter, r *http.Request) {
	// 비디오 파일을 읽어서 클라이언트로 전송
	videoID := mux.Vars(r)["id"]
	fmt.Println("serveVideosHandler : " + r.Method)
	fmt.Println(videoID)
	if r.Method == http.MethodGet {
		videoFilePath := fmt.Sprintf("./segmentVideo/segmentVideo%s.mp4", videoID)
		//각 URL에 알맞는 비디오 지정
		http.ServeFile(w, r, videoFilePath)
	}
}

func makeAutoSummarizeVideoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		util.EnableCors(&w)
		return
	}
	if r.Method == http.MethodGet {

		util.EditVideo("./originalVideo")
		res := "hello"
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(res))
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

//Go서버와 Flask서버의 통신 구현

var videoSummary = false

func startJobHandler(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)
	resp, err := http.Get("http://localhost:8070/start")
	if err != nil {
		http.Error(w, "Failed to connect to Flask server", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	// Flask 서버의 응답 출력
	fmt.Fprint(w, string(body))
}

// TODO 비디오 요약 만드는 코드 집어넣기
// 오래 걸리는 일 즉 Flask 작업
func requestVideoSummaryHandler(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)
	fmt.Fprint(w, "Job completed by Flask server")
	fmt.Println("Complete!")
	rootVideoPath := "./originalVideo"
	util.EditVideo(rootVideoPath)
	videoSummary = true
}

func checkvideoSummaryHandler(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)
	if videoSummary {
		fmt.Fprint(w, `{"status": "completed"}`)
	} else {
		fmt.Fprint(w, `{"status": "running"}`)
	}
	fmt.Println(videoSummary)
}

func postVideoHandler(w http.ResponseWriter, r *http.Request) {
	videoSummary = false
	if r.Method == http.MethodOptions {
		util.EnableCors(&w)

	}
	if r.Method == http.MethodPost {
		util.EnableCors(&w)
		file, fileHeader, err := r.FormFile("video")
		fileExtension := strings.ToLower(strings.Split(fileHeader.Filename, ".")[1])
		fileName := "original" + "." + fileExtension
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()
		uploadDir := "originalVideo/"
		dst, err := os.Create(uploadDir + fileName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// 업로드된 파일의 내용을 새 파일에 복사
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "File %s uploaded successfully.", fileName)
	}
}
