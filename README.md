# Simple API with MongoDB backend

Hopefully, this app helps to understand how to create an API
with **MongoDB** backend using **Golang**.

I also included topics such as:
- Anonymous functions, Go routines, Channels, Interfaces and Generics.
- Graceful shutdown mechanism.
- Dependency injection.
- Unit testing.

## Use case

With this API an **Author** can submit a message in a **Post**. An author 
can also make **Comments** on that or any other post.

## Code organization

Main folders are:
- [app](./app/): containing the API logic.
- [models](./models/): containing the DB models.

## Data Schema

For this purpose I need to model `one-to-one` relationships between:
- Post and Author.
- Comment and Author.

And `one-to-many` relationship between:
- Post and Comments.

So basically, I'm using 2 entities: `Post` and `Comment` 
with an *Author*'s name in each of them. Whereas a 
Post consists of the following fields:
- `Id`
- `Content`
- `Author`

and a Comment consists of the following fields:
- `Id`
- `Content`
- `Author`
- `PostId` (Using reverse reference to avoid limitation
of a big list of Ids in the Post document)

## Prerequisites

This app has been tested with:
- Git =2.37
- GNU Make =3.81
- Docker =20.10.20
- MongoDB docker image =6.0.3
- Go =1.19

## How to play with this

### 1. Run the App
Clone this repo to your local machine:
```
git clone https://github.com/gjbastidas/GoSimpleAPIWithMongoDB.git
```

Change to the GoSimpleAPIWithMongoDB directory:
```
cd GoSimpleAPIWithMongoDB
```

Export the following environment variables:
```shell
export DB_USERNAME="replace with MongoDB username" #i.e admin
export DB_PASSWORD="replace with MongoDB user password"  #i.e secret
export DB_HOST="replace with MongoDB hostname"  #i.e db
export DB_PORT="replace with MongoDB port"  #i.e 27017
```

Run the app:
```shell
make app-run
```

This last command should trigger the local pipeline to:
1. Run linters and tests
2. Set and run the MongoDB docker container
3. Build the API docker image and run the API container

### 2. Interact with the API
The simplest way is by issuing CURL commands from your terminal:

Create a post
```shell
curl -X POST http://localhost:8088/post/ \
  -H 'Content-Type: application/json' \
  -d '{"content": "my first post","author": "some author"}'
```

Get a post
```shell
curl http://localhost:8088/post/<<replace with id>>
```

Update a post
```shell
curl -X PUT http://localhost:8088/post/<<replace with id>> \
  -H 'Content-Type: application/json' \
  -d '{"content": "updated post","author": "some author"}'
```

Delete a post
```shell
curl -X DELETE http://localhost:8088/post/<<replace with id>>
```

If you're using **VS Code**, with the [Rest Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) integration,
I already included a script you could use [here](./scripts/check.http)

TODO Swagger

## Delete everything

Run:
```
make app-delete
```

## Have fun!
... and feel free to add your comments or feedback, open PRs, etc.
