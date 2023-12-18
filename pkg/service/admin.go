package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	music "github.com/bear1278/MusicWave"
	"github.com/bear1278/MusicWave/pkg/repository"
	"github.com/signintech/gopdf"
	"mime/multipart"
	"os"
)

const (
	pathPDF            = "./reports/report.pdf"
	pathExcel          = "./reports/report.xlsx"
	pathJSON           = "./reports/export.json"
	fontPath           = "C:\\Windows\\Fonts\\arial.ttf"
	usersTable         = "users"
	genresTable        = "genres"
	userGenreTable     = "user_genre"
	playlistTable      = "playlists"
	userPlaylistTable  = "user_playlist"
	trackTable         = "tracks"
	playlistTrackTable = "playlist_track"
	albumTable         = "Albums"
	artistTable        = "Artists"
	artistAlbumTable   = "Artist_album"
	artistGenreTable   = "Artist_Genre"
	userAlbumTable     = "User_Album"
	userArtistTable    = "User_Artist"
	adminHistory       = "adminhistory"
)

type AdminServiceImpl struct {
	repo repository.AdminRepo
}

func NewAdminServiceImpl(repo repository.AdminRepo) *AdminServiceImpl {
	return &AdminServiceImpl{repo: repo}
}

func (a *AdminServiceImpl) CheckAdmin(userID int64) (bool, error) {
	adminID, err := a.repo.SelectAdminId()
	if err != nil {
		return false, err
	}
	if adminID != userID {
		return false, err
	} else {
		return true, err
	}
}

func (a *AdminServiceImpl) GetHistory() ([]music.History, error) {
	return a.repo.SelectHistory()
}

func (a *AdminServiceImpl) GetReportPDF() (string, error) {
	var content string
	history, err := a.repo.SelectHistory()
	if err != nil {
		return "", err
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	err = pdf.AddTTFFont("arial", fontPath)
	if err != nil {

		return "", err
	}

	err = pdf.SetFont("arial", "", 14)
	if err != nil {

		return "", err
	}

	for _, h := range history {
		content = fmt.Sprintf("ID: %d\nName: %s\nType: %s\nReason: %s\nDeletedDate: %s\n",
			h.ID, h.Name, h.Type, h.Reason, h.DeletedDate.Format("2006-01-02 15:04:05"))

		if pdf.GetY() > 800 {
			pdf.AddPage() // Добавляем новую страницу
		}
		width, err := pdf.MeasureTextWidth(content)
		if width > 550 {
			content = fmt.Sprintf("ID: %d\nName: %s\nType: %s \nDeletedDate: %s\n",
				h.ID, h.Name, h.Type, h.DeletedDate.Format("2006-01-02 15:04:05"))
			err = pdf.Cell(nil, content)

			if err != nil {
				return "", err
			}
			pdf.Br(20)
			content = fmt.Sprintf("Reason: %s\n", h.Reason)
			err = pdf.Cell(nil, content)

			if err != nil {
				return "", err
			}
			pdf.Br(20)
		} else {
			if err != nil {
				return "", err
			}
			err = pdf.Cell(nil, content)
		}

		if err != nil {
			return "", err
		}
		pdf.Br(20)
	}

	err = pdf.WritePdf(pathPDF)
	if err != nil {

		return "", err
	}
	return pathPDF, err
}

func (a *AdminServiceImpl) GetReportExcel() (string, error) {
	file := excelize.NewFile()

	// Создаем новый лист
	sheetName := "Report"
	index := file.NewSheet(sheetName)

	// Устанавливаем заголовки
	headers := []string{"ID", "Name", "Type", "Reason", "DeletedDate"}
	for colNum, header := range headers {
		cell := excelize.ToAlphaString(colNum+1) + "1"
		file.SetCellValue(sheetName, cell, header)
	}

	// Добавляем данные
	history, err := a.repo.SelectHistory()
	if err != nil {
		return "", err
	}

	for rowNum, record := range history {
		colNum := 1
		file.SetCellValue(sheetName, excelize.ToAlphaString(colNum)+fmt.Sprintf("%d", rowNum+2), record.ID)
		colNum++
		file.SetCellValue(sheetName, excelize.ToAlphaString(colNum)+fmt.Sprintf("%d", rowNum+2), record.Name)
		colNum++
		file.SetCellValue(sheetName, excelize.ToAlphaString(colNum)+fmt.Sprintf("%d", rowNum+2), record.Type)
		colNum++
		file.SetCellValue(sheetName, excelize.ToAlphaString(colNum)+fmt.Sprintf("%d", rowNum+2), record.Reason)
		colNum++
		file.SetCellValue(sheetName, excelize.ToAlphaString(colNum)+fmt.Sprintf("%d", rowNum+2), record.DeletedDate.Format("2006-01-02 15:04:05"))
	}

	// Устанавливаем активный лист
	file.SetActiveSheet(index)

	// Сохраняем файл
	if err := file.SaveAs(pathExcel); err != nil {
		return "", err
	}
	return pathExcel, err
}

func (a *AdminServiceImpl) GetDBInJSON() (string, error) {
	resultUsersTable, err := a.repo.SelectTable(usersTable)
	if err != nil {
		return "", err
	}
	resultGenresTable, err := a.repo.SelectTable(genresTable)
	if err != nil {
		return "", err
	}
	resultUserGenreTable, err := a.repo.SelectTable(userGenreTable)
	if err != nil {
		return "", err
	}
	resultPlaylistTable, err := a.repo.SelectTable(playlistTable)
	if err != nil {
		return "", err
	}
	resultUserPlaylistTable, err := a.repo.SelectTable(userPlaylistTable)
	if err != nil {
		return "", err
	}
	resultTrackTable, err := a.repo.SelectTable(trackTable)
	if err != nil {
		return "", err
	}
	resultPlaylistTrackTable, err := a.repo.SelectTable(playlistTrackTable)
	if err != nil {
		return "", err
	}
	resultAlbumTable, err := a.repo.SelectTable(albumTable)
	if err != nil {
		return "", err
	}
	resultArtistTable, err := a.repo.SelectTable(artistTable)
	if err != nil {
		return "", err
	}
	resultArtistAlbumTable, err := a.repo.SelectTable(artistAlbumTable)
	if err != nil {
		return "", err
	}
	resultArtistGenreTable, err := a.repo.SelectTable(artistGenreTable)
	if err != nil {
		return "", err
	}
	resultUserAlbumTable, err := a.repo.SelectTable(userAlbumTable)
	if err != nil {
		return "", err
	}
	resultUserArtistTable, err := a.repo.SelectTable(userArtistTable)
	if err != nil {
		return "", err
	}
	resultAdminHistory, err := a.repo.SelectTable(adminHistory)
	if err != nil {
		return "", err
	}

	// Объединение результатов
	result := map[string]interface{}{
		usersTable:         resultUsersTable,
		genresTable:        resultGenresTable,
		userGenreTable:     resultUserGenreTable,
		playlistTable:      resultPlaylistTable,
		userPlaylistTable:  resultUserPlaylistTable,
		trackTable:         resultTrackTable,
		playlistTrackTable: resultPlaylistTrackTable,
		albumTable:         resultAlbumTable,
		artistTable:        resultArtistTable,
		artistAlbumTable:   resultArtistAlbumTable,
		artistGenreTable:   resultArtistGenreTable,
		userAlbumTable:     resultUserAlbumTable,
		userArtistTable:    resultUserArtistTable,
		adminHistory:       resultAdminHistory,
	}

	// Запись данных в JSON файл
	file, err := os.Create(pathJSON)
	if err != nil {
		return "", err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(result); err != nil {
		return "", err
	}

	return pathJSON, err
}

func (a *AdminServiceImpl) ImportDBInJSON(file *multipart.FileHeader) error {
	var data map[string]interface{}
	var users []music.User
	var albums []music.Album
	var artists []music.Artist
	var genres []music.Genre
	var playlists []music.Playlist
	var tracks []music.Track
	var history []music.History
	var ags []music.ArtistGenre
	var aas []music.ArtistAlbum
	var uas []music.UserArtist
	var uals []music.UserAlbum
	var pts []music.PlaylistTrack
	var ugs []music.UserGenre
	var ups []music.UserPlaylist
	uploadedFile, err := file.Open()
	if err != nil {
		return err
	}
	defer uploadedFile.Close()

	decoder := json.NewDecoder(uploadedFile)
	if err = decoder.Decode(&data); err != nil {
		return err
	}

	var valid bool = false

	if usersData, ok := data[usersTable].([]interface{}); ok {
		valid = true
		for _, userData := range usersData {
			userJSON, _ := json.Marshal(userData)
			var userJson music.UserJSON
			var user music.User
			if err := json.Unmarshal(userJSON, &userJson); err == nil {
				err = user.SetFromJSON(userJson)
				users = append(users, user)
			}
			if err != nil {
				return err
			}
		}
		for _, u := range users {
			err = a.repo.InsertUser(u)
			if err != nil {
				return err
			}
		}
	}

	if artistsData, ok := data[artistTable].([]interface{}); ok {
		valid = true
		for _, artistData := range artistsData {
			JSON, _ := json.Marshal(artistData)
			var artistJson music.ArtistJSON
			var artist music.Artist
			if err := json.Unmarshal(JSON, &artistJson); err == nil {
				err = artist.SetFromJSON(artistJson)
				artists = append(artists, artist)
			}
			if err != nil {
				return err
			}
		}
		for _, art := range artists {
			err = a.repo.InsertArtist(art)
			if err != nil {
				return err
			}
		}
	}

	if albumsData, ok := data[albumTable].([]interface{}); ok {
		valid = true
		for _, albumData := range albumsData {
			albumJSON, _ := json.Marshal(albumData)
			var albumJson music.AlbumJSON
			var album music.Album
			if err := json.Unmarshal(albumJSON, &albumJson); err == nil {
				err = album.SetFromJSON(albumJson)
				albums = append(albums, album)
			}
			if err != nil {
				return err
			}
		}
		for _, al := range albums {
			err = a.repo.InsertAlbum(al)
			if err != nil {
				return err
			}
		}
	}

	if Data, ok := data[genresTable].([]interface{}); ok {
		valid = true
		for _, genreData := range Data {
			JSON, _ := json.Marshal(genreData)
			var genreJson music.GenreJSON
			var genre music.Genre
			if err := json.Unmarshal(JSON, &genreJson); err == nil {
				err = genre.SetFromJSON(genreJson)
				genres = append(genres, genre)
			}
			if err != nil {
				return err
			}
		}
		for _, g := range genres {
			err = a.repo.InsertGenre(g)
			if err != nil {
				return err
			}
		}
	}

	if Data, ok := data[playlistTable].([]interface{}); ok {
		valid = true
		for _, playlistData := range Data {
			JSON, _ := json.Marshal(playlistData)
			var playlistJson music.PlaylistJSON
			var playlist music.Playlist
			if err := json.Unmarshal(JSON, &playlistJson); err == nil {
				err = playlist.SetFromJSON(playlistJson)
				playlists = append(playlists, playlist)
			}
			if err != nil {
				return err
			}
		}
		for _, p := range playlists {
			err = a.repo.InsertPlaylist(p)
			if err != nil {
				return err
			}
		}
	}

	if Data, ok := data[trackTable].([]interface{}); ok {
		valid = true
		for _, trackData := range Data {
			JSON, _ := json.Marshal(trackData)
			var trackJson music.TrackJSON
			var track music.Track
			if err := json.Unmarshal(JSON, &trackJson); err == nil {
				err = track.SetFromJSON(trackJson)
				tracks = append(tracks, track)
			}
			if err != nil {
				return err
			}
		}
		for _, t := range tracks {
			err = a.repo.InsertTrack(t)
			if err != nil {
				return err
			}
		}
	}

	if Data, ok := data[adminHistory].([]interface{}); ok {
		valid = true
		for _, record := range Data {
			JSON, _ := json.Marshal(record)
			var historyJson music.HistoryJSON
			var hist music.History
			if err := json.Unmarshal(JSON, &historyJson); err == nil {
				err = hist.SetFromJSON(historyJson)
				history = append(history, hist)
			}
			if err != nil {
				return err
			}
		}
		for _, ah := range history {
			err = a.repo.InsertHistory(ah)
			if err != nil {
				return err
			}
		}
	}

	if Data, ok := data[artistGenreTable].([]interface{}); ok {
		valid = true
		for _, record := range Data {
			JSON, _ := json.Marshal(record)
			var agJson music.ArtistGenreJSON
			var ag music.ArtistGenre
			if err := json.Unmarshal(JSON, &agJson); err == nil {
				err = ag.SetFromJSON(agJson)
				ags = append(ags, ag)
			}
			if err != nil {
				return err
			}
		}
		for _, ag := range ags {
			err = a.repo.InsertArtistGenre(ag)
			if err != nil {
				return err
			}
		}
	}

	if Data, ok := data[artistAlbumTable].([]interface{}); ok {
		valid = true
		for _, record := range Data {
			JSON, _ := json.Marshal(record)
			var aaJson music.ArtistAlbumJSON
			var aa music.ArtistAlbum
			if err = json.Unmarshal(JSON, &aaJson); err == nil {
				err = aa.SetFromJSON(aaJson)
				aas = append(aas, aa)
			}
			if err != nil {
				return err
			}
		}
		for _, aa := range aas {
			err = a.repo.InsertArtistAlbum(aa)
			if err != nil {
				return err
			}
		}
	}

	if Data, ok := data[userArtistTable].([]interface{}); ok {
		valid = true
		for _, record := range Data {
			JSON, _ := json.Marshal(record)
			var uaJson music.UserArtistJSON
			var ua music.UserArtist
			if err = json.Unmarshal(JSON, &uaJson); err == nil {
				err = ua.SetFromJSON(uaJson)
				uas = append(uas, ua)
			}
			if err != nil {
				return err
			}
		}
		for _, ua := range uas {
			err = a.repo.InsertUserArtist(ua)
			if err != nil {
				return err
			}
		}
	}

	if Data, ok := data[userAlbumTable].([]interface{}); ok {
		valid = true
		for _, record := range Data {
			JSON, _ := json.Marshal(record)
			var ualJson music.UserAlbumJSON
			var ual music.UserAlbum
			if err = json.Unmarshal(JSON, &ualJson); err == nil {
				err = ual.SetFromJSON(ualJson)
				uals = append(uals, ual)
			}
			if err != nil {
				return err
			}
		}
		for _, ual := range uals {
			err = a.repo.InsertUserAlbum(ual)
			if err != nil {
				return err
			}
		}
	}

	if Data, ok := data[playlistTrackTable].([]interface{}); ok {
		valid = true
		for _, record := range Data {
			JSON, _ := json.Marshal(record)
			var ptJson music.PlaylistTrackJSON
			var pt music.PlaylistTrack
			if err = json.Unmarshal(JSON, &ptJson); err == nil {
				err = pt.SetFromJSON(ptJson)
				pts = append(pts, pt)
			}
			if err != nil {
				return err
			}
		}
		for _, pt := range pts {
			err = a.repo.InsertPlaylistTrack(pt)
			if err != nil {
				return err
			}
		}
	}
	if Data, ok := data[userGenreTable].([]interface{}); ok {
		valid = true
		for _, record := range Data {
			JSON, _ := json.Marshal(record)
			var ugJson music.UserGenreJSON
			var ug music.UserGenre
			if err = json.Unmarshal(JSON, &ugJson); err == nil {
				err = ug.SetFromJSON(ugJson)
				ugs = append(ugs, ug)
			}
			if err != nil {
				return err
			}
		}
		for _, ug := range ugs {
			err = a.repo.InsertUserGenre(ug)
			if err != nil {
				return err
			}
		}
	}
	if Data, ok := data[userPlaylistTable].([]interface{}); ok {
		valid = true
		for _, record := range Data {
			JSON, _ := json.Marshal(record)
			var upJson music.UserPlaylistJSON
			var up music.UserPlaylist
			if err = json.Unmarshal(JSON, &upJson); err == nil {
				err = up.SetFromJSON(upJson)
				ups = append(ups, up)
			}
			if err != nil {
				return err
			}
		}
		for _, up := range ups {
			err = a.repo.InsertUserPlaylist(up)
			if err != nil {
				return err
			}
		}
	}

	if !valid {
		return errors.New("invalid json file")
	}

	return err
}

func (a *AdminServiceImpl) GetGenrePopularity() ([]music.Genre, error) {
	genres, err := a.repo.SelectGenrePopularity()
	if err != nil {
		return nil, err
	}
	var maxPopul int64 = genres[0].Popularity
	for k, _ := range genres {
		if maxPopul <= genres[k].Popularity {
			maxPopul = genres[k].Popularity
		}
	}
	for k, _ := range genres {
		p := (float64(genres[k].Popularity) / float64(maxPopul)) * float64(100)
		genres[k].Popularity = int64(p)
	}
	return genres, err
}

func (a *AdminServiceImpl) GetGenreDiversity() ([]music.Genre, error) {
	genres, err := a.repo.SelectGenreDiversity()
	if err != nil {
		return nil, err
	}
	var sum float64 = 0
	for k, _ := range genres {
		sum += genres[k].Diversity
	}
	for k, _ := range genres {
		p := (genres[k].Diversity / sum) * float64(100)
		genres[k].Diversity = p
	}
	return genres, err
}

func (a *AdminServiceImpl) GetArtistPopularity() ([]music.Artist, error) {
	artists, err := a.repo.SelectArtistPopularity()
	if err != nil {
		return nil, err
	}
	return artists, err

}
