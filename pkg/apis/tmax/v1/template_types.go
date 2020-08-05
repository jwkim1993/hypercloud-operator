package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +kubebuilder:validation:XPreserveUnknownFields
// +kubebuilder:validation:XEmbeddedResource

type ObjectSpec struct {
	Fields metav1.FieldsV1 `json:"fields,omitempty"`
}

// +kubebuilder:validation:XPreserveUnknownFields
type PlanSpec struct {
	Id                     string            `json:"id"`
	Name                   string            `json:"name"`
	Description            string            `json:"description,omitempty"`
	Metadata               map[string]string `json:"metadata,omitempty"`
	Free                   bool              `json:"free,omitempty"`
	Bindable               bool              `json:"bindable,omitempty"`
	PlanUpdateable         bool              `json:"plan_updateable"`
	Schemas                Schemas           `json:"schemas,omitempty"`
	MaximumPollingDuration int               `json:"maximum_polling_duration,omitempty"`
	MaintenanceInfo        MaintenanceInfo   `json:"maintenance_info,omitempty"`
}

type MaintenanceInfo struct {
	Version     string `json:"version"`
	Description string `json:"description,omitempty"`
}
type Schemas struct {
	ServiceInstance ServiceInstanceSchema `json:"service_instance,omitempty"`
	ServiceBinding  ServiceBindingSchema  `json:"service_binding,omitempty"`
}

type ServiceInstanceSchema struct {
	Create SchemaParameters `json:"create,omitempty"`
	Update SchemaParameters `json:"update,omitempty"`
}

type ServiceBindingSchema struct {
	Create SchemaParameters `json:"create,omitempty"`
}

type SchemaParameters struct {
	Parameters map[string]string `json:"parameters,omitempty"`
}
type ParamSpec struct {
	Description string `json:"description,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	From        string `json:"from,omitempty"`
	Generate    string `json:"generate,omitempty"`
	Name        string `json:"name"`
	Required    bool   `json:"required,omitempty"`
	Value       string `json:"value,omitempty"`
	ValueType   string `json:"valueType,omitempty"`
}

type LabelSpec struct {
	AdditionalProperties string `json:"additionalProperties,omitempty"`
}

type MetadataSpec struct {
	GenerateName      string `json:"generateName,omitempty"`
	Name              string `json:"name,omitempty"`
	metav1.ObjectMeta `json:"type,omitempty"`
}

// +kubebuilder:resource:shortName="tp"
// TemplateSpec defines the desired state of Template
type TemplateSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	OperatorStartTime   string       `json:"operatorStartTime,omitempty"`
	ApiVersion          string       `json:"apiVersion,omitempty"`
	Kind                string       `json:"kind,omitempty"`
	Labels              LabelSpec    `json:"labels,omitempty"`
	Message             string       `json:"message,omitempty"`
	ShortDescription    string       `json:"shortDescription,omitempty"`
	LongDescription     string       `json:"LongDescription,omitempty"`
	UrlDescription      string       `json:"UrlDescription,omitempty"`
	MarkDownDescription string       `json:"markDownDescription,omitempty"`
	Provider            string       `json:"provider,omitempty"`
	ImageUrl            string       `json:"imageUrl,omitempty"`
	Recommand           bool         `json:"recommend,omitempty"`
	Tags                []string     `json:"tags,omitempty"`
	ObjectKinds         []string     `json:"objectKinds,omitempty"`
	Metadata            MetadataSpec `json:"metadata,omitempty"`
	Objects             []ObjectSpec `json:"objects,omitempty"`
	Plans               []PlanSpec   `json:"plans,omitempty"`
	Parameters          []ParamSpec  `json:"parameters"`
}

// TemplateStatus defines the observed state of Template
type TemplateStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Template is the Schema for the templates API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=templates,scope=Namespaced
type Template struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              TemplateSpec `json:"spec,omitempty"`
	// Status TemplateStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TemplateList contains a list of Template
type TemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Template `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Template{}, &TemplateList{})
}
