// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package ets

import (
	"fmt"
	"log"
)

func Example() {
	archive, err := OpenExportArchive("my-project.knxproj")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer archive.Close()

	for _, projFile := range archive.ProjectFiles {
		projInfo, err := projFile.Decode()
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Println("Project", projInfo.ID, projInfo.Name)

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
	}

	// Output:
	// Project P-XXXX MyProject
	// Installation MyInstallation
}
