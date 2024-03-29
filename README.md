# stuttgart-things/sthingsCli

providing golang building blocks for the use in command line interface modules

## MODULES

### GIT

<details><summary>CLONE GIT REPOSITORY, GET FILE LIST + READ FILE</summary>

```go
gitRepository         := "https://github.com/stuttgart-things/kaeffken.git"
gitBranch             := "main"
gitCommitID           := "09de9ff7b5c76aff8bb32f68cfb0bbe49cd5a7a8"

repo, cloneStatus := sthingsCli.CloneGitRepository(gitRepository, gitBranch, gitCommitID, nil)

if cloneStatus {
    fileList, directoryList = sthingsCli.GetFileListFromGitRepository("", repo)
    fmt.Println(fileList, directoryList)
}

readMe := sthingsCli.ReadFileContentFromGitRepo(repo, "README.md")
fmt.Println(readMe)
```

</details>


## TASKFILE

```bash
task: Available tasks for this project:
* lint:       Lint code
* tag:        commit, push & tag the module
* test:       Test code
```

## LICENSE

<details><summary><b>APACHE 2.0</b></summary>

Copyright 2023 patrick hermann.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

</details>

Author Information
------------------
Patrick Hermann, stuttgart-things 04/2023
