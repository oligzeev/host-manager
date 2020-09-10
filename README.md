#### Docker
* `docker build -t host-manager -f build/docker/Dockerfile .`
* `docker tag host-manager:latest 172.30.171.98:5000/openshift/host-manager:1`

#### OpenShift
* `oc policy add-role-to-user registry-viewer admin`
* `oc policy add-role-to-user registry-editor admin`
* `oc get svc -n default docker-registry`
* `docker tag host-manager:latest 172.30.171.98:5000/openshift/host-manager:1`
* `docker login -u admin -p $(oc whoami -t) 172.30.171.98:5000`
* `docker push 172.30.171.98:5000/openshift/host-manager:1`
* `oc new-app --docker-image 172.30.171.98:5000/openshift/host-manager:1 --name host-manager`
* `oc expose dc/host-manager --port=8080`

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