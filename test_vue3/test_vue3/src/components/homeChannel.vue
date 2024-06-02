<template>
  <div class="home-channel">
    <div class="channel-icons">
      <div
        class="avatar-container"
        :style="{ 'padding-right': '20px', 'padding-left': '10px' }"
      >
        <el-avatar
          src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png"
          :size="50"
        ></el-avatar>
        <span class="avatar-text">动态</span>
      </div>
      <div class="avatar-container">
        <el-avatar
          src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png"
          :size="50"
        ></el-avatar>
        <span class="avatar-text">热门</span>
      </div>
    </div>

    <div>
      <div class="channel-items-left">
        <div v-for="data in classifyList" :key="data.id" class="classify">
          <el-button size="medium" round>{{ data.name }}</el-button>
        </div>
      </div>

      <div class="channel-items-right">
        <el-link target="_blank" :underline="false" :icon="ElIconTickets"
          >专栏<el-icon class="el-icon--left"><ElIconTickets /></el-icon></el-link
        >
        <el-link target="_blank" :underline="false" 
          >活动 <el-icon class="el-icon--left"><ElIconFlag /></el-icon></el-link
        >
        <el-link target="_blank" :underline="false" 
          >社区中心<el-icon class="el-icon--left"><ElIconComment /></el-icon></el-link
        >
        <el-link target="_blank" :underline="false" :icon="ElIconPlatform"
          >新歌热榜
        </el-link>
      </div>
    </div>
  </div>
</template>

<script>
import {
  Tickets as ElIconTickets,
  Flag as ElIconFlag,
  Comment as ElIconComment,
  Platform as ElIconPlatform,
} from '@element-plus/icons-vue'
export default {
  data() {
    return {
      classifyList: [
        { id: 1, name: '娱乐' },
        { id: 2, name: '娱乐' },
        { id: 3, name: '娱乐' },
        { id: 4, name: '娱乐' },
        { id: 5, name: '娱乐' },
        { id: 6, name: '娱乐' },
        { id: 7, name: '娱乐' },
        { id: 8, name: '娱乐' },
        { id: 9, name: '娱乐' },
        { id: 10, name: '娱乐' },
        { id: 11, name: '娱乐' },
        { id: 12, name: '娱乐' },
        { id: 13, name: '娱乐' },
        { id: 14, name: '娱乐' },
        { id: 15, name: '娱乐' },
        { id: 16, name: '娱乐' },
        { id: 17, name: '娱乐' },
        { id: 18, name: '娱乐' },
        { id: 19, name: '娱乐' },
        { id: 20, name: '娱乐' },
        { id: 21, name: '娱乐' },
        { id: 22, name: '娱乐' },
      ]
    }
  },
}
</script>

<style scoped>
* {
  margin: 0;
  padding: 0;
  display: flex;
  justify-content: flex-start;
  font-size: 15px;
}

.channel-icons {
  width: 160px;
  height: 100px;
  align-items: center;
}

.avatar-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-right: 10px;
}

.avatar-text {
  margin-top: 5px; /* 可根据需要调整文字与图片的间距 */
}

.channel-items-left {
  padding-right: 0px;
  width: 1150px;
  height: 80px;
  margin-left: 10px;
  flex-wrap: wrap;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 16px;
  border: 1px solid rgb(216, 215, 215);
  border-left: none;
  padding-top: 4px;
  border-bottom: none;
  border-top: none;
}

.classify {
  height: 50px;
  justify-content: center;
}

.el-button {
  height: 35px;
  margin-right: 10px;
  width: 90px;
  background-color: rgb(246, 247, 248, 0.9);
  align-items: center;
  justify-content: center;
  letter-spacing: 2px;
  border-radius: 6px;
  font-weight: 10px;
}

.channel-items-right {
  width: 260px;
  height: 75px;
  flex-wrap: wrap;
  margin-top: 25px;
  padding-left: 30px;
}
.el-link {
  margin-right: 15px;
  font-size: 15px;
  width: 80px;
}
</style>
