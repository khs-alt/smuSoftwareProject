<template>
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
        <video class="video" ref="video" @loadeddata="extractImages" @loadedmetadata="getTotalTime"
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
        <img src="@/assets/save.png" @click.prevent="downloadVideo" />
        <img src="@/assets/back.png" @click.prevent="skipTime(-3)" />
        <img :src="isPlaying ? require('@/assets/pause.png') : require('@/assets/play.png')"
          @click.prevent="togglePlayPause" />
        <img src="@/assets/forward.png" @click.prevent="skipTime(3)" />
        <img src="@/assets/trim.png" @click="split" />
      </div>
      <!--타임라인 패널-->
      <div class="timeline-panel">
        <div class="playhead">
          <img src="@/assets/playhead.jpeg" :style="{ left: markerPosition + 'px' }" user-drag:none draggable="true"
            @mousedown.prevent="startDragPlayhead" @mousemove.prevent="dragPlayhead" @mouseup.prevent="stopDrag" />
        </div>
        <!--타임라인 영역-->
        <div class="timeline" v-for="(arr, i) in timeline" :key="i">
          <div class="timeline-label" v-for="(img, j) in arr.imgArr" :key="j"
            :style="{ left: ((j + segmentIndex[i]) * timelineImageWidth) + 'px', top: 30 + 'px' }">
            <span class="time-label">{{ img.time }}s</span>
          </div>
          <div class="timeline-arr" :style="{ top: 45 + 'px', left: (segmentIndex[i] * timelineImageWidth) + 'px' }"
            @mouseover="addborder" @mouseleave="removeBorder">
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



//import axios from 'axios';


export default {
  data() {
    return {
      selectedFile: null,   // 파일 선택 여부
      timeline: [], //타임라인
      imgArr: [], //타임라인 이미지 추출 배열
      segmentIndex: [], //재생헤드 분할 위치
      loadedVideo: false, //배열 추출 여부
      isPlaying: false, //동영상 재생 여부
      markerPosition: 0, //마커 위치
      timelineImageWidth: 100, // 타임라인 이미지 간격 조절
      playheadInterval: null, // 재생 헤드 업데이트를 위한 인터벌 변수 추가
      CTime: 0,
      Duration: 0,
    };
  },

  methods: {
    //drag & drop
    handleFileSelect(event) {
      if (event.type === 'drop') {
        event.preventDefault();
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
      this.$refs.video.src = URL.createObjectURL(this.selectedFile);
      this.timeline = [];
      this.imgArr = [];
      this.loadedVideo = false;
      URL.revokeObjectURL(this.selectedFile);
    },

    //재생헤드 이미지 추출
    async extractImages() {
      const video = this.$refs.video;
      const canvas = document.createElement('canvas');
      const ctx = canvas.getContext('2d');

      const duration = Math.floor(video.duration);

      for (let time = 0; time <= duration; time++) {
        video.currentTime = time;
        await video.play();
        canvas.width = video.videoWidth;
        canvas.height = video.videoHeight;
        ctx.drawImage(video, 0, 0, canvas.width, canvas.height);
        const imageDataURL = canvas.toDataURL('image/jpeg'); // 이미지 파일로 변환
        this.imgArr.push({ url: imageDataURL, time });
        await video.pause(); //잦은 재생, 정지로 에러 발생해 추가
      }
      this.segmentIndex[0] = 0;
      this.timeline.push({ imgArr: this.imgArr });
      video.currentTime = 0;
      this.loadedVideo = true;
    },

    addborder(event) {
      event.currentTarget.classList.add('hovered');
    },
    removeBorder(event) {
      event.currentTarget.classList.remove('hovered');
    },

    //현재 재생 시간 표시
    updateTime() {
      const video = this.$refs.video;
      this.CTime = video.currentTime;
    },

    //영상의 총 시간 표시
    getTotalTime() {
      const video = this.$refs.video;
      if (video.readyState >= 2) {
        this.Duration = video.duration;
      }
    },

    //키보드 제어
    handleKeydown(event) {
      const k = event.keyCode;
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

    //영상 다운로드
    downloadVideo() {
      if (this.selectedFile != null) {
        const video = this.$refs.video.src;
        const url = video;
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
      const video = this.$refs.video;
      if (video.readyState >= 2) {
        if (video.paused) {
          video.play();
          this.isPlaying = true;
        } else {
          video.pause();
          this.isPlaying = false;
        }
      }
    },

    //뒤로가기, 건너뛰기
    skipTime(seconds) {
      if (this.selectedFile != null) {
        const video = this.$refs.video;
        const max = video.duration;
        video.currentTime += seconds;
        this.markerPosition = (video.currentTime / video.duration) * this.timelineImageWidth * max;
        if (this.isPlaying) {
          video.play();
          this.isPlaying = true;
        }
      }
    },

    //영상 분할
    split() {
      const video = this.$refs.video;
      const index = Math.round(video.currentTime); //기준 시간
      const tl = this.timeline;
      //axios.post(video.currentTime);

      for (let i = 0; i < tl.length; i++) {
        if (tl[i].imgArr[0].time < index && tl[i].imgArr[tl[i].imgArr.length - 1].time > index) {
          const newVideoSegments = tl[i].imgArr.splice((index - tl[i].imgArr[0].time), (index + tl[i].imgArr.length));
          tl.splice(i + 1, 0, { imgArr: newVideoSegments });
          this.segmentIndex.splice(i + 1, 0, index);
        }
      }
      console.log(tl);
    },

    //재생 끝날 시 자동 멈춤
    hanldeVideoEnded() {
      const video = this.$refs.video;
      this.isPlaying = false;
      video.currentTime = 0;
      this.markerPosition = 0;
    },

    // 재생헤드 드래그 시작
    startDragPlayhead(event) {
      if (this.selectedFile != null && this.loadedVideo) {
        const video = this.$refs.video;
        this.isDragging = true;
        this.startX = event.clientX;
        this.initialTime = this.$refs.video.currentTime;
        this.isPlaying = false;
        video.pause();
        document.addEventListener('mousemove', this.dragPlayhead);
        document.addEventListener('mouseup', this.stopDrag);
      }
    },

    // 재생헤드 드래그 중
    dragPlayhead(event) {
      if (this.isDragging) {
        const deltaX = event.clientX - this.startX;
        const video = this.$refs.video;
        const max = video.duration;
        const newTime = this.initialTime + (deltaX / (this.timelineImageWidth * max)) * video.duration;
        if (newTime >= 0 && newTime <= video.duration) {
          video.currentTime = newTime;
          this.markerPosition = (video.currentTime / video.duration) * this.timelineImageWidth * max;
        }
      }
    },

    // 재생헤드 드래그 정지
    stopDrag() {
      if (this.isDragging) {
        const video = this.$refs.video;
        const max = video.duration;
        this.markerPosition = (video.currentTime / video.duration) * this.timelineImageWidth * max;
        this.isDragging = false;
        document.removeEventListener('mousemove', this.dragPlayhead);
        document.removeEventListener('mouseup', this.stopDrag);
      }
    },

    //재생헤드 영상 시간에 맞춰 움직이게
    playhead() {
      if (this.isPlaying && !this.isDragging) {
        const video = this.$refs.video;
        const max = video.duration;
        this.markerPosition = (video.currentTime / video.duration) * this.timelineImageWidth * max;
      }
    },

    //이미지 클릭 재생 시간 변경
    seekToTime(time) {
      const video = this.$refs.video;
      video.currentTime = time;
      this.isPlaying = true;
      video.play();
    },
  },

  mounted() {
    window.addEventListener('keydown', this.handleKeydown);
    setInterval(this.updateTime, 1);
    setInterval(this.playhead, 1);
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