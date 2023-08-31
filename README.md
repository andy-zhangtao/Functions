# Functions
Some serverless functions

## Data struct

+ Weaviate 
```json
{
    "class": "diary",
    "description": "the work content for user",
    "vectorizer": "text2vec-openai",
    "properties": [
        {
            "dataType": [
                "text"
            ],
            "description": "content that will be vectorized",
            "name": "content"
        },
        {
            "dataType": [
                "text"
            ],
            "description": "the user name",
            "name": "user"
        },
        {
            "dataType": [
                "text"
            ],
            "description": "the date of the diary",
            "name": "date"
        },
    ],
    "vectorIndexConfig":{
        "ef": 100
    }
}
```