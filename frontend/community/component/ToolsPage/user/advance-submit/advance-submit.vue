<template>
	<view class="container">
		<view class="title">建议上报</view>
		<input class="input" v-model="title" placeholder="请输入标题" />
		<textarea class="textarea" v-model="content" placeholder="请输入建议内容" />
		<button class="submit-btn" @click="handleSubmit">提交</button>
	</view>
</template>

<script setup>
import {ref} from 'vue'
import {SubmitAdvance} from '@/util/tools/advance'
import { GetCommunityUid } from '@/util/tools/member'
import { onPullDownRefresh } from '@dcloudio/uni-app'

const title = ref('')
const content = ref('')

// 下拉刷新：表单页无数据可刷，直接结束刷新动画
onPullDownRefresh(() => {
	uni.stopPullDownRefresh()
})

const handleSubmit = async ()=>{
	if(!title.value.trim() || !content.value.trim()){
		uni.showToast({ title: '请填写标题和内容', icon: 'none' })
		return
	}
	// 拿当前用户所在社区号
	const communityUid = await GetCommunityUid()
	if(!communityUid){
		uni.showToast({ title: '你还未加入社区', icon: 'none' })
		return
	}
	try{
		await SubmitAdvance(communityUid, title.value.trim(), content.value.trim())
		uni.showToast({ title: '提交成功', icon: 'success' })
		title.value = ''
		content.value = ''
		setTimeout(() => uni.navigateBack(), 1500)
	}catch(e){
		uni.showToast({ title: '提交失败', icon: 'none' })
	}
}
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
	margin-bottom: 30rpx;
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
.input {
	background: #fff;
	border: 2rpx solid #e8eaed;
	border-radius: 16rpx;
	padding: 24rpx;
	font-size: 28rpx;
	margin-bottom: 20rpx;
	box-shadow: 0 4rpx 16rpx rgba(29, 33, 41, 0.04);
}
.textarea {
	background: #fff;
	border: 2rpx solid #e8eaed;
	border-radius: 16rpx;
	padding: 24rpx;
	font-size: 28rpx;
	height: 300rpx;
	margin-bottom: 30rpx;
	box-shadow: 0 4rpx 16rpx rgba(29, 33, 41, 0.04);
}
.submit-btn {
	background: linear-gradient(135deg, #4a7dff, #6ea3ff);
	color: #fff;
	border-radius: 999rpx;
	font-size: 30rpx;
	box-shadow: 0 8rpx 20rpx rgba(74, 125, 255, 0.3);
}
</style>
