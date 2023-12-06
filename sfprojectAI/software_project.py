from moviepy.editor import VideoFileClip
from moviepy.video.io.ffmpeg_tools import ffmpeg_extract_audio
from datetime import timedelta
import os
import whisper
import openai
from tqdm.auto import tqdm
import re

# openai 엄지호 api key 
openai_api_key = "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"
openai.api_key = openai_api_key

gpt_model = "gpt-4-1106-preview"
#gpt_model = "gpt-4"
temperature = 0.3

# 모델 로드
model = whisper.load_model("medium")
#먼저 meduim 모델을 로드하고, 추후 공간이 되면 large 모델을 로드합니다.
#model = whisper.load_model("large-v3")
print("Whisper model loaded.")

def format_time(seconds):
    """초를 'HH:MM:SS.mmm' 형태로 변환"""
    return '{:02}:{:02}:{:02}.{:03}'.format(
        int(seconds // 3600),
        int(seconds % 3600 // 60),
        int(seconds % 60),
        int((seconds % 1) * 1000)
    )


# 가져욜 audio_file_path와 자막을 저장할 output_directory
def transcribe_audio(audio_file_path, output_directory):
    transcribe = model.transcribe(audio=audio_file_path)
    segments = transcribe['segments']

    # 경로를 설정하여 자막 파일 이름을 포함합니다.
    srt_filename = os.path.join(output_directory, "subtitle_"+ audio_file_path[-5:-4]+ ".srt")
    with open(srt_filename, 'w', encoding='utf-8') as srt_file:
        for segment in segments:
            start_time = format_time(segment['start'])
            end_time = format_time(segment['end'])
            text = segment['text']
            srt_file.write(f"{start_time} - {end_time}\n{text}\n\n")
    
    return srt_filename

# 영상 시간만 자르는 함수
def extract_times(text):
    # 정규 표현식을 사용하여 시간 정보 추출
    pattern1 = r'\d{2}:\d{2}:\d{2}.\d{3} --> \d{2}:\d{2}:\d{2}.\d{3}'
    pattern2 = r'\d{2}:\d{2}:\d{2}.\d{3} - \d{2}:\d{2}:\d{2}.\d{3}'
    pattern3 = r'\d{2}:\d{2}:\d{2}.\d{3}부터 \d{2}:\d{2}:\d{2}.\d{3}'
    pattern4 = r'\d{2}:\d{2}:\d{2}.\d{3} ~ \d{2}:\d{2}:\d{2}.\d{3}'
    pattern5 = r'\d{2}:\d{2}:\d{2}.\d{3} to \d{2}:\d{2}:\d{2}.\d{3}'
    
    matches1 = re.findall(pattern1, text)
    matches2 = re.findall(pattern2, text)
    matches3 = re.findall(pattern3, text)
    matches4 = re.findall(pattern4, text)
    matches5 = re.findall(pattern5, text)
    
    match1 = [m1.replace(" --> ", "-") for m1 in matches1]
    match2 = [m2.replace(" - ", "-") for m2 in matches2]
    match3 = [m3.replace("부터 ", "-") for m3 in matches3]
    match4 = [m4.replace(" ~ ", "-") for m4 in matches4] 
    match5 = [m5.replace(" to ", "-") for m5 in matches5]
    
    times = match1 + match2 + match3 + match4 + match5
    
    # times = [t.replace(" - ", "-") for t in match]

    return times

# summerized_subtitle_file_path: 생성된 txt 파일을 저장할 directory path, summerized_timeline: 위에서 auto_editing_video의 return 값
# 저장 형식은 f"Start time: {start_time}, End time: {end_time}\n"
def save_summerized_timeline(video_file_path, summerized_subtitle_file_path, summerized_timeline_list):
    if os.path.exists(video_file_path):
        filename = "summerized_subtitle" +".csv" # video_file_path[-5:-4] + ".csv"
        print(f"파일명 '{filename}'으로 저장됨.")
    else:
        filename = "summerized_subtitle_default.csv"

    # 타임라인만 파일로 만들어서 저장
    filename = os.path.join(summerized_subtitle_file_path, filename)
    
    # 리스트의 모든 항목을 하나의 문자열로 연결하고 작은 따옴표를 제거
    csv_timeline = ','.join(summerized_timeline_list).replace('\'', '')
    
    # 파일에 저장
    with open(filename, 'w') as file:
        file.write(csv_timeline)

# text로 생성된 자막 받아와서 GPT에 요약 요청하는 함수
def get_summerize_text(text):
  messages1 = []
  content1 = """
     다음 글은 시사뉴스 자막인데 자막이 진행되는 시간과 자막 내용을 줄거야.
     그러면 너는 전체 자막 내용을 요약해줘
  """  + text

  completion1 = openai.ChatCompletion.create(
      model=gpt_model,
      messages = [
          {"role":"user", "content": content1}
      ],
      temperature = temperature
  )
  chat_response1 = completion1.choices[0].message.content
  print("Sumerizing Subtitle Succesful ==================================================")

  messages2 = []
  content2 = """
    내가 첫 번째로 요약된 내용의 글을 줄거야 그러면 너는 2번째로 받은 글에서 요약된 내용을 찾아서 자막 타임라인을 알려줘. 
    타임라인은 '\d{2}:\d{2}:\d{2}.\d{3}' 이런 형식으로 저장되어 있어.
    몇 분 몇 초부터 몇 분 몇 초인지 알려줘.
  """ + chat_response1 + """다음 글에서 위의 내용을 찾아서 타임라인에 매칭시켜줘.""" + text
    
  completion2 = openai.ChatCompletion.create(
      model = gpt_model,
      messages = [
          {"role":"user", "content": content2}
      ],
      temperature = temperature
  )
  chat_response2 = completion2.choices[0].message.content
  print("Matching Timeline Succesful ==================================================")
  # 요약되고 타임라인 매칭된 자막 return
  return chat_response2

# 마지막에 이 함수 하나만 실행하면 모든게 가능하게 만들기 
# video_file_path - 비디오 파일이 있는 경로 및 이름 .mp4로 끝나야 함, audio_file_path - 오디오 저장 경로 및 이름 .mp3로 끝나야 함
# subtitle_file_path - 생성된 자막 저장할 경로(폴더), summerized_subtitle_file_path - 요약된 자막 저장 경로(폴더)
def auto_editing_video(video_file_path, audio_file_path, subtitle_file_path, summerized_subtitle_file_path):
    video = VideoFileClip(video_file_path)
    # video to audio
    ffmpeg_extract_audio(video_file_path, audio_file_path)
    print("Complete generating audio file")
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

    # 요약된 자막 및 타임라인 출력
    print("타임라인에 매칭된 요약된 자막\n" + summerized_text)
    
    # 요약된 자막에서 타임라인만 추출
    times = extract_times(summerized_text)
    print("타임라인만 추출\n")
    print(times)

    # 요약된 자막 summerized_subtitle_file_path에 저장
    save_summerized_timeline(video_file_path, summerized_subtitle_file_path, times)
    
    # # 추출된 시간 출력
    # for start_time, end_time in times:
    #     print(f"Start Time: {start_time}, End Time: {end_time}")

    # 2차원 리스트로 추출된 시간 return
    return times


# # 함수 사용 예시
# video_file_path = "./video1.mp4" # .mp4 파일로 끝나는 path
# audio_file_path = "./audio1.mp3" # .mp3 파일로 끝나는 path
# subtitle_file_path = "subtitle" # 폴더로 끝나는 path
# summerized_subtitle_file_path = "summerized_subtitle" # 폴더로 끝나는 path
# print("================================================")

# result_times = auto_editing_video(video_file_path, audio_file_path, subtitle_file_path, summerized_subtitle_file_path)
# print(result_times)
# strarr = ""
# for i in result_times:
#     strarr += str(i)
# summerized_timeline_path = "summerized_timeline"

# save_summerized_timeline(summerized_timeline_path, strarr)