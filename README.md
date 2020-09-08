#### Swagger
[Swaggo/swag][swaggo] is used to generate swagger.

To generate swagger:
* `go get -u github.com/swaggo/swag/cmd/swag`
* `<GO_PATH>/bin/swag init --dir ./ --generalInfo ./cmd/host-manager/main.go --output ./api/swagger`
* `http://<HOST>:<PORT>/swagger/index.html`

#### Configuration
[Envconfig][envconfig] is used to read application configuration

[envconfig]: https://github.com/kelseyhightower/envconfig
[swaggo]: https://github.com/swaggo/swag