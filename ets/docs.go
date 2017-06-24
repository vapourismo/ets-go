// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

/*
Package ets provides functions to extract information from ETS .knxproj and .knxprod exports.

Opening the archive

You can open exported projects and product databases.

	archive, err := ets.OpenExportArchive("my-project.knxproj")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Make sure to close the archive eventually.
	defer archive.Close()

OpenExportArchive will scan for project and manufacturer files inside the given export archive.
Project and manufacturer files will be stored in ProjectFiles and ManufacturerFiles respectively.

	for _, manuFile := range archive.ManufacturerFiles {
		fmt.Println(manuFile.ContentID)
	}

	for _, projFile := range archive.ProjectFiles {
		fmt.Println(projFile.ProjectID)
	}

Decoding files

Not all files within the export might be relevant to you. Therefore, no files are decoded
automatically.

	for _, projFile := range archive.ProjectFiles {
		projInfo, err := projFile.Decode()
		if err != nil {
			log.Println(err)
			continue
		}

		// Variable projInfo contains the project info described in the projFile.
		fmt.Println("Project", projInfo.Name)
	}

This decodes only the basic project information. The actual project is contained in the installation
files. Each installation file can contain multiple project installations. Yes, this is weird.

	for _, instFile := range projFile.InstallationFiles {
		proj, err := instFile.Decode()
		if err != nil {
			log.Println(err)
			continue
		}

		for _, inst := range proj.Installations {
			fmt.Println("Installation", inst.Name)
		}
	}

*/
package ets
