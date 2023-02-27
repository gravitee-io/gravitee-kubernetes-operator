## [0.5.1](https://github.com/gravitee-io/gravitee-kubernetes-operator/compare/0.5.0...0.5.1) (2023-02-27)


### Bug Fixes

* restore and deprecate v0.4.0 status fields ([536e806](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/536e806b567cf6c36d7ed64cf615f12a6804cecc))

# [0.5.0](https://github.com/gravitee-io/gravitee-kubernetes-operator/compare/0.4.0...0.5.0) (2023-02-17)


### Bug Fixes

* remove automatic plan creation ([98b78a8](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/98b78a85b4b426c473a7d3504b6283311dc5d761))

### Features

* add ingress events on create, update and delete ([1df9534](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/1df9534b1c408e9491e9d8815f17b6ca16ffafea))
* allow to listen for resources in a namespace ([5cbdf0d](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/5cbdf0d824f36226c847a3876d6e28518baf03bb))
* allow users to customize the manager env ([37c8644](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/37c8644bed16fad2fc1552cc7f9e4b141da45fe3))
* handle ingress with multiple hosts ([e56b5ac](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/e56b5acb3d8b2078cc104f438e871b0cacdb065b))
* release the operator as a helm chart ([b182920](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/b182920ff8ffab87002f96236f933f64c7ed7b23))
* set definition context on create and updates ([520f710](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/520f710502c5cfc14984a51e675c960df4deb8da))

# [0.5.0](https://github.com/gravitee-io/gravitee-kubernetes-operator/compare/0.4.0...0.5.0) (2023-02-15)


### Bug Fixes

* remove automatic plan creation ([98b78a8](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/98b78a85b4b426c473a7d3504b6283311dc5d761))

### Features

* add ingress events on create, update and delete ([1df9534](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/1df9534b1c408e9491e9d8815f17b6ca16ffafea))
* allow to listen for resources in a namespace ([5cbdf0d](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/5cbdf0d824f36226c847a3876d6e28518baf03bb))
* allow users to customize the manager env ([37c8644](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/37c8644bed16fad2fc1552cc7f9e4b141da45fe3))
* handle ingress with multiple hosts ([e56b5ac](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/e56b5acb3d8b2078cc104f438e871b0cacdb065b))
* release the operator as a helm chart ([b182920](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/b182920ff8ffab87002f96236f933f64c7ed7b23))
* set definition context on create and updates ([520f710](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/520f710502c5cfc14984a51e675c960df4deb8da))

# [0.4.0](https://github.com/gravitee-io/gravitee-kubernetes-operator/compare/0.3.0...0.4.0) (2022-12-07)


### Bug Fixes

* error log typo ([a7377ec](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/a7377ec9ba2535307a3d435fa165fb7ed52ca629))


### Features

* add DEV_MODE logging option ([d1cae84](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/d1cae8487ad7627651e20026e40776087a3ff614))
* use `message` and `timestamp` keys in log ([15b75d4](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/15b75d483520e06eb245b4af8671d9f768564955))

# [0.3.0](https://github.com/gravitee-io/gravitee-kubernetes-operator/compare/0.2.0...0.3.0) (2022-11-23)


### Bug Fixes

* make metadata value optional ([737aa3d](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/737aa3dd31124c99ea9df6d1d13256a6dd76a024))


### Features

* reference api resource from name and namespace ([ac749ca](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/ac749caf7796ffcf9e7f44a532a28e20d56809bf))

# [0.2.0](https://github.com/gravitee-io/gravitee-kubernetes-operator/compare/0.1.1...0.2.0) (2022-11-09)


### Bug Fixes

* add rbac marker for secret lists ([9ed5735](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/9ed5735300acd5d208b485573a4915d0151bed6f))
* import api with disabled health check ([8698633](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/869863348960b00d7775088f7b988e0ae97a1e7f))
* import API with logging ([5b28322](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/5b2832235a4b57451c0aeabede356fd068014b50))
* import api with several endpoint groups ([c308730](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/c308730b5b1d66e375319d85646b254826f1c391))


### Features

* reference secret in context ([6d3acf6](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/6d3acf66b277fb00407096b0c862d472b93f45a3))

## [0.1.1](https://github.com/gravitee-io/gravitee-kubernetes-operator/compare/0.1.0...0.1.1) (2022-10-28)

### Features

* reconcile api resources on context updates ([c820c14](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/c820c1472d050e3676f3ff5823d1d530f31b5852))
  
### Bug Fixes

* add enabled in health check model ([c3098e3](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/c3098e3dd7e375c72697a14b64b6b0aaf3d94dd0))
* align endpoint mapping with apim ([568c879](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/568c8795a22345334a01273d115de7609043fac4))
* change fail over data type ([0fbe2bd](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/0fbe2bdd607fc431b92e428f94954e08a4fbe2a0))
* import api with life cycle state ([5293ddd](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/5293dddf5aee6f7373f3013e9bbeba7525ffd77c))
* merge create and update of api definition ([1e677f0](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/1e677f0de588eb4a37b1f59fd8ba384fcfbc6b52))
* rename cors fields to match v3 definition ([3ebd4d0](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/3ebd4d0cd3ee6f545f51e27e6fd087bfa618f7d5))

# [0.1.0](https://github.com/gravitee-io/gravitee-kubernetes-operator/compare/0.0.0...0.1.0) (2022-10-07)

### Features

* add events on api resource ([da695a7](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/da695a721e58ff5187484c258bb41ea8d9591434))
* add organization and environment to management context ([869be0d](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/869be0dc8cffbfc083e201b310a698921684423c))
* delete an api definition ([8b763be](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/8b763be49ec779fcdbd7682bbf41b4815060c4ea))
* start and stop api ([a58756f](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/a58756f276f06ec9e72de36847c6408719552895))
* create and update an api definition ([005ece9](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/005ece9c61744c5a3ebb1a449cbb935bfa1deb18))
