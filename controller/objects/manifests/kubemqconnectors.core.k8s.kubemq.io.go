package manifests

const(
	KubemqConnectorsCustomResourceDefinition =`
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: kubemqconnectors.core.k8s.kubemq.io
spec:
  additionalPrinterColumns:
    - JSONPath: .status.type
      name: Type
      type: string
    - JSONPath: .status.image
      name: Image
      type: string
    - JSONPath: .status.api
      name: API
      type: string
    - JSONPath: .status.status
      name: Status
      type: string
  group: core.k8s.kubemq.io
  names:
    kind: KubemqConnector
    listKind: KubemqConnectorList
    plural: kubemqconnectors
    singular: kubemqconnector
  scope: Namespaced
  subresources:
    scale:
      labelSelectorPath: .status.selector
      specReplicasPath: .spec.replicas
      statusReplicasPath: .status.replicas
    status: {}
  validation:
    openAPIV3Schema:
      description: KubemqConnector is the Schema for the kubemqconnectors API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: KubemqConnectorSpec defines the desired state of KubemqConnector
          properties:
            config:
              type: string
            image:
              type: string
            node_port:
              format: int32
              type: integer
            replicas:
              format: int32
              minimum: 0
              type: integer
            service_type:
              type: string
            type:
              type: string
          required:
            - config
            - type
          type: object
        status:
          description: KubemqConnectorStatus defines the observed state of KubemqConnector
          properties:
            api:
              type: string
            image:
              type: string
            replicas:
              format: int32
              type: integer
            status:
              type: string
            type:
              type: string
          required:
            - api
            - image
            - replicas
            - status
            - type
          type: object
      type: object
  version: v1alpha1
  versions:
    - name: v1alpha1
      served: true
      storage: true
`
)
