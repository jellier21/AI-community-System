<template>
	<view>
		<view class="title">设置</view>
		<view class="set-container">
			<view class="account-setting" @click="GoToAccountPage">
				<view class="title">
					账号设置
				</view>
				<view class="user-container">
					<view class="user-header">
						<image :src="avatar || '/static/logo.png'" mode="aspectFill"></image>
					</view>
					<view class="user-info">
						<text class="nick">{{ nickname }}</text>
						<text class="community-name">{{ Is_community ? community_name + '成员' : '未加入社区' }}</text>
					</view>
				</view>
			</view>
			<view class="community-container">
				<view class="title">
					我的社区
				</view>
				<view class="Yes_community" v-if="Is_community">
					<view class="currectCommunity">
						<text>当前社区：{{community_name}}</text>
					</view>
					<view class="join-info" v-if="community_type === 'join'">
						<text>创建者：{{ leader_nick || community_leader }}</text>
					</view>
					<view class="quit-community" @click="QuitCommunity">
						<text>{{ community_type === 'leader' ? '解散社区' : '退出社区' }}</text>
						<image src="/static/right.png" mode="aspectFit"></image>
					</view>
				</view>
				<view class="No_community" v-else-if="!loading">
					当前未加入社区
					<view class="Join-community" @click="showJoinDialog = true">
						<text>加入社区</text>
						<image src="/static/right.png" mode="aspectFit"></image>
					</view>

					<view class="create-community" @click="GoToCreateCommunity">
						<text>创建社区</text>
						<image src="/static/right.png" mode="aspectFit"></image>
					</view>
				</view>
				<view class="No_community" v-else>
					加载中...
				</view>
			</view>
			<view class="quit-account" @click="QuitAccount">
				退出登录
			</view>
		</view>
		<view class="Quit-Community-Content" v-if="Quit_Prove" @click="CloseQuitProve">
			<view class="UserQuit" v-if="community_type === 'join'" @click.stop="noop">
				<text>您是否决定退出社区？</text>
				<view class="btns">
					<view class="quit-btn user-c" @click="CloseQuitProve">取消</view>
					<view class="quit-btn user-q" @click="ConfirmQuitCommunity">退出</view>
				</view>
			</view>
			<view class="LeaderQuit" v-else @click.stop="noop">
				<text>您是社区负责人，请输入“确认”二字确认解散社区</text>
				<input type="text" class="leader-input" v-model="prove_leader" placeholder="请输入“确认”二字"/>
				<view class="btns">
					<view class="quit-btn leader-c" @click="CloseQuitProve">取消</view>
					<view class="quit-btn leader-q" @click="ConfirmQuitCommunity">解散</view>
				</view>
			</view>
		</view>
		<!-- 输入社区号加入社区 -->
		<view class="join-mask" v-if="showJoinDialog" @click="showJoinDialog = false">
			<view class="join-dialog" @click.stop="noop">
				<text class="join-title">输入社区号</text>
				<input class="join-input" v-model="inputCommunityUid"
				       placeholder="请输入负责人给你的社区号" />
				<view class="join-btns">
					<text class="join-btn cancel" @click="showJoinDialog = false">取消</text>
					<text class="join-btn ok" @click="confirmJoinCommunity">加入</text>
				</view>
			</view>
		</view>

	</view>
</template>

<script setup>
import { ref } from 'vue';
import { onShow } from '@dcloudio/uni-app';
import { BASEURL_NAME } from '@/util/configuration.js';
import {request} from '@/util/request.js'

const GoToAccountPage = ()=>{
	uni.navigateTo({
		url:'/component/settingPage/myPage/myPage'
	})
}

const GoToCreateCommunity = async ()=>{
	uni.navigateTo({
		url: '/component/settingPage/create_community/create_community'
	})
}
//返回昵称
const nickname = ref(uni.getStorageSync('user_nickname') || '加载中')
const avatar = ref('')
const UploadUserAvatar = async ()=>{
	const account = uni.getStorageSync('user_account') || ''
	if (!account) return
	try{
		const res = await request({
			url:'/api/user/query/avatar',
			data:{ account }
		})
		if (res.code === 200 && res.avatar) {
			avatar.value = res.avatar
		}
	} catch(err) {}
}

const UploadUserNickname = async ()=>{
	const account = uni.getStorageSync('user_account') || ''
	if (!account) return
	try{
		const res = await request({
			url:'/api/user/query/nick',
			data:{ account }
		})
		if (res.code === 200 && res.nickname) {
			nickname.value = res.nickname
			uni.setStorageSync('user_nickname', res.nickname)
		}
	} catch(err) {
		console.log("err:",err)
	}
}

//返回管理员昵称
const leader_nick= ref('')
const UploadLeaderNickname = async ()=>{
	const account = community_leader.value || ''
	if (!account) return
	try{
		const res = await request({
			url:'/api/user/query/nick',
			data:{ account }
		})
		if (res.code === 200 && res.nickname) {
			leader_nick.value = res.nickname
		}
	} catch(err) {
		console.log("err:",err)
	}
}
//退出登录
const QuitAccount = ()=>{
	uni.showModal({
		title:'提示',
		content:'是否确定退出登录？',
		success: (res) => {
			if(res.confirm){
				uni.removeStorageSync("jwtToken")
				uni.removeStorageSync("user_account")
				uni.reLaunch({
					url:"/pages/index/index"
				})
			} else if (res.cancel) {
				console.log('=========取消退出登录======')
			}
		}
	})

}

// ============ 输入社区号加入社区 ============
const showJoinDialog = ref(false)
const inputCommunityUid = ref('')

const confirmJoinCommunity = async ()=>{
	const uid = inputCommunityUid.value.trim()
	if(!uid){
		uni.showToast({ title: '请输入社区号', icon: 'none' })
		return
	}
	try{
		await request({
			url: '/api/community/query/member/join',
			method: 'POST',
			data: { community_uid: uid }
		})
		uni.showToast({ title: '加入成功', icon: 'success' })
		showJoinDialog.value = false
		inputCommunityUid.value = ''
		CheckCommunityStatus()
	}catch(e){
		uni.showToast({ title: '加入失败，可能未被邀请', icon: 'none' })
	}
}


//================检查加入||创建社区情况====================
const Is_community = ref(false)
const loading = ref(true)
const community_name = ref('')
const community_type = ref('')    // leader=自己创建 | join=加入他人
const community_leader = ref('')  // 加入他人社区时，显示创建者账号

const CheckCommunityStatus = ()=>{
	const account = uni.getStorageSync('user_account') || ''
	const token = uni.getStorageSync('jwtToken') || ''

	// 未登录或无令牌时，不请求接口，直接视为未加入
	if (!account || !token) {
		Is_community.value = false
		community_name.value = ''
		community_type.value = ''
		community_leader.value = ''
		loading.value = false
		return
	}

	uni.request({
		url: `${BASEURL_NAME}/api/community/query/have`,
		method: 'GET',
		data: { account },
		header: { Authorization: token },
		success: (res) => {
			if (res.statusCode === 200) {
				const info = res.data?.info
				if (info) {
					Is_community.value = true
					community_name.value = info.community_name || ''
					community_type.value = info.type || ''
					community_leader.value = info.leader || ''
					// 拿到创建者账号后再查其昵称
					UploadLeaderNickname()
				} else {
					Is_community.value = false
					leader_nick.value = ''
				}
			} else if (res.statusCode === 401) {
				// 令牌失效，清除登录态并回登录页
				uni.removeStorageSync('jwtToken')
				uni.removeStorageSync('user_account')
				uni.reLaunch({ url: '/pages/index/index' })
			} else {
				// 未加入任何社区（后端返回 500 + error "未查询到用户加入社区"）
				Is_community.value = false
				community_name.value = ''
				community_type.value = ''
				community_leader.value = ''
				leader_nick.value = ''
			}
		},
		fail: () => {
			uni.showToast({ title: '网络异常', icon: 'none' })
			Is_community.value = false
		},
		complete: () => {
			loading.value = false
		}
	})
}

const Quit_Prove = ref(false)
const prove_leader = ref('')

// 阻止冒泡空函数（弹窗内容区点击不关闭遮罩）
const noop = () => {}

// 关闭退出/解散确认弹层
const CloseQuitProve = ()=>{
	Quit_Prove.value = false
	prove_leader.value = ''
}

// 打开退出/解散确认弹层
const QuitCommunity = ()=>{
	Quit_Prove.value = true
	prove_leader.value = ''
}

// 确认退出/解散：POST /api/community/delete（JWT 鉴权）
const ConfirmQuitCommunity = ()=>{
	const isLeader = community_type.value === 'leader'
	// 管理员解散需输入"确认"校验
	if (isLeader && prove_leader.value.trim() !== '确认') {
		uni.showToast({ title: '请输入"确认"二字', icon: 'none' })
		return
	}
	const account = uni.getStorageSync('user_account') || ''
	request({
		url: '/api/community/delete',
		method: 'POST',
		data: { account }
	}).then((res) => {
		uni.showToast({ title: isLeader ? '社区已解散' : '已退出社区', icon: 'success' })
		// 清空社区状态并关闭弹层
		Is_community.value = false
		community_name.value = ''
		community_type.value = ''
		community_leader.value = ''
		leader_nick.value = ''
		CloseQuitProve()
	}).catch((err) => {
		console.log("err:",err)
		uni.showToast({ title: '操作失败', icon: 'none' })
	})
}

// 每次进入设置页都刷新社区状态和昵称（创建/加入/修改昵称后返回可即时更新）
onShow(() => {
	loading.value = true
	CheckCommunityStatus()
	UploadUserNickname()
	UploadUserAvatar()
})


</script>

<style lang="scss" scoped>
page {
	background: #f4f5f7;
	min-height: 100vh;
}

.setting-page {
	padding-top: var(--status-bar-height);
	padding-bottom: calc(40rpx + env(safe-area-inset-bottom));

	/* ============ 页面标题 ============ */
	> .title {
		height: 100rpx;
		line-height: 100rpx;
		text-align: center;
		font-size: 36rpx;
		font-weight: 600;
		color: #1d2129;
		background: #ffffff;
		border-bottom: 1rpx solid #f0f1f3;
	}
}

.set-container {
	padding: 24rpx;
	display: flex;
	flex-direction: column;
	gap: 24rpx;
}

/* ============ 通用卡片 ============ */
.account-setting,
.community-container {
	background: #ffffff;
	border-radius: 24rpx;
	padding: 28rpx;
	box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);

	> .title {
		font-size: 26rpx;
		color: #86909c;
		margin-bottom: 20rpx;
	}
}

/* ============ 账号设置 ============ */
.account-setting {
	&:active {
		background: #fafbfc;
	}

	.user-container {
		display: flex;
		flex-direction: row;
		align-items: center;

		.user-header {
			width: 110rpx;
			height: 110rpx;
			border-radius: 50%;
			overflow: hidden;
			background: #f2f3f5;
			border: 2rpx solid #e5e6eb;
			flex-shrink: 0;

			image {
				width: 100%;
				height: 100%;
			}
		}

		.user-info {
			margin-left: 24rpx;
			flex: 1;
			min-width: 0;
			display: flex;
			flex-direction: column;

			.nick {
				font-size: 34rpx;
				font-weight: 600;
				color: #1d2129;
			}

			.community-name {
				margin-top: 10rpx;
				font-size: 24rpx;
				color: #86909c;
			}
		}
	}
}

/* ============ 我的社区 ============ */
.community-container {
	/* 通用行样式：左文案 + 右箭头 */
	.quit-community,
	.Join-community,
	.create-community {
		display: flex;
		flex-direction: row;
		align-items: center;
		justify-content: space-between;
		height: 96rpx;
		padding: 0 4rpx;
		border-top: 1rpx solid #f5f6f8;

		text {
			font-size: 30rpx;
			color: #1d2129;
		}

		image {
			width: 28rpx;
			height: 28rpx;
			opacity: 0.4;
		}

		&:active {
			background: #fafbfc;
		}
	}

	.currectCommunity,
	.join-info {
		padding: 16rpx 4rpx;
		font-size: 28rpx;
		color: #4e5969;
	}

	.currectCommunity {
		font-weight: 500;
		color: #1d2129;
	}

	.join-info {
		padding-top: 0;
		font-size: 26rpx;
		color: #86909c;
	}

	.quit-community {
		margin-top: 12rpx;

		text {
			color: #f53f3f;
		}

		image {
			opacity: 0.5;
		}
	}

	.No_community {
		font-size: 28rpx;
		color: #86909c;
		padding: 8rpx 4rpx 0;

		.Join-community {
			border-top: none;
			padding-top: 0;
		}
	}
}

/* ============ 退出登录 ============ */
.quit-account {
	height: 96rpx;
	line-height: 96rpx;
	text-align: center;
	font-size: 30rpx;
	color: #f53f3f;
	background: #ffffff;
	border-radius: 24rpx;
	box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);

	&:active {
		background: #fff1f0;
	}
}

/* ============ 退出/解散确认弹层 ============ */
.Quit-Community-Content {
	position: fixed;
	left: 0;
	top: 0;
	right: 0;
	bottom: 0;
	z-index: 900;
	background: rgba(0, 0, 0, 0.5);
	display: flex;
	align-items: center;
	justify-content: center;

	.UserQuit,
	.LeaderQuit {
		width: 560rpx;
		background: #ffffff;
		border-radius: 24rpx;
		padding: 40rpx 32rpx 28rpx;
		display: flex;
		flex-direction: column;

		> text {
			font-size: 30rpx;
			color: #1d2129;
			line-height: 1.6;
			text-align: center;
		}
	}

	.leader-input {
		margin-top: 28rpx;
		height: 84rpx;
		padding: 0 20rpx;
		background: #f7f8fa;
		border: 1rpx solid #e5e6eb;
		border-radius: 12rpx;
		font-size: 28rpx;
		color: #1d2129;
	}

	.btns {
		margin-top: 32rpx;
		display: flex;
		flex-direction: row;
		border-top: 1rpx solid #f0f1f3;

		.quit-btn {
			flex: 1;
			height: 88rpx;
			line-height: 88rpx;
			text-align: center;
			font-size: 30rpx;

			&:active {
				background: #f7f8fa;
			}
		}

		.user-c,
		.leader-c {
			color: #86909c;
			border-right: 1rpx solid #f0f1f3;
		}

		.user-q,
		.leader-q {
			color: #f53f3f;
			font-weight: 600;
		}
	}
}

/* ============ 加入社区弹层 ============ */
.join-mask {
	position: fixed;
	left: 0;
	top: 0;
	right: 0;
	bottom: 0;
	z-index: 950;
	background: rgba(0, 0, 0, 0.5);
	display: flex;
	align-items: center;
	justify-content: center;

	.join-dialog {
		width: 580rpx;
		background: #ffffff;
		border-radius: 24rpx;
		padding: 36rpx 32rpx 0;
		display: flex;
		flex-direction: column;

		.join-title {
			font-size: 32rpx;
			font-weight: 600;
			color: #1d2129;
			text-align: center;
		}

		.join-input {
			margin-top: 28rpx;
			height: 88rpx;
			padding: 0 20rpx;
			background: #f7f8fa;
			border: 1rpx solid #e5e6eb;
			border-radius: 12rpx;
			font-size: 28rpx;
			color: #1d2129;
		}

		.join-btns {
			margin-top: 32rpx;
			display: flex;
			flex-direction: row;
			border-top: 1rpx solid #f0f1f3;

			.join-btn {
				flex: 1;
				height: 88rpx;
				line-height: 88rpx;
				text-align: center;
				font-size: 30rpx;

				&:active {
					background: #f7f8fa;
				}
			}

			.cancel {
				color: #86909c;
				border-right: 1rpx solid #f0f1f3;
			}

			.ok {
				color: #4a7dff;
				font-weight: 600;
			}
		}
	}
}
</style>
