import {ref} from 'vue'
import {request} from '@/util/request.js'

const _memberList_ = ref([
	{
		name: "666",
		account: "2455555",
		join_is: false
	},
	{
		name: "546",
		account: "465sa4d65",
		join_is: true
	},
	{
		name: "s9d",
		account: "sda56655",
		join_is: false
	},
	{
		name: "萨达",
		account: "5s4d6ad",
		join_is: true
	},
])

const GetCommunityUid = async ()=>{
//--------------------------------------------获取社区Uid
	let account = uni.getStorageSync("user_account")
	try{
		const uid_res = await request({
			url: '/api/community/query/uid',
			data:{account}
		})
		console.log(uid_res)
		return uid_res.info
	} catch(err){
		console.log("err:",err)
	}
}

const GetMemberList = async (_community_uid_)=>{
	//--------------------------------------------获取成员状态
	try{
		const res = await request({
			url: '/api/community/query/member/list',
			data:{uid:_community_uid_}
		})
		const list = res.info || []
		_memberList_.value = []            // 先清空旧数据（含初始 mock）
		const myAccount = uni.getStorageSync("user_account")
		list.forEach((item)=>{
			_memberList_.value.push({
				name: item.account === myAccount ? "我" : item.name,
				account: item.account,
				join_is: item.is_join === 1,   // 后端返回 int 1/0，转成 boolean 给模板用
			})
		})
	} catch(err){
		console.log("err:",err)
	}
}

const RemoveMembers = async (community_uid, accounts)=>{
	//删除成员
	try{
		const res = await request({
			url: '/api/community/query/member/remove',
			method: 'POST',
			data: {
				community_uid: community_uid,
				accounts: accounts.join(',')
			}
		})
		return res
	} catch(err){
		console.log("err:",err)
		throw err
	}
}

const AddMember = async (community_uid, account, member_name)=>{
	try{
		const res = await request({
			url: '/api/community/query/member/add',
			method: 'POST',
			data: {
				community_uid: community_uid,
				account: account,
				member_name: member_name
			}
		})
		return res
	} catch(err){
		console.log("err:",err)
		throw err
	}
}

const JoinCommunityByUid = async (community_uid)=>{
	try{
		const res = await request({
			url: '/api/community/query/member/join',
			method: 'POST',
			data: { community_uid: community_uid }
		})
		return res
	} catch(err){
		console.log("err:",err)
		throw err
	}
}


export {
	_memberList_,
	GetCommunityUid,
	GetMemberList,
	RemoveMembers,
	AddMember,
	JoinCommunityByUid
}
