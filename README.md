## deepseen-backend

A backend for the [Deepseen](https://github.com/peterdee/deepseen-desktop) project

Stack: [Golang](https://golang.org), [Fiber](https://gofiber.io), [MongoDB](https://github.com/mongodb/mongo-go-driver), [Redis](https://github.com/go-redis/redis), [JWT](https://github.com/dgrijalva/jwt-go)

DEV: http://localhost:1337

STAGE: https://deepseen-backend.herokuapp.com

### Deploy

Golang **v1.15.X** is required

```
git clone https://github.com/peterdee/deepseen-backend
cd ./deepseen-backend
```

### Environment variables

The `.env` file is required, see [.env.example](.env.example) for details

### Launch

```
go run ./
```

Can be launched with [AIR](https://github.com/cosmtrek/air), see [run.sh](run.sh) for details

### Heroku

The `staging` branch is auto-deployed to Heroku

### License
[MIT](LICENSE)
