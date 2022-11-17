package model

//Pod的状态
type Pod struct {
	ID int64 `gorm:"primary_key;not_null;auto_increment" json:"id"`
	PodName string `gorm:"unique_index;not_null" json:"pod_name"`
	PodNamespace string `json:"pod_namespace"`
	//POD 所属团队
	PodTeamID int64 `json:"pod_team_id"`
	//POD 使用的CPU最小值
	PodCpuMin float32 `json:"pod_cpu_min"`
	//POD 使用的CPU最大值
	PodCpuMax float32 `json:"pod_cpu_max"`
	//副本数量
	PodRelicas int32 `json:"pod_relicas"`
	//POD 使用的内存最小值
	PodMemoryMin float32 `json:"pod_memory_min"`
	//POD 使用的内存最大值
	PodMemoryMax float32 `json:"pod_memory_max"`
	//POD 开放的端口
	PodPort[] PodPort `gorm:"ForeignKey:PodID" json:"pod_port"`
	//POD 使用的环境变量
	PodEnv[] PodEnv `gorm:"ForeignKey:PodID" json:"pod_env"`
	//镜像拉取策略
	PodPullPolicy string `json:"pod_pull_policy"`
	//重启策略
	PodRestart string `json:"pod_restart"`
	///pod发布策略
	PodType string `json:"pod_type"`
	//使用的镜像名称+tag
	PodImage string `json:"pod_image"`
	//@TODO 挂盘
	//@TODO 域名设置
}
