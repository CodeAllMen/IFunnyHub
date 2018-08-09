package page

import (
	"github.com/MobileCPX/IFunnyHub/models/content"
)

func VideoContent(id string) map[string]interface{} {
	video := content.GetItemlById(id)
	return map[string]interface{}{
		"Src":        video.Source,
		"Img":        video.Img,
		"Title":      video.Title,
		"Time":       video.Create,
		"Like":       video.Like,
		"Dislike":    video.Dislike,
		"YouMaylike": content.GetRandomItem(1, id),
	}
}

func GameContent(id string) map[string]interface{} {
	game := content.GetItemlById(id)
	return map[string]interface{}{
		"Src":        game.Source,
		"Title":      game.Title,
		"Time":       game.Create,
		"Like":       game.Like,
		"Dislike":    game.Dislike,
		"YouMaylike": content.GetRandomItem(2, id),
	}
}
func PictureContent(id string) map[string]interface{} {
	picture := content.GetItemlById(id)
	return map[string]interface{}{
		"Id":         picture.Id,
		"Src":        picture.Source,
		"Title":      picture.Title,
		"Time":       picture.Create,
		"Like":       picture.Like,
		"Dislike":    picture.Dislike,
		"YouMaylike": content.GetRandomItem(3, id),
	}
}
func RingtoneContent(id string) map[string]interface{} {
	video := content.GetItemlById(id)
	return map[string]interface{}{
		"Src":        video.Source,
		"Title":      video.Title,
		"Time":       video.Create,
		"Like":       video.Like,
		"Dislike":    video.Dislike,
		"YouMaylike": content.GetRandomItem(4, id),
	}
}
