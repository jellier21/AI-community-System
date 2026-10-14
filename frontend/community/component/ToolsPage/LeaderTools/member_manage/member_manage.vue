<template>
	<view>
		<view class="title">
			<text>社区号:{{community_Number}}</text>
			<view class="copy-btn" @click="uni.setClipboardData({
				data:community_Number
			})">
				<image src="/static/tools/copy.png" mode="aspectFit"></image>
				<text>点击复制社区号分享给成员加入社区</text>
			</view>
		</view>
		<view class="member-list">
			<view class="list-header">
				<text class="member-content-title">成员列表管理：</text>
				<view class="member-manage-tool">
					<text class="manage" @click="openAddDialog">添加成员</text>
					<text class="manage" style="margin-left: 16rpx;" @click="clickGroupManageBtn()">批量管理</text>
					<view class="manage-item">
						<view class="delete" v-if="group_manager" @click="doBatchDelete">
							<text>删除({{selectedAccounts.length}})</text>
							<image src="/static/tools/delete.png" mode="aspectFit"></image>
						</view>
					</view>
				</view>
			</view>
			<view class="member-item" v-for="(item,index) in _memberList_" :key="item.account">
				<view class="one-item">
					<view class="member-name">
						{{item.name}}
					</view>
					<view class="member-account">
						{{item.account}}
					</view>
					<view class="member-Is">
						<text v-if="item.join_is">已加入</text>
						<text v-else>未加入</text>
					</view>
					<view class="groupDel-Content" v-if="group_manager" @click="toggleSelect(item.account)">
						<image v-if="!selectedAccounts.includes(item.account)"
						       src="/static/tools/no_select.png" mode="aspectFit"></image>
						<image v-else
						       src="/static/tools/select.png" mode="aspectFit"></image>
					</view>
				</view>
			</view>
			<!-- 添加成员弹窗 -->
			<view class="dialog-mask" v-if="showAddDialog" @click="showAddDialog = false">
				<view class="dialog-content" @click.stop>
					<text class="dialog-title">添加成员</text>
					<input class="dialog-input" v-model="newAccount"
						   placeholder="输入对方账号" />
					<input class="dialog-input" v-model="newMemberName"
						   placeholder="输入成员名称（聊天里显示的名字）" />
					<view class="dialog-btns">
						<text class="dialog-btn cancel" @click="showAddDialog = false">取消</text>
						<text class="dialog-btn confirm" @click="confirmAddMember">确定</text>
					</view>
				</view>
			</view>
		</view>
	</view>
</template>

<script setup>
import {ref} from 'vue'
import {onShow} from '@dcloudio/uni-app'
import { _memberList_,GetCommunityUid,GetMemberList,RemoveMembers,AddMember } from '@/util/tools/member'


const community_Number = ref('sss')

//-------批量删除-------
const group_manager = ref(false)
const selectedAccounts = ref([])

const clickGroupManageBtn = ()=>{
	group_manager.value = !group_manager.value
	selectedAccounts.value = []          // 切换模式时清空勾选
}

// 勾选/取消勾选
const toggleSelect = (account)=>{
	const myAccount = uni.getStorageSync("user_account")
	if(account === myAccount) return     // 自己不能被选
	const i = selectedAccounts.value.indexOf(account)
	if(i >= 0){
		selectedAccounts.value.splice(i, 1)
	}else{
		selectedAccounts.value.push(account)
	}
}

// 批量删除
const doBatchDelete = async ()=>{
	if(selectedAccounts.value.length === 0){
		uni.showToast({ title: '请先选择要删除的成员', icon: 'none' })
		return
	}
	uni.showModal({
		title: '确认删除',
		content: `确定移除选中的 ${selectedAccounts.value.length} 名成员？`,
		success: async (r)=>{
			if(!r.confirm) return
			try{
				await RemoveMembers(community_Number.value, selectedAccounts.value)
				uni.showToast({ title: '删除成功', icon: 'success' })
				group_manager.value = false
				selectedAccounts.value = []
				await GetMemberList(community_Number.value)   // 刷新列表
			}catch(e){
				uni.showToast({ title: '删除失败', icon: 'none' })
			}
		}
	})
}

//------------------------添加成员-----------------------------
const showAddDialog = ref(false)
const newAccount = ref('')
const newMemberName = ref('')

const openAddDialog = ()=>{
	newAccount.value = ''
	newMemberName.value = ''
	showAddDialog.value = true
}

const confirmAddMember = async ()=>{
	const account = newAccount.value.trim()
	const name = newMemberName.value.trim()
	if(!account){
		uni.showToast({ title: '请输入对方账号', icon: 'none' })
		return
	}
	if(!name){
		uni.showToast({ title: '请输入成员名称', icon: 'none' })
		return
	}
	try{
		await AddMember(community_Number.value, account, name)
		uni.showToast({ title: '已邀请，等待对方加入', icon: 'success' })
		showAddDialog.value = false
		await GetMemberList(community_Number.value)
	}catch(e){
		uni.showToast({ title: '添加失败，对方可能还为注册小程序账号', icon: 'none' })
	}
}



onShow(async()=>{
	community_Number.value = await GetCommunityUid()
	await GetMemberList(community_Number.value)
})
</script>

<style lang="scss" scoped>
page {
	background: #f4f5f7;
}

/* ============ 顶部社区号卡片 ============ */
.title {
	margin: 24rpx;
	padding: 36rpx 28rpx;
	background: linear-gradient(135deg, #4a7dff 0%, #6ea3ff 100%);
	border-radius: 24rpx;
	box-shadow: 0 10rpx 30rpx rgba(74, 125, 255, 0.25);
	display: flex;
	flex-direction: column;
	align-items: center;
	color: #ffffff;

	> text:first-child {
		font-size: 42rpx;
		font-weight: 600;
		letter-spacing: 2rpx;
	}

	.copy-btn {
		margin-top: 28rpx;
		display: flex;
		flex-direction: row;
		align-items: center;
		padding: 14rpx 28rpx;
		border-radius: 999rpx;
		background: rgba(255, 255, 255, 0.2);
		transition: all 0.15s ease;

		image {
			width: 28rpx;
			height: 28rpx;
			margin-right: 12rpx;
		}

		text {
			font-size: 24rpx;
			color: #ffffff;
		}

		&:active {
			background: rgba(255, 255, 255, 0.35);
			transform: scale(0.96);
		}
	}
}

/* ============ 成员列表容器 ============ */
.member-list {
	margin: 0 24rpx 40rpx;
	background: #ffffff;
	border-radius: 24rpx;
	overflow: hidden;
	box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);

	/* ---------- 头部 ---------- */
	.list-header {
		display: flex;
		flex-direction: row;
		align-items: center;
		justify-content: space-between;
		padding: 28rpx 28rpx 24rpx;
		border-bottom: 1rpx solid #f0f1f3;

		.member-content-title {
			font-size: 30rpx;
			font-weight: 600;
			color: #222222;
		}
	}

	.member-manage-tool {
		display: flex;
		flex-direction: row;
		align-items: center;

		.manage {
			font-size: 26rpx;
			color: #4a7dff;
			padding: 8rpx 18rpx;
			border-radius: 999rpx;
			background: rgba(74, 125, 255, 0.08);

			&:active {
				background: rgba(74, 125, 255, 0.18);
			}
		}

		.manage-item {
			margin-left: 16rpx;

			.delete {
				display: flex;
				flex-direction: row;
				align-items: center;
				padding: 8rpx 18rpx;
				border-radius: 999rpx;
				background: rgba(245, 63, 63, 0.08);

				text {
					font-size: 26rpx;
					color: #f53f3f;
					margin-right: 8rpx;
				}

				image {
					width: 26rpx;
					height: 26rpx;
				}

				&:active {
					background: rgba(245, 63, 63, 0.16);
					opacity: 0.85;
				}
			}
		}
	}

	/* ---------- 每一行 ---------- */
	.member-item {
		.one-item {
			display: flex;
			flex-direction: row;
			align-items: center;
			padding: 26rpx 28rpx;
			border-bottom: 1rpx solid #f5f6f8;

			&:active {
				background: #fafbfc;
			}
		}

		&:last-child .one-item {
			border-bottom: none;
		}

		.member-name {
			flex: 1;
			min-width: 0;
			font-size: 30rpx;
			color: #222222;
			overflow: hidden;
			white-space: nowrap;
			text-overflow: ellipsis;
		}

		.member-account {
			width: 220rpx;
			margin-left: 16rpx;
			font-size: 24rpx;
			color: #999999;
			overflow: hidden;
			white-space: nowrap;
			text-overflow: ellipsis;
		}

		.member-Is {
			width: 130rpx;
			margin-left: 12rpx;
			text-align: center;

			text {
				display: inline-block;
				padding: 6rpx 16rpx;
				border-radius: 8rpx;
				font-size: 22rpx;
			}
		}

		.groupDel-Content {
			width: 56rpx;
			height: 56rpx;
			margin-left: 8rpx;
			display: flex;
			align-items: center;
			justify-content: center;

			image {
				width: 40rpx;
				height: 40rpx;
			}

			&:active {
				opacity: 0.6;
			}
		}
	}
}
/* ============ 添加成员弹窗 ============ */
.dialog-mask {
	position: fixed;
	top: 0; left: 0; right: 0; bottom: 0;
	background: rgba(0,0,0,0.4);
	display: flex; align-items: center; justify-content: center;
	z-index: 999;
}
.dialog-content {
	width: 600rpx;
	background: #fff;
	border-radius: 20rpx;
	padding: 40rpx;
}
.dialog-title {
	font-size: 32rpx;
	font-weight: 600;
	text-align: center;
	display: block;
	margin-bottom: 30rpx;
}
.dialog-input {
	border: 1rpx solid #ddd;
	border-radius: 12rpx;
	padding: 20rpx;
	font-size: 28rpx;
	margin-bottom: 24rpx;
}
.dialog-btns {
	display: flex;
	justify-content: space-around;
	margin-top: 10rpx;
}
.dialog-btn {
	font-size: 30rpx;
	padding: 16rpx 50rpx;
}
.dialog-btn.confirm { color: #4a7dff; }
.dialog-btn.cancel { color: #999; }

</style>
