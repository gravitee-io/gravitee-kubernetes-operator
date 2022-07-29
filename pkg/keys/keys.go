package keys

// Kubernetes Ingresses.
const (
	IngressLabel                = "gravitee.io/ingress"
	IngressClassAnnotation      = "kubernetes.io/ingress.class"
	IngressClassAnnotationValue = "graviteeio"
	IngressTemplateAnnotation   = "gravitee.io/template"
	IngressFinalizer            = "finalizers.gravitee.io/ingress"
)

// Gravitee.io CRDs.
const (
	CrdGroup   = "gravitee.io"
	CrdVersion = "v1alpha1"

	CrdManagementContextResource = "managementcontext"
	CrdApiDefinitionResource     = "apidefinitions"
	CrdApiDefinitionTemplate     = "template"
)

// Kubernetes Finalizers.
const (
	ApiDefinitionFinalizer         = "finalizers.gravitee.io/apidefinition"
	ApiDefinitionTemplateFinalizer = "finalizers.gravitee.io/apidefinitiontemplate"
)
