# Meli URL Shortener

### How to Up

```bash
cp .env.example .env

docker compose up -d

bash docs/dynamo-db/init.sh
```

### Documentation

After run, PUT one URL into http://localhost:8080

```bash
curl --request POST \
  --url http://localhost:8080/ \
  --header 'Content-Type: multipart/form-data' \
  --form url=https://www.mercadolivre.com.br/
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
```


### Pending Features

 - [ ] Short URL Redirect
 - [ ] Short URL Delete
 - [ ] Short URL Stats
 - [ ] Deployment
 - [ ] API Documentation
 - [ ] Swagger UI
 - [ ] Tests

### Issues

 - [ ] Worker Container not working because AWS Credentials cache problems
 - [ ] Add tests
 - [ ] Add README.md

