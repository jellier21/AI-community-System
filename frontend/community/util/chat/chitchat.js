import { ref } from "vue";
import { request } from '@/util/request.js'
import { WS_BASEURL_NAME } from '@/util/configuration.js'

// ============ 状态 ============
// 消息列表：{ type: 'msg'|'system', account, nickname, content, time }
const _chitchat_ = ref([])
const myAvatar = ref('')
const setMyAvatar = (url) => { myAvatar.value = url }


let socketTask = null        // uni.connectSocket 返回的 SocketTask
let isConnected = false      // WebSocket 是否已建立
let isManualClose = false    // 是否主动关闭（主动关闭时禁止自动重连）
let roomUid = ''             // 当前房间 uid
let myAccount = ''           // 当前登录账号
let reconnectCount = 0       // 已重连次数
let reconnectTimer = null    // 重连定时器句柄

const MAX_RECONNECT = 3      // 最大重连次数
const RECONNECT_BASE_DELAY = 2000 // 重连基础间隔(ms)，逐次翻倍

// ============ 连接管理 ============

// 建立 WebSocket 连接（携带 JWT + 房间参数）
const connectChatSocket = (uid, account) => {
	roomUid = uid
	myAccount = account

	// 已有连接先关闭，避免重复建连
	if (socketTask) {
		closeChatSocket()
	}
	// 在旧连接 close() 异步触发 onClose 之前重置标志：
	// 旧连接的 onClose 读到 isManualClose=true 不会重连，新连接断线才允许重连
	isManualClose = false

	const token = uni.getStorageSync('jwtToken')
	// token 双通道：header（App/小程序）+ query（H5 原生 WebSocket 无法自定义 header，走后端 query 兜底）
	const url = `${WS_BASEURL_NAME}/ws?account=${encodeURIComponent(account)}&roomUid=${encodeURIComponent(uid)}${token ? '&token=' + encodeURIComponent(token) : ''}`

	socketTask = uni.connectSocket({
		url,
		header: token ? { Authorization: token } : {},
		complete: () => {}
	})

	// 连接建立成功
	socketTask.onOpen(() => {
		isConnected = true
		reconnectCount = 0
		console.log('聊天连接已建立:', url)
	})

	// 收到服务端消息
	socketTask.onMessage((res) => {
		handleSocketMessage(res.data)
	})

	// 连接出错
	socketTask.onError((err) => {
		console.error('聊天连接错误:', err)
	})

	// 连接关闭 → 主动关闭不重连；异常断开才重连
	socketTask.onClose(() => {
		isConnected = false
		if (isManualClose) {
			console.log('主动关闭连接，不自动重连')
			return
		}
		tryReconnect()
	})
}

// ============ 消息收发 ============

// 处理服务端推送的消息（JSON）
const handleSocketMessage = (raw) => {
	try {
		const msg = JSON.parse(raw)
		if (msg.type === 'msg' || msg.type === 'message') {
			// 普通聊天消息
			_chitchat_.value.push({
				type: 'msg',
				account: msg.account,
				nickname: msg.nickname || msg.account,
				avatar: msg.avatar || '',
				content: msg.content,
				time: msg.time
			})

		} else {
			// join / leave / system 统一作为居中系统提示
			// _chitchat_.value.push({
			// 	type: 'system',
			// 	content: msg.content || '',
			// 	time: msg.time
			// })
			console.log(msg)
		}
	} catch (e) {
		// 非 JSON 的纯文本消息，兜底展示
		_chitchat_.value.push({
			type: 'system',
			content: raw,
			time: ''
		})
	}
}

// 发送消息：先本地回显，再推送给服务端
// 说明：后端广播会跳过发送者，必须本地追加一条，否则自己看不到消息
const sendChatMessage = (content) => {
	const text = (content || '').trim()
	if (!text) {
		return Promise.resolve(false)
	}
	if (!socketTask || !isConnected) {
		uni.showToast({ title: '连接未建立', icon: 'none' })
		return Promise.reject(new Error('socket 未连接'))
	}

	// 本地回显
	_chitchat_.value.push({
		type: 'msg',
		account: myAccount,
		nickname: myAccount,
		avatar: myAvatar.value,
		content: text,
		time: formatNow()
	})


	return new Promise((resolve, reject) => {
		socketTask.send({
			data: JSON.stringify({ type: 'msg', content: text }),
			success: () => resolve(true),
			fail: (err) => {
				console.error('消息发送失败:', err)
				reject(err)
			}
		})
	})
}

// ============ 历史消息 ============

// 进入房间时加载最近消息（默认 20 条）
const UploadHistoryMessage = async (uid) => {
	try {
		const res = await request({
			url: '/api/ws/messages/recent',
			data: { room_uid: uid, limit: 20 }
		})
		if (res.code === 200 && Array.isArray(res.data)) {
			_chitchat_.value = res.data.map((m) => ({
				type: 'msg',
				account: m.account,
				nickname: m.nickname || m.account,
				avatar: m.avatar || '',
				content: m.content,
				time: m.created_at
			}))

			return true
		}
		return false
	} catch (e) {
		console.error('加载最近消息失败:', e)
		return false
	}
}

// ============ 断开与重连 ============

// 主动关闭连接（页面卸载时调用）
const closeChatSocket = () => {
	isManualClose = true
	if (reconnectTimer) {
		clearTimeout(reconnectTimer)
		reconnectTimer = null
	}
	if (socketTask) {
		try {
			socketTask.close({ code: 1000, reason: 'leave' })
		} catch (e) {
			console.error('关闭连接异常:', e)
		}
		socketTask = null
	}
	isConnected = false
}

// 断线重连（有限次数，间隔递增）
const tryReconnect = () => {
	if (reconnectCount >= MAX_RECONNECT) {
		console.log('重连次数已达上限，停止重连')
		return
	}
	reconnectCount++
	const delay = RECONNECT_BASE_DELAY * reconnectCount
	console.log(`连接断开，${delay / 1000}s 后尝试第 ${reconnectCount} 次重连...`)
	reconnectTimer = setTimeout(() => {
		connectChatSocket(roomUid, myAccount)
	}, delay)
}

// ============ 工具 ============

// 当前时间格式化：YYYY-MM-DD HH:mm:ss
const formatNow = () => {
	const d = new Date()
	const p = (n) => (n < 10 ? '0' + n : '' + n)
	return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`
}

export {
	_chitchat_,
	connectChatSocket,
	sendChatMessage,
	UploadHistoryMessage,
	closeChatSocket,
	setMyAvatar
}

