from flask import Flask, request, jsonify
from flask_cors import CORS
import time
import threading, requests
from main_func import *
from software_project import *

app = Flask(__name__)

# CORS 설정: 모든 도메인에서 접근을 허용합니다. 실제 배포 환경에서는 더 제한적인 설정을 사용해야 합니다.
CORS(app)
@app.route('/summarisevideo', methods=['GET'])
def summarise_video():
    # POST 요청의 데이터를 받아올 수 있습니다.
    # 이 예제에서는 데이터를 사용하지 않고 "hello"를 반환합니다.
    return jsonify({"message": "Hello, video summarization!"})

@app.route('/receive_data', methods=['POST'])
def receive_data():
    data = request.json  # JSON 데이터를 받음

    # 백그라운드 스레드를 시작하여 장기 실행 작업 수행
    long_running_task()

    background_thread = threading.Thread(target=long_running_task)
    background_thread.start()

    # 즉시 응답을 클라이언트에게 보내고, 작업 결과는 나중에 처리
    response_data = {"message": "Request received by B server", "received_data": data}
    return jsonify(response_data)


#TODO 모든 변수 초기화 해야함 오디오파일 비디오파일 자막파일
def long_running_job():
    # 시간이 오래 걸리는 작업

    #TODO: video 패스 go 백엔드에 path에 맞춰서 수정
    video_file_path = "../sfprojectBack/originalVideo/original.mp4" # .mp4 파일로 끝나는 path
    audio_file_path = "./audio1.mp3" # .mp3 파일로 끝나는 path
    subtitle_file_path = "subtitle" # 폴더로 끝나는 path
    summerized_subtitle_file_path = "summerized_subtitle" # 폴더로 끝나는 path
    print("================================================")

    result_times = auto_editing_video(video_file_path, audio_file_path, subtitle_file_path, summerized_subtitle_file_path)
    print(result_times)
    # strarr = ""
    # for i in result_times:
    #     strarr += str(i)
    # summerized_timeline_path = "summerized_timeline"

    # save_summerized_timeline(summerized_timeline_path, strarr)
    # 작업 완료 후 A 서버에게 요청 보내기
    response = requests.get("http://localhost:8000/video_summary")
    print(response.status_code)

@app.route('/start')
def start_job():
    # 별도의 스레드에서 시간이 오래 걸리는 작업 실행
    threading.Thread(target=long_running_job).start()
    return jsonify({"status": "Job started, checking status..."})


if __name__ == '__main__':
    app.run(port=8070)

