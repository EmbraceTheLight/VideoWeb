<template>
  <div>
    <div class="video-container">
      <video ref="videoRef" src="/video/过掉守门员.mp4"></video>
    </div>
    <button @click="sendDanmu">发送弹幕</button>
  </div>
</template>

<script lang="ts">
import 'plyr/dist/plyr.css'
import Danmaku from 'danmaku'
import Plyr from 'plyr'



export default {
  data() {
    return {
      canvasRef: null,
      videoRef: null,
      danmaku: null,
      plyrPlayer: null,
      barrgeList: [
        { content: '弹幕1', time: 10, userId: '123456', userName: '用户1' },
        { content: '弹幕2', time: 10, userId: '123457', userName: '用户2' },
      ]
    }
  },
  mounted() {
    this.canvasRef = this.$refs.canvasRef
    this.videoRef = this.$refs.videoRef

    if (this.canvasRef && this.videoRef) {
      this.plyrPlayer = new Plyr(this.videoRef, {
        captions: { active: true, update: true },
      })
      this.danmaku = new Danmaku({
        // 必填。用于显示弹幕的「舞台」会被添加到该容器中。
        container: document.getElementById('my-container'),

        // 媒体可以是 <video> 或 <audio> 元素，如果未提供，会变成实时模式。
        media: document.getElementById('my-media'),

        // 预设的弹幕数据数组，在媒体模式中使用。在 emit API 中有说明格式。
        comments: [],

        // 支持 DOM 引擎和 canvas 引擎。canvas 引擎比 DOM 更高效，但相对更耗内存。
        // 完整版本中默认为 DOM 引擎。
        engine: 'canvas',

        // 弹幕速度，也可以用 speed API 设置。
        speed: 144
      });
    }
  },

  methods: {
    sendDanmu() {
      if (this.danmaku) {
        this.danmaku.emit({
          text: 'example',

          // 默认为 rtl（从右到左），支持 ltr、rtl、top、bottom。
          mode: 'rtl',

          // 弹幕显示的时间，单位为秒。
          // 在使用媒体模式时，如果未设置，会默认为音视频的当前时间；实时模式不需要设置。
          time: 233.3,

          // 在使用 DOM 引擎时，Danmaku 会为每一条弹幕创建一个 <div> 节点，
          // style 对象会直接设置到 `div.style` 上，按 CSS 规则来写。
          // 例如：
          style: {
            fontSize: '20px',
            color: '#ffffff',
            border: '1px solid #337ab7',
            textShadow: '-1px -1px #000, -1px 1px #000, 1px -1px #000, 1px 1px #000'
          },
        });


      }
    }
  }
}
</script>

<style scoped>
.video-container {
  position: relative;
  width: 1120px;
  height: 630px;
}


button {
  position: relative;
  z-index: 2;
}
</style>
