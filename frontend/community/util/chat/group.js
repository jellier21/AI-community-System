import { ref } from "vue";
import { request } from '@/util/request.js'

const group_list = ref([])

const updateGroupList = async ()=>{
	try {
		const res = await request({
			url: '/api/chat/group/list' ,
			data:{
				account: uni.getStorageSync('user_account')
			}
			})
		console.log('room_list:', res)
		if(res.code === 200){
			group_list.value = []
			let temp_list = res.room_list
			for(let i = 0;i<temp_list.length;i++){
				group_list.value.push({
					type:"group",
					data:{
						room_id: temp_list[i].room_uid,
						room_name: temp_list[i].room_name,
						leader: temp_list[i].leader,
					},
					created_at: temp_list[i].created_at
				})
			}
		}
		return false
	} catch (err) {
		console.log('err:', err)
		return true
	}
}

// ====== 私聊 ======
const createPrivateChat = async (friendAccount) => {
	return await request({
		url: '/api/chat/group/private/create',
		method: 'POST',
		data: {
			account: uni.getStorageSync('user_account'),
			friend_account: friendAccount
		}
	})
}

// ====== 群聊 ======
const createGroupChat = async (groupName, members) => {
	return await request({
		url: '/api/chat/group/create',
		method: 'POST',
		data: {
			account: uni.getStorageSync('user_account'),
			group_name: groupName,
			members: members.join(',')
		}
	})
}

// ====== 好友 ======
const getFriendList = async () => {
	const res = await request({
		url: '/api/chat/group/friends',
		data: { account: uni.getStorageSync('user_account') }
	})
	return res.friends || []
}

const deleteFriend = async (roomUid) => {
	return await request({
		url: '/api/chat/group/friend/delete',
		method: 'POST',
		data: { room_uid: roomUid }
	})
}

export {
	group_list,
	updateGroupList,
	createPrivateChat,
	createGroupChat,
	getFriendList,
	deleteFriend
}
