import { LoadFile } from './wailsjs/go/main/App.js';

let filePathElement = document.getElementById("filePath");
let contentElement = document.getElementById("content");

window.loadFile = function () {
    let path = filePathElement.value;

    if (path === "") {
        alert("Please enter a file path");
        return;
    }

    try {
        LoadFile(path)
            .then((html) => {
                contentElement.innerHTML = html;
            })
            .catch((err) => {
                console.error(err);
                contentElement.innerHTML = `<div style="color: red;">Error loading file: ${err}</div>`;
            });
    } catch (err) {
        console.error(err);
        contentElement.innerHTML = `<div style="color: red;">Critical Error: ${err}</div>`;
    }
};

// Auto Reload
// Auto Reload
window.runtime.EventsOn("file:changed", (path) => {
    console.log("File changed event received:", path);
    console.log("Current input value:", filePathElement.value);

    // Normalize paths for comparison (Windows issue)
    const normPath = path.replace(/\\/g, '/');
    const normInput = filePathElement.value.replace(/\\/g, '/');

    console.log("Normalized Path:", normPath);
    console.log("Normalized Input:", normInput);

    // Reload if the changed file matches current input
    if (normInput === normPath) {
        console.log("Path match! Reloading...");
        loadFile();
    } else {
        console.log("Path mismatch. Not reloading.");
    }
});

// Debug Runtime
console.log("Wails Runtime:", window.runtime);

// Scrollbar Toggle
window.runtime.EventsOn("toggle-scrollbar", () => {
    document.body.classList.toggle("show-scrollbar");
    console.log("Scrollbar toggled, classes:", document.body.className);
});

// File Dropped via Wails Runtime
if (window.runtime && window.runtime.OnFileDrop) {
    window.runtime.OnFileDrop((x, y, paths) => {
        console.log("File dropped (Wails):", paths);
        if (paths.length > 0) {
            filePathElement.value = paths[0];
            loadFile();
        }
    });
} else {
    console.warn("window.runtime.OnFileDrop is not defined");
}

// Manual Drag & Drop handlers removed to test Wails Native Drop exclusively.
// If Wails EnableFileDrop + DisableWebViewDrop works, we don't need JS events.
