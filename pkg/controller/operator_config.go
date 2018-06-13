package controller

import (
	"encoding/json"
	"fmt"

	"github.com/zalando-incubator/postgres-operator/pkg/util/config"
	"github.com/zalando-incubator/postgres-operator/pkg/util/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Controller) readOperatorConfigurationFromCRD(configObjectName string) (*config.OperatorConfiguration, error) {
	var (
		config config.OperatorConfiguration
	)

	req := c.KubeClient.CRDREST.Get().
		Name(configObjectName).
		Namespace(c.opConfig.WatchedNamespace).
		Resource(constants.OperatorConfigCRDResource).
		VersionedParams(&metav1.ListOptions{ResourceVersion: "0"}, metav1.ParameterCodec)

	data, err := req.DoRaw()
	if err != nil {
		return nil, fmt.Errorf("could not get operator configuration object %s: %v", configObjectName, err)
	}
	if err = json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("could not unmarshal operator configuration object %s, %v", configObjectName, err)
	}

	return &config, nil
}

// importConfigurationFromCRD is a transitional function that converts CRD configuration to the one based on the configmap
func (c *Controller) importConfigurationFromCRD(fromCRD *config.OperatorConfigurationData) *config.Config {
	result := &config.Config{}

	result.EtcdHost = fromCRD.EtcdHost
	result.DockerImage = fromCRD.DockerImage
	result.Workers = fromCRD.Workers
	result.MinInstances = fromCRD.MinInstances
	result.MaxInstances = fromCRD.MaxInstances
	result.ResyncPeriod = fromCRD.ResyncPeriod

	result.SuperUsername = fromCRD.PostgresUsersConfiguration.SuperUsername
	result.ReplicationUsername = fromCRD.PostgresUsersConfiguration.ReplicationUsername

	result.PodServiceAccountName = fromCRD.Kubernetes.PodServiceAccountName
	result.PodServiceAccountDefinition = fromCRD.Kubernetes.PodServiceAccountDefinition
	result.PodTerminateGracePeriod = fromCRD.Kubernetes.PodTerminateGracePeriod
	result.WatchedNamespace = fromCRD.Kubernetes.WatchedNamespace
	result.PDBNameFormat = fromCRD.Kubernetes.PDBNameFormat
	result.SecretNameTemplate = fromCRD.Kubernetes.SecretNameTemplate
	result.OAuthTokenSecretName = fromCRD.Kubernetes.OAuthTokenSecretName
	result.InfrastructureRolesSecretName = fromCRD.Kubernetes.InfrastructureRolesSecretName
	result.PodRoleLabel = fromCRD.Kubernetes.PodRoleLabel
	result.ClusterLabels = fromCRD.Kubernetes.ClusterLabels
	result.ClusterNameLabel = fromCRD.Kubernetes.ClusterNameLabel
	result.NodeReadinessLabel = fromCRD.Kubernetes.NodeReadinessLabel

	result.DefaultCPURequest = fromCRD.PostgresPodResources.DefaultCPURequest
	result.DefaultMemoryRequest = fromCRD.PostgresPodResources.DefaultMemoryRequest
	result.DefaultCPULimit = fromCRD.PostgresPodResources.DefaultCPULimit
	result.DefaultMemoryLimit = fromCRD.PostgresPodResources.DefaultMemoryLimit

	result.ResourceCheckInterval = fromCRD.Timeouts.ResourceCheckInterval
	result.ResourceCheckTimeout = fromCRD.Timeouts.ResourceCheckTimeout
	result.PodLabelWaitTimeout = fromCRD.Timeouts.PodLabelWaitTimeout
	result.PodDeletionWaitTimeout = fromCRD.Timeouts.PodDeletionWaitTimeout
	result.ReadyWaitInterval = fromCRD.Timeouts.ReadyWaitInterval
	result.ReadyWaitTimeout = fromCRD.Timeouts.ReadyWaitTimeout

	result.DbHostedZone = fromCRD.LoadBalancer.DbHostedZone
	result.EnableMasterLoadBalancer = fromCRD.LoadBalancer.EnableMasterLoadBalancer
	result.EnableReplicaLoadBalancer = fromCRD.LoadBalancer.EnableReplicaLoadBalancer
	result.MasterDNSNameFormat = fromCRD.LoadBalancer.MasterDNSNameFormat
	result.ReplicaDNSNameFormat = fromCRD.LoadBalancer.ReplicaDNSNameFormat

	result.WALES3Bucket = fromCRD.AWSGCP.WALES3Bucket
	result.LogS3Bucket = fromCRD.AWSGCP.LogS3Bucket
	result.KubeIAMRole = fromCRD.AWSGCP.KubeIAMRole

	result.DebugLogging = fromCRD.OperatorDebug.DebugLogging
	result.EnableDBAccess = fromCRD.OperatorDebug.EnableDBAccess
	result.EnableTeamsAPI = fromCRD.TeamsAPI.EnableTeamsAPI
	result.TeamsAPIUrl = fromCRD.TeamsAPI.TeamsAPIUrl
	result.TeamAPIRoleConfiguration = fromCRD.TeamsAPI.TeamAPIRoleConfiguration
	result.EnableTeamSuperuser = fromCRD.TeamsAPI.EnableTeamSuperuser
	result.TeamAdminRole = fromCRD.TeamsAPI.TeamAdminRole
	result.PamRoleName = fromCRD.TeamsAPI.PamRoleName

	result.APIPort = fromCRD.LoggingRESTAPI.APIPort
	result.RingLogLines = fromCRD.LoggingRESTAPI.RingLogLines
	result.ClusterHistoryEntries = fromCRD.LoggingRESTAPI.ClusterHistoryEntries

	result.ScalyrAPIKey = fromCRD.Scalyr.ScalyrAPIKey
	result.ScalyrImage = fromCRD.Scalyr.ScalyrImage
	result.ScalyrServerURL = fromCRD.Scalyr.ScalyrServerURL
	result.ScalyrCPURequest = fromCRD.Scalyr.ScalyrCPURequest
	result.ScalyrMemoryRequest = fromCRD.Scalyr.ScalyrMemoryRequest
	result.ScalyrCPULimit = fromCRD.Scalyr.ScalyrCPULimit
	result.ScalyrMemoryLimit = fromCRD.Scalyr.ScalyrMemoryLimit

	return result
}
