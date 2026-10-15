<template>
	<view class="page">
		<view class="header">
			<text class="title">AI 周报</text>
			<text class="sub">基于本周公告与成员健康咨询自动生成</text>
		</view>

		<button class="gen-btn" :disabled="loading" @click="generate">
			{{ loading ? '生成中…' : '🔄 生成本周周报' }}
		</button>

		<view class="report-area" v-if="report">
			<rich-text :nodes="markdownToHtml(report)"></rich-text>
		</view>
		<view class="empty" v-else-if="!loading">
			<text>点击上方按钮生成周报</text>
		</view>
	</view>
</template>

<script setup>
import {ref} from 'vue'
import {request} from '@/util/request.js'
import { onPullDownRefresh } from '@dcloudio/uni-app'

const loading = ref(false)
const report = ref('')

// 下拉刷新：周报需手动生成，直接结束刷新动画
onPullDownRefresh(() => {
	uni.stopPullDownRefresh()
})

const generate = async ()=>{
	if (loading.value) return
	loading.value = true
	report.value = ''
	uni.showLoading({ title: 'AI 生成中…', mask: true })
	try {
		const res = await request({
			url: '/api/report/generate',
			method: 'POST',
			data: {},
			timeout: 120000
		})
		// Go 返回 { report: { code, report } }
		const text = (res && res.report && res.report.report) ? res.report.report : '生成失败，请重试'
		report.value = text
	} catch (err) {
		console.log('生成周报失败:', err)
		uni.showToast({ title: '生成失败', icon: 'none' })
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
.header{
	padding: 16rpx 8rpx 24rpx;
}
.title{
	display: block;
	font-size: 40rpx;
	font-weight: 700;
	color: #222;
}
.sub{
	display: block;
	font-size: 26rpx;
	color: #8a94a6;
	margin-top: 8rpx;
}
.gen-btn{
	width: 100%;
	height: 88rpx;
	line-height: 88rpx;
	margin: 0;
	padding: 0;
	font-size: 30rpx;
	color: #fff;
	background: linear-gradient(135deg, #4a7dff, #6ea3ff);
	border-radius: 999rpx;
	box-shadow: 0 8rpx 20rpx rgba(74, 125, 255, 0.3);
	&::after{ border: none; }
	&:active{ opacity: 0.85; }
	&[disabled]{
		background: #e9edf3;
		color: #b6c0cf;
		box-shadow: none;
	}
}
.report-area{
	margin-top: 24rpx;
	background: #fff;
	border-radius: 24rpx;
	padding: 28rpx;
	box-shadow: 0 8rpx 24rpx rgba(29, 33, 41, 0.05);
	font-size: 28rpx;
	line-height: 1.8;
	color: #333;
	word-break: break-all;
}
.empty{
	margin-top: 100rpx;
	text-align: center;
	color: #9aa5b4;
	font-size: 26rpx;
}
</style>
