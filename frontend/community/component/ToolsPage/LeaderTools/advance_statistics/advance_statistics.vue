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

const list = ref([])

onMounted(async ()=>{
	const communityUid = await GetCommunityUid()
	if(communityUid){
		list.value = await GetAdvanceList(communityUid)
	}
})
</script>

<style scoped>
.container {
	padding: 24rpx;
}
.title {
	font-size: 36rpx;
	font-weight: 600;
	color: #222;
	margin-bottom: 24rpx;
}
.item {
	background: #fff;
	border-radius: 16rpx;
	padding: 24rpx;
	margin-bottom: 16rpx;
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
