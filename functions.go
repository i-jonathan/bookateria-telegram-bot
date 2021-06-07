package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
)

func add(list []query, element query) []query {
	for index, query := range list {
		if query.UserID == element.UserID {
			if query.ChatID == element.ChatID {
				list[index].Type = element.Type
				list[index].MessageID = element.MessageID
				return list
			}
		}
	}
	return append(list, element)
}

func get(list []query, chatID int) (query, bool) {
	for _, query := range list {
		if query.ChatID == chatID {
			return query, true
		}
	}
	return query{}, false
}

func remove(list []query, query query) []query {
	for index, user := range list {
		if user.UserID == query.UserID {
			return append(list[:index], list[index+1:]...)
		}
	}
	return list
}


func getToken(chatID int) (string, error) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB: 0,
	})

	token, err := client.Get(ctx, strconv.Itoa(chatID)).Result()

	if err != nil {
		return "", err
	}

	return token, nil
}