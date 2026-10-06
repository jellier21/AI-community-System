<template>
	<view class="page">
		<view class="header">
			<text class="title">建立群聊</text>
		</view>
		<view class="group-name">
			<input class="name-input" v-model="groupName" placeholder="输入群聊名称" />
		</view>
		<view class="section-title">选择好友（{{ selected.length }}人已选）</view>
		<view class="list">
			<view class="friend-item" v-for="f in list" :key="f.room_uid" @click="toggle(f)">
				<view class="checkbox" :class="{checked: selected.includes(f.friend_account)}">
					<text v-if="selected.includes(f.friend_account)">✓</text>
				</view>
				<image class="avatar" :src="f.friend_avatar || '/static/logo.png'" mode="aspectFill"></image>
				<view class="info">
					<text class="nick">{{ f.friend_nick || f.friend_account }}</text>
					<text class="account">{{ f.friend_account }}</text>
				</view>
			</view>
			<view class="empty" v-if="list.length === 0">
				<text>还没有好友，先去发起私聊吧</text>
			</view>
		</view>
		<view class="bottom-bar">
			<button class="create-btn" @click="doCreate" :disabled="!groupName || selected.length === 0">
				创建群聊
			</button>
		</view>
	</view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { getFriendList, createGroupChat } from '@/util/chat/group'

const list = ref([])
const selected = ref([])
const groupName = ref('')

onLoad(async () => {
	list.value = await getFriendList()
})

const toggle = (f) => {
	const idx = selected.value.indexOf(f.friend_account)
	if (idx >= 0) {
		selected.value.splice(idx, 1)
	} else {
		selected.value.push(f.friend_account)
	}
}

const doCreate = async () => {
	if (!groupName.value.trim()) {
		uni.showToast({ title: '请输入群名', icon: 'none' })
		return
	}
	if (selected.value.length === 0) {
		uni.showToast({ title: '请至少选一个好友', icon: 'none' })
		return
	}
	try {
		const r = await createGroupChat(groupName.value.trim(), selected.value)
		if (r.code === 200) {
			uni.showToast({ title: '群聊已创建', icon: 'success' })
			setTimeout(() => {
				uni.redirectTo({
					url: '/component/chatPage/chat/chat?uid=' + r.room_uid
				})
			}, 600)
		}
	} catch(e) {
		uni.showToast({ title: '创建失败', icon: 'none' })
	}
}
</script>

<style lang="scss" scoped>
page { background: #f5f6f8; }
.page { min-height: 100vh; padding-bottom: 140rpx; }
.header {
	background: #fff; padding: 20rpx 32rpx;
	text-align: center; font-size: 34rpx; font-weight: 600;
}
.group-name { padding: 24rpx; }
.name-input {
	background: #fff; border-radius: 12rpx;
	padding: 20rpx 24rpx; font-size: 28rpx;
}
.section-title {
	padding: 10rpx 24rpx; font-size: 24rpx; color: #999;
}
.list { padding: 0 24rpx; }
.friend-item {
	display: flex; align-items: center;
	background: #fff; border-radius: 16rpx;
	padding: 20rpx 24rpx; margin-bottom: 16rpx;
}
.checkbox {
	width: 40rpx; height: 40rpx; border: 2rpx solid #ccc;
	border-radius: 50%; margin-right: 20rpx;
	display: flex; align-items: center; justify-content: center;
	font-size: 26rpx; color: #fff;
}
.checkbox.checked {
	background: #4a7dff; border-color: #4a7dff;
}
.avatar {
	width: 72rpx; height: 72rpx; border-radius: 50%;
	margin-right: 20rpx; background: #f0f0f0;
}
.info { flex: 1; }
.nick { display: block; font-size: 28rpx; font-weight: 600; color: #222; }
.account { display: block; font-size: 22rpx; color: #999; margin-top: 4rpx; }
.empty {
	text-align: center; color: #999; padding: 100rpx 0; font-size: 26rpx;
}
.bottom-bar {
	position: fixed; bottom: 0; left: 0; right: 0;
	padding: 20rpx 32rpx; background: #fff;
	border-top: 1rpx solid #eee;
}
.create-btn {
	background: #4a7dff; color: #fff;
	border-radius: 999rpx; font-size: 30rpx;
}
.create-btn[disabled] {
	background: #c8d4ff; color: #fff;
}
</style>
