package feed

import "context"

const ynetTelegramChannel = "@ynetalerts"

func GetYnet(ctx context.Context) ([]string, error) {
	messages, err := fetchTelegramChannelMessages(ctx, ynetTelegramChannel, 10)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
