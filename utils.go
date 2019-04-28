package makerbot

import (
	"encoding/binary"
	"strconv"
	"strings"
	"time"
)

// https://gist.github.com/alexmcroberts/219127816e7a16c7bd70
type epochTime time.Time

func (t epochTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

func (t *epochTime) UnmarshalJSON(s []byte) (err error) {
	r := strings.Replace(string(s), `"`, ``, -1)

	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(t) = time.Unix(q/1000, 0)
	return
}

func (t epochTime) String() string { return time.Time(t).String() }

// CameraFrameFormat specifies which format a camera frame is in
type CameraFrameFormat uint32

const (
	// CameraFrameFormatInvalid means that the camera frame is invalid(?)
	CameraFrameFormatInvalid CameraFrameFormat = iota
	// CameraFrameFormatYUYV means that the camera frame is in YUYV format
	CameraFrameFormatYUYV
	// CameraFrameFormatJPEG means that the camera frame is in JPEG format
	CameraFrameFormatJPEG
)

// CameraFrame is a single camera snapshot returned by the printer
type CameraFrame struct {
	Data     []byte
	Metadata *CameraFrameMetadata
}

// CameraFrameMetadata holds information about a camera frame returned
// by the printer
type CameraFrameMetadata struct {
	FileSize uint32            // Frame's file size in bytes
	Width    uint32            // Frame's width in pixels
	Height   uint32            // Frame's height in pixels
	Format   CameraFrameFormat // Format that the frame is in (invalid, YUYV, JPEG)
}

func unpackCameraFrameMetadata(packed []byte) CameraFrameMetadata {
	return CameraFrameMetadata{
		FileSize: binary.BigEndian.Uint32(packed[0:4]) - 16,
		Width:    binary.BigEndian.Uint32(packed[4:8]),
		Height:   binary.BigEndian.Uint32(packed[8:12]),
		Format:   CameraFrameFormat(binary.BigEndian.Uint32(packed[12:16])),
	}
}

func parseInfoFields(inf *[]string) *map[string]string {
	var fields map[string]string

	for _, field := range *inf {
		spl := strings.Split(field, "=")
		fields[spl[0]] = fields[spl[1]]
	}

	return &fields
}
