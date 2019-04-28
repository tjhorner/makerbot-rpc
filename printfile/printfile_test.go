package printfile_test

import (
	"testing"

	"github.com/tjhorner/makerbot-rpc/printfile"
)

const (
	file = "test/box.makerbot"

	expectedUUID             = "3c997805-959b-414f-85b5-c45872a11b78"
	expectedThumbnails       = 3
	expectedToolpathCommands = 9158
)

func TestParseFile(t *testing.T) {
	file, err := printfile.ParseFile(file)
	if err != nil {
		t.Error(err)
	}

	if file.Metadata.UUID != expectedUUID {
		t.Errorf("metadata UUID is wrong; wanted: %s, got: %s\n", expectedUUID, file.Metadata.UUID)
	}

	if len(*file.ThumbnailSizes) != expectedThumbnails {
		t.Errorf("print thumbnails size is wrong; wanted: %d, got: %d\n", expectedThumbnails, len(*file.ThumbnailSizes))
	}

	if len(*file.Toolpath) != expectedToolpathCommands {
		t.Errorf("toolpath commands size is wrong; wanted: %d, got: %d\n", expectedToolpathCommands, len(*file.Toolpath))
	}

	if len(*file.Toolpath) != file.Metadata.TotalCommands {
		t.Errorf("toolpath commands size is inconsistent; wanted: %d, got: %d\n", file.Metadata.TotalCommands, len(*file.Toolpath))
	}
}

func TestGetFileMetadata(t *testing.T) {
	metadata, err := printfile.GetFileMetadata(file)
	if err != nil {
		t.Error(err)
	}

	if metadata.UUID != expectedUUID {
		t.Errorf("metadata UUID is wrong; wanted: %s, got: %s\n", expectedUUID, metadata.UUID)
	}
}

func TestGetFileThumbnails(t *testing.T) {
	thumbs, err := printfile.GetFileThumbnails(file)
	if err != nil {
		t.Error(err)
	}

	if len(*thumbs) != expectedThumbnails {
		t.Errorf("print thumbnails size is wrong; wanted: %d, got: %d\n", expectedThumbnails, len(*thumbs))
	}
}

func TestGetFileToolpath(t *testing.T) {
	tp, err := printfile.GetFileToolpath(file)
	if err != nil {
		t.Error(err)
	}

	if len(*tp) != expectedToolpathCommands {
		t.Errorf("toolpath commands size is wrong; wanted: %d, got: %d\n", expectedToolpathCommands, len(*tp))
	}
}
