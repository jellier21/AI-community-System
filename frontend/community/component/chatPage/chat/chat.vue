<template>
	<view class="container">
		<!-- 消息列表区 -->
		<scroll-view class="chitchat-content" scroll-y :scroll-into-view="scrollIntoView" scroll-with-animation>
			<view v-for="(item,index) in _chit_" :key="index" :id="'msg-' + index">
				<!-- 系统提示（join / leave / system） -->
				<view class="system-item" v-if="item.type === 'system'">
					{{ item.content }}
				</view>
				<!-- 聊天气泡 -->
				<view class="chitchat-item" v-else>
					<view class="msg-row me-row" v-if="item.account === myAccount">
						<view class="bubble-wrap me-bubble">
							<view class="me-item">{{ item.content }}</view>
						</view>
						<image class="avatar" :src="item.avatar || '/static/logo.png'" mode="aspectFill"
						       @click="goPerson(item.account)"></image>
					</view>
					<view class="msg-row" v-else>
						<image class="avatar" :src="item.avatar || '/static/logo.png'" mode="aspectFill"
						       @click="goPerson(item.account)"></image>
						<view class="bubble-wrap">
							<view class="another-nick" v-if="item.nickname">{{ item.nickname }}</view>
							<view class="another-item">{{ item.content }}</view>
						</view>
					</view>
				</view>

			</view>
		</scroll-view>

		<!-- 底部输入区 -->
		<view class="input-message-area">
			<view class="input-text">
				<input
					type="text"
					v-model="message_text"
					confirm-type="send"
					@confirm="sendMessage"
					placeholder="输入消息..."
				/>
			</view>
			<view class="send-message">
				<view class="send-text-btn">
					<button @click="sendMessage">发送</button>
				</view>
			</view>
		</view>
	</view>
</template>

<script>
import { request } from '@/util/request.js';
	import {
		setMyAvatar,
		_chitchat_,
		connectChatSocket,
		sendChatMessage,
		UploadHistoryMessage,
		closeChatSocket
	} from '@/util/chat/chitchat';
export default{
	data(){
		return{
			message_text: '',
			_chit_: _chitchat_,
			myAccount: '',
			roomUid: '',
			leader: '',
			photo_file_text: false,
			scrollIntoView: ''
		}
	},
	watch: {
		// 消息列表变化（新消息/历史加载）后滚动到底部
		_chit_: {
			handler() {
				this.scrollToBottom()
			},
			deep: true
		}
	},
	methods: {
		openPhotoFile() {
			this.photo_file_text = !this.photo_file_text
		},
		// 发送消息
		sendMessage() {
			const content = this.message_text
			if (!content || !content.trim()) {
				return
			}
			sendChatMessage(content)
				.then(() => {
					this.message_text = ''
					uni.hideKeyboard()
				})
				.catch(() => {})
		},
		// 滚动到最新一条消息
		scrollToBottom() {
			this.$nextTick(() => {
				const len = _chitchat_.value.length
				if (len > 0) {
					this.scrollIntoView = 'msg-' + (len - 1)
				}
			})
		},
		goPerson(account) {
			uni.navigateTo({
				url: '/component/chatPage/person/person?account=' + account
			})
		},

	},
	onLoad(option){
		console.log("room uid:", option.uid)
		console.log("room leader:", option.leader)
		this.myAccount = uni.getStorageSync('user_account')
		this.roomUid = option.uid
		this.leader = option.leader || ''
		// 加载自己的头像
		request({ url: '/api/user/query/avatar', data: { account: this.myAccount } }).then((res) => {
			if (res.code === 200 && res.avatar) {
				setMyAvatar(res.avatar)
			}
		})
	},
	onReady(){
		// 1. 先加载最近历史消息
		UploadHistoryMessage(this.roomUid).then((ok) => {
			console.log('历史消息加载' + (ok ? '成功' : '失败'))
		})
		// 2. 建立 WebSocket 实时连接
		connectChatSocket(this.roomUid, this.myAccount)
	},
	onUnload(){
		// 离开页面时关闭连接
		closeChatSocket()
	}
}
</script>

<style lang="scss" scoped>
	/* ========== 全局：让页面可被 scroll-view 撑开 ========== */
	.container {
		display: flex;
		flex-direction: column;
		height: 100vh;
		box-sizing: border-box;
		background: #f2f3f5;
	}

	/* ========== 消息列表区 ========== */
	.chitchat-content {
		flex: 1;
		height: 0;              /* 关键：配合 flex:1，小程序/APP 下才能正常滚动 */
		box-sizing: border-box;
		padding: 20rpx 24rpx 8rpx;
	}

	/* 系统提示（join / leave / system） */
	.system-item {
		align-self: center;
		margin: 12rpx auto 20rpx;
		padding: 6rpx 20rpx;
		max-width: 72%;
		text-align: center;
		font-size: 24rpx;
		line-height: 1.4;
		color: #ffffff;
		background: rgba(0, 0, 0, 0.18);
		border-radius: 8rpx;
		word-break: break-all;
	}

	/* 每条消息的容器 */
	.chitchat-item {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		margin-bottom: 24rpx;
	}

	/* 自己的气泡（右侧） */
	.me-item {
		align-self: flex-end;
		max-width: 67%;
		position: relative;
		padding: 18rpx 24rpx;
		font-size: 30rpx;
		line-height: 1.55;
		color: #1a1a1a;
		background: #95ec69;
		border-radius: 12rpx;
		word-break: break-word;
		white-space: pre-wrap;
		box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);

		/* 小尖角 */
		&::after {
			content: '';
			position: absolute;
			top: 26rpx;
			right: -14rpx;
			border-width: 10rpx 0 10rpx 16rpx;
			border-style: solid;
			border-color: transparent transparent transparent #95ec69;
		}
	}

	/* 对方的气泡（左侧） */
	.another-item {
		align-self: flex-start;
		max-width: 67%;
		position: relative;
		padding: 18rpx 24rpx;
		font-size: 30rpx;
		line-height: 1.55;
		color: #1a1a1a;
		background: #ffffff;
		border-radius: 12rpx;
		word-break: break-word;
		white-space: pre-wrap;
		box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);

		&::after {
			content: '';
			position: absolute;
			top: 26rpx;
			left: -14rpx;
			border-width: 10rpx 16rpx 10rpx 0;
			border-style: solid;
			border-color: transparent #ffffff transparent transparent;
		}
	}

	/* 对方昵称 */
	.another-nick {
		margin-bottom: 8rpx;
		font-size: 24rpx;
		line-height: 1.2;
		color: #9a9a9a;
	}
	.msg-row {
		display: flex;
		align-items: flex-start;
		gap: 16rpx;
		width: 100%;
		margin-bottom: 24rpx;
		box-sizing: border-box;
	}
	.msg-row.me-row {
		justify-content: flex-end;
	}
	.msg-row .avatar {
		width: 72rpx;
		height: 72rpx;
		border-radius: 50%;
		flex-shrink: 0;
		background: #f0f0f0;
	}
	.bubble-wrap {
		max-width: 68%;
	}


	/* ========== 底部输入区 ========== */
	.input-message-area {
		position: relative;
		display: flex;
		align-items: center;
		padding: 16rpx 20rpx;
		padding-bottom: calc(16rpx + env(safe-area-inset-bottom));
		background: #f7f7f7;
		border-top: 1rpx solid #e4e4e4;
	}

	.input-text {
		flex: 1;
		min-width: 0;

		input {
			width: 100%;
			height: 72rpx;
			padding: 0 28rpx;
			box-sizing: border-box;
			font-size: 28rpx;
			color: #1a1a1a;
			background: #ffffff;
			border-radius: 36rpx;
		}
	}

	.send-message {
		flex-shrink: 0;
		margin-left: 16rpx;

		/* 去掉 uni-button 默认边框与样式 */
		button {
			padding: 0;
			margin: 0;
			line-height: normal;
			background: transparent;

			&::after {
				border: none;
			}
		}
	}

	/* “+” 号按钮 */
	.send-file {
		button {
			width: 72rpx;
			height: 72rpx;
			line-height: 68rpx;
			text-align: center;
			font-size: 44rpx;
			font-weight: 300;
			color: #555555;
			background: #ffffff;
			border: 1rpx solid #dcdcdc;
			border-radius: 50%;
			box-sizing: border-box;
		}
	}

	/* 发送按钮 */
	.send-text-btn {
		button {
			height: 72rpx;
			padding: 0 32rpx;
			line-height: 72rpx;
			font-size: 28rpx;
			color: #ffffff;
			background: #07c160;
			border-radius: 36rpx;

			&:active {
				background: #06ad56;
			}
		}
	}

	/* ========== 图片/拍照面板（悬浮在输入区上方） ========== */
	.file-type {
		position: absolute;
		left: 0;
		right: 0;
		bottom: 100%;
		display: flex;
		align-items: center;
		padding: 32rpx 40rpx calc(32rpx + env(safe-area-inset-bottom));
		background: #f7f7f7;
		border-top: 1rpx solid #e4e4e4;

		.take-photo,
		.album-file {
			width: 120rpx;
			height: 120rpx;
			margin-right: 40rpx;
			display: flex;
			align-items: center;
			justify-content: center;
			background: #ffffff;
			border-radius: 16rpx;
			box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.06);

			image {
				width: 64rpx;
				height: 64rpx;
			}

			&:active {
				background: #ececec;
			}
		}
	}
</style>
