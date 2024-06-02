<template>
  <div class="person">
    <div class="user_info">
      <div :style="{'font-size': '16px',color:'black'}">这是八个字的用户名</div>
      <div class="shells">
        <el-text :style="{'font-size': '12px',marginRight: '10px'}">硬币：876.0</el-text>
        <el-text :style="{'font-size': '12px'}">b币：0</el-text>
      </div>
    </div>
    <div class="level">
      <div class="top_level">
        <el-icon size="30px">
          <SvgIcon iconName="icon-ic_userlevel_4" />
        </el-icon>
        <el-progress :percentage="50" :style="{'width': '185px'}" :show-text="false" color='orange' :stroke-width="2" />
        <el-icon size="30px"  style="margin-left: 5px">
          <SvgIcon iconName="icon-ic_userlevel_5" />
        </el-icon>
      </div>
    </div>
    <div class=" top_banner">
      <div class="info">
        <div class="number">43</div>
        <div class="text">关注</div>
      </div>
      <div class="info">
        <div class="number">0</div>
        <div class="text">粉丝</div>
      </div>
      <div class="info">
        <div class="number">5</div>
        <div class="text">动态</div>
      </div>
    </div>

    <div class="footer">
      <el-button round>个人中心</el-button>
      <el-button round>投稿管理</el-button>
      <el-button round>推荐服务</el-button>
      <div class="divider"></div>
      <el-button round @click="exit_login">退出登录</el-button>
    </div>
  </div>
</template>

<script>
import SvgIcon from './iconfont/SvgIcon.vue';

export default {
  methods: {
    exit_login() {
      localStorage.removeItem('token')
      localStorage.removeItem('userId')
      window.location.reload()
    },
  },
  data() {
    return {
      userLevel: '',
      followers: '',
      username: '',
      signature: '',
      shells: ''
    }
  },
  mounted() {
    const token = localStorage.getItem('token')
    const headers = {
      'Authorization': token,
    }
    console.log(token)
    if (token) {
      this.isLogin = true
      console.log('test1')
      axios.get('/yanxi//User/User-detail', { headers }).then(res => {
        const data = res.data.data;
        this.userLevel = data.userLevel
        this.followers = data.followers
        this.username = data.userName

        this.signature = data.signature
        this.shells = data.shells
        console.log(data)
      })
    }
  }
}
</script>

<style scoped>
.person {
  position: absolute;
  left: -185px;
  top: 35px;
  width: 300px;
  height: 450px;
  background: white;
  display: inline-block;
  border:1px solid #e2e2e2;
  z-index: 1;
  border-radius: 10px;
}
.info{
  display: grid;
  grid-template-rows: 25px 15px; /* 让行高自适应内容 */
}
.top_banner {
  position: relative;
  top: 70px;
  display: grid;
  grid-template-columns: repeat(3, auto); /* 让列宽自适应内容 */
  grid-template-rows: 70; /* 让行高自适应内容 */
  gap: 10px; /* 添加列之间的间距 */
  margin-left: 20px;
}
.number {
  font-size: 17px;
  font-weight: bold;
  text-align: center;
}
.text {
  font-size: 11px;
  text-align: center; /* 让文字居中对齐 */
}

.footer {
  position: relative;
  top: 100px;
  display: grid;
  grid-template-rows: 35px 35px 35px 15px;
  gap: 5px;
}
.el-button {
  border: none;
}

.el-button:hover {
  background: rgb(223, 226, 227);
  color: none;
}
.divider {
  border: none;
  border-bottom: 1px solid rgb(31, 30, 30, 0.3);
  width: 100%;
  height: 5px;
}

.user_info {
 display: grid;
 justify-content: center;
 margin-top:30px;
 grid-template-rows: 25px 35px;
}

.shells{
  margin-top:5px;
  display: flex;
  justify-content: center;
}

.level{
  display: grid;
  grid-template-rows: 20px 20px;
  justify-content: center;
}

.top_level {
  display: flex;
  align-items: center;

}
</style>
