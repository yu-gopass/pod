package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/yu-gopass/pod/domain/model"
)

//创建需要实现的接口
type IPodRepository interface {
	//初始化表
	InitTable() error
	//根据ID查找数据
	FindPodByID(int64) (*model.Pod,error)
	//创建一条Pod 数据
	CreatePod(*model.Pod) (int64,error)
	//根据ID删除一条 pod 数据
	DeletePodByID(int64) error
	//修改一条数据
	UpdatePod(*model.Pod) error
	//查找pod所有数据
	FindAll() ([]model.Pod,error)
}

//创建 PodRepository
func NewPodRepository(db *gorm.DB) IPodRepository {
	return &PodRepository{mysqlDb:db}
}

type PodRepository struct {
	mysqlDb *gorm.DB
}

//初始化表
func (u *PodRepository)InitTable() error  {
	return u.mysqlDb.CreateTable(&model.Pod{},&model.PodEnv{},
	&model.PodPort{}).Error
}

//根据ID查找Pod信息
func (u *PodRepository)FindPodByID(podID int64) (pod *model.Pod,err error)  {
	pod = &model.Pod{}
	return pod, u.mysqlDb.Preload("PodEnv").Preload("PodPort").First(pod,podID).Error
}

//创建pod
func (u *PodRepository) CreatePod(pod *model.Pod) (int64, error) {
	return pod.ID, u.mysqlDb.Create(pod).Error
}

//根据 ID 删除 pod
func (u *PodRepository) DeletePodByID(podID int64) error{
	tx := u.mysqlDb.Begin()
	//遇到问题回滚
	defer func() {
		if r := recover();r!=nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	//没有问题，彻底删除pod
	if err := u.mysqlDb.Where("id = ?", podID).Delete(&model.Pod{}).Error; err!=nil{
		tx.Rollback()
		return err
	}
	//彻底删除 podenv
	if err := u.mysqlDb.Where("pod_id = ?", podID).Delete(&model.PodEnv{}).Error; err!=nil{
		tx.Rollback()
		return err
	}
	//彻底删除 podport
	if err := u.mysqlDb.Where("pod_id = ?", podID).Delete(&model.PodPort{}).Error; err!=nil{
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

//更新 pod
func (u *PodRepository) UpdatePod(pod *model.Pod) error {
	return u.mysqlDb.Model(pod).Update(pod).Error
}

//获取结果集合
func (u *PodRepository) FindAll() (podAll []model.Pod, err error) {
	return podAll, u.mysqlDb.Find(&podAll).Error
}