/**
 * Common JavaScript for all pages
 * Handles mobile menu, canvas header, workspace client
 */

// Error handler (loaded before this script)
let errorHandler;
if (typeof ErrorHandler !== 'undefined') {
  errorHandler = new ErrorHandler();
} else {
  // Fallback if error handler not loaded
  errorHandler = {
    logError: (msg, err, ctx) => {
      // Silent in production, console.error in development
      if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
        console.error(ctx ? `${ctx}: ${msg}` : msg, err);
      }
    }
  };
}

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
 * Initialize navbar tracking info (persists across all pages)
 */
function initCanvasHeader() {
  const baseURL = window.location.origin;

  // Update navbar tracking info from sessionStorage immediately (for fast page loads)
  const storedClientName = sessionStorage.getItem('clientName');
  const storedClientWarning = sessionStorage.getItem('clientWarning') === 'true';
  const navbarClientName = document.getElementById('navbarClientName');
  const navbarClientWarning = document.getElementById('navbarClientWarning');

  if (navbarClientName && storedClientName) {
    navbarClientName.textContent = storedClientName;
  }
  if (navbarClientWarning) {
    navbarClientWarning.style.display = storedClientWarning ? 'inline' : 'none';
  }

  // Make navbar client name clickable for override (double-click shows dropdown)
  if (navbarClientName) {
    navbarClientName.addEventListener('dblclick', async () => {
      // Fetch available clients
      try {
        const response = await fetch(`${baseURL}/api/clients`);
        if (!response.ok) {
          throw new Error(`HTTP ${response.status}`);
        }
        const data = await response.json();

        if (!data.success || !data.clients || data.clients.length === 0) {
          alert('No clients available');
          return;
        }

        // Create dropdown menu
        const currentName = navbarClientName.textContent.replace('✓ ', '');
        const dropdown = document.createElement('select');
        dropdown.className = 'client-select-dropdown';
        dropdown.style.cssText = `
          position: fixed;
          z-index: 10000;
          padding: 8px 12px;
          border: 1px solid var(--border-color, #ccc);
          border-radius: var(--radius-md, 4px);
          background-color: var(--bg-card, #fff);
          color: var(--text-primary, #000);
          font-size: var(--font-size-sm, 14px);
          font-family: inherit;
          cursor: pointer;
          box-shadow: 0 4px 12px rgba(0,0,0,0.2);
          min-width: 220px;
          max-height: 300px;
          outline: none;
        `;

        // Add default option
        const defaultOption = document.createElement('option');
        defaultOption.value = '';
        defaultOption.textContent = '-- Select Client --';
        dropdown.appendChild(defaultOption);

        // Add client options
        data.clients.forEach(client => {
          const option = document.createElement('option');
          option.value = client.name;
          option.textContent = client.name;
          if (client.name === currentName) {
            option.selected = true;
          }
          dropdown.appendChild(option);
        });

        // Position dropdown near the client name (use fixed positioning for scroll safety)
        const rect = navbarClientName.getBoundingClientRect();
        const scrollY = window.scrollY || window.pageYOffset;
        const scrollX = window.scrollX || window.pageXOffset;

        // Position below the client name, but adjust if near bottom of viewport
        let top = rect.bottom + scrollY + 4;
        const viewportHeight = window.innerHeight;
        if (rect.bottom + 300 > viewportHeight) {
          // Position above instead
          top = rect.top + scrollY - 300 - 4;
        }

        dropdown.style.top = `${Math.max(4, top)}px`;
        dropdown.style.left = `${rect.left + scrollX}px`;

        // Add to document
        document.body.appendChild(dropdown);
        dropdown.focus();

        // Handle selection
        const handleSelection = async () => {
          const selectedName = dropdown.value;
          if (!selectedName) {
            document.body.removeChild(dropdown);
            return;
          }

          // Remove dropdown
          document.body.removeChild(dropdown);

          // Send override request
          try {
            const overrideResponse = await fetch(`${baseURL}/api/client/override`, {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json',
              },
              body: JSON.stringify({ client_name: selectedName })
            });

            if (!overrideResponse.ok) {
              const errorText = await overrideResponse.text();
              throw new Error(`HTTP ${overrideResponse.status}: ${errorText}`);
            }

            const overrideData = await overrideResponse.json();
            if (overrideData.success) {
              // Update navbar client name
              navbarClientName.textContent = overrideData.client_name || selectedName;
              sessionStorage.setItem('clientName', overrideData.client_name || selectedName);
              // Show success message briefly before reload
              navbarClientName.textContent = '✓ ' + navbarClientName.textContent;
              setTimeout(() => {
                // Refresh page to update connection
                window.location.reload();
              }, 500);
            } else {
              alert('Failed to override client: ' + (overrideData.error || 'Unknown error'));
            }
          } catch (error) {
            errorHandler.logError('Error overriding client', error, 'ClientOverride');
            let errorMsg = 'Error overriding client';
            if (error.message) {
              errorMsg = error.message;
            }
            alert(errorMsg);
          }
        };

        // Handle selection on change
        dropdown.addEventListener('change', handleSelection);

        // Handle escape key to close
        const handleEscape = (e) => {
          if (e.key === 'Escape') {
            if (document.body.contains(dropdown)) {
              document.body.removeChild(dropdown);
            }
            document.removeEventListener('keydown', handleEscape);
            document.removeEventListener('click', handleClickOutside);
          }
        };

        // Handle click outside to close
        const handleClickOutside = (e) => {
          if (!dropdown.contains(e.target) && e.target !== navbarClientName) {
            if (document.body.contains(dropdown)) {
              document.body.removeChild(dropdown);
            }
            document.removeEventListener('keydown', handleEscape);
            document.removeEventListener('click', handleClickOutside);
          }
        };

        // Add event listeners
        document.addEventListener('keydown', handleEscape);
        setTimeout(() => {
          document.addEventListener('click', handleClickOutside);
        }, 100);
      } catch (error) {
        errorHandler.logError('Error fetching clients', error, 'ClientList');
        alert('Error fetching client list: ' + (error.message || 'Unknown error'));
      }
    });
  }

  // Fetch installation info and client info
  fetch(`${baseURL}/api/installation/info`)
    .then(response => response.json())
    .then(data => {
      console.log('[common.js] Installation info received:', data);
      // Update navbar elements (primary) and legacy canvas-header elements (fallback)
      const navbarClientNameEl = document.getElementById('navbarClientName');
      const navbarCanvasNameEl = document.getElementById('navbarCanvasName');
      const navbarWarningEl = document.getElementById('navbarClientWarning');
      const legacyClientNameEl = document.getElementById('clientNameDisplay');
      const legacyWarningEl = document.getElementById('clientWarning');

      // Update canvas name if available
      if (navbarCanvasNameEl) {
        if (data.canvas_name) {
          navbarCanvasNameEl.textContent = data.canvas_name;
          sessionStorage.setItem('canvasName', data.canvas_name);
        } else if (data.canvas_id) {
          navbarCanvasNameEl.textContent = data.canvas_id.substring(0, 8) + '...';
        } else {
          navbarCanvasNameEl.textContent = '...';
        }
      }

      // Check connection status first - if connected, we have a valid client
      if (data.connected && data.client_id) {
        // We're connected - show client name if available, otherwise show client ID or installation name
        const displayName = data.client_name || data.client_id || data.installation_name || 'Connected';
        if (navbarClientNameEl) {
          navbarClientNameEl.textContent = displayName;
          sessionStorage.setItem('clientName', displayName);
        }
        // Only show warning if we have client_id but no client_name (client might not have a name)
        if (navbarWarningEl) {
          if (data.client_id && !data.client_name) {
            navbarWarningEl.style.display = 'inline';
            navbarWarningEl.textContent = '(No name)';
          } else {
            navbarWarningEl.style.display = 'none';
          }
        }
        if (legacyClientNameEl) legacyClientNameEl.textContent = displayName;
        if (legacyWarningEl) {
          if (data.client_id && !data.client_name) {
            legacyWarningEl.style.display = 'inline';
            legacyWarningEl.textContent = '(Client has no name)';
          } else {
            legacyWarningEl.style.display = 'none';
          }
        }
        sessionStorage.setItem('clientWarning', data.client_id && !data.client_name ? 'true' : 'false');
      } else if (data.client_name) {
        // Client exists on server - show its name (but not connected yet)
        if (navbarClientNameEl) {
          navbarClientNameEl.textContent = data.client_name;
          sessionStorage.setItem('clientName', data.client_name);
        }
        if (navbarWarningEl) navbarWarningEl.style.display = 'none';
        if (legacyClientNameEl) legacyClientNameEl.textContent = data.client_name;
        if (legacyWarningEl) legacyWarningEl.style.display = 'none';
        sessionStorage.setItem('clientWarning', 'false');
      } else if (data.installation_name) {
        // Client not found or not connected - show installation name
        if (navbarClientNameEl) {
          navbarClientNameEl.textContent = data.installation_name;
          sessionStorage.setItem('clientName', data.installation_name);
        }
        if (navbarWarningEl) {
          navbarWarningEl.style.display = 'inline';
          navbarWarningEl.textContent = data.client_id ? '(No name)' : '(Not found)';
        }
        if (legacyClientNameEl) legacyClientNameEl.textContent = data.installation_name;
        if (legacyWarningEl) {
          legacyWarningEl.style.display = 'inline';
          legacyWarningEl.textContent = data.client_id ? '(Client has no name)' : '(Client not found on server)';
        }
        sessionStorage.setItem('clientWarning', 'true');
      } else {
        // No info available
        if (navbarClientNameEl) {
          navbarClientNameEl.textContent = 'Unknown';
          sessionStorage.setItem('clientName', 'Unknown');
        }
        if (navbarWarningEl) navbarWarningEl.style.display = 'none';
        if (legacyClientNameEl) legacyClientNameEl.textContent = 'Unknown';
        if (legacyWarningEl) legacyWarningEl.style.display = 'none';
        sessionStorage.setItem('clientWarning', 'false');
      }
    })
    .catch(error => {
      errorHandler.logError('Failed to fetch installation info', error, 'NavbarTracking');
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

  // Check if this is a page refresh on the home page
  const isHomePage = window.location.pathname === '/' || window.location.pathname === '/index.html' || window.location.pathname === '/main.html';
  const wasRefreshed = performance.navigation.type === performance.navigation.TYPE_RELOAD ||
                       (performance.getEntriesByType && performance.getEntriesByType('navigation')[0]?.type === 'reload');

  // If home page was refreshed, force disconnect and restart
  if (isHomePage && wasRefreshed) {
    console.log('[common.js] Home page refresh detected - forcing disconnect and restart');

    // Disconnect existing connection
    workspaceClient.disconnect();

    // Call restart API to restart canvas service
    fetch(`${baseURL}/api/canvas/restart`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
    })
      .then(response => response.json())
      .then(data => {
        if (data.success) {
          console.log('[common.js] Canvas service restarted successfully');
          // Wait a moment for service to restart, then connect
          setTimeout(() => {
            workspaceClient.connect(baseURL);
          }, 1000);
        } else {
          console.error('[common.js] Failed to restart canvas service:', data.error);
          // Still try to connect even if restart failed
          workspaceClient.connect(baseURL);
        }
      })
      .catch(error => {
        console.error('[common.js] Error calling restart API:', error);
        // Still try to connect even if restart failed
        workspaceClient.connect(baseURL);
      });
  } else {
    // Normal page load - just connect
    workspaceClient.connect(baseURL);
  }

  workspaceClient.on('connected', () => {
    // Update status to connected - no need to fetch installation info again
    // The canvas_update events will provide client_name and canvas_name
    // and we already fetched installation info once on page load
    updateStatus('connected', 'Connected');
  });

  workspaceClient.on('disconnected', () => {
    updateStatus('disconnected', 'Disconnected');
  });

  workspaceClient.on('reconnecting', (attempts) => {
    updateStatus('connecting', `Reconnecting (${attempts})...`);
  });

  workspaceClient.on('canvas_update', (data) => {
    console.log('[common.js] canvas_update event received:', data);

    // Update client name if provided
    if (data.client_name) {
      const navbarClientName = document.getElementById('navbarClientName');
      if (navbarClientName) {
        navbarClientName.textContent = data.client_name;
        sessionStorage.setItem('clientName', data.client_name);
      }
    }

    // Update canvas name if provided
    if (data.canvas_name) {
      const navbarCanvasName = document.getElementById('navbarCanvasName');
      if (navbarCanvasName) {
        navbarCanvasName.textContent = data.canvas_name;
        sessionStorage.setItem('canvasName', data.canvas_name);
      }
    }
  });

  function updateStatus(status, text) {
    // Update navbar status (primary)
    const navbarIndicator = document.getElementById('navbarStatusIndicator');
    const navbarStatusText = document.getElementById('navbarStatusText');

    // Update legacy status elements (fallback for pages that still have them)
    const legacyIndicator = document.getElementById('statusIndicator');
    const legacyStatusText = document.getElementById('statusText');

    if (navbarIndicator) navbarIndicator.className = `navbar-status-indicator ${status}`;
    if (navbarStatusText) navbarStatusText.textContent = text;

    if (legacyIndicator) legacyIndicator.className = `status-indicator ${status}`;
    if (legacyStatusText) legacyStatusText.textContent = text;

    // Store status in sessionStorage for persistence
    sessionStorage.setItem('connectionStatus', status);
    sessionStorage.setItem('connectionStatusText', text);
  }

  // Restore status from sessionStorage on page load
  const storedStatus = sessionStorage.getItem('connectionStatus');
  const storedStatusText = sessionStorage.getItem('connectionStatusText');
  if (storedStatus && storedStatusText) {
    updateStatus(storedStatus, storedStatusText);
  }
}

