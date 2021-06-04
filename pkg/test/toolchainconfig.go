package test

import (
	toolchainv1alpha1 "github.com/codeready-toolchain/api/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ToolchainConfigOptionFunc func(config *toolchainv1alpha1.ToolchainConfig)

type ToolchainConfigOption interface {
	Apply(config *toolchainv1alpha1.ToolchainConfig)
}

type ToolchainConfigOptionImpl struct {
	toApply []ToolchainConfigOptionFunc
}

func (option *ToolchainConfigOptionImpl) Apply(config *toolchainv1alpha1.ToolchainConfig) {
	for _, apply := range option.toApply {
		apply(config)
	}
}

func (option *ToolchainConfigOptionImpl) addFunction(funcToAdd ToolchainConfigOptionFunc) {
	option.toApply = append(option.toApply, funcToAdd)
}

type AutomaticApprovalOption struct {
	*ToolchainConfigOptionImpl
}

func AutomaticApproval() *AutomaticApprovalOption {
	o := &AutomaticApprovalOption{
		ToolchainConfigOptionImpl: &ToolchainConfigOptionImpl{},
	}
	o.addFunction(func(config *toolchainv1alpha1.ToolchainConfig) {
		config.Spec.Host.AutomaticApproval = toolchainv1alpha1.AutomaticApproval{}
	})
	return o
}

func (o AutomaticApprovalOption) Enabled() AutomaticApprovalOption {
	o.addFunction(func(config *toolchainv1alpha1.ToolchainConfig) {
		val := true
		config.Spec.Host.AutomaticApproval.Enabled = &val
	})
	return o
}

func (o AutomaticApprovalOption) Disabled() AutomaticApprovalOption {
	o.addFunction(func(config *toolchainv1alpha1.ToolchainConfig) {
		val := false
		config.Spec.Host.AutomaticApproval.Enabled = &val
	})
	return o
}

type DeactivationOption struct {
	*ToolchainConfigOptionImpl
}

func Deactivation() *DeactivationOption {
	o := &DeactivationOption{
		ToolchainConfigOptionImpl: &ToolchainConfigOptionImpl{},
	}
	o.addFunction(func(config *toolchainv1alpha1.ToolchainConfig) {
		config.Spec.Host.Deactivation = toolchainv1alpha1.Deactivation{}
	})
	return o
}

func (o DeactivationOption) DeactivatingNotificationDays(days int) DeactivationOption {
	o.addFunction(func(config *toolchainv1alpha1.ToolchainConfig) {
		config.Spec.Host.Deactivation.DeactivatingNotificationDays = &days
	})
	return o
}

type PerMemberClusterOption func(map[string]int)

func PerMemberCluster(name string, value int) PerMemberClusterOption {
	return func(clusters map[string]int) {
		clusters[name] = value
	}
}

func (o AutomaticApprovalOption) ResourceCapThreshold(defaultThreshold int, perMember ...PerMemberClusterOption) AutomaticApprovalOption {
	o.addFunction(func(config *toolchainv1alpha1.ToolchainConfig) {
		config.Spec.Host.AutomaticApproval.ResourceCapacityThreshold.DefaultThreshold = &defaultThreshold
		config.Spec.Host.AutomaticApproval.ResourceCapacityThreshold.SpecificPerMemberCluster = map[string]int{}
		for _, add := range perMember {
			add(config.Spec.Host.AutomaticApproval.ResourceCapacityThreshold.SpecificPerMemberCluster)
		}
	})
	return o
}

func (o AutomaticApprovalOption) MaxUsersNumber(overall int, perMember ...PerMemberClusterOption) AutomaticApprovalOption {
	o.addFunction(func(config *toolchainv1alpha1.ToolchainConfig) {
		config.Spec.Host.AutomaticApproval.MaxNumberOfUsers.Overall = &overall
		config.Spec.Host.AutomaticApproval.MaxNumberOfUsers.SpecificPerMemberCluster = map[string]int{}
		for _, add := range perMember {
			add(config.Spec.Host.AutomaticApproval.MaxNumberOfUsers.SpecificPerMemberCluster)
		}
	})
	return o
}

func NewToolchainConfig(options ...ToolchainConfigOption) *toolchainv1alpha1.ToolchainConfig {
	toolchainConfig := &toolchainv1alpha1.ToolchainConfig{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: HostOperatorNs,
			Name:      "config",
		},
	}
	for _, option := range options {
		option.Apply(toolchainConfig)
	}
	return toolchainConfig
}