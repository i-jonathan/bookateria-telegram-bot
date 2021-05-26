package main

func add(list []query, element query) []query {
	for index, query := range list {
		if query.User == element.User {
			if query.ChatID == element.ChatID {
				list[index].Type = element.Type
				list[index].MessageID = element.MessageID
				return list
			}
		}
	}
	return append(list, element)
}

func get(list []query, chat_id int) (query, bool) {
	for _, query := range list {
		if query.ChatID == chat_id {
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
