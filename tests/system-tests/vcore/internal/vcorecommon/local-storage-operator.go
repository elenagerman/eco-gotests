package vcorecommon

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/openshift-kni/eco-goinfra/pkg/mco"
	"github.com/openshift-kni/eco-goinfra/pkg/nodes"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/internal/remote"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/vcore/internal/ocpcli"

	"github.com/openshift-kni/eco-goinfra/pkg/lso"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/internal/await"
	lsov1 "github.com/openshift/local-storage-operator/api/v1"
	lsov1alpha1 "github.com/openshift/local-storage-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift-kni/eco-goinfra/pkg/pod"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/internal/apiobjectshelper"
	. "github.com/openshift-kni/eco-gotests/tests/system-tests/vcore/internal/vcoreinittools"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/vcore/internal/vcoreparams"
)

// VerifyLSOSuite container that contains tests for LSO verification.
func VerifyLSOSuite() {
	Describe(
		"LSO validation",
		Label(vcoreparams.LabelVCoreLSO), func() {
			It(fmt.Sprintf("Verifies %s namespace exists", vcoreparams.LSONamespace),
				Label("lso"), VerifyLSONamespaceExists)

			It("Verify Local Storage Operator successfully installed",
				Label("lso"), reportxml.ID("59491"), VerifyLSODeployment)

			It("Verify localvolumeset instance exists",
				Label("lso"), reportxml.ID("74918"), VerifyLocalVolumeSet)

			It("Apply taints to the ODF nodes",
				Label("lso"), reportxml.ID("74916"), LabelODFNodesAndSetTaints)
		})
}

// VerifyLSONamespaceExists asserts namespace for Local Storage Operator exists.
func VerifyLSONamespaceExists(ctx SpecContext) {
	err := apiobjectshelper.VerifyNamespaceExists(APIClient, vcoreparams.LSONamespace, time.Second)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to pull namespace %q; %v", vcoreparams.LSONamespace, err))
} // func VerifyLSONamespaceExists (ctx SpecContext)

// VerifyLSODeployment asserts Local Storage Operator successfully installed.
func VerifyLSODeployment(ctx SpecContext) {
	err := apiobjectshelper.VerifyOperatorDeployment(APIClient,
		vcoreparams.LSOName,
		vcoreparams.LSOName,
		vcoreparams.LSONamespace,
		time.Minute)
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("operator deployment %s failure in the namespace %s; %v",
			vcoreparams.LSOName, vcoreparams.LSONamespace, err))

	glog.V(vcoreparams.VCoreLogLevel).Infof("Confirm that LSO %s pod was deployed and running in %s namespace",
		vcoreparams.LSOName, vcoreparams.LSONamespace)

	lsoPods, err := pod.ListByNamePattern(APIClient, vcoreparams.LSOName, vcoreparams.LSONamespace)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("No %s pods were found in %s namespace due to %v",
		vcoreparams.LSOName, vcoreparams.LSONamespace, err))
	Expect(len(lsoPods)).ToNot(Equal(0), fmt.Sprintf("The list of pods %s found in namespace %s is empty",
		vcoreparams.LSOName, vcoreparams.LSONamespace))

	lsoPod := lsoPods[0]
	lsoPodName := lsoPod.Object.Name

	err = lsoPod.WaitUntilReady(time.Second)
	if err != nil {
		lsoPodLog, _ := lsoPod.GetLog(600*time.Second, vcoreparams.LSOName)
		glog.Fatalf("%s pod in %s namespace in a bad state: %s",
			lsoPodName, vcoreparams.LSONamespace, lsoPodLog)
	}
} // func VerifyLSODeployment (ctx SpecContext)

// VerifyLocalVolumeSet asserts localvolumeset instance exists.
func VerifyLocalVolumeSet(ctx SpecContext) {
	glog.V(vcoreparams.VCoreLogLevel).Infof("Create localvolumeset instance %s in namespace %s if not found",
		vcoreparams.ODFLocalVolumeSetName, vcoreparams.LSONamespace)

	var err error

	localVolumeSetObj := lso.NewLocalVolumeSetBuilder(APIClient,
		vcoreparams.ODFLocalVolumeSetName,
		vcoreparams.LSONamespace)

	nodeSelector := corev1.NodeSelector{NodeSelectorTerms: []corev1.NodeSelectorTerm{{
		MatchExpressions: []corev1.NodeSelectorRequirement{{
			Key:      "kubernetes.io/hostname",
			Operator: "In",
			Values:   odfNodesList,
		}}},
	}}

	deviceInclusionSpec := lsov1alpha1.DeviceInclusionSpec{
		DeviceTypes:                []lsov1alpha1.DeviceType{lsov1alpha1.RawDisk},
		DeviceMechanicalProperties: []lsov1alpha1.DeviceMechanicalProperty{lsov1alpha1.NonRotational},
	}

	tolerations := []corev1.Toleration{{
		Key:      "node.ocs.openshift.io/storage",
		Operator: "Equal",
		Value:    "true",
		Effect:   "NoSchedule",
	}}

	_, err = localVolumeSetObj.WithNodeSelector(nodeSelector).
		WithStorageClassName(vcoreparams.StorageClassName).
		WithVolumeMode(lsov1.PersistentVolumeBlock).
		WithFSType("ext4").
		WithMaxDeviceCount(int32(42)).
		WithDeviceInclusionSpec(deviceInclusionSpec).
		WithTolerations(tolerations).Create()
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("failed to create localvolumeset %s in namespace %s "+
		"due to %v", vcoreparams.ODFLocalVolumeSetName, vcoreparams.LSONamespace, err))

	pvLabel := fmt.Sprintf("storage.openshift.com/owner-name=%s", vcoreparams.ODFLocalVolumeSetName)

	err = await.WaitUntilPersistentVolumeCreated(APIClient,
		3,
		15*time.Minute,
		metav1.ListOptions{LabelSelector: pvLabel})
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("failed to create persistentVolumes due to %v", err))
} // func VerifyLocalVolumeSet (ctx SpecContext)

// LabelODFNodesAndSetTaints asserts ODF nodes taints configuration.
func LabelODFNodesAndSetTaints(ctx SpecContext) {
	glog.V(vcoreparams.VCoreLogLevel).Infof("Create new mcp %s", VCoreConfig.OdfMCPName)
	odfMcp := mco.NewMCPBuilder(APIClient, VCoreConfig.OdfMCPName)

	if !odfMcp.Exists() {
		odfMCPTemplateName := "odf-mcp.yaml"
		varsToReplace := make(map[string]interface{})
		varsToReplace["MCPName"] = VCoreConfig.OdfMCPName

		workingDir, err := os.Getwd()
		Expect(err).ToNot(HaveOccurred(), err)

		templateDir := filepath.Join(workingDir, vcoreparams.TemplateFilesFolder)

		err = ocpcli.ApplyConfig(
			filepath.Join(templateDir, odfMCPTemplateName),
			filepath.Join(vcoreparams.ConfigurationFolderPath, odfMCPTemplateName),
			varsToReplace)
		Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to create mcp %s", VCoreConfig.OdfMCPName))

		err = odfMcp.WaitForUpdate(3 * time.Minute)
		Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to update mcp %s", VCoreConfig.OdfMCPName))
	}

	for _, odfNode := range odfNodesList {
		currentODFNode, err := nodes.Pull(APIClient, odfNode)
		Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("failed to retrieve node %s object due to %v",
			odfNode, err))

		glog.V(vcoreparams.VCoreLogLevel).Infof("Change node %s role to the %s", odfNode, VCoreConfig.OdfMCPName)

		_, err = currentODFNode.
			WithNewLabel("custom-label/used", "").
			WithNewLabel("cluster.ocs.openshift.io/openshift-storage", "").
			WithNewLabel("node-role.kubernetes.io/infra", "").
			WithNewLabel("node-role.kubernetes.io/odf", "").
			RemoveLabel("node-role.kubernetes.io/worker", "").Update()
		Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("failed to update labels for the node %s due to %v",
			odfNode, err))

		glog.V(vcoreparams.VCoreLogLevel).Infof("Insure taints applyed to the %s node", odfNode)

		applyTaintsCmd := fmt.Sprintf(
			"oc adm taint node %s node.ocs.openshift.io/storage=true:NoSchedule --overwrite=true --kubeconfig=%s",
			odfNode, VCoreConfig.KubeconfigPath)
		_, err = remote.ExecCmdOnHost(VCoreConfig.Host, VCoreConfig.User, VCoreConfig.Pass, applyTaintsCmd)
		Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("failed to execute %s script due to %v",
			applyTaintsCmd, err))
	}

	glog.V(vcoreparams.VCoreLogLevel).Infof("Wait for the mcp %s to update", VCoreConfig.OdfMCPName)
	time.Sleep(3 * time.Second)

	err := odfMcp.WaitForUpdate(3 * time.Minute)
	Expect(err).To(BeNil(), fmt.Sprintf("Failed to create mcp %s", VCoreConfig.OdfMCPName))
} // func LabelODFNodesAndSetTaints (ctx SpecContext)
