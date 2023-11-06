package util

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ExtractNonOverlappingTimes(startTimes []string, endTimes []string) ([]string, []string) {
	n := len(startTimes)
	var newStartTimes []string
	var newEndTimes []string
	newStartTimes = append(newStartTimes, startTimes[0])
	for i := 0; i < n-1; i++ {
		if startTimes[i+1] != endTimes[i] {
			newStartTimes = append(newStartTimes, startTimes[i+1])
			newEndTimes = append(newEndTimes, endTimes[i])
		}
	}
	newEndTimes = append(newEndTimes, endTimes[n-1])
	return newStartTimes, newEndTimes
}

// 원본 비디오가 들어오면 segment Video를 만듭니다
// 3초 단위로 비디오를 만들며 각 비디오는 video1, video2..이런식으로 만들어집니다.
// 이후 생성된 비디오는 "./segmentVideo"에 일괄 저장됩니다.
func SegmentVideo(rootPath string) error {

	inputVideoFile, err := ReadDirFile(rootPath)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	inputVideoFile = rootPath + "/" + inputVideoFile

	// 세그먼트 저장 디렉토리
	outputSegmentDirectory := "./segmentVideo"

	// 세그먼트의 시간 간격 (예: 10초)
	segmentDuration := "3"

	// 세그먼트 디렉토리가 없다면 생성
	if _, err := os.Stat(outputSegmentDirectory); os.IsNotExist(err) {
		os.Mkdir(outputSegmentDirectory, os.ModeDir)
	}

	// ffmpeg 명령 실행
	ffmpegPath := "ffmpeg"
	cmd := exec.Command(ffmpegPath,
		"-i", inputVideoFile,
		"-c:v", "copy",
		"-c:a", "copy",
		"-segment_time", segmentDuration,
		"-f", "segment",
		"-reset_timestamps", "1",
		"-map", "0",
		outputSegmentDirectory+"/segment%d.mp4",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println("비디오 세그먼트 생성이 완료되었습니다.")
	return nil
}

func EditVideo(inputFilePath string) {

	// 입력 동영상 파일 및 출력 디렉토리 설정
	inputVideoFile, err := ReadDirFile(inputFilePath)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	inputVideoFile = inputFilePath + "/" + inputVideoFile

	outputDirectory := "output"
	DeleteFilesInFolder(outputDirectory)

	// // 시간대를 받아오기 (예: 00:01:10~00:01:20, 00:02:10~00:02:30)
	// timelineFilePath := "..\\sfprojectAI\\summerized_subtitle\\summerized_subtitle_1.csv"
	// fileData, err := ioutil.ReadFile(timelineFilePath)
	// if err != nil {
	// 	fmt.Println(err)
	// 	panic(err)
	// }

	// timelines := string(fileData)
	// splitTimes := strings.Split(timelines, ",")
	// fmt.Println(splitTimes)

	file, err := os.Open("..\\sfprojectAI\\summerized_subtitle\\summerized_subtitle.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// CSV 리더 생성
	reader := csv.NewReader(file)
	splitTimes, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	//convertMillisecondsToTime(splitTimes)
	// 각 시간대로 동영상을 자르고 저장
	var startTimes []string
	var endTimes []string
	fmt.Println(splitTimes)
	for _, timeRange := range splitTimes[0] {
		splitTime := strings.Split(timeRange, "-")
		startTime := splitTime[0]
		endTime := splitTime[1]
		fmt.Println(startTime, endTime)
		startTime = strings.TrimLeft(startTime, ` '`)
		endTime = strings.TrimRight(endTime, `'`)
		fmt.Println(startTime, endTime)
		startTimes = append(startTimes, startTime)
		endTimes = append(endTimes, endTime)
	}

	startTimes, endTimes = ExtractNonOverlappingTimes(startTimes, endTimes)
	num := len(startTimes)
	fmt.Print("startTimes, endTimes")
	fmt.Println(startTimes, endTimes)
	for i := 0; i < num; i++ {
		fileName := fmt.Sprintf("output_%02d.mp4", i+1)

		outputFile := filepath.Join(outputDirectory, fileName)
		splitCmd := exec.Command("ffmpeg",
			"-i", inputVideoFile,
			"-ss", startTimes[i],
			"-to", endTimes[i],
			"-c:v", "libx264",
			"-c:a", "aac",
			"-force_key_frames", startTimes[i],
			outputFile)
		err := splitCmd.Run()
		if err != nil {
			fmt.Printf("분할 오류 (시간대 %d): %s\n", i+1, err)
			return
		}
		fmt.Printf("시간대 %d: %s에서 %s까지의 동영상을 분할 완료\n", i+1, startTimes[i], endTimes[i])
	}

	// 분할된 동영상 파일들을 리스트업
	var fileList []string
	fileListTxt := "filelist.txt"
	for i := 1; ; i++ {
		fileName := fmt.Sprintf("output_%02d.mp4", i)
		outputFile := filepath.Join(outputDirectory, fileName)
		_, err := os.Stat(outputFile)
		if os.IsNotExist(err) {
			break
		}
		fileList = append(fileList, "file '"+outputFile+"'")
	}

	// 파일 목록을 filelist.txt에 저장
	fileListContents := strings.Join(fileList, "\n")
	err = ioutil.WriteFile(fileListTxt, []byte(fileListContents), 0644)
	if err != nil {
		fmt.Println("파일 목록 생성 오류:", err)
		return
	}

	// 분할된 동영상 파일을 다시 합치기 위한 FFmpeg 명령
	concatCmd := exec.Command("ffmpeg",
		"-y",
		"-f", "concat",
		"-safe", "0",
		"-i", "filelist.txt",
		"-c:v", "libx264",
		"-c:a", "aac",
		"output_final1.mp4")
	//concatCmd.Dir = outputDirectory
	err = concatCmd.Run()
	if err != nil {
		fmt.Println("합치기 오류:", err)
		return
	}
	fmt.Println("동영상이 성공적으로 합쳐졌습니다.")
}
