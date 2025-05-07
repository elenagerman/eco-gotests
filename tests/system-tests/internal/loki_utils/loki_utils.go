package loki_utils

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/configmap"
	"github.com/openshift-kni/eco-goinfra/pkg/secret"
	"github.com/openshift-kni/eco-goinfra/pkg/service"
	"github.com/openshift-kni/eco-goinfra/pkg/storage"
	//"io"
	corev1 "k8s.io/api/core/v1"
)

const (
	apiPath        = "/api/logs/v1/"
	queryPath      = "/loki/api/v1/query"
	queryRangePath = "/loki/api/v1/query_range"
)

// CheckODF check if the ODF is installed in the cluster or not
// here only checks the storageclasses and svc/s3
func CheckODF(apiClient *clients.Settings) error {
	expectedSC := []string{
		"openshift-storage.noobaa.io",
		"ocs-external-storagecluster-ceph-rbd",
		"ocs-external-storagecluster-cephfs"}

	for _, scName := range expectedSC {
		glog.V(100).Infof("Checking if %s exists", scName)

		_, err := storage.PullClass(apiClient, scName)

		if err != nil {
			glog.V(100).Infof("Error finding %s storageclass; %v", scName, err)

			return fmt.Errorf("error finding %s storageclass; %w", scName, err)
		}

		glog.V(100).Infof("Found %s storageclass", scName)
	}

	_, err := service.Pull(apiClient, "s3", "openshift-storage")

	if err != nil {
		glog.V(100).Infof("Error pulling s3 service; %v", err)

		return fmt.Errorf("error pulling s3 service; %w", err)
	}

	return nil
}

type lokiClient struct {
	address     string //Server address.
	bearerToken string //adds the Authorization header to API requests for authentication purposes.
	retries     int    //How many times to retry each query when getting an error response from Loki.
	queryTags   string //adds X-Query-Tags header to API requests.
	quiet       bool   //Suppress query metadata.
}

// newLokiClient initializes a lokiClient with server address
func newLokiClient(routeAddress string) *lokiClient {
	client := &lokiClient{}
	client.address = routeAddress
	client.retries = 5
	client.quiet = true

	return client
}

// retry sets how many times to retry each query
func (c *lokiClient) retry(retry int) *lokiClient {
	nc := *c
	nc.retries = retry

	return &nc
}

// withToken sets the token used to do query
func (c *lokiClient) withToken(bearerToken string) *lokiClient {
	nc := *c
	nc.bearerToken = bearerToken

	return &nc
}

//func (c *lokiClient) getHTTPRequestHeader() (http.Header, error) {
//	h := make(http.Header)
//
//	h.Set("User-Agent", "loki-logcli")
//
//	if c.queryTags != "" {
//		h.Set("X-Query-Tags", c.queryTags)
//	}
//
//	if c.bearerToken != "" {
//		h.Set("Authorization", "Bearer "+c.bearerToken)
//	}
//
//	return h, nil
//}

//func (c *lokiClient) doRequest(path, query string, out interface{}) error {
//	h, err := c.getHTTPRequestHeader()
//	if err != nil {
//		return err
//	}
//
//	resp, err := doHTTPRequest(h, c.address, path, query, "GET", c.quiet, c.retries, nil, 200)
//	if err != nil {
//		return err
//	}
//	return json.Unmarshal(resp, out)
//}

//func (c *lokiClient) doQuery(path string, query string) (*lokiQueryResponse, error) {
//	var err error
//	var r lokiQueryResponse
//
//	if err = c.doRequest(path, query, &r); err != nil {
//		return nil, err
//	}
//
//	return &r, nil
//}
//
//// query uses the /api/v1/query endpoint to execute an instant query
//// lc.query("application", "sum by(kubernetes_namespace_name)(count_over_time({kubernetes_namespace_name=\"multiple-containers\"}[5m]))", 30, false, time.Now())
//func (c *lokiClient) query(tenant string, queryStr string, limit int, forward bool, time time.Time) (*lokiQueryResponse, error) {
//	direction := func() string {
//		if forward {
//			return "FORWARD"
//		}
//		return "BACKWARD"
//	}
//	qsb := newQueryStringBuilder()
//	qsb.setString("query", queryStr)
//	qsb.setInt("limit", int64(limit))
//	qsb.setInt("time", time.UnixNano())
//	qsb.setString("direction", direction())
//	var logPath string
//	if len(tenant) > 0 {
//		logPath = apiPath + tenant + queryRangePath
//	} else {
//		logPath = queryRangePath
//	}
//	return c.doQuery(logPath, qsb.encode())
//}
//
//// queryRange uses the /api/v1/query_range endpoint to execute a range query
//// tenant: application, infrastructure, audit
//// queryStr: string to filter logs, for example: "{kubernetes_namespace_name="test"}"
//// limit: max log count
//// start: Start looking for logs at this absolute time(inclusive), e.g.: time.Now().Add(time.Duration(-1)*time.Hour) means 1 hour ago
//// end: Stop looking for logs at this absolute time (exclusive)
//// forward: true means scan forwards through logs, false means scan backwards through logs
//func (c *lokiClient) queryRange(tenant string, queryStr string, limit int, start, end time.Time, forward bool) (*lokiQueryResponse, error) {
//	direction := func() string {
//		if forward {
//			return "FORWARD"
//		}
//		return "BACKWARD"
//	}
//	params := newQueryStringBuilder()
//	params.setString("query", queryStr)
//	params.setInt32("limit", limit)
//	params.setInt("start", start.UnixNano())
//	params.setInt("end", end.UnixNano())
//	params.setString("direction", direction())
//	var logPath string
//	if len(tenant) > 0 {
//		logPath = apiPath + tenant + queryRangePath
//	} else {
//		logPath = queryRangePath
//	}
//
//	return c.doQuery(logPath, params.encode())
//}
//
//func (c *lokiClient) searchLogsInLoki(tenant, query string) (*lokiQueryResponse, error) {
//	res, err := c.queryRange(tenant, query, 5, time.Now().Add(time.Duration(-1)*time.Hour), time.Now(), false)
//	return res, err
//}
//
//func (c *lokiClient) waitForLogsAppearByQuery(tenant, query string) error {
//	return wait.PollUntilContextTimeout(context.Background(), 10*time.Second, 300*time.Second, true, func(context.Context) (done bool, err error) {
//		logs, err := c.searchLogsInLoki(tenant, query)
//		if err != nil {
//			glog.V(100).Infof("\ngot err when searching logs: %v, retrying...\n", err)
//			return false, nil
//		}
//		if len(logs.Data.Result) > 0 {
//			glog.V(100).Infof(`find logs by %s`, query)
//			return true, nil
//		}
//		return false, nil
//	})
//}
//
//func (c *lokiClient) searchByKey(tenant, key, value string) (*lokiQueryResponse, error) {
//	res, err := c.searchLogsInLoki(tenant, "{"+key+"=\""+value+"\"}")
//	return res, err
//}
//
//func (c *lokiClient) waitForLogsAppearByKey(tenant, key, value string) {
//	err := wait.PollUntilContextTimeout(context.Background(), 10*time.Second, 300*time.Second, true, func(context.Context) (done bool, err error) {
//		logs, err := c.searchByKey(tenant, key, value)
//		if err != nil {
//			glog.V(100).Infof("\ngot err when searching logs: %v, retrying...\n", err)
//			return false, nil
//		}
//		if len(logs.Data.Result) > 0 {
//			glog.V(100).Infof(`find logs by {%s="%s"}`, key, value)
//			return true, nil
//		}
//		return false, nil
//	})
//	assertWaitPollNoErr(err, fmt.Sprintf(`can't find logs by {%s="%s"} in last 5 minutes`, key, value))
//}
//
//func assertWaitPollNoErr(e error, msg string) {
//	if e == nil {
//		return
//	}
//	var err error
//	if strings.Compare(e.Error(), "timed out waiting for the condition") == 0 || strings.Compare(e.Error(), "context deadline exceeded") == 0 {
//		err = fmt.Errorf("case: %v\nerror: %s", CurrentSpecReport().FullText(), msg)
//	} else {
//		err = fmt.Errorf("case: %v\nerror: %s", CurrentSpecReport().FullText(), e.Error())
//	}
//	Expect(err).NotTo(HaveOccurred())
//
//}
//func (c *lokiClient) searchByNamespace(tenant, projectName string) (*lokiQueryResponse, error) {
//	res, err := c.searchLogsInLoki(tenant, "{kubernetes_namespace_name=\""+projectName+"\"}")
//	return res, err
//}

// extractLogEntities extract the log entities from loki query response, designed for checking the content of log data in Loki
//func extractLogEntities(lokiQueryResult *lokiQueryResponse) []LogEntity {
//	var lokiLogs []LogEntity
//
//	for _, res := range lokiQueryResult.Data.Result {
//		for _, value := range res.Values {
//			lokiLog := LogEntity{}
//			// only process log data, drop timestamp
//			json.Unmarshal([]byte(convertInterfaceToArray(value)[1]), &lokiLog)
//			lokiLogs = append(lokiLogs, lokiLog)
//		}
//	}
//
//	return lokiLogs
//}
//
//func convertInterfaceToArray(t interface{}) []string {
//	var data []string
//	switch reflect.TypeOf(t).Kind() {
//		case reflect.Slice, reflect.Array:
//			s := reflect.ValueOf(t)
//			for i := 0; i < s.Len(); i++ {
//				data = append(data, fmt.Sprint(s.Index(i)))
//			}
//		}
//
//	return data
//}

//type queryStringBuilder struct {
//	values url.Values
//}
//
//func newQueryStringBuilder() *queryStringBuilder {
//	return &queryStringBuilder{
//		values: url.Values{},
//	}
//}
//
//func (b *queryStringBuilder) setString(name, value string) {
//	b.values.Set(name, value)
//}
//
//func (b *queryStringBuilder) setInt(name string, value int64) {
//	b.setString(name, strconv.FormatInt(value, 10))
//}
//
//func (b *queryStringBuilder) setInt32(name string, value int) {
//	b.setString(name, strconv.Itoa(value))
//}
//
//// encode returns the URL-encoded query string based on key-value
//// parameters added to the builder calling Set functions.
//func (b *queryStringBuilder) encode() string {
//	return b.values.Encode()
//}

//
//// queryAlertManagerForLokiAlerts() queries user-workload alert-manager if isUserWorkloadAM parameter is true.
//// All active alerts should be returned when querying Alert Managers
//func queryAlertManagerForActiveAlerts(oc *exutil.CLI, token string, isUserWorkloadAM bool, alertName string, timeInMinutes int) {
//	var err error
//	if !isUserWorkloadAM {
//		alertManagerRoute := getRouteAddress(oc, "openshift-monitoring", "alertmanager-main")
//		h := make(http.Header)
//		h.Add("Content-Type", "application/json")
//		h.Add("Authorization", "Bearer "+token)
//		params := url.Values{}
//		err = wait.PollUntilContextTimeout(context.Background(), 30*time.Second, time.Duration(timeInMinutes)*time.Minute, true, func(context.Context) (done bool, err error) {
//			resp, err := doHTTPRequest(h, "https://"+alertManagerRoute, "/api/v2/alerts", params.Encode(), "GET", true, 5, nil, 200)
//			if err != nil {
//				return false, err
//			}
//			if strings.Contains(string(resp), alertName) {
//				return true, nil
//			}
//			glog.V(100).Infof("Waiting for alert %s to be in Firing state", alertName)
//			return false, nil
//		})
//
//	} else {
//		userWorkloadAlertManagerURL := "https://alertmanager-user-workload.openshift-user-workload-monitoring.svc:9095/api/v2/alerts"
//		authBearer := " \"Authorization: Bearer " + token + "\""
//		cmd := "curl -k -H" + authBearer + " " + userWorkloadAlertManagerURL
//		err = wait.PollUntilContextTimeout(context.Background(), 30*time.Second, time.Duration(timeInMinutes)*time.Minute, true, func(context.Context) (done bool, err error) {
//			alerts, err := exutil.RemoteShPod(oc, "openshift-monitoring", "prometheus-k8s-0", "/bin/sh", "-x", "-c", cmd)
//			if err != nil {
//				return false, err
//			}
//			if strings.Contains(string(alerts), alertName) {
//				return true, nil
//			}
//			glog.V(100).Infof("Waiting for alert %s to be in Firing state", alertName)
//			return false, nil
//		})
//	}
//
//	assertWaitPollNoErr(err, fmt.Sprintf("Alert %s is not firing after %d minutes", alertName, timeInMinutes))
//}
//
//func doHTTPRequest(header http.Header, address, path, query, method string, quiet bool, attempts int, requestBody io.Reader, expectedStatusCode int) ([]byte, error) {
//	us, err := buildURL(address, path, query)
//	if err != nil {
//		return nil, err
//	}
//	if !quiet {
//		glog.V(100).Infof("the URL is: %s", us)
//	}
//
//	req, err := http.NewRequest(strings.ToUpper(method), us, requestBody)
//	if err != nil {
//		return nil, err
//	}
//
//	req.Header = header
//
//	tr := &http.Transport{
//		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
//	}
//
//	client := &http.Client{Transport: tr}
//
//	var resp *http.Response
//	success := false
//
//	for attempts > 0 {
//		attempts--
//
//		resp, err = client.Do(req)
//		if err != nil {
//			glog.V(100).Infof("error sending request %v", err)
//			continue
//		}
//		if resp.StatusCode != expectedStatusCode {
//			buf, _ := io.ReadAll(resp.Body) // nolint
//			glog.V(100).Infof("Error response from server: %s %s (%v), attempts remaining: %d", resp.Status, string(buf), err, attempts)
//			if err := resp.Body.Close(); err != nil {
//				glog.V(100).Infof("error closing body: %v", err)
//			}
//			// sleep 5 second before doing next request
//			time.Sleep(5 * time.Second)
//			continue
//		}
//		success = true
//		break
//	}
//	if !success {
//		return nil, fmt.Errorf("run out of attempts while querying the server")
//	}
//
//	defer func() {
//		if err := resp.Body.Close(); err != nil {
//			glog.V(100).Infof("error closing body: %v", err)
//		}
//	}()
//	return io.ReadAll(resp.Body)
//}
//
//// buildURL concats a url `http://foo/bar` with a path `/buzz`.
//func buildURL(u, p, q string) (string, error) {
//	url, err := url.Parse(u)
//	if err != nil {
//		return "", err
//	}
//	url.Path = path.Join(url.Path, p)
//	url.RawQuery = q
//	return url.String(), nil
//}

func getLokiBucketData(apiClient *clients.Settings, configmapName, nsName string) (string, string, string, error) {
	glog.V(100).Infof("Get Loki bucket configmap %s in namespace %s object", configmapName, nsName)

	bucketHostKey := "BUCKET_HOST"
	bucketNameKey := "BUCKET_NAME"
	bucketPortKey := "BUCKET_PORT"

	lokiBucketConfigmap, err := configmap.Pull(apiClient, configmapName, nsName)

	if err != nil {
		glog.V(100).Infof("Error pulling Loki bucket configmap %s from the namespace %s: %v",
			configmapName, nsName, err)

		return "", "", "", fmt.Errorf("error pulling Loki bucket configmap %s from the namespace %s: %w",
			configmapName, nsName, err)
	}

	glog.V(100).Infof("Get Loki bucket host and name from configmap %s in namespace %s",
		configmapName, nsName)

	cmData := lokiBucketConfigmap.Object.Data

	if len(cmData) == 0 {
		glog.V(100).Infof("Loki bucket configmap %s in namespace %s has no data",
			configmapName, nsName)

		return "", "", "", fmt.Errorf("loki bucket configmap %s in namespace %s has no data",
			configmapName, nsName)
	}

	glog.V(100).Infof("cmData: %s", cmData)

	bucketHost := cmData[bucketHostKey]

	if bucketHost == "" {
		glog.V(100).Infof("%s value not found in the configmap %s in namespace %s",
			bucketHostKey, configmapName, nsName)

		return "", "", "", fmt.Errorf("%s value not found in the configmap %s in namespace %s",
			bucketHostKey, configmapName, nsName)
	}

	glog.V(100).Infof("bucketHost: %s", bucketHost)

	bucketName := cmData[bucketNameKey]

	if bucketName == "" {
		glog.V(100).Infof("%s value not found in the configmap %s in namespace %s",
			bucketNameKey, configmapName, nsName)

		return "", "", "", fmt.Errorf("%s value not found in the configmap %s in namespace %s",
			bucketNameKey, configmapName, nsName)
	}

	glog.V(100).Infof("bucketName: %s", bucketName)

	bucketPort := cmData[bucketPortKey]

	if bucketPort == "" {
		glog.V(100).Infof("%s value not found in the configmap %s in namespace %s",
			bucketPortKey, configmapName, nsName)

		return "", "", "", fmt.Errorf("%s value not found in the configmap %s in namespace %s",
			bucketPortKey, configmapName, nsName)
	}

	glog.V(100).Infof("bucketPort: %s", bucketPort)

	return bucketHost, bucketName, bucketPort, nil
}

func getLokiBucketAccessKeyData(apiClient *clients.Settings, secretName, nsName string) (string, string, error) {
	glog.V(100).Infof("Get Loki bucket secret %s in namespace %s object", secretName, nsName)

	accessKeyIDKey := "AWS_ACCESS_KEY_ID"
	secretAccessKeyKey := "AWS_SECRET_ACCESS_KEY"

	lokiBucketSecret, err := secret.Pull(apiClient, secretName, nsName)

	if err != nil {
		glog.V(100).Infof("Error pulling Loki bucket secret %s from the namespace %s: %v",
			secretName, nsName, err)

		return "", "", fmt.Errorf("error pulling Loki bucket secret %s from the namespace %s: %w",
			secretName, nsName, err)
	}

	glog.V(100).Infof("Get Loki bucket accessKeyID and secretAccessKey from secret %s in namespace %s",
		secretName, nsName)

	secretData := lokiBucketSecret.Object.Data

	if len(secretData) == 0 {
		glog.V(100).Infof("Loki bucket secret %s in namespace %s has no data",
			secretName, nsName)

		return "", "", fmt.Errorf("loki bucket secret %s in namespace %s has no data", secretName, nsName)
	}

	glog.V(100).Infof("secretData: %s", secretData)

	accessKeyID := string(secretData[accessKeyIDKey])

	if accessKeyID == "" {
		glog.V(100).Infof("%s value not found in the secret %s in namespace %s",
			accessKeyIDKey, secretName, nsName)

		return "", "", fmt.Errorf("%s value not found in the secret %s in namespace %s",
			accessKeyIDKey, secretName, nsName)
	}

	glog.V(100).Infof("accessKeyID: %s", accessKeyID)

	secretAccessKey := string(secretData[secretAccessKeyKey])

	if secretAccessKey == "" {
		glog.V(100).Infof("%s value not found in the secret %s in namespace %s",
			secretAccessKeyKey, secretName, nsName)

		return "", "", fmt.Errorf("%s value not found in the secret %s in namespace %s",
			secretAccessKeyKey, secretName, nsName)
	}

	glog.V(100).Infof("secretAccessKey: %s", secretAccessKey)

	return accessKeyID, secretAccessKey, nil
}

// CreateLokiStackSecret creating lokiStack secret object.
func CreateLokiStackSecret(apiClient *clients.Settings, secretName, objectBucketClaim, nsName string) error {
	lokiSecretObj := secret.NewBuilder(apiClient, secretName, nsName, corev1.SecretTypeOpaque)

	if lokiSecretObj.Exists() {
		err := lokiSecretObj.Delete()

		if err != nil {
			glog.V(100).Infof("Failed to delete loki secret %s from namespace %s; %v",
				secretName, nsName, err)

			return fmt.Errorf("failed to delete loki secret %s from namespace %s; %w",
				secretName, nsName, err)
		}
	}

	glog.V(100).Infof("Create loki secret %s in namespace %s", secretName, nsName)

	bucketHost, bucketName, bucketPort, err := getLokiBucketData(apiClient, objectBucketClaim, nsName)

	if err != nil {
		glog.V(100).Infof("Error fetching Loki bucket data from configmap %s in namespace %s: %v",
			bucketName, nsName, err)

		return fmt.Errorf("error fetching Loki bucket data from configmap %s in namespace %s: %w",
			bucketName, nsName, err)
	}

	glog.V(100).Infof("DATA: %s, %s, %s", bucketHost, bucketName, bucketPort)

	accessKeyID, secretAccessKey, err := getLokiBucketAccessKeyData(
		apiClient, objectBucketClaim, nsName)

	if err != nil {
		glog.V(100).Infof("Error fetching Loki bucket access key data from secret %s in namespace %s: %v",
			bucketName, nsName, err)

		return fmt.Errorf("error fetching Loki bucket access key data from secret %s in namespace %s: %w",
			bucketName, nsName, err)
	}

	glog.V(100).Infof("SECRET DATA: %s, %s", accessKeyID, secretAccessKey)

	stringData := map[string]string{
		"bucketnames":       bucketName,
		"endpoint":          fmt.Sprintf("https://%s:%s", bucketHost, bucketPort),
		"region":            "",
		"access_key_id":     accessKeyID,
		"access_key_secret": secretAccessKey,
	}

	_, err = lokiSecretObj.WithStringData(stringData).Create()

	if err != nil {
		glog.V(100).Infof("Error creating lokiStack secret %s in namespace %s: %v",
			secretName, nsName, err)

		return fmt.Errorf("error creating lokiStack secret %s in namespace %s: %w",
			secretName, nsName, err)
	}

	return nil
}
