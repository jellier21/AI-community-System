<template>
	<view class="page">
		<!-- 社区头像上传区 -->
		<view class="avatar-upload" @click="chooseAvatar">
			<image v-if="avatarPath" :src="avatarPath" mode="aspectFill" class="avatar-img"></image>
			<view v-else class="avatar-placeholder">
				<image src="/static/chat/push.png" mode="aspectFit"></image>
				<text>上传社区头像</text>
			</view>
		</view>

		<!-- 社区信息表单 -->
		<view class="form-card">
			<view class="form-item">
				<text class="form-label">社区名称</text>
				<input
					class="form-input"
					type="text"
					v-model="communityName"
					placeholder="请输入社区名称"
					:maxlength="20"
				/>
			</view>
			<view class="form-item">
				<text class="form-label">社区介绍</text>
				<textarea
					class="form-textarea"
					v-model="description"
					placeholder="介绍一下你的社区（选填）"
					:maxlength="50"
				/>
			</view>
		</view>

		<!-- 创建按钮 -->
		<button class="create-btn" :disabled="submitting" @click="createCommunity">
			{{ submitting ? '创建中...' : '创建社区' }}
		</button>
	</view>
</template>

<script setup>
	import { ref } from 'vue';
	import { onLoad, onPullDownRefresh } from '@dcloudio/uni-app';
	import { BASEURL_NAME } from '@/util/configuration.js';

	const account = ref('')
	const communityName = ref('')
	const description = ref('')
	const avatarPath = ref('')   // 头像仅本地预览（后端创建接口暂无头像字段）
	const submitting = ref(false)

	onLoad(() => {
		account.value = uni.getStorageSync('user_account') || ''
	})

	// 下拉刷新：表单页无数据可刷，直接结束刷新动画
	onPullDownRefresh(() => {
		uni.stopPullDownRefresh()
	})

	// 选择社区头像（本地预览）
	const chooseAvatar = () => {
		uni.chooseImage({
			count: 1,
			sizeType: ['compressed'],
			success: (res) => {
				avatarPath.value = res.tempFilePaths[0]
			}
		})
	}

	// 创建社区：POST /api/community/create（JWT 鉴权）
	const createCommunity = () => {
		const name = communityName.value.trim()
		if (!name) {
			uni.showToast({ title: '请输入社区名称', icon: 'none' })
			return
		}
		if (submitting.value) return
		submitting.value = true

		const token = uni.getStorageSync('jwtToken')
		uni.request({
			url: `${BASEURL_NAME}/api/community/create`,
			method: 'POST',
			header: {
				'Content-Type': 'application/x-www-form-urlencoded',
				...(token ? { Authorization: token } : {})
			},
			data: {
				account: account.value,
				name: name,
				description: description.value.trim()
			},
			success: (res) => {
				if (res.statusCode === 200) {
					uni.showToast({ title: '创建成功', icon: 'success' })
					setTimeout(() => {
						uni.navigateBack()
					}, 800)
				} else if (res.statusCode === 401) {
					uni.removeStorageSync('jwtToken')
					uni.reLaunch({ url: '/pages/index/index' })
				} else {
					// 展示后端具体错误（如"该用户已存在社区"）
					uni.showToast({ title: res.data?.error || '创建失败', icon: 'none' })
				}
			},
			fail: () => {
				uni.showToast({ title: '网络异常', icon: 'none' })
			},
			complete: () => {
				submitting.value = false
			}
		})
	}
</script>

<style lang="scss" scoped>
	page {
		background-color: #f5f6f7;
	}

	.page {
		min-height: 100vh;
		padding: 60rpx 40rpx;
		box-sizing: border-box;
	}

	/* ---------- 头像上传区 ---------- */
	.avatar-upload {
		width: 200rpx;
		height: 200rpx;
		margin: 0 auto 60rpx;
		border-radius: 50%;
		background: #fff;
		display: flex;
		align-items: center;
		justify-content: center;
		overflow: hidden;
		box-sizing: border-box;
		border: 2rpx dashed #a9c4ff;
		box-shadow: 0 8rpx 24rpx rgba(74, 125, 255, 0.12);

		.avatar-img {
			width: 100%;
			height: 100%;
		}

		.avatar-placeholder {
			display: flex;
			flex-direction: column;
			align-items: center;

			image {
				width: 64rpx;
				height: 64rpx;
				margin-bottom: 8rpx;
			}

			text {
				font-size: 22rpx;
				color: #9a9fa8;
			}
		}
	}

	/* ---------- 表单 ---------- */
	.form-card {
		background: #fff;
		border-radius: 24rpx;
		padding: 0 32rpx;
		margin-bottom: 60rpx;
		box-shadow: 0 8rpx 24rpx rgba(29, 33, 41, 0.05);

		.form-item {
			padding: 32rpx 0;
			border-bottom: 1rpx solid #f0f0f0;

			&:last-child {
				border-bottom: none;
			}

			.form-label {
				display: block;
				font-size: 28rpx;
				color: #1a1a1a;
				margin-bottom: 16rpx;
			}

			.form-input {
				width: 100%;
				height: 72rpx;
				padding: 0 24rpx;
				background: #f5f6f7;
				border: 2rpx solid #eef0f3;
				border-radius: 16rpx;
				font-size: 28rpx;
				color: #1a1a1a;
				box-sizing: border-box;
			}

			.form-textarea {
				width: 100%;
				height: 160rpx;
				padding: 20rpx 24rpx;
				background: #f5f6f7;
				border: 2rpx solid #eef0f3;
				border-radius: 16rpx;
				font-size: 28rpx;
				color: #1a1a1a;
				box-sizing: border-box;
			}
		}
	}

	/* ---------- 创建按钮 ---------- */
	.create-btn {
		width: 100%;
		height: 92rpx;
		line-height: 92rpx;
		margin: 0;
		padding: 0;
		font-size: 32rpx;
		font-weight: 500;
		color: #fff;
		background: linear-gradient(135deg, #4a7dff, #6ea3ff);
		border-radius: 999rpx;
		box-sizing: border-box;
		box-shadow: 0 10rpx 24rpx rgba(74, 125, 255, 0.3);

		&::after {
			border: none;
		}

		&[disabled] {
			background: #c8d4ff;
			color: #fff;
			box-shadow: none;
		}
	}
</style>
