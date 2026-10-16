<template>
	<view class="home-page">
		<view class="Welcome">
			<view class="greet">
				<view class="time-domain">
					{{ domain }}
				</view>
				<view class="community-greet">
					<text>欢迎使用 邻友小圈</text>
				</view>
			</view>
			<view class="header-image">
				<image src="/static/logo.jpeg" mode="aspectFill"></image>
			</view>
		</view>
		<view class="ai-user-service">
			<navigator url="/component/AI_page/mondo/mondo" class="ai-answer">
				<image src="/static/home/QA_dialog.png" mode="aspectFit"></image>
				<text>ai百科问答</text>
			</navigator>
			<navigator url="/component/AI_page/medical/medical" class="ai-medical-service">
				<image src="/static/home/medical.png" mode="aspectFit"></image>
				<text>ai医疗咨询</text>
			</navigator>
		</view>
		<view class="announcement" v-if="latestNotice.title" @click="goNoticeList">
			<view class="announce-title">
				<image src="/static/home/notice.png" mode="aspectFit"></image>
				{{ latestNotice.title }}
			</view>
			<view class="announce-content">{{ latestNotice.content }}</view>
		</view>
		<view class="announcement" v-else>
			暂无公告
		</view>

		<view class="recommend">
			<view class="recommend-title">为您推荐</view>
			<view class="recommend-container" v-if="latestShares.length">
				<view class="rec-item" v-for="item in latestShares" :key="item.id" @click="goShare">
					<view class="rec-head">
						<image class="rec-avatar" :src="item.avatar || '/static/logo.png'" mode="aspectFill"></image>
						<view class="rec-head-info">
							<text class="rec-author">{{ item.name }}</text>
							<text class="rec-time">{{ item.create_at }}</text>
						</view>
					</view>
					<view class="rec-content">{{ item.content }}</view>
					<image v-if="firstImage(item.images)" class="rec-thumb" :src="firstImage(item.images)" mode="aspectFill"></image>
				</view>
				<view class="rec-more" @click="goShare">
						<text>查看全部动态</text>
						<image class="rec-more-icon" src="/static/right.png" mode="aspectFit"></image>
					</view>
			</view>
			<view class="recommend-container rec-empty" v-else>
				暂无动态，快去社区圈发布第一条吧
			</view>
		</view>
	</view>
</template>

<script setup>
import { onMounted, onUnmounted, ref } from 'vue';
import { onPullDownRefresh } from '@dcloudio/uni-app';
const domain = ref()
const GreetTime = ()=>{
	if(currentHour.value < 6){
		domain.value = '夜猫子好啊~'
	} else if(currentHour.value < 12){
		domain.value = '早上好啊~'
	} else if(currentHour.value < 18){
		domain.value = '下午好啊~'
	} else {
		domain.value = '晚上好啊~'
	}
}

let timer = null
const getCurrentHours = () => {
	return new Date().getHours()
}
const currentHour = ref()

//=====================公告=========================
import { GetLatestNotice } from '@/util/tools/notice'
import { GetCommunityUid } from '@/util/tools/member'
import { GetLatestShares } from '@/util/tools/share'

const latestNotice = ref({})

const loadLatestNotice = async () => {
	const account = uni.getStorageSync('user_account')
	if (!account) return
	// 先拿当前用户所在社区号
	const communityUid = await GetCommunityUid()
	if (!communityUid) return
	latestNotice.value = await GetLatestNotice(communityUid)
}

const goNoticeList = () => {
	// 有需要再跳转到公告列表页，暂时只展示
}

//=====================推荐动态=========================
const latestShares = ref([])

const firstImage = (images) => {
	if (!images) return ''
	return images.split(',')[0]
}

const loadLatestShares = async () => {
	try {
		latestShares.value = await GetLatestShares()
	} catch (e) {
		console.log('加载推荐动态失败:', e)
	}
}

const goShare = () => {
	uni.navigateTo({ url: '/component/ToolsPage/user/share/share' })
}



onMounted(() => {
	currentHour.value = getCurrentHours()
	timer = setInterval(() => {
		currentHour.value = getCurrentHours()
		GreetTime()
	}, 1000)
	loadLatestNotice()
	loadLatestShares()
})


onUnmounted(() => {
	if(timer)
		clearInterval(timer)
})

// 下拉刷新：刷新问候语、最新公告与推荐动态
onPullDownRefresh(async () => {
	GreetTime()
	await loadLatestNotice()
	await loadLatestShares()
	uni.stopPullDownRefresh()
})
</script>

<style lang="scss" scoped>
page {
	background: #f4f5f7;
	min-height: 100vh;
}

.home-page {
	padding: calc(var(--status-bar-height) + 24rpx) 24rpx calc(40rpx + env(safe-area-inset-bottom));
	display: flex;
	flex-direction: column;
	gap: 24rpx;
}

/* ============ 顶部欢迎区 ============ */
.Welcome {
	position: relative;
	overflow: hidden;
	display: flex;
	flex-direction: row;
	align-items: center;
	justify-content: space-between;
	padding: 44rpx 32rpx;
	border-radius: 28rpx;
	background: linear-gradient(135deg, #4a7dff 0%, #7aa8ff 100%);
	box-shadow: 0 12rpx 32rpx rgba(74, 125, 255, 0.28);

	/* 装饰圆 */
	&::before {
		content: '';
		position: absolute;
		right: -80rpx;
		top: -80rpx;
		width: 280rpx;
		height: 280rpx;
		border-radius: 50%;
		background: rgba(255, 255, 255, 0.14);
	}

	&::after {
		content: '';
		position: absolute;
		right: 96rpx;
		bottom: -120rpx;
		width: 220rpx;
		height: 220rpx;
		border-radius: 50%;
		background: rgba(255, 255, 255, 0.08);
	}

	.greet {
		position: relative;
		z-index: 1;
		flex: 1;
		min-width: 0;
		display: flex;
		flex-direction: column;

		.time-domain {
			font-size: 44rpx;
			font-weight: 700;
			color: #ffffff;
			letter-spacing: 2rpx;
			line-height: 1.3;
		}

		.community-greet {
			margin-top: 14rpx;

			text {
				font-size: 24rpx;
				color: rgba(255, 255, 255, 0.85);
			}
		}
	}

	.header-image {
		position: relative;
		z-index: 1;
		width: 140rpx;
		height: 140rpx;
		margin-left: 24rpx;
		border-radius: 50%;
		overflow: hidden;
		background: rgba(255, 255, 255, 0.25);
		border: 4rpx solid rgba(255, 255, 255, 0.45);
		flex-shrink: 0;

		image {
			width: 100%;
			height: 100%;
		}
	}
}

/* ============ AI 入口 ============ */
.ai-user-service {
	display: flex;
	flex-direction: row;
	gap: 24rpx;

	.ai-answer,
	.ai-medical-service {
		flex: 1;
		background: #ffffff;
		border-radius: 24rpx;
		padding: 32rpx 0 26rpx;
		display: flex;
		flex-direction: column;
		align-items: center;
		box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
		transition: all 0.15s ease;

		image {
			width: 88rpx;
			height: 88rpx;
			padding: 16rpx;
			box-sizing: border-box;
			background: linear-gradient(135deg, #eef3ff, #dbe8ff);
			border-radius: 50%;
			margin-bottom: 18rpx;
		}

		text {
			font-size: 27rpx;
			color: #1d2129;
			font-weight: 500;
		}

		&:active {
			transform: scale(0.96);
			background: #fafbfc;
		}
	}
}

/* ============ 公告 ============ */
.announcement {
	position: relative;
	background: #ffffff;
	border-radius: 24rpx;
	padding: 28rpx;
	box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);

	&::before {
		content: '';
		position: absolute;
		left: 0;
		top: 28rpx;
		bottom: 28rpx;
		width: 6rpx;
		border-radius: 0 6rpx 6rpx 0;
		background: linear-gradient(180deg, #4a7dff, #7aa8ff);
	}

	&:active {
		background: #fafbfc;
	}

	.announce-title {
		display: flex;
		flex-direction: row;
		align-items: center;
		font-size: 30rpx;
		font-weight: 600;
		color: #1d2129;

		image {
			width: 32rpx;
			height: 32rpx;
			margin-right: 12rpx;
			flex-shrink: 0;
		}

		/* 标题过长省略 */
		overflow: hidden;
		white-space: nowrap;
		text-overflow: ellipsis;
	}

	.announce-content {
		margin-top: 16rpx;
		font-size: 26rpx;
		line-height: 1.6;
		color: #86909c;
		/* 最多两行 */
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}

	/* 暂无公告态 */
	&:not(:has(.announce-title)) {
		text-align: center;
		font-size: 26rpx;
		color: #c9cdd4;
		padding: 40rpx 0;
	}
}

/* ============ 推荐区 ============ */
.recommend {
	.recommend-title {
		font-size: 32rpx;
		font-weight: 600;
		color: #1d2129;
		margin: 8rpx 0 20rpx;
		padding-left: 6rpx;
	}

	.recommend-container {
		background: #ffffff;
		border-radius: 24rpx;
		padding: 24rpx;
		box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
	}

	.rec-empty {
		min-height: 200rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 26rpx;
		color: #c9cdd4;
	}

	.rec-item {
		padding: 20rpx 12rpx;
		border-radius: 12rpx;
		&:not(:last-child) {
			border-bottom: 1rpx solid #f2f3f5;
		}

		&:active {
			background: #fafbfc;
		}
	}

	.rec-head {
		display: flex;
		flex-direction: row;
		align-items: center;
		margin-bottom: 12rpx;
	}

	.rec-avatar {
		width: 56rpx;
		height: 56rpx;
		border-radius: 50%;
		margin-right: 16rpx;
		background: #f0f0f0;
		flex-shrink: 0;
	}

	.rec-head-info {
		flex: 1;
		display: flex;
		flex-direction: column;
		overflow: hidden;
	}

	.rec-author {
		font-size: 26rpx;
		font-weight: 600;
		color: #1d2129;
	}

	.rec-time {
		font-size: 22rpx;
		color: #a9aeb8;
		margin-top: 4rpx;
	}

	.rec-content {
		font-size: 27rpx;
		line-height: 1.55;
		color: #4e5969;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}

	.rec-thumb {
		width: 100%;
		height: 240rpx;
		border-radius: 16rpx;
		margin-top: 16rpx;
		background: #f2f3f5;
	}

	.rec-more {
		padding-top: 20rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 26rpx;
		color: #4a7dff;

		.rec-more-icon {
			width: 24rpx;
			height: 24rpx;
			margin-left: 8rpx;
			opacity: 0.7;
		}
	}
}
</style>
