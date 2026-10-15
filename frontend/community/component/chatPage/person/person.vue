<template>
	<view class="page">
		<view class="profile-card">
			<image class="avatar" :src="profile.avatar || '/static/logo.png'" mode="aspectFill"></image>
			<view class="nickname">{{ profile.nickname || '未设置昵称' }}</view>
			<view class="account" @click="copyAccount">账号：{{ account }} <text class="copy-icon">点击复制</text></view>

			<view class="signature" v-if="profile.signature">{{ profile.signature }}</view>
			<view class="signature placeholder" v-else>这个人很懒，什么都没写</view>
		</view>
	</view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad, onPullDownRefresh } from '@dcloudio/uni-app'
import { request } from '@/util/request.js'

const account = ref('')
const profile = ref({
	nickname: '',
	avatar: '',
	signature: ''
})

onLoad((option) => {
	account.value = option.account || ''
	if (account.value) {
		loadProfile()
	}
})

// 下拉刷新：重新拉取对方资料
onPullDownRefresh(async () => {
	if (account.value) {
		await loadProfile()
	}
	uni.stopPullDownRefresh()
})

const loadProfile = async () => {
	try {
		const res = await request({
			url: '/api/user/query/profile',
			data: { account: account.value }
		})
		if (res.code === 200 && res.profile) {
			profile.value = res.profile
		}
	} catch (e) {
		console.error('查询用户资料失败:', e)
	}
}
const copyAccount = () => {
	uni.setClipboardData({
		data: account.value,
		success: () => {
			uni.showToast({ title: '账号已复制', icon: 'success' })
		}
	})
}

</script>

<style lang="scss" scoped>
page {
	background: #f5f6f7;
}
.page {
	min-height: 100vh;
}
.profile-card {
	position: relative;
	display: flex;
	flex-direction: column;
	align-items: center;
	margin: 24rpx;
	padding: 80rpx 40rpx 60rpx;
	background: #fff;
	border-radius: 28rpx;
	box-shadow: 0 8rpx 24rpx rgba(29, 33, 41, 0.06);
	overflow: hidden;

	/* 顶部渐变装饰 */
	&::before {
		content: '';
		position: absolute;
		left: 0;
		right: 0;
		top: 0;
		height: 200rpx;
		background: linear-gradient(135deg, #4a7dff, #7aa8ff);
	}

	.avatar {
		position: relative;
		width: 160rpx;
		height: 160rpx;
		border-radius: 50%;
		background: #f0f0f0;
		border: 6rpx solid #ffffff;
		box-shadow: 0 8rpx 20rpx rgba(74, 125, 255, 0.25);
		margin-bottom: 28rpx;
	}
	.nickname {
		position: relative;
		font-size: 40rpx;
		font-weight: 600;
		color: #1a1a1a;
		margin-bottom: 12rpx;
	}
	.account {
		position: relative;
		font-size: 26rpx;
		color: #9a9fa8;
		margin-bottom: 24rpx;
	}
	.copy-icon {
		font-size: 22rpx;
		color: #4a7dff;
		margin-left: 12rpx;
	}

	.signature {
		position: relative;
		font-size: 28rpx;
		color: #6b7280;
		max-width: 80%;
		text-align: center;
		line-height: 1.6;
	}
	.signature.placeholder {
		color: #c0c4cc;
	}
}
</style>
