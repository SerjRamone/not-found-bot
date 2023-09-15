// Package app core bot code
package app

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/SerjRamone/not-found-bot/config"
	"github.com/SerjRamone/not-found-bot/internal/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vitali-fedulov/images4"
	"go.uber.org/zap"
)

// Run ...
func Run(cfg *config.Config) {
	l := logger.Get()

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		l.Fatal("can't create bot instance", zap.Error(err))
	}

	l.Info("autorized on account", zap.String("bot_name", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		// checks only channel posts and ignore others
		if update.ChannelPost != nil {
			if isTargetChan(update.ChannelPost.Chat.ID, cfg.TargetChannels) {
				if len(update.ChannelPost.Photo) > 0 {
					l.Info(
						"new image from channel",
						zap.String("channel_name", update.ChannelPost.Chat.UserName),
					)

					fileID := update.ChannelPost.Photo[0].FileID

					// Get file info via FileID
					fileInfo, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
					if err != nil {
						l.Error("can't get file", zap.Error(err))
					}

					// Получите URL для скачивания файла
					fileURL := fileInfo.Link(bot.Token)

					// Определите путь и имя файла, куда сохранить изображение
					savePath := filepath.Join("/tmp", fileID+".jpg")
					err = os.MkdirAll(filepath.Dir(savePath), os.ModePerm)
					if err != nil {
						l.Error("can't create dirrectories", zap.Error(err))
					}

					// Скачайте и сохраните файл
					response, err := http.Get(fileURL)
					if err != nil {
						l.Error("http.Get() error", zap.Error(err))
					}

					file, err := os.Create(savePath)
					if err != nil {
						l.Error("can't create file", zap.Error(err))
					}

					_, err = io.Copy(file, response.Body)
					if err != nil {
						l.Error("can't copy file", zap.Error(err))
					}
					response.Body.Close()
					file.Close()

					// Photos to compare.
					// Open files (discarding errors here).
					img1, _ := images4.Open(cfg.PathToTargetImage)
					img2, _ := images4.Open(savePath)

					// Icons are compact image representations (image "hashes"). Name "hash" is reserved for "true" hashes in package imagehash.
					icon1 := images4.Icon(img1)
					icon2 := images4.Icon(img2)

					// Comparison. Images are not used directly. Icons are used instead, because they have tiny
					// memory footprint and fast to compare. If you need to include images rotated right and left use func Similar90270.
					if images4.Similar(icon1, icon2) {
						l.Info("🔴  images are similar")

						err = deletePost(
							bot,
							update.ChannelPost.Chat.ID,
							update.ChannelPost.MessageID,
						)
						if err != nil {
							l.Error("deleting message error", zap.Error(err))
						} else {
							l.Info("🔥  post deleted")
						}
					} else {
						l.Info("🟢  images are distinct")
					}

					// remove downloaded image
					err = os.Remove(savePath)
					if err != nil {
						l.Error("deleting file error", zap.Error(err))
					}

				}
			}
		}
	}
}

// deletePost from channel via messageID
func deletePost(bot *tgbotapi.BotAPI, chatID int64, messageID int) error {
	deleteMessageConfig := tgbotapi.DeleteMessageConfig{
		ChatID:    chatID,
		MessageID: messageID,
	}

	_, err := bot.Request(deleteMessageConfig)
	if err != nil {
		return err
	}
	return nil
}

// isTargetChan return true if targetChannels contents c
func isTargetChan(c int64, targetChannels []int64) bool {
	for _, v := range targetChannels {
		if c == v {
			return true
		}
	}
	return false
}
