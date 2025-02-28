package clusterstatus

import (
	"fmt"
	"time"

	datahub_v1alpha1 "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"

	cluster_status_entity "github.com/containers-ai/alameda/datahub/pkg/entity/influxdb/cluster_status"
	"github.com/containers-ai/alameda/datahub/pkg/repository/influxdb"
	"github.com/containers-ai/alameda/datahub/pkg/utils"
	"github.com/containers-ai/alameda/pkg/utils/log"
	influxdb_client "github.com/influxdata/influxdb/client/v2"
)

var (
	scope = log.RegisterScope("influxdb_repo_node_measurement", "InfluxDB repository node measurement", 0)
)

type NodeRepository struct {
	influxDB *influxdb.InfluxDBRepository
}

func (nodeRepository *NodeRepository) IsTag(column string) bool {
	for _, tag := range cluster_status_entity.NodeTags {
		if column == string(tag) {
			return true
		}
	}
	return false
}

func NewNodeRepository(influxDBCfg *influxdb.Config) *NodeRepository {
	return &NodeRepository{
		influxDB: &influxdb.InfluxDBRepository{
			Address:  influxDBCfg.Address,
			Username: influxDBCfg.Username,
			Password: influxDBCfg.Password,
		},
	}
}

// AddAlamedaNodes add node information to database
func (nodeRepository *NodeRepository) AddAlamedaNodes(alamedaNodes []*datahub_v1alpha1.Node) error {
	points := []*influxdb_client.Point{}
	for _, alamedaNode := range alamedaNodes {
		tags := map[string]string{
			string(cluster_status_entity.NodeName): alamedaNode.Name,
		}
		fields := map[string]interface{}{
			string(cluster_status_entity.NodeInCluster): true,
		}
		if pt, err := influxdb_client.NewPoint(string(Node), tags, fields, influxdb.ZeroTime); err == nil {
			points = append(points, pt)
		} else {
			scope.Error(err.Error())
		}
	}
	nodeRepository.influxDB.WritePoints(points, influxdb_client.BatchPointsConfig{
		Database: string(influxdb.ClusterStatus),
	})
	return nil
}

func (nodeRepository *NodeRepository) RemoveAlamedaNodes(alamedaNodes []*datahub_v1alpha1.Node) error {
	points := []*influxdb_client.Point{}
	for _, alamedaNode := range alamedaNodes {
		// SELECT * FROM node WHERE "name"='%s' AND in_cluster=true ORDER BY time ASC LIMIT 1
		cmd := fmt.Sprintf("SELECT * FROM %s WHERE \"%s\"='%s' AND \"%s\"=%s ORDER BY time ASC LIMIT 1",
			string(Node), string(cluster_status_entity.NodeName), alamedaNode.Name,
			string(cluster_status_entity.NodeInCluster), "true")
		if results, _ := nodeRepository.influxDB.QueryDB(cmd, string(influxdb.ClusterStatus)); len(results) == 1 && len(results[0].Series) == 1 {
			newFields := map[string]interface{}{}
			newTags := map[string]string{}
			originTime := time.Now()
			for columnIdx, column := range results[0].Series[0].Columns {
				colVal := results[0].Series[0].Values[0][columnIdx]
				if column == influxdb.Time && colVal != nil {
					originTime, _ = utils.ParseTime(colVal.(string))
				} else if nodeRepository.IsTag(column) && column != influxdb.Time && colVal != nil {
					newTags[column] = colVal.(string)
				} else if !nodeRepository.IsTag(column) {
					if column == string(cluster_status_entity.NodeInCluster) {
						newFields[column] = false
					} else {
						newFields[column] = results[0].Series[0].Values[0][columnIdx]
					}
				}
			}

			if pt, err := influxdb_client.NewPoint(string(Node), newTags, newFields, originTime); err == nil {
				points = append(points, pt)
			} else {
				scope.Error(err.Error())
			}
		}
	}
	nodeRepository.influxDB.WritePoints(points, influxdb_client.BatchPointsConfig{
		Database: string(influxdb.ClusterStatus),
	})
	return nil
}

func (nodeRepository *NodeRepository) ListAlamedaNodes() ([]*influxdb.InfluxDBEntity, error) {
	entities := []*influxdb.InfluxDBEntity{}
	cmd := fmt.Sprintf("SELECT * FROM %s WHERE \"%s\"=%s", string(Node), string(cluster_status_entity.NodeInCluster), "true")
	scope.Infof(fmt.Sprintf("Query nodes in cluster: %s", cmd))
	if results, _ := nodeRepository.influxDB.QueryDB(cmd, string(influxdb.ClusterStatus)); len(results) == 1 && len(results[0].Series) > 0 {
		for _, value := range results[0].Series[0].Values {
			newFields := map[string]interface{}{}
			newTags := map[string]string{}
			entity := influxdb.InfluxDBEntity{}
			for columnIdx, column := range results[0].Series[0].Columns {
				if column == influxdb.Time && value[columnIdx] != nil {
					entity.Time, _ = utils.ParseTime(value[columnIdx].(string))
				} else if nodeRepository.IsTag(column) && value[columnIdx] != nil {
					newTags[column] = value[columnIdx].(string)
				} else if !nodeRepository.IsTag(column) {
					newFields[column] = value[columnIdx]
				}
			}
			entity.Tags = newTags
			entity.Fields = newFields

			entities = append(entities, &entity)
		}
	}
	return entities, nil
}
