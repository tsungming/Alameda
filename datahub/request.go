package datahub

import (
	"time"

	"github.com/containers-ai/alameda/datahub/pkg/dao"
	prediction_dao "github.com/containers-ai/alameda/datahub/pkg/dao/prediction"
	score_dao "github.com/containers-ai/alameda/datahub/pkg/dao/score"
	"github.com/containers-ai/alameda/datahub/pkg/metric"
	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	"github.com/golang/protobuf/ptypes"
)

type datahubListPodMetricsRequestExtended struct {
	datahub_v1alpha1.ListPodMetricsRequest
}

func (r datahubListPodMetricsRequestExtended) validate() error {
	return nil
}

type datahubListNodeMetricsRequestExtended struct {
	datahub_v1alpha1.ListNodeMetricsRequest
}

func (r datahubListNodeMetricsRequestExtended) validate() error {
	return nil
}

type datahubCreatePodPredictionsRequestExtended struct {
	datahub_v1alpha1.CreatePodPredictionsRequest
}

func (r datahubCreatePodPredictionsRequestExtended) validate() error {
	return nil
}

func (r datahubCreatePodPredictionsRequestExtended) daoContainerPredictions() []*prediction_dao.ContainerPrediction {

	var (
		containerPredictions []*prediction_dao.ContainerPrediction
	)

	for _, datahubPodPrediction := range r.PodPredictions {

		podNamespace := ""
		podName := ""
		if datahubPodPrediction.GetNamespacedName() != nil {
			podNamespace = datahubPodPrediction.GetNamespacedName().GetNamespace()
			podName = datahubPodPrediction.GetNamespacedName().GetName()
		}

		for _, datahubContainerPrediction := range datahubPodPrediction.GetContainerPredictions() {
			containerName := datahubContainerPrediction.GetName()

			for _, rawData := range datahubContainerPrediction.GetPredictedRawData() {

				containerPrediction := prediction_dao.ContainerPrediction{
					Namespace:     podNamespace,
					PodName:       podName,
					ContainerName: containerName,
					Predictions:   make(map[metric.ContainerMetricType][]metric.Sample),
				}

				samples := []metric.Sample{}
				for _, datahubSample := range rawData.GetData() {
					time, err := ptypes.Timestamp(datahubSample.GetTime())
					if err != nil {
						scope.Error(" failed: " + err.Error())
					}
					sample := metric.Sample{
						Timestamp: time,
						Value:     datahubSample.GetNumValue(),
					}
					samples = append(samples, sample)
				}

				var metricType metric.ContainerMetricType
				switch rawData.GetMetricType() {
				case datahub_v1alpha1.MetricType_CPU_USAGE_SECONDS_PERCENTAGE:
					metricType = metric.TypeContainerCPUUsageSecondsPercentage
				case datahub_v1alpha1.MetricType_MEMORY_USAGE_BYTES:
					metricType = metric.TypeContainerMemoryUsageBytes
				}
				containerPrediction.Predictions[metricType] = samples

				containerPredictions = append(containerPredictions, &containerPrediction)
			}
		}
	}

	return containerPredictions
}

type datahubCreateNodePredictionsRequestExtended struct {
	datahub_v1alpha1.CreateNodePredictionsRequest
}

func (r datahubCreateNodePredictionsRequestExtended) validate() error {
	return nil
}

func (r datahubCreateNodePredictionsRequestExtended) daoNodePredictions() []*prediction_dao.NodePrediction {

	var (
		NodePredictions []*prediction_dao.NodePrediction
	)

	for _, datahubNodePrediction := range r.NodePredictions {

		nodeName := datahubNodePrediction.GetName()
		isScheduled := datahubNodePrediction.GetIsScheduled()

		for _, rawData := range datahubNodePrediction.GetPredictedRawData() {

			samples := []metric.Sample{}
			for _, datahubSample := range rawData.GetData() {
				time, err := ptypes.Timestamp(datahubSample.GetTime())
				if err != nil {
					scope.Error(" failed: " + err.Error())
				}
				sample := metric.Sample{
					Timestamp: time,
					Value:     datahubSample.GetNumValue(),
				}
				samples = append(samples, sample)
			}

			NodePrediction := prediction_dao.NodePrediction{
				NodeName:    nodeName,
				IsScheduled: isScheduled,
				Predictions: make(map[metric.NodeMetricType][]metric.Sample),
			}

			var metricType metric.ContainerMetricType
			switch rawData.GetMetricType() {
			case datahub_v1alpha1.MetricType_CPU_USAGE_SECONDS_PERCENTAGE:
				metricType = metric.TypeNodeCPUUsageSecondsPercentage
			case datahub_v1alpha1.MetricType_MEMORY_USAGE_BYTES:
				metricType = metric.TypeNodeMemoryUsageBytes
			}
			NodePrediction.Predictions[metricType] = samples

			NodePredictions = append(NodePredictions, &NodePrediction)
		}
	}

	return NodePredictions
}

type datahubListPodPredictionsRequestExtended struct {
	request *datahub_v1alpha1.ListPodPredictionsRequest
}

func (r datahubListPodPredictionsRequestExtended) daoListPodPredictionsRequest() prediction_dao.ListPodPredictionsRequest {

	var (
		namespace      string
		podName        string
		queryCondition dao.QueryCondition
	)

	if r.request.GetNamespacedName() != nil {
		namespace = r.request.GetNamespacedName().GetNamespace()
		podName = r.request.GetNamespacedName().GetName()
	}

	queryCondition = datahubQueryConditionExtend{r.request.GetQueryCondition()}.daoQueryCondition()
	listContainerPredictionsRequest := prediction_dao.ListPodPredictionsRequest{
		Namespace:      namespace,
		PodName:        podName,
		QueryCondition: queryCondition,
	}

	return listContainerPredictionsRequest
}

type datahubListNodePredictionsRequestExtended struct {
	request *datahub_v1alpha1.ListNodePredictionsRequest
}

func (r datahubListNodePredictionsRequestExtended) daoListNodePredictionsRequest() prediction_dao.ListNodePredictionsRequest {

	var (
		nodeNames      []string
		queryCondition dao.QueryCondition
	)

	for _, nodeName := range r.request.GetNodeNames() {
		nodeNames = append(nodeNames, nodeName)
	}
	queryCondition = datahubQueryConditionExtend{r.request.GetQueryCondition()}.daoQueryCondition()
	listNodePredictionsRequest := prediction_dao.ListNodePredictionsRequest{
		NodeNames:      nodeNames,
		QueryCondition: queryCondition,
	}

	return listNodePredictionsRequest
}

type datahubListSimulatedSchedulingScoresRequestExtended struct {
	request *datahub_v1alpha1.ListSimulatedSchedulingScoresRequest
}

func (r datahubListSimulatedSchedulingScoresRequestExtended) daoLisRequest() score_dao.ListRequest {

	var (
		queryCondition dao.QueryCondition
	)

	queryCondition = datahubQueryConditionExtend{r.request.GetQueryCondition()}.daoQueryCondition()
	listRequest := score_dao.ListRequest{
		queryCondition,
	}

	return listRequest
}

type datahubQueryConditionExtend struct {
	queryCondition *datahub_v1alpha1.QueryCondition
}

func (d datahubQueryConditionExtend) metricDAOQueryCondition() dao.QueryCondition {

	var (
		queryStartTime      *time.Time
		queryEndTime        *time.Time
		queryStepTime       *time.Duration
		queryTimestampOrder int
		queryLimit          int
		queryCondition      = dao.QueryCondition{}
	)

	if d.queryCondition == nil {
		return queryCondition
	}

	if d.queryCondition.GetTimeRange() != nil {
		timeRange := d.queryCondition.GetTimeRange()
		if timeRange.GetStartTime() != nil {
			tmpTime, _ := ptypes.Timestamp(timeRange.GetStartTime())
			queryStartTime = &tmpTime
		}
		if timeRange.GetEndTime() != nil {
			tmpTime, _ := ptypes.Timestamp(timeRange.GetEndTime())
			queryEndTime = &tmpTime
		}
		if timeRange.GetStep() != nil {
			tmpTime, _ := ptypes.Duration(timeRange.GetStep())
			queryStepTime = &tmpTime
		}

		switch d.queryCondition.GetOrder() {
		case datahub_v1alpha1.QueryCondition_ASC:
			queryTimestampOrder = dao.Asc
		case datahub_v1alpha1.QueryCondition_DESC:
			queryTimestampOrder = dao.Desc
		default:
			queryTimestampOrder = dao.Asc
		}

		queryLimit = int(d.queryCondition.GetLimit())
	}
	queryTimestampOrder = int(d.queryCondition.GetOrder())
	queryLimit = int(d.queryCondition.GetLimit())

	queryCondition = dao.QueryCondition{
		StartTime:      queryStartTime,
		EndTime:        queryEndTime,
		StepTime:       queryStepTime,
		TimestampOrder: queryTimestampOrder,
		Limit:          queryLimit,
	}
	return queryCondition
}

func (d datahubQueryConditionExtend) daoQueryCondition() dao.QueryCondition {

	var (
		queryStartTime      *time.Time
		queryEndTime        *time.Time
		queryStepTime       *time.Duration
		queryTimestampOrder int
		queryLimit          int
		queryCondition      = dao.QueryCondition{}
	)

	if d.queryCondition == nil {
		return queryCondition
	}

	if d.queryCondition.GetTimeRange() != nil {
		timeRange := d.queryCondition.GetTimeRange()
		if timeRange.GetStartTime() != nil {
			tmpTime, _ := ptypes.Timestamp(timeRange.GetStartTime())
			queryStartTime = &tmpTime
		}
		if timeRange.GetEndTime() != nil {
			tmpTime, _ := ptypes.Timestamp(timeRange.GetEndTime())
			queryEndTime = &tmpTime
		}
		if timeRange.GetStep() != nil {
			tmpTime, _ := ptypes.Duration(timeRange.GetStep())
			queryStepTime = &tmpTime
		}

		switch d.queryCondition.GetOrder() {
		case datahub_v1alpha1.QueryCondition_ASC:
			queryTimestampOrder = dao.Asc
		case datahub_v1alpha1.QueryCondition_DESC:
			queryTimestampOrder = dao.Desc
		default:
			queryTimestampOrder = dao.Asc
		}

		queryLimit = int(d.queryCondition.GetLimit())
	}
	queryTimestampOrder = int(d.queryCondition.GetOrder())
	queryLimit = int(d.queryCondition.GetLimit())

	queryCondition = dao.QueryCondition{
		StartTime:      queryStartTime,
		EndTime:        queryEndTime,
		StepTime:       queryStepTime,
		TimestampOrder: queryTimestampOrder,
		Limit:          queryLimit,
	}
	return queryCondition
}
