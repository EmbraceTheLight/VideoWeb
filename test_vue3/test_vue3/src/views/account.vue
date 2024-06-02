<template>
  <div class="root">
    <div class="header">
      <homeNavbar></homeNavbar>
      <div class="bar"></div>
    </div>
    <div class="container">
      <el-tabs tab-position="left" type="border-card" v-model="activeName">
        <el-tab-pane label="首页" name="home">
          <div>基本资料</div>
        </el-tab-pane>
        <el-tab-pane label="我的信息" name="info">
          <div :style="{marginLeft: '50px'}">
            <el-text>用户名：</el-text>
            <el-input v-model="username" :style="{width: '225px',marginLeft: '10px'}"></el-input>
          </div>
          <div :style="{marginLeft: '50px',marginTop: '20px'}">
            <el-text>用户ID：</el-text>
            <el-text :style="{marginLeft: '10px'}">567892</el-text>
          </div>
          <div style="display: flex; align-items: center; margin-left: 40px; margin-top: 20px;">
            <el-text style="margin-right: 10px;position: relative;bottom:35px;">我的签名：</el-text>
            <el-input style="width: 615px;" type="textarea" v-model="signature" rows="4"
              placeholder="请输入你的签名"></el-input>
          </div>
          <div class="footerBtn">
            <el-button
              :style="{width: '110px',height: '40px',background:'rgb(0, 160, 216)',color:'white'}">保存</el-button>
          </div>
        </el-tab-pane>
        <el-tab-pane label="我的头像" name="avatar">        
<div v-if="isavatar">
    <canvas id="myCanvas" width="250" height="250"></canvas>
    <el-avatar class="change-avatar" :size="50" @click="changeAvatar">更换头像</el-avatar>
     <el-avatar class="my-avatar" :size="100">user</el-avatar>
</div>
        </el-tab-pane>
        <el-tab-pane label="账号安全" name="security">
          <div :style="{marginTop: '10px'}">
            <el-breadcrumb>
              <el-breadcrumb-item @click="backhome">home</el-breadcrumb-item>
              <el-breadcrumb-item v-show="!isHome">{{ cur_func }}</el-breadcrumb-item>
            </el-breadcrumb>
          </div>

          <div v-if="isHome">
            <div class="item" style="margin-Top: 30px;">
              <el-avatar :size="40" :style="{marginLeft: '20px',background:'rgb(0, 160, 216)'}"></el-avatar>
              <el-text :style="{marginLeft: '5px',fontSize: '13px'}">修改密码</el-text>
              <el-button round
                :style="{marginLeft: '300px',width: '80px',height: '30px',background:'rgb(0, 160, 216)',color:'white'}"
                @click="handleClick('修改密码')">去修改</el-button>
            </div>
            <div class="item">
              <el-avatar :size="40" :style="{ marginLeft: '20px', background: 'rgb(0, 160, 216)' }"></el-avatar>
              <el-text :style="{ marginLeft: '5px', fontSize: '13px' }">更改邮箱</el-text>
              <el-button round
                :style="{ marginLeft: '300px', width: '80px', height: '30px', background: 'rgb(0, 160, 216)', color: 'white' }"
                @click="handleClick('更改邮箱')">去更改</el-button>
            </div>
            <div class="item">
              <el-avatar :size="40" :style="{ marginLeft: '20px', background: 'rgb(0, 160, 216)' }"></el-avatar>
              <el-text :style="{ marginLeft: '5px', fontSize: '13px' }">重置密码</el-text>
              <el-button round
                :style="{ marginLeft: '300px', width: '80px', height: '30px', background: 'rgb(0, 160, 216)', color: 'white' }"
                @click="handleClick('重置密码')">去重置</el-button>
            </div>
          </div>

          <div class="updatePwd" v-if="cur_func === '修改密码'&&!isHome">
            <div :style="{marginTop: '30px'}">
              <el-text>输入旧密码：</el-text>
              <el-input type="password" v-model="oldPwd" placeholder="请输入旧密码" :show-password="true"
                :style="{width: '225px',marginLeft: '10px'}"></el-input>
            </div>
            <div :style="{ marginTop: '30px' }">
              <el-text>输入新密码：</el-text>
              <el-input type="password" v-model="password" placeholder="请输入新密码" :show-password="true"
                :style="{ width: '225px', marginLeft: '10px' }"></el-input>
            </div>
            <div :style="{ marginTop: '30px' }">
              <el-text>重复新密码：</el-text>
              <el-input type="password" v-model="repassword" :show-password="true" placeholder="请重复输入新密码"
                :style="{ width: '225px', marginLeft: '10px' }"></el-input>
            </div>
            <div class="footerBtn">
              <el-button
                :style="{ width: '110px', height: '40px', background: 'rgb(0, 160, 216)', color: 'white' }">保存</el-button>
            </div>
          </div>


          <div v-if="cur_func === '更改邮箱' && !isHome" class="updataEmail">
            <div :style="{ marginTop: '30px' }">
              <el-text>输入新邮箱：</el-text>
              <el-input v-model="newEmail" placeholder="请输入新邮箱"
                :style="{ width: '225px', marginLeft: '10px' }"></el-input>
            </div>
            <div :style="{ marginTop: '30px' }">
              <el-text>输入验证码：</el-text>
              <el-input v-model="code" placeholder="请输入验证码" :style="{ width: '225px', marginLeft: '10px' }">
                <template #suffix>
                  <el-button
                    :style="{width: '60px',height: '25px',background:'rgb(0, 160, 216)',color:'white'}">获取</el-button>
                </template>
              </el-input>
            </div>
            <div class="footerBtn">
              <el-button
                :style="{ width: '110px', height: '40px', background: 'rgb(0, 160, 216)', color: 'white' }">保存</el-button>
            </div>
          </div>

          <div class="resetPwd" v-if="cur_func === '重置密码'&&!isHome">
            <div :style="{ marginTop: '30px' }">
              <el-text>输入邮箱：</el-text>
              <el-input v-model="newEmail" placeholder="请输入邮箱"
                :style="{ width: '225px', marginLeft: '10px' }"></el-input>
            </div>
            <div :style="{ marginTop: '30px' }">
              <el-text>输入新密码：</el-text>
              <el-input type="password" v-model="password" placeholder="请输入新密码" :show-password="true"
                :style="{ width: '225px', marginLeft: '10px' }"></el-input>
            </div>
            <div :style="{ marginTop: '30px' }">
              <el-text>重复输入密码：</el-text>
              <el-input type="password" v-model="repassword" placeholder="请重复输入密码" :show-password="true"
                :style="{ width: '225px', marginLeft: '10px' }"></el-input>
            </div>
            <div :style="{ marginTop: '30px' }">
              <el-text>输入验证码：</el-text>
              <el-input v-model="code" placeholder="请输入验证码" :style="{ width: '225px', marginLeft: '10px' }">
                <template #suffix>
                  <el-button
                    :style="{ width: '60px', height: '25px', background: 'rgb(0, 160, 216)', color: 'white' }">获取</el-button>
                </template>
              </el-input>
            </div>
            <div class="footerBtn">
              <el-button
                :style="{ width: '110px', height: '40px', background: 'rgb(0, 160, 216)', color: 'white' }">保存</el-button>
            </div>
          </div>
        </el-tab-pane>
        <el-tab-pane label="个人中心" name="center">
          <div>安全中心</div>
        </el-tab-pane>
      </el-tabs>
    </div>
    <div class="footer"></div>
  </div>
</template>

<script>
import homeNavbar from '@/components/homeNavbar.vue'
export default {
  components: {
    homeNavbar,
  },
  data() {
    return {
      activeName: 'home',
      username: '这是用户名',
      signature: '这是我的签名',
      cur_func: '重置密码',
      funcText: ['重置密码', '修改邮箱', '修改密码'],
      isHome: true,
      newEmail: '',
      code: '',
      password: '',
      repassword: '',
      oldPwd: '',
      isavatar:true,
    }
  },
  methods: {
    handleClick(name) {
      this.cur_func = name;
      this.isHome = false;
      console.log('cur_func:' + this.cur_func + 'isHome:' + this.isHome)
    },
    changeAvatar(){
      this.isavatar=false;
    },
    backhome() {
      this.isHome = true;
      this.cur_func = '';
      this.newEmail = '';
      this.code = '';
      this.password = '';
      this.repassword = '';
      this.oldPwd = '';
    }
  },
  mounted() {
    var canvas = document.getElementById('myCanvas');
    const userid = localStorage.getItem('userId');
    var context = canvas.getContext('2d');
    context.strokeStyle = '#00a0d8';
    // 绘制圆形
    context.beginPath();
    context.arc(125, 125, 100, 0, 2 * Math.PI);
    context.stroke();

    const token = localStorage.getItem('token');
    const headers = {
      'Authorization': token,
    };
    console.log(userid);
   }
  }
</script>

<style scoped>
* {
  margin: 0;
  padding: 0;
}
        canvas {
          z-index:-1;
          pointer-events: none;
          margin-left:300px;
        }
.bar {
  background: rgb(0, 160, 216);
  height: 85px;
  margin-top:30px;
  pointer-events: none;
}
.header {
  display: grid;
  grid-template-rows: 35px 85px;
  border:1px solid #ccc;
}
.container {
  margin-top: 30px;
  position: relative;
  border: 1px solid #ccc;
  width: 980px;
  height: 1200px;
  margin-left: calc(50vw - 490px);
  display: grid;
  grid-template-columns: 155px calc(100%-155px);
}

.el-menu-item {
  width: 155px;
}
.footer {
  margin-top: 50px;
  height: 200px;
}

.footerBtn {
  display: flex;
  justify-content: center;
  margin-top: 50px;
}

.item{
  width:500px;
  height: 65px;
  margin-left:20px;
  display: flex;
  align-items: center;
  margin-bottom:30px;
  border-radius:5px;
  border:1px solid #ccc;
  box-shadow: 1px 1px 5px #ccc;
}

.updatePwd{
  display:grid;
  align-items: center;
  justify-content: center;
}

.updataEmail{
    display: grid;
      align-items: center;
      justify-content: center;
}

.resetPwd{
  display: grid;
  align-items: center;
  justify-content: center;
}

.tab-avatar{
  display: flex;
  justify-content: center;
  align-items: center;
}

.change-avatar{
  background-color: rgb(0, 160, 216);
  position:relative;
  right:250px;
  bottom:130px;
  font-size:13px;
}

 
.change-avatar:hover{
  cursor: pointer;
}

.my-avatar{
  position:relative;
  bottom:120px;
  right:225px;
  font-size:13px;
}


</style>
