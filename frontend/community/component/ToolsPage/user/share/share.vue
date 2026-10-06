<template>
	<view class="container">
		<!-- 顶部导航栏 -->
		<view class="navbar">
			<text class="navbar-title">社区圈</text>
			<view class="add-btn" @click="openPublish">
				<text class="add-icon">+</text>
			</view>
		</view>

		<!-- 社区圈信息流 -->
		<view class="feed">
			<view class="item" v-for="item in list" :key="item.id">
				<view class="item-head">
					<image class="avatar" :src="item.avatar || '/static/logo.png'" mode="aspectFill"></image>
					<view class="head-info">
						<view class="author">{{ item.account === myAccount ? '我' : item.name }}</view>
						<view class="time">{{ item.create_at }}</view>
					</view>
				</view>
				<view class="content">{{ item.content }}</view>
				<view class="img-grid" v-if="item.images">
					<image v-for="(img,i) in item.images.split(',')"
					       :key="i" :src="img" mode="aspectFill"
					       @click="previewImage(item.images.split(','), i)" />
				</view>
				<view class="actions" v-if="item.account === myAccount">
					<text class="action-btn edit" @click="startEdit(item)">编辑</text>
					<text class="action-btn del" @click="handleDelete(item.id)">删除</text>
				</view>
			</view>
			<view class="empty" v-if="list.length === 0">
				<text>还没有动态，点右上角 + 发布第一条吧</text>
			</view>
		</view>

		<!-- 发布弹窗（底部弹出） -->
		<view class="mask" v-if="showPublish" @click="closePublish">
			<view class="publish-sheet" @click.stop>
				<view class="sheet-handle"></view>
				<view class="sheet-title">分享新鲜事</view>
				<textarea class="input" v-model="newContent" placeholder="分享美好生活..." />
				<view class="img-preview-row">
					<view class="img-item" v-for="(img,i) in newImages" :key="i">
						<image :src="img" mode="aspectFill" />
						<text class="img-del" @click="newImages.splice(i,1)">×</text>
					</view>
					<view class="img-add" @click="pickImages" v-if="newImages.length < 9">
						<text>+</text>
					</view>
				</view>
				<button class="btn-publish" @click="handlePublish">发布</button>
			</view>
		</view>

		<!-- 编辑弹窗 -->
		<view class="mask" v-if="editing" @click="editing=null">
			<view class="publish-sheet" @click.stop>
				<view class="sheet-handle"></view>
				<view class="sheet-title">编辑动态</view>
				<textarea class="input" v-model="editContent" />
				<view class="img-preview-row">
					<view class="img-item" v-for="(img,i) in editImages" :key="i">
						<image :src="img" mode="aspectFill" />
						<text class="img-del" @click="editImages.splice(i,1)">×</text>
					</view>
					<view class="img-add" @click="pickEditImages" v-if="editImages.length < 9">
						<text>+</text>
					</view>
				</view>
				<view class="dialog-btns">
					<text class="action-btn cancel" @click="editing=null">取消</text>
					<text class="action-btn edit" @click="confirmEdit">保存</text>
				</view>
			</view>
		</view>
	</view>
</template>

<script setup>
import {ref} from 'vue'
import {onShow} from '@dcloudio/uni-app'
import {
	PublishShare, GetCommunityShares, DeleteShare, UpdateShare,
	chooseAndUploadImages
} from '@/util/tools/share'
import { GetCommunityUid } from '@/util/tools/member'

const showPublish = ref(false)
const newContent = ref('')
const newImages = ref([])
const list = ref([])
const editing = ref(null)
const editContent = ref('')
const editImages = ref([])
const myAccount = ref(uni.getStorageSync('user_account'))

const loadList = async ()=>{
	const uid = await GetCommunityUid()
	if(uid) list.value = await GetCommunityShares(uid)
}

const openPublish = ()=>{
	newContent.value = ''
	newImages.value = []
	showPublish.value = true
}
const closePublish = ()=>{
	showPublish.value = false
}

const pickImages = async ()=>{
	const urls = await chooseAndUploadImages(9 - newImages.value.length)
	newImages.value.push(...urls)
}

const pickEditImages = async ()=>{
	const urls = await chooseAndUploadImages(9 - editImages.value.length)
	editImages.value.push(...urls)
}

const handlePublish = async ()=>{
	if(!newContent.value.trim() && newImages.value.length === 0){
		uni.showToast({ title: '说点什么吧', icon: 'none' })
		return
	}
	const communityUid = await GetCommunityUid()
	try{
		await PublishShare(
			communityUid,
			newContent.value.trim(),
			newImages.value.join(',')
		)
		uni.showToast({ title: '发布成功', icon: 'success' })
		showPublish.value = false
		loadList()
	}catch(e){
		uni.showToast({ title: '发布失败', icon: 'none' })
	}
}

const handleDelete = async (id)=>{
	uni.showModal({
		title: '确认删除',
		success: async (r)=>{
			if(!r.confirm) return
			try{
				await DeleteShare(id)
				uni.showToast({ title: '已删除' })
				loadList()
			}catch(e){
				uni.showToast({ title: '删除失败', icon: 'none' })
			}
		}
	})
}

const startEdit = (item)=>{
	editing.value = item
	editContent.value = item.content
	editImages.value = item.images ? item.images.split(',') : []
}

const confirmEdit = async ()=>{
	try{
		await UpdateShare(
			editing.value.id,
			editContent.value.trim(),
			editImages.value.join(',')
		)
		uni.showToast({ title: '已保存' })
		editing.value = null
		loadList()
	}catch(e){
		uni.showToast({ title: '保存失败', icon: 'none' })
	}
}

const previewImage = (urls, current)=>{
	uni.previewImage({ urls, current })
}

onShow(loadList)
</script>

<style scoped>
.container {
	min-height: 100vh;
	background: #f4f5f7;
}

/* 顶部导航 */
.navbar {
	position: sticky; top: 0; z-index: 100;
	background: #fff;
	padding: 20rpx 32rpx;
	display: flex; align-items: center; justify-content: space-between;
	box-shadow: 0 2rpx 12rpx rgba(0,0,0,0.04);
}
.navbar-title {
	font-size: 36rpx; font-weight: 700; color: #1a1a1a;
}
.add-btn {
	width: 64rpx; height: 64rpx; border-radius: 50%;
	background: linear-gradient(135deg, #4a7dff, #6ea3ff);
	display: flex; align-items: center; justify-content: center;
	box-shadow: 0 4rpx 12rpx rgba(74,125,255,0.3);
}
.add-icon {
	color: #fff; font-size: 40rpx; font-weight: 300; line-height: 1;
	margin-top: -4rpx;
}

/* 信息流 */
.feed { padding: 20rpx 24rpx; }
.item {
	background: #fff; border-radius: 20rpx;
	padding: 24rpx; margin-bottom: 20rpx;
	box-shadow: 0 2rpx 12rpx rgba(0,0,0,0.03);
}
.item-head {
	display: flex; align-items: center; margin-bottom: 16rpx;
}
.avatar {
	width: 72rpx; height: 72rpx; border-radius: 50%;
	margin-right: 16rpx;
	background: #f0f0f0;
}

.head-info { flex: 1; }
.author { font-size: 28rpx; font-weight: 600; color: #222; }
.time { font-size: 22rpx; color: #999; margin-top: 4rpx; }
.content { font-size: 29rpx; color: #333; line-height: 1.6; margin-bottom: 16rpx; }
.img-grid {
	display: flex; flex-wrap: wrap; gap: 12rpx; margin-bottom: 16rpx;
}
.img-grid image {
	width: 210rpx; height: 210rpx; border-radius: 12rpx;
}
.actions {
	display: flex; gap: 40rpx;
	padding-top: 16rpx; border-top: 1rpx solid #f5f6f8;
}
.action-btn { font-size: 26rpx; color: #4a7dff; }
.action-btn.cancel { color: #999; }
.action-btn.del { color: #f53f3f; }

.empty {
	text-align: center; color: #999; padding: 120rpx 0; font-size: 26rpx;
}

/* 底部弹窗 */
.mask {
	position: fixed; top:0; left:0; right:0; bottom:0;
	background: rgba(0,0,0,0.4);
	z-index: 999;
	display: flex; align-items: flex-end;
}
.publish-sheet {
	width: 100%;
	background: #fff;
	border-radius: 32rpx 32rpx 0 0;
	padding: 20rpx 32rpx 48rpx;
	box-sizing: border-box;
}
.sheet-handle {
	width: 60rpx; height: 6rpx; background: #ddd;
	border-radius: 3rpx; margin: 0 auto 24rpx;
}
.sheet-title {
	font-size: 32rpx; font-weight: 600; color: #222;
	text-align: center; margin-bottom: 24rpx;
}
.input {
	width: 100%; box-sizing: border-box;
	background: #f7f8fa; border-radius: 16rpx;
	padding: 24rpx; font-size: 28rpx; height: 200rpx;
	margin-bottom: 20rpx;
}
.img-preview-row {
	display: flex; flex-wrap: wrap; gap: 16rpx; margin-bottom: 24rpx;
}
.img-item {
	width: 150rpx; height: 150rpx; position: relative;
}
.img-item image {
	width: 100%; height: 100%; border-radius: 12rpx;
}
.img-del {
	position: absolute; top: -10rpx; right: -10rpx;
	width: 36rpx; height: 36rpx; line-height: 36rpx;
	background: rgba(0,0,0,0.6); color: #fff;
	border-radius: 50%; text-align: center; font-size: 24rpx;
}
.img-add {
	width: 150rpx; height: 150rpx;
	border: 2rpx dashed #ccc; border-radius: 12rpx;
	display: flex; align-items: center; justify-content: center;
	font-size: 60rpx; color: #999;
}
.btn-publish {
	background: linear-gradient(135deg, #4a7dff, #6ea3ff);
	color: #fff; border-radius: 999rpx;
	font-size: 30rpx; padding: 20rpx 0;
}
.dialog-btns {
	display: flex; justify-content: flex-end; gap: 40rpx; margin-top: 10rpx;
}
</style>
