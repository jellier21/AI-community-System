const TurnToLogin = ()=>{
	uni.reLaunch({
		url:'/pages/index/index'
	})
}

const GoHome = ()=>{
	uni.reLaunch({
		url:'/pages/home/home'
	})
}


export {
	TurnToLogin,
	GoHome
}
