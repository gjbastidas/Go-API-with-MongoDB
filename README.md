# Simple API with MongoDB backend

Hopefully, this app helps to understand how to create an API
with **MongoDB** backend using **Golang**.

I also included topics such as:
- Anonymous functions, Go routines, Channels, Interfaces and Generics
- Graceful shutdown mechanism
- Dependency injection
- Unit testing

## Use case

With this API an **Author** can submit a message in a **Post**. 
An author can also make **Comments** on that or any other post.

## Code organization

Main folders are:
- [app](./app/): containing the API logic
- [models](./models/): containing the DB models

## Data Schema

For this purpose I need to model `one-to-one` relationships between:
- Post and Author 
- Comment and Author

And `one-to-many` relationship between:
- Post and Comments

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
- GNU Make =3.81
- Docker =20.10.20
- MongoDB docker image =6.0.3
- Go =1.19

## How to play with this

There are **default values** for environment variables set 
in the [env.go](./env/env.go) file. 

Run the app:
```shell
make go-run
```

This command should trigger the local pipeline to:
1. Run linters, tests and build the code
2. Run the MongoDB docker container
3. TODO Run the API docker container

## Delete everything

TODO

## Last comments

If you need further environment customization 
you'll need to first export them in your terminal:

```shell
export SVR_ADDR="replace with custom server address and port"
export DB_USERNAME="replace with MongoDB username"
export DB_PASSWORD="replace with MongoDB user password"
...
```

Then remove the respective comment in the `docker-run` section 
in the [Makefile](./Makefile) in order to export them to the docker container

## Have fun!
... and feel free to add your comments and/or feedback
