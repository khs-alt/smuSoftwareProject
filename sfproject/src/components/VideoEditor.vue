<template>
  <div class="summerizedVideo" v-if="isSummerize">
    <div class="sVideo">
      <video class="sumVideo" ref="sVideo"></video>
    </div>
  </div>
  <div class="loadingContainer" v-if="isLoading">
    <div class="load">
      <FadeLoader />
    </div>
  </div>
  <div class="videoEditor" @keydown.prevent="handleKeydown">
    <!--상단부-->
    <div class="top">
      <!--파일 업로드 영역-->
      <div @drop="handleFileSelect" @dragover.prevent class="drop-area">
        <input type="file" ref="fileInput" @change="handleFileSelect" accept="video/*" style="display: none;">
        <p>동영상 파일을 여기로 드래그 앤 드롭하세요.</p>
        <p>또는</p>
        <button @click="openFileInput">파일 선택</button>
      </div>
      <!--비디오 뷰어 패널-->
      <div class="viewer-panel">
        <video class="video" :src="videoSrc" ref="video" @loadeddata="extractImages" @loadedmetadata="getTotalTime"
          @ended="hanldeVideoEnded"></video>
      </div>
    </div>
    <!--하단부-->
    <div class="bottom">
      <!--재생시간-->
      <div class="video-time">
        <span>{{ CTime.toFixed(2) }}</span> / <span>{{ Duration.toFixed(2) }}</span>
      </div>
      <!--컨트롤 패널-->
      <div class="control-panel">
        <!--컨트롤 버튼 영억-->
        <img src="@/assets/edit.png" @click.prevent="clickEditBtn" />
        <!-- <img src="@/assets/save.png" @click.prevent="downloadVideo" /> -->
        <img src="@/assets/back.png" @click.prevent="skipTime(-3)" />
        <img :src="isPlaying ? require('@/assets/pause.png') : require('@/assets/play.png')"
          @click.prevent="togglePlayPause" />
        <img src="@/assets/forward.png" @click.prevent="skipTime(3)" />
        <img src="@/assets/trim.png" @click.prevent="split" />
        <img src="@/assets/remove.png" @click.prevent="removeVideo">
      </div>
      <!--타임라인 패널-->
      <div class="timeline-panel">
        <div class="playhead">
          <img src="@/assets/playhead.jpeg" :style="{ left: markerPosition + 'px' }" user-drag:none draggable="true"
            @mousedown.prevent="startDragPlayhead" @mousemove.prevent="dragPlayhead" @mouseup.prevent="stopDrag" />
        </div>
        <!--타임라인 영역-->
        <div class="timeline-label" v-for="(label, i) in timeLabel" :key="i"
          :style="{ left: i * timelineImageWidth + 'px', top: 30 + 'px' }">
          <span class="time-label">{{ label }}s</span>
        </div>
        <div class="timeline" v-for="(arr, i) in timeline" :key="i">
          <div class="timeline-arr" :style="{ top: 45 + 'px', left: (segmentIndex[i] * timelineImageWidth) + 'px' }"
            @mouseover="addborder" @mouseleave="removeBorder" @click="clickTimeline($event, i)">
            <div class="timeline-image" v-for="(img, j) in arr.imgArr" :key="j"
              :style="{ left: ((j + segmentIndex[i]) * timelineImageWidth) + 'px' }">
              <img :src="img.url" alt="Image">
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import FadeLoader from 'vue-spinner/src/FadeLoader.vue'

export default {
  components: { FadeLoader },
  data() {
    return {
      selectedFile: null,   // 파일 선택 여부
      videoSrc: null,
      timeline: [], //타임라인
      imgArr: [], //타임라인 이미지 추출 배열
      timeLabel: [],
      segmentIndex: [], //재생헤드 분할 위치
      selectedTimelineIndex: -1, //타임라인 클릭 여부
      loadedVideo: false, //배열 추출 여부
      isPlaying: false, //동영상 재생 여부
      markerPosition: 0, //마커 위치
      timelineImageWidth: 100, // 타임라인 이미지 간격 조절
      playheadInterval: null, // 재생 헤드 업데이트를 위한 인터벌 변수 추가
      CTime: 0, //현재 시간
      Duration: 0,  //총 영상 시간
      baseUrl: "http://34.64.149.210:8000",
      modalCheck: false,
      isLoading: false,
      closeBtn: false,
      content: "",
      pollingInterval: null,
      isSummerize: false,
    };
  },

  methods: {

    //drag & drop
    handleFileSelect(event) {
      event.preventDefault();
      if (event.type === 'drop') {
        this.selectedFile = event.dataTransfer.files[0];
      } else if (event.type === 'change') {
        this.selectedFile = event.target.files[0];
      }
      if (this.selectedFile != null) {
        this.markerPosition = 0;
        this.loadVideo();
      }
    },

    //버튼 눌러 파일 업로드
    openFileInput() {
      this.$refs.fileInput.click();
    },

    //영상 로드
    loadVideo() {
      this.videoSrc = URL.createObjectURL(this.selectedFile);
      this.timeline = [];
      this.imgArr = [];
      this.loadedVideo = false;
      URL.revokeObjectURL(this.selectedFile);
      this.postVideo();
    },

    // 비디오 업로드 axios
    postVideo() {
      var formdata = new FormData();
      formdata.append("video", this.selectedFile);
      axios
        .post(this.baseUrl + "/uploadVideo", formdata, {
          // headers: {
          //   'Content-Type' : 'multipart/form-data'
          // }
        })
        .then((response) => {
          console.log(response.data);
          alert(response.data)
        })
        .catch((error) => {
          console.error(error);
        })
    },


    //재생헤드 이미지 추출
    async extractImages() {
      const v = this.$refs.video;
      console.log(v);
      const canvas = document.createElement('canvas');
      const ctx = canvas.getContext('2d');

      //TODO:
      ctx.drawImage(v, 0, 0);
      v.setAttribute('crossorigin', 'anonymous');

      const duration = Math.floor(v.duration);

      for (let time = 0; time <= duration; time++) {
        v.currentTime = time;
        await v.play();
        canvas.width = v.videoWidth;
        canvas.height = v.videoHeight;
        ctx.drawImage(v, 0, 0, canvas.width, canvas.height);
        const imageDataURL = canvas.toDataURL('image/jpeg'); // 이미지 파일로 변환
        this.imgArr.push({ url: imageDataURL, time });
        this.timeLabel[time] = time;
        await v.pause(); //잦은 재생, 정지로 에러 발생해 추가
      }
      this.segmentIndex[0] = 0;
      this.timeLabel[this.timeLabel.length] = this.timeLabel.length;
      this.timeline.push({ imgArr: this.imgArr });
      v.currentTime = 0;
      this.loadedVideo = true;
    },

    //스트리밍 기능
    streamingVideo() {
      const v = this.$refs.video;
      const mediaSource = new MediaSource();
      v.src = window.URL.createObjectURL(mediaSource);

      mediaSource.addEventListener('sourceopen', function () {
        const sourceBuffer = mediaSource.addSourceBuffer('video/mp4');

        axios.get(this.baseUrl + 'segment1.mp4', { responseType: 'arraybuffer' }) //세그먼트 비디오 파일 경로
          .then(response => {
            sourceBuffer.appendBuffer(new Uint8Array(response.data));
          })
          .then(() => {
            v.play();
          })
          .catch(error => {
            console.error('Error loading video segment: ', error);
          });
      });
    },

    //테두리 생성
    addborder(event) {
      if (this.selectedTimelineIndex === -1) {
        event.currentTarget.classList.add('hovered');
      }
    },

    //테두리 제거
    removeBorder(event) {
      if (this.selectedTimelineIndex === -1) {
        event.currentTarget.classList.remove('hovered');
      }
    },

    //현재 재생 시간 표시
    updateTime() {
      const v = this.$refs.video;
      this.CTime = v.currentTime;
    },

    //영상의 총 시간 표시
    getTotalTime() {
      const v = this.$refs.video;
      if (v.readyState >= 2) {
        this.Duration = v.duration;
      }
    },

    //키보드 제어
    handleKeydown(e) {
      const k = e.keyCode;
      if (this.selectedFile != null) {
        //오른쪽 방향키(정밀조절)
        if (k === 39) {
          this.skipTime(0.01);
        }
        //왼쪽 방향키(정밀조절)
        else if (k === 37) {
          this.skipTime(-0.01);
        }
        //스페이스바(재생/정지)
        else if (k === 32) {
          this.togglePlayPause();
        }
      }
    },

    //요약 버튼 클릭
    clickEditBtn() {
      //if (this.selectedFile != null) {
      this.isLoading = true;
      this.editVideo();
      //  }
    },

    //비디오 요약
    editVideo() {
      axios
        .get(this.baseUrl + "/start")
        .then((response) => {
          console.log(response);
          this.isLoading = true;
          this.pollingInterval = setInterval(this.pollJobStatus, 10000);
          console.log(this.videoSrc)
          //this.videoSrc = null;
        })
        .catch((error) => {
          this.isLoading = false;
          console.error(error);
        });
    },
    //TODO:
    async pollJobStatus() {
      console.log("We are in pollJobStatus!!")
      await axios
        .get(this.baseUrl + "/check_video_summary")
        .then(response => {
          if (response.data.status === "completed") {
            this.message = "Job completed";
            this.isLoading = false;
            clearInterval(this.pollingInterval);
            //TODO:
            this.timeline = [];
            this.selectedTimelineIndex = -1;
            this.timeLabel = [];
            this.segmentIndex = [];
            this.videoSrc = "http://localhost:8000/summerized_video";
            //this.isSummerize = true;
          }
        })
        .catch((error) => {
          //this.isLoading = false;
          console.error(error);
        });
    },

    //영상 다운로드
    downloadVideo() {
      if (this.selectedFile != null) {
        const v = this.$refs.video.src;
        const url = v;
        const a = document.createElement('a');
        const fName = this.selectedFile.name;
        a.href = url;
        a.download = fName;
        a.click();
        a.remove();
      }
    },

    //재생 여부에 따라 버튼 변경
    togglePlayPause() {
      const v = this.$refs.video;
      if (v.readyState >= 2) {
        if (v.paused) {
          v.play();
          this.isPlaying = true;
        } else {
          v.pause();
          this.isPlaying = false;
        }
      }
    },

    //시간 전송
    postTime(currentTime) {
      axios
        .post(this.baseUrl + "/postTime", currentTime)
        .then((response) => {
          console.log(response.data);
        })
        .catch((error) => {
          console.error(error);
        })
    },

    //뒤로가기, 건너뛰기
    skipTime(seconds) {
      if (this.selectedFile != null) {
        const v = this.$refs.video;
        const max = v.duration;
        v.currentTime += seconds;
        this.markerPosition = (v.currentTime / v.duration) * this.timelineImageWidth * max;
        this.postTime(v.currentTime);
        if (this.isPlaying) {
          v.play();
          this.isPlaying = true;
        }
      }
    },

    //영상 분할
    split() {
      const v = this.$refs.video;
      const index = Math.round(v.currentTime); //기준 시간
      const tl = this.timeline;

      for (let i = 0; i < tl.length; i++) {
        if (tl[i].imgArr[0].time < index && tl[i].imgArr[(tl[i].imgArr.length - 1)].time + 1 > index) {
          const newVideoSegments = tl[i].imgArr.splice((index - tl[i].imgArr[0].time), index + (tl[i].imgArr.length - 1));
          tl.splice(i + 1, 0, { imgArr: newVideoSegments });
          this.segmentIndex.splice(i + 1, 0, index);
        }
      }
    },

    //클릭된 타임라인 선택
    clickTimeline(e, index) {
      if (index === this.selectedTimelineIndex) {
        this.selectedTimelineIndex = -1;
      }
      else if (this.selectedTimelineIndex === -1) {
        this.selectedTimelineIndex = index;
        e.currentTarget.classList.add('hovered');
      }
    },

    //비디오 삭제
    removeVideo() {
      const v = this.videoSrc;
      if (this.segmentIndex.length > 0 && this.selectedTimelineIndex != -1) {
        let len = this.timeline[this.selectedTimelineIndex].imgArr.length - 1;
        axios
          .post(this.baseUrl + "/cutVideo", {
            "startTime": this.timeline[this.selectedTimelineIndex].imgArr[0].time,
            "endTime": this.timeline[this.selectedTimelineIndex].imgArr[len].time
          })
          .then((response) => {
            // TODO: 편집된 영상 다시 보내주는 기능 구현 필요
            console.log(response.data);
          })
          .catch((error) => {
            console.error(error);
          })
        this.timeline.splice(this.selectedTimelineIndex, 1);
        this.segmentIndex.splice(this.selectedTimelineIndex, 1);
        this.selectedTimelineIndex = -1;
        this.Duration = v.duration
      }
    },

    //재생 끝날 시 자동 멈춤
    hanldeVideoEnded() {
      const v = this.$refs.video;
      this.isPlaying = false;
      v.currentTime = 0;
      this.markerPosition = 0;
    },

    // 재생헤드 드래그 시작
    startDragPlayhead(e) {
      if (this.selectedFile != null && this.loadedVideo) {
        const v = this.$refs.video;
        this.isDragging = true;
        this.startX = e.clientX;
        this.initialTime = this.$refs.video.currentTime;
        this.isPlaying = false;
        v.pause();
        document.addEventListener('mousemove', this.dragPlayhead);
        document.addEventListener('mouseup', this.stopDrag);
      }
    },

    // 재생헤드 드래그 중
    dragPlayhead(e) {
      if (this.isDragging) {
        const deltaX = e.clientX - this.startX;
        const v = this.$refs.video;
        const max = v.duration;
        const newTime = this.initialTime + (deltaX / (this.timelineImageWidth * max)) * v.duration;
        if (newTime >= 0 && newTime <= v.duration) {
          v.currentTime = newTime;
          this.markerPosition = (v.currentTime / v.duration) * this.timelineImageWidth * max;
        }
      }
    },

    // 재생헤드 드래그 중지
    stopDrag() {
      if (this.isDragging) {
        const v = this.$refs.video;
        const max = v.duration;
        this.markerPosition = (v.currentTime / v.duration) * this.timelineImageWidth * max;
        this.isDragging = false;
        this.postTime(v.currentTime);
        document.removeEventListener('mousemove', this.dragPlayhead);
        document.removeEventListener('mouseup', this.stopDrag);
      }
    },

    //재생헤드 영상 시간에 맞춰 움직이게
    playhead() {
      if (this.isPlaying && !this.isDragging) {
        const v = this.$refs.video;
        const max = v.duration;
        this.markerPosition = (v.currentTime / v.duration) * this.timelineImageWidth * max;
      }
    },
  },

  mounted() {
    window.addEventListener('keydown', this.handleKeydown);
    setInterval(this.updateTime, 1);
    setInterval(this.playhead, 1);
    this.extractImages();
    this.playhead();
    this.dragPlayhead();
    this.startDragPlayhead();
    this.stopDrag();
  },

  beforeUnmount() {
    window.removeEventListener('keydown', this.handleKeydown);
    clearInterval(this.updateTime);
    clearInterval(this.playhead);
  },
};

</script>

<style>
.videoEditor {
  width: 100%;
  height: 100%;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-use-select: none;
  user-select: none;
}

/* .modal {
  position: fixed;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.8);
  z-index: 10;
}
/* 
/* modal or popup
.modalContainer {
  position: relative;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 550px;
  background: #fff;
  border-radius: 10px;
  padding: 20px;
  box-sizing: border-box;
  z-index: 11;
} */

.loadingContainer {
  position: fixed;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.8);
  z-index: 10;
}

.load {
  z-index: 11;
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  box-shadow: rgba(0, 0, 0, 0.6);
}

.summerizedVideo {
  position: fixed;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.8);
  z-index: 10;
}

.sVideo {
  z-index: 11;
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  box-shadow: rgba(0, 0, 0, 0.6);
}

.top {
  display: flex;
  flex-wrap: nowrap;
  justify-content: center;
  align-items: center;
  width: 90%;
  height: 320px;
  padding-top: 10px;
  margin: 20px;
}

.drop-area {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  border: 2px dashed #444444;
  padding: 5px;
  width: auto;
  height: 100%;
  cursor: pointer;
  text-align: center;
  margin-right: 20px;
}

.viewer-panel {
  width: auto;
  height: 100%;
  background-color: #444444;
}

.video {
  width: auto;
  height: 320px;
}

.bottom {
  position: relative;
  height: 300px;
  justify-content: center;
  flex-wrap: wrap;
  margin-top: 20px;
}

.video-time {
  text-align: center;
  color: black;
  font-size: 16px;
  margin-top: 20px;
}

.control-panel {
  position: relative;
}

.control-panel img {
  width: 50px;
}

.timeline-panel {
  position: relative;
  top: 40px;
  width: 100%;
  height: 200px;
  background-color: #444444;
  overflow: auto;
  flex-wrap: nowrap;
}

.playhead {
  position: absolute;
  height: 100%;
  z-index: 4;
}

.playhead img {
  position: absolute;
  cursor: grab;
  width: 10px;
  height: 100%;
  left: 0px;
}

.timeline-label {
  position: absolute;
  margin-left: 5px;
  transform: translateX(-50%);
  white-space: nowrap;
  color: white;
  font-size: 12px;
  z-index: 3;
}

.timeline {
  position: absolute;
  width: 100%;
  height: 100%;
  display: flex;
  white-space: nowrap;
}

.timeline-arr {
  position: absolute;
  width: auto;
  height: 100px;
  z-index: 1;
}

.timeline-arr.hovered {
  border: 2px solid yellow;
  z-index: 2;
}

.timeline-image {
  display: inline-block;
}

.timeline-image img {
  width: 100px;
  height: 100px;
  -webkit-user-drag: none;
}

/* 모바일 화면(세로)에서 화면 축소 */
@media screen and (max-width: 768px) {
  .videoEditor {
    width: 100%;
    height: 100%;
  }

  .top {
    display: flex;
    flex-direction: column;
    justify-content: center;
    width: 90%;
    height: auto;
  }

  .drop-area {
    font-size: 12px;
    width: 90%;
    height: 20%;
    margin-right: 0px;
  }

  .viewer-panel {
    width: 100%;
    height: auto;
  }

  .video {
    width: 100%;
    height: auto;
  }

  .drop-area,
  .viewer-panel,
  .control-panel,
  .video-time,
  .timeline-panel {
    margin-top: 10px;
  }

  .bottom {
    position: relative;
    height: 300px;
    justify-content: center;
    flex-wrap: wrap;
    margin-top: 20px;
  }

  .video-time {
    font-size: 15px;
  }

  .control-panel img {
    width: 30px;
    height: auto;
  }


  .timeline-panel {
    margin-top: -20px;
    width: 100%;
    height: 150px;
  }

  .timeline-label {
    margin-left: 5px;
    position: absolute;
    top: 10px;
    transform: translateX(-50%);
    white-space: nowrap;
    color: white;
    font-size: 10px;
  }

  .timeline-image {
    top: 25px;
  }
}
</style>