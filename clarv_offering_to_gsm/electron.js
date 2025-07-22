const { app, BrowserWindow } = require('electron');
const path = require('path');
const isDev = require('electron-is-dev');

let mainWindow;

function createWindow() {
  console.log('Creating main window...');

  mainWindow = new BrowserWindow({
    width: 1000,
    height: 800,
    // show: false,
    webPreferences: {
      nodeIntegration: true,
      contextIsolation: false,
    },
  });

  const startURL = isDev
    ? 'http://localhost:3000'
    : `file://${path.join(__dirname, 'build', 'index.html')}`;

  console.log('Loading URL:', startURL);
  mainWindow.loadURL(startURL);

  mainWindow.once('ready-to-show', () => {
    console.log('Window ready to show');
    mainWindow.show();

    if (isDev) {
      mainWindow.webContents.openDevTools({ mode: 'detach' }); // opens devtools in separate window
    }
  });

  mainWindow.webContents.on('did-finish-load', () => {
    console.log('React app finished loading');
  });

  mainWindow.on('closed', () => {
    console.log('Main window closed');
    mainWindow = null;
  });

  mainWindow.webContents.on('crashed', () => {
    console.error('Renderer process crashed');
  });
}

app.on('ready', () => {
  console.log('App ready');
  createWindow();
});

app.on('window-all-closed', () => {
  console.log('All windows closed');
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('activate', () => {
  console.log('App activated');
  if (mainWindow === null) {
    createWindow();
  }
});


// .\node_modules\.bin\electron.cmd . --enable-logging And npm start
// To run the app

// go run main.go crypto.go 
// to run backend