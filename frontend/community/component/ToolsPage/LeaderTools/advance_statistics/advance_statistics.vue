<template>
	<view class="container">
		<view class="title">建议列表</view>
		<view class="item" v-for="item in list" :key="item.id">
			<view class="item-title">{{ item.title }}</view>
			<view class="item-content">{{ item.content }}</view>
			<view class="item-meta">来自：{{ item.account }}</view>
		</view>
		<view class="empty" v-if="list.length === 0">暂无建议</view>
	</view>
</template>

<script setup>
import {ref, onMounted} from 'vue'
import {GetAdvanceList} from '@/util/tools/advance'
import {GetCommunityUid} from '@/util/tools/member'
import { onPullDownRefresh } from '@dcloudio/uni-app'

const list = ref([])

onMounted(async ()=>{
	const communityUid = await GetCommunityUid()
	if(communityUid){
		list.value = await GetAdvanceList(communityUid)
	}
})

// 下拉刷新：刷新建议列表
onPullDownRefresh(async () => {
	const communityUid = await GetCommunityUid()
	if(communityUid){
		list.value = await GetAdvanceList(communityUid)
	}
	uni.stopPullDownRefresh()
})
</script>

<style scoped>
.container {
	padding: 24rpx;
}
.title {
	position: relative;
	padding-left: 20rpx;
	font-size: 36rpx;
	font-weight: 600;
	color: #222;
	margin-bottom: 24rpx;
}
.title::before {
	content: '';
	position: absolute;
	left: 0;
	top: 50%;
	transform: translateY(-50%);
	width: 6rpx;
	height: 30rpx;
	border-radius: 6rpx;
	background: linear-gradient(180deg, #4a7dff, #7aa8ff);
}
.item {
	background: #fff;
	border-radius: 20rpx;
	padding: 24rpx;
	margin-bottom: 16rpx;
	box-shadow: 0 8rpx 24rpx rgba(29, 33, 41, 0.05);
}
.item-title {
	font-size: 30rpx;
	font-weight: 600;
	color: #222;
	margin-bottom: 10rpx;
}
.item-content {
	font-size: 28rpx;
	color: #555;
	margin-bottom: 10rpx;
}
.item-meta {
	font-size: 22rpx;
	color: #999;
}
.empty {
	text-align: center;
	color: #999;
	padding: 60rpx 0;
	font-size: 26rpx;
}
</style>
