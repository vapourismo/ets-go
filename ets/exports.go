// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package ets

import (
	"archive/zip"
	"path"
	"regexp"
)

// InstallationFile is a file that contains zero or more project installations.
type InstallationFile struct {
	*zip.File
}

// Decode the file in order to retrieve the project inside it.
func (i *InstallationFile) Decode() (p *Project, err error) {
	r, err := i.Open()
	if err != nil {
		return
	}

	p, err = DecodeProject(r)
	r.Close()

	return
}

// ProjectFile is a file that contains project information.
type ProjectFile struct {
	*zip.File

	ProjectID         string
	InstallationFiles []InstallationFile
}

// Decode the file in order to retrieve the project info inside it.
func (pf *ProjectFile) Decode() (pi *ProjectInfo, err error) {
	r, err := pf.Open()
	if err != nil {
		return
	}

	pi, err = DecodeProjectInfo(r)
	r.Close()

	return
}

var projectFileBaseRe = regexp.MustCompile("^\\d.xml$")

func newProjectFile(archive *zip.ReadCloser, metaFile *zip.File) (projFile ProjectFile) {
	projectDir := path.Dir(metaFile.Name)

	projFile.File = metaFile
	projFile.ProjectID = projectDir

	// Search for the project installation file.
	for _, file := range archive.File {
		if path.Dir(file.Name) == projectDir && projectFileBaseRe.MatchString(path.Base(file.Name)) {
			projFile.InstallationFiles = append(projFile.InstallationFiles, InstallationFile{file})
		}
	}

	return
}

// ManufacturerFile is a manufacturer file.
type ManufacturerFile struct {
	*zip.File

	ManufacturerID string
	ContentID      string
}

// ExportArchive is a handle to an exported archive (.knxproj or .knxprod).
type ExportArchive struct {
	archive *zip.ReadCloser

	ProjectFiles      []ProjectFile
	ManufacturerFiles []ManufacturerFile
}

// OpenExportArchive opens the exported archive located at given path.
func OpenExportArchive(path string) (*ExportArchive, error) {
	archive, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}

	ex := &ExportArchive{archive: archive}

	if err = ex.findFiles(); err != nil {
		archive.Close()
		return nil, err
	}

	return ex, nil
}

var (
	projectMetaFileRe  = regexp.MustCompile("^(p|P)-([0-9a-zA-Z]+)/(p|P)roject.xml$")
	manufacturerFileRe = regexp.MustCompile("^(m|M)-([0-9a-zA-Z]+)/(m|M)-([0-9a-zA-Z]+)([^.]+).xml$")

	// TODO: Figure out if '/' is a universal path seperator in ZIP files.
)

func (ex *ExportArchive) findFiles() error {
	for _, file := range ex.archive.File {
		if projectMetaFileRe.MatchString(file.Name) {
			ex.ProjectFiles = append(ex.ProjectFiles, newProjectFile(ex.archive, file))
		} else if matches := manufacturerFileRe.FindStringSubmatch(file.Name); matches != nil {
			ex.ManufacturerFiles = append(ex.ManufacturerFiles, ManufacturerFile{
				File:           file,
				ManufacturerID: "M-" + matches[2],
				ContentID:      "M-" + matches[4] + matches[5],
			})
		}
	}

	return nil
}

// Close the archive handle.
func (ex *ExportArchive) Close() error {
	return ex.archive.Close()
}
