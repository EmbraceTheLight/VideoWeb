<template>
  <div class="recPlate">
    <div class="recommend-swiper">
      <!-- 轮播组件，根据swiperList中的数据循环展示轮播项 -->
      <el-carousel trigger="click" height="440px" interval="8000">
        <el-carousel-item v-for="item in swiperList" :key="item.id">
          <div class="title">{{ item.title }}</div>
          <div><img :src="item.src" /></div>
          <div class="blur-overlay"></div>
        </el-carousel-item>
      </el-carousel>
    </div>

    <div class="test">
      <div class="recommand-video" v-for="video in videoList.slice(0, 6)" :key="video.VID">
        <!-- 根据videoList中的数据循环展示推荐视频 -->
        <el-image :src="video.Cover" class="image-with-overlay">
        </el-image>
        <div class="video-info">
          <el-button text :style="{ color: 'white' ,marginLeft: '10px'}">
            <el-icon size="18px" :style="{ marginRight: '5px' }">
              <SvgIcon iconName="icon-bofangshu" />
            </el-icon>
            {{ video.CntViews }}
          </el-button>
          <el-button text :style="{ color: 'white', marginLeft: '10px' }">
            <el-icon size="18px" :style="{ marginRight: '5px' }">
              <SvgIcon iconName="icon-danmushu" />
            </el-icon>
            {{ video.CntBarrages }}
          </el-button>
          <el-button text :style="{ color: 'white', marginLeft: 'calc(10vw )' }">
            {{ video.Duration }}
          </el-button>
        </div>
        <div class=" video-title">
          <el-tooltip :content="video.Title" placement="bottom">
            <el-link :underline="false" class="video-title-link" :href="`#/video/VideoID=${video.VID}`">{{
              video.Title }}</el-link>
          </el-tooltip>
        </div>
        <div class="video-footer">
          <el-link :underline="false">
            <el-icon size="15px" :style="{ marginRight: '5px' }">
              <SvgIcon iconName="icon-UPzhu" />
            </el-icon>
            {{ video.VCName }}
          </el-link>
        </div>
      </div>
    </div>

    <!-- 根据videoList中的数据切片后循环展示推荐视频 -->
    <div class="recommand-video" v-for="video in videoList.slice(6)" :key="video.VID">
      <el-image :src="video.Cover" class="image-with-overlay">
      </el-image>
      <div class="video-info">
        <el-button text :style="{ color: 'white', marginLeft: '10px' }">
          <el-icon size="18px" :style="{ marginRight: '5px' }">
            <SvgIcon iconName="icon-bofangshu" />
          </el-icon>
          {{ video.CntViews }}
        </el-button>
        <el-button text :style="{ color: 'white', marginLeft: '10px' }">
          <el-icon size="18px" :style="{ marginRight: '5px' }">
            <SvgIcon iconName="icon-danmushu" />
          </el-icon>
          {{ video.CntBarrages }}
        </el-button>
        <el-button text :style="{ color: 'white', marginLeft: 'calc(10vw )' }">
          {{ video.Duration }}
        </el-button>
      </div>
      <div class=" video-title">
        <el-tooltip :content="video.Title" placement="bottom">
          <el-link :underline="false" class="video-title-link" :href="`#/video/VideoID=${video.VID}`">{{
            video.Title }}</el-link>
        </el-tooltip>
      </div>
      <div class="video-footer">
        <el-link :underline="false">
          <el-icon size="15px" :style="{ marginRight: '5px' }">
            <SvgIcon iconName="icon-UPzhu" />
          </el-icon>
          {{ video.VCName }}
        </el-link>
      </div>
    </div>
  </div>
</template>

<script>
import { Search as ElIconSearch } from '@element-plus/icons-vue'
import bg1 from '../assets/bg.jpg';
import bg2 from '../assets/bg2.jpg';
import bg3 from '../assets/bg3.jpg';
import jsonData from '@/assets/response_1716983420154.json';
import axios from 'axios';
export default {
  data() {
    return {
      swiperList: [
        { id: 1, src: bg1, title: '标题1' },
        { id: 2, src: bg2, title: '标题2' },
        { id: 3, src: bg3, title: '标题3' },
      ],
      videoList: [],
      dataLoaded: false,
      ElIconSearch,
    }
  },
  mounted() {
    const data = null;
    console.log('开始获取数据')
    axios.get('/yanxi/video/VideoList')
      .then(res => {
        this.dataLoaded = true;
        console.log('获取数据成功');
        console.log(res.data.data);
        this.videoList = res.data.data.map(item => {
          return {
            ...item,
            Cover: `data:image/png;base64,${item.Cover}`,
            Duration: this.formatTime(item.Duration)
          }
        });
        console.log(this.videoList);
      })
    const token = localStorage.getItem('token');
  },
  methods: {
    formatTime(timeString) {
      if (timeString.startsWith('00:')) {  // 检查时间字符串是否以 "00:" 开头
        return timeString.substring(3);  // 去除代表小时的数字以及后面的冒号
      }
      return timeString;  // 如果不是以 "00:" 开头，则保持原样返回
    },
    RouterSkip() {
      // 在你的组件方法中，使用 $router.push 进行路由跳转
      window.location.href = '/video/VideoID=123456';
    }
  }
}
</script>

<style scoped>
* {
  margin: 0;
  padding: 0;
}
.recPlate {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-start;
  height: auto;
}
.recommend-swiper {
  width: 635px;
  height: 460px;
  margin-right: 10px;
  margin-top: 10px;
  margin-bottom: 80px;
  z-index: 0;
}
img {
  width: 635px;
}
.blur-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 90px; /* 调整模糊层的高度 */
  background: linear-gradient(
    to bottom,
    rgba(255, 255, 255, 0),
    rgba(255, 255, 255, 1)
  );
  backdrop-filter: blur(30px); /* 调整模糊程度 */
}

.title {
  display: flex;
  position: relative;
  top: 390px;
  left: 20px;
  color: black;
  font-size: 20px;
  font-weight: 12px;
  z-index: 1;
}
.recommend-video {
  display: flex;
  flex-wrap: wrap;
}
.el-image {
  width: 310px;
  height: 170px;
  padding-top: 35px;
  padding-right: 20px;
  border-radius: 5px;
}
.video-title {
  height: 45px;
}
 .video-title .el-link{
  height: 50px;
  width:280px;
  font-size:15px;
  display: flex;
  justify-content: flex-start;
  flex-wrap: wrap ;
 }
.video-footer {
  height: 30px;
  display: flex;
  justify-content: flex-start;
  font-weight: 100;
  font-size: 13px;
}
.test {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
}

.image-with-overlay {
  position: relative;
}

.image-with-overlay::after {
  content: "";
  position: absolute;
  bottom: 0;
  left: 0;
  width: 94%;
  height: 20%;
  /* 使用 CSS 渐变创建透明黑色背景 */
  background: linear-gradient(to bottom, rgba(0, 0, 0, 0), rgba(0, 0, 0, 0.7));
  /* 设置透明黑色的颜色和透明度 */
}

.el-button:hover {
  background-color: transparent !important;
  color: none !important;
  box-shadow: none !important;
}

.video-info{
  position: absolute;
  transform: translateY(-100%);
  z-index:999;
}

.video-title-link{
  color: black;
  font-size: 15px;
}

.video-title-link:hover{
  color: #409eff;
}
</style>
