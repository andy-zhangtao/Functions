# Weaviate Plugin

```json
{
    "plugin_key": 1,
    "name": "weaviate",
    "descript": "This is a Weaviate plugin",
    "module": "",
    "input": [
        {
            "Name": "content",
            "Value": {
                "type": "string",
                "description": "The content of need store into weaviate"
            }
        },
        {
            "Name": "username",
            "Value": {
                "type": "string",
                "description": "The user name"
            }
        },
        {
            "Name": "date",
            "Value": {
                "type": "string",
                "description": "The date of the content. YYYY-MM-DD format"
            }
        }
    ],
    "reference": {
        "up":-1,
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