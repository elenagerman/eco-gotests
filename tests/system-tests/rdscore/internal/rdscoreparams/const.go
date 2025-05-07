package rdscoreparams

const (
	// Label is used to select tests for RDS Core setup.
	Label = "rdscore"

	// RDSCoreLogLevel configures logging level for RDS Core related tests.
	RDSCoreLogLevel = 90

	// NMStateInstanceName is a name of the NMState instance.
	NMStateInstanceName = "nmstate"

	// MachineConfidDaemonPodSelector is a a label selector for all machine-config-daemon pods.
	MachineConfidDaemonPodSelector = "k8s-app=machine-config-daemon"

	// LabelValidatePerformanceProfile is a test selector for performance profile validation.
	LabelValidatePerformanceProfile = "rds-core-performance-profile"

	// MachineConfigDaemonContainerName is a name of container within machine-config-daemon pod.
	MachineConfigDaemonContainerName = "machine-config-daemon"

	// LabelValidateNMState a label to select tests for NMState validation.
	LabelValidateNMState = "rds-core-validate-nmstate"

	// ConditionTypeReadyString constant to fix linter warning.
	ConditionTypeReadyString = "Ready"

	// ConstantTrueString constant to fix linter warning.
	ConstantTrueString = "True"

	// LabelValidateSRIOV a label to select tests for SR-IOV validation.
	LabelValidateSRIOV = "rds-core-validate-sriov"

	// MetalLBOperatorNamespace MetalLB operator namespace.
	MetalLBOperatorNamespace = "metallb-system"

	// MetalLBFRRPodSelector pod selector for MetalLB-FRR pods.
	MetalLBFRRPodSelector = "app=frr-k8s"

	// MetalLBFRRContainerName name of the FRR container within a pod.
	MetalLBFRRContainerName = "frr"

	// COONamespace is a cluster-observability operator namespace.
	COONamespace = "openshift-cluster-observability-operator"

	// COOName is a cluster-observability operator name.
	COOName = "cluster-observability-operator"

	// COOServiceAccount is a cluster-observability serviceAccount name.
	COOServiceAccount = "prometheus-k8s"

	// COOSANamespace is a cluster-observability serviceAccount namespace.
	COOSANamespace = "openshift-monitoring"

	// ODFNamespace is an odf namespace.
	ODFNamespace = "openshift-storage"

	// CLONamespace is a clusterlogging operator namespace.
	CLONamespace = "openshift-logging"

	// CLOName is a clusterlogging operator name.
	CLOName = "cluster-logging"

	// CLODeploymentName is a clusterlogging operator deployment name.
	CLODeploymentName = "cluster-logging-operator"

	// CLOInstanceName is a clusterlogging instance name.
	CLOInstanceName = "instance"

	// LokiNamespace is a loki operator namespace.
	LokiNamespace = "openshift-operators-redhat"

	// LokiOperatorSubscriptionName is a loki operator subscription name.
	LokiOperatorSubscriptionName = "loki-operator"

	// LokiOperatorDeploymentName is a loki operator deployment name.
	LokiOperatorDeploymentName = "loki-operator-controller-manager"

	// LokiStackName is a lokiStack instance name.
	LokiStackName = "logging-loki"

	// KubeconfigPath is a path to the kubeconfig file on the hypervisor machine.
	KubeconfigPath = "logging-loki"
)
