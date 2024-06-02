<template>
    <div>
        <homeNavbar></homeNavbar>
        <div><el-divider></el-divider></div>
        <div class="search-container">
            <el-input placeholder="请输入关键字搜索" v-model="searchText" clearable @keydown.native.enter="handleSearch" >
                <template #prefix>
                    <el-icon color="rgb(50,192,240)" size="20px">
                        <Search />
                    </el-icon>
                </template>
                <template #suffix>
                    <el-button
                        :style="{ 'background-color': 'rgb(50,192,240)', width: '100px', height: '40px', color: 'white' }"@click="handleSearch">搜索</el-button>
                </template>
            </el-input>
        </div>

        <el-tabs v-model="activeName" :style="{ 'width': '90vw', height: 'calc(100vh)' }">
            <el-tab-pane name="first">
                <template #label>
                    <span style="font-size: 17px;margin-right: 30px,padding-left 10px;">综合</span>
                </template>
                <div :style="{ marginTop: '25px' }">
                    <el-button v-for="(button, index) in buttons" :key="index"
                        :class="{ 'selected': selectedButton === index }" @click="selectButton(index,'videoSort')" text
                        :size="large" :style="{ paddingLeft: '10px', paddingRight: '10px',height: '30px' }" >{{
                        button.text
                        }}</el-button>
                </div>
                <div class="videoList">
                    <div class="recommend-video" v-for="video in videoList1.slice(0, 6)" :key="video.VID">
                        <!-- 根据videoList中的数据循环展示推荐视频 -->
                        <el-image :src="video.Cover" class="image-with-overlay">
                        </el-image>
                        <div class="video-info">
                            <el-button text :style="{ color: 'white', marginLeft: '10px'}">
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
                            <el-button text :style="{ color: 'white', marginLeft: 'calc(7vw )' }">
                                {{ video.Duration }}
                            </el-button>
                        </div>
                        <div class="video-title">
                            <el-link :underline="false" :style="{ fontSize: '14px', color: 'black' }"
                                :href="`#/video/VideoID=${video.VID}`">{{ video.Title }}</el-link>
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
                <div class="pagination">
                    <el-pagination background layout="prev, pager, next" :total="totalPage" />
                </div>
            </el-tab-pane>
            <el-tab-pane name="second">
                <template #label>
                    <el-badge :value="100" :max="99" :offset="[15, 10]" color="rgb(241,242,243)"
                        :badge-style="{ color: 'black' }" :style="{ marginRight: '30px' }">
                        <span style="font-size: 17px;">视频</span></el-badge>
                </template>
                <div>
                    <div :style="{ marginTop: '25px' }">
                        <el-button v-for="(button, index) in buttons" :key="index"
                            :class="{ 'selected': selectedButton === index }" @click="selectButton(index,'videoSort')" text
                            :size="large" :style="{ paddingLeft: '10px', paddingRight: '10px', height: '30px' }">{{
                            button.text
                            }}</el-button>
                    </div>
                </div>
                <div class="videoList">
                    <div class="recommend-video" v-for="video in videoList1.slice(0, 6)" :key="video.VID">
                        <!-- 根据videoList中的数据循环展示推荐视频 -->
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
                            <el-button text :style="{ color: 'white', marginLeft: 'calc(7vw )' }">
                                {{ video.Duration }}
                            </el-button>
                        </div>
                        <div class="video-title">
                            <el-link :underline="false" :style="{ fontSize: '14px', color: 'black' }"
                                :href="`#/video/VideoID=53540871860293`">{{ video.Title }}</el-link>
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
                <div class="pagination">
                    <el-pagination background layout="prev, pager, next" :total="totalPage" />
                </div>
            </el-tab-pane>
            <el-tab-pane name="third">
                <template #label>
                    <el-badge :value="100" :max="99" :offset="[15, 10]" color="rgb(241,242,243)"
                        :badge-style="{ color: 'black' }" :style="{ marginRight: '30px' }">
                        <span style="font-size: 17px;">用户</span></el-badge>
                </template>
                <div :style="{ marginTop: '25px' }">
                    <el-button v-for="(button, index) in buttonsUser" :key="index"
                        :class="{ 'selected': selectedButton === index }" @click="selectButton(index,'userSort')" text
                        :size="large" :style="{ paddingLeft: '10px', paddingRight: '10px', height: '30px' }">{{
                        button.text
                        }}</el-button>
                </div>

                <div class="userList">
                    <div class="search-user" >
                        <el-avatar :size="80">user</el-avatar>
                        <div class="user-info">
                            <div :style="{ display: 'flex',alignItems: 'center' }">
                                <el-text
                                    :style="{ fontSize: '18px', fontWeight: 'bold',color: 'black' }">这是用户名</el-text>
                                <el-icon size="25px" :style="{ marginLeft: '10px' }">
                                    <SvgIcon iconName="icon-ic_userlevel_6" />
                                </el-icon>
                            </div>
                            <div>
                                <el-text :style="{ fontSize: '13px', color: 'rgb(97,102,109)' }">2103粉丝</el-text>
                                <el-text :style="{ fontSize: '13px', color: 'rgb(97,102,109)',marginLeft: '10px' }">
                                    77个视频</el-text>
                                <el-text :style="{ fontSize: '13px', color: 'rgb(97,102,109)', marginLeft: '10px' }">
                                    你无敌了孩子</el-text>
                            </div>
                            <div>
                                <el-button
                                    :style="{ backgroundColor: 'rgb(0,174,236)', color: 'white', width: '100px', height: '35px',marginTop: '5px',borderRadius: '5px' }">+关注</el-button>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="pagination">
                    <el-pagination background layout="prev, pager, next" :total="totalPage" />
                </div>
            </el-tab-pane>
        </el-tabs>
    </div>
</template>

<script>
import homeNavbar from "@/components/homeNavbar.vue";
import { Search } from '@element-plus/icons-vue';
import jsonData from '@/assets/response_1716983420154.json';
import axios from "axios";

export default {
    components: {
        homeNavbar,
        Search
    },
    data() {
        return {
            searchText: "",
            activeName: "first",
            buttons: [
                { text: '综合升序' },
                { text: '最多播放' },
                { text: '最新发布' },
                { text: '最多弹幕' },
                { text: '最多收藏' }
            ],
            buttonsUser: [
                { text: '默认排序' },
                { text: '粉丝数由高到低' },
                { text: '粉丝数由低到高' },
                { text: 'Lv等级由高到低' },
                { text: 'Lv等级由低到高' }
            ],
            videoSort: ['default',  'mostPlay','newest', 'mostBarrage','mostFavorite'],
            videoList1: [],
            selectedButton: 0,
            totalPage: 280,
            keyword: "",
            type: 'videoSort',
            sortOrder: 'default'
        }
    },
    methods: {
        selectButton(index,string) {
            this.selectedButton = index;
            this.type = string
          console.log(index)
            if (this.type === 'videoSort') {
                this.sortOrder = this.videoSort[index];
                console.log(this.sortOrder);
                const token = localStorage.getItem('token');
              axios.get('/yanxi/video/SearchVideos', {
                headers: {
                  'Authorization': token
                },
                params: {
                  videoNums: 10,
                  offset: 0,
                  key: this.keyword,
                  sortOrder: this.sortOrder

                }
              }).then(res => {
                const videodata = res.data.data;
                this.videoList1 = videodata.map(item => ({
                  ...item,
                  Cover: `data:image/png;base64,${item.Cover}`,
                  Duration: this.formatTime(item.Duration)
                }));
                console.log(videodata)
              })
            }
        },
        formatTime(timeString) {
            if (timeString.startsWith('00:')) {  // 检查时间字符串是否以 "00:" 开头
                return timeString.substring(3);  // 去除代表小时的数字以及后面的冒号
            }
            return timeString;  // 如果不是以 "00:" 开头，则保持原样返回
        },
      handleSearch() {
          console.log(this.searchText);
          window.location.href = `#/search/keyword=${this.searchText}`;
        window.location.reload();
      }
    },
    mounted() {
        const data = jsonData.data;
        const token = localStorage.getItem('token');
        this.keyword = this.$route.params.keyword;
        console.log(this.keyword);

            if (this.type === 'videoSort') {
                axios.get('/yanxi/video/SearchVideos', {
                    headers: {
                        'Authorization': token
                    },
                    params: {
                        videoNums: 10,
                        offset: 0,
                        key: this.keyword,
                        sortOrder: this.sortOrder

                    }
                }).then(res => {
                    const videodata = res.data.data;
                  this.videoList1 = videodata.map(item => ({
                    ...item,
                    Cover: `data:image/png;base64,${item.Cover}`,
                    Duration: this.formatTime(item.Duration)
                  }));
                     console.log(videodata)
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
.recPlate {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-start;
  height: auto;
}
.homeNavbar {
    height: 20px;
}

.el-divider {
    position: relative;
    transform: translateY(10px);
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.search-container {
    margin-top: 10px;
    display: flex;
    height: 110px;
    justify-content: center;
    align-items: center;
}

.el-input {
    width: 650px;
    height: 50px;
    font-size: 16px;
}

::v-deep .el-input div {
    background-color: rgb(246, 247, 248);
}

::v-deep .el-input:focus div {
    background-color: white;
}

.el-tabs {
    margin-left: 100px;
}

.selected {
    background-color: rgb(223, 246, 253);
    color: rgb(0, 174, 236);
}

.el-button:hover {
    color: rgb(0, 174, 236);
    background-color: inherit !important;
}

.recommend-video {
    margin-top:15px;
    display: flex;
    flex-wrap: wrap;
    height:230px;
    width:270px;
    margin-right:10px;
}

.el-image {
    width: 270px;
    height: 150px;
    padding-right: 20px;
    border-radius: 5px;
}

.video-title {
    height: 45px;
}

.video-title .el-link {
    height: 50px;
    width: 230px;
    font-size: 15px;
    display: flex;
    justify-content: flex-start;
    flex-wrap: wrap;
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
    background: linear-gradient(to bottom, rgba(0, 0, 0, 0), rgba(0, 0, 0, 0.7));
}

.video-info {
    position: absolute;
    transform: translateY(calc(13vh));
}

.videoList{
    display: flex;
    flex-wrap: wrap;
    align-items: center;
}
.pagination{
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100px;
}

.el-pagination{
  width:200px;
}

.userList{
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    margin-top: 40px;
}

.search-user{
    display: flex;
    align-items: center;
    margin-right: 200px;
}

.user-info{
    margin-left: 10px;
    display: grid;
    row-gap: 5px;
    grid-template-rows: 25px 15px 55px;
}
</style>
