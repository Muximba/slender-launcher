package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// Web service do Muximba (mesmo login.php que serve cacheinfo/eventschedule).
// Ajustar se o endpoint mudar.
const muximbaWebService = "https://app.muximba.dev/login.php"

// ServerStatusInfo é o status do servidor exposto pro frontend.
type ServerStatusInfo struct {
	Online        bool `json:"online"`
	PlayersOnline int  `json:"playersOnline"`
	TotalPlayers  int  `json:"totalPlayers"`
}

func (a *App) postWebService(payload string) ([]byte, error) {
	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Post(muximbaWebService, "application/json", bytes.NewReader([]byte(payload)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ServerStatus busca players online via cacheinfo do web service.
func (a *App) ServerStatus() ServerStatusInfo {
	data, err := a.postWebService(`{"type":"cacheinfo"}`)
	if err != nil {
		a.logger.Errorf("ServerStatus error: %v", err)
		return ServerStatusInfo{Online: false}
	}
	var parsed struct {
		PlayersOnline int `json:"playersonline"`
		TotalPlayers  int `json:"totalPlayers"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		a.logger.Errorf("ServerStatus decode: %v", err)
		return ServerStatusInfo{Online: false}
	}
	return ServerStatusInfo{
		Online:        true,
		PlayersOnline: parsed.PlayersOnline,
		TotalPlayers:  parsed.TotalPlayers,
	}
}

// NewsItem é uma entrada de news/aviso do servidor.
type NewsItem struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Content string `json:"content"`
	URL     string `json:"url"`
}

// News busca a lista de news do web service. Retorna vazio se o servidor
// ainda nao tiver o endpoint (type=news) — pronto pra plugar quando o site tiver.
func (a *App) News() []NewsItem {
	data, err := a.postWebService(`{"type":"news"}`)
	if err != nil {
		a.logger.Errorf("News error: %v", err)
		return []NewsItem{}
	}
	var parsed struct {
		News []NewsItem `json:"news"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return []NewsItem{}
	}
	return parsed.News
}
