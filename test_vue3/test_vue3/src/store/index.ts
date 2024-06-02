import { createStore } from 'vuex'

interface IRootState {
    token: string,
    isLogin: boolean
}

const store = createStore<IRootState>({   //定义state的类型限定
	
})

export default store 
