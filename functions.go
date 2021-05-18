package main

func add(list []query, element query) []query {
	for index, query := range list {
		if query.User == element.User {
			if query.Chat_ID == element.Chat_ID {
				list[index].Type = element.Type
				list[index].Message_ID = element.Message_ID
				return list
			}
		}
	}
	return append(list, element)
}

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
