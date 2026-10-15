<template>
	<view class="mondo-content">

		<!-- 消息展示区 -->
		<scroll-view class="Message-Area" scroll-y :scroll-into-view="scrollIntoView" scroll-with-animation>
			<view v-for="(msg, index) in chatList" :key="index" :id="'msg-' + index" class="chat-msg" :class="msg.role">
				<image v-if="msg.image" :src="msg.image" mode="aspectFill" class="msg-image"></image>
				<view class="msg-bubble" v-if="msg.role === 'ai'">
					<rich-text :nodes="markdownToHtml(msg.content)"></rich-text>
				</view>
				<view class="msg-bubble" v-else>{{ msg.content }}</view>
			</view>
			<view v-if="loading" :id="'msg-' + chatList.length" class="chat-msg ai">
				<view class="msg-bubble loading-bubble">
					<text>正在思考中</text>
					<view class="loading-dots">
						<view class="dot"></view>
						<view class="dot"></view>
						<view class="dot"></view>
					</view>
				</view>
			</view>
		</scroll-view>

		<!-- 底部输入面板 -->
		<view class="input-panel">
			<!-- 附件暂存区 -->
			<view class="FileTempArea">
				<view v-if="selectedFile" class="file-item">
					<image :src="selectedFile" mode="aspectFill"></image>
					<view class="del" @click="removeFile">×</view>
				</view>
			</view>

			<!-- 输入行 -->
			<view class="inputArea">
				<view class="input">
					<input v-model="text" type="text" placeholder="请输入您想问的问题…" placeholder-class="ph"
						@confirm="sendQuestion" />
					<button @click="ImageOp = !ImageOp">+</button>
				</view>
				<button :disabled="!canSend" @click="sendQuestion">发送</button>
			</view>
		</view>

		<!-- 图片选择：遮罩 + 选项 -->
		<view class="select-mask" v-if="ImageOp" @click="ImageOp = false"></view>
		<view class="image-select-op" v-if="ImageOp">
			<view class="op-item" @click="chooseFromAlbum">📁 从相册选择</view>
			<view class="op-item" @click="takePhoto">📷 拍照上传</view>
		</view>
	</view>
</template>
<script setup>
import {ref, computed} from 'vue'
import {SubmitTeacherQuestion} from '@/util/tools/mondo.js'
import { onPullDownRefresh } from '@dcloudio/uni-app'

const chatList = ref([])        // 聊天消息
const text = ref('')            // 输入框内容
const selectedFile = ref('')    // 选中的图片本地路径
const loading = ref(false)      // 是否在等待 AI 回复
const ImageOp = ref(false)      // 图片选择框开关
const scrollIntoView = ref('')  // 滚动目标

const canSend = computed(() => text.value.trim() !== '' || selectedFile.value !== '')

// 从相册选图
const chooseFromAlbum = ()=>{
	uni.chooseImage({
		count: 1,
		success: (res) => {
			selectedFile.value = res.tempFilePaths[0]
			ImageOp.value = false
		}
	})
}

// 拍照
const takePhoto = ()=>{
	uni.chooseImage({
		count: 1,
		sourceType: ['camera'],
		success: (res) => {
			selectedFile.value = res.tempFilePaths[0]
			ImageOp.value = false
		}
	})
}

// 删除暂存图片
const removeFile = ()=>{ selectedFile.value = '' }

// markdown → HTML（AI 回答渲染）
const markdownToHtml = (content = '')=>{
	let html = content
		.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
	html = html
		.replace(/^### (.*)$/gm, '<h4>$1</h4>')
		.replace(/^## (.*)$/gm, '<h3>$1</h3>')
		.replace(/^# (.*)$/gm, '<h2>$1</h2>')
		.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
		.replace(/`([^`]+)`/g, '<code>$1</code>')
		.replace(/^[-*] (.*)$/gm, '• $1')
		.replace(/^(\d+)\. (.*)$/gm, '<br/>$1. $2')
		.replace(/\n/g, '<br/>')
	return html
}

// 滚动到底部
const scrollToBottom = ()=>{
	setTimeout(()=>{
		const len = chatList.value.length
		if (loading.value) scrollIntoView.value = 'msg-' + len
		else if (len > 0) scrollIntoView.value = 'msg-' + (len - 1)
	}, 50)
}

// 发送
const sendQuestion = async ()=>{
	if (!canSend.value) {
		uni.showToast({ title: '请先输入内容或选择图片', icon: 'none' })
		return
	}
	const question = text.value
	const file = selectedFile.value
	chatList.value.push({ role: 'user', content: question || '[图片]', image: file })
	text.value = ''
	selectedFile.value = ''
	loading.value = true
	scrollToBottom()

	try {
		const res = await SubmitTeacherQuestion(question, file)
		// Go 端返回 { response: { code, answer, image_desc } }
		const answer = (res && res.response && res.response.answer)
			? res.response.answer
			: '未获取到有效回复'
		chatList.value.push({ role: 'ai', content: answer })
	} catch (err) {
		console.log('发送失败:', err)
		const msg = (err && err.message) ? err.message : '发送失败，请重试'
		chatList.value.push({ role: 'ai', content: '发送失败：' + msg })
	} finally {
		loading.value = false
		scrollToBottom()
	}
}

// 下拉刷新：聊天页无可刷数据，直接结束刷新动画
onPullDownRefresh(() => {
	uni.stopPullDownRefresh()
})
</script>


<style lang="scss" scoped>
	page {
		height: 100%;
		background: #f2f4f7;
	}

	/* ========== 1. 消息展示区 ========== */
	.Message-Area {
		height: 100vh;
		box-sizing: border-box;
		/* 底部输入栏是 fixed，用 padding 把内容顶上去，避免被遮挡 */
		padding: 24rpx 24rpx calc(280rpx + env(safe-area-inset-bottom));
		overflow-y: auto;
	}

	/* ========== 2. 底部输入面板（固定吸底） ========== */
	.input-panel {
		position: fixed;
		left: 0;
		right: 0;
		bottom: 0;
		z-index: 99;
		background: #ffffff;
		border-top: 1rpx solid #ebedf0;
		border-radius: 24rpx 24rpx 0 0;
		box-shadow: 0 -6rpx 24rpx rgba(0, 0, 0, 0.04);
		padding: 20rpx 24rpx calc(20rpx + constant(safe-area-inset-bottom));
		padding: 20rpx 24rpx calc(20rpx + env(safe-area-inset-bottom));
	}

	/* ========== 3. 文件暂存区（附件缩略图，位于输入框上方） ========== */
	.FileTempArea {
		display: flex;
		flex-wrap: wrap;
		gap: 16rpx;
		max-height: 240rpx;
		overflow-y: auto;
		margin-bottom: 16rpx;

		&:empty {              /* 没有附件时不占位 */
			margin-bottom: 0;
			max-height: 0;
		}

		/* 动态插入的附件项示例：<view class="file-item"><image/></view> */
		.file-item {
			position: relative;
			width: 120rpx;
			height: 120rpx;
			border-radius: 16rpx;
			background: #f2f4f7;
			overflow: hidden;

			image { width: 100%; height: 100%; }

			.del {
				position: absolute;
				top: 0;
				right: 0;
				width: 36rpx;
				height: 36rpx;
				line-height: 32rpx;
				text-align: center;
				font-size: 24rpx;
				color: #fff;
				background: rgba(0, 0, 0, 0.5);
				border-radius: 0 16rpx 0 16rpx;
			}
		}
	}

	/* ========== 4. 输入行：输入框 + 发送 ========== */
	.inputArea {
		display: flex;
		align-items: center;
		gap: 20rpx;
	}

	/* --- 输入框（右侧内嵌圆形 + 按钮） --- */
	.input {
		position: relative;
		flex: 1;
		min-width: 0;
		display: flex;
		align-items: center;
		height: 80rpx;
		padding: 0 96rpx 0 32rpx;   /* 右侧留 96rpx 给 + 按钮 */
		background: #f4f6f9;
		border-radius: 40rpx;

		input {
			flex: 1;
			height: 80rpx;
			font-size: 28rpx;
			color: #222;
			background: transparent;
		}

		/* 圆形 + 按钮：定位到输入框右侧 */
		button {
			position: absolute;
			right: 8rpx;
			top: 50%;
			transform: translateY(-50%);
			width: 72rpx;
			height: 72rpx;
			margin: 0;
			padding: 0;
			display: flex;
			align-items: center;
			justify-content: center;
			font-size: 44rpx;
			line-height: 1;
			font-weight: 300;
			color: #6b7280;
			background: #ffffff;
			border: 2rpx solid #e5e7eb;
			border-radius: 50%;
			box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.06);

			&::after { border: none; }   /* 去掉小程序 button 默认边框 */
			&:active { background: #f0f2f5; }
		}
	}

	/* --- 发送按钮（胶囊） --- */
	.inputArea > button {
		flex-shrink: 0;
		height: 80rpx;
		padding: 0 40rpx;
		margin: 0;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 28rpx;
		line-height: 1;
		color: #ffffff;
		background: linear-gradient(135deg, #4b8cff, #2f6bff);
		border-radius: 40rpx;
		box-shadow: 0 6rpx 16rpx rgba(47, 107, 255, 0.3);

		&::after { border: none; }
		&:active { opacity: 0.85; transform: scale(0.97); }

		/* 空内容时的禁用态：<button :disabled="!text"> */
		&[disabled] {
			color: #b6c0cf;
			background: #e9edf3;
			box-shadow: none;
		}
	}

		.mondo-content{
			height: 100vh;
			display: flex;
			flex-direction: column;
		}
		/* 消息区改 scroll-view 适配：占满中间 */
		.Message-Area{
			flex: 1;
			height: 0;
			padding-bottom: 320rpx;   /* 给固定输入面板让位 */
		}
		.chat-msg{
			margin-bottom: 20rpx;
			display: flex;
			flex-direction: column;
		}
		.chat-msg.user{ align-items: flex-end; }
		.chat-msg.ai{ align-items: flex-start; }
		.msg-bubble{
			max-width: 70%;
			padding: 20rpx 24rpx;
			border-radius: 16rpx;
			font-size: 28rpx;
			line-height: 1.6;
			word-break: break-all;
		}
		.chat-msg.user .msg-bubble{
			background: linear-gradient(135deg, #4a7dff, #6ea3ff);
			color: #fff;
			box-shadow: 0 6rpx 16rpx rgba(74, 125, 255, 0.22);
		}
		.chat-msg.ai .msg-bubble{
			background: #fff;
			color: #333;
			box-shadow: 0 4rpx 16rpx rgba(29, 33, 41, 0.05);
		}
		.msg-image{
			width: 160rpx; height: 160rpx;
			border-radius: 12rpx;
			margin-bottom: 8rpx;
		}

		/* 思考中动画 */
		.loading-bubble{ display: flex; align-items: center; gap: 12rpx; }
		.loading-dots{ display: flex; gap: 8rpx; }
		.loading-dots .dot{
			width: 12rpx; height: 12rpx; border-radius: 50%;
			background: #999; animation: dot-blink 1.2s infinite ease-in-out;
		}
		.loading-dots .dot:nth-child(2){ animation-delay: 0.2s; }
		.loading-dots .dot:nth-child(3){ animation-delay: 0.4s; }
		@keyframes dot-blink{
			0%, 80%, 100%{ opacity: 0.3; transform: scale(0.8); }
			40%{ opacity: 1; transform: scale(1); }
		}

		/* 图片选择：遮罩 + 底部弹出框 */
		.select-mask{
			position: fixed; top: 0; left: 0; right: 0; bottom: 0;
			z-index: 98; background: rgba(0,0,0,0.3);
		}
		.image-select-op{
			position: fixed; left: 0; right: 0; bottom: 0;
			z-index: 99;
			background: #fff;
			border-radius: 24rpx 24rpx 0 0;
			padding: 24rpx 24rpx calc(24rpx + env(safe-area-inset-bottom));
			display: flex; flex-direction: column; gap: 20rpx;
		}
		.op-item{
			height: 88rpx; line-height: 88rpx; text-align: center;
			background: #f4f6f9; border-radius: 16rpx;
			font-size: 30rpx; color: #333;
		}
		.ph{ color: #a8b0bd; }

</style>
