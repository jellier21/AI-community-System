<template>
	<view class="page">
		<view class="list">
			<view class="friend-item" v-for="f in list" :key="f.room_uid">
				<image class="avatar" :src="f.friend_avatar || '/static/logo.png'" mode="aspectFill"
				       @click="goChat(f)"></image>
				<view class="info" @click="goChat(f)">
					<text class="nick">{{ f.friend_nick || f.friend_account }}</text>
					<text class="account">{{ f.friend_account }}</text>
				</view>
				<text class="del-btn" @click="confirmDelete(f)">删除</text>
			</view>
			<view class="empty" v-if="list.length === 0">
				<text>还没有好友，点右上角+发起私聊吧</text>
			</view>
		</view>
	</view>
</template>

<script setup>
import { ref } from 'vue'
import { onShow, onPullDownRefresh } from '@dcloudio/uni-app'
import { getFriendList, deleteFriend } from '@/util/chat/group'

const list = ref([])

const loadList = async () => {
	list.value = await getFriendList()
}

onShow(loadList)

// 下拉刷新：刷新好友列表
onPullDownRefresh(async () => {
	await loadList()
	uni.stopPullDownRefresh()
})

const goChat = (f) => {
	uni.navigateTo({
		url: '/component/chatPage/chat/chat?uid=' + f.room_uid
	})
}

const confirmDelete = (f) => {
	uni.showModal({
		title: '删除好友',
		content: `确定删除 ${f.friend_nick || f.friend_account}？聊天记录也会一起删除`,
		success: async (r) => {
			if (!r.confirm) return
			try {
				await deleteFriend(f.room_uid)
				uni.showToast({ title: '已删除', icon: 'success' })
				loadList()
			} catch(e) {
				uni.showToast({ title: '删除失败', icon: 'none' })
			}
		}
	})
}
</script>

<style lang="scss" scoped>
page { background: #f5f6f8; }
.page { min-height: 100vh; }
.header {
	background: #fff; padding: 20rpx 32rpx;
	text-align: center; font-size: 34rpx; font-weight: 600;
}
.list { padding: 20rpx 24rpx; }
.friend-item {
	display: flex; align-items: center;
	background: #fff; border-radius: 20rpx;
	padding: 20rpx 24rpx; margin-bottom: 16rpx;
	box-shadow: 0 8rpx 24rpx rgba(29,33,41,0.05);
}
.avatar {
	width: 80rpx; height: 80rpx; border-radius: 50%;
	margin-right: 20rpx; background: #f0f0f0;
	border: 2rpx solid #e8eefc;
}
.info { flex: 1; }
.nick { display: block; font-size: 30rpx; font-weight: 600; color: #222; }
.account { display: block; font-size: 24rpx; color: #999; margin-top: 6rpx; }
.del-btn {
	font-size: 26rpx; color: #f53f3f;
	padding: 10rpx 24rpx;
	border-radius: 999rpx;
	background: #fff1f0;
}
.empty {
	text-align: center; color: #999; padding: 120rpx 0; font-size: 26rpx;
}
</style>
