# URL Shortener

### How to Up

```bash
cp .env.example .env

docker compose up -d

bash docs/dynamo-db/init.sh
```

### Documentation

After run, POST your complete URL into http://localhost:8080

```bash
curl --request POST \
  --url http://localhost:8080/ \
  --header 'Content-Type: multipart/form-data' \
  --form url=https://www.mercadolivre.com.br/

curl --request GET \
  --url http://localhost:8080/{short}/stats

curl --request DELETE \
  --url http://localhost:8080/{short}
```

### Hexagonal Diagram

```mermaid
classDiagram
    class `ports/Input/HttpServer` {
        <<interface>>
        +Start()
        +Shutdown()
        +Options() Options
    }

    class `ports/Input/Consumer` {
        <<interface>>
        +Start()
        +GetRouter() Router
    }

    class `ports/Output/EventDispatcher` {
        <<interface>>
        +Dispatch(context, to, msg) error
        +Close()
    }

    class `ports/Output/UrlRepository` {
        <<interface>>
        +UrlFromShort(context, short) url
        +StatsFromShort(context, short) map
    }

    class App {
    }
    
    class FastHttp {
    }
    
    class Kafka {
    }
    
    class DynamoDb {
    }

    App --> `ports/Input/HttpServer`
    App --> `ports/Input/Consumer`
    App --> `ports/Output/EventDispatcher`
    FastHttp <|-- `ports/Input/HttpServer`
    Kafka <|-- `ports/Input/Consumer`
    Kafka <|-- `ports/Output/EventDispatcher`
    DynamoDb <|-- `ports/Output/EventDispatcher`
    DynamoDb <|-- `ports/Output/UrlRepository`
```


### Pending Features

 - [ ] Deployment
 - [ ] Complete Documentation
 - [ ] Swagger UI
 - [ ] Tests

### Issues

 - [ ] Worker Container not working because AWS Credentials cache problems
