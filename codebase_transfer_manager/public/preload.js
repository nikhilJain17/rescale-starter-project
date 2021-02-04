const {
    contextBridge,
    ipcRenderer
} = require("electron");

// Expose protected methods that allow the renderer process to use
// the ipcRenderer without exposing the entire object
contextBridge.exposeInMainWorld(
    "api", {
        send: (channel, data) => {
            // whitelist channels
            let validChannels = ["toMain", "notify", "upload", "download"];

            // Send an event from the Renderer process to the Main process
            if (validChannels.includes(channel)) {
                ipcRenderer.send(channel, data);
            }
        },
        receive: (channel, func) => {
            let validChannel = (channel == 'fromMain') || (channel == 'allFiles');
            
            // Attach a listener to the Renderer process
            if (validChannel) {
                // Deliberately strip event as it includes `sender` 
                console.log('[Middleman] Attaching listener', channel);
                ipcRenderer.on(channel, (event, ...args) => func(...args));
            }
        }
    }
);
