## [0.8.4](https://github.com/gravitee-io/gravitee-kubernetes-operator/compare/0.8.3...0.8.4) (2023-09-28)


### Bug Fixes

* let users disable flow-steps ([a50d237](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/a50d23793dc60f519029af3dfc50e44d5a1ca247))
* unmarshal int values in GenericMap ([6ba35fc](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/6ba35fcb666cbaf920cb172547e7005465e5ec6b))
* update kube-rbac-proxy version ([2a27e41](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/2a27e41f5b73468d29bbac8fc462f77935aa085a))

## [0.8.3](https://github.com/gravitee-io/gravitee-kubernetes-operator/compare/0.8.2...0.8.3) (2023-09-15)


### Bug Fixes

* use PUT when setting definition context on APIM ([979b6d2](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/979b6d215814f37b6882151a4f0cc7473473c5b9))

## [0.8.2](https://github.com/gravitee-io/gravitee-kubernetes-operator/compare/0.8.1...0.8.2) (2023-09-14)


### Bug Fixes

* disable flows when `enabled` is false ([29d7438](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/29d7438ad71cfcdea463a2e153e933277a2c0d87))
* rollout on helm upgrades when config changes ([ee274b5](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/ee274b565661feb4f8330d7f6123a8a3bc4cced3))

## [0.8.1](https://github.com/gravitee-io/gravitee-kubernetes-operator/compare/0.8.0...0.8.1) (2023-09-11)


### Bug Fixes

* reconcile applications on generation change only ([b04e935](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/b04e93592717933ab1b80a4b47d44c5277acb2e8))
* set config namespace to release namespace ([69aab41](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/69aab4130f3d4ea306f66fb2c21334bd8fd75f7d))

# [0.8.0](https://github.com/gravitee-io/gravitee-kubernetes-operator/compare/0.7.0...0.8.0) (2023-08-10)


### Bug Fixes

* add flow id property ([b46cc43](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/b46cc430e2c79edc55800dca541e1d99492ecdf8))
* allow insecure skip verify in HTTPClient ([76f5268](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/76f52688e19ce7d6be576cee6d77a8ca99d4c8db))
* handle UI exports with endpoint level healthchecks ([3332554](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/33325543cf101c0d933f840df8953c2ab95613c5))
* management context check for non-local APIs ([68f58d7](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/68f58d74e12302b7a4569321c52e516d632aa139))
* set env and org id in application status on updates ([f978120](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/f978120cd2875764d698b3c3516fdfed032c7ab6))


### Features

* allow custom manager image and tag ([2529646](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/25296467e68b9258fc78a795d8d232718652ef24))

# [0.7.0](https://github.com/gravitee-io/gravitee-kubernetes-operator/compare/0.6.0...0.7.0) (2023-06-08)


### Bug Fixes

* remove hard-coded "keystore" key ([8ba89a5](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/8ba89a5b1a2ee9105f9611c12954a1ebfb191fce))
* wrong data type while unmarshalling ([8c4e0e9](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/8c4e0e9a03a9d1d4d8006065fb4d04101e6b06d6))


### Features

* template resolver ([5a93e6a](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/5a93e6a91442cda82e9a40b1268d12a5ba482207))

# [0.6.0](https://github.com/gravitee-io/gravitee-kubernetes-operator/compare/0.5.1...0.6.0) (2023-05-26)


### Bug Fixes

* add applications rbac ([1347acd](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/1347acd630ae840837466d8a180a6b744752e28c))
* add CRDs to the docker image ([9b75ebf](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/9b75ebfaf1151375a0fabe08dacf6493a9201b31))
* add deletion finalizer to context secret ([346501e](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/346501e350541fc69468ee7ec13e15402e2b1b19))
* add finalizer to api resources ([9ab9b17](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/9ab9b178cc4933d3474c32b145b7e9a8fc4d192c))
* add rbac for crd patch ([3da603c](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/3da603c05ffb65a3b2d1f99ff0f97a4cc444edef))
* add support for ingressClassName field ([e914b50](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/e914b504904f2ea44b89e088848be8228791e2ba))
* remove local flag on ApiResource ([5456bf5](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/5456bf540aa674128dd89c2908b6b125f8ba6bd0))
* resolve few bug ([46812f4](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/46812f4aea5897b923bc9512b701e6e483f6cc86))
* resolve race conditions on helm deletion ([22240d0](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/22240d09098739854ddeaab6b83cd6b0c51117d9))
* restore namespaces in resource refs ([9089861](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/9089861c635204ca2a96d766d1e49471ecbd1885))
* support different key type ([7d948c7](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/7d948c7369d2d293e6295abc68002ddefe06a333))


### Features

* add local flag in all samples ([4856932](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/4856932057248ef7c903605b3d7d5650a8588e52))
* api definition template ([9f6dc4c](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/9f6dc4c35d68c11703517864bb9904405aa1cdfb))
* application CRD ([0195f25](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/0195f25cac6c4d41be3dce1f7b1ffa029e7dc2b2))
* define ApiDefinition visibility in Kubenetes clusters ([70e92ee](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/70e92eef790f26f0eb4797fb10c45cbf8d60a72c))
* define custom not found response templates ([18c62b4](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/18c62b4dddd1542962c417a6f2d4a6cb11d2153f))
* handle ingress resources with multiple hosts ([1e40555](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/1e40555f5c39943e4154024703b8b1610ce42168))
* handle ingress tls option ([c66e023](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/c66e023a7152ecc1ae767d89a75619031204f52c))
* patch resource definitions on startup ([a523075](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/a52307530a1e3e705435f45fdfdf314f619b8bd2))
* use a role for configs and secrets if namespaced ([ca3d58b](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/ca3d58bf98dfbce34ae54bb4a668a06cc7c95bd7))

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
