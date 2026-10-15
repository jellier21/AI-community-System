<template>
	<view class="page">
		<!-- 顶部资料区 -->
		<view class="profile">
			<view class="avatar" @click="onAvatarClick">
				<image :src="avatar || '/static/logo.jpeg'" mode="aspectFill"></image>
			</view>

			<view class="nickname">{{ nickname || '未设置昵称' }}</view>
			<view class="account">账号：{{ account }}</view>
			<view class="signature" v-if="signature">{{ signature }}</view>
		</view>

		<!-- 设置列表 -->
		<view class="cell-group">
			<view class="cell" @click="openEdit('nick')">
				<text class="cell-label">昵称</text>
				<view class="cell-value">
					<text class="cell-text">{{ nickname || '未设置' }}</text>
					<image src="/static/right.png" mode="aspectFit"></image>
				</view>
			</view>
			<view class="cell" @click="openEdit('signature')">
				<text class="cell-label">个性签名</text>
				<view class="cell-value">
					<text class="cell-text" :class="{ 'is-empty': !signature }">{{ signature || '未设置' }}</text>
					<image src="/static/right.png" mode="aspectFit"></image>
				</view>
			</view>
			<view class="cell" @click="openEdit('password')">
				<text class="cell-label">修改密码</text>
				<view class="cell-value">
					<image src="/static/right.png" mode="aspectFit"></image>
				</view>
			</view>
		</view>

		<!-- 编辑弹层 -->
		<view class="mask" v-if="editType" @click="closeEdit">
			<view class="dialog" @click.stop>
				<view class="dialog-title">{{ editTitle }}</view>

				<!-- 修改密码：当前密码 + 新密码 + 确认新密码 -->
				<block v-if="editType === 'password'">
					<input
						class="dialog-input"
						type="password"
						v-model="oldPassword"
						placeholder="请输入当前密码"
						:focus="true"
					/>
					<input
						class="dialog-input"
						type="password"
						v-model="newPassword"
						placeholder="请输入新密码"
						style="margin-top: 20rpx;"
					/>
					<input
						class="dialog-input"
						type="password"
						v-model="confirmPassword"
						placeholder="请再次输入新密码"
						style="margin-top: 20rpx;"
						confirm-type="done"
						@confirm="confirmEdit"
					/>
				</block>

				<!-- 昵称 / 签名：单输入框 -->
				<input
					v-else
					class="dialog-input"
					type="text"
					v-model="editValue"
					:placeholder="editPlaceholder"
					:maxlength="editType === 'nick' ? 20 : 50"
					:focus="true"
					confirm-type="done"
					@confirm="confirmEdit"
				/>

				<view class="dialog-btns">
					<button class="btn-cancel" @click="closeEdit">取消</button>
					<button class="btn-confirm" @click="confirmEdit">保存</button>
				</view>
			</view>
		</view>
	</view>
</template>

<script setup>
	import { ref } from 'vue';
	import { onLoad, onPullDownRefresh } from '@dcloudio/uni-app';
	import { request } from '@/util/request.js';
	import { uploadOneImage } from '@/util/tools/share.js';

	// ===== 用户信息 =====
	const account = ref('')
	const nickname = ref('')      // 昵称以后端查询结果为准，不拿账号兜底
	const signature = ref('')
	const avatar = ref('')


	// ===== 编辑弹层状态 =====
	const editType = ref('')      // '' | 'nick' | 'signature' | 'password'
	const editTitle = ref('')
	const editPlaceholder = ref('')
	const editValue = ref('')
	// 修改密码三字段：当前密码 / 新密码 / 确认新密码
	const oldPassword = ref('')
	const newPassword = ref('')
	const confirmPassword = ref('')

	// 查询昵称：GET /api/user/query/nick?account=xxx
	const fetchNickname = async () => {
		try {
			const res = await request({
				url: '/api/user/query/nick',
				data: { account: account.value }
			})
			if (res.code === 200 && res.nickname) {
				nickname.value = res.nickname
				uni.setStorageSync('user_nickname', res.nickname)
			} else {
				nickname.value = ''
			}
		} catch (e) {
			console.error('查询昵称失败:', e)
			nickname.value = ''
		}
	}

	const fetchAvatar = async ()=>{
		try{
			const res = await request({
				url: '/api/user/query/avatar',
				data: { account: account.value }
			})
			if(res.code === 200 && res.avatar){
				avatar.value = res.avatar
			}
		}catch(e){}
	}

	// 以服务端资料为准拉取签名（避免本地 storage 残留上一个账号的签名）
	const fetchProfile = async () => {
		try {
			const res = await request({
				url: '/api/user/query/profile',
				data: { account: account.value }
			})
			if (res.code === 200 && res.profile) {
				signature.value = res.profile.signature || ''
				nickname.value = res.profile.nickname || nickname.value
				if (res.profile.signature) {
					uni.setStorageSync('user_signature', res.profile.signature)
				} else {
					uni.removeStorageSync('user_signature')
				}
			}
		} catch (e) {
			console.error('查询用户资料失败:', e)
		}
	}

	const onAvatarClick = ()=>{
		uni.showModal({
			title: '修改头像',
			content: '是否上传新头像？',
			success: async (r)=>{
				if(!r.confirm) return
				uni.chooseImage({
					count: 1,
					success: async (chooseRes)=>{
						try{
							uni.showLoading({ title: '上传中...' })
							const avatarUrl = await uploadOneImage(chooseRes.tempFilePaths[0])
							await request({
								url: '/api/user/amend/avatar',
								method: 'POST',
								data: { account: account.value, avatar: avatarUrl }
							})
							uni.hideLoading()
							uni.showToast({ title: '头像已更新', icon: 'success' })
							avatar.value = avatarUrl
						}catch(e){
							uni.hideLoading()
							const msg = (e && e.message) ? String(e.message) : '上传失败'
							uni.showToast({ title: msg.length > 18 ? msg.slice(0, 18) : msg, icon: 'none' })
						}
					}
				})
			}
		})
	}


	onLoad(() => {
		account.value = uni.getStorageSync('user_account') || ''
		// 本地缓存仅作首屏占位，最终以服务端 profile 为准
		signature.value = uni.getStorageSync('user_signature') || ''
		fetchNickname()
		fetchAvatar()
		fetchProfile()
	})

	// 下拉刷新：重新拉取昵称、头像与签名
	onPullDownRefresh(async () => {
		await fetchNickname()
		await fetchAvatar()
		await fetchProfile()
		uni.stopPullDownRefresh()
	})


	// ===== 弹层操作 =====
	const openEdit = (type) => {
		editType.value = type
		editValue.value = ''
		oldPassword.value = ''
		newPassword.value = ''
		confirmPassword.value = ''
		if (type === 'nick') {
			editTitle.value = '修改昵称'
			editPlaceholder.value = '请输入新昵称'
			editValue.value = nickname.value || ''
		} else if (type === 'signature') {
			editTitle.value = '编辑个性签名'
			editPlaceholder.value = '请输入个性签名'
			editValue.value = signature.value
		} else if (type === 'password') {
			editTitle.value = '修改密码'
		}
	}

	const closeEdit = () => {
		editType.value = ''
		editValue.value = ''
		oldPassword.value = ''
		newPassword.value = ''
		confirmPassword.value = ''
	}

	const confirmEdit = () => {
		if (editType.value === 'password') {
			submitPassword()
			return
		}
		const value = editValue.value.trim()
		if (!value) {
			uni.showToast({ title: '内容不能为空', icon: 'none' })
			return
		}
		if (editType.value === 'nick') {
			submitNick(value)
		} else if (editType.value === 'signature') {
			submitSignature(value)
		}
	}

	// ===== 调用后端接口 =====

	// POST /api/user/amend/nick
	const submitNick = async (value) => {
		try {
			const res = await request({
				url: '/api/user/amend/nick',
				method: 'POST',
				data: { account: account.value, nickname: value }
			})
			if (res.code === 200) {
				nickname.value = value
				uni.setStorageSync('user_nickname', value)
				uni.showToast({ title: '昵称修改成功', icon: 'success' })
				closeEdit()
			} else {
				uni.showToast({ title: res.error || '修改失败', icon: 'none' })
			}
		} catch (e) {
			console.error('修改昵称失败:', e)
		}
	}

	// POST /api/user/amend/signature
	const submitSignature = async (value) => {
		try {
			const res = await request({
				url: '/api/user/amend/signature',
				method: 'POST',
				data: { account: account.value, signature: value }
			})
			if (res.code === 200) {
				signature.value = value
				uni.setStorageSync('user_signature', value)
				uni.showToast({ title: '签名修改成功', icon: 'success' })
				closeEdit()
			} else {
				uni.showToast({ title: res.error || '修改失败', icon: 'none' })
			}
		} catch (e) {
			console.error('修改签名失败:', e)
		}
	}

	// POST /api/user/amend/password
	// 后端契约：password_old=当前密码（bcrypt 校验身份），password_new=新密码（写入）
	const submitPassword = async () => {
		const oldPwd = oldPassword.value.trim()
		const newPwd = newPassword.value.trim()
		const confirmPwd = confirmPassword.value.trim()

		if (!oldPwd || !newPwd || !confirmPwd) {
			uni.showToast({ title: '请填写完整', icon: 'none' })
			return
		}
		if (newPwd !== confirmPwd) {
			uni.showToast({ title: '两次输入的新密码不一致', icon: 'none' })
			return
		}
		if (newPwd === oldPwd) {
			uni.showToast({ title: '新密码不能与当前密码相同', icon: 'none' })
			return
		}
		try {
			const res = await request({
				url: '/api/user/amend/password',
				method: 'POST',
				data: {
					account: account.value,
					password_old: oldPwd,
					password_new: newPwd
				}
			})
			if (res.code === 200) {
				uni.showToast({ title: '密码修改成功', icon: 'success' })
				closeEdit()
			} else {
				uni.showToast({ title: res.error || '修改失败', icon: 'none' })
			}
		} catch (e) {
			console.error('修改密码失败:', e)
		}
	}
</script>

<style lang="scss" scoped>
	page {
		background-color: #f5f6f7;
	}

	.page {
		min-height: 100vh;
		padding-bottom: 60rpx;
		box-sizing: border-box;
	}

	/* ---------- 顶部资料区 ---------- */
	.profile {
		position: relative;
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 100rpx 40rpx 80rpx;
		background: #fff;
		overflow: hidden;

		/* 顶部渐变装饰 */
		&::before {
			content: '';
			position: absolute;
			left: 0;
			right: 0;
			top: 0;
			height: 220rpx;
			background: linear-gradient(135deg, #4a7dff, #7aa8ff);
		}

		.avatar {
			position: relative;
			width: 140rpx;
			height: 140rpx;
			border-radius: 50%;
			overflow: hidden;
			border: 6rpx solid #ffffff;
			box-shadow: 0 8rpx 20rpx rgba(74, 125, 255, 0.25);
			margin-bottom: 28rpx;

			image {
				width: 100%;
				height: 100%;
			}
		}

		.nickname {
			position: relative;
			font-size: 40rpx;
			font-weight: 600;
			color: #1a1a1a;
			margin-bottom: 10rpx;
		}

		.account {
			position: relative;
			font-size: 24rpx;
			color: #9a9fa8;
			margin-bottom: 16rpx;
		}

		.signature {
			position: relative;
			font-size: 26rpx;
			color: #6b7280;
			max-width: 80%;
			text-align: center;
			word-break: break-all;
		}
	}

	/* ---------- 设置列表 ---------- */
	.cell-group {
		margin: 24rpx;
		background: #fff;
		border-radius: 24rpx;
		overflow: hidden;
		box-shadow: 0 8rpx 24rpx rgba(29, 33, 41, 0.05);

		.cell {
			display: flex;
			justify-content: space-between;
			align-items: center;
			padding: 34rpx 32rpx;
			border-bottom: 1rpx solid #f0f0f0;

			&:last-child {
				border-bottom: none;
			}

			&:active {
				background: #fafafa;
			}

			.cell-label {
				font-size: 30rpx;
				color: #1a1a1a;
			}

			.cell-value {
				display: flex;
				align-items: center;
				min-width: 0;
				max-width: 55%;

				.cell-text {
					font-size: 28rpx;
					color: #8a8f99;
					margin-right: 8rpx;
					overflow: hidden;
					text-overflow: ellipsis;
					white-space: nowrap;

					&.is-empty {
						color: #c0c4cc;
					}
				}

				image {
					width: 24rpx;
					height: 24rpx;
					flex-shrink: 0;
				}
			}
		}
	}

	/* ---------- 编辑弹层 ---------- */
	.mask {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		z-index: 999;
		background: rgba(0, 0, 0, 0.4);
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.dialog {
		width: 600rpx;
		background: #fff;
		border-radius: 24rpx;
		padding: 48rpx 40rpx 40rpx;
		box-sizing: border-box;
		box-shadow: 0 16rpx 48rpx rgba(29, 33, 41, 0.12);

		.dialog-title {
			font-size: 32rpx;
			font-weight: 600;
			color: #1a1a1a;
			text-align: center;
			margin-bottom: 40rpx;
		}

		.dialog-input {
			width: 100%;
			height: 88rpx;
			padding: 0 24rpx;
			background: #f5f6f7;
			border-radius: 12rpx;
			font-size: 28rpx;
			color: #1a1a1a;
			box-sizing: border-box;
		}

		.dialog-btns {
			display: flex;
			margin-top: 40rpx;

			button {
				flex: 1;
				height: 88rpx;
				line-height: 88rpx;
				margin: 0;
				padding: 0;
				font-size: 30rpx;
				border-radius: 12rpx;
				box-sizing: border-box;

				& + button {
					margin-left: 20rpx;
				}

				&::after {
					border: none;
				}
			}

			.btn-cancel {
				background: #f0f2f5;
				color: #6b7280;
			}

			.btn-confirm {
				background: linear-gradient(135deg, #4a7dff, #6ea3ff);
				color: #fff;
			}
		}
	}
</style>
