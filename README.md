#### Introduction
This simple application introduces the simplest usages of a couple of common libraries.
The application registers hosts and provide get-all and get-by-id operations via rest.
There're 2 ways to register a host into the application: environment variables and openshift's routes

To register a host via environment variables just add variable with preconfigured prefix.
To register a host via openshift's routes:
* If you start application outside openshift cluster just update config/kube-config.yaml or provide your config-file via ENV_CONFIG_PATH
* If you start application inside openshift cluster you have to provide a namespace (mapping.namespace configuration value)

**Note**: all default configuration variables placed in config/host-manager.yaml (could be changed via ENV_CONFIG_PATH)

#### Docker
There're 2 docker-files (see build directory):
* docker1 - single-stage docker-file
* docker2 - multi-stage docker-file

To build an image just: `docker build -t host-manager -f build/dockerX/Dockerfile .`

#### OpenShift
To push the image to internal openshift's repository:
* Provide to the user required grants: `oc policy add-role-to-user registry-viewer admin` and `oc policy add-role-to-user registry-editor admin`
* Get current internal registry's ip address: `oc get svc -n default docker-registry`
* Tag the image with current URL: `docker tag host-manager:latest <REGISTRY-IP>:5000/openshift/host-manager:<VERSION>`
* Login to internal registry: `docker login -u admin -p $(oc whoami -t) <REGISTRY-IP>:5000`
* Push the image: `docker push <REGISTRY-IP>:5000/openshift/host-manager:<VERSION>`

To create new application: 
* Create new deployment configuration: `oc new-app --docker-image <REGISTRY-IP>:5000/openshift/host-manager:<VERSION> --name host-manager`
* Add service account to deployment configuration which is has grants to read routes in the namespace:
  * `oc create sa host-manager`
  * `oc adm policy add-role-to-user view -z host-manager`
  * `oc patch dc/host-manager -p '{"spec":{"template":{"spec":{"serviceAccount":"host-manager"}}}}'`
* Create new service: `oc expose dc/host-manager --port=8080`
* Create new route: `oc expose svc/host-manager`
* To add a host via environment variable: `oc set env dc/host-manager APP_HOST_HOST1=host1:8080`
* Add namespace scanning if required: `oc set env dc/host-manager APP_MAPPING_NAMESPACE=<NAMESPACE>`

#### Tracing
For development purposes it's easy to use all-in-one jaeger version: 
* Start all-in-one jaeger: `docker run --name jgr -p5775:5775/udp -p6831:6831/udp -p16686:16686 -d jaegertracing/all-in-one:latest`
* Go to: `http://localhost:16686/search`

#### Swagger
[Swaggo/swag][swaggo] is used to generate swagger.

To generate swagger:
* `go get -u github.com/swaggo/swag/cmd/swag`
* `<GO_PATH>/bin/swag init --dir ./ --generalInfo ./cmd/host-manager/main.go --output ./api/swagger`
* `http://<HOST>:<PORT>/swagger/index.html`

#### Testing
* `go test -coverprofile=out/test.html ./internal/service/mapping/`
* `go tool cover -html=out/test.html`

#### Configuration
There's 2 options to configure application: [Envconfig][envconfig] and configuration file (config/host-manager.yaml by default) 

#### PProf
* `https://github.com/gin-contrib/pprof`
* `go tool pprof http://localhost:8080/debug/pprof/profile?seconds=5`

[envconfig]: https://github.com/kelseyhightower/envconfig
[swaggo]: https://github.com/swaggo/swag