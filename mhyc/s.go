package mhyc

import "net/http"

//GET https://docater1.cn/index.php?g=Game&m=GameInterface&game_session=0e0f0787ece611ef27ac216976a0cff4&a=game_recent&ua=2

type GameSession struct {
	sessionID string
}

func NewSession(sid string) *GameSession {
	return &GameSession{sessionID: sid}
}

// GetGameRecent 获取游戏的根本信息
func (s *GameSession) GetGameRecent() {
	http.NewRequest(http.MethodGet, "https://docater1.cn/index.php?g=Game&m=GameInterface&game_session="+s.sessionID+"&a=game_recent&ua=2", nil)
}
