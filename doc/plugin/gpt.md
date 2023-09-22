# GPT Plugin

```json
{
    "id": 1,
    "name": "plugin1",
    "descript": "This is a GPT plugin",
    "module": "gpt",
    "input": [
        {
            "Name": "prompt",
            "Value": {
                "Description": "What is the meaning of life?"
            }
        },
        {
            "Name": "max_tokens",
            "Value": {
                "Description": "50"
            }
        },
        {
            "Name": "temperature",
            "Value": {
                "Description": "0.7"
            }
        },
        {
            "Name": "model",
            "Value": {
                "Description": "davinci"
            }
        }
    ],
    "reference": {
        "up": "",
        "down": []
    },
    "invoke_type": "sync",
    "invoke_url": "http://localhost:8080/invoke",
    "stage": {
        "type": "text",
        "value": "hello world"
    }
}
```