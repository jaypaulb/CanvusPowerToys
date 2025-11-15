/**
 * RCU Page JavaScript
 * Handles RCU configuration and status
 */

document.addEventListener('DOMContentLoaded', () => {
  loadRCUConfig();
  loadRCUStatus();
  initRCUForm();
});

/**
 * Load RCU configuration
 */
async function loadRCUConfig() {
  const baseURL = window.location.origin;

  try {
    const response = await fetch(`${baseURL}/api/rcu/config`);
    if (response.ok) {
      const config = await response.json();
      document.getElementById('rcuEnabled').checked = config.enabled || false;
      document.getElementById('rcuPort').value = config.port || '';
      document.getElementById('rcuTimeout').value = config.timeout || '';
    }
  } catch (error) {
    // Only log in development
    if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
      console.error('Error loading RCU config:', error);
    }
  }
}

/**
 * Load RCU status
 */
async function loadRCUStatus() {
  const baseURL = window.location.origin;
  const statusDiv = document.getElementById('rcuStatus');

  try {
    const response = await fetch(`${baseURL}/api/rcu/status`);
    if (response.ok) {
      const status = await response.json();
      renderRCUStatus(status);
    } else {
      statusDiv.innerHTML = '<p class="text-muted">Unable to load RCU status</p>';
    }
  } catch (error) {
    // Only log in development
    if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
      console.error('Error loading RCU status:', error);
    }
    statusDiv.innerHTML = '<p class="text-muted">Error loading RCU status</p>';
  }
}

/**
 * Render RCU status
 */
function renderRCUStatus(status) {
  const statusDiv = document.getElementById('rcuStatus');

  statusDiv.innerHTML = `
    <div class="form-group">
      <label class="input-label">Status</label>
      <p>${status.connected ? 'Connected' : 'Disconnected'}</p>
    </div>
    <div class="form-group">
      <label class="input-label">Last Update</label>
      <p>${status.last_update ? new Date(status.last_update).toLocaleString() : 'Never'}</p>
    </div>
  `;
}

/**
 * Initialize RCU form
 */
function initRCUForm() {
  const form = document.getElementById('rcuForm');
  const testBtn = document.getElementById('testRcuBtn');

  form.addEventListener('submit', handleSaveConfig);
  testBtn.addEventListener('click', handleTestConnection);
}

/**
 * Handle save configuration
 */
async function handleSaveConfig(e) {
  e.preventDefault();

  const baseURL = window.location.origin;
  const formData = {
    enabled: document.getElementById('rcuEnabled').checked,
    port: parseInt(document.getElementById('rcuPort').value) || 8080,
    timeout: parseInt(document.getElementById('rcuTimeout').value) || 30
  };

  try {
    const response = await fetch(`${baseURL}/api/rcu/config`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(formData)
    });

    if (response.ok) {
      alert('Configuration saved successfully');
      loadRCUStatus();
    } else {
      alert('Failed to save configuration');
    }
  } catch (error) {
    // Only log in development
    if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
      console.error('Error saving RCU config:', error);
    }
    alert('Error saving configuration');
  }
}

/**
 * Handle test connection
 */
async function handleTestConnection() {
  const baseURL = window.location.origin;

  try {
    const response = await fetch(`${baseURL}/api/rcu/test`, {
      method: 'POST'
    });

    if (response.ok) {
      const result = await response.json();
      alert(result.success ? 'Connection test successful' : 'Connection test failed');
    } else {
      alert('Connection test failed');
    }
  } catch (error) {
    // Only log in development
    if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
      console.error('Error testing RCU connection:', error);
    }
    alert('Error testing connection');
  }
}

