package server

import (
	"Service/Tools/JWT"
	"Service/server/advance"
	"Service/server/chat"
	"Service/server/community"
	"Service/server/image"
	"Service/server/notice"
	"Service/server/report"
	"Service/server/share"
	"Service/server/user"
	"database/sql"
	"github.com/gin-gonic/gin"
)

func StartServer(r *gin.Engine, db *sql.DB) {
	user.StartUserServeAPI(r, db)
	JWT.JwtExpirationMonitoring(r)
	chat.WebsocketServer(r, db)
	community.HandleCommunity(r, db)
	notice.HandleNotice(r, db)
	advance.HandleAdvance(r, db)
	share.HandleShare(r, db)
	image.HandleImageQiNiu(r, db)
	report.HandleReport(r, db)

}
