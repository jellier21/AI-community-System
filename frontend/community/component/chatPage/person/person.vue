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
import { onLoad } from '@dcloudio/uni-app'
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
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 100rpx 40rpx 60rpx;
	background: #fff;

	.avatar {
		width: 160rpx;
		height: 160rpx;
		border-radius: 50%;
		background: #f0f0f0;
		margin-bottom: 28rpx;
	}
	.nickname {
		font-size: 40rpx;
		font-weight: 600;
		color: #1a1a1a;
		margin-bottom: 12rpx;
	}
	.account {
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
