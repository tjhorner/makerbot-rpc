// Package printfile is a library for parsing .makerbot print files.
package printfile

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"image/png"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var thumbnailRegex = regexp.MustCompile(`thumbnail_([0-9]+)x([0-9]+)\.png`)

// Parse parses everything (metadata, thumbnails, toolpath) from a
// .makerbot file given a zip.ReadCloser
func Parse(r *zip.ReadCloser) (*MakerBotFile, error) {
	metadata, err := parseMetadata(r)
	if err != nil {
		return nil, err
	}

	toolpath, err := parseToolpath(r)
	if err != nil {
		return nil, err
	}

	thumbnails := parseThumbnails(r)
	if err != nil {
		return nil, err
	}

	return &MakerBotFile{
		ThumbnailSizes: thumbnails,
		Metadata:       metadata,
		Toolpath:       toolpath,
	}, nil
}

// GetMetadata grabs just the Metadata of a .makerbot file given a zip.ReadCloser
func GetMetadata(r *zip.ReadCloser) (*Metadata, error) {
	return parseMetadata(r)
}

// GetToolpath grabs just the Toolpath of a .makerbot file given a zip.ReadCloser
func GetToolpath(r *zip.ReadCloser) (*Toolpath, error) {
	return parseToolpath(r)
}

// GetThumbnails grabs just the []Thumbnail of a .makerbot file given a zip.ReadCloser
func GetThumbnails(r *zip.ReadCloser) *[]Thumbnail {
	return parseThumbnails(r)
}

// ParseFile parses everything (metadata, thumbnails, toolpath) from a
// .makerbot file given a filepath
func ParseFile(filename string) (*MakerBotFile, error) {
	rc, err := zip.OpenReader(filename)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return Parse(rc)
}

// GetFileMetadata grabs just the Metadata of a .makerbot file given a filepath
func GetFileMetadata(filename string) (*Metadata, error) {
	rc, err := zip.OpenReader(filename)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return parseMetadata(rc)
}

// GetFileToolpath grabs just the Toolpath of a .makerbot file given a filepath
func GetFileToolpath(filename string) (*Toolpath, error) {
	rc, err := zip.OpenReader(filename)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return parseToolpath(rc)
}

// GetFileThumbnails grabs just the []Thumbnail of a .makerbot file given a filepath
func GetFileThumbnails(filename string) (*[]Thumbnail, error) {
	rc, err := zip.OpenReader(filename)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return parseThumbnails(rc), nil
}

func parseMetadata(zr *zip.ReadCloser) (*Metadata, error) {
	var meta Metadata
	found := false

	for _, f := range zr.File {
		if f.Name != "meta.json" {
			continue
		}

		found = true

		fc, err := f.Open()
		if err != nil {
			return nil, err
		}

		dec := json.NewDecoder(fc)
		dec.Decode(&meta)
		break
	}

	if !found {
		return nil, errors.New("parseMetadata: malformed .makerbot file; does not have metadata")
	}

	return &meta, nil
}

func parseToolpath(zr *zip.ReadCloser) (*Toolpath, error) {
	var tp Toolpath
	found := false

	for _, f := range zr.File {
		if f.Name != "print.jsontoolpath" {
			continue
		}

		found = true

		fc, err := f.Open()
		if err != nil {
			return nil, err
		}

		dec := json.NewDecoder(fc)
		dec.Decode(&tp)
		break
	}

	if !found {
		return nil, errors.New("parseToolpath: malformed .makerbot file; does not have toolpath")
	}

	return &tp, nil
}

func parseThumbnail(f *zip.File) (*Thumbnail, error) {
	fc, err := f.Open()
	if err != nil {
		// Errors are negligible here
		return nil, err
	}

	matches := thumbnailRegex.FindStringSubmatch(f.Name)
	if len(matches) < 3 {
		return nil, err
	}

	width, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, err
	}

	height, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, err
	}

	img, err := png.DecodeConfig(fc)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(fc)
	if err != nil {
		return nil, err
	}

	thumb := Thumbnail{
		TargetHeight: height,
		TargetWidth:  width,
		ActualHeight: img.Height,
		ActualWidth:  img.Width,
		Data:         data,
	}

	return &thumb, nil
}

func parseThumbnails(zr *zip.ReadCloser) *[]Thumbnail {
	var thumbnails []Thumbnail

	for _, f := range zr.File {
		if !strings.HasPrefix(f.Name, "thumbnail_") || !strings.HasSuffix(f.Name, ".png") {
			continue
		}

		thumb, err := parseThumbnail(f)
		if err != nil {
			continue
		}

		thumbnails = append(thumbnails, *thumb)
	}

	return &thumbnails
}
