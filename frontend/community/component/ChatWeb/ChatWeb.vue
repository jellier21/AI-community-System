<template>
	<view class="chat-container">
		<view class="user-tool">
			<view class="switch-tools-sign" @click="startPushList">
				<image src="/static/chat/push.png" mode="aspectFit"></image>
			</view>
			<view class="shade-index" @click="closeToolList" v-show="tool_list">
				<view class="tool-list">
					<view class="tool-build-connect" @click="openPrivateDialog">私聊</view>
					<view class="tool-build-group" @click="goCreateGroup">建立群聊</view>
					<view class="friend-manage" @click="goFriendList">好友管理</view>
				</view>
			</view>
		</view>
		<view class="notice-container">
			<view class="notice-list" v-for="(item,index) in group_list">
				<view class="notice-item" v-if="item.type === 'message'">
					<navigator url="/component/chatPage/chat/chat">{{ item.data.room_name }}</navigator>
				</view>
				<!-- <view class="chat-item" v-else>
					<navigator url="/component/chatPage/chat/chat">{{ item.data.room_name }}</navigator>
				</view> -->
				<navigator :url="`/component/chatPage/chat/chat?uid=${item.data.room_id}&leader=${item.data.leader}`" class="chat-item" v-else>
					{{ item.data.room_name }}
				</navigator>
			</view>
		</view>
	</view>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { group_list, updateGroupList, createPrivateChat } from '@/util/chat/group';

onMounted(async () => {
	await updateGroupList()
})

//=========小工具按钮===========
const tool_list = ref(false)
const startPushList = ()=>{
	tool_list.value = !tool_list.value
}
const closeToolList = ()=>{
	tool_list.value = false
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

  /* ---------- 右上角工具区 ---------- */
  .user-tool {
    position: fixed;
    top: calc(var(--status-bar-height) + 24rpx);
    left: 24rpx;
    z-index: 999;

    /* 触发按钮：必须高于遮罩，否则打开后点不到 */
    .switch-tools-sign {
      position: relative;
      z-index: 3;
      width: 76rpx;
      height: 76rpx;
      border-radius: 50%;
      background: #ffffff;
      box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.12);
      display: flex;
      align-items: center;
      justify-content: center;
      transition: transform 0.2s ease;

      &:active {
        transform: scale(0.92);
      }

      image {
        width: 40rpx;
        height: 40rpx;
        transition: transform 0.25s ease;
      }

      /* 展开时图标旋转 45°，做“+”→“×”的过渡 */
      &.is-open {
        background: #3b7cff;
        image {
          transform: rotate(45deg);
        }
      }
    }

    /* ---------- 全屏遮罩 ---------- */
    .shade-index {
      position: fixed;
      top: 0;
      right: 0;
      bottom: 0;
      left: 0;
      z-index: 2;
      background: rgba(226, 226, 226, 0.4);
      /* 小程序端遮罩不响应滚动穿透，H5 可用 backdrop-filter */
      backdrop-filter: blur(2rpx);
      animation: shade-in 0.18s ease-out;

      /* ---------- 下拉菜单 ---------- */
      .tool-list {
        position: absolute;
        top: 112rpx;
        left: 24rpx;
        width: 248rpx;
        padding: 8rpx 0;
        background: #ffffff;
        border-radius: 18rpx;
        box-shadow: 0 12rpx 36rpx rgba(0, 0, 0, 0.18);
        transform-origin: top right;
        animation: menu-in 0.18s ease-out;

        /* 顶部小三角 */
        &::before {
          content: '';
          position: absolute;
          top: -12rpx;
          right: 26rpx;
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

  /* ---------- 列表区 ---------- */
  .notice-container {
    padding-top: 112rpx;

    .notice-list {
      margin-bottom: 20rpx;
      background: #ffffff;
      border-radius: 18rpx;
      box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
      overflow: hidden;

      &:active {
        background: #fafafa;
      }
    }

    /* 新消息通知 */
    .notice-item {
      position: relative;
      padding: 30rpx 72rpx 30rpx 40rpx;
      font-size: 30rpx;
      font-weight: 500;
      color: #222222;

      &::before {
        content: '';
        position: absolute;
        left: 0;
        top: 0;
        bottom: 0;
        width: 8rpx;
        background: #ff7a45;
      }

      &::after {
        content: '';
        position: absolute;
        right: 32rpx;
        top: 50%;
        transform: translateY(-50%);
        width: 16rpx;
        height: 16rpx;
        border-radius: 50%;
        background: #fa5151;
      }
    }

    /* 普通会话 */
    .chat-item {
      position: relative;
      padding: 30rpx 72rpx 30rpx 40rpx;
      font-size: 30rpx;
      color: #666666;

      &::after {
        content: '';
        position: absolute;
        right: 32rpx;
        top: 50%;
        width: 14rpx;
        height: 14rpx;
        border-top: 3rpx solid #c8c9cc;
        border-right: 3rpx solid #c8c9cc;
        transform: translateY(-50%) rotate(45deg);
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
