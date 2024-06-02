<template>
  <div>
    <home_navbar></home_navbar>
    <div><el-divider></el-divider></div>

    <div class="video_detail">
      <div class="video_header">
        <div class="title">这是视频标题</div>
        <div class="video_info">
          <div class="info_item">
            <img src="../assets/icons/播放量.png" class="icons" />
            <div class="text">21.5万</div>
          </div>
          <div class="info_item">
            <img src="../assets/icons/弹幕.png" class="icons" :style="{ height: '15px', color: '#707070' }" />
            <div class="text">722</div>
          </div>
          <div class="info_item">
            <div class="text">2024-03-20 18:01:14</div>
          </div>
          <div class="info_item">
            <img src="../assets/icons/禁止.png" class="icons" :style="{ height: '12px', width: '12px' }" />
            <div class="text">未经作者授权,禁止转载</div>
          </div>
        </div>
      </div>
      <div class="authorPart">
        <el-avatar src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png"
          class="bigAvatar"></el-avatar>
        <div class="rightPart">
          <div class="authorTop">
            <div class="authorName">电竞瓜瓜乐</div>
            <el-link :underline="false" :icon="ElIconMessage">发消息</el-link>
          </div>
          <div class="authorDes">
            别看了,这就是用来进行占位测试的,这是用来占位测试的,这是用来占位测试的
          </div>
          <div class="authorFooter">
            <el-button :style="{
              color: 'rgb(0, 174, 236)',
              border: '1px rgb(0, 174, 236) solid',
            }">充电</el-button>
            <el-button type="primary" :style="{ 'background-color': 'rgb(0, 174, 236)' }">+关注</el-button>
          </div>
        </div>
      </div>
      <video_player :videoID="videoID" v-if="videoID !== null"></video_player>
      <div class="barrageList">
        <el-collapse accordion>
          <el-collapse-item>
            <template slot="title"> 弹幕列表 </template>
            <div class="barrageContext">
              <div class="barrageTop">
                <p>时间</p>
                <p>弹幕内容</p>
                <p>发送时间</p>
              </div>
              <div class="listItem" v-for="item in BarrageList" :key="item.id">
                <p>{{ item.time }}</p>
                <p>{{ item.content }}</p>
                <p>{{ item.date }}</p>
              </div>
              <el-button class="footerBtn">查看历史弹幕</el-button>
            </div>
          </el-collapse-item>
        </el-collapse>
      </div>

      <div class="videoRest">
        <div class="videoFooter">
          <div class="videoFooterLeft">
            <div :style="{ width: '230px', marginLeft: '10px' }">64人在看，已装填985条弹幕</div>
            <el-icon size="25px" :style="{ marginLeft: '0px' }">
              <SvgIcon iconName="icon-bofangqi-danmukai" />
            </el-icon>

            <el-icon size="25px" :style="{ marginLeft: '20px', marginRight: '10px' }">
              <SvgIcon iconName="icon-danmushezhi" />
            </el-icon>

            <el-input :style="{ width: '750px', marginRight: '10px' }" maxLength="100" placeholder="输入弹幕内容"
              v-model="newBarrage">
              <template #prefix>
                <el-icon size="16px" class="icon-zihao">
                  <SvgIcon iconName="icon-zihao" />
                </el-icon>
              </template>
              <template #append>
                <el-button type="primary" size="large"
                  :style="{ 'background-color': 'rgb(0, 174, 236)', color: 'white' }">发送</el-button>
              </template>
            </el-input>
          </div>
        </div>
        <div class="interactions">
          <div class="footerBtn">
            <el-link :underline="false">
              <el-icon size="35px">
                <SvgIcon iconName="icon-zan" />
              </el-icon>
              1212
            </el-link>
          </div>
          <div class="footerBtn">
            <el-link :underline="false">
              <el-icon size="35px">
                <SvgIcon iconName="icon-Bbi" />
              </el-icon>
              1212
            </el-link>
          </div>
          <div class="footerBtn">
            <el-link :underline="false">
              <el-icon size="35px">
                <SvgIcon iconName="icon-shoucang" />
              </el-icon>
              1212
            </el-link>
          </div>
          <div class="footerBtn">
            <el-link :underline="false">
              <el-icon size="35px">
                <SvgIcon iconName="icon-arrow-" />
              </el-icon>
              1212
            </el-link>
          </div>
        </div>
        <div class="videoDescription">
          <div class="description">
            <el-text size="large" type="info" :style="{ color: 'black' }"
              line-clamp="4">文本自动换行工具可用于编程代码、串口调试数据、网站访问记录（*.log文件）等未分行数据分析，将未分行文本按所需换行字符进行自动换行处理。

              输入【换行符】和【待要换行文本】，点击“处理”按钮，将自动返回换行后文本。多次点击“处理”按钮，可调节换行间距。

              注：

              （1）【大框】输入待处理文本，【小框】输入换行符，当换行符为空时，将会在【每个】字之间增加换行符。

              （2）可选择换行符【前】或换行符【后】开始换行，可以选择【删除】或【保留】换行符。

              （3）换后符为【标点符号】时，应注意区分【全角】和【半角】。

              【举例说明】

              原始文本为：FF FF 55 35 10 94 84 6A 6C FF 55 AA 06 18 19 20 0D 64 FF 55 AA 06 18 19 21 0D 65 FF 55 AA 06 18 19
              22 0D 66 FF 55 AA 06 18 19 23 0D 67

              以“FF 55”开始换行处理。将原始文本录入【大框】，将换行字符录入【小框】，点击确认后。</el-text>
          </div>

          <div class="videoTags">
            <el-tag type="info" round="true" effect="light" size="large" v-for="item in videoTags" :key="item"
              :style="{ marginRight: '15px', height: '25px', paddingLeft: '5px', paddingRight: '5px', fontSize: '14px' }">{{
              item }}</el-tag>
          </div>
        </div>
        <div class="commentsContainer">
          <div class="commentsTabs">
            <span :style="{ fontSize: '23px' }">评论</span>
            <span :style="{ fontSize: '12px', marginLeft: '3px' ,fontWeight: '50',marginRight: '60px' }">698</span>
            <el-link :underline="false"
              :style="{ fontSize: '15px', color: 'rgb(148, 153, 160)',width: '40px' }">最热</el-link>
            <el-divider direction="vertical"
              :style="{ position: 'relative', transform: 'translate(-25px,-0px)',  height: '15px',borderWidth:'1px',borderColor:'rgb(148, 153, 160)' }"></el-divider>
            <el-link :underline="false"
              :style="{ fontSize: '15px', color: 'rgb(148, 153, 160)', width: '40px'  }">最新</el-link>
          </div>
          <div class="AddComment">
            <el-avatar :size="45"> user </el-avatar>
            <el-input :style="{ width: '1000px', marginLeft: '20px',height: '50px' }" maxLength="100"
              placeholder="留下友善的评论吧!" class="custom-input"> </el-input>
          </div>
          <div class="comments">
            <el-avatar :size="45">user</el-avatar>
            <el-text :style="{marginLeft: '10px'}">七个字的用户名</el-text>
            <el-icon size="25px" :style="{ marginLeft: '5px' }">
              <SvgIcon iconName="icon-ic_userlevel_4" />
            </el-icon>
            <el-text
              :style="{marginLeft: '55px',fontSize:'16px',flex:'1 0 calc(100% - 55px)',flexWrap:'wrap',wordBreak:'break-all'}">swani算是我看比赛里面进步最大的教练了，从开始的叫个暂停把麦一拉让呼吸说到科隆时候队伍丢分及时叫暂停讲战术开始拿分，本来还以为能一起更进一步，也不知道怎么把人家搞得直接退了，新来的taz目前还没有看见什么明显的妙手，但是选手的状态年龄估计也等不到他像swani那样蜕变的一天了</el-text>
            <div class="commentFooter">
              <el-text size="small" :style="{color:'rgb(148, 153, 160)',marginTop:'5px'}">2021-03-22 11:29</el-text>
              <div :style="{display:'flex',alignItems:'center',marginLeft:'20px',marginTop:'8px'}">
                <el-icon size="15px" color='rgb(148, 153, 160)'>
                  <SvgIcon iconName="icon-zan" />
                </el-icon>
                <el-text size="small"
                  :style="{color:'rgb(148, 153, 160)',marginLeft:'3px',marginBottom:'3px'}">121</el-text>
                <el-icon size="15px" color='rgb(148, 153, 160)' :style="{marginLeft:'20px',marginBottom:'3px'}">
                  <SvgIcon iconName="icon-cai" />
                </el-icon>
                <el-text size="small" class="replyBtn">回复</el-text>
              </div>
            </div>
          </div>
          <div class="reply">
            <el-avatar :size="25">use</el-avatar>
            <el-text :style="{ marginLeft: '10px' }" size="small">七个字的用户名</el-text>
            <el-icon size="25px" :style="{ marginLeft: '5px' }">
              <SvgIcon iconName="icon-ic_userlevel_4" />
            </el-icon>
            <el-text
              :style="{ marginLeft: '25px', fontSize: '16px'}">辅助不跟，跟了也保不住的环境。要求射手一个人面对中野射辅甚至五人联动还能抗压不死，同时还要求发育起来打输出！</el-text>
            <el-text :style="{ marginLeft: '35px', fontSize: '16px' ,flex: '1 0 calc(100% - 35px)' }">你们都这么要求了，加强射手不是应该的吗！？怎么有些人就开始闹起来了</el-text>
            <div class="commentFooter" :style="{ marginLeft: '35px' }">
              <el-text size="small" :style="{ color: 'rgb(148, 153, 160)', marginTop: '5px' }">2021-03-22
                11:29</el-text>
              <div :style="{ display: 'flex', alignItems: 'center', marginLeft: '20px', marginTop: '8px' }">
                <el-icon size="15px" color='rgb(148, 153, 160)'>
                  <SvgIcon iconName="icon-zan" />
                </el-icon>
                <el-text size="small"
                  :style="{ color: 'rgb(148, 153, 160)', marginLeft: '3px', marginBottom: '3px' }">121</el-text>
                <el-icon size="15px" color='rgb(148, 153, 160)' :style="{ marginLeft: '20px', marginBottom: '3px' }">
                  <SvgIcon iconName="icon-cai" />
                </el-icon>
                <el-text size="small" class="replyBtn">回复</el-text>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { Message as ElIconMessage } from '@element-plus/icons-vue'
import home_navbar from '@/components/homeNavbar.vue'
import axios from 'axios'
import video_player from '@/components/videoPlayer.vue'
import { transform } from 'typescript';
export default {
  data() {
    return {
      newBarrge: '',
      videoID: null,
      BarrageList: [
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 1,
        },
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 2,
        },
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 3,
        },
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 4,
        },
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 5,
        },
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 6,
        },
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 7,
        },
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 8,
        },
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 9,
        },
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 10,
        },
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 11,
        },
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 12,
        },
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 13,
        },
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 14,
        },
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 15,
        },
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 16,
        },
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 17,
        },
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 18,
        },
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 19,
        },
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 20,
        },
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 21,
        },
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 22,
        },
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 23,
        },
        {
          time: '0:00',
          date: '03-22 11:29',
          content: '这是一条弹幕',
          id: 24,
        },
      ],
      ElIconMessage,
      videoTags: ['标签1', '标签2', '标签3'],
      newestChecked: false,
      hotestChecked: false
    }
  },
  components: {
    home_navbar,
    video_player,
  },
  mounted() {
    this.videoID = this.$route.params.id;
  }
}
</script>

<style scoped>
* {
  margin: 0;
  padding: 0;
}

.home_navbar {
  height: 30px;
}

.video_detail {
  display: grid;
  margin-top: 40px;
  grid-template-columns: 1120px 410px;
  margin-left: 200px;
  grid-template-rows: 50px 630px;
  column-gap: 30px;
  row-gap: 26px;
}

.el-divider {
  position: relative;
  transform: translateY(30px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.video-info {
  margin-top: 0px;
}

.icons {
  width: 17px;
  height: 13px;
}

.title {
  font-size: 23px;
  margin-bottom: 6px;
}

.video_info {
  width: 450px;
  height: 30px;
  font-size: 13px;
  color: rgb(148, 153, 160);
  display: grid;
  grid-auto-flow: column;
  grid-template-columns: 65 50 120 160;
  column-gap: 0px;
  align-items: center;
}

.info_item {
  display: flex;
  align-items: center;
  column-gap: 5px;
}

.text {
  vertical-align: middle;
}

.bigAvatar {
  width: 50px;
  height: 50px;
}

.authorPart {
  display: grid;
  grid-template-columns: 50px 325px;
  column-gap: 25px;
}

.rightPart {
  display: grid;
  grid-template-rows: 20px 15px 30px;
  row-gap: 5px;
  height: 75px;
}

.authorName {
  display: flex;
  font-size: 15px;
  color: rgb(230, 25, 103);
  letter-spacing: 1px;
}

.authorTop {
  display: grid;
  grid-template-columns: 80px 60px;
  column-gap: 10px;
}

.el-link {
  align-items: center;
  display: grid;
  grid-template-columns: 14px 50px;
  font-size: 14px;
}

.authorDes {
  font-size: 13px;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  font-weight: 50;
}

.authorFooter {
  height: 30px;
  display: grid;
  grid-template-columns: 75px 230px;
}

.authorFooter .el-button {
  height: 25px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.barrageList {
  width: 400px;
  height: 100px;
}

.barrageTop {
  display: grid;
  grid-template-columns: 50px 235px 65px;
  column-gap: 15px;
}

.listItem {
  height: 20px;
  display: grid;
  grid-template-columns: 50px 235px 65px;
  column-gap: 15px;
}

.barrageTop p {
  font-size: 12px;
  font-weight: 50;
}

.listItem p {
  font-size: 12px;
  font-weight: 200;
}

.barrageContext {
  display: grid;
  row-gap: 5px;
  height: 590px;
  overflow: auto;
}


.videoRest {
  display: grid;
  grid-template-rows: 50px 65px 1fr;
  width: 100%;
  height: 600px;
  row-gap: 0px;
}

.videoFooter {
  width: 1100px;
  height: 50px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  display: flex;
  height: 50px;
  margin-top: -10px;
  border-radius: 5px;
}

.videoFooterLeft {
  font-size: 14px;
  font-weight: 150;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-zihao:hover {
  color: rgb(0, 174, 236);
  cursor: pointer;
}

.interactions {
  width: 1100px;
  height: 65px;
  display: grid;
  grid-template-columns: 85px 85px 85px 85px;
  column-gap: 30px;
  align-items: center;
  border-bottom: 1px solid #e8e8e8;
}

.footerBtn {
  width: 85px;
  height: 45px;
  display: flex;
  align-items: center;
  margin-left: 20px;
}

.footerBtn .el-link {
  display: flex;
  align-items: center;
  justify-content: center;
  padding-left: 10px;
}

.footerBtn .el-link .el-icon {
  margin-right: 5px;
}

.videoDescription {
  border-bottom: 1px solid #e8e8e8;
  max-height: 150px;
  width: 1100px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  margin-top: 20px;
  /* 设置高度为100%以填充整个父容器 */
}

.description {
  max-height: 90px;
  /* 设置该元素在垂直方向上占据剩余空间 */
  overflow-y: auto;
  /* 当文本内容超出容器高度时显示滚动条 */
}

.videoTags {
  display: flex;

  /* 标签向左对齐 */
  align-items: center;
  margin-bottom: 20px;
}

.commentsContainer {
  width: 1100px;
  border: 1px solid #e8e8e8;
  margin-bottom:20px;
}

.commentsTabs{
  display: flex;
  align-items: center;
}

 .AddComment{
  margin-top: 20px;
  margin-left:10px;
 }
 
::v-deep  .el-input div{
  background-color: rgb(241,242,243) !important;
}

.comments{
  display: flex;
  border:1px solid #e8e8e8;
  margin-left: 10px;
  align-items: center;
  flex-wrap: wrap;
}

 .commentFooter{
  display:flex;
  align-items: center;
  margin-left:55px;
 }

 .replyBtn:hover{
  color:rgb(0, 174, 236);
  cursor: pointer;
 }

 .replyBtn{
  margin-bottom: 3px;
  margin-left: 20px;
  color:rgb(148, 153, 160)
 }

  .reply{
    margin-left: 50px;
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    border: 1px solid #e8e8e8;
  }
</style>
