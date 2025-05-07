package loki_utils

// LogEntity the entity of log data
type LogEntity struct {
	Kubernetes struct {
		Annotations       map[string]string `json:"annotations,omitempty"`
		ContainerID       string            `json:"container_id,omitempty"`
		ContainerImage    string            `json:"container_image"`
		ContainerImageID  string            `json:"container_image_id,omitempty"`
		ContainerIOStream string            `json:"container_iostream,omitempty"`
		ContainerName     string            `json:"container_name"`
		FlatLabels        []string          `json:"flat_labels"`
		Host              string            `json:"host"`
		Lables            map[string]string `json:"labels,omitempty"`
		MasterURL         string            `json:"master_url,omitempty"`
		NamespaceID       string            `json:"namespace_id"`
		NamespaceLabels   map[string]string `json:"namespace_labels,omitempty"`
		NamespaceName     string            `json:"namespace_name"`
		PodID             string            `json:"pod_id"`
		PodIP             string            `json:"pod_ip,omitempty"`
		PodName           string            `json:"pod_name"`
		PodOwner          string            `json:"pod_owner"`
	} `json:"kubernetes,omitempty"`
	Systemd struct {
		SystemdT struct {
			SystemdInvocationID string `json:"SYSTEMD_INVOCATION_ID"`
			BootID              string `json:"BOOT_ID"`
			GID                 string `json:"GID"`
			CmdLine             string `json:"CMDLINE"`
			PID                 string `json:"PID"`
			SystemSlice         string `json:"SYSTEMD_SLICE"`
			SelinuxContext      string `json:"SELINUX_CONTEXT"`
			UID                 string `json:"UID"`
			StreamID            string `json:"STREAM_ID"`
			Transport           string `json:"TRANSPORT"`
			Comm                string `json:"COMM"`
			EXE                 string
			SystemdUnit         string `json:"SYSTEMD_UNIT"`
			CapEffective        string `json:"CAP_EFFECTIVE"`
			MachineID           string `json:"MACHINE_ID"`
			SystemdCgroup       string `json:"SYSTEMD_CGROUP"`
		} `json:"t"`
		SystemdU struct {
			SyslogIdntifier string `json:"SYSLOG_IDENTIFIER"`
			SyslogFacility  string `json:"SYSLOG_FACILITY"`
		} `json:"u"`
	} `json:"systemd,omitempty"`
	ViaqMsgID string `json:"viaq_msg_id,omitempty"`
	Level     string `json:"level"`
	LogSource string `json:"log_source"`
	LogType   string `json:"log_type,omitempty"`
	Message   string `json:"message"`
	Docker    struct {
		ContainerID string `json:"container_id"`
	} `json:"docker,omitempty"`
	HostName  string `json:"hostname"`
	TimeStamp string `json:"@timestamp"`
	File      string `json:"file,omitempty"`
	OpenShift struct {
		ClusterID string            `json:"cluster_id,omitempty"`
		Sequence  int64             `json:"sequence"`
		Labels    map[string]string `json:"labels,omitempty"`
	} `json:"openshift,omitempty"`
	PipelineMetadata struct {
		Collector struct {
			ReceivedAt string `json:"received_at"`
			Name       string `json:"name"`
			InputName  string `json:"inputname"`
			Version    string `json:"version"`
			IPaddr4    string `json:"ipaddr4"`
		} `json:"collector"`
	} `json:"pipeline_metadata,omitempty"`
	Structured struct {
		Level        string `json:"level,omitempty"`
		StringNumber string `json:"StringNumber,omitempty"`
		Message      string `json:"message,omitempty"`
		Number       int    `json:"Number,omitempty"`
		Layer1       string `json:"Layer1,omitempty"`
		FooColonBar  string `json:"foo:bar,omitempty"`
		FooDotBar    string `json:"foo.bar,omitempty"`
		BraceItem    string `json:"{foobar},omitempty"`
		BracketItem  string `json:"[foobar],omitempty"`
		Layer2       struct {
			Name string `json:"name,omitempty"`
			Tips string `json:"tips,omitempty"`
		} `json:"layer2,omitempty"`
	} `json:"structured,omitempty"`
}

type lokiQueryResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Stream *struct {
				DetectedLevel           string `json:"detected_level,omitempty"`
				K8sContainerName        string `json:"k8s_container_name,omitempty"`
				K8sNamespaceName        string `json:"k8s_namespace_name,omitempty"`
				K8sNodeName             string `json:"k8s_node_name,omitempty"`
				K8sPodName              string `json:"k8s_pod_name,omitempty"`
				K8sPodUID               string `json:"k8s_pod_uid,omitempty"`
				LogType                 string `json:"log_type,omitempty"`
				Tag                     string `json:"tag,omitempty"`
				FluentdThread           string `json:"fluentd_thread,omitempty"`
				KubernetesContainerName string `json:"kubernetes_container_name,omitempty"`
				KubernetesHost          string `json:"kubernetes_host,omitempty"`
				KubernetesNamespaceName string `json:"kubernetes_namespace_name,omitempty"`
				KubernetesPodName       string `json:"kubernetes_pod_name,omitempty"`
				LogIOStream             string `json:"log_iostream,omitempty"`
				LogSource               string `json:"log_source,omitempty"`
				ObservedTimestamp       string `json:"observed_timestamp,omitempty"`
				OpenshiftClusterID      string `json:"openshift_cluster_id,omitempty"`
				OpenshiftClusterUID     string `json:"openshift_cluster_uid,omitempty"`
				OpenshiftLogSource      string `json:"openshift_log_source,omitempty"`
				OpenshiftLogType        string `json:"openshift_log_type,omitempty"`
				SeverityText            string `json:"severity_text,omitempty"`
			} `json:"stream,omitempty"`
			Metric *struct {
				LogType                 string `json:"log_type,omitempty"`
				KubernetesContainerName string `json:"kubernetes_container_name,omitempty"`
				KubernetesHost          string `json:"kubernetes_host,omitempty"`
				KubernetesNamespaceName string `json:"kubernetes_namespace_name,omitempty"`
				KubernetesPodName       string `json:"kubernetes_pod_name,omitempty"`
			} `json:"metric,omitempty"`
			Values []interface{} `json:"values,omitempty"`
			Value  interface{}   `json:"value,omitempty"`
		} `json:"result"`
		Stats struct {
			Summary struct {
				BytesProcessedPerSecond int     `json:"bytesProcessedPerSecond"`
				LinesProcessedPerSecond int     `json:"linesProcessedPerSecond"`
				TotalBytesProcessed     int     `json:"totalBytesProcessed"`
				TotalLinesProcessed     int     `json:"totalLinesProcessed"`
				ExecTime                float32 `json:"execTime"`
			} `json:"summary"`
			Store struct {
				TotalChunksRef        int `json:"totalChunksRef"`
				TotalChunksDownloaded int `json:"totalChunksDownloaded"`
				ChunksDownloadTime    int `json:"chunksDownloadTime"`
				HeadChunkBytes        int `json:"headChunkBytes"`
				HeadChunkLines        int `json:"headChunkLines"`
				DecompressedBytes     int `json:"decompressedBytes"`
				DecompressedLines     int `json:"decompressedLines"`
				CompressedBytes       int `json:"compressedBytes"`
				TotalDuplicates       int `json:"totalDuplicates"`
			} `json:"store"`
			Ingester struct {
				TotalReached       int `json:"totalReached"`
				TotalChunksMatched int `json:"totalChunksMatched"`
				TotalBatches       int `json:"totalBatches"`
				TotalLinesSent     int `json:"totalLinesSent"`
				HeadChunkBytes     int `json:"headChunkBytes"`
				HeadChunkLines     int `json:"headChunkLines"`
				DecompressedBytes  int `json:"decompressedBytes"`
				DecompressedLines  int `json:"decompressedLines"`
				CompressedBytes    int `json:"compressedBytes"`
				TotalDuplicates    int `json:"totalDuplicates"`
			} `json:"ingester"`
		} `json:"stats"`
	} `json:"data"`
}

type ThanosQueryResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Name       string `json:"__name__,omitempty"`
				Endpoint   string `json:"endpoint,omitempty"`
				Instance   string `json:"instance,omitempty"`
				Job        string `json:"job,omitempty"`
				Namespace  string `json:"namespace,omitempty"`
				Pod        string `json:"pod,omitempty"`
				Prometheus string `json:"prometheus,omitempty"`
				Service    string `json:"service,omitempty"`
			} `json:"metric,omitempty"`
			Value interface{} `json:"value,omitempty"`
		} `json:"result"`
	} `json:"data"`
	Analysis struct{}
}
