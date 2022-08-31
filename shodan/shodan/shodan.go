package shodan

// 定义基本的信息
const BASEURL = "https://api.shodan.io"

type Client struct {
	apiKey string
}

func New(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
	}
}
