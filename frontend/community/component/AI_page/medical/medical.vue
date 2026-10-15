<template>
	<!-- 点击遮罩关闭图片选择框 -->
	<view class="select-mask" v-if="ImageOp" @click="ImageOp = false"></view>

	<view class="medical-content">
		<scroll-view class="medical-chat-area" scroll-y :scroll-into-view="scrollIntoView" scroll-with-animation>
			<view v-for="(msg, index) in chatList" :key="index" :id="'msg-' + index" class="chat-msg" :class="msg.role">
				<image v-if="msg.image" :src="msg.image" mode="aspectFill" class="msg-image"></image>
				<view class="msg-bubble" v-if="msg.role === 'ai'">
					<rich-text :nodes="markdownToHtml(msg.content)"></rich-text>
				</view>
				<view class="msg-bubble" v-else>{{ msg.content }}</view>

			</view>
			<view v-if="loading" :id="'msg-' + chatList.length" class="chat-msg ai">
				<view class="msg-bubble loading-bubble">
					<text>正在诊断中</text>
					<view class="loading-dots">
						<view class="dot"></view>
						<view class="dot"></view>
						<view class="dot"></view>
					</view>
				</view>
			</view>
		</scroll-view>


		<view class="user-medical-input-area">
			<view class="chat-image-show">
				<view v-if="selectedFile" class="preview-wrap">
					<image :src="selectedFile" mode="aspectFill" class="preview-image"></image>
					<view class="preview-remove" @click="removeImage">×</view>
				</view>
			</view>


			<view class="medical-input-area">
				<view class="chat-input">
					<input type="text" placeholder="请描述您的症状…" placeholder-class="input-placeholder" v-model="formQuestion"/>
				</view>
				<view class="image-operator-btn" @click="ImageOPBtnswitch">
					<image src="/static/chat/push.png" mode="aspectFit"></image>
				</view>
				<view class="submit-medical-data">
					<button @click="sendQuestion">发送</button>
				</view>
			</view>

			<view class="image-select-op" v-if="ImageOp">
				<view class="select-file" @click="chooseImage">
					<image src="/static/chat/album.png" mode="aspectFit"></image>
				</view>
				<view class="take-medical-photo" @click="takePhoto">
					<image src="/static/chat/take_photo.png" mode="aspectFit"></image>
				</view>
			</view>
		</view>
	</view>
</template>

<script setup>
import {ref} from 'vue'
import {SubmitMyQuestion} from '@/util/tools/medical.js'
const ImageOp = ref(false)
const loading = ref(false)   // 是否正在诊断

const ImageOPBtnswitch = ()=>{
	ImageOp.value = !ImageOp.value
}

const formQuestion = ref('')
const selectedFile = ref('')   // 选中的图片本地路径，没选就是空字符串

// 删除暂存的图片
const removeImage = ()=>{
	selectedFile.value = ''
}


// 从相册选图
const chooseImage = ()=>{
	uni.chooseImage({
		count: 1,
		success: (res) => {
			selectedFile.value = res.tempFilePaths[0]
			ImageOp.value = false          // ← 选完关闭选择框
			console.log('选图成功:', selectedFile.value)
		},
		fail: () => {
			selectedFile.value = ''
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
			ImageOp.value = false          // ← 选完关闭选择框
		},
		fail: () => { selectedFile.value = '' }
	})
}


const chatList = ref([])   // 聊天消息列表

const scrollIntoView = ref('')   // 滚动目标消息 id
const scrollToBottom = () => {
	setTimeout(() => {
		const len = chatList.value.length
		if (loading.value) {
			scrollIntoView.value = 'msg-' + len   // 滚到 loading 位置
		} else if (len > 0) {
			scrollIntoView.value = 'msg-' + (len - 1)
		}
	}, 50)
}

// 极简 markdown → HTML（覆盖 AI 回答常用格式）
const markdownToHtml = (text = '') => {
	let html = text
		.replace(/&/g, '&amp;')
		.replace(/</g, '&lt;')
		.replace(/>/g, '&gt;')
	html = html
		// 标题 ### / ## / #
		.replace(/^### (.*)$/gm, '<h4>$1</h4>')
		.replace(/^## (.*)$/gm, '<h3>$1</h3>')
		.replace(/^# (.*)$/gm, '<h2>$1</h2>')
		// 加粗 **文字**
		.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
		// 行内代码 `code`
		.replace(/`([^`]+)`/g, '<code>$1</code>')
		// 无序列表 - 或 * 开头
		.replace(/^[-*] (.*)$/gm, '• $1')
		// 数字列表 1. 2. 3.
		.replace(/^(\d+)\. (.*)$/gm, '<br/>$1. $2')
		// 换行
		.replace(/\n/g, '<br/>')
	return html
}



const sendQuestion = async ()=>{
	if (formQuestion.value.trim() === '' && !selectedFile.value) {
		uni.showToast({ title: '请先输入内容或选择图片', icon: 'none' })
		return
	}
	const question = formQuestion.value
	const file = selectedFile.value
		chatList.value.push({ role: 'user', content: question || '[图片]', image: file })
		formQuestion.value = ''
		selectedFile.value = ''
		loading.value = true              // 开始诊断 → 显示动画
		scrollToBottom()                  // ← 滚到 loading
		try {
			const res = await SubmitMyQuestion(question, file)
			if (res && res.response) {
				chatList.value.push({
					role: 'ai',
					content: `【症状】${res.response.symptom}\n【诊断】${res.response.diagnosis}`
				})
			} else {
				chatList.value.push({ role: 'ai', content: '未获取到有效回复' })
			}
			scrollToBottom()              // ← AI 回复后滚到底
		} catch (err) {
			console.log('发送失败:', err)
			chatList.value.push({ role: 'ai', content: '发送失败，请重试' })
			scrollToBottom()
		} finally {
			loading.value = false
			scrollToBottom()              // ← 关掉 loading 后再滚一次，去掉多余空白
		}

}



</script>

<style scoped>
	page{
		height: 100%;
	}

.medical-content{
	width: 100%;
	height: 100vh;
	display: flex;
	flex-direction: column;   /* ← 纵向排列，chat 区占满高度 */
}
.select-mask{
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	z-index: 40;
	background: transparent;
}

.medical-chat-area{
	width: 100%;
	flex: 1;
	height: 0;              /* ← scroll-view 配合 flex 必须加 height:0 才能滚动 */
	padding: 20rpx 24rpx calc(140rpx + env(safe-area-inset-bottom));
	box-sizing: border-box;
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

.chat-msg.user .msg-bubble{ background: #3b7cff; color: #ffffff; }
.chat-msg.ai .msg-bubble{ background: #ffffff; color: #333333; }
.msg-image{
	width: 160rpx;
	height: 160rpx;
	border-radius: 12rpx;
	margin-bottom: 8rpx;
}

.loading-bubble{
	display: flex;
	align-items: center;
	gap: 12rpx;
}
.loading-dots{
	display: flex;
	gap: 8rpx;
}
.loading-dots .dot{
	width: 12rpx;
	height: 12rpx;
	border-radius: 50%;
	background: #999999;
	animation: dot-blink 1.2s infinite ease-in-out;
}
.loading-dots .dot:nth-child(2){ animation-delay: 0.2s; }
.loading-dots .dot:nth-child(3){ animation-delay: 0.4s; }
@keyframes dot-blink{
	0%, 80%, 100%{ opacity: 0.3; transform: scale(0.8); }
	40%{ opacity: 1; transform: scale(1); }
}


.user-medical-input-area{
	width: 100%;
	min-height: 120rpx;
	overflow: hidden;
	position: fixed;
	bottom: 0;
	z-index: 50;             /* ← 高于遮罩，输入区不被遮罩挡住 */
	background-color: aliceblue;
}


.medical-input-area{
	width: 100%;
	height: 120rpx;
	display: flex;
	justify-content: center;
	align-items: center;
	padding: 0 10rpx;      /* 左右留点边距 */
	box-sizing: border-box;
}

.chat-image-show{
	padding: 10rpx 20rpx 0;
}
.preview-wrap{
	position: relative;
	width: 140rpx;
	height: 140rpx;
	border-radius: 12rpx;
	overflow: hidden;
}
.preview-image{
	width: 100%;
	height: 100%;
}
.preview-remove{
	position: absolute;
	top: 0;
	right: 0;
	width: 40rpx;
	height: 40rpx;
	line-height: 40rpx;
	text-align: center;
	background: rgba(0,0,0,0.5);
	color: #ffffff;
	font-size: 28rpx;
	border-radius: 0 0 0 12rpx;
}


.chat-input{
	flex: 1;
	height: 75rpx;
	min-width: 0;          /* 关键：允许 flex 子项被压缩，防止 input 撑破 */
}
.chat-input input{
	width: 100%;
	height: 100%;
	background-color: white;
	border-radius: 25rpx;
	padding: 0 20rpx;
	box-sizing: border-box;
}

.image-operator-btn{
	width: 75rpx;
	height: 75rpx;
	flex-shrink: 0;        /* 不允许被压缩 */
	margin-left: 10rpx;
}
.image-operator-btn image{
	width: 100%;
	height: 100%;
}

.submit-medical-data{
	width: 175rpx;
	height: 75rpx;
	border-radius: 30rpx;
	overflow: hidden;
	flex-shrink: 0;
	margin-left: 10rpx;
}
.submit-medical-data button{
	width: 100%;
	height: 100%;
	background-color: greenyellow;
}

.image-select-op{
	display: flex;
	justify-content: space-evenly;
}
.select-file, .take-medical-photo{
	width: 30%;
}
.select-file image, .take-medical-photo image{
	width: 100%;
}
</style>
