/**
 * Common JavaScript for all pages
 * Handles mobile menu, canvas header, workspace client
 */

document.addEventListener('DOMContentLoaded', () => {
  initMobileMenu();
  initCanvasHeader();
  initWorkspaceClient();
});

/**
 * Initialize mobile menu
 */
function initMobileMenu() {
  const mobileMenuToggle = document.getElementById('mobileMenuToggle');
  const mobileMenu = document.getElementById('mobileMenu');
  
  if (mobileMenuToggle && mobileMenu) {
    mobileMenuToggle.addEventListener('click', () => {
      mobileMenu.classList.toggle('open');
    });
    
    // Close menu when clicking outside
    document.addEventListener('click', (e) => {
      if (!mobileMenu.contains(e.target) && !mobileMenuToggle.contains(e.target)) {
        mobileMenu.classList.remove('open');
      }
    });
  }
}

/**
 * Initialize canvas header
 */
function initCanvasHeader() {
  const baseURL = window.location.origin;
  
  // Fetch installation info
  fetch(`${baseURL}/api/installation/info`)
    .then(response => response.json())
    .then(data => {
      if (data.installation_name) {
        const element = document.getElementById('installationName');
        if (element) element.textContent = data.installation_name;
      }
    })
    .catch(error => {
      console.error('Error fetching installation info:', error);
    });
  
  // Fetch current canvas info
  fetch(`${baseURL}/api/canvas/info`)
    .then(response => response.json())
    .then(data => {
      if (data.canvas_name) {
        const element = document.getElementById('canvasName');
        if (element) element.textContent = data.canvas_name;
      }
    })
    .catch(error => {
      console.error('Error fetching canvas info:', error);
    });
}

/**
 * Initialize workspace client
 */
function initWorkspaceClient() {
  if (typeof WorkspaceClient === 'undefined') {
    return;
  }
  
  const workspaceClient = new WorkspaceClient();
  const baseURL = window.location.origin;
  
  workspaceClient.on('connected', () => {
    updateStatus('connected', 'Connected');
  });
  
  workspaceClient.on('disconnected', () => {
    updateStatus('disconnected', 'Disconnected');
  });
  
  workspaceClient.on('reconnecting', (attempts) => {
    updateStatus('connecting', `Reconnecting (${attempts})...`);
  });
  
  workspaceClient.on('canvas_update', (data) => {
    if (data.canvas_name) {
      const element = document.getElementById('canvasName');
      if (element) element.textContent = data.canvas_name;
    }
  });
  
  function updateStatus(status, text) {
    const indicator = document.getElementById('statusIndicator');
    const statusText = document.getElementById('statusText');
    
    if (indicator) indicator.className = `status-indicator ${status}`;
    if (statusText) statusText.textContent = text;
  }
  
  workspaceClient.connect(baseURL);
}

