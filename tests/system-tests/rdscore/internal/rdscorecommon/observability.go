package rdscorecommon

import (
	"encoding/json"
	"fmt"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/nodes"
	"github.com/openshift-kni/eco-goinfra/pkg/pod"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/internal/apiobjectshelper"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/internal/loki_utils"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/internal/remote"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/rdscore/internal/rdscoreparams"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math/rand"
	"time"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/openshift-kni/eco-gotests/tests/system-tests/rdscore/internal/rdscoreinittools"
)

const (
	prometheusDeploymentLabel = "prometheus=k8s"
)

var (
	clusterObservabilityOperatorsSuite = []string{
		"obo-prometheus-operator",
		"obo-prometheus-operator-admission-webhook",
		"observability-operator",
		"perses-operator"}

	prometheusPodSelector = metav1.ListOptions{LabelSelector: prometheusDeploymentLabel}
)

// CreateLokiStackSecret asserts lokiStack secret exists.
func CreateLokiStackSecret() {
	By("Verify LokiStack Operator successfully deployed")

	glog.V(rdscoreparams.RDSCoreLogLevel).Infof("Verify LokiStack operator namespace %s defined",
		rdscoreparams.LokiNamespace)

	err := apiobjectshelper.VerifyNamespaceExists(APIClient, rdscoreparams.LokiNamespace, time.Second)

	if err != nil {
		glog.V(rdscoreparams.RDSCoreLogLevel).Infof(
			"%s namespace not found. Skipping LokiStack secret configuration stage...", rdscoreparams.LokiNamespace)

		Skip("LokiStack namespace not found. Skipping LokiStack secret configuration stage...")
	}

	glog.V(rdscoreparams.RDSCoreLogLevel).Infof("Verify LokiStack operator deployment validity")

	err = apiobjectshelper.VerifyOperatorDeployment(APIClient,
		rdscoreparams.LokiOperatorSubscriptionName,
		rdscoreparams.LokiOperatorDeploymentName,
		rdscoreparams.LokiNamespace,
		time.Minute)

	if err != nil {
		glog.V(rdscoreparams.RDSCoreLogLevel).Infof(
			"LokiStack operator deployment not found or broken. Skipping LokiStack secret configuration stage...")

		Skip("LokiStack operator deployment not found or broken. Skipping LokiStack secret configuration stage...")
	}

	By("Creating LokiStack secret")

	err = loki_utils.CreateLokiStackSecret(
		APIClient,
		RDSCoreConfig.LokiSecretName,
		RDSCoreConfig.LokiObjectBucketClaim,
		rdscoreparams.CLONamespace)
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("LokiStack secret %s failure in the namespace %s; %v",
			RDSCoreConfig.LokiSecretName, rdscoreparams.CLONamespace, err))
}

// ClusterObservabilityOperatorMonitoredByPlatform asserts cluster-observability operator is monitored by platform.
func ClusterObservabilityOperatorMonitoredByPlatform() {
	By("Verify Cluster-Observability Operator successfully deployed")

	glog.V(rdscoreparams.RDSCoreLogLevel).Infof("Verify Cluster-Observability operator namespace %s defined",
		rdscoreparams.COONamespace)

	err := apiobjectshelper.VerifyNamespaceExists(APIClient, rdscoreparams.CLONamespace, time.Second)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to pull namespace %q; %v",
		rdscoreparams.CLONamespace, err))

	glog.V(rdscoreparams.RDSCoreLogLevel).Infof("Verify Cluster-Observability operator deployment validity")

	for _, deploymentName := range clusterObservabilityOperatorsSuite {
		err := apiobjectshelper.VerifyOperatorDeployment(APIClient,
			rdscoreparams.COOName,
			deploymentName,
			rdscoreparams.COONamespace,
			time.Minute)
		Expect(err).ToNot(HaveOccurred(),
			fmt.Sprintf("Cluster-Observability operator deployment %s failure in the namespace %s; %v",
				deploymentName, rdscoreparams.COONamespace, err))
	}

	By("Verify Cluster-Observability Operator is monitored by platform")

	glog.V(rdscoreparams.RDSCoreLogLevel).Infof("Create SA %s token", rdscoreparams.COOServiceAccount)

	getTokenCmd := fmt.Sprintf("oc create token %s -n %s --kubeconfig=%s",
		rdscoreparams.COOServiceAccount, rdscoreparams.COOSANamespace, RDSCoreConfig.HypervisorKubeconfig)

	saToken, err := remote.ExecCmdOnHost(
		RDSCoreConfig.HypervisorHost,
		RDSCoreConfig.HypervisorUser,
		RDSCoreConfig.HypervisorPass,
		getTokenCmd)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to execute cmd %s; %v. \noutput: %v",
		getTokenCmd, err, saToken))

	glog.V(rdscoreparams.RDSCoreLogLevel).Infof("Generated SA token: %s", saToken)

	glog.V(rdscoreparams.RDSCoreLogLevel).Infof("Get prometheus test pod")

	podObjects, err := pod.List(APIClient, rdscoreparams.COOSANamespace, prometheusPodSelector)
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("Failed to retrieve pods list from namespace %s with label %v: %v",
			rdscoreparams.COOSANamespace, prometheusPodSelector, err))
	Expect(len(podObjects)).To(BeNumerically(">", 0),
		fmt.Sprintf("No pods found in namespace %s with label %v",
			rdscoreparams.COOSANamespace, prometheusPodSelector))

	prometheusTestPod := podObjects[0]

	glog.V(rdscoreparams.RDSCoreLogLevel).Infof(
		"Check that Cluster-Observability Operator is monitored by platform")

	cmdToRun := fmt.Sprintf("oc -n %s --kubeconfig=%s exec -c %s %s -- curl -k -H "+
		"\"Authorization: Bearer %s\" 'https://thanos-querier.openshift-monitoring.svc:9091/api/v1/query?' "+
		"--data-urlencode 'query=up{namespace=\"%s\"}==1' | jq",
		rdscoreparams.COOSANamespace,
		RDSCoreConfig.HypervisorKubeconfig,
		prometheusTestPod.Object.Spec.Containers[0].Name,
		prometheusTestPod.Definition.Name,
		saToken,
		rdscoreparams.COONamespace)

	glog.V(100).Infof("Execute command: %q", cmdToRun)

	var ctx SpecContext

	Eventually(func() bool {
		output, err := remote.ExecCmdOnHost(
			RDSCoreConfig.HypervisorHost,
			RDSCoreConfig.HypervisorUser,
			RDSCoreConfig.HypervisorPass,
			cmdToRun)

		if err != nil {
			glog.V(rdscoreparams.RDSCoreLogLevel).Infof("Error running command from within a pod %q: %v",
				prometheusTestPod.Object.Name, err)

			return false
		}

		glog.V(rdscoreparams.RDSCoreLogLevel).Infof("Command's output:\n\t%v", output)

		var result *loki_utils.ThanosQueryResponse

		err = json.Unmarshal(([]byte)(output), &loki_utils.ThanosQueryResponse{})

		glog.V(rdscoreparams.RDSCoreLogLevel).Infof("Unmarshal result:\n\t%v", result)

		return true
	}).WithContext(ctx).WithPolling(5*time.Second).WithTimeout(1*time.Minute).Should(BeTrue(),
		"Failed to run command from within pod")
}

func verifyCLONamespaceExists() {
	err := apiobjectshelper.VerifyNamespaceExists(APIClient, rdscoreparams.CLONamespace, time.Second)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to pull namespace %q; %v",
		rdscoreparams.CLONamespace, err))
}

func verifyLokiNamespaceExists() {
	err := apiobjectshelper.VerifyNamespaceExists(APIClient, rdscoreparams.LokiNamespace, time.Second)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to pull namespace %q; %v",
		rdscoreparams.LokiNamespace, err))
}

func verifyLokiDeployment() {
	err := apiobjectshelper.VerifyOperatorDeployment(APIClient,
		rdscoreparams.LokiOperatorSubscriptionName,
		rdscoreparams.LokiOperatorDeploymentName,
		rdscoreparams.LokiNamespace,
		time.Minute)
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("Loki operator deployment %s failure in the namespace %s; %v",
			rdscoreparams.LokiOperatorDeploymentName, rdscoreparams.LokiNamespace, err))
}

func VerifyCLODeployment() {
	err := apiobjectshelper.VerifyOperatorDeployment(APIClient,
		rdscoreparams.CLOName,
		rdscoreparams.CLODeploymentName,
		rdscoreparams.CLONamespace,
		time.Minute)
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("operator deployment %s failure in the namespace %s; %v",
			rdscoreparams.CLOName, rdscoreparams.CLONamespace, err))
}

func getRandomWorkerNode(apiClient *clients.Settings) (*nodes.Builder, error) {
	glog.V(100).Infof("Get random workers nodes list")

	workersList, err := nodes.List(apiClient, RDSCoreConfig.WorkerLabelListOption)

	if err != nil {
		glog.Errorf("Failed to list workers: %v", err)

		return nil, fmt.Errorf("failed to list workers: %v", err)
	}

	glog.V(100).Infof("Get random worker node")

	testNode := rand.Intn(len(workersList))

	glog.V(100).Infof("Chosen node for test: %s", workersList[testNode].Definition.Name)

	return workersList[testNode], nil
}
