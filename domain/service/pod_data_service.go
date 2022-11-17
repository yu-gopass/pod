package service

import (
	"context"
	"github.com/yu-gopass/pod/domain/model"
	"github.com/yu-gopass/pod/domain/repository"
	"github.com/yu-gopass/pod/proto/pod"
	v1 "k8s.io/api/apps/v1"
	v13 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strconv"
)

type IPodDataService interface {
	AddPod(*model.Pod) (int64, error)
	DeletePod(int64) error
	UpdatePod(*model.Pod) error
	FindPodByID(int64) (*model.Pod, error)
	FindAllPod() ([]model.Pod, error)
	CreateToK8s(*pod.PodInfo) error
	DeleteFromK8s(*model.Pod) error
	UpdateToK8s(*pod.PodInfo) error
}

func NewPodDataService(podRepository repository.PodRepository, clientset *kubernetes.Clientset) IPodDataService {
	return &PodDataService{
		PodRepository: podRepository,
		K8sClientSet:  clientset,
		deployment:    v1.Deployment{},
	}
}

type PodDataService struct {
	PodRepository repository.IPodRepository
	K8sClientSet  *kubernetes.Clientset
	deployment    *v1.Deployment
}

func (u *PodDataService) CreateToK8s(podInfo *pod.PodInfo) (err error) {
	u.SetDeployment(podInfo)
	if _,err = u.K8sClientSet.AppsV1().Deployments(podInfo.PodNamespace).Get(context.TODO(),podInfo.PodName,v12.GetOptions{});err!=nil{

	}
}

func (u *PodDataService) SetDeployment(podInfo *pod.PodInfo) {
	deployment := &v1.Deployment{}
	deployment.TypeMeta = v12.TypeMeta{
		Kind:       "deployment",
		APIVersion: "v1",
	}
	deployment.ObjectMeta = v12.ObjectMeta{
		Name:      podInfo.PodName,
		Namespace: podInfo.PodNamespace,
		Labels: map[string]string{
			"app-name": podInfo.PodName,
			"author":   "yuge",
		},
	}
	deployment.Name = podInfo.PodName
	deployment.Spec = v1.DeploymentSpec{
		Replicas: &podInfo.PodReplicas,
		Selector: &v12.LabelSelector{
			MatchLabels: map[string]string{
				"app-name": podInfo.PodName,
			},
			MatchExpressions: nil,
		},
		Template: v13.PodTemplateSpec{
			ObjectMeta: v12.ObjectMeta{
				Labels: map[string]string{
					"app-name": podInfo.PodName,
				},
			},
			Spec: v13.PodSpec{
				Containers: []v13.Container{
					{
						Name:      podInfo.PodName,
						Image:     podInfo.PodImages,
						Ports:     u.getContainerPort(podInfo),
						Env:       u.getEnv(podInfo),
						Resources: u.getResources(podInfo),
						ImagePullPolicy:u.getImagePullPolicy(podInfo),
					},
				},
			},
		},
		Strategy:                v1.DeploymentStrategy{},
		MinReadySeconds:         0,
		RevisionHistoryLimit:    nil,
		Paused:                  false,
		ProgressDeadlineSeconds: nil,
	}
	u.deployment = deployment
}

//获取容器端口
func (u *PodDataService) getContainerPort(podInfo *pod.PodInfo) (containerPort []v13.ContainerPort) {
	for _, v := range podInfo.PodPort {
		containerPort = append(containerPort, v13.ContainerPort{
			Name:          "port-" + strconv.FormatInt(int64(v.ContainerPort), 10),
			ContainerPort: v.ContainerPort,
			Protocol:      u.getProtocol(v.Protocol),
		})
	}
	return
}

//获取协议
func (u *PodDataService) getProtocol(protocol string) v13.Protocol {
	switch protocol {
	case "TCP":
		return "TCP"
	case "UDP":
		return "UDP"
	case "SCTP":
		return "SCTP"
	default:
		return "TCP"
	}
}

//获取环境变量
func (u *PodDataService) getEnv(podInfo *pod.PodInfo) (envVar []v13.EnvVar) {
	for _, v := range podInfo.PodEnv {
		envVar = append(envVar, v13.EnvVar{
			Name:      v.EnvKey,
			Value:     v.EnvValue,
			ValueFrom: nil,
		})
	}
	return
}

//获取资源限制
func (u *PodDataService) getResources(podInfo *pod.PodInfo) (source v13.ResourceRequirements) {
	//最大能够使用多少资源
	source.Limits = v13.ResourceList{
		"cpu":    resource.MustParse(strconv.FormatFloat(float64(podInfo.PodCpuMax), 'f', 6, 64)),
		"memory": resource.MustParse(strconv.FormatFloat(float64(podInfo.PodMemoryMax), 'f', 6, 64)),
	}
	//请求值
	source.Requests = v13.ResourceList{
		"cpu":    resource.MustParse(strconv.FormatFloat(float64(podInfo.PodCpuMax), 'f', 6, 64)),
		"memory": resource.MustParse(strconv.FormatFloat(float64(podInfo.PodMemoryMax), 'f', 6, 64)),
	}
	return
}

//获取拉取策略
func (u *PodDataService) getImagePullPolicy(podInfo *pod.PodInfo) v13.PullPolicy {
	switch podInfo.PodPullPolicy {
	case "Always":
		return "Always"
	case "Never":
		return "Never"
	case "IfNotPresent":
		return "IfNotPresent"
	default:
		return "Always"
	}
}

//添加pod
func (u *PodDataService) AddPod(pod2 *model.Pod) (int64, error) {
	return u.PodRepository.CreatePod(pod2)
}

//删除
func (u *PodDataService) DeletePod(podID int64) error {
	return u.PodRepository.DeletePodByID(podID)
}

//更新
func (u *PodDataService) UpdatePod(pod2 *model.Pod) error {
	return u.PodRepository.UpdatePod(pod2)
}

//单个ID查找
func (u *PodDataService) FindPodByID(podID int64) (*model.Pod, error) {
	return u.PodRepository.FindPodByID(podID)
}

//查找所有
func (u *PodDataService) FindAllPod() ([]model.Pod, error) {
	return u.PodRepository.FindAll()
}
