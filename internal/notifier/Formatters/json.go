package Formatters

import (
	"encoding/json"
	"gold/internal/analyzer"
)

type JSONFormatter struct {
}

type JSONInfo struct {
	ServiceName string `json:"service_name"`
	Timestamp   int64  `json:"timestamp"`

	Goroutines     int `json:"goroutines"`
	Alerts         int `json:"alerts"`
	TotalClusters  int `json:"total_clusters"`
	MaxClusterSize int `json:"max_cluster_size"`

	TotalBlockEvents int `json:"total_block_events"`
	TotalMutexEvents int `json:"total_mutex_events"`

	Hotspots []string `json:"hotspots"`
	// top 3 clusters by count
}
type JSONAlert struct {
	ServiceName string   `json:"service_name"`
	Type        string   `json:"type"`
	Hash        string   `json:"hash"`
	Frames      []string `json:"frames"`
	Description string   `json:"description"`
	Advice      string   `json:"advice"`
	PrevCount   int      `json:"prev_count"`
	NewCount    int      `json:"new_count"`
	Rate        float64  `json:"rate"`
	Cycles      int64    `json:"cycles"`
	Timestamp   int64    `json:"timestamp"`
}

func (f *JSONFormatter) FormatAlert(alert *analyzer.Alert) ([]byte, error) {
	jsonAlert := ConvertToJSONAlert(alert)
	buf, err := json.Marshal(jsonAlert)
	if err != nil {
		return nil, err
	}
	return buf, nil

}
func (f *JSONFormatter) FormatInfo(info *analyzer.Info) ([]byte, error) {
	jsonInfo := ConvertToJSONInfo(info)
	buf, err := json.Marshal(jsonInfo)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func ConvertToJSONAlert(alert *analyzer.Alert) JSONAlert {
	return JSONAlert{
		ServiceName: alert.ServiceName,
		Type:        alert.Type,
		Hash:        alert.Hash,
		Frames:      alert.Frames,
		Description: alert.Description,
		Advice:      alert.Advice,
		PrevCount:   alert.PrevCount,
		NewCount:    alert.NewCount,
		Rate:        alert.Rate,
		Cycles:      alert.Cycles,
		Timestamp:   alert.Timestamp,
	}
}
func ConvertToJSONInfo(info *analyzer.Info) JSONInfo {
	return JSONInfo{
		ServiceName:      info.ServiceName,
		Timestamp:        info.Timestamp,
		Goroutines:       info.Goroutines,
		Alerts:           info.Alerts,
		TotalClusters:    info.TotalClusters,
		MaxClusterSize:   info.MaxClusterSize,
		TotalBlockEvents: info.TotalBlockEvents,
		TotalMutexEvents: info.TotalMutexEvents,
		Hotspots:         info.Hotspots,
	}
}
