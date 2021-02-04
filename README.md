# rescale-starter-project
starter project for rescale (codebase sp21)

## Server 
The server has 3 endpoints:  
```
/upload
```
where you attach a file as a multipart form. The key should be "uploadFile" and 
the value will be your file. This will return a fileID on success, which you use to download your file.

```
/file/{filei\Id}.{extension}
```
where you download files. It expects the fileId returned from upload, as well as the extension of the original file.

```
/getallfiles
```
is where you can view a list of fileIds and their associated filenames.
