<script>
	export default {
		onLaunch: function() {
			console.log('App Launch')
		},
		onShow: function() {
			console.log('App Show')
		},
		onHide: function() {
			console.log('App Hide')
		}
	}
</script>

<style>
	/*每个页面公共css */
	page {
		background-color: #f4f6f9;
		font-family: "PingFang SC", "Helvetica Neue", Helvetica, Arial, "Microsoft YaHei", sans-serif;
		color: #1d2129;
		-webkit-font-smoothing: antialiased;
	}

	/* 通用按压反馈 */
	.press {
		transition: transform 0.15s ease, opacity 0.15s ease;
	}
	.press:active {
		transform: scale(0.97);
		opacity: 0.92;
	}
</style>
