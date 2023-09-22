# How to invoke?

## How to invoke workflow?

```curl
curl --location 'https://xxxx/api/workflow?id=12345' \
--header 'Content-Type: application/json' \
--data '{
    "action": 1,
    "user": "zhangtao",
    "name": "",
    "question": "请记录今天的工作内容: 我完成了Father的初步设计和调试工作。"
}'
```

> `id` is the workflow id. 

> The `name` is a optional value, needn't set it.

> user is the user's name, it will be store in weaviate. 