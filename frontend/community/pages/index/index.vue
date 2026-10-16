<template>
	<view class="container">
		<view class="content">

			<view class="app-title">邻友小圈</view>
			<view class="app-subtitle">连接邻里 · 服务生活</view>

			<view class="status_card">
				<!-- 切换 -->
				<view class="card-tabs">
					<view class="tab" :class="{active: card_status}" @click="card_status = true">登录</view>
					<view class="tab" :class="{active: !card_status}" @click="card_status = false">注册</view>
				</view>

				<!-- 登录 -->
				<view class="login_card" v-if="card_status">
					<input type="text" placeholder="请输入您的账号" v-model="L_account"/>
					<input type="password" placeholder="请输入您的密码" v-model="L_password"/>
					<button @click="LoginSkipHome">登录</button>
					<text class="switch-tip" @click="switchCardStatus">没有账号？点我注册···</text>
				</view>

				<!-- 注册 -->
				<view class="register_card" v-else>
					<view class="avatar-picker" @click="chooseAvatar">
						<image v-if="avatarPath" :src="avatarPath" mode="aspectFill"></image>
						<text v-else class="avatar-placeholder">选择头像</text>
					</view>
					<input type="text" placeholder="请输入您所要注册的账号" v-model="_r_account_"/>
					<input type="password" placeholder="请输入您的密码" v-model="_r_password_"/>
					<input type="password" placeholder="请再次输入您的密码" v-model="re_password"/>
					<button @click="RegisterUser">注册</button>
					<text class="switch-tip" @click="switchCardStatus">已有账号？点我登录···</text>
				</view>
			</view>

		</view>
	</view>
</template>

<script setup>
	import {ref} from 'vue'
	import {BASEURL_NAME} from "@/util/configuration"
	import { TurnToLogin, GoHome } from "@/util/tools.js"
	import { checkJwtTimer } from '@/util/request'
	import { request } from '@/util/request.js'
	import { uploadOneImage } from '@/util/tools/share.js'
	import { onPullDownRefresh } from '@dcloudio/uni-app'

	if(uni.getStorageSync('jwtToken')){
			console.log('检测到已有 <令牌>')
			;(async () => {
				const Error = await checkJwtTimer()
				if(!Error){
					GoHome()
				}
			})()
		}



	// if(uni.getStorageSync('jwtToken')){
	// 	// Temp code
	// 	GoHome()
	// }

	// 下拉刷新：重新校验登录态，Token 有效则直接进入首页
	onPullDownRefresh(async () => {
		if (uni.getStorageSync('jwtToken')) {
			const Error = await checkJwtTimer()
			if (!Error) {
				uni.stopPullDownRefresh()
				GoHome()
				return
			}
		}
		uni.stopPullDownRefresh()
	})

	const card_status = ref(true)
	const switchCardStatus = ()=>{
		card_status.value = !card_status.value
	}

	const L_account = ref('')
	const L_password = ref('')

	const LoginSkipHome = ()=>{
		if (!L_account.value || !L_password.value){
			uni.showToast({ title: '请输入账号和密码', icon: 'none' })
			return
		}
		uni.request({
			url: `${BASEURL_NAME}/api/user/login`,
			method: 'POST',
			header: {
				"Content-Type": "application/x-www-form-urlencoded"
			},
			data: {
				account: L_account.value,
				password: L_password.value
			},
			success: (res) => {
				console.log(res.data)
				if (res.data.code === 200){
					uni.showToast({ title: '登录成功', icon: 'success' })
					uni.setStorageSync('jwtToken', res.data.token)
					uni.setStorageSync('user_account', L_account.value)
					// 清除上一个账号残留的本地资料，避免串号
					uni.removeStorageSync('user_signature')
					uni.removeStorageSync('user_nickname')
					setTimeout(() => {
						uni.switchTab({
							url:'/pages/home/home',
							fail: (err) => {
								console.error('跳转失败:', err)
								uni.showToast({ title: '页面跳转失败，请重试', icon: 'none' })
							}
						})
					}, 600)
				}	else {
					uni.showToast({ title: res.data?.error || '登录失败', icon: 'none' })
				}
			},
			fail: (err) => {
				console.log(err)
				uni.showToast({ title: '网络异常', icon: 'none' })
			}
		})
	}

	const _r_account_ = ref('')
	const _r_password_ = ref('')
	const re_password = ref('')
	const avatarPath = ref('')

	const chooseAvatar = ()=>{
		uni.chooseImage({
			count: 1,
			success: (res)=>{
				avatarPath.value = res.tempFilePaths[0]
			}
		})
	}


	const RegisterUser = ()=>{
		if (!_r_account_.value || !_r_password_.value){
			uni.showToast({ title: '请输入账号和密码', icon: 'none' })
			return
		}
		if (_r_password_.value !== re_password.value){
			uni.showToast({ title: '两次输入的密码不一致', icon: 'none' })
			return
		}
		uni.request({
			url: `${BASEURL_NAME}/api/user/register`,
			method: 'POST',
			header: {
				"Content-Type": "application/x-www-form-urlencoded"
			},
			data: {
				account: _r_account_.value,
				password: _r_password_.value
			},
						success: async (res) => {
							console.log(res.data)
							if (res.data.code === 200){
								uni.showToast({ title: '注册成功', icon: 'success' })
								// 自动登录拿 token
								uni.request({
								    url: `${BASEURL_NAME}/api/user/login`,
								    method: 'POST',
								    header: {
								        "Content-Type": "application/x-www-form-urlencoded"
								    },
								    data: { account: _r_account_.value, password: _r_password_.value },
								    success: async (loginRes)=>{
										if(loginRes.data.code === 200){
											uni.setStorageSync('jwtToken', loginRes.data.token)
											uni.setStorageSync('user_account', _r_account_.value)
											// 新账号登录，清掉任何上一个账号的残留缓存
											uni.removeStorageSync('user_signature')
											uni.removeStorageSync('user_nickname')
											// 如果选了头像，上传七牛云
											if(avatarPath.value){
												try{
													uni.showLoading({ title: '上传头像...' })
													const avatarUrl = await uploadOneImage(avatarPath.value)
													await request({
														url: '/api/user/amend/avatar',
														method: 'POST',
														data: { account: _r_account_.value, avatar: avatarUrl }
													})
													uni.hideLoading()
												}catch(e){
													console.log('头像上传失败:', e)
													uni.hideLoading()
												}
											}

										}
										setTimeout(()=>{
											uni.switchTab({ url: '/pages/home/home' })
										}, 800)
									}
								})
							}	else {

					uni.showToast({ title: res.data?.error || '注册失败', icon: 'none' })
				}
			},
			fail: (err) => {
				console.log(err)
				uni.showToast({ title: '网络异常', icon: 'none' })
			}
		})
	}
</script>

<style scoped>
	page{
		height: 100%;
		background: #f5f6f8;
	}

	.container{
		min-height: 100vh;
		background: linear-gradient(180deg, #edf2ff 0%, #f5f6f8 60%, #f5f6f8 100%);
		box-sizing: border-box;
	}

	.content{
		min-height: 100vh;
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 0 48rpx;
		box-sizing: border-box;
	}

	/* ===== 标题 ===== */
	.app-title{
		margin-top: 120rpx;
		margin-bottom: 12rpx;
		font-size: 44rpx;
		font-weight: 700;
		color: #1f2329;
		letter-spacing: 2rpx;
		background: linear-gradient(135deg, #4a7dff, #7aa8ff);
		-webkit-background-clip: text;
		background-clip: text;
		-webkit-text-fill-color: transparent;
	}

	.app-subtitle{
		margin-bottom: 60rpx;
		font-size: 24rpx;
		color: #8a94a6;
		letter-spacing: 4rpx;
	}

	/* ===== 卡片 ===== */
	.status_card{
		width: 100%;
		background: #ffffff;
		border-radius: 24rpx;
		padding: 44rpx 40rpx 36rpx;
		box-sizing: border-box;
		box-shadow: 0 16rpx 48rpx rgba(74, 125, 255, 0.12);
	}

	/* ===== 切换 tabs（朴素下划线式） ===== */
	.card-tabs{
		display: flex;
		justify-content: center;
		margin-bottom: 40rpx;
	}
	.tab{
		position: relative;
		width: 160rpx;
		height: 72rpx;
		line-height: 72rpx;
		text-align: center;
		font-size: 30rpx;
		color: #8a94a6;
	}
	.tab.active{
		color: #1f2329;
		font-weight: 600;
	}
	.tab.active::after{
		content: '';
		position: absolute;
		left: 50%;
		bottom: 0;
		transform: translateX(-50%);
		width: 48rpx;
		height: 6rpx;
		border-radius: 6rpx;
		background: #3a7afe;
	}

	/* ===== 输入框 ===== */
	.login_card input,
	.register_card input{
		width: 100%;
		height: 88rpx;
		padding: 0 24rpx;
		box-sizing: border-box;
		margin-bottom: 24rpx;
		background: #ffffff;
		border: 1rpx solid #e4e7ec;
		border-radius: 12rpx;
		font-size: 28rpx;
		color: #1f2329;
	}
	.login_card input:focus,
	.register_card input:focus{
		border-color: #3a7afe;
	}

	/* ===== 按钮 ===== */
	.login_card button,
	.register_card button{
		width: 100%;
		height: 88rpx;
		line-height: 88rpx;
		margin: 8rpx 0 0;
		padding: 0;
		background: linear-gradient(135deg, #4a7dff, #6ea3ff);
		color: #ffffff;
		font-size: 30rpx;
		border-radius: 12rpx;
		box-shadow: 0 8rpx 20rpx rgba(74, 125, 255, 0.3);
	}
	.login_card button::after,
	.register_card button::after{ border: none; }
	.login_card button:active,
	.register_card button:active{
		background: linear-gradient(135deg, #2f6bff, #4a7dff);
	}

	/* ===== 切换提示 ===== */
	.switch-tip{
		display: block;
		text-align: center;
		font-size: 24rpx;
		color: #8a94a6;
		margin-top: 32rpx;
	}

	/* ===== 头像 ===== */
	.avatar-picker{
		width: 136rpx;
		height: 136rpx;
		border-radius: 50%;
		background: #f5f6f8;
		margin: 0 auto 32rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		overflow: hidden;
	}
	.avatar-picker image{ width: 100%; height: 100%; }
	.avatar-placeholder{ font-size: 22rpx; color: #9aa5b4; }
</style>
