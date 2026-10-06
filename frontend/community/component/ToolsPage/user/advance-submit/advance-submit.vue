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

const title = ref('')
const content = ref('')

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
	font-size: 36rpx;
	font-weight: 600;
	color: #222;
	margin-bottom: 30rpx;
}
.input {
	background: #fff;
	border-radius: 12rpx;
	padding: 24rpx;
	font-size: 28rpx;
	margin-bottom: 20rpx;
}
.textarea {
	background: #fff;
	border-radius: 12rpx;
	padding: 24rpx;
	font-size: 28rpx;
	height: 300rpx;
	margin-bottom: 30rpx;
}
.submit-btn {
	background: #4a7dff;
	color: #fff;
	border-radius: 999rpx;
	font-size: 30rpx;
}
</style>
