package capsule

import (
	"fmt"
	"log"
	"os"

	diskfs "github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/disk"
	"github.com/diskfs/go-diskfs/filesystem"
	"github.com/diskfs/go-diskfs/filesystem/iso9660"
)

// ImageManager ... Create/destroy/manipulate disk images
type ImageManager struct {
	VolumeLabel string
	FsDisk      *disk.Disk
	Fs          filesystem.FileSystem
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

	// the following line is required for an ISO, which may have logical block sizes
	// only of 2048, 4096, 8192
	i.FsDisk.LogicalBlocksize = 2048
	fspec := disk.FilesystemSpec{
		Partition:   0,
		FSType:      filesystem.TypeISO9660,
		VolumeLabel: i.VolumeLabel,
	}
	newFs, err := i.FsDisk.CreateFilesystem(fspec)
	handler.HandleErrors(err)
	i.Fs = newFs

	rw, err := i.Fs.OpenFile("demo.txt", os.O_CREATE|os.O_RDWR)
	content := []byte("demo")
	_, err = rw.Write(content)
	handler.HandleErrors(err)

	iso, ok := i.Fs.(*iso9660.FileSystem)
	if !ok {
		handler.HandleErrors(fmt.Errorf("not an iso9660 filesystem"))
	}
	err = iso.Finalize(iso9660.FinalizeOptions{})
	handler.HandleErrors(err)
}
