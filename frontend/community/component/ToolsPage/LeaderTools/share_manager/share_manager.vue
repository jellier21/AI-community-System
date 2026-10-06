<template>
	<view class="container">
		<view class="section-title">本社区动态管理</view>
		<view class="item" v-for="item in list" :key="item.id">
			<view class="content">{{ item.content }}</view>
			<view class="img-grid" v-if="item.images">
				<image v-for="(img,i) in item.images.split(',')"
				       :key="i" :src="img" mode="aspectFill"
				       @click="previewImage(item.images.split(','), i)" />
			</view>
			<view class="meta">
				<text>发布者：{{ item.name }}</text>
				<text class="time">{{ item.create_at }}</text>
			</view>
			<view class="actions">
				<text class="action-btn del" @click="handleDelete(item.id)">删除该动态</text>
			</view>
		</view>
		<view class="empty" v-if="list.length === 0">暂无动态</view>
	</view>
</template>

<script setup>
import {ref} from 'vue'
import {onShow} from '@dcloudio/uni-app'
import { GetCommunityShares, DeleteShare } from '@/util/tools/share'
import { GetCommunityUid } from '@/util/tools/member'

const list = ref([])

const loadList = async ()=>{
	const uid = await GetCommunityUid()
	if(uid) list.value = await GetCommunityShares(uid)
}

const handleDelete = async (id)=>{
	uni.showModal({
		title: '确认删除',
		content: '删除后不可恢复',
		success: async (r)=>{
			if(!r.confirm) return
			try{
				await DeleteShare(id)
				uni.showToast({ title: '已删除' })
				loadList()
			}catch(e){
				uni.showToast({ title: '删除失败', icon: 'none' })
			}
		}
	})
}

const previewImage = (urls, current)=>{
	uni.previewImage({ urls, current })
}


onShow(loadList)
</script>

<style scoped>
.container { padding: 24rpx; }
.section-title {
	font-size: 32rpx; font-weight: 600; color: #222; margin-bottom: 20rpx;
}
.item {
	background: #fff; border-radius: 16rpx; padding: 24rpx; margin-bottom: 16rpx;
}
.content { font-size: 28rpx; color: #333; margin-bottom: 10rpx; }
.meta {
	display: flex; justify-content: space-between;
	font-size: 22rpx; color: #999; margin-bottom: 10rpx;
}
.actions { margin-top: 10rpx; }
.action-btn.del { font-size: 26rpx; color: #f53f3f; }
.img-grid {
	display: flex; flex-wrap: wrap; gap: 12rpx; margin-bottom: 10rpx;
}
.img-grid image { width: 200rpx; height: 200rpx; border-radius: 12rpx; }

.empty { text-align: center; color: #999; padding: 60rpx 0; }
</style>
