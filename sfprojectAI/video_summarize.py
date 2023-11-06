from moviepy.editor import VideoFileClip
from moviepy.video.io.ffmpeg_tools import ffmpeg_extract_audio
from datetime import timedelta
import os
import whisper
import openai

# openai 엄지호 api key 
openai_api_key = "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"
openai.api_key = openai_api_key

import os

gpt_model = "gpt-3.5-turbo-16k"
#gpt_model = "gpt-4"
temperature = 0.3
model = whisper.load_model("medium")

# 가져욜 audio_file_path와 자막을 저장할 output_directory
def transcribe_audio(audio_file_path, output_directory):
    print("Whisper model loaded.")
    transcribe = model.transcribe(audio=audio_file_path)
    segments = transcribe['segments']

    # 경로를 설정하여 자막 파일 이름을 포함합니다.
    srt_filename = os.path.join(output_directory, "output.srt")  
    with open(srt_filename, 'w', encoding='utf-8') as srt_file:
        for segment in segments:
            start_time = str(timedelta(seconds=int(segment['start'])))
            end_time = str(timedelta(seconds=int(segment['end'])))
            text = segment['text']
            srt_file.write(f"{start_time} --> {end_time}\n{text}\n\n")

    return srt_filename

import re

# 영상 시간만 자르는 함수
def extract_times(text):
    # 정규 표현식을 사용하여 시간 정보 추출
    pattern1 = r'(\d+:\d+:\d+) --> (\d+:\d+:\d+)'
    pattern2 = r'(\d+:\d+:\d+) - (\d+:\d+:\d+)'

    matches1 = re.findall(pattern1, text)
    matches2 = re.findall(pattern2, text)

    matches = matches1 + matches2
    # 추출된 시간 정보를 리스트에 저장
    times = []
    for match in matches:
        times.append((match[0], match[1]))
    return times

import os

# summerized_subtitle_file_path: 생성된 txt 파일을 저장할 directory path, summerized_timeline: 위에서 auto_editing_video의 return 값
# 저장 형식은 f"Start time: {start_time}, End time: {end_time}\n"
def save_summerized_timeline(summerized_subtitle_file_path, summerized_timeline):
    file_count = 0
    if os.path.exists(summerized_subtitle_file_path) and os.path.isdir(summerized_subtitle_file_path):
        # 지정한 경로가 존재하고 디렉토리인지 확인
        files = os.listdir(summerized_subtitle_file_path)  # 디렉토리 내의 모든 파일 목록을 가져옴
        file_count = len(files)
        print(f"디렉토리 '{summerized_subtitle_file_path}'에는 {file_count}개의 파일이 있습니다.")
    else:
        print(f"'{summerized_subtitle_file_path}'는 존재하지 않거나 디렉토리가 아닙니다.")
        print("다음 위치에 폴더를 생성합니다.")
        print(summerized_subtitle_file_path)
        os.makedirs(summerized_subtitle_file_path, exist_ok=True)

    # 타임라인만 파일로 만들어서 저장
    filename = os.path.join(summerized_subtitle_file_path , (str(file_count) + ".txt"))
    with open(filename, 'w') as file:
        file.write(summerized_timeline)
        # for (start_time, end_time) in summerized_timeline:
        #     file.write(f"{start_time}~{end_time}\n")




# text로 생성된 자막 받아와서 GPT에 요약 요청하는 함수
def get_summerize_text(text):
  messages1 = []
  content1 = """
     다음 글은 시사뉴스 자막인데 자막이 진행되는 시간과 자막 내용을 줄거야.
     그러면 너는 먼저 자막 내용을 요약해줘
  """  + text

  completion1 = openai.ChatCompletion.create(
      model=gpt_model,
      messages = [
          {"role":"user", "content": content1}
      ],
      temperature = temperature
  )
  chat_response1 = completion1.choices[0].message.content
  print("==================================================")

  messages2 = []
  content2 = """
    내가 첫 번째로 요약된 내용의 글을 줄거야 그러면 너는 2번째로 받은 글에서 요약된 내용을 찾아서 자막 타임라인을 알려줘. 몇 분 몇 초부터 몇 분 몇 초인지.
  """ + chat_response1 + """다음 글에서 위의 내용을 찾아서 타임라인에 매칭시켜줘.""" + text
  completion2 = openai.ChatCompletion.create(
      model = gpt_model,
      messages = [
          {"role":"user", "content": content2}
      ],
      temperature = temperature
  )
  chat_response2 = completion2.choices[0].message.content
  # 요약되고 타임라인 매칭된 자막 return
  return chat_response2

# 마지막에 이 함수 하나만 실행하면 모든게 가능하게 만들기 
# video_file_path - 비디오 파일이 있는 경로 및 이름 .mp4로 끝나야 함, audio_file_path - 오디오 저장 경로 및 이름 .mp3로 끝나야 함
# subtitle_file_path - 생성된 자막 저장할 경로(폴더), summerized_subtitle_file_path - 요약된 자막 저장 경로(폴더)
def auto_editing_video(video_file_path, audio_file_path, subtitle_file_path, summerized_subtitle_file_path):
    video = VideoFileClip(video_file_path)
    # video to audio
    ffmpeg_extract_audio(video_file_path, audio_file_path)
    # srt subtitles
    srt_output = transcribe_audio(audio_file_path, subtitle_file_path)
    print(f"SRT 자막 파일이 생성되었습니다: {srt_output}")
    # SRT 파일 경로 설정
    srt_file_path = srt_output
    
    # SRT 파일 열기
    with open(srt_file_path, 'r', encoding='utf-8') as srt_file:
        srt_contents = srt_file.read()
    
    # 생성된 자막 요약 시작 및 타임라인 매칭
    summerized_text = get_summerize_text(srt_contents)

    # 요약된 자막 summerized_subtitle_file_path에 저장
    save_summerized_timeline(summerized_subtitle_file_path, summerized_text)
    
    # 요약된 자막 및 타임라인 출력
    print(summerized_text)
    
    # 요약된 자막에서 타임라인만 추출
    times = extract_times(summerized_text)
    
    # 추출된 시간 출력
    for start_time, end_time in times:
        print(f"Start Time: {start_time}, End Time: {end_time}")

    # 2차원 리스트로 추출된 시간 return
    return times

# 함수 사용 예시
video_file_path = "./video1.mp4"
audio_file_path = "./audio1.mp3"
subtitle_file_path = "subtitle"
summerized_subtitle_file_path = "summerized_subtitle"
print("================================================")
# 타임라인만 표시된 2차원 리스트가 result_times에 저장됨 -> 가끔 GPT에서 요약을 잘못하면 타임라인 생성이 안됨 -> 개선 필요
result_times = auto_editing_video(video_file_path, audio_file_path, subtitle_file_path, summerized_subtitle_file_path)
print(result_times)
# strarr = ""
# for i in result_times:
#     strarr += str(i)
# summerized_timeline_path = "./summerized_timeline"

# save_summerized_timeline(summerized_timeline_path, strarr)
