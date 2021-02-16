# rescale-starter-project
## Assignment
You've been tasked by Codebase to build a Codebase Transfer Manager! 

This is a system
for uploading files to their platform, downloading files from their platform, and
viewing what files have been uploaded. 

Your job is to implement file upload, file download, and viewing all files.

After implementing this, feel free to get creative! You can improve on anything you see fit!

Here are some ideas:

* notifications for when your file is done downloading
* notifications for when your file is done uploading
* beautiful UI
* multiple screens

## Directory Structure

```
fileserver.go
```
This is the fileserver that you can ping to upload files, download files, and view file IDs
for the files you have uploaded. 

It simulates the Rescale platform. You don't need to edit this file.

```
tmp
```
This is where the fileserver stores the files. You shouldn't have to touch this directory.

```
codebase_transfer_manager
```
This is the Electron app you will be editing.

```
codebase_transfer_manager/public/electron.js
```
This is the Electron part of your app. Put your Electron specific code here, such as creating windows,
sending HTTP requests, and callbacks for the main Electron process.

```
codebase_transfer_manager/public/index.html
```
This is the html file, but the only element is a root div which holds the React components.

```
codebase_transfer_manager/public/public.js
```
This is the contextBridge file, which connects the frontend (React) and backend (Electron). It exposes
an API that you can use in React (with window.api).

```
codebase_transfer_manager/src
```
This directory contains the React components.

```
codebase_transfer_manager/src/App.js
```
This is the root React component. If you make more React components, you should probably import them and display them here. 

```
codebase_transfer_manager/src/index.js
```
This file displays the root App component in index.html.

## Server 
After downloading go, you can run the server with
```
go run fileserver.go
```

The server has 4 endpoints:  
```
/upload
```
where you attach a file as a multipart form. The key should be "uploadFile" and the value will be your file, 
which will be attached as multipart/form-data in the body of the request. This will return a fileID on success, 
which you use to download your file later.

```
Here you download files. Pass in the fileId from earlier as a parameter. 
The key should be "file" and the value should be the fileId.
```
where you download files. It expects the fileId returned from upload.

```
/getallfiles
```
is where you can view a list of fileIds and their associated filenames.

```
/hello
```
to get a greeting message. You can use this to check if you are pinging the server properly.

## Getting Started
Install all the node modules with
```
npm install
```
Run the Electron app with
```
npm run
```
You may need to refresh the app window (CMD + R) or (CTRL + R) if you see a blank white screen.

