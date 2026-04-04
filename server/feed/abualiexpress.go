package feed

import "context"

const abuTelegramChannel = "@abualiexpress"

func GetAbuAliExpress(ctx context.Context) ([]string, error) {
	messages, err := fetchTelegramChannelMessages(ctx, abuTelegramChannel, 10)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
