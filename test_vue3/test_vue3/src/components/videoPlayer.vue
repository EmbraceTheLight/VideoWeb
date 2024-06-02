<template>
  <div id="player">
    <video id="videoPlayer" ref="videoRef" ></video>
    <div id="danmaku-container" style="position: absolute; top: 0; left: 0; width: 100%; height: 100%; pointer-events: none;z-index:2147483647;"></div>
  </div>
</template>

<script>
import 'plyr/dist/plyr.css';
import Plyr from 'plyr'
import axios from 'axios'
import dashjs from 'dashjs'
import Danmaku from 'danmaku'
import { vi } from 'element-plus/es/locales.mjs';
export default {
  data() {
    return {
      basePath: '', // 用于存放视频地址
      player: null, // 用于存放Plyr实例
      videoRef: null,
      danmaku: null,
      danmuPool: null,
      barrgeList: [
        { content: '弹幕1', time: 10, userId: '123456', userName: '用户1', isSended: false },
        { content: '弹幕2', time: 10, userId: '123457', userName: '用户2', isSended: false},
      ], // 弹幕列表
      isSended: false, // 是否已发送弹幕
    }
  },
  mounted() {
    const videoID = this.videoID;
    console.log('Video ID:', videoID);
    const dash = dashjs.MediaPlayer().create();
    const video = document.getElementById('videoPlayer');
    const barrageList = this.barrgeList;
    this.videoRef = this.$refs.videoRef;
    const videoElement = this.$refs.videoRef;
    const token = localStorage.getItem('token');
    const headers = {
      'Authorization': token
    }
    axios.get(`http://localhost:51233/video/${videoID}/VideoDetail`,{headers})
      .then((response) => {
        const data = response.data;
        this.basePath = data.basePath;
        console.log('Base Path:', this.basePath);
        // 初始化 Dash.js 播放器
        dash.initialize(
          video,
          `http://localhost:51233/video/OfferMpd?filePath=${this.basePath}`,
          true
        );

        // 初始化 Plyr 播放器
        this.player = new Plyr(video, {
          captions: { active: true, update: true },
        });

        const danmakuOptions = {
          container: document.getElementById('danmaku-container'),
          media: this.videoRef,
          comments: [
            // 你可以提前加载一些弹幕数据，比如：
            { text: '弹幕1', mode: 'rtl', time: 10, style: { fontSize: '20px', color: '#ffffff' } },
            { text: '弹幕2', mode: 'rtl', time: 15, style: { fontSize: '20px', color: '#ffffff' } },
          ]
          // 可以根据需要添加其他选项
        };
        const danmaku = new Danmaku(danmakuOptions);

        document.addEventListener('fullscreenchange', adjustDanmaku);
        document.addEventListener('webkitfullscreenchange', adjustDanmaku); // 兼容WebKit内核浏览器
        document.addEventListener('mozfullscreenchange', adjustDanmaku); // 兼容Firefox
        document.addEventListener('MSFullscreenChange', adjustDanmaku); // 兼容IE

        function adjustDanmaku() {
          const videoPlayer = document.getElementById('player');
          const danmakuContainer = document.getElementById('danmaku-container');

          if (document.fullscreenElement || document.webkitFullscreenElement || document.mozFullScreenElement || document.msFullscreenElement) {
            danmaku.resize(window.innerWidth, window.innerHeight);
          } else {
            danmaku.resize(window.innerWidth, window.innerHeight);
          }
        }

        // 监听视频播放事件
        // 监听 Dash.js 的片段加载事件

        dash.on(dashjs.MediaPlayer.events.FRAGMENT_LOADING_STARTED, (e) => {
          const request = e.request;

          if (request) {
            //console.log("Original URL:", request.url);

            // 解析原始 URL
            let originalUrl = new URL(request.url);

            // 获取文件名
            let fileName = originalUrl.pathname.split('/').pop();

            // 构建新的 URL
            let newUrl = new URL(originalUrl.origin);

            // 构建路径
            let newPath = '/video/DASHStreamTransmission';

            // 设置新的路径
            newUrl.pathname = newPath;

            // 拼接 filePath 查询参数
            let filePath = `${this.basePath}/${fileName}`;
            newUrl.searchParams.append('filePath', filePath);

            // 保留原始 URL 的其他查询参数
            if (originalUrl.search.length > 1) {
              newUrl.search = originalUrl.search;
            }

            // 将修改后的 URL 应用到请求对象
            request.url = newUrl.toString();

            //console.log("Modified URL:", request.url);
          }
        });
      })
      .catch((error) => {
        console.error('Axios请求出错：', error);
      });
  },
  props: {
    videoID: String,
  }
}
</script>

<style scoped>

#player {
  position: relative; /* 确保这是定位上下文 */
  width: 1120px; /* 自定义宽度 */
  height: 630px; /* 自定义高度 */
  overflow: hidden; /* 隐藏溢出内容 */
}

</style>
