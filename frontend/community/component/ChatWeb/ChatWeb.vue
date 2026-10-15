<template>
	<view class="chat-container">
		<!-- 顶部品牌横幅 -->
		<view class="chat-header">
			<view class="switch-tools-sign" :class="{ 'is-open': tool_list }" @click="startPushList">
				<image src="/static/chat/push.png" mode="aspectFit"></image>
			</view>
			<view class="header-right">
				<text class="chat-header-title">消息</text>
				<text class="chat-header-sub">{{ group_list.length }} 个会话</text>
			</view>
		</view>

		<!-- 会话列表 -->
		<view class="notice-container">
			<view class="notice-list" v-for="(item,index) in group_list" :key="index">
				<view class="chat-item" hover-class="chat-item-hover" @click="goChat(item)">
					<view class="chat-avatar" :class="item.type === 'message' ? 'avatar-message' : 'avatar-group'">
						{{ item.data.room_name.slice(0, 1) }}
					</view>
					<view class="chat-name">{{ item.data.room_name }}</view>
				</view>
			</view>
		</view>

		<!-- 全屏遮罩 + 下拉菜单 -->
		<view class="shade-index" @click="closeToolList" v-show="tool_list">
			<view class="tool-list">
				<view class="tool-build-connect" @click="openPrivateDialog">私聊</view>
				<view class="tool-build-group" @click="goCreateGroup">建立群聊</view>
				<view class="friend-manage" @click="goFriendList">好友管理</view>
			</view>
		</view>
	</view>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { group_list, updateGroupList, createPrivateChat } from '@/util/chat/group';
import { onPullDownRefresh } from '@dcloudio/uni-app';

onMounted(async () => {
	await updateGroupList()
})

// 下拉刷新：刷新会话/群聊列表
onPullDownRefresh(async () => {
	await updateGroupList()
	uni.stopPullDownRefresh()
})

//=========小工具按钮===========
const tool_list = ref(false)
const startPushList = ()=>{
	tool_list.value = !tool_list.value
}
const closeToolList = ()=>{
	tool_list.value = false
}

// 进入会话
const goChat = (item) => {
	if (item.type === 'message') {
		uni.navigateTo({ url: '/component/chatPage/chat/chat' })
	} else {
		uni.navigateTo({
			url: `/component/chatPage/chat/chat?uid=${item.data.room_id}&leader=${item.data.leader}`
		})
	}
}

// 私聊：弹窗输入对方账号
const openPrivateDialog = () => {
	uni.showModal({
		title: '发起私聊',
		editable: true,
		placeholderText: '输入对方账号',
		success: async (res) => {
			closeToolList()
			if (!res.confirm || !res.content) return
			try {
				const r = await createPrivateChat(res.content.trim())
				if (r.code === 200) {
					uni.showToast({ title: '已创建', icon: 'success' })
					setTimeout(() => {
						uni.navigateTo({
							url: '/component/chatPage/chat/chat?uid=' + r.room_uid
						})
					}, 600)
				} else {
					uni.showToast({ title: r.error || '创建失败', icon: 'none' })
				}
			} catch(e) {
				uni.showToast({ title: '创建失败', icon: 'none' })
			}
		}
	})
}

const goCreateGroup = () => {
	closeToolList()
	uni.navigateTo({ url: '/component/chatPage/create_group/create_group' })
}

const goFriendList = () => {
	closeToolList()
	uni.navigateTo({ url: '/component/chatPage/friend_list/friend_list' })
}
</script>

<style lang="scss" scoped>
.chat-container {
  min-height: 100vh;
  box-sizing: border-box;
  padding: calc(var(--status-bar-height) + 24rpx) 24rpx calc(40rpx + env(safe-area-inset-bottom));
  background: #f5f6f8;

  /* ---------- 顶部品牌横幅 ---------- */
  .chat-header {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 20rpx;
    padding: 40rpx 36rpx;
    border-radius: 28rpx;
    background: linear-gradient(135deg, #4a7dff 0%, #7aa8ff 100%);
    box-shadow: 0 12rpx 32rpx rgba(74, 125, 255, 0.28);

    .header-right {
      display: flex;
      flex-direction: column;
      align-items: flex-end;

      .chat-header-title {
        font-size: 44rpx;
        font-weight: 700;
        color: #ffffff;
        letter-spacing: 2rpx;
      }

      .chat-header-sub {
        margin-top: 12rpx;
        font-size: 24rpx;
        color: rgba(255, 255, 255, 0.85);
      }
    }

    /* 右上角“+”按钮 */
    .switch-tools-sign {
      width: 72rpx;
      height: 72rpx;
      border-radius: 50%;
      background: rgba(255, 255, 255, 0.94);
      box-shadow: 0 4rpx 14rpx rgba(0, 0, 0, 0.15);
      display: flex;
      align-items: center;
      justify-content: center;
      transition: transform 0.2s ease;
      flex-shrink: 0;

      &:active {
        transform: scale(0.92);
      }

      image {
        width: 36rpx;
        height: 36rpx;
        transition: transform 0.25s ease;
      }

      /* 展开时图标旋转 45°，做“+”→“×”的过渡 */
      &.is-open {
        background: #ffffff;
        image {
          transform: rotate(45deg);
        }
      }
    }
  }

  /* ---------- 会话列表 ---------- */
  .notice-container {
    padding-top: 8rpx;

    .notice-list {
      margin-bottom: 20rpx;
      background: #ffffff;
      border-radius: 20rpx;
      box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
      overflow: hidden;
    }

    .chat-item {
      display: flex;
      align-items: center;
      padding: 28rpx;
      transition: background 0.15s ease;

      /* 会话头像 */
      .chat-avatar {
        width: 84rpx;
        height: 84rpx;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #ffffff;
        font-size: 32rpx;
        font-weight: 600;
        margin-right: 24rpx;
        flex-shrink: 0;

        &.avatar-message {
          background: linear-gradient(135deg, #4a7dff, #7aa8ff);
        }

        &.avatar-group {
          background: linear-gradient(135deg, #22b573, #5ad095);
        }
      }

      /* 会话名称 */
      .chat-name {
        flex: 1;
        min-width: 0;
        font-size: 30rpx;
        color: #333333;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
      }

      /* 右侧箭头 */
      &::after {
        content: '';
        width: 14rpx;
        height: 14rpx;
        border-top: 3rpx solid #c8c9cc;
        border-right: 3rpx solid #c8c9cc;
        transform: rotate(45deg);
        margin-left: 12rpx;
        flex-shrink: 0;
      }
    }
  }

  /* ---------- 全屏遮罩 + 下拉菜单 ---------- */
  .shade-index {
    position: fixed;
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
    z-index: 999;
    background: rgba(226, 226, 226, 0.4);
    animation: shade-in 0.18s ease-out;

    .tool-list {
      position: absolute;
      top: calc(var(--status-bar-height) + 144rpx);
      left: 60rpx;
      width: 248rpx;
      padding: 8rpx 0;
      background: #ffffff;
      border-radius: 18rpx;
      box-shadow: 0 12rpx 36rpx rgba(0, 0, 0, 0.18);
      transform-origin: top left;
      animation: menu-in 0.18s ease-out;

      /* 顶部小三角 */
      &::before {
        content: '';
        position: absolute;
        top: -12rpx;
        left: 36rpx;
        width: 0;
        height: 0;
        border-left: 12rpx solid transparent;
        border-right: 12rpx solid transparent;
        border-bottom: 14rpx solid #ffffff;
      }

      .tool-build-connect,
      .tool-build-group,
      .friend-manage {
        position: relative;
        padding: 24rpx 24rpx 24rpx 76rpx;
        font-size: 28rpx;
        color: #333333;
        line-height: 1.2;

        &:active {
          background: #f2f3f5;
        }

        /* 左侧圆形图标底色 */
        &::before {
          content: '';
          position: absolute;
          left: 24rpx;
          top: 50%;
          transform: translateY(-50%);
          width: 36rpx;
          height: 36rpx;
          border-radius: 50%;
        }

        /* 图标字符 */
        &::after {
          position: absolute;
          left: 24rpx;
          top: 50%;
          transform: translateY(-50%);
          width: 36rpx;
          text-align: center;
          font-size: 22rpx;
          font-weight: bold;
          color: #ffffff;
        }
      }

      .tool-build-connect {
        &::before { background: #3b7cff; }
        &::after  { content: '+'; }
      }

      .tool-build-group {
        &::before { background: #22b573; }
        &::after  { content: '群'; font-size: 20rpx; }
      }

      .friend-manage {
        margin-top: 6rpx;
        border-top: 1rpx solid #f0f0f0;
        &::before { background: #ff9f2e; }
        &::after  { content: '管'; font-size: 20rpx; }
      }
    }
  }
}

@keyframes shade-in {
  from { opacity: 0; }
  to   { opacity: 1; }
}

@keyframes menu-in {
  from {
    opacity: 0;
    transform: scale(0.9) translateY(-8rpx);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}
</style>

<!-- 全局按压反馈：hover-class 在微信端需全局样式才生效 -->
<style lang="scss">
.chat-item-hover {
  background: #f7f8fa;
}
</style>
