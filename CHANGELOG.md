# [1.0.0-beta.5](https://github.com/gravitee-io/gravitee-kubernetes-operator/compare/1.0.0-beta.4...1.0.0-beta.5) (2024-01-05)


### Features

* introduce pem registry ([3ca58f6](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/3ca58f6a8930731b346e438863a3099986f5e776))

# [1.0.0-beta.4](https://github.com/gravitee-io/gravitee-kubernetes-operator/compare/1.0.0-beta.3...1.0.0-beta.4) (2024-01-05)


### Bug Fixes

* allow support for http2 requests on ingresses ([cc05c9a](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/cc05c9a449e8f4da983b1728211d0227207d3b61))

# [1.0.0-beta.3](https://github.com/gravitee-io/gravitee-kubernetes-operator/compare/1.0.0-beta.2...1.0.0-beta.3) (2023-12-21)


### Features

* allow to configure service monitor from helm ([4bab3cb](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/4bab3cb64502cd41a226b3833dab987dfa316156))
* allow to configure service monitor from helm ([1cd2818](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/1cd281888f2c17a141b8a43532ab574a47f8231b))

# [1.0.0-beta.2](https://github.com/gravitee-io/gravitee-kubernetes-operator/compare/1.0.0-beta.1...1.0.0-beta.2) (2023-12-21)


### Bug Fixes

* default manager image tag to chart version ([fcad3b3](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/fcad3b36a7553517cd1a41f090385831dc871529))

# [1.0.0-beta.1](https://github.com/gravitee-io/gravitee-kubernetes-operator/compare/0.5.0...1.0.0-beta.1) (2023-12-15)


### Bug Fixes

* add applications rbac ([1347acd](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/1347acd630ae840837466d8a180a6b744752e28c))
* add CRDs to the docker image ([9b75ebf](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/9b75ebfaf1151375a0fabe08dacf6493a9201b31))
* add deletion finalizer to context secret ([346501e](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/346501e350541fc69468ee7ec13e15402e2b1b19))
* add finalizer to api resources ([9ab9b17](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/9ab9b178cc4933d3474c32b145b7e9a8fc4d192c))
* add flow id property ([0b802e8](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/0b802e80dfe9d7e9fa3648dbe38529306be8db89))
* add missing config key, update API Definition CRD ([e21c229](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/e21c2293f9e4a2e4a9a5c2078caf147efec1305a))
* add rbac for crd patch ([3da603c](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/3da603c05ffb65a3b2d1f99ff0f97a4cc444edef))
* add support for ingressClassName field ([e914b50](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/e914b504904f2ea44b89e088848be8228791e2ba))
* allow insecure skip verify in HTTPClient ([3f600bb](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/3f600bbdc5cee09ca422fdb6b5df126aaf5467ab))
* disable flows when `enabled` is false ([1565797](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/1565797efcaf223ea72c902bcac865c4ef368920))
* enable/disable api when it is not local ([fa7c6f8](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/fa7c6f8743ef99973b5249b0f3d16009d8c78df3))
* ensure ns is set on context secret requests ([3995fae](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/3995faeb85246fb8c8eb40ee09b330e84b84fe62))
* handle UI exports with endpoint level healthchecks ([3cc152a](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/3cc152ad7b5bb446ae584ac318d9a3f433ca58f6))
* **helm:** avoid conflicting type in manager config map ([0fcb28e](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/0fcb28e76f0dc39029a7d4b5c5ae170b548e789b))
* let users disable flow-steps ([d38a116](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/d38a11662c1a8d1c9ed8f8c567d2883d2d519c36))
* management context check for non-local APIs ([55237dd](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/55237dd350e744c9000b905e2acf1577af2ef753))
* reconcile applications on generation change only ([f5ca10e](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/f5ca10ed4eb6cc968054057b722c9369de5d4799))
* reconcile ingress when the template annotation is updated ([686d6c7](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/686d6c769171681ac61510943639e0e098630c6e))
* remove hard-coded "keystore" key ([c33a634](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/c33a6349bf757ae08218eae7c77f9c3b7caba4fd))
* remove local flag on ApiResource ([5456bf5](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/5456bf540aa674128dd89c2908b6b125f8ba6bd0))
* resolve few bug ([46812f4](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/46812f4aea5897b923bc9512b701e6e483f6cc86))
* resolve race conditions on helm deletion ([22240d0](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/22240d09098739854ddeaab6b83cd6b0c51117d9))
* restore and deprecate v0.4.0 status fields ([536e806](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/536e806b567cf6c36d7ed64cf615f12a6804cecc))
* restore namespaces in resource refs ([9089861](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/9089861c635204ca2a96d766d1e49471ecbd1885))
* rollout on helm upgrades when config changes ([1416f5e](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/1416f5ec63373a47bac11a85afb915b7fa09ef4b))
* set config namespace to release namespace ([27e7c58](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/27e7c58cffa76354db3c51c432490812a5a4757b))
* set env and org id in application status on updates ([5aaabbd](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/5aaabbdb09e9e965261a4d1baed980de6b63a689))
* support different key type ([7d948c7](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/7d948c7369d2d293e6295abc68002ddefe06a333))
* unmarshal int values in GenericMap ([a7f3e7c](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/a7f3e7cf95783b04bb7af8f3fa658c1fa8c0f93b))
* update kube-rbac-proxy version ([d5c4a4e](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/d5c4a4e9ae7890ce098292ccf73014d721c6e74d))
* use annotation for ingress templates everywhere ([e22ef46](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/e22ef46fe565a9b4089152b23fed581f8cdd3712))
* use PUT when setting definition context on APIM ([b46adf5](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/b46adf53f164e66858fff7a5f833217fdf5a3512))
* wrong data type while unmarshalling ([7e72a47](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/7e72a477ef120579c19439423c1e884e872a5408))


### Features

* add local flag in all samples ([4856932](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/4856932057248ef7c903605b3d7d5650a8588e52))
* allow custom manager image and tag ([5ae681d](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/5ae681d81df907eeb490fc1aa48e98221efb4829))
* allow install without any cluster role ([519ee04](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/519ee048fe7797b5e467fe1aae3461efaf10464f))
* api definition template ([9f6dc4c](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/9f6dc4c35d68c11703517864bb9904405aa1cdfb))
* application CRD ([0195f25](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/0195f25cac6c4d41be3dce1f7b1ffa029e7dc2b2))
* bring support for APIM v4 API definitions ([d52fc81](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/d52fc81efe9c6b14b7d57950fdf6e761113c8670))
* define ApiDefinition visibility in Kubenetes clusters ([70e92ee](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/70e92eef790f26f0eb4797fb10c45cbf8d60a72c))
* define custom not found response templates ([18c62b4](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/18c62b4dddd1542962c417a6f2d4a6cb11d2153f))
* define resources through values ([#494](https://github.com/gravitee-io/gravitee-kubernetes-operator/issues/494)) ([0a92276](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/0a92276e06796433340844eba19f8c899b89c870))
* handle ingress resources with multiple hosts ([1e40555](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/1e40555f5c39943e4154024703b8b1610ce42168))
* handle ingress tls option ([c66e023](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/c66e023a7152ecc1ae767d89a75619031204f52c))
* **helm:** add helm standard labels ([0301ad7](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/0301ad7171e05fcb830fdc8dd880cdd2c22e6acd))
* make http client timeout configurable ([4b6a125](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/4b6a125389464aed4940244f0d80197fd0a342f3))
* patch resource definitions on startup ([a523075](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/a52307530a1e3e705435f45fdfdf314f619b8bd2))
* template resolver ([ca359cc](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/ca359cc68c260a9b8161be44315180e5529ff30f))
* use a role for configs and secrets if namespaced ([ca3d58b](https://github.com/gravitee-io/gravitee-kubernetes-operator/commit/ca3d58bf98dfbce34ae54bb4a668a06cc7c95bd7))


### BREAKING CHANGES

* This version requires APIM 4.x and upper

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
