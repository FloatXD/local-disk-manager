package scheduler

import (
	"fmt"

	"github.com/hwameistor/local-disk-manager/pkg/csi/volumemanager"

	"github.com/hwameistor/local-disk-manager/pkg/csi/diskmanager"
	"github.com/hwameistor/local-disk-manager/pkg/csi/driver/identity"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	storagev1lister "k8s.io/client-go/listers/storage/v1"
)

// DiskVolumeSchedulerPlugin implement the Scheduler interface
// defined in github.com/hwameistor/scheduler/pkg/scheduler/scheduler.go: Scheduler
type DiskVolumeSchedulerPlugin struct {
	diskNodeHandler   diskmanager.DiskManager
	diskVolumeHandler volumemanager.VolumeManager

	boundVolumes     []string
	pendingVolumes   []*v1.PersistentVolumeClaim
	tobeScheduleNode *v1.Node
	scLister         storagev1lister.StorageClassLister
}

func NewDiskVolumeSchedulerPlugin(scLister storagev1lister.StorageClassLister) *DiskVolumeSchedulerPlugin {
	return &DiskVolumeSchedulerPlugin{
		diskNodeHandler:   diskmanager.NewLocalDiskManager(),
		diskVolumeHandler: volumemanager.NewLocalDiskVolumeManager(),
		scLister:          scLister,
	}
}

// Filter whether the node meets the storage requirements of pod runtime.
// The following two types of situations need to be met at the same time:
//1. If the pod uses a created volume, we need to ensure that the volume is located at the scheduled node.
//2. If the pod uses a pending volume, we need to ensure that the scheduled node can meet the requirements of the volume.
func (s *DiskVolumeSchedulerPlugin) Filter(boundVolumes []string, pendingVolumes []*v1.PersistentVolumeClaim, node *v1.Node) (bool, error) {
	s.filterFor(boundVolumes, pendingVolumes, node)

	// step1: filter bounded volumes
	ok, err := s.filterExistVolumes()
	if err != nil {
		log.WithError(err).Errorf("failed to filter node %s for bounded volumes %v due to error: %s", node.GetName(), boundVolumes, err.Error())
		return false, err
	}
	if !ok {
		log.Infof("node %s is not suitable because of bounded volumes %v is already located on the other node", node.GetName(), boundVolumes)
		return false, nil
	}

	// step2: filter pending volumes
	ok, err = s.filterPendingVolumes()
	if err != nil {
		log.WithError(err).Infof("failed to filter node %s for pending volumes due to error: %s", node.GetName(), err.Error())
		return false, err
	}
	if !ok {
		log.Infof("node %s is not suitable", node.GetName())
		return false, nil
	}

	return true, nil
}

// Reserve disk needed by the volumes
func (s *DiskVolumeSchedulerPlugin) Reserve(pendingVolumes []*v1.PersistentVolumeClaim, node string) error {
	log.WithFields(log.Fields{"node": node, "volumes": pendingVolumes}).Debug("reserving disk")
	for _, pvc := range pendingVolumes {
		diskReq, err := s.convertPVCToDiskRequest(pvc, node)
		if err != nil {
			return err
		}
		if err = s.diskNodeHandler.ReserveDiskForVolume(diskReq, pvc.GetName()); err != nil {
			return err
		}
	}
	return nil
}

// Unreserve disk reserved by the volumes on the node
func (s *DiskVolumeSchedulerPlugin) Unreserve(pendingVolumes []*v1.PersistentVolumeClaim, node string) error {
	log.WithFields(log.Fields{"node": node, "volumes": pendingVolumes}).Debug("unreserving disk")
	for _, pvc := range pendingVolumes {
		if err := s.diskNodeHandler.UnReserveDiskForPVC(pvc.GetName()); err != nil {
			return err
		}
	}
	return nil
}

func (s *DiskVolumeSchedulerPlugin) filterFor(boundVolumes []string, pendingVolumes []*v1.PersistentVolumeClaim, node *v1.Node) {
	s.boundVolumes = boundVolumes
	s.pendingVolumes = pendingVolumes
	s.tobeScheduleNode = node
	s.removeDuplicatePVC()
}

func (s *DiskVolumeSchedulerPlugin) removeDuplicatePVC() {
	pvcMap := map[string]*v1.PersistentVolumeClaim{}
	pvcCopy := s.pendingVolumes
	for i, pvc := range pvcCopy {
		if _, ok := pvcMap[pvc.GetName()]; ok {
			s.pendingVolumes = append(s.pendingVolumes[:i], s.pendingVolumes[i+1:]...)
		}
	}
}

// filterExistVolumes compare the tobe scheduled node is equal to the node where volume already located at
func (s *DiskVolumeSchedulerPlugin) filterExistVolumes() (bool, error) {
	for _, name := range s.boundVolumes {
		volume, err := s.diskVolumeHandler.GetVolumeInfo(name)
		if err != nil {
			log.WithError(err).Errorf("failed to get volume %s info", name)
			return false, err
		}
		if volume.AttachNode != s.tobeScheduleNode.GetName() {
			log.Infof("bounded volume is located at node %s,so node %s is not suitable", volume.AttachNode,
				s.tobeScheduleNode.GetName())
			return false, nil
		}
	}

	return true, nil
}

func (s *DiskVolumeSchedulerPlugin) convertPVCToDiskRequest(pvc *v1.PersistentVolumeClaim, node string) (diskmanager.Disk, error) {
	sc, err := s.getParamsFromStorageClass(pvc)
	if err != nil {
		log.WithError(err).Errorf("failed to parse params from StorageClass")
		return diskmanager.Disk{}, err
	}

	storage := pvc.Spec.Resources.Requests[v1.ResourceStorage]
	return diskmanager.Disk{
		AttachNode: node,
		Capacity:   storage.Value(),
		DiskType:   sc.DiskType,
	}, nil
}

func (s *DiskVolumeSchedulerPlugin) getParamsFromStorageClass(volume *v1.PersistentVolumeClaim) (*StorageClassParams, error) {
	// sc here can't be empty,
	// more info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#class-1
	if volume.Spec.StorageClassName == nil {
		return nil, fmt.Errorf("storageclass in pvc %s can't be empty", volume.GetName())
	}

	sc, err := s.scLister.Get(*volume.Spec.StorageClassName)
	if err != nil {
		return nil, err
	}

	return parseParams(sc.Parameters), nil
}

// filterPendingVolumes select free disks for pending pvc
func (s *DiskVolumeSchedulerPlugin) filterPendingVolumes() (bool, error) {
	var reqDisks []diskmanager.Disk
	for _, pvc := range s.pendingVolumes {
		disk, err := s.convertPVCToDiskRequest(pvc, s.tobeScheduleNode.GetName())
		if err != nil {
			return false, err
		}
		reqDisks = append(reqDisks, disk)
	}

	return s.diskNodeHandler.FilterFreeDisks(reqDisks)
}

func (s *DiskVolumeSchedulerPlugin) CSIDriverName() string {
	return identity.DriverName
}
