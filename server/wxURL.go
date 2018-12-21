package server

func wxURL() {
	wx.GET("/login", Login)
	wx.GET("/bingtg", JWTAuth(), BindTg)
	wx.GET("/UnBindTg", JWTAuth(), UnBindTg)
	wx.GET("/GetMyGifs", JWTAuth(), GetMyGifs)
}
