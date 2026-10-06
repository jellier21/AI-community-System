<template>
	<view>
		<view class="tool-container">
			<text class="title">工具箱</text>
			<view class="tool-area">
				<view class="medical-serve" @click="uni.navigateTo({
					url:'/component/AI_page/medical/medical'
				})">
					<image src="/static/tools/medical.png" mode="aspectFit"></image>
					<text>医疗服务</text>
				</view>
				<view class="knowledge-serve" @click="uni.navigateTo({
					url:'/component/AI_page/mondo/mondo'
				})">
					<image src="/static/tools/knownledge.png" mode="aspectFit"></image>
					<text>生活小知识</text>
				</view>
				<view class="understand-image" @click="uni.navigateTo({
					url:'/component/ToolsPage/user/MapReading/MapReading'
				})">
					<image src="/static/tools/mapReading.png" mode="aspectFit"></image>
					<text>拍照识图</text>
				</view>

			</view>
		</view>
		<view class="member-tools">
			<view class="submit-advance" @click="uni.navigateTo({
				url: '/component/ToolsPage/user/advance-submit/advance-submit'
			})">
				<image src="/static/tools/advance.png" mode="aspectFit"></image>
				<text>建议上报</text>
			</view>
			<view class="community-file" @click="uni.navigateTo({
				url:'/component/ToolsPage/user/share/share'
			})">
				<image src="/static/tools/community.png" mode="aspectFit"></image>
				<text>社区圈</text>
			</view>
		</view>
		<view class="special-manager" v-show="Is_Leader">
			<text class="title">管理员工具箱</text>
			<view class="manager-container">
				<view class="member" @click="uni.navigateTo({url:'/component/ToolsPage/LeaderTools/member_manage/member_manage'})">
					<image src="/static/tools/member-manage.png" mode="aspectFit"></image>
					<text>成员管理</text>
				</view>
				<view class="notice" @click="uni.navigateTo({
					url: '/component/ToolsPage/LeaderTools/newNotice/newNotice'
				})">
					<image src="/static/tools/notice.png" mode="aspectFit"></image>
					<text>公告管理</text>
				</view>
				<view class="suggest" @click="uni.navigateTo({
					url:'/component/ToolsPage/LeaderTools/advance_statistics/advance_statistics'
				})">
					<image src="/static/tools/suggestion.png" mode="aspectFit"></image>
					<text>建议统计</text>
				</view>
				<view class="weekly-news" @click="uni.navigateTo({
					url:'/component/ToolsPage/LeaderTools/FinalReport/FinalReport'
				})">
					<image src="/static/tools/week_news.png" mode="aspectFit"></image>
					<text>AI周报</text>
				</view>

				<view class="weekly-file" @click="uni.navigateTo({
					url:'/component/ToolsPage/LeaderTools/share_manager/share_manager'
				})">
					<image src="/static/tools/file.png" mode="aspectFit"></image>
					<text>社区圈管理</text>
				</view>
			</view>
		</view>
	</view>
</template>

<script setup>
import {ref} from 'vue'
import {request} from '@/util/request.js'
import { onShow } from '@dcloudio/uni-app'



//是否为负责人
const Is_Leader = ref(false)
const IsLeaderCheck = async ()=>{
	try {
		const IsLeader_res = await request({
			url: "/api/community/query/isleader",
			data:{
				account: uni.getStorageSync("user_account")
			}
		})
		console.log(IsLeader_res)
		if(IsLeader_res.info === 'leader'){
			Is_Leader.value = true
		}

	} catch (err) {
		console.log("error:",err)
	}
}

onShow(()=>{
	IsLeaderCheck()
})
</script>

<style lang="scss">
	page {
		background: #f4f6f9;
		padding-top: var(--status-bar-height);
	}

	/* ---------- 公共卡片容器 ---------- */
	.tool-container,
	.member-tools,
	.special-manager {
		margin: 24rpx;
		padding: 28rpx 24rpx;
		background: #ffffff;
		border-radius: 24rpx;
		box-shadow: 0 8rpx 24rpx rgba(0, 0, 0, 0.05);
	}

	/* ---------- 标题 ---------- */
	.title {
		display: block;
		font-size: 32rpx;
		font-weight: 600;
		color: #222;
		line-height: 1.2;
		margin-bottom: 24rpx;
		padding-left: 16rpx;
		border-left: 8rpx solid #3a7afe;
	}

	/* ---------- 工具项公共样式（图标上、文字下） ---------- */
	.tool-area > view,
	.member-tools > view,
	.manager-container > view {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 28rpx 12rpx;
		background: #f7f9fc;
		border-radius: 18rpx;
		transition: all 0.2s;
	}

	.tool-area > view:active,
	.member-tools > view:active,
	.manager-container > view:active {
		background: #eaf1ff;
		transform: scale(0.97);
	}

	.tool-area image,
	.member-tools image,
	.manager-container image {
		width: 76rpx;
		height: 76rpx;
		margin-bottom: 14rpx;
	}

	.tool-area text,
	.member-tools text,
	.manager-container text {
		font-size: 26rpx;
		color: #333;
	}

	/* ---------- 第一行：医疗服务 / 生活小知识 / 拍照识图 ---------- */
	.tool-area {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		gap: 20rpx;
	}

	/* ---------- 建议提交 / 社区 ---------- */
	.member-tools {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: 20rpx;
	}

	/* ---------- 管理员工具箱：6 列栅格 ---------- */
	.manager-container {
		display: grid;
		grid-template-columns: repeat(6, 1fr); /* 一行 6 个单位 */
		gap: 20rpx;
	}

	/* 第一行 3 个：每个占 2 个单位 */
	.manager-container .member,
	.manager-container .notice,
	.manager-container .suggest {
		grid-column: span 2;
	}

	/* 第二行 2 个：每个占 3 个单位，正好拉满整行 */
	.manager-container .weekly-news,
	.manager-container .weekly-file {
		grid-column: span 3;
		/* 宽卡片改成「图标左 + 文字右」，视觉更饱满 */
		flex-direction: row;
		justify-content: flex-start;
		padding: 28rpx 32rpx;
	}

	.manager-container .weekly-news image,
	.manager-container .weekly-file image {
		margin-bottom: 0;
		margin-right: 20rpx;
	}

	.manager-container .weekly-news text,
	.manager-container .weekly-file text {
		font-size: 28rpx;
	}
</style>
