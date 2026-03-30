package feed

func GetAbuAliExpress() ([]string, error) {
	messages, err := fetchTelegramChannelMessages("abualiexpress", 10)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
