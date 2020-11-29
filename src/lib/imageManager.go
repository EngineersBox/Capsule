package capsule

import (
	"fmt"
	"log"
	"os"

	diskfs "github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/disk"
	"github.com/diskfs/go-diskfs/filesystem"
	"github.com/diskfs/go-diskfs/filesystem/iso9660"
	"github.com/google/uuid"
)

// BlockSizes ... Logical block sizes
type BlockSizes struct {
	Bs2048 int64
	Bs4096 int64
	Bs8192 int64
}

var blockSizes BlockSizes = BlockSizes{2048, 4096, 8192}

const basePartitionIdentifier string = "BASE_FS_DECL.txt"

// ImageManager ... Create/destroy/manipulate disk images
type ImageManager struct {
	ID          uuid.UUID
	VolumeLabel string
	FsDisk      *disk.Disk
	Fs          filesystem.FileSystem
}

// CreateISOFileSystem ... Create an ISO9660 compatible filespace with a single partition
func (i *ImageManager) CreateISOFileSystem() {
	// ISO block sizes can only be one of 2048, 4096, 8192
	i.FsDisk.LogicalBlocksize = blockSizes.Bs2048
	fspec := disk.FilesystemSpec{
		Partition:   0,
		FSType:      filesystem.TypeISO9660,
		VolumeLabel: i.VolumeLabel,
	}
	newFs, err := i.FsDisk.CreateFilesystem(fspec)
	handler.HandleErrors(err)
	i.Fs = newFs
}

// InitialiseFsPartitionIdentifier ... Create a FS identifier from ID + VolumeLabel + FsDisk.Table.Type and write at root dir "/"
func (i *ImageManager) InitialiseFsPartitionIdentifier() {
	rw, err := i.Fs.OpenFile(basePartitionIdentifier, os.O_CREATE|os.O_RDWR)
	content := []byte(string(i.ID.String() + i.VolumeLabel + i.FsDisk.Table.Type()))
	_, err = rw.Write(content)
	handler.HandleErrors(err)
}

// CreateIso ... Create an ISO9660 disk image
func (i *ImageManager) CreateIso(diskImg string) {
	if diskImg == "" {
		log.Fatal("must have a valid path for diskImg")
	}
	var diskSize int64
	diskSize = 10 * 1024 * 1024 // 10 MB
	newDisk, err := diskfs.Create(diskImg, diskSize, diskfs.Raw)
	handler.HandleErrors(err)
	i.FsDisk = newDisk

	i.CreateISOFileSystem()

	iso, ok := i.Fs.(*iso9660.FileSystem)
	if !ok {
		handler.HandleErrors(fmt.Errorf("not an iso9660 filesystem"))
	}
	err = iso.Finalize(iso9660.FinalizeOptions{})
	handler.HandleErrors(err)
}
