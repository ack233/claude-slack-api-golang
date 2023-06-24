# claude slack api golang
通过 Slack API 来使用 Claude


## 配置

### 1.创建一个 slack 应用 

1. 访问 https://api.slack.com/apps ` > ` 点击 `From an app manifest` ` > ` 选择一个工作区 ` > ` 下一步  

2. 转到  `OAuth & Permissions` > `Scopes`, 并添加以下范围：:

Bot Token Scopes:  
| OAuth Scope | Reasons |
| ------------| ------------|
| `channels:history` | Read conversation history and find claude's conversations. |
| `users:read` | Find out the user id of claude and talk to it. |


User Token Scopes:  
| OAuth Scope | Reasons |
| ------------| ------------|
| `chat:write` | Send messages on a user’s behalf. |

3. 转到 `OAuth & Permissions` > `Redirect URLs`, 然后添加回调URL: `http://your_server_ip:5000/callback`  
对应后续配置文件的SLACK_OAUTH_REDIRECT_URI字段.

 `SLACK_BOT_TOKEN` 位置在 `Bot User OAuth Token`.

4. 转到 `Event Subscriptions` > `Enable Events` > `Subscribe to bot events`, 添加事件:

| Event Name | Reasons |
| ------------| ------------|
| `message.channels` | Receive claude's reply. |

5. 转到 `Socket Mode` > `Enable Socket Mode`.

6. 转到 `Baisc Information` > `App-Level Tokens` > `Generate Tokens and Scopes`: 

Token Name: `SocketMode`
Scope: `connection:write`

这里会的到 `SLACK_APP_TOKEN` 配置.

### 2. 设置你的配置文件 config.yaml

| Key                        | Description                                              | Example                           |
| -------------------------- | ---------------------------------------------------------|-----------------------------------|
| port                       | Server listen port                                       | 5000                              |                   |
| ENCRYPTION_KEY             | Key used to encrypt user's access_token                  | a_random_string                         |
| SLACK_APP_TOKEN            | App token for Slack API                                  | xapp-1-A05xxxxxLGN8-517xxxxxxxxxx-b853xxxxxxxxxxxxxxxd34850xxxxxxxxxxbed3084|
| SLACK_BOT_TOKEN            | Bot token for Slack API                                  | xoxb-517xxxxxxxxxx-51xxxxxxxx-xxxxxxxxxxxxxxxxxxxxxxx|
| SLACK_OAUTH_CLIENT_ID      | Client ID for Slack API OAuth 2.0 authentication flow    | 1234567890                        |
| SLACK_OAUTH_CLIENT_SECRET  | Client secret for Slack API OAuth 2.0 authentication flow| AbCdEfGhIjKlMnOpQrStUvWxYz123456 |
| SLACK_OAUTH_REDIRECT_URI   | Redirect URI for Slack API OAuth 2.0 authentication flow | http://your_server_ip:5000/callback   |

### 3.  启动项目并获取access_token
1.执行` go run main.go` 启动项目

2.访问http://your_server_ip:5000/login，将您的应用程序授权给工作区。用户将获得一个加密的ACCESS_TOKEN。

3.邀请您的应用和 Claude 加入工作区的某个频道, 获取频道id (点击左上角 #xx 字样, 弹窗底部就是频道id)

### 4. 发送请求

```
import json
import uuid
import requests

def interact_with_server(channel_id, access_token, prompt, conversation_id=None):
    payload = {
        "action": "next",
        "messages": [
            {
                "id": str(uuid.uuid4()),
                "role": "user",
                "author": {
                    "role": "user"
                },
                "content": {
                    "content_type": "text",
                    "parts": [
                        prompt
                    ]
                }
            }
        ],
        "conversation_id": conversation_id,
        "parent_message_id": str(uuid.uuid4()),
        "channel_id": channel_id,
        "model": "claude-unknown-version"
    }

    headers = {
        'Authorization': f'Bearer {access_token}',
        'Content-Type': 'application/json'
    }

    response = requests.post("https://ip:port/backend-api/conversation", headers=headers, json=payload, timeout=1000)
    response.raise_for_status()

    for line in response.iter_lines():
        if not line or line is None:
            continue
        if "data:" in str(line):
            line = line[5:]
        if "[DONE]" in str(line):
            break

        try:
            line = json.loads(line)
        except json.decoder.JSONDecodeError as e:
            print(line)
            print(e)
            continue

        conversation_id = line["conversation_id"]
        message = line["message"]["content"]["parts"][0]
        yield (conversation_id, message)

# Example usage
channel_id = 'xxxxxxxx' 
access_token = 'xxxxxxxxxxxx......xxxxxxxxx' # 
conversation_id = None

#  call
for conversation_id, message in interact_with_server(channel_id, access_token, "晚上好?", conversation_id):
    print(f"Received message: {message}")
```

##  加telegram claude 群体验
https://t.me/claude00000



<br>
<br>

## 参考:
- https://github.com/LlmKira/claude-in-slack-server





