package processors

import (
	"bitbucket.org/daforester/go-audience/apps/priceisright/app/sockets"
	"bitbucket.org/daforester/go-audience/apps/priceisright/models"
	"bitbucket.org/daforester/go-audience/apps/priceisright/repositories"
	"encoding/json"
	"errors"
	"mime/multipart"
	"os"
	"strings"
)

type SaveImageUploadsProcessorInterface interface {
	SaveImageUploads(rounds []*models.Round, files []*multipart.FileHeader, path string) error
}

type SaveImageUploadsProcessor struct {
	gameRepository repositories.GameRepository
}

func (P SaveImageUploadsProcessor) New(gameRepository repositories.GameRepository) *SaveImageUploadsProcessor {
	return &SaveImageUploadsProcessor{
		gameRepository,
	}
}

func (P *SaveImageUploadsProcessor) SaveImageUploads(rounds []*models.Round, files []*multipart.FileHeader, path string) error {
	var p string
	if len(rounds) > len(files) {
		return errors.New("mismatching array lengths")
	}
	if path[len(path)-1] == '/' {
		p = path
	} else {
		p = path + "/"
	}
	for i, round := range rounds {
		var ext = ""

		if files[i] == nil {
			continue
		}

		switch strings.ToLower(files[i].Header.Get("Content-Type")) {
		case "image/png":
			ext = "png"
		case "image/webp":
			ext = "webp"
		default:
			continue
		}

		fo, _ := files[i].Open()
		b := make([]byte, files[i].Size)
		_, _ = fo.Read(b)

		err := os.WriteFile(p+round.UUID+"."+ext, b, 0644)
		if err != nil {
			return err
		}

		round.Image = round.UUID + "." + ext
	}

	return nil
}

func (P *SaveImageUploadsProcessor) roundsUpdatedMessage(rounds []*models.Round) error {
	wsCollection := sockets.NewCollection("priceIsRight")

	data := struct {
		Rounds []*models.Round
	}{
		rounds,
	}

	msg, err := json.Marshal(struct {
		Command string
		Data    interface{}
	}{
		"ROUNDS_UPDATED",
		data,
	})

	if err != nil {
		return errors.New("unable to marshal socket message")
	}

	socks := wsCollection.HasTag(string(models.PLAYER_CLASS_LEADER))
	wsCollection.SendText(string(msg), socks...)

	socks = wsCollection.HasTag(string(models.PLAYER_CLASS_DISPLAY))
	wsCollection.SendText(string(msg), socks...)

	return nil
}
