<template>
	<view class="page">

		<!-- 图片预览区 -->
		<view class="preview-area" @click="chooseFromAlbum">
			<image v-if="imagePath" :src="imagePath" mode="aspectFit" class="preview-img"></image>
			<view v-else class="placeholder">
				<text class="plus">＋</text>
				<text class="tip">点击上传图片</text>
			</view>
		</view>

		<!-- 操作按钮 -->
		<view class="btn-group">
			<button class="btn camera" @click="takePhoto">📷 拍照</button>
			<button class="btn album" @click="chooseFromAlbum">🖼 从相册选择</button>
		</view>
		<button class="btn recognize" :disabled="!imagePath || loading" @click="startRecognize">
			{{ loading ? '识别中…' : '开始识别' }}
		</button>

		<!-- 识别结果 -->
		<view class="result-area" v-if="result">
			<view class="result-title">识别结果</view>
			<view class="result-content">
				<rich-text :nodes="markdownToHtml(result)"></rich-text>
			</view>
		</view>
	</view>
</template>

<script setup>
import {ref} from 'vue'
import {uploadFile} from '@/util/request.js'
import { onPullDownRefresh } from '@dcloudio/uni-app'

const imagePath = ref('')      // 选中图片本地路径
const loading = ref(false)     // 识别中
const result = ref('')         // 识别结果

// 下拉刷新：工具页无可刷数据，直接结束刷新动画
onPullDownRefresh(() => {
	uni.stopPullDownRefresh()
})

// 拍照
const takePhoto = ()=>{
	uni.chooseImage({
		count: 1,
		sourceType: ['camera'],
		success: (res) => {
			imagePath.value = res.tempFilePaths[0]
			result.value = ''
		}
	})
}

// 从相册选择
const chooseFromAlbum = ()=>{
	uni.chooseImage({
		count: 1,
		success: (res) => {
			imagePath.value = res.tempFilePaths[0]
			result.value = ''
		}
	})
}

// 开始识别
const startRecognize = async ()=>{
	if (!imagePath.value || loading.value) return
	loading.value = true
	result.value = ''
	uni.showLoading({ title: '识别中…', mask: true })
	try {
		const res = await uploadFile({
			url: '/api/agent/map/recognize',
			filePath: imagePath.value,
			name: 'file',
			formData: { account: uni.getStorageSync('user_account') },
			timeout: 120000
		})
		// Go 端返回 { response: { code, result } }
		const answer = (res && res.response && res.response.result)
			? res.response.result
			: '识别失败，请重试'
		result.value = answer
	} catch (err) {
		console.log('识别失败:', err)
		const msg = (err && err.message) ? err.message : '识别失败'
		uni.showToast({ title: msg.length > 18 ? msg.slice(0, 18) : msg, icon: 'none' })
	} finally {
		loading.value = false
		uni.hideLoading()
	}
}

// markdown → HTML
const markdownToHtml = (content = '')=>{
	let html = content
		.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
	html = html
		.replace(/^### (.*)$/gm, '<h4>$1</h4>')
		.replace(/^## (.*)$/gm, '<h3>$1</h3>')
		.replace(/^# (.*)$/gm, '<h2>$1</h2>')
		.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
		.replace(/`([^`]+)`/g, '<code>$1</code>')
		.replace(/^[-*] (.*)$/gm, '• $1')
		.replace(/^(\d+)\. (.*)$/gm, '<br/>$1. $2')
		.replace(/\n/g, '<br/>')
	return html
}
</script>

<style lang="scss" scoped>
.page{
	min-height: 100vh;
	background: #f4f6f9;
	padding: 24rpx;
	box-sizing: border-box;
}
.preview-area{
	width: 100%;
	height: 480rpx;
	background: #ffffff;
	border-radius: 24rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	overflow: hidden;
	border: 2rpx dashed #a9c4ff;
	box-shadow: 0 8rpx 24rpx rgba(74, 125, 255, 0.08);
}
.preview-img{
	width: 100%;
	height: 100%;
}
.placeholder{
	display: flex;
	flex-direction: column;
	align-items: center;
	color: #9aa5b4;
}
.plus{
	font-size: 80rpx;
	font-weight: 200;
	line-height: 1;
}
.tip{
	font-size: 26rpx;
	margin-top: 12rpx;
}
.btn-group{
	display: flex;
	gap: 20rpx;
	margin-top: 24rpx;
}
.btn{
	flex: 1;
	height: 84rpx;
	line-height: 84rpx;
	margin: 0;
	padding: 0;
	font-size: 30rpx;
	border-radius: 20rpx;
	&::after{ border: none; }
	&:active{ opacity: 0.85; }
}
.btn.camera{
	background: linear-gradient(135deg, #4a7dff, #6ea3ff);
	color: #fff;
	box-shadow: 0 8rpx 20rpx rgba(74, 125, 255, 0.25);
}
.btn.album{
	background: #ffffff;
	color: #333;
	border: 2rpx solid #e5e7eb;
	box-shadow: 0 4rpx 16rpx rgba(29, 33, 41, 0.04);
}
.btn.recognize{
	width: 100%;
	flex: none;
	margin-top: 24rpx;
	background: linear-gradient(135deg, #4b8cff, #2f6bff);
	color: #fff;
	box-shadow: 0 8rpx 20rpx rgba(47, 107, 255, 0.3);
	&[disabled]{
		background: #e9edf3;
		color: #b6c0cf;
		box-shadow: none;
	}
}
.result-area{
	margin-top: 24rpx;
	background: #ffffff;
	border-radius: 24rpx;
	padding: 24rpx;
	box-shadow: 0 8rpx 24rpx rgba(29, 33, 41, 0.05);
}
.result-title{
	font-size: 30rpx;
	font-weight: 600;
	color: #222;
	margin-bottom: 16rpx;
	padding-left: 16rpx;
	border-left: 8rpx solid #3a7afe;
}
.result-content{
	font-size: 28rpx;
	line-height: 1.7;
	color: #444;
	word-break: break-all;
}
</style>
