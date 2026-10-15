<template>
	<view class="notice-container">
		<!-- 发布公告 -->
		<view class="issue-notice">
			<view class="section-title">发布新公告</view>
			<input class="input-title" v-model="newTitle" placeholder="请输入公告标题" />
			<textarea class="input-content" v-model="newContent" placeholder="请输入公告内容" />
			<button class="publish-btn" @click="handlePublish">发布</button>
		</view>

		<!-- 历史公告 -->
		<view class="notice-history">
			<view class="section-title">历史公告</view>
			<view class="notice-item" v-for="item in noticeList" :key="item.uid">
				<view class="notice-title">{{ item.title }}</view>
				<view class="notice-content">{{ item.content }}</view>
				<view class="notice-time">{{ item.create_at }}</view>
			</view>
			<view class="empty" v-if="noticeList.length === 0">暂无公告</view>
		</view>
	</view>
</template>

<script setup>
import {ref, onMounted} from 'vue'
import {onShow, onPullDownRefresh} from '@dcloudio/uni-app'
import {PublishNotice, GetNoticeList} from '@/util/tools/notice'
import { GetCommunityUid } from '@/util/tools/member'

const community_Number = ref('')
const newTitle = ref('')
const newContent = ref('')
const noticeList = ref([])

const loadList = async ()=>{
	if(!community_Number.value) return
	noticeList.value = await GetNoticeList(community_Number.value)
}

const handlePublish = async ()=>{
	if(!newTitle.value.trim() || !newContent.value.trim()){
		uni.showToast({ title: '请填写标题和内容', icon: 'none' })
		return
	}
	try{
		await PublishNotice(community_Number.value, newTitle.value.trim(), newContent.value.trim())
		uni.showToast({ title: '发布成功', icon: 'success' })
		newTitle.value = ''
		newContent.value = ''
		loadList()
	}catch(e){
		uni.showToast({ title: '发布失败', icon: 'none' })
	}
}

onShow(async ()=>{
	// 复用 member.js 里的 GetCommunityUid
	community_Number.value = await GetCommunityUid()
	loadList()
})

// 下拉刷新：刷新历史公告列表
onPullDownRefresh(async () => {
	community_Number.value = await GetCommunityUid()
	await loadList()
	uni.stopPullDownRefresh()
})
</script>

<style scoped>
.notice-container {
	padding: 24rpx;
}
.section-title {
	position: relative;
	padding-left: 20rpx;
	font-size: 32rpx;
	font-weight: 600;
	color: #222;
	margin-bottom: 20rpx;
}
.section-title::before {
	content: '';
	position: absolute;
	left: 0;
	top: 50%;
	transform: translateY(-50%);
	width: 6rpx;
	height: 28rpx;
	border-radius: 6rpx;
	background: linear-gradient(180deg, #4a7dff, #7aa8ff);
}
.issue-notice {
	background: #fff;
	border-radius: 24rpx;
	padding: 28rpx;
	margin-bottom: 30rpx;
	box-shadow: 0 8rpx 24rpx rgba(29, 33, 41, 0.05);
}
.input-title {
	border: 2rpx solid #e8eaed;
	border-radius: 16rpx;
	padding: 20rpx;
	font-size: 28rpx;
	background: #fafbfc;
	margin-bottom: 20rpx;
}
.input-content {
	border: 2rpx solid #e8eaed;
	border-radius: 16rpx;
	padding: 20rpx;
	font-size: 28rpx;
	background: #fafbfc;
	height: 200rpx;
	margin-bottom: 20rpx;
}
.publish-btn {
	background: linear-gradient(135deg, #4a7dff, #6ea3ff);
	color: #fff;
	border-radius: 999rpx;
	font-size: 28rpx;
	box-shadow: 0 8rpx 20rpx rgba(74, 125, 255, 0.3);
}
.notice-history {
	background: #fff;
	border-radius: 24rpx;
	padding: 28rpx;
	box-shadow: 0 8rpx 24rpx rgba(29, 33, 41, 0.05);
}
.notice-item {
	padding: 20rpx 0;
	border-bottom: 1rpx solid #f0f1f3;
}
.notice-title {
	font-size: 30rpx;
	font-weight: 600;
	color: #222;
	margin-bottom: 10rpx;
}
.notice-content {
	font-size: 28rpx;
	color: #555;
	margin-bottom: 10rpx;
	line-height: 1.6;
}
.notice-time {
	font-size: 22rpx;
	color: #999;
}
.empty {
	text-align: center;
	color: #999;
	font-size: 26rpx;
	padding: 40rpx 0;
}
</style>
