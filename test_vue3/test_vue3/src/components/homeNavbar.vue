<template>
  <div class="home_navbar" id="navbar">
    <div class="left-container">
      <el-menu mode="horizontal" :ellipsis="false" :style="{ border: 'none', 'margin-right': '0px' }"
        active-text-color="#ffd04b">
        <el-menu-item :style="{ 'padding-left': '0px' }" class="no-bounce"><el-icon size="20px">
            <SvgIcon iconName="icon-Bzhan" />
          </el-icon><el-link :underline="false" href="#/home">首页</el-link></el-menu-item>
        <el-menu-item class="text">番剧</el-menu-item>
        <el-menu-item class="text">游戏中心</el-menu-item>
        <el-menu-item class="text">会员购</el-menu-item>
        <el-menu-item class="text">漫画</el-menu-item>
        <el-menu-item class="text">赛事</el-menu-item>
      </el-menu>
    </div>

    <div class="search-container">
      <search_bar v-show="!isSearch"></search_bar>
    </div>

    <div class="right-container">

      <el-menu mode="horizontal" :style="{ border: 'none' }" :ellipsis="false">
        <el-menu-item>
          <div class="avatar-container" @mouseenter="handleMouseEnter" @mouseleave="handleMouseLeave">
            <el-avatar ref="avatar" @click.native="handleAvatar" class="neg-login" v-if="!isLogin">
              登录
            </el-avatar>
            <el-avatar ref="avatar" :style="{
              Transform: isHovered ? 'scale(2.1) translate(-10px, 15px)' : 'none',
            }" @click.native="handleAvatar" v-if="isLogin" :src="avatar">
            </el-avatar>
            <person_widget v-if="isHovered && isLogin"></person_widget>
          </div>
        </el-menu-item>
        <el-menu-item>
          <el-badge class="menu-item">
            <el-icon size="25px" class="centered-icon">
              <SvgIcon iconName="icon-wodedahuiyuan" />
            </el-icon>
            大会员
          </el-badge>
        </el-menu-item>
        <el-menu-item>
          <el-badge class="menu-item">
            <el-icon size="20px" class="centered-icon">
              <SvgIcon iconName="icon-sixin" />
            </el-icon>
            消息
          </el-badge>
        </el-menu-item>
        <el-menu-item>
          <el-badge is-dot :offset="[-6, 20]" class="menu-item">
            <el-icon size="20px" class="centered-icon">
              <SvgIcon iconName="icon-fengche" />
            </el-icon>
            动态
          </el-badge>
        </el-menu-item>
        <el-menu-item>
          <el-badge class="menu-item">
            <el-icon size="20px" class="centered-icon">
              <SvgIcon iconName="icon-shoucang" />
            </el-icon>
            收藏
          </el-badge>
        </el-menu-item>
        <el-menu-item>
          <el-badge class="menu-item">
            <el-icon size="20px" class="centered-icon">
              <SvgIcon iconName="icon-lishijilu1" />
            </el-icon>
            历史
          </el-badge>
        </el-menu-item>
        <el-menu-item>
          <el-badge class="menu-item">
            <el-icon size="20px" class="centered-icon">
              <SvgIcon iconName="icon-chuangzuozhongxin" />
            </el-icon>
            创作中心
          </el-badge>
        </el-menu-item>
      </el-menu>
    </div>
  </div>
</template>

<script>
import { HomeFilled as ElIconSHome } from '@element-plus/icons-vue'
import { Message } from '@element-plus/icons-vue'
import search_bar from '@/components/searchBar.vue'
import person_widget from './personWidget.vue'
import { el } from 'element-plus/es/locales.mjs';
import SvgIcon from './iconfont/SvgIcon.vue';
import axios from 'axios';
export default {
  components: {
    search_bar,
    person_widget,
    ElIconSHome,
  },
  data() {
    return {
      isHovered: false,
      isLogin: false,
      avatar: '',
      isSearch:false,
    }
  },
  methods: {
    handleAvatar() {
      this.$emit('loginEvent')
    },
    handleMouseEnter() {
      this.isHovered = true
    },
    handleMouseLeave() {
      this.isHovered = false
    },
  },
  mounted() {
    const token = localStorage.getItem('token');
    const keyword = this.$route.params.keyword;
    console.log("keyword--->",keyword)
    if (keyword){
        this.isSearch = true
    }
    const headers = {
      'Authorization': token,
    }
    console.log(token)
    if (token) {
      this.isLogin = true
      console.log('test1')
      axios.get('/yanxi//User/User-detail', { headers }).then(res => {
        const data = res.data.data;
        localStorage.setItem('avatar', data.avatar);
         this.avatar = `data:image/png;base64,${data.avatar}`;
      })
    }
  }
}
</script>

<style scoped>
* {
  margin: 0;
  padding: 0;
}

.home_navbar {
  display: flex;
  justify-content: space-around;
  z-index: 1;
}

.fixed {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  z-index: 999;
  background: white;
}

div {
  display: flex;
  justify-content: space-around;
  z-index: 1;
}

.el-menu {
  border-bottom: none;
  display: flex;
  background: none;
  transition: none;
}

.el-menu-item {
  padding-left: 0px;
  /* 根据需要调整 */
  font-size: 14px;
  height: 60px;
  flex: 1;
  margin-right: 25px;
  transition: none;
  color: rgb(10, 13, 14, 0.8);
}

.el-menu-item:hover {
  background: none !important;
}

.el-menu-item:focus {
  background: none !important;
}

.left-container {
  min-width: 450px;
}

.right-container {
  min-width: 450px;
  margin-left: 20px;
  display: flex;
  justify-content: flex-start;
}

.el-avatar {
  margin-right: 20px;
  margin-top: 0px;
  z-index: 2;
}

.el-avatar {
  position: absolute;
  top:-15px;
  overflow: hidden;
  transition: all 0.3s linear;
  width: 38px;
  height: 38px;
}

.overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  z-index: 1;
}

.avatar-container {
  position: relative;
}

.neg-login {
  background: rgb(0, 174, 236);
  color: white;
  font-size: 12px;
  font-weight: bold;
  text-align: center;
}


.el-button .el-icon {
  font-size: 20px;
}

.icons:hover {
  animation: bounce 0.5s;
}

.el-button:hover {
  background-color: transparent !important;
  color: none !important;
  box-shadow: none !important;
}

@keyframes bounce {

  0%,
  100% {
    transform: translateY(0);
  }

  50% {
    transform: translateY(-5px);
  }
}

.menu-item {
  justify-content: center;
  display: grid;
  grid-template-rows: 20px 15px;
  row-gap: 5px;
  padding-top: 15px;
  align-items: center;
}
 
 .el-menu-item:hover {
  background-color: transparent !important;
    /* 设置背景色为透明 */
    color: inherit !important;
    /* 使用继承的颜色 */
 }

 .el-menu-item:not(.no-bounce):hover .el-icon{
   animation:  bounce 0.5s;
 }

 .text:hover{
  animation: bounce 0.5s;
 }
.centered-icon {
  margin: 0 auto;
}

</style>
