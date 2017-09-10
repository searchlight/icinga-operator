package operator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/appscode/go/runtime"
	"github.com/appscode/pat"
	tapi "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/icinga"
)

const (
	PathParamNamespace = ":namespace"
	PathParamType      = ":type"
	PathParamName      = ":name"
)

type AcknowledgeRequest struct {
	ObjectName string
	Author     string
	Comment    string
}

func Acknowledge(client *icinga.Client, w http.ResponseWriter, r *http.Request) {
	defer runtime.HandleCrash()

	params, found := pat.FromContext(r.Context())
	if !found {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	namespace := params.Get(PathParamNamespace)
	if namespace == "" {
		http.Error(w, "Missing parameter "+PathParamNamespace, http.StatusBadRequest)
		return
	}

	alertType := params.Get(PathParamType)
	if alertType == "" {
		http.Error(w, "Missing parameter "+PathParamType, http.StatusBadRequest)
		return
	}

	alertName := params.Get(PathParamName)
	if alertName == "" {
		http.Error(w, "Missing parameter "+PathParamName, http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var ackRequest AcknowledgeRequest
	if err := json.Unmarshal(body, &ackRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	host := &icinga.IcingaHost{
		ObjectName:     ackRequest.ObjectName,
		AlertNamespace: namespace,
	}

	switch alertType {
	case tapi.ResourceTypePodAlert:
		host.Type = icinga.TypePod
	case tapi.ResourceTypeNodeAlert:
		host.Type = icinga.TypeNode
	case tapi.ResourceTypeClusterAlert:
		host.Type = icinga.TypeCluster
	}

	hostName, err := host.Name()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mp := make(map[string]interface{})
	mp["type"] = "Service"
	mp["filter"] = fmt.Sprintf(`service.name == "%s" && host.name == "%s"`, alertName, hostName)
	mp["comment"] = ackRequest.Comment
	mp["notify"] = true
	mp["author"] = ackRequest.Author

	jsonStr, err := json.Marshal(mp)
	if err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}
	resp := client.Actions("acknowledge-problem").Update([]string{}, string(jsonStr)).Do()
	if resp.Status == 200 {
		http.Error(w, "Problem acknowledged", http.StatusOK)
		return
	}

	http.Error(w, "Failed to acknowledge", http.StatusOK)
	return
}
