<template>
	<view class="home-page">
		<view class="Welcome">
			<view class="greet">
				<view class="time-domain">
					{{ domain }}
				</view>
				<view class="community-greet">
					<text>欢迎使用 智慧社区 AI 服务平台</text>
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
			<view class="recommend-container">推送内容</view>
		</view>
	</view>
</template>

<script setup>
import { onMounted, onUnmounted, ref } from 'vue';
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



onMounted(() => {
	currentHour.value = getCurrentHours()
	timer = setInterval(() => {
		currentHour.value = getCurrentHours()
		GreetTime()
	}, 1000)
	loadLatestNotice()
})


onUnmounted(() => {
	if(timer)
		clearInterval(timer)
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
	display: flex;
	flex-direction: row;
	align-items: center;
	justify-content: space-between;
	padding: 40rpx 32rpx;
	border-radius: 28rpx;
	background: linear-gradient(135deg, #4a7dff 0%, #7aa8ff 100%);
	box-shadow: 0 12rpx 32rpx rgba(74, 125, 255, 0.25);

	.greet {
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
	background: #ffffff;
	border-radius: 24rpx;
	padding: 28rpx;
	box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);

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
		min-height: 240rpx;
		padding: 28rpx;
		font-size: 28rpx;
		color: #86909c;
		box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
	}
}
</style>
