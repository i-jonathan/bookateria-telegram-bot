package main

func get(list []query, chat_id int) (query, bool) {
	for _, query := range list {
		if query.Chat_ID == chat_id {
			return query, true
		}
	}
	return query{}, false
}

func remove(list []query, query query) []query {
	for index, user := range list {
		if user.User == query.User {
			return append(list[:index], list[index+1:]...)
		}
	}
	return list
}
