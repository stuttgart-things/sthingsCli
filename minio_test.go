/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

// func TestGetObjectsFromMinioBucket(t *testing.T) {

// 	os.Setenv("MINIO_ACCESS_KEY", "sthings")
// 	os.Setenv("MINIO_SECRET_KEY", "TOBESET")
// 	os.Setenv("MINIO_ADDR", "artifacts.automation.sthings-vsphere.labul.sva.de")
// 	os.Setenv("MINIO_SECURE", "true")

// 	bucket := "ansiblerun"
// 	file := "ansibleRuntroja.yaml"
// 	created, minioClient := CreateMinioClient()

// 	if created {
// 		files := GetObjectsFromMinioBucket(minioClient, bucket)
// 		fmt.Println(files)

// 		if !DownloadObjectFromMinioBucket(minioClient, bucket, file, "/tmp/"+file) {
// 			fmt.Println("FILE NOT DOWNLOADED")
// 		}
// 		uploaded, fileSize := UploadObjectToMinioBucket(minioClient, bucket, "/tmp/pod.yaml", "pod.yaml")

// 		if uploaded {
// 			log.Printf("SUCCESSFULLY UPLOADED OF SIZE %d\n", fileSize)
// 		}

// 	}

// }
